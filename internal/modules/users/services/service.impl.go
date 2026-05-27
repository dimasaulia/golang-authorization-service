package services

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/open-suite/authorization/internal/entities"
	teamMemberServices "github.com/open-suite/authorization/internal/modules/teammembers/services"
	userRoleServices "github.com/open-suite/authorization/internal/modules/userroles/services"
	"github.com/open-suite/authorization/internal/modules/users/dto"
	"github.com/open-suite/authorization/internal/modules/users/repositories"
	"github.com/open-suite/authorization/internal/platform/config"
	"github.com/open-suite/authorization/internal/platform/freeipa"
	"github.com/open-suite/authorization/internal/platform/keycloak"
	"github.com/open-suite/authorization/internal/platform/logger"
	"github.com/open-suite/authorization/internal/platform/mailer"
	"github.com/open-suite/authorization/internal/shared"
	"golang.org/x/crypto/bcrypt"
)

const (
	userTypeInternal = "internal"
	userTypeExternal = "external"
	userTypeService  = "service"

	userStatusInvited = "invited"
	userStatusActive  = "active"

	purposeEmailVerification = "email_verification"
	purposePasswordSetup     = "password_setup"

	keycloakProvider = "keycloak"
	freeIpaProvider  = "freeipa"
)

var ErrInvalidRequest = errors.New("invalid user request")

type UserServiceImpl struct {
	UserRepository    repositories.UserRepository
	UserRoleService   userRoleServices.UserRoleService
	TeamMemberService teamMemberServices.TeamMemberService
	cfg               config.Config
	keycloak          keycloak.Client
	freeipa           freeipa.Client
	mailer            mailer.Mailer
	log               *logger.LayerLogger
}

func NewUserService(repository repositories.UserRepository, userRoleService userRoleServices.UserRoleService, teamMemberService teamMemberServices.TeamMemberService, cfg config.Config, keycloakClient keycloak.Client, freeIPAClient freeipa.Client, mailer mailer.Mailer, appLogger *logger.Logger) UserService {
	return &UserServiceImpl{
		UserRepository:    repository,
		UserRoleService:   userRoleService,
		TeamMemberService: teamMemberService,
		cfg:               cfg,
		keycloak:          keycloakClient,
		freeipa:           freeIPAClient,
		mailer:            mailer,
		log:               appLogger.Layer("service.users"),
	}
}

func (s *UserServiceImpl) Find(ctx context.Context, params shared.ListParams) ([]entities.User, error) {
	end := s.log.Start(ctx, "Find")
	items, err := s.UserRepository.Find(ctx, params)
	end(err, "count", len(items))
	return items, err
}

func (s *UserServiceImpl) FindByID(ctx context.Context, id int64) (*entities.User, error) {
	end := s.log.Start(ctx, "FindByID", "id", id)
	item, err := s.UserRepository.FindByID(ctx, id)
	end(err)
	return item, err
}

func (s *UserServiceImpl) Create(ctx context.Context, request dto.CreateUserRequest) (*dto.UserResponse, error) {
	end := s.log.Start(ctx, "Create")
	request = normalizeCreateRequest(request)
	if err := validateCreateRequest(request); err != nil {
		end(err)
		return nil, err
	}

	var passwordHash string
	if strings.TrimSpace(request.Password) != "" {
		hash, err := hashPassword(request.Password)
		if err != nil {
			end(err)
			return nil, err
		}
		passwordHash = hash
	}

	var token string
	input := repositories.CreateUserInput{
		User: entities.User{
			OrganizationId: request.OrganizationId,
			Username:       request.Username,
			Email:          request.Email,
			DisplayName:    request.DisplayName,
			Type:           request.Type,
			Status:         request.Status,
		},
		PasswordHash:       passwordHash,
		MustChangePassword: valueOrDefault(request.MustChangePassword, strings.TrimSpace(request.Password) != ""),
	}

	if request.SendInvitation || request.Status == userStatusInvited {
		token = newToken()
		input.User.Status = userStatusInvited
		input.VerificationCode = &repositories.CreateVerificationCodeInput{
			Purpose:   purposePasswordSetup,
			CodeHash:  hashCode(token),
			ExpiresAt: time.Now().Add(72 * time.Hour),
		}
	}

	item, err := s.UserRepository.Create(ctx, input)
	if err != nil {
		end(err)
		return nil, err
	}

	if len(request.RoleIds) > 0 {
		organizationID := request.OrganizationId
		if _, err := s.UserRoleService.AssignRolesToUser(ctx, item.ID, request.RoleIds, &organizationID, nil); err != nil {
			s.deleteCreatedUser(ctx, item.ID)
			end(err)
			return nil, err
		}
	}
	if len(request.TeamIds) > 0 {
		if _, err := s.TeamMemberService.AssignTeamsToUser(ctx, item.ID, request.TeamIds); err != nil {
			s.deleteCreatedUser(ctx, item.ID)
			end(err)
			return nil, err
		}
	}

	freeIPAUID, err := s.createFreeIPAUser(ctx, item, request.Password)
	if err != nil {
		s.deleteCreatedUser(ctx, item.ID)
		end(err)
		return nil, err
	}

	provisionedTo := []string{freeIpaProvider}
	if s.keycloak.Enabled() {
		_, err = s.createKeycloakUser(ctx, item, request.Password, input.MustChangePassword, item.Status == userStatusActive)
		if err != nil {
			s.deleteFreeIPAUser(ctx, freeIPAUID)
			s.deleteCreatedUser(ctx, item.ID)
			end(err)
			return nil, err
		}
		provisionedTo = append(provisionedTo, keycloakProvider)
	}

	if token != "" {
		if err := s.sendPasswordSetupEmail(ctx, item.Email, token); err != nil {
			s.log.Warn(ctx, "send_password_setup_email.failed", "user_id", item.ID, "error", err.Error())
		}
	}

	end(nil)
	return toUserResponse(item, input.MustChangePassword, provisionedTo, request.RoleIds, request.TeamIds), nil
}

func (s *UserServiceImpl) Signup(ctx context.Context, request dto.SignupUserRequest) (*dto.UserResponse, error) {
	end := s.log.Start(ctx, "Signup")
	request.Username = strings.TrimSpace(request.Username)
	request.Email = strings.ToLower(strings.TrimSpace(request.Email))
	request.DisplayName = strings.TrimSpace(request.DisplayName)
	if request.OrganizationId == 0 || request.Username == "" || request.Email == "" || request.DisplayName == "" || strings.TrimSpace(request.Password) == "" {
		end(ErrInvalidRequest)
		return nil, ErrInvalidRequest
	}

	passwordHash, err := hashPassword(request.Password)
	if err != nil {
		end(err)
		return nil, err
	}

	token := newToken()
	item, err := s.UserRepository.Create(ctx, repositories.CreateUserInput{
		User: entities.User{
			OrganizationId: request.OrganizationId,
			Username:       request.Username,
			Email:          request.Email,
			DisplayName:    request.DisplayName,
			Type:           userTypeExternal,
			Status:         userStatusActive,
		},
		PasswordHash:       passwordHash,
		MustChangePassword: false,
		VerificationCode: &repositories.CreateVerificationCodeInput{
			Purpose:   purposeEmailVerification,
			CodeHash:  hashCode(token),
			ExpiresAt: time.Now().Add(24 * time.Hour),
		},
	})
	if err != nil {
		end(err)
		return nil, err
	}

	createdInKeycloak := s.keycloak.Enabled()
	if createdInKeycloak {
		if _, err := s.createKeycloakUser(ctx, item, request.Password, false, true); err != nil {
			s.deleteCreatedUser(ctx, item.ID)
			end(err)
			return nil, err
		}
	}

	if err := s.sendEmailVerification(ctx, item.Email, token); err != nil {
		s.log.Warn(ctx, "send_email_verification.failed", "user_id", item.ID, "error", err.Error())
	}

	end(nil)
	provisionedTo := []string{}
	if createdInKeycloak {
		provisionedTo = append(provisionedTo, keycloakProvider)
	}
	return toUserResponse(item, false, provisionedTo, nil, nil), nil
}

func (s *UserServiceImpl) SignupWithGoogle(ctx context.Context, request dto.GoogleSignupRequest) (*dto.UserResponse, error) {
	end := s.log.Start(ctx, "SignupWithGoogle")
	request.Username = strings.TrimSpace(request.Username)
	request.Email = strings.ToLower(strings.TrimSpace(request.Email))
	request.DisplayName = strings.TrimSpace(request.DisplayName)
	request.ProviderUserId = strings.TrimSpace(request.ProviderUserId)
	if request.OrganizationId == 0 || request.ProviderUserId == "" || request.Username == "" || request.Email == "" || request.DisplayName == "" {
		end(ErrInvalidRequest)
		return nil, ErrInvalidRequest
	}

	username := request.Username
	email := request.Email
	item, err := s.UserRepository.Create(ctx, repositories.CreateUserInput{
		User: entities.User{
			OrganizationId: request.OrganizationId,
			Username:       request.Username,
			Email:          request.Email,
			DisplayName:    request.DisplayName,
			Type:           userTypeExternal,
			Status:         userStatusActive,
		},
		Identity: &entities.UserIdentity{
			Provider:       "google",
			ProviderUserId: request.ProviderUserId,
			Username:       &username,
			Email:          &email,
			IsPrimary:      true,
		},
	})
	if err != nil {
		end(err)
		return nil, err
	}

	end(nil)
	return toUserResponse(item, false, nil, nil, nil), nil
}

func (s *UserServiceImpl) VerifyEmail(ctx context.Context, request dto.VerifyEmailRequest) error {
	end := s.log.Start(ctx, "VerifyEmail")
	code := strings.TrimSpace(request.Code)
	if code == "" {
		end(ErrInvalidRequest)
		return ErrInvalidRequest
	}

	item, err := s.UserRepository.FindVerificationCode(ctx, purposeEmailVerification, hashCode(code))
	if err != nil {
		end(err)
		return err
	}

	if err := s.UserRepository.MarkEmailVerified(ctx, item.UserID); err != nil {
		end(err)
		return err
	}

	err = s.UserRepository.UseVerificationCode(ctx, item.ID)
	end(err)
	return err
}

func (s *UserServiceImpl) Update(ctx context.Context, id int64, request dto.UpdateUserRequest) (*entities.User, error) {
	end := s.log.Start(ctx, "Update", "id", id)

	data := map[string]any{}
	if request.OrganizationId != nil {
		data["organization_id"] = *request.OrganizationId
	}
	if request.Username != nil {
		data["username"] = *request.Username
	}
	if request.Email != nil {
		data["email"] = *request.Email
	}
	if request.DisplayName != nil {
		data["display_name"] = *request.DisplayName
	}
	if request.Type != nil {
		data["type"] = *request.Type
	}
	if request.Status != nil {
		data["status"] = *request.Status
	}

	item, err := s.UserRepository.Update(ctx, id, data)
	end(err)
	return item, err
}

func (s *UserServiceImpl) Delete(ctx context.Context, id int64) error {
	end := s.log.Start(ctx, "Delete", "id", id)
	err := s.UserRepository.Delete(ctx, id)
	end(err)
	return err
}

func normalizeCreateRequest(request dto.CreateUserRequest) dto.CreateUserRequest {
	request.Username = strings.TrimSpace(request.Username)
	request.Email = strings.ToLower(strings.TrimSpace(request.Email))
	request.DisplayName = strings.TrimSpace(request.DisplayName)
	request.Type = strings.TrimSpace(request.Type)
	request.Status = strings.TrimSpace(request.Status)
	if request.Type == "" {
		request.Type = userTypeInternal
	}
	if request.Status == "" {
		request.Status = userStatusActive
	}
	return request
}

func validateCreateRequest(request dto.CreateUserRequest) error {
	if request.OrganizationId == 0 || request.Username == "" || request.Email == "" || request.DisplayName == "" {
		return ErrInvalidRequest
	}
	if !isAllowed(request.Type, userTypeInternal, userTypeExternal, userTypeService) {
		return ErrInvalidRequest
	}
	if !isAllowed(request.Status, userStatusInvited, userStatusActive, "suspended", "disabled") {
		return ErrInvalidRequest
	}
	if request.Status == userStatusInvited && strings.TrimSpace(request.Password) != "" {
		return ErrInvalidRequest
	}
	return nil
}

func isAllowed(value string, allowed ...string) bool {
	for _, item := range allowed {
		if value == item {
			return true
		}
	}
	return false
}

func hashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func newToken() string {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(bytes)
}

func hashCode(code string) string {
	sum := sha256.Sum256([]byte(code))
	return hex.EncodeToString(sum[:])
}

func valueOrDefault(value *bool, fallback bool) bool {
	if value == nil {
		return fallback
	}
	return *value
}

func (s *UserServiceImpl) sendEmailVerification(ctx context.Context, email string, token string) error {
	link := fmt.Sprintf("%s/api/v1/users/verify-email?code=%s", s.cfg.App.PublicURL, token)
	return s.mailer.Send(ctx, email, "Verify your email", "Please verify your email using this link:\n\n"+link)
}

func (s *UserServiceImpl) sendPasswordSetupEmail(ctx context.Context, email string, token string) error {
	link := fmt.Sprintf("%s/password/setup?code=%s", s.cfg.App.PublicURL, token)
	return s.mailer.Send(ctx, email, "Set up your password", "Please set your password using this link:\n\n"+link)
}

func (s *UserServiceImpl) deleteCreatedUser(ctx context.Context, userID int64) {
	if err := s.UserRepository.Delete(ctx, userID); err != nil {
		s.log.Warn(ctx, "create.rollback_failed", "user_id", userID, "error", err.Error())
	}
}

func (s *UserServiceImpl) createKeycloakUser(ctx context.Context, user *entities.User, password string, temporaryPassword bool, enabled bool) (string, error) {
	keycloakUserID, err := s.keycloak.CreateUser(ctx, keycloak.CreateUserInput{
		Username:          user.Username,
		Email:             user.Email,
		DisplayName:       user.DisplayName,
		Enabled:           enabled,
		EmailVerified:     user.EmailVerifiedAt != nil,
		Password:          password,
		TemporaryPassword: temporaryPassword,
	})
	if err != nil {
		return "", err
	}

	username := user.Username
	email := user.Email
	if err := s.UserRepository.LinkIdentity(ctx, entities.UserIdentity{
		UserId:         user.ID,
		Provider:       keycloakProvider,
		ProviderUserId: keycloakUserID,
		Username:       &username,
		Email:          &email,
		IsPrimary:      true,
	}); err != nil {
		_ = s.keycloak.DeleteUser(ctx, keycloakUserID)
		return "", err
	}

	return keycloakUserID, nil
}

func (s *UserServiceImpl) createFreeIPAUser(ctx context.Context, user *entities.User, password string) (string, error) {
	uid, err := s.freeipa.CreateUser(ctx, freeipa.CreateUserInput{
		Username:    user.Username,
		Email:       user.Email,
		DisplayName: user.DisplayName,
		Password:    password,
	})
	if err != nil {
		return "", err
	}

	username := user.Username
	email := user.Email
	if err := s.UserRepository.LinkIdentity(ctx, entities.UserIdentity{
		UserId:         user.ID,
		Provider:       freeIpaProvider,
		ProviderUserId: uid,
		Username:       &username,
		Email:          &email,
		IsPrimary:      false,
	}); err != nil {
		_ = s.freeipa.DeleteUser(ctx, uid)
		return "", err
	}

	return uid, nil
}

func (s *UserServiceImpl) deleteFreeIPAUser(ctx context.Context, uid string) {
	if uid == "" {
		return
	}
	if err := s.freeipa.DeleteUser(ctx, uid); err != nil {
		s.log.Warn(ctx, "freeipa.rollback_failed", "uid", uid, "error", err.Error())
	}
}

func toUserResponse(user *entities.User, mustChangePassword bool, provisionedTo []string, roleIDs []int64, teamIDs []int64) *dto.UserResponse {
	return &dto.UserResponse{
		ID:                 user.ID,
		OrganizationId:     user.OrganizationId,
		Username:           user.Username,
		Email:              user.Email,
		DisplayName:        user.DisplayName,
		Type:               user.Type,
		Status:             user.Status,
		MustChangePassword: mustChangePassword,
		ProvisionedTo:      provisionedTo,
		RoleIds:            roleIDs,
		TeamIds:            teamIDs,
	}
}

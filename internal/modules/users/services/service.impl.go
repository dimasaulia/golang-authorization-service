package services

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"html/template"
	"net/url"
	"strings"
	"time"

	"github.com/open-suite/authorization/internal/entities"
	teamMemberServices "github.com/open-suite/authorization/internal/modules/teammembers/services"
	userRoleServices "github.com/open-suite/authorization/internal/modules/userroles/services"
	"github.com/open-suite/authorization/internal/modules/users/dto"
	"github.com/open-suite/authorization/internal/modules/users/repositories"
	"github.com/open-suite/authorization/internal/platform/config"
	"github.com/open-suite/authorization/internal/platform/freeipa"
	"github.com/open-suite/authorization/internal/platform/i18n"
	"github.com/open-suite/authorization/internal/platform/keycloak"
	"github.com/open-suite/authorization/internal/platform/logger"
	"github.com/open-suite/authorization/internal/platform/mailer"
	"github.com/open-suite/authorization/internal/platform/redis"
	"github.com/open-suite/authorization/internal/shared"
	"github.com/open-suite/authorization/internal/shared/requestctx"
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
	redis             *redis.Redis
	translator        *i18n.Translator
	log               *logger.LayerLogger
}

func NewUserService(repository repositories.UserRepository, userRoleService userRoleServices.UserRoleService, teamMemberService teamMemberServices.TeamMemberService, cfg config.Config, keycloakClient keycloak.Client, freeIPAClient freeipa.Client, mailer mailer.Mailer, redisClient *redis.Redis, translator *i18n.Translator, appLogger *logger.Logger) UserService {
	return &UserServiceImpl{
		UserRepository:    repository,
		UserRoleService:   userRoleService,
		TeamMemberService: teamMemberService,
		cfg:               cfg,
		keycloak:          keycloakClient,
		freeipa:           freeIPAClient,
		mailer:            mailer,
		redis:             redisClient,
		translator:        translator,
		log:               appLogger.Layer("service.users"),
	}
}

func (s *UserServiceImpl) Find(ctx context.Context, params shared.ListParams) ([]entities.User, error) {
	end := s.log.Start(ctx, "Find")
	items, err := s.UserRepository.Find(ctx, params)
	end(err, "count", len(items))
	return items, err
}

func (s *UserServiceImpl) FindByID(ctx context.Context, id int64) (*dto.UserResponse, error) {
	end := s.log.Start(ctx, "FindByID", "id", id)
	item, err := s.UserRepository.FindByID(ctx, id)
	if err != nil {
		end(err)
		return nil, err
	}

	response, err := s.toUserResponseWithAssignments(ctx, item, false, nil)
	end(err)
	return response, err
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

	if _, err := s.createFreeIPAUser(ctx, item, request.Password); err != nil {
		s.deleteCreatedUser(ctx, item.ID)
		end(err)
		return nil, err
	}

	provisionedTo := []string{freeIpaProvider}

	if token != "" {
		if err := s.sendPasswordSetupEmail(ctx, item.Email, item.DisplayName, token, request.SetupPasswordURL); err != nil {
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

	if err := s.sendEmailVerification(ctx, item.Email, item.DisplayName, token); err != nil {
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

func (s *UserServiceImpl) SetupPassword(ctx context.Context, request dto.SetupPasswordRequest) error {
	end := s.log.Start(ctx, "SetupPassword")
	println("SETUP PASSWORD => start")
	code := strings.TrimSpace(request.Code)
	password := strings.TrimSpace(request.Password)
	if code == "" || password == "" {
		println("SETUP PASSWORD => invalid request")
		end(ErrInvalidRequest)
		return ErrInvalidRequest
	}

	codeHash := hashCode(code)
	println("SETUP PASSWORD => payload code:", code)
	println("SETUP PASSWORD => computed code_hash:", codeHash)

	verificationCode, err := s.UserRepository.FindVerificationCode(ctx, purposePasswordSetup, codeHash)
	if err != nil {
		println("SETUP PASSWORD => FindVerificationCode error:", err.Error())
		s.log.Error(ctx, "setup_password.verification_code_lookup_failed", err, "purpose", purposePasswordSetup)
		end(err)
		return err
	}
	println("SETUP PASSWORD => verification_code_id:", verificationCode.ID)
	println("SETUP PASSWORD => verification_code_user_id:", verificationCode.UserID)
	println("SETUP PASSWORD => verification_code_purpose:", verificationCode.Purpose)
	s.log.Debug(ctx, "setup_password.verification_code_found", "verification_code_id", verificationCode.ID, "user_id", verificationCode.UserID, "purpose", verificationCode.Purpose, "expires_at", verificationCode.ExpiresAt)

	user, err := s.UserRepository.FindByID(ctx, verificationCode.UserID)
	if err != nil {
		println("SETUP PASSWORD => FindByID error for user_id:", verificationCode.UserID, "error:", err.Error())
		end(err)
		return err
	}
	println("SETUP PASSWORD => user_found:", user.ID)

	passwordHash, err := hashPassword(password)
	if err != nil {
		println("SETUP PASSWORD => hashPassword error:", err.Error())
		end(err)
		return err
	}

	if err := s.syncPasswordSetupProviders(ctx, user, password); err != nil {
		println("SETUP PASSWORD => syncPasswordSetupProviders error:", err.Error())
		end(err)
		return err
	}

	if err := s.UserRepository.UpdateCredential(ctx, user.ID, passwordHash, false); err != nil {
		println("SETUP PASSWORD => UpdateCredential error:", err.Error())
		end(err)
		return err
	}
	if _, err := s.UserRepository.Update(ctx, user.ID, map[string]any{"status": userStatusActive}); err != nil {
		println("SETUP PASSWORD => Update user status error:", err.Error())
		end(err)
		return err
	}

	if err := s.UserRepository.UseVerificationCode(ctx, verificationCode.ID); err != nil {
		println("SETUP PASSWORD => UseVerificationCode error:", err.Error())
		end(err)
		return err
	}

	println("SETUP PASSWORD => success user_id:", user.ID)
	end(nil, "user_id", user.ID)
	return nil
}

func (s *UserServiceImpl) ResendVerificationEmail(ctx context.Context, request dto.ResendVerificationEmailRequest) error {
	end := s.log.Start(ctx, "ResendVerificationEmail")
	email := strings.ToLower(strings.TrimSpace(request.Email))
	if email == "" {
		end(ErrInvalidRequest)
		return ErrInvalidRequest
	}

	user, err := s.UserRepository.FindByEmail(ctx, email)
	if err != nil {
		end(err)
		return err
	}
	if user.EmailVerifiedAt != nil {
		end(ErrInvalidRequest)
		return ErrInvalidRequest
	}

	token := newToken()
	if err := s.UserRepository.CreateVerificationCode(ctx, user.ID, repositories.CreateVerificationCodeInput{
		Purpose:   purposeEmailVerification,
		CodeHash:  hashCode(token),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}); err != nil {
		end(err)
		return err
	}

	err = s.sendEmailVerification(ctx, user.Email, user.DisplayName, token)
	end(err)
	return err
}

func (s *UserServiceImpl) ResendInvitation(ctx context.Context, id int64, request dto.ResendInvitationRequest) error {
	end := s.log.Start(ctx, "ResendInvitation", "id", id)
	if id == 0 {
		end(ErrInvalidRequest)
		return ErrInvalidRequest
	}

	user, err := s.UserRepository.FindByID(ctx, id)
	if err != nil {
		end(err)
		return err
	}
	if user.Status != userStatusInvited {
		end(ErrInvalidRequest)
		return ErrInvalidRequest
	}

	token := newToken()
	if err := s.UserRepository.CreateVerificationCode(ctx, user.ID, repositories.CreateVerificationCodeInput{
		Purpose:   purposePasswordSetup,
		CodeHash:  hashCode(token),
		ExpiresAt: time.Now().Add(72 * time.Hour),
	}); err != nil {
		end(err)
		return err
	}

	err = s.sendPasswordSetupEmail(ctx, user.Email, user.DisplayName, token, request.SetupPasswordURL)
	end(err)
	return err
}

func (s *UserServiceImpl) Update(ctx context.Context, id int64, request dto.UpdateUserRequest) (*dto.UserResponse, error) {
	end := s.log.Start(ctx, "Update", "id", id)

	user, err := s.UserRepository.FindByID(ctx, id)
	if err != nil {
		end(err)
		return nil, err
	}

	nextUsername := user.Username
	nextEmail := user.Email
	nextDisplayName := user.DisplayName
	data := map[string]any{}
	if request.OrganizationId != nil {
		data["organization_id"] = *request.OrganizationId
	}
	if request.Username != nil {
		nextUsername = strings.TrimSpace(*request.Username)
		if nextUsername == "" {
			end(ErrInvalidRequest)
			return nil, ErrInvalidRequest
		}
		data["username"] = nextUsername
	}
	if request.Email != nil {
		nextEmail = strings.ToLower(strings.TrimSpace(*request.Email))
		if nextEmail == "" {
			end(ErrInvalidRequest)
			return nil, ErrInvalidRequest
		}
		data["email"] = nextEmail
	}
	if request.DisplayName != nil {
		nextDisplayName = strings.TrimSpace(*request.DisplayName)
		if nextDisplayName == "" {
			end(ErrInvalidRequest)
			return nil, ErrInvalidRequest
		}
		data["display_name"] = nextDisplayName
	}
	if request.Type != nil {
		data["type"] = *request.Type
	}
	if request.Status != nil {
		data["status"] = *request.Status
	}

	password := ""
	passwordHash := ""
	if request.Password != nil {
		password = strings.TrimSpace(*request.Password)
		if password == "" {
			end(ErrInvalidRequest)
			return nil, ErrInvalidRequest
		}
		passwordHash, err = hashPassword(password)
		if err != nil {
			end(err)
			return nil, err
		}
	}

	providerSyncNeeded := request.Username != nil || request.Email != nil || request.DisplayName != nil || password != ""
	if providerSyncNeeded {
		identities, err := s.UserRepository.FindIdentitiesByUserID(ctx, id)
		if err != nil {
			end(err)
			return nil, err
		}
		updatedProviderIDs := map[string]string{}
		for _, identity := range identities {
			switch identity.Provider {
			case keycloakProvider:
				if s.keycloak.Enabled() {
					if err := s.keycloak.UpdateUser(ctx, identity.ProviderUserId, keycloak.UpdateUserInput{
						Username:    nextUsername,
						Email:       nextEmail,
						DisplayName: nextDisplayName,
						Password:    password,
					}); err != nil {
						end(err, "provider", keycloakProvider)
						return nil, err
					}
				}
			case freeIpaProvider:
				if s.freeipa.Enabled() {
					uid, err := s.freeipa.UpdateUser(ctx, identity.ProviderUserId, freeipa.UpdateUserInput{
						Username:    nextUsername,
						Email:       nextEmail,
						DisplayName: nextDisplayName,
						Password:    password,
					})
					if err != nil {
						end(err, "provider", freeIpaProvider)
						return nil, err
					}
					if uid != "" {
						updatedProviderIDs[freeIpaProvider] = uid
					}
				}
			}
		}
		for _, identity := range identities {
			providerUserID := identity.ProviderUserId
			if value := updatedProviderIDs[identity.Provider]; value != "" {
				providerUserID = value
			}
			if err := s.UserRepository.UpdateIdentityProfile(ctx, id, identity.Provider, providerUserID, nextUsername, nextEmail); err != nil {
				end(err, "provider", identity.Provider)
				return nil, err
			}
		}
	}

	if passwordHash != "" {
		mustChangePassword := valueOrDefault(request.MustChangePassword, false)
		if err := s.UserRepository.UpdateCredential(ctx, id, passwordHash, mustChangePassword); err != nil {
			end(err)
			return nil, err
		}
	}

	item, err := s.UserRepository.Update(ctx, id, data)
	if err != nil {
		end(err)
		return nil, err
	}

	if request.RoleIds != nil {
		organizationID := item.OrganizationId
		if _, err := s.UserRoleService.ReplaceRolesForUser(ctx, item.ID, *request.RoleIds, &organizationID, nil); err != nil {
			end(err)
			return nil, err
		}
	}
	if request.TeamIds != nil {
		if _, err := s.TeamMemberService.ReplaceTeamsForUser(ctx, item.ID, *request.TeamIds); err != nil {
			end(err)
			return nil, err
		}
	}

	if err := s.resetUserAccessCache(ctx, id); err != nil {
		end(err)
		return nil, err
	}

	response, err := s.toUserResponseWithAssignments(ctx, item, valueOrDefault(request.MustChangePassword, false), nil)
	end(err)
	return response, err
}

func (s *UserServiceImpl) syncPasswordSetupProviders(ctx context.Context, user *entities.User, password string) error {
	identities, err := s.UserRepository.FindIdentitiesByUserID(ctx, user.ID)
	if err != nil {
		return err
	}

	identityByProvider := map[string]entities.UserIdentity{}
	for _, identity := range identities {
		identityByProvider[identity.Provider] = identity
	}

	if s.freeipa.Enabled() {
		if identity, exists := identityByProvider[freeIpaProvider]; exists {
			uid, err := s.freeipa.UpdateUser(ctx, identity.ProviderUserId, freeipa.UpdateUserInput{
				Username:    user.Username,
				Email:       user.Email,
				DisplayName: user.DisplayName,
				Password:    password,
			})
			if err != nil {
				return err
			}
			if uid != "" {
				if err := s.UserRepository.UpdateIdentityProfile(ctx, user.ID, freeIpaProvider, uid, user.Username, user.Email); err != nil {
					return err
				}
			}
		} else {
			if _, err := s.createFreeIPAUser(ctx, user, password); err != nil {
				return err
			}
		}
	}

	if s.keycloak.Enabled() {
		if identity, exists := identityByProvider[keycloakProvider]; exists {
			if err := s.keycloak.UpdateUser(ctx, identity.ProviderUserId, keycloak.UpdateUserInput{
				Username:    user.Username,
				Email:       user.Email,
				DisplayName: user.DisplayName,
				Password:    password,
			}); err != nil {
				return err
			}
			// if err := s.UserRepository.UpdateIdentityProfile(ctx, user.ID, keycloakProvider, identity.ProviderUserId, user.Username, user.Email); err != nil {
			// 	return err
			// }
		} else {
			if _, err := s.createKeycloakUser(ctx, user, password, false, true); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *UserServiceImpl) Delete(ctx context.Context, id int64) error {
	end := s.log.Start(ctx, "Delete", "id", id)

	if _, err := s.UserRepository.FindByID(ctx, id); err != nil {
		end(err)
		return err
	}

	identities, err := s.UserRepository.FindIdentitiesByUserID(ctx, id)
	if err != nil {
		end(err)
		return err
	}
	for _, identity := range identities {
		switch identity.Provider {
		case keycloakProvider:
			if err := s.deleteKeycloakUser(ctx, identity.ProviderUserId); err != nil {
				end(err, "provider", keycloakProvider)
				return err
			}
		case freeIpaProvider:
			if err := s.deleteFreeIPAUser(ctx, identity.ProviderUserId); err != nil {
				end(err, "provider", freeIpaProvider)
				return err
			}
		}
	}

	err = s.UserRepository.Delete(ctx, id)
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

func (s *UserServiceImpl) sendEmailVerification(ctx context.Context, email string, displayName string, token string) error {
	link := fmt.Sprintf("%s/api/v1/users/verify-email?code=%s", s.cfg.App.PublicURL, token)
	emailContent, err := s.renderUserEmail(ctx, userEmailTemplateData{
		SubjectKey:      "mail.users.verify_email.subject",
		TitleKey:        "mail.users.verify_email.title",
		PreheaderKey:    "mail.users.verify_email.preheader",
		GreetingKey:     "mail.users.greeting",
		IntroKey:        "mail.users.verify_email.intro",
		ButtonKey:       "mail.users.verify_email.button",
		FooterKey:       "mail.users.verify_email.footer",
		LinkLabelKey:    "mail.users.verify_email.link_label",
		ExpiresLabelKey: "mail.users.verify_email.expires_label",
		ExpiresInKey:    "mail.users.verify_email.expires_in",
		DisplayName:     displayName,
		Email:           email,
		ActionURL:       link,
	})
	if err != nil {
		return err
	}
	return s.mailer.SendHTML(ctx, email, emailContent.Subject, emailContent.TextBody, emailContent.HTMLBody)
}

func (s *UserServiceImpl) sendPasswordSetupEmail(ctx context.Context, email string, displayName string, token string, setupPasswordURL string) error {
	link := s.passwordSetupLink(token, setupPasswordURL)
	emailContent, err := s.renderUserEmail(ctx, userEmailTemplateData{
		SubjectKey:      "mail.users.invitation.subject",
		TitleKey:        "mail.users.invitation.title",
		PreheaderKey:    "mail.users.invitation.preheader",
		GreetingKey:     "mail.users.greeting",
		IntroKey:        "mail.users.invitation.intro",
		ButtonKey:       "mail.users.invitation.button",
		FooterKey:       "mail.users.invitation.footer",
		LinkLabelKey:    "mail.users.invitation.link_label",
		ExpiresLabelKey: "mail.users.invitation.expires_label",
		ExpiresInKey:    "mail.users.invitation.expires_in",
		DisplayName:     displayName,
		Email:           email,
		ActionURL:       link,
	})
	if err != nil {
		return err
	}
	return s.mailer.SendHTML(ctx, email, emailContent.Subject, emailContent.TextBody, emailContent.HTMLBody)
}

func (s *UserServiceImpl) passwordSetupLink(token string, setupPasswordURL string) string {
	setupPasswordURL = strings.TrimSpace(setupPasswordURL)
	if setupPasswordURL == "" {
		setupPasswordURL = fmt.Sprintf("%s/password/setup", s.cfg.App.PublicURL)
	}

	parsed, err := url.Parse(setupPasswordURL)
	if err != nil {
		return fmt.Sprintf("%s?code=%s", setupPasswordURL, url.QueryEscape(token))
	}
	values := parsed.Query()
	values.Set("code", token)
	parsed.RawQuery = values.Encode()
	return parsed.String()
}

type userEmailTemplateData struct {
	SubjectKey      string
	TitleKey        string
	PreheaderKey    string
	GreetingKey     string
	IntroKey        string
	ButtonKey       string
	FooterKey       string
	LinkLabelKey    string
	ExpiresLabelKey string
	ExpiresInKey    string
	DisplayName     string
	Email           string
	ActionURL       string
}

type renderedUserEmail struct {
	Subject  string
	TextBody string
	HTMLBody string
}

type userEmailView struct {
	Subject      string
	Preheader    string
	Title        string
	Greeting     string
	Intro        string
	Button       string
	Footer       string
	LinkLabel    string
	ExpiresLabel string
	ExpiresIn    string
	ActionURL    string
	PublicURL    string
}

func (s *UserServiceImpl) renderUserEmail(ctx context.Context, data userEmailTemplateData) (renderedUserEmail, error) {
	language := requestctx.Language(ctx)
	recipientName := strings.TrimSpace(data.DisplayName)
	if recipientName == "" {
		recipientName = data.Email
	}

	params := map[string]string{
		"name": recipientName,
		"link": data.ActionURL,
	}
	view := userEmailView{
		Subject:      s.translate(language, data.SubjectKey, nil),
		Preheader:    s.translate(language, data.PreheaderKey, nil),
		Title:        s.translate(language, data.TitleKey, nil),
		Greeting:     s.translate(language, data.GreetingKey, map[string]string{"name": recipientName}),
		Intro:        s.translate(language, data.IntroKey, nil),
		Button:       s.translate(language, data.ButtonKey, nil),
		Footer:       s.translate(language, data.FooterKey, nil),
		LinkLabel:    s.translate(language, data.LinkLabelKey, nil),
		ExpiresLabel: s.translate(language, data.ExpiresLabelKey, nil),
		ExpiresIn:    s.translate(language, data.ExpiresInKey, nil),
		ActionURL:    data.ActionURL,
		PublicURL:    s.cfg.App.PublicURL,
	}

	textBody := strings.Join([]string{
		view.Greeting,
		"",
		view.Intro,
		"",
		view.LinkLabel + ": " + data.ActionURL,
		view.ExpiresLabel + ": " + view.ExpiresIn,
		"",
		view.Footer,
		"",
		s.translate(language, "mail.users.plain_link_hint", params),
	}, "\n")

	body := bytes.Buffer{}
	tmpl, err := template.New("user_email").Parse(userEmailHTMLTemplate)
	if err != nil {
		return renderedUserEmail{}, err
	}
	if err := tmpl.Execute(&body, view); err != nil {
		return renderedUserEmail{}, err
	}

	return renderedUserEmail{
		Subject:  view.Subject,
		TextBody: textBody,
		HTMLBody: body.String(),
	}, nil
}

func (s *UserServiceImpl) translate(language string, key string, params map[string]string) string {
	if s.translator == nil {
		return key
	}
	return s.translator.Translate(language, key, params)
}

const userEmailHTMLTemplate = `<!doctype html>
	<html lang="en">
	<head>
	<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<title>{{.Subject}}</title>
	</head>
	<body style="margin:0;padding:0;background:#f4f7fb;color:#172033;font-family:Arial,Helvetica,sans-serif;">
	<div style="display:none;max-height:0;overflow:hidden;opacity:0;color:transparent;">{{.Preheader}}</div>
	<table role="presentation" width="100%" cellspacing="0" cellpadding="0" style="width:100%;background:#f4f7fb;padding:32px 12px;">
		<tr>
		<td align="center">
			<table role="presentation" width="100%" cellspacing="0" cellpadding="0" style="width:100%;max-width:640px;background:#ffffff;border:1px solid #dde5f0;border-radius:8px;overflow:hidden;">
			<tr>
				<td style="background:#10243f;padding:28px 32px;">
				<div style="font-size:13px;line-height:18px;letter-spacing:0;color:#9bd3ff;font-weight:700;">Open Suite Authorization</div>
				<h1 style="margin:12px 0 0;font-size:26px;line-height:34px;color:#ffffff;font-weight:700;">{{.Title}}</h1>
				</td>
			</tr>
			<tr>
				<td style="padding:32px;">
				<p style="margin:0 0 16px;font-size:16px;line-height:24px;color:#172033;">{{.Greeting}}</p>
				<p style="margin:0 0 26px;font-size:16px;line-height:26px;color:#3b4658;">{{.Intro}}</p>
				<table role="presentation" cellspacing="0" cellpadding="0" style="margin:0 0 26px;">
					<tr>
					<td bgcolor="#1769e0" style="border-radius:6px;">
						<a href="{{.ActionURL}}" style="display:inline-block;padding:14px 22px;font-size:15px;line-height:20px;color:#ffffff;text-decoration:none;font-weight:700;">{{.Button}}</a>
					</td>
					</tr>
				</table>
				<table role="presentation" width="100%" cellspacing="0" cellpadding="0" style="width:100%;background:#f7f9fc;border:1px solid #e3e9f2;border-radius:8px;margin:0 0 24px;">
					<tr>
					<td style="padding:18px 20px;">
						<div style="font-size:13px;line-height:18px;color:#637083;font-weight:700;">{{.ExpiresLabel}}</div>
						<div style="margin-top:4px;font-size:15px;line-height:22px;color:#172033;">{{.ExpiresIn}}</div>
					</td>
					</tr>
				</table>
				<p style="margin:0 0 8px;font-size:13px;line-height:20px;color:#637083;">{{.LinkLabel}}</p>
				<p style="margin:0;word-break:break-all;font-size:13px;line-height:20px;color:#1769e0;"><a href="{{.ActionURL}}" style="color:#1769e0;text-decoration:underline;">{{.ActionURL}}</a></p>
				</td>
			</tr>
			<tr>
				<td style="padding:22px 32px;background:#eef3f9;border-top:1px solid #dde5f0;">
				<p style="margin:0;font-size:13px;line-height:20px;color:#526174;">{{.Footer}}</p>
				</td>
			</tr>
			</table>
		</td>
		</tr>
	</table>
	</body>
	</html>`

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

func (s *UserServiceImpl) deleteKeycloakUser(ctx context.Context, userID string) error {
	if strings.TrimSpace(userID) == "" {
		return nil
	}
	if !s.keycloak.Enabled() {
		s.log.Warn(ctx, "keycloak.delete_skipped", "user_id", userID, "reason", "disabled")
		return nil
	}
	if err := s.keycloak.DeleteUser(ctx, userID); err != nil {
		s.log.Warn(ctx, "keycloak.delete_failed", "user_id", userID, "error", err.Error())
		return err
	}
	return nil
}

func (s *UserServiceImpl) deleteFreeIPAUser(ctx context.Context, uid string) error {
	if uid == "" {
		return nil
	}
	if !s.freeipa.Enabled() {
		s.log.Warn(ctx, "freeipa.delete_skipped", "uid", uid, "reason", "disabled")
		return nil
	}
	if err := s.freeipa.DeleteUser(ctx, uid); err != nil {
		s.log.Warn(ctx, "freeipa.delete_failed", "uid", uid, "error", err.Error())
		return err
	}
	return nil
}

func (s *UserServiceImpl) toUserResponseWithAssignments(ctx context.Context, user *entities.User, mustChangePassword bool, provisionedTo []string) (*dto.UserResponse, error) {
	roles, err := s.UserRoleService.FindAssignedRolesByUserID(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	teams, err := s.TeamMemberService.FindAssignedTeamsByUserID(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	roleIDs := make([]int64, 0, len(roles))
	for _, role := range roles {
		roleIDs = append(roleIDs, role.ID)
	}
	teamIDs := make([]int64, 0, len(teams))
	for _, team := range teams {
		teamIDs = append(teamIDs, team.ID)
	}

	response := toUserResponse(user, mustChangePassword, provisionedTo, roleIDs, teamIDs)
	response.Roles = roles
	response.Teams = teams
	return response, nil
}

func (s *UserServiceImpl) resetUserAccessCache(ctx context.Context, userID int64) error {
	if s.redis == nil || s.redis.Client == nil {
		return nil
	}

	keys := []string{
		fmt.Sprintf("authz:apps:user:%d", userID),
		fmt.Sprintf("authz:me:user:%d", userID),
	}

	pattern := fmt.Sprintf("authz:access:user:%d:app:*", userID)
	var cursor uint64
	for {
		items, nextCursor, err := s.redis.Client.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return err
		}
		keys = append(keys, items...)
		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	return s.redis.Client.Del(ctx, keys...).Err()
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

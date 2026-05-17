package auditlogs

import (
	"net/http"

	"github.com/open-suite/authorization/internal/modules/auditlogs/controllers"
)

type AuditLogModuleImpl struct {
	AuditLogController controllers.AuditLogController
}

func NewAuditLogModule(controller controllers.AuditLogController) *AuditLogModuleImpl {
	return &AuditLogModuleImpl{
		AuditLogController: controller,
	}
}

func (m *AuditLogModuleImpl) Name() string {
	return "audit_logs"
}

func (m *AuditLogModuleImpl) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/audit-logs", m.AuditLogController.Find)
	mux.HandleFunc("GET /api/v1/audit-logs/{id}", m.AuditLogController.FindByID)
	mux.HandleFunc("POST /api/v1/audit-logs", m.AuditLogController.Create)
	mux.HandleFunc("PUT /api/v1/audit-logs/{id}", m.AuditLogController.Update)
	mux.HandleFunc("DELETE /api/v1/audit-logs/{id}", m.AuditLogController.Delete)
}

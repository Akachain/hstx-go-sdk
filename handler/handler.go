package handler

// Handler ...
type Handler struct {
	SuperAdminHandler *SuperAdminHandler
	AdminHandler      *AdminHandler
	ProposalHandler   *ProposalHandler
	ApprovalHandler   *ApprovalHandler
}

// InitHandler ...
func (h *Handler) InitHandler() {
	h.SuperAdminHandler = new(SuperAdminHandler)
	h.AdminHandler = new(AdminHandler)
	h.ProposalHandler = new(ProposalHandler)
	h.ApprovalHandler = new(ApprovalHandler)
}

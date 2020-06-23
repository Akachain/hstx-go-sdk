package model

// ApprovalTable - Table name
const ApprovalTable = "HSTX_APPROVAL"

// Approval contain a Super Admin's signature to Approve or Reject a Proposal
type Approval struct {
	ApprovalID string `json:"ApprovalID"`	// set
	ProposalID string `json:"ProposalID"`	// args[0] proposalID
	ApproverID string `json:"ApproverID"`	// args[0] approverID
	Challenge  string `json:"Challenge"`	// args[0] singned challenge
	Signature  string `json:"Signature"`	// args[0] signature
	Message    string `json:"Message"`		// args[0] singned Message
	Status     string `json:"Status"`		// args[0] approval status: Approved/Rejected
	CreatedAt  string `json:"CreatedAt"`	// set
}

package model

// AdminTable - Table name
const AdminTable = "HSTX_ADMIN"

// Admin who is able be create a Proposal
type Admin struct {
	AdminID string `json:"AdminID"`
	Name    string `json:"Name"`
	Status  string `json:"Status"`
}

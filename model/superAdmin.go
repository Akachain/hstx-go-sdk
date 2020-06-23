package model

// SuperAdminTable - Table name
const SuperAdminTable = "HSTX_SUPER_ADMIN"

// SuperAdmin , who has permission to approve or reject a proposal, is a member in the Quorum
type SuperAdmin struct {
	SuperAdminID string `json:"SuperAdminID"`	// args[0] keyhandle of yubikey and application
	Name         string `json:"Name"`			// args[0] name
	PublicKey    string `json:"PublicKey"`		// args[0] publickey of yubikey (format: pem)
	Status       string `json:"Status"`			// args[0] A/I (active/inactive)
}

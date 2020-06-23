package main

import (
	"encoding/json"
	"testing"

	"github.com/Akachain/hstx-go-sdk/model"
	"gotest.tools/assert"

	"github.com/Akachain/akc-go-sdk/common"
	"github.com/Akachain/akc-go-sdk/util"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var superAdminID string
var adminID string
var proposalID string
var approvalID string
var stub = setupMock(false)

func setupMock(isDropDB bool) *util.MockStubExtend {

	// Initialize mockstubextend
	cc := new(Chaincode)
	stub := util.NewMockStubExtend(shim.NewMockStub("TestMockStub", cc), cc)

	// Create a new database, Drop old database
	db, _ := util.NewCouchDBHandlerWithConnectionAuthentication(isDropDB)
	stub.SetCouchDBConfiguration(db)
	return stub
}

func TestCreateSuperAdmin(t *testing.T) {
	common.Logger.SetLevel(shim.LogDebug)
	
	stub = setupMock(true)

	superAdminID = "SYduDJxe6-MeAqyzGYqUB9LXK0e79o63OH2Tp7npcGdMG_IfaN6WAqfIfs388HlHjW9PIE2tP7MPGxzof6406g"

	superAdmin := model.SuperAdmin{
		SuperAdminID: superAdminID,
		Name:         "TestSuperAdmin" + superAdminID,
		PublicKey:    "-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEeSohIzvASeU0nZADWDFS/WI/0+AU\n8tRQrGtRtaZ2hv6JCdq+BK1euL4VFPqVw5ddYbw0+t4b5zQDah7qrP41ug==\n-----END PUBLIC KEY-----\n",
		Status:       "A",
	}

	superAdminBytes, _ := json.Marshal(superAdmin)

	// Create a new Super Admin
	response := util.MockInvokeTransaction(t, stub, [][]byte{[]byte("CreateSuperAdmin"), superAdminBytes})
	var result map[string]interface{}
	json.Unmarshal([]byte(response), &result)
	if result["status"] != nil {
		common.Logger.Debug("Error: ", result["msg"])
	}

	var createdSuperAdmin model.SuperAdmin
	json.Unmarshal([]byte(response), &createdSuperAdmin)
	common.Logger.Debugf("%v", createdSuperAdmin)

	assert.Equal(t, superAdmin.SuperAdminID, createdSuperAdmin.SuperAdminID)
	assert.Equal(t, superAdmin.Name, createdSuperAdmin.Name)
	assert.Equal(t, superAdmin.PublicKey, createdSuperAdmin.PublicKey)
	assert.Equal(t, superAdmin.Status, createdSuperAdmin.Status)

	// Check if the created data exists
	compositeKey, _ := stub.CreateCompositeKey(model.SuperAdminTable, []string{superAdmin.SuperAdminID})
	state, _ := stub.GetState(compositeKey)

	var stateSuperAdmin model.SuperAdmin
	json.Unmarshal([]byte(state), &stateSuperAdmin)

	// Check if the created user information is correct
	assert.Equal(t, superAdmin.SuperAdminID, stateSuperAdmin.SuperAdminID)
	assert.Equal(t, superAdmin.Name, stateSuperAdmin.Name)
	assert.Equal(t, superAdmin.PublicKey, stateSuperAdmin.PublicKey)
	assert.Equal(t, superAdmin.Status, stateSuperAdmin.Status)
}

func TestCreateAdmin(t *testing.T) {
	common.Logger.SetLevel(shim.LogDebug)
	
	if stub == nil {
		stub = setupMock(false)
	}

	admin := model.Admin{
		Name: "SonPH",
	}

	adminBytes, _ := json.Marshal(admin)

	// Create a new Admin
	response := util.MockInvokeTransaction(t, stub, [][]byte{[]byte("CreateAdmin"), adminBytes})
	var result map[string]interface{}
	json.Unmarshal([]byte(response), &result)
	if result["status"] != nil {
		common.Logger.Debug("Error: ", result["msg"])
	}

	var createdAdmin model.Admin
	json.Unmarshal([]byte(response), &createdAdmin)
	common.Logger.Debugf("%v", createdAdmin)

	assert.Equal(t, admin.Name, createdAdmin.Name)
	assert.Equal(t, "Active", createdAdmin.Status)

	// Check if the created data exists
	compositeKey, _ := stub.CreateCompositeKey(model.AdminTable, []string{createdAdmin.AdminID})
	state, _ := stub.GetState(compositeKey)

	var stateAdmin model.Admin
	json.Unmarshal([]byte(state), &stateAdmin)

	// Check if the created user information is correct
	assert.Equal(t, createdAdmin.AdminID, stateAdmin.AdminID)
	assert.Equal(t, admin.Name, stateAdmin.Name)
	assert.Equal(t, createdAdmin.Status, stateAdmin.Status)

	adminID = stateAdmin.AdminID
}

func TestCreateProposal(t *testing.T) {
	common.Logger.SetLevel(shim.LogDebug)
	
	if stub == nil {
		stub = setupMock(false)
	}

	proposal := model.Proposal{
		CreatedBy: "Admin1",
		Message: "Chuyển 1 tỷ cho anh Long sex",
		QuorumNumber: 1,
	}

	proposalBytes, _ := json.Marshal(proposal)

	// Create a new Proposal
	response := util.MockInvokeTransaction(t, stub, [][]byte{[]byte("CreateProposal"), proposalBytes})
	var result map[string]interface{}
	json.Unmarshal([]byte(response), &result)
	if result["status"] != nil {
		common.Logger.Debug("Error: ", result["msg"])
	}

	var createdProposal model.Proposal
	json.Unmarshal([]byte(response), &createdProposal)
	common.Logger.Debugf("%v", createdProposal)

	assert.Equal(t, proposal.CreatedBy, createdProposal.CreatedBy)
	assert.Equal(t, proposal.Message, createdProposal.Message)
	assert.Equal(t, "Pending", createdProposal.Status)

	// Check if the created data exists
	compositeKey, _ := stub.CreateCompositeKey(model.ProposalTable, []string{createdProposal.ProposalID})
	state, _ := stub.GetState(compositeKey)

	var stateProposal model.Proposal
	json.Unmarshal([]byte(state), &stateProposal)

	// Check if the created user information is correct
	assert.Equal(t, createdProposal.ProposalID, stateProposal.ProposalID)
	assert.Equal(t, proposal.Message, stateProposal.Message)
	assert.Equal(t, proposal.CreatedBy, stateProposal.CreatedBy)
	assert.Equal(t, createdProposal.Status, stateProposal.Status)
	assert.Equal(t, createdProposal.CreatedAt, stateProposal.CreatedAt)
	assert.Equal(t, createdProposal.UpdatedAt, stateProposal.UpdatedAt)

	proposalID = stateProposal.ProposalID
}

func TestCreateApproval(t *testing.T) {
	common.Logger.SetLevel(shim.LogDebug)
	
	if stub == nil {
		stub = setupMock(false)
	}

	approval := model.Approval{
		ProposalID: proposalID,
		ApproverID: superAdminID,
		Challenge: "Q2h1eeG7g24gMSB04bu3IGNobyBhbmggTG9uZyBz4bq9",
		Signature: "MEUCIQDaeCaR1s33mIse0i5hq69srNfBHTajF47pibvJoLY7KwIgD/eerEBfeOTizizYIinZ7WLJ90UFpS6IBuEe0sqpYaM=",
		Message: "EWJcgjcFkTyobMwx1j6p2xtEmKE8H/ctu5jr4jVilUgBAAAABSpqAxztD4zx7PrOdW1ur4ba/eE1yjhPtYY2QCNcOz4N",
		Status: "Approved",
		// Status: "Rejected",
	}

	approvalBytes, _ := json.Marshal(approval)

	// Create a new Approval
	response := util.MockInvokeTransaction(t, stub, [][]byte{[]byte("CreateApproval"), approvalBytes})
	
	var result map[string]interface{}
	json.Unmarshal([]byte(response), &result)
	if result["status"] != nil {
		common.Logger.Debug("Error: ", result["msg"])
	} else {
		common.Logger.Debug(response)
	}

	var createdApproval model.Approval
	json.Unmarshal([]byte(response), &createdApproval)

	assert.Equal(t, approval.ProposalID, createdApproval.ProposalID)
	assert.Equal(t, approval.ApproverID, createdApproval.ApproverID)
	assert.Equal(t, approval.Challenge, createdApproval.Challenge)
	assert.Equal(t, approval.Signature, createdApproval.Signature)
	assert.Equal(t, approval.Message, createdApproval.Message)
	assert.Equal(t, approval.Status, createdApproval.Status)

	// Check if the created data exists
	compositeKey, _ := stub.CreateCompositeKey(model.ApprovalTable, []string{createdApproval.ProposalID, createdApproval.ApproverID})
	state, _ := stub.GetState(compositeKey)

	var stateApproval model.Approval
	json.Unmarshal([]byte(state), &stateApproval)

	// Check if the created user information is correct
	assert.Equal(t, createdApproval.ApprovalID, stateApproval.ApprovalID)
	assert.Equal(t, approval.ProposalID, stateApproval.ProposalID)
	assert.Equal(t, approval.ApproverID, stateApproval.ApproverID)
	assert.Equal(t, approval.Challenge, stateApproval.Challenge)
	assert.Equal(t, approval.Signature, stateApproval.Signature)
	assert.Equal(t, approval.Message, stateApproval.Message)
	assert.Equal(t, approval.Status, stateApproval.Status)
}

func TestCommitProposal(t *testing.T) {
	common.Logger.SetLevel(shim.LogDebug)
	
	if stub == nil {
		stub = setupMock(false)
	}

	// Create a new Approval
	response := util.MockInvokeTransaction(t, stub, [][]byte{[]byte("CommitProposal"), []byte(proposalID)})
	
	var result map[string]interface{}
	json.Unmarshal([]byte(response), &result)
	if result["status"] != nil {
		common.Logger.Debug("Error: ", result["msg"])
	} else {
		common.Logger.Debug(response)
	}

	var proposal model.Proposal
	json.Unmarshal([]byte(response), &proposal)

	assert.Equal(t, proposalID, proposal.ProposalID)

	// Check if the created data exists
	compositeKey, _ := stub.CreateCompositeKey(model.ProposalTable, []string{proposalID})
	state, _ := stub.GetState(compositeKey)

	var stateProposal model.Proposal
	json.Unmarshal([]byte(state), &stateProposal)

	// Check if the created user information is correct
	assert.Equal(t, proposal.ProposalID, stateProposal.ProposalID)
	assert.Equal(t, proposal.CreatedBy, stateProposal.CreatedBy)
	assert.Equal(t, proposal.Message, stateProposal.Message)
	assert.Equal(t, proposal.Status, stateProposal.Status)
	assert.Equal(t, proposal.CreatedAt, stateProposal.CreatedAt)
	assert.Equal(t, proposal.UpdatedAt, stateProposal.UpdatedAt)
}
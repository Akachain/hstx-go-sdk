package handler

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/Akachain/akc-go-sdk/common"
	"github.com/Akachain/akc-go-sdk/util"
	"github.com/Akachain/hstx-go-sdk/model"
	hUtil "github.com/Akachain/hstx-go-sdk/utils"
	"github.com/hyperledger/fabric/bccsp/utils"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/mitchellh/mapstructure"
)

// ApprovalHandler ...
type ApprovalHandler struct{}

// CreateApproval ...
func (sah *ApprovalHandler) CreateApproval(stub shim.ChaincodeStubInterface, approvalStr string) (result *string, err error) {
	common.Logger.Debugf("Input-data sent to CreateApproval func: %+v\n", approvalStr)

	// Check role: SuperAdmin
	// err = hUtil.IsSuperAdmin(stub)
	// if err != nil {
	// 	return nil, fmt.Errorf("%s %s", err.Error(), common.GetLine())
	// }

	// Parse approvalStr to approval
	approval := new(model.Approval)
	err = json.Unmarshal([]byte(approvalStr), approval)
	if err != nil { // Return error: Can't unmarshal json
		return nil, fmt.Errorf("%s %s %s", common.ResCodeDict[common.ERR3], err.Error(), common.GetLine())
	}

	// Check SuperAdmin's status
	err = sah.checkApproverStatus(stub, approval.ApproverID)
	if err != nil {
		return nil, fmt.Errorf("%s %s %s", "This approver is not active", err.Error(), common.GetLine())
	}

	// Get proposal by approval.ProposalID
	proposalStr, err := new(ProposalHandler).GetProposalByID(stub, approval.ProposalID)
	if err != nil {
		return nil, fmt.Errorf("%s %s %s", common.ResCodeDict[common.ERR4], err.Error(), common.GetLine())
	}

	var proposal model.Proposal
	err = json.Unmarshal([]byte(*proposalStr), &proposal)
	if err != nil {
		return nil, fmt.Errorf("%s %s %s", common.ResCodeDict[common.ERR3], err.Error(), common.GetLine())
	}

	// Check whether the proposal was rejected or not
	if strings.Compare("Rejected", proposal.Status) == 0 {
		return nil, fmt.Errorf("%s %s", "The proposal was rejected", common.GetLine())
	}

	// Check this approver hasn't signed the proposal
	compositeKey, _ := stub.CreateCompositeKey(model.ApprovalTable, []string{approval.ProposalID, approval.ApproverID})
	rs, err := stub.GetState(compositeKey)
	if err != nil { // Return error: Fail to get data
		return nil, fmt.Errorf("%s %s %s", common.ResCodeDict[common.ERR4], err.Error(), common.GetLine())
	}
	if len(rs) > 0 { // Return error: Only signing once
		return nil, fmt.Errorf("%s %s %s", common.ResCodeDict[common.ERR9], "This proposal had already been approved", common.GetLine())
	}

	// Verify signature with the singed message
	err = sah.verifySignature(stub, approval.ApproverID, approval.Signature, approval.Message)
	if err != nil { // Return error: Verify error
		return nil, fmt.Errorf("%s %s %s", common.ResCodeDict[common.ERR8], err.Error(), common.GetLine())
	}

	// Set approval.ApprovalID & approval.CreatedAt
	approval.ApprovalID = hUtil.GenerateDocumentID(stub)
	timestamp, err := stub.GetTxTimestamp()
	if err != nil {
		return nil, fmt.Errorf("%s %s %s", common.ResCodeDict[common.ERR4], err.Error(), common.GetLine())
	}
	formatedTime := time.Unix(timestamp.Seconds, 0)
	approval.CreatedAt = formatedTime.String()

	// Create Approval
	common.Logger.Infof("Creating Approval: %+v\n", approval)
	err = util.Createdata(stub, model.ApprovalTable, []string{approval.ProposalID, approval.ApproverID}, &approval)
	if err != nil { // Return error: Fail to insert data
		return nil, fmt.Errorf("%s %s %s", common.ResCodeDict[common.ERR5], err.Error(), common.GetLine())
	}

	// Update proposal if necessary
	sah.updateProposal(stub, approval)

	bytes, err := json.Marshal(approval)
	if err != nil { // Return error: Can't marshal json
		return nil, fmt.Errorf("%s %s %s", common.ResCodeDict[common.ERR3], err.Error(), common.GetLine())
	}
	temp := ""
	result = &temp
	*result = string(bytes)

	return result, nil
}

// GetAllApproval ...
func (sah *ApprovalHandler) GetAllApproval(stub shim.ChaincodeStubInterface) (result *string, err error) {
	res := util.GetAllData(stub, new(model.Approval), model.ApprovalTable)
	if res.Status == 200 {
		return &res.Message, nil
	}
	return nil, fmt.Errorf("%s %s %s", common.ResCodeDict[common.ERR4], res.Message, common.GetLine())
}

// GetApprovalByID ...
func (sah *ApprovalHandler) GetApprovalByID(stub shim.ChaincodeStubInterface, approvalID string) (result *string, err error) {
	common.Logger.Debugf("Input-data sent to GetApprovalByID func: %+v\n", approvalID)

	rawApproval, err := util.Getdatabyid(stub, approvalID, model.ApprovalTable)
	if err != nil {
		return nil, fmt.Errorf("%s %s %s", common.ResCodeDict[common.ERR4], err.Error(), common.GetLine())
	}

	approval := new(model.Approval)
	mapstructure.Decode(rawApproval, approval)

	bytes, err := json.Marshal(approval)
	if err != nil { // Return error: Can't marshal json
		return nil, fmt.Errorf("%s %s %s", common.ResCodeDict[common.ERR3], err.Error(), common.GetLine())
	}
	temp := ""
	result = &temp
	*result = string(bytes) 

	return result, nil
}

// UpdateApproval ...
func (sah *ApprovalHandler) UpdateApproval(stub shim.ChaincodeStubInterface, approvalStr string) (result *string, err error) {
	common.Logger.Debugf("Input-data sent to UpdateApproval func: %+v\n", approvalStr)

	err = hUtil.IsSuperAdmin(stub)
	if err != nil {
		return nil, fmt.Errorf("%s %s", err.Error(), common.GetLine())
	}

	newApproval := new(model.Approval)
	err = json.Unmarshal([]byte(approvalStr), newApproval)
	if err != nil { // Return error: Can't unmarshal json
		return nil, fmt.Errorf("%s %s %s", common.ResCodeDict[common.ERR3], err.Error(), common.GetLine())
	}

	if len(newApproval.ApprovalID) == 0 {
		return nil, fmt.Errorf("%s %s", "This ApprovalID can't be empty", common.GetLine())
	}

	// Get approval information
	rawApproval, err := util.Getdatabyid(stub, newApproval.ApprovalID, model.ApprovalTable)
	if err != nil {
		return nil, fmt.Errorf("%s %s %s", common.ResCodeDict[common.ERR4], err.Error(), common.GetLine())
	}

	approval := new(model.Approval)
	mapstructure.Decode(rawApproval, approval)

	// Filter fields needed to update
	newApprovalValue := reflect.ValueOf(newApproval).Elem()
	approvalValue := reflect.ValueOf(approval).Elem()
	for i := 0; i < newApprovalValue.NumField(); i++ {
		fieldName := newApprovalValue.Type().Field(i).Name
		if len(newApprovalValue.Field(i).String()) > 0 {
			field := approvalValue.FieldByName(fieldName)
			if field.CanSet() {
				field.SetString(newApprovalValue.Field(i).String())
			}
		}
	}

	err = util.Changeinfo(stub, model.ApprovalTable, []string{approval.ApprovalID}, approval)
	if err != nil { // Return error: Fail to Update data
		return nil, fmt.Errorf("%s %s %s", common.ResCodeDict[common.ERR5], err.Error(), common.GetLine())
	}

	bytes, err := json.Marshal(approval)
	if err != nil { // Return error: Can't marshal json
		return nil, fmt.Errorf("%s %s %s", common.ResCodeDict[common.ERR3], err.Error(), common.GetLine())
	}
	temp := ""
	result = &temp
	*result = string(bytes)

	return result, nil
}

// checkApproverStatus func to check whether the SuperAdmin is active or inactive
func (sah *ApprovalHandler) checkApproverStatus(stub shim.ChaincodeStubInterface, approverID string) error {
	// Get approver by approval.ApproverID
	superAdminStr, err := new(SuperAdminHandler).GetSuperAdminByID(stub, approverID)
	if err != nil {
		return fmt.Errorf("%s %s %s", common.ResCodeDict[common.ERR4], err.Error(), common.GetLine())
	}

	var superAdmin model.SuperAdmin
	err = json.Unmarshal([]byte(*superAdminStr), &superAdmin)
	if err != nil {
		return fmt.Errorf("%s %s %s", common.ResCodeDict[common.ERR3], err.Error(), common.GetLine())
	}

	// Check SuperAdmin's status
	if superAdmin.Status != "A" && superAdmin.Status != "Active" {
		return fmt.Errorf("%s %s", "This approver is not active", common.GetLine())
	}
	// If the SuperAdmin is active, return nil
	return nil
}

// verifySignature ...
func (sah *ApprovalHandler) verifySignature(stub shim.ChaincodeStubInterface, approverID string, signature string, message string) error {
	if len(approverID) == 0 {
		return errors.New("approverID is empty")
	}

	//get superAdmin information
	rawSuperAdmin, err := util.Getdatabyid(stub, approverID, model.SuperAdminTable)
	if err != nil {
		return err
	}

	superAdmin := new(model.SuperAdmin)
	mapstructure.Decode(rawSuperAdmin, superAdmin)

	// Start verify
	pkBytes := []byte(superAdmin.PublicKey)
	pkBlock, _ := pem.Decode(pkBytes)
	if pkBlock == nil {
		return errors.New("can't decode public key")
	}

	rawPk, err := x509.ParsePKIXPublicKey(pkBlock.Bytes)
	if err != nil {
		return err
	}

	pk := rawPk.(*ecdsa.PublicKey)

	// SIGNATURE
	signatureByte, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return err
	}

	R, S, err := utils.UnmarshalECDSASignature(signatureByte)
	if err != nil {
		return err
	}

	// DATA
	dataByte, err := base64.StdEncoding.DecodeString(message)
	if err != nil {
		return err
	}

	hash := sha256.Sum256(dataByte)
	var hashData = hash[:]

	// VERIFY
	checkedResult := ecdsa.Verify(pk, hashData, R, S)

	if checkedResult {
		return nil
	}
	return errors.New("verifying failed")
}

// updateProposal ...
func (sah *ApprovalHandler) updateProposal(stub shim.ChaincodeStubInterface, approval *model.Approval) error {
	rawProposal, err := util.Getdatabyid(stub, approval.ProposalID, model.ProposalTable)
	if err != nil {
		return err
	}

	proposal := new(model.Proposal)
	mapstructure.Decode(rawProposal, proposal)

	if strings.Compare(approval.Status, "Rejected") == 0 {
		if strings.Compare(proposal.Status, "Commited") != 0 {
			proposal.Status = approval.Status
			proposal.UpdatedAt = approval.CreatedAt
			bytes, err := json.Marshal(proposal)
			if err != nil {
				return err
			}
			new(ProposalHandler).UpdateProposal(stub, string(bytes))
		}
		return nil
	}

	resIterator, err := hUtil.GetContainKey(stub, model.ApprovalTable, approval.ProposalID)
	if err != nil {
		return err
	}
	defer resIterator.Close()
	count := 0
	if approval.Status == "Approved" {
		count ++
	}
	for resIterator.HasNext() {
		stateIterator, err := resIterator.Next()
		if err != nil {
			return err
		}
		approvalState := new(model.Approval)
		err = json.Unmarshal(stateIterator.Value, approvalState)
		if err != nil { // Convert JSON error
			return err
		}
		
		if strings.Compare("Approved", approvalState.Status) == 0 {
			count++
		} 
	}
	// Check approved number >= proposal.QuorumNumber to update the Proposal's satatus
	if count >= proposal.QuorumNumber {
		rawProposal, err := util.Getdatabyid(stub, approval.ProposalID, model.ProposalTable)
		if err != nil {
			return err
		}
		proposal := new(model.Proposal)
		mapstructure.Decode(rawProposal, proposal)
		if strings.Compare(proposal.Status, "Pending") == 0 {
			proposal.Status = "Approved"
			proposal.UpdatedAt = approval.CreatedAt
			bytes, err := json.Marshal(proposal)
			if err != nil {
				return err
			}
			new(ProposalHandler).UpdateProposal(stub, string(bytes))
		}
	}
	return nil
}
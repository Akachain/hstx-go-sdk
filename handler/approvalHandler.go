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
func (sah *ApprovalHandler) CreateApproval(stub shim.ChaincodeStubInterface, approvalStr string, criteria int) (result *string, err error) {
	common.Logger.Debugf("Input-data sent to CreateApproval func: %+v\n", approvalStr)

	approval := new(model.Approval)
	err = json.Unmarshal([]byte(approvalStr), approval)
	if err != nil { // Return error: Can't unmarshal json
		return nil, fmt.Errorf("%s %s %s", common.ResCodeDict[common.ERR3], err.Error(), common.GetLine())
	}

	rs, err := hUtil.GetByTwoColumns(stub, model.ApprovalTable, "ProposalID", fmt.Sprintf("\"%s\"", approval.ProposalID), "ApproverID", fmt.Sprintf("\"%s\"", approval.ApproverID))
	if err != nil { // Return error: Fail to get data
		return nil, fmt.Errorf("%s %s %s", common.ResCodeDict[common.ERR4], err.Error(), common.GetLine())
	}
	if rs.HasNext() { // Return error: Only signing once
		return nil, fmt.Errorf("%s %s %s", common.ResCodeDict[common.ERR9], "This proposal had already been approved", common.GetLine())
	}

	approval.ApprovalID = stub.GetTxID()

	// Verify signature with the singed message
	err = sah.verifySignature(stub, approval.ApproverID, approval.Signature, approval.Message)
	if err != nil { // Return error: Verify error
		return nil, fmt.Errorf("%s %s %s", common.ResCodeDict[common.ERR8], err.Error(), common.GetLine())
	}

	// Set approval's status
	if len(approval.Status) == 0 {
		approval.Status = "Approved"
	}

	common.Logger.Infof("Creating Approval: %+v\n", approval)
	err = util.Createdata(stub, model.ApprovalTable, []string{approval.ApprovalID}, &approval)
	if err != nil { // Return error: Fail to insert data
		return nil, fmt.Errorf("%s %s %s", common.ResCodeDict[common.ERR5], err.Error(), common.GetLine())
	}

	// Update proposal if necessary
	sah.updateProposal(stub, approval, criteria)

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
	return nil, fmt.Errorf("%s %s %s", common.ResCodeDict[common.ERR4], err.Error(), common.GetLine())
}

// GetApprovalByID ...
func (sah *ApprovalHandler) GetApprovalByID(stub shim.ChaincodeStubInterface, approvalID string) (result *string, err error) {
	common.Logger.Debugf("Input-data sent to GetApprovalByID func: %+v\n", approvalID)

	res := util.GetDataByID(stub, approvalID, new(model.Approval), model.ApprovalTable)
	if res.Status == 200 {
		return &res.Message, nil
	}
	return nil, fmt.Errorf("%s %s %s", common.ResCodeDict[common.ERR4], err.Error(), common.GetLine())
}

// UpdateApproval ...
func (sah *ApprovalHandler) UpdateApproval(stub shim.ChaincodeStubInterface, approvalStr string) (result *string, err error) {
	common.Logger.Debugf("Input-data sent to UpdateApproval func: %+v\n", approvalStr)

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
func (sah *ApprovalHandler) updateProposal(stub shim.ChaincodeStubInterface, approval *model.Approval, criteria int) error {
	if strings.Compare(approval.Status, "Rejected") == 0 {
		rawProposal, err := util.Getdatabyid(stub, approval.ProposalID, model.ProposalTable)
		if err != nil {
			return err
		}

		proposal := new(model.Proposal)
		mapstructure.Decode(rawProposal, proposal)

		if strings.Compare(proposal.Status, "Pending") == 0 {
			proposal.Status = approval.Status
			proposal.CreatedAt = approval.CreatedAt
			bytes, err := json.Marshal(proposal)
			if err != nil {
				return err
			}
			new(ProposalHandler).UpdateProposal(stub, string(bytes))
		}
		return nil
	}

	resIterator, err := hUtil.GetByOneColumn(stub, model.ApprovalTable, "ProposalID", fmt.Sprintf("\"%s\"", approval.ProposalID))
	if err != nil {
		return err
	}
	defer resIterator.Close()
	count := 0
	for resIterator.HasNext() {
		_, err := resIterator.Next()
		if err != nil {
			return err
		}
		count++
	}
	if count >= criteria {
		rawProposal, err := util.Getdatabyid(stub, approval.ProposalID, model.ProposalTable)
		if err != nil {
			return err
		}
		proposal := new(model.Proposal)
		mapstructure.Decode(rawProposal, proposal)
		if strings.Compare(proposal.Status, "Pending") == 0 {
			proposal.Status = "Approved"
			proposal.CreatedAt = approval.CreatedAt
			bytes, err := json.Marshal(proposal)
			if err != nil {
				return err
			}
			new(ProposalHandler).UpdateProposal(stub, string(bytes))
		}
	}
	return nil
}

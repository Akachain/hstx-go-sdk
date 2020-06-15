package handler

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/Akachain/akc-go-sdk/common"
	"github.com/Akachain/akc-go-sdk/util"
	"github.com/Akachain/hstx-go-sdk/model"
	hUtil "github.com/Akachain/hstx-go-sdk/utils"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/mitchellh/mapstructure"
)

// ProposalHandler ...
type ProposalHandler struct{}

// CreateProposal ...
func (sah *ProposalHandler) CreateProposal(stub shim.ChaincodeStubInterface, proposalStr string) (result *string, err error) {
	common.Logger.Debugf("Input-data sent to CreateProposal func: %+v\n", proposalStr)

	proposal := new(model.Proposal)
	err = json.Unmarshal([]byte(proposalStr), proposal)
	if err != nil { // Return error: Can't unmarshal json
		return nil, fmt.Errorf("%s %s %s", common.ResCodeDict[common.ERR3], err.Error(), common.GetLine())
	}

	proposal.ProposalID = stub.GetTxID()
	proposal.Status = "Pending"

	common.Logger.Infof("Create Proposal: %+v\n", proposal)
	err = util.Createdata(stub, model.ProposalTable, []string{proposal.ProposalID}, &proposal)
	if err != nil { // Return error: Fail to insert data
		return nil, fmt.Errorf("%s %s %s", common.ResCodeDict[common.ERR5], err.Error(), common.GetLine())
	}

	bytes, err := json.Marshal(proposal)
	if err != nil { // Return error: Can't marshal json
		return nil, fmt.Errorf("%s %s %s", common.ResCodeDict[common.ERR3], err.Error(), common.GetLine())
	}
	temp := ""
	result = &temp
	*result = string(bytes)

	return result, nil
}

// GetAllProposal ...
func (sah *ProposalHandler) GetAllProposal(stub shim.ChaincodeStubInterface) (result *string, err error) {
	res := util.GetAllData(stub, new(model.Proposal), model.ProposalTable)
	if res.Status == 200 {
		return &res.Message, nil
	}
	return nil, fmt.Errorf("%s %s %s", common.ResCodeDict[common.ERR4], err.Error(), common.GetLine())
}

// GetProposalByID ...
func (sah *ProposalHandler) GetProposalByID(stub shim.ChaincodeStubInterface, proposalID string) (result *string, err error) {
	common.Logger.Debugf("Input-data sent to GetProposalByID func: %+v\n", proposalID)

	res := util.GetDataByID(stub, proposalID, new(model.Proposal), model.ProposalTable)
	if res.Status == 200 {
		return &res.Message, nil
	} else {
		return nil, fmt.Errorf("%s %s %s", common.ResCodeDict[common.ERR4], err.Error(), common.GetLine())
	}
}

// GetPendingProposalBySuperAdminID ...
func (sah *ProposalHandler) GetPendingProposalBySuperAdminID(stub shim.ChaincodeStubInterface, superAdminID string) (result *string, err error) {
	common.Logger.Debugf("Input-data sent to GetPendingProposalBySuperAdminID func: %+v\n", superAdminID)

	var proposalList []model.Proposal

	queryStr := fmt.Sprintf("{\"selector\": {\"_id\": {\"$regex\": \"%s\"},\"$or\": [{\"Status\": \"Pending\"},{\"Status\": \"Approved\"}]}}", model.ProposalTable)
	resultsIterator, err := stub.GetQueryResult(queryStr)
	if err != nil {
		return nil, fmt.Errorf("%s %s %s", common.ResCodeDict[common.ERR4], err.Error(), common.GetLine())
	}
	defer resultsIterator.Close()
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("%s %s %s", common.ResCodeDict[common.ERR4], err.Error(), common.GetLine())
		}

		proposal := new(model.Proposal)
		err = json.Unmarshal(queryResponse.Value, proposal)
		if err != nil { // Convert JSON error
			return nil, fmt.Errorf("%s %s %s", common.ResCodeDict[common.ERR3], err.Error(), common.GetLine())
		}
		proposalList = append(proposalList, *proposal)
	}

	for i := len(proposalList) - 1; i >= 0; i-- {
		proposal := proposalList[i]
		rs, err := hUtil.GetByTwoColumns(stub, model.ApprovalTable, "ProposalID", fmt.Sprintf("\"%s\"", proposal.ProposalID), "ApproverID", fmt.Sprintf("\"%s\"", superAdminID))
		if err != nil {
			return nil, fmt.Errorf("%s %s %s", common.ResCodeDict[common.ERR4], err.Error(), common.GetLine())
		}
		if rs.HasNext() {
			proposalList[i] = proposalList[len(proposalList)-1]
			proposalList = proposalList[:len(proposalList)-1]
		}
	}

	bytes, err := json.Marshal(proposalList)
	if err != nil { // Return error: Can't marshal json
		return nil, fmt.Errorf("%s %s %s", common.ResCodeDict[common.ERR3], err.Error(), common.GetLine())
	}
	temp := ""
	result = &temp
	*result = string(bytes)

	return result, nil
}

//UpdateProposal ...
func (sah *ProposalHandler) UpdateProposal(stub shim.ChaincodeStubInterface, proposalStr string) (result *string, err error) {
	common.Logger.Debugf("Input-data sent to UpdateProposal func: %+v\n", proposalStr)

	newProposal := new(model.Proposal)
	err = json.Unmarshal([]byte(proposalStr), newProposal)
	if err != nil { // Return error: Can't unmarshal json
		return nil, fmt.Errorf("%s %s %s", common.ResCodeDict[common.ERR3], err.Error(), common.GetLine())
	}

	if len(newProposal.ProposalID) == 0 {
		return nil, fmt.Errorf("%s %s", "This ApprovalID can't be empty", common.GetLine())
	}

	//get proposal information
	rawProposal, err := util.Getdatabyid(stub, newProposal.ProposalID, model.ProposalTable)
	if err != nil {
		return nil, fmt.Errorf("%s %s %s", common.ResCodeDict[common.ERR4], err.Error(), common.GetLine())
	}

	proposal := new(model.Proposal)
	mapstructure.Decode(rawProposal, proposal)

	// Filter fields needed to update
	newProposalValue := reflect.ValueOf(newProposal).Elem()
	proposalValue := reflect.ValueOf(proposal).Elem()
	for i := 0; i < newProposalValue.NumField(); i++ {
		fieldName := newProposalValue.Type().Field(i).Name
		if len(newProposalValue.Field(i).String()) > 0 {
			field := proposalValue.FieldByName(fieldName)
			if field.CanSet() {
				field.SetString(newProposalValue.Field(i).String())
			}
		}
	}

	err = util.Changeinfo(stub, model.ProposalTable, []string{proposal.ProposalID}, proposal)
	if err != nil { // Return error: Fail to Update data
		return nil, fmt.Errorf("%s %s %s", common.ResCodeDict[common.ERR5], err.Error(), common.GetLine())
	}

	bytes, err := json.Marshal(proposal)
	if err != nil { // Return error: Can't marshal json
		return nil, fmt.Errorf("%s %s %s", common.ResCodeDict[common.ERR3], err.Error(), common.GetLine())
	}
	temp := ""
	result = &temp
	*result = string(bytes)

	return result, nil
}

//CommitProposal ...
func (sah *ProposalHandler) CommitProposal(stub shim.ChaincodeStubInterface, proposalID string, criteria int) (result *string, err error) {
	common.Logger.Debugf("Input-data sent to CommitProposal func: %+v\n", proposalID)

	resIterator, err := hUtil.GetByOneColumn(stub, model.ApprovalTable, "ProposalID", fmt.Sprintf("\"%s\"", proposalID))
	if err != nil {
		return nil, fmt.Errorf("%s %s %s", common.ResCodeDict[common.ERR4], err.Error(), common.GetLine())
	}
	defer resIterator.Close()

	count := 0
	for resIterator.HasNext() {
		_, err := resIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("%s %s %s", common.ResCodeDict[common.ERR4], err.Error(), common.GetLine())
		}
		count++
	}
	if count < criteria {
		return nil, fmt.Errorf("%s %s", "not enought signature", common.GetLine())
	}

	// Get proposal information
	rawProposal, err := util.Getdatabyid(stub, proposalID, model.ProposalTable)
	if err != nil {
		return nil, fmt.Errorf("%s %s %s", common.ResCodeDict[common.ERR4], err.Error(), common.GetLine())
	}

	proposal := new(model.Proposal)
	mapstructure.Decode(rawProposal, proposal)

	proposal.Status = "Committed"
	timestamp, err := stub.GetTxTimestamp()
	if err != nil {
		return nil, fmt.Errorf("%s %s %s", common.ResCodeDict[common.ERR4], err.Error(), common.GetLine())
	}
	updatedTime := time.Unix(timestamp.Seconds, 0)
	proposal.UpdatedAt = updatedTime.String()

	err = util.Changeinfo(stub, model.ProposalTable, []string{proposal.ProposalID}, proposal)
	if err != nil { // Return error: Fail to Update data
		return nil, fmt.Errorf("%s %s %s", common.ResCodeDict[common.ERR5], err.Error(), common.GetLine())
	}

	bytes, err := json.Marshal(proposal)
	if err != nil { // Return error: Can't marshal json
		return nil, fmt.Errorf("%s %s %s", common.ResCodeDict[common.ERR3], err.Error(), common.GetLine())
	}
	temp := ""
	result = &temp
	*result = string(bytes)

	return result, nil
}

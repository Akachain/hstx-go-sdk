package main

import (
	"fmt"
	"strconv"

	"github.com/Akachain/akc-go-sdk/common"
	hdl "github.com/Akachain/hstx-go-sdk/handler"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// Chaincode struct
type Chaincode struct {
}

var handler = new(hdl.Handler)

// Init method is called when the Chain code" is instantiated by the blockchain network
func (s *Chaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

// Invoke method is called as a result of an application request to run the chain code
// The calling application program has also specified the particular smart contract function to be called, with arguments
func (s *Chaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	common.Logger.Info("########### Hstx Invoke ###########")

	handler.InitHandler()

	// Retrieve the requested Smart Contract function and arguments
	functionName, args := stub.GetFunctionAndParameters()

	router := map[string]func(shim.ChaincodeStubInterface, []string) pb.Response{
		"CreateSuperAdmin": createSuperAdmin,
		"CreateAdmin": createAdmin,
		"CreateProposal":   createProposal,
		"CreateApproval":   createApproval,
		"CommitProposal":   commitProposal,
		// "UpdateSuperAdmin": handler.SuperAdminHandler.UpdateSuperAdmin,
		// "UpdateAdmin":      handler.AdminHandler.UpdateAdmin,
		// "UpdateProposal":   handler.ProposalHandler.UpdateProposal,
		// "UpdateApproval":   handler.ApprovalHandler.UpdateApproval,
	}

	invokeFunc := router[functionName]
	if invokeFunc != nil {
		return invokeFunc(stub, args)
	}
	return s.Query(stub)
}

// Query callback representing the query of a chaincode
func (s *Chaincode) Query(stub shim.ChaincodeStubInterface) pb.Response {
	common.Logger.Info("########### Hstx Query ###########")

	// Retrieve the requested Smart Contract function and arguments
	functionName, args := stub.GetFunctionAndParameters()

	router := map[string]func(shim.ChaincodeStubInterface, []string) pb.Response{
		// "GetAllSuperAdmin":                 handler.SuperAdminHandler.GetAllSuperAdmin,
		// "GetSuperAdminByID":                handler.SuperAdminHandler.GetSuperAdminByID,
		// "GetAllAdmin":                      handler.AdminHandler.GetAllAdmin,
		// "GetAdminByID":                     handler.AdminHandler.GetAdminByID,
		// "GetAllProposal":                   handler.ProposalHandler.GetAllProposal,
		// "GetProposalByID":                  handler.ProposalHandler.GetProposalByID,
		// "GetPendingProposalBySuperAdminID": handler.ProposalHandler.GetPendingProposalBySuperAdminID,
		// "GetAllApproval":                   handler.ApprovalHandler.GetAllApproval,
		// "GetApprovalByID":                  handler.ApprovalHandler.GetApprovalByID,
	}

	queryFunc := router[functionName]
	return queryFunc(stub, args)
}

// createSuperAdmin
func createSuperAdmin(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	superAdminStr := args[0]

	created, err := handler.SuperAdminHandler.CreateSuperAdmin(stub, superAdminStr)
	if err != nil {
		// Returning error: Can't create data
		return common.RespondError(common.ResponseError{
			ResCode: common.ERR4,
			Msg:     fmt.Sprintf("%s %s %s", common.ResCodeDict[common.ERR4], err.Error(), common.GetLine()),
		})
	}

	// Returning success with payload is created data
	resSuc := common.ResponseSuccess{
		ResCode: common.SUCCESS,
		Msg:     common.ResCodeDict[common.SUCCESS],
		Payload: *created,
	}
	return common.RespondSuccess(resSuc)
}

// createAdmin
func createAdmin(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	adminStr := args[0]

	created, err := handler.AdminHandler.CreateAdmin(stub, adminStr)
	if err != nil {
		// Returning error: Can't create data
		return common.RespondError(common.ResponseError{
			ResCode: common.ERR4,
			Msg:     fmt.Sprintf("%s %s %s", common.ResCodeDict[common.ERR4], err.Error(), common.GetLine()),
		})
	}

	// Returning success with payload is created data
	resSuc := common.ResponseSuccess{
		ResCode: common.SUCCESS,
		Msg:     common.ResCodeDict[common.SUCCESS],
		Payload: *created,
	}
	return common.RespondSuccess(resSuc)
}

// createProposal
func createProposal(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	proposalStr := args[0]

	created, err := handler.ProposalHandler.CreateProposal(stub, proposalStr)
	if err != nil {
		// Returning error: Can't create data
		return common.RespondError(common.ResponseError{
			ResCode: common.ERR4,
			Msg:     fmt.Sprintf("%s %s %s", common.ResCodeDict[common.ERR4], err.Error(), common.GetLine()),
		})
	}

	// Returning success with payload is created data
	resSuc := common.ResponseSuccess{
		ResCode: common.SUCCESS,
		Msg:     common.ResCodeDict[common.SUCCESS],
		Payload: *created,
	}
	return common.RespondSuccess(resSuc)
}

// createApproval
func createApproval(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	approvalStr := args[0]
	criteria, err := strconv.Atoi(args[1])
	if err != nil {
		// Returning error: Can't convert data
		return common.RespondError(common.ResponseError{
			ResCode: common.ERR3,
			Msg:     fmt.Sprintf("%s %s %s", common.ResCodeDict[common.ERR3], err.Error(), common.GetLine()),
		})
	}

	created, err := handler.ApprovalHandler.CreateApproval(stub, approvalStr, criteria)
	if err != nil {
		// Returning error: Can't create data
		return common.RespondError(common.ResponseError{
			ResCode: common.ERR4,
			Msg:     fmt.Sprintf("%s %s %s", common.ResCodeDict[common.ERR4], err.Error(), common.GetLine()),
		})
	}

	// Returning success with payload is created data
	resSuc := common.ResponseSuccess{
		ResCode: common.SUCCESS,
		Msg:     common.ResCodeDict[common.SUCCESS],
		Payload: *created,
	}
	return common.RespondSuccess(resSuc)
}

// commitProposal
func commitProposal(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	approvalStr := args[0]
	criteria, err := strconv.Atoi(args[1])
	if err != nil {
		// Returning error: Can't convert data
		return common.RespondError(common.ResponseError{
			ResCode: common.ERR3,
			Msg:     fmt.Sprintf("%s %s %s", common.ResCodeDict[common.ERR3], err.Error(), common.GetLine()),
		})
	}

	created, err := handler.ProposalHandler.CommitProposal(stub, approvalStr, criteria)
	if err != nil {
		// Returning error: Can't create data
		return common.RespondError(common.ResponseError{
			ResCode: common.ERR4,
			Msg:     fmt.Sprintf("%s %s %s", common.ResCodeDict[common.ERR4], err.Error(), common.GetLine()),
		})
	}

	// Returning success with payload is created data
	resSuc := common.ResponseSuccess{
		ResCode: common.SUCCESS,
		Msg:     common.ResCodeDict[common.SUCCESS],
		Payload: *created,
	}
	return common.RespondSuccess(resSuc)
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {
	// Create a new Chain code
	err := shim.Start(new(Chaincode))
	if err != nil {
		fmt.Printf("Error creating new Chain code: %s", err)
	}
}

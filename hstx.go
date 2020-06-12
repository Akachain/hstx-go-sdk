package main

import (
	"fmt"

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
		// "CreateSuperAdmin": handler.SuperAdminHandler.CreateSuperAdmin,
		// "UpdateSuperAdmin": handler.SuperAdminHandler.UpdateSuperAdmin,
		// "CreateAdmin":      handler.AdminHandler.CreateAdmin,
		// "UpdateAdmin":      handler.AdminHandler.UpdateAdmin,
		// "CreateProposal":   handler.ProposalHandler.CreateProposal,
		// "UpdateProposal":   handler.ProposalHandler.UpdateProposal,
		// "CommitProposal":   handler.ProposalHandler.CommitProposal,
		// "CreateApproval":   handler.ApprovalHandler.CreateApproval,
		// "UpdateApproval":   handler.ApprovalHandler.UpdateApproval,
	}

	invokeFunc := router[functionName]
	if invokeFunc != nil {
		return invokeFunc(stub, args)
	} else {
		return s.Query(stub)
	}
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

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {
	// Create a new Chain code
	err := shim.Start(new(Chaincode))
	if err != nil {
		fmt.Printf("Error creating new Chain code: %s", err)
	}
}

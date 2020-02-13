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
		"CreateSuperAdmin": handler.SuperAdminHanler.CreateSuperAdmin,
		"UpdateSuperAdmin": handler.SuperAdminHanler.UpdateSuperAdmin,
		"CreateAdmin":      handler.AdminHanler.CreateAdmin,
		"UpdateAdmin":      handler.AdminHanler.UpdateAdmin,
		"CreateProposal":   handler.ProposalHanler.CreateProposal,
		"UpdateProposal":   handler.ProposalHanler.UpdateProposal,
		"CommitProposal":   handler.ProposalHanler.CommitProposal,
		"CreateApproval":   handler.ApprovalHanler.CreateApproval,
		"UpdateApproval":   handler.ApprovalHanler.UpdateApproval,
	}

	invokeFunc := router[functionName]
	return invokeFunc(stub, args)
}

// Query callback representing the query of a chaincode
func (s *Chaincode) Query(stub shim.ChaincodeStubInterface) pb.Response {
	common.Logger.Info("########### Hstx Query ###########")

	// Retrieve the requested Smart Contract function and arguments
	functionName, args := stub.GetFunctionAndParameters()

	router := map[string]func(shim.ChaincodeStubInterface, []string) pb.Response{
		"GetAllSuperAdmin":                 handler.SuperAdminHanler.GetAllSuperAdmin,
		"GetSuperAdminByID":                handler.SuperAdminHanler.GetSuperAdminByID,
		"GetAllAdmin":                      handler.AdminHanler.GetAllAdmin,
		"GetAdminByID":                     handler.AdminHanler.GetAdminByID,
		"GetAllProposal":                   handler.ProposalHanler.GetAllProposal,
		"GetProposalByID":                  handler.ProposalHanler.GetProposalByID,
		"GetPendingProposalBySuperAdminID": handler.ProposalHanler.GetPendingProposalBySuperAdminID,
		"GetAllApproval":                   handler.ApprovalHanler.GetAllApproval,
		"GetApprovalByID":                  handler.ApprovalHanler.GetApprovalByID,
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

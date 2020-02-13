package utils

import (
	"fmt"

	"github.com/Akachain/akc-go-sdk/common"
	"github.com/hyperledger/fabric/core/chaincode/lib/cid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// GetCertID func to get Certtificate ID of current user
func GetCertID(stub shim.ChaincodeStubInterface) (*string, error) {
	id, err := cid.GetID(stub)
	if err != nil {
		return nil, fmt.Errorf("%s %s %s", "Can't get Certificate ID.", err.Error(), common.GetLine())
	}
	return &id, nil
}

// GetMSPID func to get MSP ID of current user
func GetMSPID(stub shim.ChaincodeStubInterface) (*string, error) {
	mspid, err := cid.GetMSPID(stub)
	if err != nil {
		return nil, fmt.Errorf("%s %s %s", "Can't get Certificate MSP ID.", err.Error(), common.GetLine())
	}
	return &mspid
}

// GetRole func to get 'role' attribute saved in current user's certificate
func GetRole(stub shim.ChaincodeStubInterface) (*string, error) {
	return GetAttributeValue(stub, "role")
}

// GetAttributeValue func to get a attribute saved in current user's certificate
func GetAttributeValue(stub shim.ChaincodeStubInterface, attrName string) (*string, error) {
	val, ok, err := cid.GetAttributeValue(stub, attrName)
	if err != nil {
		return nil, fmt.Errorf("%s %s %s", "Can't get attribute in Certificate", err.Error(), common.GetLine())
	}
	if !ok {
		return nil, fmt.Errorf("%s %s %s", "Can't get attribute in Certificate", "This attr doesn't exist", common.GetLine())
	}
	return &val
}

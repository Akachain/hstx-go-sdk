package utils

import (
	"fmt"
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"

	"github.com/Akachain/akc-go-sdk/common"
	"github.com/hyperledger/fabric/bccsp/utils"
	"github.com/hyperledger/fabric/core/chaincode/lib/cid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// GetCertID func to get Certtificate ID of current user
func GetCertID(stub shim.ChaincodeStubInterface) (*string, error) {
	id, err := cid.GetID(stub)
	if err != nil {
		return nil, fmt.Errorf("Can't get Certificate ID. Cause: %s %s", err.Error(), common.GetLine())
	}
	return &id, nil
}

// GetMSPID func to get MSP ID of current user
func GetMSPID(stub shim.ChaincodeStubInterface) (*string, error) {
	mspid, err := cid.GetMSPID(stub)
	if err != nil {
		return nil, fmt.Errorf("Can't get MSP ID. Cause: %s %s", err.Error(), common.GetLine())
	}
	return &mspid, nil
}

// GetRole func to get 'role' attribute saved in current user's certificate
func GetRole(stub shim.ChaincodeStubInterface) (*string, error) {
	return GetAttributeValue(stub, "role")
}

// GetAttributeValue func to get a attribute saved in current user's certificate
func GetAttributeValue(stub shim.ChaincodeStubInterface, attrName string) (*string, error) {
	val, ok, err := cid.GetAttributeValue(stub, attrName)
	if err != nil {
		return nil, fmt.Errorf("Can't get attribute '%s' in the Certificate. Cause: %s %s", attrName, err.Error(), common.GetLine())
	}
	if !ok {
		return nil, fmt.Errorf("Can't get attribute '%s' in the Certificate. Cause: %s", attrName, common.GetLine())
	}
	return &val, nil
}

// GetByOneColumn func to get information
func GetByOneColumn(stub shim.ChaincodeStubInterface, table string, column string, value interface{}) (resultsIterator shim.StateQueryIteratorInterface, err error) {
	queryString := fmt.Sprintf("{\"selector\": {\"_id\": {\"$regex\": \"%s\"},\"%s\": %v}}", table, column, value)
	common.Logger.Info(queryString)
	resultsIterator, err = stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	return resultsIterator, nil
}

// GetByTwoColumns func to get information
func GetByTwoColumns(stub shim.ChaincodeStubInterface, table string, column1 string, value1 interface{}, column2 string, value2 interface{}) (resultsIterator shim.StateQueryIteratorInterface, err error) {
	queryString := fmt.Sprintf("{\"selector\": {\"_id\": {\"$regex\": \"%s\"},\"%s\": %v, \"%s\": %v}}", table, column1, value1, column2, value2)
	common.Logger.Info(queryString)
	resultsIterator, err = stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	return resultsIterator, nil
}

// verifySignature func base on approver's public key, signature and message that was singed
func verifySignature(stub shim.ChaincodeStubInterface, publicKey string, signature string, message string) error {
	if len(publicKey) == 0 {
		return fmt.Errorf("publicKey is empty %s", common.GetLine())
	}
	if len(signature) == 0 {
		return fmt.Errorf("signature is empty %s", common.GetLine())
	}

	// Start verify
	pkBytes := []byte(publicKey)
	pkBlock, _ := pem.Decode(pkBytes)
	if pkBlock == nil {
		return fmt.Errorf("Can't decode public key %s", common.GetLine())
	}

	rawPk, err := x509.ParsePKIXPublicKey(pkBlock.Bytes)
	if err != nil {
		return err
	}

	pk := rawPk.(*ecdsa.PublicKey)

	// SIGNATURE
	signaturebyte, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return err
	}

	R, S, err := utils.UnmarshalECDSASignature(signaturebyte)
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
	checksign := ecdsa.Verify(pk, hashData, R, S)

	if checksign {
		return nil
	}
	return errors.New("Verify failed")
}
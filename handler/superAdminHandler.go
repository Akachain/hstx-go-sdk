package handler

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/Akachain/akc-go-sdk/common"
	"github.com/Akachain/akc-go-sdk/util"
	"github.com/Akachain/hstx-go-sdk/model"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/mitchellh/mapstructure"
)

// SuperAdminHandler ...
type SuperAdminHandler struct{}

// CreateSuperAdmin ...
func (sah *SuperAdminHandler) CreateSuperAdmin(stub shim.ChaincodeStubInterface, superAdminStr string) (result *string, err error) {
	common.Logger.Debugf("Input-data sent to CreateSuperAdmin func: %+v\n", superAdminStr)

	superAdmin := new(model.SuperAdmin)
	err = json.Unmarshal([]byte(superAdminStr), superAdmin)
	if err != nil { // Return error: Can't unmarshal json
		return nil, fmt.Errorf("%s %s %s", common.ResCodeDict[common.ERR3], err.Error(), common.GetLine())
	}

	superAdmin.Status = "Active"

	common.Logger.Infof("Create SuperAdmin: %+v\n", superAdmin)
	err = util.Createdata(stub, model.SuperAdminTable, []string{superAdmin.SuperAdminID}, &superAdmin)
	if err != nil { // Return error: Fail to insert data
		return nil, fmt.Errorf("%s %s %s", common.ResCodeDict[common.ERR5], err.Error(), common.GetLine())
	}

	bytes, err := json.Marshal(superAdmin)
	if err != nil { // Return error: Can't marshal json
		return nil, fmt.Errorf("%s %s %s", common.ResCodeDict[common.ERR3], err.Error(), common.GetLine())
	}
	temp := ""
	result = &temp
	*result = string(bytes)

	return result, nil
}

// GetAllSuperAdmin ...
func (sah *SuperAdminHandler) GetAllSuperAdmin(stub shim.ChaincodeStubInterface) (result *string, err error) {
	res := util.GetAllData(stub, new(model.SuperAdmin), model.SuperAdminTable)
	if res.Status == 200 {
		return &res.Message, nil
	}
	return nil, fmt.Errorf("%s %s %s", common.ResCodeDict[common.ERR4], err.Error(), common.GetLine())
}

// GetSuperAdminByID ...
func (sah *SuperAdminHandler) GetSuperAdminByID(stub shim.ChaincodeStubInterface, superAdminID string) (result *string, err error) {
	common.Logger.Debugf("Input-data sent to SuperAdminHandler func: %+v\n", superAdminID)

	res := util.GetDataByID(stub, superAdminID, new(model.SuperAdmin), model.SuperAdminTable)
	if res.Status == 200 {
		return &res.Message, nil
	} else {
		return nil, fmt.Errorf("%s %s %s", common.ResCodeDict[common.ERR4], err.Error(), common.GetLine())
	}
}

//UpdateSuperAdmin ...
func (sah *SuperAdminHandler) UpdateSuperAdmin(stub shim.ChaincodeStubInterface, superAdminStr string) (result *string, err error) {
	common.Logger.Debugf("Input-data sent to UpdateSuperAdmin func: %+v\n", superAdminStr)

	newSuperAdmin := new(model.SuperAdmin)
	err = json.Unmarshal([]byte(superAdminStr), newSuperAdmin)
	if err != nil { // Return error: Can't unmarshal json
		return nil, fmt.Errorf("%s %s %s", common.ResCodeDict[common.ERR3], err.Error(), common.GetLine())
	}

	if len(newSuperAdmin.SuperAdminID) == 0 {
		return nil, fmt.Errorf("%s %s", "This ApprovalID can't be empty", common.GetLine())
	}

	// Get superAdmin information
	rawSuperAdmin, err := util.Getdatabyid(stub, newSuperAdmin.SuperAdminID, model.SuperAdminTable)
	if err != nil {
		return nil, fmt.Errorf("%s %s %s", common.ResCodeDict[common.ERR4], err.Error(), common.GetLine())
	}

	superAdmin := new(model.SuperAdmin)
	mapstructure.Decode(rawSuperAdmin, superAdmin)

	// Filter fields needed to update
	newSuperAdminValue := reflect.ValueOf(newSuperAdmin).Elem()
	superAdminValue := reflect.ValueOf(superAdmin).Elem()
	for i := 0; i < newSuperAdminValue.NumField(); i++ {
		fieldName := newSuperAdminValue.Type().Field(i).Name
		if len(newSuperAdminValue.Field(i).String()) > 0 {
			field := superAdminValue.FieldByName(fieldName)
			if field.CanSet() {
				field.SetString(newSuperAdminValue.Field(i).String())
			}
		}
	}

	err = util.Changeinfo(stub, model.SuperAdminTable, []string{superAdmin.SuperAdminID}, superAdmin)
	if err != nil {
		return nil, fmt.Errorf("%s %s %s", common.ResCodeDict[common.ERR5], err.Error(), common.GetLine())
	}

	bytes, err := json.Marshal(superAdmin)
	if err != nil { // Return error: Can't marshal json
		return nil, fmt.Errorf("%s %s %s", common.ResCodeDict[common.ERR3], err.Error(), common.GetLine())
	}
	temp := ""
	result = &temp
	*result = string(bytes)

	return result, nil
}

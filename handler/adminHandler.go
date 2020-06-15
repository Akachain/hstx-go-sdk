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

// AdminHandler ...
type AdminHandler struct{}

// CreateAdmin ...
func (sah *AdminHandler) CreateAdmin(stub shim.ChaincodeStubInterface, adminStr string) (result *string, err error) {
	common.Logger.Debugf("Input-data sent to CreateAdmin func: %+v\n", adminStr)

	admin := new(model.Admin)
	err = json.Unmarshal([]byte(adminStr), admin)
	if err != nil { // Return error: Can't unmarshal json
		return nil, fmt.Errorf("%s %s %s", common.ResCodeDict[common.ERR3], err.Error(), common.GetLine())
	}

	admin.AdminID = stub.GetTxID()
	admin.Status = "Active"

	common.Logger.Infof("Create Admin: %+v\n", admin)
	err = util.Createdata(stub, model.AdminTable, []string{admin.AdminID}, &admin)
	if err != nil { // Return error: Fail to insert data
		return nil, fmt.Errorf("%s %s %s", common.ResCodeDict[common.ERR5], err.Error(), common.GetLine())
	}

	bytes, err := json.Marshal(admin)
	if err != nil { // Return error: Can't marshal json
		return nil, fmt.Errorf("%s %s %s", common.ResCodeDict[common.ERR3], err.Error(), common.GetLine())
	}
	temp := ""
	result = &temp
	*result = string(bytes)

	return result, nil
}

// GetAllAdmin ...
func (sah *AdminHandler) GetAllAdmin(stub shim.ChaincodeStubInterface) (result *string, err error) {
	res := util.GetAllData(stub, new(model.Admin), model.AdminTable)
	if res.Status == 200 {
		return &res.Message, nil
	}
	return nil, fmt.Errorf("%s %s %s", common.ResCodeDict[common.ERR4], err.Error(), common.GetLine())
}

// GetAdminByID ...
func (sah *AdminHandler) GetAdminByID(stub shim.ChaincodeStubInterface, adminID string) (result *string, err error) {
	common.Logger.Debugf("Input-data sent to GetAdminByID func: %+v\n", adminID)

	res := util.GetDataByID(stub, adminID, new(model.Admin), model.AdminTable)
	if res.Status == 200 {
		return &res.Message, nil
	} else {
		return nil, fmt.Errorf("%s %s %s", common.ResCodeDict[common.ERR4], err.Error(), common.GetLine())
	}
}

//UpdateAdmin ...
func (sah *AdminHandler) UpdateAdmin(stub shim.ChaincodeStubInterface, adminStr string) (result *string, err error) {
	common.Logger.Debugf("Input-data sent to UpdateAdmin func: %+v\n", adminStr)

	newAdmin := new(model.Admin)
	err = json.Unmarshal([]byte(adminStr), newAdmin)
	if err != nil { // Return error: Can't unmarshal json
		return nil, fmt.Errorf("%s %s %s", common.ResCodeDict[common.ERR3], err.Error(), common.GetLine())
	}

	if len(newAdmin.AdminID) == 0 {
		return nil, fmt.Errorf("%s %s", "This ApprovalID can't be empty", common.GetLine())
	}

	// Get admin information
	rawAdmin, err := util.Getdatabyid(stub, newAdmin.AdminID, model.AdminTable)
	if err != nil {
		return nil, fmt.Errorf("%s %s %s", common.ResCodeDict[common.ERR4], err.Error(), common.GetLine())
	}

	admin := new(model.Admin)
	mapstructure.Decode(rawAdmin, admin)

	// Filter fields needed to update
	newAdminValue := reflect.ValueOf(newAdmin).Elem()
	adminValue := reflect.ValueOf(admin).Elem()
	for i := 0; i < newAdminValue.NumField(); i++ {
		fieldName := newAdminValue.Type().Field(i).Name
		if len(newAdminValue.Field(i).String()) > 0 {
			field := adminValue.FieldByName(fieldName)
			if field.CanSet() {
				field.SetString(newAdminValue.Field(i).String())
			}
		}
	}

	err = util.Changeinfo(stub, model.AdminTable, []string{admin.AdminID}, admin)
	if err != nil { // Return error: Fail to Update data
		return nil, fmt.Errorf("%s %s %s", common.ResCodeDict[common.ERR5], err.Error(), common.GetLine())
	}

	bytes, err := json.Marshal(admin)
	if err != nil { // Return error: Can't marshal json
		return nil, fmt.Errorf("%s %s %s", common.ResCodeDict[common.ERR3], err.Error(), common.GetLine())
	}
	temp := ""
	result = &temp
	*result = string(bytes)

	return result, nil
}

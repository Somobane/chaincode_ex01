/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/

package main

import (
	"errors"
	"fmt"
	"strconv"
	"time"
	"math/rand"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	
)


// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

//iQOS Packaging Line Changes starts ----------------------------------------------------------------
// Shipping Case will have:
	// DeviceId : AssemblyID (Both Primary and Secondary)
	// DeviceType: Type of Device (Charger or Holder) (Both Primary and Secondary)
	// CaseId: UniqueID to be genrated
	// TestResult: Overall Test Status
	// TestDate: Date and time Tested
	// PackagingLineState: Overall status of the Packaging Line
	// PakagingDate: Date when package was packed
type ShippingCase struct{	
	CaseId int `json:"caseId"`
	PrimaryDeviceId string `json:"primaryDeviceId"`
	PrimaryDeviceType string `json:"primaryDeviceType"`
	SecondaryDeviceId string `json:"secondaryDeviceId"`
	SecondaryDeviceType string `json:"secondaryDeviceType"`
	PackagingDate string `json:"packagingDate"`
	TestResult string `json:"testResult"`
	TestDate string `json:"testDate"`
	PackagingLineState string `json:"packagingLineState"`
	}

// Retrieve Package line Status-- not implemented
type RetreivePackageLineStatus struct{	
	PackagingLineState string `json:"packagingLineState"`
}
//iQOS Changes ends ----------------------------------------------------------------


func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Printf("Init called, initializing chaincode")
	
	// Create application Table
	err1 = stub.CreateTable("ShippingCase", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "caseId", Type: shim.ColumnDefinition_INT, Key: true},
		&shim.ColumnDefinition{Name: "primaryDeviceId", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "primaryDeviceType", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "secondaryDeviceId", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "secondaryDeviceType", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "packagingDate", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "testResult", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "testDate", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "packagingLineState", Type: shim.ColumnDefinition_STRING, Key: false},

	})
	if err1 != nil {
		return nil, errors.New("Failed creating ShippingCase.")
	}


	return nil, nil
}


func (t *SimpleChaincode) createShippingCase(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

if len(args) != 4 {
			return nil, fmt.Errorf("Incorrect number of arguments. Expecting 9. Got: %d.", len(args))
		}
		
		var r int;
		r := rand.New(rand.NewSource(99));

		var t string;
		t := time.Now().Day() + "-" + time.Now().Month() + "-" + time.Now().Year();

		caseId := r.Int31()
		primaryDeviceId:=args[0]
		primaryDeviceType :=args[1]
		secondaryDeviceId :=args[2]
		secondaryDeviceType:=args[3]
		packagingDate:= t
		testResult:= "Within Range"
		testDate:= t;
		packagingLineState:= "QA"
	
		// Insert a row
		ok, err := stub.InsertRow("ShippingCase", shim.Row{
			Columns: []*shim.Column{
				&shim.Column{Value: &shim.Column_Int_{Int_: caseId}},
				&shim.Column{Value: &shim.Column_String_{String_: primaryDeviceId}},
				&shim.Column{Value: &shim.Column_String_{String_: primaryDeviceType}},
				&shim.Column{Value: &shim.Column_String_{String_: secondaryDeviceId}},
				&shim.Column{Value: &shim.Column_String_{String_: secondaryDeviceType}},
				&shim.Column{Value: &shim.Column_String_{String_: packagingDate}},
				&shim.Column{Value: &shim.Column_String_{String_: testResult}},
				&shim.Column{Value: &shim.Column_String_{String_: testDate}},
				&shim.Column{Value: &shim.Column_String_{String_: packagingLineState}},
				}})

		if err != nil {
			return nil, err 
		}
		if !ok && err == nil {
			return nil, errors.New("Row already exists.")
		}
			
		return nil, nil

}


// Invoke callback representing the invocation of a chaincode
// This chaincode will manage two accounts A and B and will transfer X units from A to B upon invoke
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Printf("Invoke called, determining function")
	
	// Handle different functions
	if function == "init" {
		fmt.Printf("Function is init")
		return t.Init(stub, function, args)
	} else if function == "createShippingCase" {
		fmt.Printf("Function is createShippingCase")
		return t.createShippingCase(stub, args)
	} 
	return nil, errors.New("Received unknown function invocation")
}

func (t* SimpleChaincode) Run(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Printf("Run called, passing through to Invoke (same function)")
	
	// Handle different functions
	if function == "createShippingCase" {
		fmt.Printf("Function is createShippingCase")
		return t.createShippingCase(stub, args)
	}   else if function == "init" {
		fmt.Printf("Function is init")
		return t.Init(stub, function, args)
	}

	return nil, errors.New("Received unknown function invocation")
}

//get the Shipping Case ID
func (t *SimpleChaincode) getCaseID(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting Shipping Case ID to query")
	}

	caseId := args[0]
	

	// Get the row pertaining to this caseId
	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_Int_{Int_: caseId}}
	columns = append(columns, col1)

	row, err := stub.GetRow("ShippingCase", columns)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get the data for the caseId " + caseId + "\"}"
		return nil, errors.New(jsonResp)
	}

	// GetRows returns empty message if key does not exist
	if len(row.Columns) == 0 {
		jsonResp := "{\"Error\":\"Failed to get the data for the caseId " + caseId + "\"}"
		return nil, errors.New(jsonResp)
	}

	 mapB, _ := json.Marshal(row)
    fmt.Println(string(mapB))
	
	return mapB, nil

}


// query queries the chaincode
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Printf("Query called, determining function")
if function == "getCaseID" { 
		t := SimpleChaincode{}
		return t.getCaseID(stub, args)
	}
	
	return nil, nil
}


func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

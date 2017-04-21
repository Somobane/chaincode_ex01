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

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

//iQOS Changes starts ----------------------------------------------------------------
// Assembly comprises of one Serial ID and multiple Batch IDs
type AssemblyLine struct{	
	AssemblyId string `json:"assemblyId"`
	SerialId string `json:"serialId"`
	OriginalFilamentBatchId string `json:"originalFilamentBatchId"`
	OriginalLedBatchId string `json:"originalLedBatchId"`
	OriginalCircuitBoardBatchId string `json:"originalCircuitBoardBatchId"`
	OriginalWireBatchId string `json:"originalWireBatchId"`
	OriginalCasingBatchId string `json:"originalCasingBatchId"`
	OriginalAdaptorBatchId string `json:"originalAdaptorBatchId"`
	OriginalStickPodBatchId string `json:"originalStickPodBatchId"`
	AssemblyStatus string `json:"assemblyStatus"`
	}

// GetAssemblyLineStatus is for storing retreived Assembly Line Status
type GetAssemblyLineStatus struct{	
	AssemblyLineStatus string `json:"assemblyLineStatus"`
}
//iQOS Changes ends ----------------------------------------------------------------


func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Printf("Init called, initializing chaincode")
	
	var A, B string    // Entities
	var Aval, Bval int // Asset holdings
	var err error
    var err1 error

	if len(args) != 4 {
		return nil, errors.New("Incorrect number of arguments. Expecting 4")
	}

	// Initialize the chaincode
	A = args[0]
	Aval, err = strconv.Atoi(args[1])
	if err != nil {
		return nil, errors.New("Expecting integer value for asset holding")
	}
	B = args[2]
	Bval, err = strconv.Atoi(args[3])
	if err != nil {
		return nil, errors.New("Expecting integer value for asset holding")
	}
	fmt.Printf("Aval = %d, Bval = %d\n", Aval, Bval)
	
	// Write the state to the ledger
	err = stub.PutState(A, []byte(strconv.Itoa(Aval)))
	if err != nil {
		return nil, err
	}

	err = stub.PutState(B, []byte(strconv.Itoa(Bval)))
	if err != nil {
		return nil, err
	}

	//iQOS Changes starts ---------------------------------------------------------------------
	/*
	// Check if table already exists
	_, err1 := stub.GetTable("AssemblyLine")
	if err1 == nil {
		// Table already exists; do not recreate
		return nil, nil
	}
*/
	// Create application Table
	err1 = stub.CreateTable("AssemblyLine", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "assemblyId", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "serialId", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "originalFilamentBatchId", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "originalLedBatchId", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "originalCircuitBoardBatchId", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "originalWireBatchId", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "originalCasingBatchId", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "originalAdaptorBatchId", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "originalStickPodBatchId", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "assemblyStatus", Type: shim.ColumnDefinition_STRING, Key: false},

	})
	if err1 != nil {
		return nil, errors.New("Failed creating AssemblyLine.")
	}
//iQOS Changes ends ---------------------------------------------------------------------

	return nil, nil
}

// Transaction makes payment of X units from A to B but 2 times the value
func (t *SimpleChaincode) invoke2(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Printf("Running invoke2")
	
	var A, B string    // Entities
	var Aval, Bval int // Asset holdings
	var X int          // Transaction value
	var err error

	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3")
	}

	A = args[0]
	B = args[1]

	// Get the state from the ledger
	// TODO: will be nice to have a GetAllState call to ledger
	Avalbytes, err := stub.GetState(A)
	if err != nil {
		return nil, errors.New("Failed to get state")
	}
	if Avalbytes == nil {
		return nil, errors.New("Entity not found")
	}
	Aval, _ = strconv.Atoi(string(Avalbytes))

	Bvalbytes, err := stub.GetState(B)
	if err != nil {
		return nil, errors.New("Failed to get state")
	}
	if Bvalbytes == nil {
		return nil, errors.New("Entity not found")
	}
	Bval, _ = strconv.Atoi(string(Bvalbytes))

	// Perform the execution
	X, err = strconv.Atoi(args[2])
	Aval = Aval - X - X 
	Bval = Bval + X + X 
	fmt.Printf("Aval = %d, Bval = %d\n", Aval, Bval)

	// Write the state back to the ledger
	err = stub.PutState(A, []byte(strconv.Itoa(Aval)))
	if err != nil {
		return nil, err
	}

	err = stub.PutState(B, []byte(strconv.Itoa(Bval)))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Transaction makes payment of X units from A to B
func (t *SimpleChaincode) invoke(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Printf("Running invoke")
	
	var A, B string    // Entities
	var Aval, Bval int // Asset holdings
	var X int          // Transaction value
	var err error

	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3")
	}

	A = args[0]
	B = args[1]

	// Get the state from the ledger
	// TODO: will be nice to have a GetAllState call to ledger
	Avalbytes, err := stub.GetState(A)
	if err != nil {
		return nil, errors.New("Failed to get state")
	}
	if Avalbytes == nil {
		return nil, errors.New("Entity not found")
	}
	Aval, _ = strconv.Atoi(string(Avalbytes))

	Bvalbytes, err := stub.GetState(B)
	if err != nil {
		return nil, errors.New("Failed to get state")
	}
	if Bvalbytes == nil {
		return nil, errors.New("Entity not found")
	}
	Bval, _ = strconv.Atoi(string(Bvalbytes))

	// Perform the execution
	X, err = strconv.Atoi(args[2])
	Aval = Aval - X
	Bval = Bval + X
	fmt.Printf("Aval = %d, Bval = %d\n", Aval, Bval)

	// Write the state back to the ledger
	err = stub.PutState(A, []byte(strconv.Itoa(Aval)))
	if err != nil {
		return nil, err
	}

	err = stub.PutState(B, []byte(strconv.Itoa(Bval)))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Deletes an entity from state
func (t *SimpleChaincode) delete(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Printf("Running delete")
	
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3")
	}

	A := args[0]

	// Delete the key from the state in ledger
	err := stub.DelState(A)
	if err != nil {
		return nil, errors.New("Failed to delete state")
	}

	return nil, nil
}

//iQOS Changes starts---------------------------------------------------------------------
//registerUser to register a user
func (t *SimpleChaincode) startAssemblyLine(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

if len(args) != 9 {
			return nil, fmt.Errorf("Incorrect number of arguments. Expecting 9. Got: %d.", len(args))
		}
		
		assemblyId:=args[0]
		serialId:=args[1]
		originalFilamentBatchId:=args[2]
		originalLedBatchId:=args[3]
		originalCircuitBoardBatchId:=args[4]
		originalWireBatchId:=args[5]
		originalCasingBatchId:=args[6]
		originalAdaptorBatchId:=args[7]
		originalStickPodBatchId:=args[8]
		assemblyStatus:= "InProgress"

		// Insert a row
		ok, err := stub.InsertRow("AssemblyLine", shim.Row{
			Columns: []*shim.Column{
				&shim.Column{Value: &shim.Column_String_{String_: assemblyId}},
				&shim.Column{Value: &shim.Column_String_{String_: serialId}},
				&shim.Column{Value: &shim.Column_String_{String_: originalFilamentBatchId}},
				&shim.Column{Value: &shim.Column_String_{String_: originalLedBatchId}},
				&shim.Column{Value: &shim.Column_String_{String_: originalCircuitBoardBatchId}},
				&shim.Column{Value: &shim.Column_String_{String_: originalWireBatchId}},
				&shim.Column{Value: &shim.Column_String_{String_: originalCasingBatchId}},
				&shim.Column{Value: &shim.Column_String_{String_: originalAdaptorBatchId}},
				&shim.Column{Value: &shim.Column_String_{String_: originalStickPodBatchId}},
				&shim.Column{Value: &shim.Column_String_{String_: assemblyStatus}},
			}})

		if err != nil {
			return nil, err 
		}
		if !ok && err == nil {
			return nil, errors.New("Row already exists.")
		}
			
		return nil, nil

}
//iQOS Changes ends---------------------------------------------------------------------

// Invoke callback representing the invocation of a chaincode
// This chaincode will manage two accounts A and B and will transfer X units from A to B upon invoke
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Printf("Invoke called, determining function")
	
	// Handle different functions
	if function == "invoke" {
		// Transaction makes payment of X units from A to B
		fmt.Printf("Function is invoke")
		return t.invoke(stub, args)
	} else if function == "init" {
		fmt.Printf("Function is init")
		return t.Init(stub, function, args)
	} else if function == "invoke2" {
		fmt.Printf("Function is invoke2")
		return t.invoke2(stub, args)
	} else if function == "startAssemblyLine" {
		fmt.Printf("Function is startAssemblyLine")
		return t.startAssemblyLine(stub, args)
	} else if function == "delete" {
		// Deletes an entity from its state
		fmt.Printf("Function is delete")
		return t.delete(stub, args)
	} else if function == "startAssemblyLine" {
		// Deletes an entity from its state
		fmt.Printf("Function is startAssemblyLine")
		return t.startAssemblyLine(stub, args)
	}

	return nil, errors.New("Received unknown function invocation")
}

func (t* SimpleChaincode) Run(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Printf("Run called, passing through to Invoke (same function)")
	
	// Handle different functions
	if function == "invoke" {
		// Transaction makes payment of X units from A to B
		fmt.Printf("Function is invoke")
		return t.invoke(stub, args)
	} else if function == "invoke2" {
		fmt.Printf("Function is invoke2")
		return t.invoke2(stub, args)
	} else if function == "startAssemblyLine" {
		fmt.Printf("Function is startAssemblyLine")
		return t.startAssemblyLine(stub, args)
	} else if function == "init" {
		fmt.Printf("Function is init")
		return t.Init(stub, function, args)
	} else if function == "delete" {
		// Deletes an entity from its state
		fmt.Printf("Function is delete")
		return t.delete(stub, args)
	}else if function == "startAssemblyLine" {
		// Deletes an entity from its state
		fmt.Printf("Function is startAssemblyLine")
		return t.startAssemblyLine(stub, args)
	}

	return nil, errors.New("Received unknown function invocation")
}

//iQOS Changes starts------------------------------------------------------------------------------------------
// Query to get Value of A/B
func (t *SimpleChaincode) getValue(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	
	
	var A string // Entities
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the person to query")
	}

	A = args[0]

	// Get the state from the ledger
	Avalbytes, err := stub.GetState(A)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + A + "\"}"
		return nil, errors.New(jsonResp)
	}

	if Avalbytes == nil {
		jsonResp := "{\"Error\":\"Nil amount for " + A + "\"}"
		return nil, errors.New(jsonResp)
	}

	jsonResp := "{\"Name\":\"" + A + "\",\"Amount\":\"" + string(Avalbytes) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)
	return Avalbytes, nil
}


//get the status against the AssemblyID
func (t *SimpleChaincode) getAssemblyLineStatus(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting AssemblyID to query")
	}

	assemblyID := args[0]
	

	// Get the row pertaining to this assemblyID
	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: assemblyID}}
	columns = append(columns, col1)

	row, err := stub.GetRow("AssemblyLine", columns)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get the data for the assemblyID " + assemblyID + "\"}"
		return nil, errors.New(jsonResp)
	}

	// GetRows returns empty message if key does not exist
	if len(row.Columns) == 0 {
		jsonResp := "{\"Error\":\"Failed to get the data for the assemblyID " + assemblyID + "\"}"
		return nil, errors.New(jsonResp)
	}

	
	
	res2E := GetAssemblyLineStatus{}
	
	res2E.AssemblyLineStatus = row.Columns[9].GetString_()
	
    //mapB, _ := json.Marshal(res2E)
    //fmt.Println(string(mapB))
	
	//return mapB, nil
	return []byte (res2E.AssemblyLineStatus), nil

}


// query queries the chaincode
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Printf("Query called, determining function")

	if function == "getValue" {
		t := SimpleChaincode{}
		return t.getValue(stub, args)		
	} else if function == "getAssemblyLineStatus" { 
		t := SimpleChaincode{}
		return t.getAssemblyLineStatus(stub, args)
	}
	
	return nil, nil
}
//iQOS Changes ends------------------------------------------------------------------------------------------

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

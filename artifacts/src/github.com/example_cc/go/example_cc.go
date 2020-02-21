/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

/*
 * The sample smart contract for documentation topic:
 * Writing Your First Blockchain Application
 */

package main

/* Imports
 * 4 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 2 specific Hyperledger Fabric specific libraries for Smart Contracts
 */
import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SimpleChaincode struct {
}

// Define the car structure, with 4 properties.  Structure tags are used by encoding/json library
type Trade struct {
	TradeId       string `json:"tradeID"`  
	FromParty     string `json:"fromParty"`
	ToParty       string `json:"toParty"`
	Amount        int    `json:"amount"`
	Status        string `json:"status"`
	Ctime         string `json:"ctime"`
}

/*
 * The Init method is called when the Smart Contract "fabcar" is instantiated by the blockchain network
 * Best practice is to have any Ledger initialization in separate function -- see initLedger()
 */
func (s *SimpleChaincode) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	_, args := APIstub.GetFunctionAndParameters()
 if len(args) != 5 {
                return shim.Error("Incorrect number of arguments. Expecting 6")
        }
        loc, _ := time.LoadLocation("America/Los_Angeles")
        current_time := time.Now().In(loc)
        currTime := current_time.Format("2006-01-02 15:04:05")
        amount, err := strconv.Atoi(args[3])
        if err != nil {
                return shim.Error("3r argument must be a numeric string")
        }

        var car = Trade{TradeId: args[0], FromParty: args[1], ToParty: args[2], Amount: amount, Status: args[4], Ctime: currTime}

        carAsBytes, _ := json.Marshal(car)
        APIstub.PutState(args[0], carAsBytes)

	return shim.Success(nil)
}

/*
 * The Invoke method is called as a result of an application request to run the Smart Contract "fabcar"
 * The calling application program has also specified the particular smart contract function to be called, with arguments
 */
func (s *SimpleChaincode) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	 if function == "createCar" {
		return s.createCar(APIstub, args)
	}else if function == "queryAllTrades" {
		return s.queryAllTrades(APIstub)
	}else if function == "updateStatus" {
		return s.updateStatus(APIstub, args)
	}else if function == "query" {
		return s.queryforTradeID(APIstub, args)
	} 

	return shim.Error("Invalid Smart Contract function name.")
}



func (s *SimpleChaincode) createCar(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 6")
	}
	loc, _ := time.LoadLocation("America/Los_Angeles")
        current_time := time.Now().In(loc)
	currTime := current_time.Format("2006-01-02 15:04:05")
	amount, err := strconv.Atoi(args[3])
	if err != nil {
		return shim.Error("3r argument must be a numeric string")
	}

	var car = Trade{TradeId: args[0], FromParty: args[1], ToParty: args[2], Amount: amount, Status: args[4], Ctime: currTime}

	carAsBytes, _ := json.Marshal(car)
	APIstub.PutState(args[0], carAsBytes)

	return shim.Success(nil)
}
func (s *SimpleChaincode) queryAllTrades(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := ""
	endKey := ""

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAllTrades:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}
func (s *SimpleChaincode) updateStatus(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	tradeAsBytes, _ := APIstub.GetState(args[0])
	trade := Trade{}

	json.Unmarshal(tradeAsBytes, &trade)
	trade.Status = "InProgress"

	tradeAsBytes, _ = json.Marshal(trade)
	APIstub.PutState(args[0], tradeAsBytes)

	return shim.Success(nil)
}
// query callback representing the query of a chaincode
func (t *SimpleChaincode) queryforTradeID(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	var A string // Entities
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the person to query")
	}

	A = args[0]

	// Get the state from the ledger
	Avalbytes, err := stub.GetState(A)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + A + "\"}"
		return shim.Error(jsonResp)
	}

	if Avalbytes == nil {
		jsonResp := "{\"Error\":\"No Response for Trade" + A + "\"}"
		return shim.Error(jsonResp)
	}

	jsonResp := "{\"SearchTrade\":\"" + A + "\",\"Response\":\"" + string(Avalbytes) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)
	return shim.Success(Avalbytes)
}


// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}

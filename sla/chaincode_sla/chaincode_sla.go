/*copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

//WARNING - this chaincode's ID is hard-coded in chaincode_example04 to illustrate one way of
//calling chaincode from a chaincode. If this example is modified, chaincode_example04.go has
//to be modified as well with the new ID of chaincode_example02.
//chaincode_example05 show's how chaincode ID can be passed in as a parameter instead of
//hard-coding.

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// ===========================================================
//  Struct 및 Constant 정의
// ===========================================================
type SimpleChaincode struct {
}

// key-value store 의 키 구분자
const ContractIDSeparator = "|"
const FIELDSEP = "|"
const ENTRYSEP = "$"


func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	var A, B string    // Entities
	var Aval, Bval int // Asset holdings
	var err error

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

	return nil, nil
}

// Transaction makes payment of X units from A to B
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	//function, args := stub.GetFunctionAndParameters()

	switch function {

	case "registerContract":
		return t.registerContract(stub, args)
	}

	return nil, errors.New("Invalid invoke function name. Expecting \"searchContractByID\" \"searchContractListByName~\"")
}

// Query callback representing the query of a chaincode
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	switch function {

	case "searchContractByID":
		return t.searchContractByID(stub, args)

	case "searchContractListByName":
		return t.searchContractListByName(stub, args)

	case "searchContractListByClient":
		return t.searchContractListByClient(stub, args)

	}
	return nil, errors.New("Invalid Query function name. Expecting \"searchContractByID\" \"searchContractListByName~\"")
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Deletes an entity from state
func (t *SimpleChaincode) delete(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	A := args[0]

	// Delete the key from the state in ledger
	err := stub.DelState(A)
	if err != nil {
		return nil, errors.New("Failed to delete state")
	}

	return nil, nil
}

// ===========================================================
//  SLAChaincodeStub 검색 함수
// ===========================================================
func (t *SimpleChaincode) registerContract(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the Value to registerContract")
	}

	contractID := args[0]
	contractName := strings.Split(args[1], ",")[1]
	contractClient := strings.Split(args[1], ",")[4]

	//1. 계약ID 등록
	stub.PutState(contractID, []byte(args[1]))

	{ //2. 계약명 등록
		var err error

		//계약명으로 기존내역조회
		contractIDsInBytes, err := stub.GetState(contractName) // 리턴값 ([]byte, error)
		if err != nil {
			return nil, errors.New("Failed to get state with" + string(contractIDsInBytes))
		}

		contractIDsInString := string(contractIDsInBytes)

		//기존내역이 없을경우 "계약명"-"계약ID목록" 등록
		if contractIDsInString == "" {
			err = stub.PutState(contractName, []byte(contractID))
			if err != nil {
				return nil, err
			}

		} else {
			err = stub.PutState(contractName, []byte(contractIDsInString+ContractIDSeparator+contractID))
			if err != nil {
				return nil, err
			}
		}
	}

	{ //3. 고객사명 등록

		var err error

		//계약명으로 기존내역조회
		contractIDsInBytes, _ := stub.GetState(contractClient) // 리턴값 ([]byte, error)
		contractIDsInString := string(contractIDsInBytes)

		//기존내역이 없을경우 "고객사명"-"계약ID목록" 등록
		if contractIDsInString == "" {
			err = stub.PutState(contractClient, []byte(contractID))
			if err != nil {
				return nil, err
			}
		} else {
			err = stub.PutState(contractClient, []byte(contractIDsInString+ContractIDSeparator+contractID))
			if err != nil {
				return nil, err
			}
		}
	}
	return nil, nil
}

// ===========================================================
//  SLAChaincodeStub 검색 함수
// ===========================================================

// SLA ID 검색
func (t *SimpleChaincode) searchContractByID(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	var dataInBytes string // Entities
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the Value to searchContractByID")
	}

	dataInBytes = args[0]
	Valuebytes, err := stub.GetState(args[0])

	if err != nil {
		return nil, errors.New("Failed to get state with" + dataInBytes)
	}

	fmt.Printf("searchbyid Response:%s\n", Valuebytes)

	return Valuebytes, nil
}

// SLA 계약명 검색
func (t *SimpleChaincode) searchContractListByName(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	var dataInBytes string // Entities
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the value to searchContractListByName")
	}

	dataInBytes = args[0]
	contractName := args[0]

	// 계약명으로 계약ID목록 조회
	contractIDsInBytes, err := stub.GetState(contractName)
	contractIDsInString := string(contractIDsInBytes)
	if err != nil {
		return nil, errors.New("Failed to get state with " + dataInBytes)
	}

	// 계약ID목록의 형태를 스트링에서 배열로 전환
	contractIDs := strings.Split(contractIDsInString, ContractIDSeparator)

	// 리턴값 초기화
	contractList := make([]string, len(contractIDs))

	// 계약ID목록으로 계약내용을 추출하여 계약목록 작성
	for i, _ := range contractIDs {
		contractInBytes, _ := stub.GetState(contractIDs[i])
		contractList[i] = string(contractInBytes)
	}

	contractListBytes := strings.Join(contractList, ContractIDSeparator)

	fmt.Printf("searchContractListByName Response:%s\n", contractListBytes)

	return []byte(contractListBytes), nil

}

// SLA 고객사명 검색
func (t *SimpleChaincode) searchContractListByClient(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	var dataInBytes string // Entities
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the value to searchContractListByClient")
	}

	dataInBytes = args[0]
	contractClient := args[0]

	contractIDsInBytes, err := stub.GetState(contractClient)
	contractIDsInString := string(contractIDsInBytes)

	if err != nil {
		return nil, errors.New("Failed to get state with " + dataInBytes)
	}

	// 계약ID목록의 형태를 스트링에서 배열로 전환
	contractIDs := strings.Split(contractIDsInString, ContractIDSeparator)

	// 리턴값 초기화
	contractList := make([]string, len(contractIDs))

	// 계약ID목록으로 계약내용을 추출하여 계약목록 작성
	for i, _ := range contractIDs {
		contractInBytes, _ := stub.GetState(contractIDs[i])
		contractList[i] = string(contractInBytes)
	}

	contractListBytes := strings.Join(contractList, ContractIDSeparator)

	fmt.Printf("searchContractListByClient Response:%s\n", contractListBytes)

	return []byte(contractListBytes), nil
}

// ===========================================================
//  SLAChaincodeStub 업데이트 함수
// ===========================================================

func (t *SimpleChaincode) updateContractId(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	var dataInBytes string // Entities
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the value to searchContractListByClient")
	}

	dataInBytes = args[0]
	contractID := args[0]

	// 기존내용 조회
	contractIDsInBytes, err := stub.GetState(contractID)
	if err != nil {
		return nil, errors.New("Failed to get state with " + string(contractIDsInBytes))
	}

	// UPDATDE 처리
	stub.PutState(contractID, []byte(args[1]))

	// 변경내용 조회
	update_value, err := stub.GetState(contractID)
	if err != nil {
		return nil, errors.New("Failed to get state with " + dataInBytes)
	}

	fmt.Printf("searchContractListByClient Response:%s\n", update_value)

	return []byte(update_value), nil
}

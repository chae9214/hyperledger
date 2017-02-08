package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// ===========================================================
//  Struct 및 Constant 정의
// ===========================================================

type SimpleChaincode struct {
}

// fraud entry 의 필드갯수와 각 필드별 인덱스
const NUM_FIELDS = 8

const IND_CID = 0
const IND_MAC = 1
const IND_UUID = 2
const IND_FINALDATE = 3
const IND_FINALTIME = 4
const IND_FDSPRODUCEDBY = 5
const IND_FDSREGISTEREDBY = 6
const IND_FDSREASON = 7

// key-value store 의 키 구분자
const PREFIX_EID = "eid_"
const PREFIX_CID = "cid_"
const PREFIX_MAC = "mac_"
const PREFIX_UUID = "uuid_"

const SEP = "|"

/*
// ===========================================================
//  Initialization 함수
// ===========================================================

func CreateStub() ChaincodeStubInterface {
	var stub ChaincodeStubInterface
	stub.kvs = make(map[string][]byte)
	return stub
}

func CreateFDSChaincodeStub() FDSChaincodeStub {
	kvs := make(map[string][]byte)
	nextEID := 1
	stub := FDSChaincodeStub{ChaincodeStubInterface{kvs}, nextEID}
	return stub
}

func CreateSLAChaincodeStub() SLAChaincodeStub {
	kvs := make(map[string][]byte)
	stub := SLAChaincodeStub{ChaincodeStubInterface{kvs}}
	return stub
}

func (stub *ChaincodeStubInterface) String() string {
	var s string

	s += "<KVS>\n"
	for k, v := range stub.kvs {
		s += "\t" + k + ": " + string(v) + "\n"
	}
	return s
}

// ===========================================================
//  Serialize / Deserialize 함수
// ===========================================================

func stringToByteArray(s string) []byte {
	return []byte(s)
}

func byteArrayToString(b []byte) string {
	return string(b)
}

func stringArrayToByteArray(slist []string) []byte {
	return stringToByteArray(strings.Join(slist, SEP))
}

func byteArrayToStringArray(b []byte) []string {
	if byteArrayToString(b) == "" {
		return []string{}
	} else {
		return strings.Split(byteArrayToString(b), SEP)
	}
}

func appendToEIDList(b []byte, eid string) []byte {
	eidKeys := byteArrayToStringArray(b)
	return stringArrayToByteArray(append(eidKeys, eid))
}

// ===========================================================
//  ChaincodeStubInterface 함수
// ===========================================================

func (stub *ChaincodeStubInterface) PutState(key string, value []byte) error {
	if value == nil {
		return errors.New("entry cannot be empty")
	} else {
		stub.kvs[key] = value
		return nil
	}
}

func (stub *ChaincodeStubInterface) GetState(key string) ([]byte, error) {
	value := stub.kvs[key]
	return value, nil
}

func (stub *ChaincodeStubInterface) DelState(key string) error {
	delete(stub.kvs, key)
	return nil
}

func (stub *ChaincodeStubInterface) GetKVSLength() int {
	return len(stub.kvs)
}
*/

// ===========================================================
//  SimpleChaincode 함수
// ===========================================================

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

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("ex02 Invoke")
	function, args := stub.GetFunctionAndParameters()
	if function == "invoke" {
		// Make payment of X units from A to B
		return t.invoke(stub, args)
	} else if function == "delete" {
		// Deletes an entity from its state
		return t.delete(stub, args)
	} else if function == "query" {
		// the old "Query" is now implemtned in invoke
		return t.query(stub, args)
	}

	return shim.Error("Invalid invoke function name. Expecting \"invoke\" \"delete\" \"query\"")
}

func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function != "query" {
		return nil, errors.New("Invalid query function name. Expecting \"query\"")
	}
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

// ===========================================================
//  FDS 등록/수정 함수
// ===========================================================

func (stub *FDSChaincodeStub) RegisterFraudEntry(fields []string) bool {
	if len(fields) != NUM_FIELDS {
		return false
	}

	eidKey := PREFIX_EID + strconv.Itoa(stub.nextEID)
	cidKey := PREFIX_CID + fields[IND_CID]
	macKey := PREFIX_MAC + fields[IND_MAC]
	uuidKey := PREFIX_UUID + fields[IND_UUID]

	stub.PutState(eidKey, stringArrayToByteArray(fields))

	b, err := stub.GetState(cidKey)
	if err == nil {
		stub.PutState(cidKey, appendToEIDList(b, eidKey))
	} else {
		return false
	}
	b, err = stub.GetState(macKey)
	if err == nil {
		stub.PutState(macKey, appendToEIDList(b, eidKey))
	} else {
		return false
	}
	b, err = stub.GetState(uuidKey)
	if err == nil {
		stub.PutState(uuidKey, appendToEIDList(b, eidKey))
	} else {
		return false
	}

	stub.nextEID++
	return true
}

// ===========================================================
//  FDS 삭제 함수
// ===========================================================

func (stub *FDSChaincodeStub) RemoveWithCID(cid string) bool {
	cidKey := PREFIX_CID + cid

	b, err := stub.GetState(cidKey)
	if err == nil {
		eidKeys := byteArrayToStringArray(b)

		for _, eidKey := range eidKeys {
			stub.DelState(eidKey)
		}
		stub.DelState(cidKey)
	} else {
		return false
	}
	return true
}

func (stub *FDSChaincodeStub) RemoveWithMAC(mac string) bool {
	macKey := PREFIX_MAC + mac

	b, err := stub.GetState(macKey)
	if err == nil {
		eidKeys := byteArrayToStringArray(b)

		for _, eidKey := range eidKeys {
			stub.DelState(eidKey)
		}
		stub.DelState(macKey)
	} else {
		return false
	}
	return true
}

func (stub *FDSChaincodeStub) RemoveWithUUID(uuid string) bool {
	uuidKey := PREFIX_UUID + uuid

	b, err := stub.GetState(uuidKey)
	if err == nil {
		eidKeys := byteArrayToStringArray(b)

		for _, eidKey := range eidKeys {
			stub.DelState(eidKey)
		}
		stub.DelState(uuidKey)
	} else {
		return false
	}
	return true
}

// ===========================================================
//  FDS 조회 함수
// ===========================================================

func (stub *FDSChaincodeStub) LookupWithCID(cid string) (entries [][]string, result bool) {
	cidKey := PREFIX_CID + cid

	b, err := stub.GetState(cidKey)
	if err == nil {
		eidKeys := byteArrayToStringArray(b)

		entries = make([][]string, len(eidKeys))
		for i, eidKey := range eidKeys {
			b, _ = stub.GetState(eidKey)
			entries[i] = byteArrayToStringArray(b)
		}
	} else {
		return entries, false
	}
	return entries, true
}

func (stub *FDSChaincodeStub) LookupWithMAC(mac string) (entries [][]string, result bool) {
	macKey := PREFIX_MAC + mac

	b, err := stub.GetState(macKey)
	if err == nil {
		eidKeys := byteArrayToStringArray(b)

		entries = make([][]string, len(eidKeys))
		for i, eidKey := range eidKeys {
			b, _ = stub.GetState(eidKey)
			entries[i] = byteArrayToStringArray(b)
		}
	} else {
		return entries, false
	}
	return entries, true
}

func (stub *FDSChaincodeStub) LookupWithUUID(uuid string) (entries [][]string, result bool) {
	uuidKey := PREFIX_UUID + uuid

	b, err := stub.GetState(uuidKey)
	if err == nil {
		eidKeys := byteArrayToStringArray(b)

		entries = make([][]string, len(eidKeys))
		for i, eidKey := range eidKeys {
			b, _ = stub.GetState(eidKey)
			entries[i] = byteArrayToStringArray(b)
		}
	} else {
		return entries, false
	}
	return entries, true
}

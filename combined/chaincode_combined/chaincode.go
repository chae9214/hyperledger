package chaincode

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

type ChaincodeStubInterface struct {
	kvs map[string][]byte
}

type FDSChaincodeStub struct {
	ChaincodeStubInterface
	nextEID int
}

type SLAChaincodeStub struct {
	ChaincodeStubInterface
}

// 사용하지 않는 struct (대신 string array 사용)
type FDSValues struct {
	Cid             string
	Mac             string
	Uuid            string
	FinalDate       string
	FinalTime       string
	FDSProducedBy   string
	FDSRegisteredBy string
	FDSReason       string
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

const SEP = "$"
const ContractIDSeparator = "|"

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
//  FDSValues 함수 ***사용하지 않음!!!***
// ===========================================================

func FDSValuesToByteArray(v FDSValues) ([]byte, error) {
	b, err := json.Marshal(v)
	return b, err
}

func ByteArrayToFDSValues(b []byte) (FDSValues, error) {
	var v FDSValues
	err := json.Unmarshal(b, &v)
	return v, err
}

func (stub *FDSChaincodeStub) RegisterFraudEntryUsingFDSValues(fields []string) bool {
	if len(fields) != 8 {
		return false
	}

	v := FDSValues{fields[0], fields[1], fields[2], fields[3], fields[4], fields[5], fields[6], fields[7]}
	b, err := FDSValuesToByteArray(v)
	if err == nil {
		stub.kvs[fields[0]] = b
		return true
	}
	return false
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

// ===========================================================
//  FDSChaincodeStub 등록/수정 함수
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
//  FDSChaincodeStub 삭제 함수
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
//  FDSChaincodeStub 조회 함수
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

// ===========================================================
//  SLAChaincodeStub 등록 함수
// ===========================================================

func (t *SimpleChaincode) registerContract(stub ChaincodeStubInterface, args []string) {

	contractID := args[0]
	contractName := strings.Split(args[1], ",")[1]
	contractClient := strings.Split(args[1], ",")[4]

	//1. 계약ID 등록
	stub.PutState(contractID, []byte(args[1]))

	{ //2. 계약명 등록
		//계약명으로 기존내역조회
		contractIDsInBytes, _ := stub.GetState(contractName) // 리턴값 ([]byte, error)
		contractIDsInString := string(contractIDsInBytes)

		//기존내역이 없을경우 "계약명"-"계약ID목록" 등록
		if contractIDsInString == "" {
			stub.PutState(contractName, []byte(contractID))
		} else {
			stub.PutState(contractName, []byte(contractIDsInString+ContractIDSeparator+contractID))
		}
	}

	{ //3. 고객사명 등록
		//계약명으로 기존내역조회
		contractIDsInBytes, _ := stub.GetState(contractClient) // 리턴값 ([]byte, error)
		contractIDsInString := string(contractIDsInBytes)

		//기존내역이 없을경우 "고객사명"-"계약ID목록" 등록
		if contractIDsInString == "" {
			stub.PutState(contractClient, []byte(contractID))
		} else {
			stub.PutState(contractClient, []byte(contractIDsInString+ContractIDSeparator+contractID))
		}
	}
}

// ===========================================================
//  SLAChaincodeStub 검색 함수
// ===========================================================

// SLA ID 검색
func (t *SimpleChaincode) searchContractByID(stub ChaincodeStubInterface, args []string) string {
	value, _ := stub.GetState(args[0])
	return string(value)
}

// SLA 계약명 검색
func (t *SimpleChaincode) searchContractListByName(stub ChaincodeStubInterface, args []string) []string {
	contractName := args[0]

	// 계약명으로 계약ID목록 조회
	contractIDsInBytes, _ := stub.GetState(contractName)
	contractIDsInString := string(contractIDsInBytes)

	// 계약ID목록의 형태를 스트링에서 배열로 전환
	contractIDs := strings.Split(contractIDsInString, ContractIDSeparator)

	// 리턴값 초기화
	contractList := make([]string, len(contractIDs))

	// 계약ID목록으로 계약내용을 추출하여 계약목록 작성
	for i, _ := range contractIDs {
		contractInBytes, _ := stub.GetState(contractIDs[i])
		contractList[i] = string(contractInBytes)
	}

	return contractList
}

// SLA 고객사명 검색
func (t *SimpleChaincode) searchContractListByClient(stub ChaincodeStubInterface, args []string) []string {
	contractClient := args[0]

	// 고객사명으로 계약ID목록 조회
	contractIDsInBytes, _ := stub.GetState(contractClient)
	contractIDsInString := string(contractIDsInBytes)

	// 계약ID목록의 형태를 스트링에서 배열로 전환
	contractIDs := strings.Split(contractIDsInString, ContractIDSeparator)

	// 리턴값 초기화
	contractList := make([]string, len(contractIDs))

	// 계약ID목록으로 계약내용을 추출하여 계약목록 작성
	for i, _ := range contractIDs {
		contractInBytes, _ := stub.GetState(contractIDs[i])
		contractList[i] = string(contractInBytes)
	}
	return contractList
}

// ===========================================================
//  SLAChaincodeStub 업데이트 함수
// ===========================================================

func (t *SimpleChaincode) updateContractId(stub ChaincodeStubInterface, args []string) string {

	contractID := args[0]

	// 기존내용 조회
	contractIDsInBytes, _ := stub.GetState(contractID) // 리턴값 ([]byte, error)
	contractIDsInString := string(contractIDsInBytes)

	// 기존내역이 없을경우 확인 여부
	if contractIDsInString == "" {
		fmt.Printf("No date :%v\n", contractIDsInString)
	}

	// UPDATDE 처리
	stub.PutState(contractID, []byte(args[1]))

	// 변경내용 조회
	update_value, _ := stub.GetState(contractID)

	return string(update_value)
}

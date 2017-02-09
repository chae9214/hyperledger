package chaincode

import (
	"bytes"
	//"fmt"
	"reflect"
	"strconv"
	"testing"
)

// ===========================================================
//  Helper Function
// ===========================================================

func initializeEntries(n int) [][]string {
	var entries = make([][]string, n)

	for i := range entries {
		entries[i] = make([]string, NUM_FIELDS)
		entries[i][IND_CID] = "cid" + strconv.Itoa(i+1)
		entries[i][IND_MAC] = "mac" + strconv.Itoa(i+1)
		entries[i][IND_UUID] = "uuid" + strconv.Itoa(i+1)
		entries[i][IND_FINALDATE] = "최종거래일자" + strconv.Itoa(i+1)
		entries[i][IND_FINALTIME] = "최종거래시간" + strconv.Itoa(i+1)
		entries[i][IND_FDSPRODUCEDBY] = "FDS제공처" + strconv.Itoa(i+1)
		entries[i][IND_FDSREGISTEREDBY] = "FDS등록처" + strconv.Itoa(i+1)
		entries[i][IND_FDSREASON] = "FDS등록사유" + strconv.Itoa(i+1)
	}
	return entries
}

// ===========================================================
//  ChaincodeStubInterface 함수 테스트
// ===========================================================

func TestPutState(t *testing.T) {
	var stub = CreateFDSChaincodeStub()

	key := "key"
	value := []byte("value")
	stub.PutState(key, value)

	expected := 1
	actual := stub.GetKVSLength()

	if !(expected == actual) {
		t.Errorf("Expected %v, but returned %v instead", expected, actual)
	}
}

func TestGetState(t *testing.T) {
	var stub = CreateFDSChaincodeStub()

	key := "key"
	value := []byte("value")
	stub.PutState(key, value)

	expected := []byte("value")
	actual, _ := stub.GetState(key)

	if !(bytes.Equal(expected, actual)) {
		t.Errorf("Expected %v, but returned %v instead", expected, actual)
	}
}

func TestPutStateAndGetState(t *testing.T) {
	var c = CreateStub()

	c.PutState("Alice", []byte("Married"))

	expected := []byte("Married")
	actual, _ := c.GetState("Alice")

	if !(bytes.Equal(expected, actual)) {
		t.Errorf("Expected %v, but returned %v instead.", expected, actual)
	}

	c.PutState("Bob", []byte("Born"))

	expected = []byte("Born")
	actual, _ = c.GetState("Bob")

	if !(bytes.Equal(expected, actual)) {
		t.Errorf("Expected %v, but returned %v instead.", expected, actual)
	}
}

func TestDelState(t *testing.T) {
	var c = CreateStub()

	c.PutState("Alice", []byte("Married"))
	c.PutState("Bob", []byte("Born"))

	c.DelState("Alice")

	expected := []byte{}
	actual, _ := c.GetState("Alice")

	if !(bytes.Equal(expected, actual)) {
		t.Errorf("Expected %v, but returned %v instead.", expected, actual)
	}

	expected = []byte("Born")
	actual, _ = c.GetState("Bob")

	if !(bytes.Equal(expected, actual)) {
		t.Errorf("Expected %v, but returned %v instead.", expected, actual)
	}
}

// ===========================================================
//  FDSChaincodeStub 등록/수정 테스트
// ===========================================================

func TestRegiserFraudEntry(t *testing.T) {
	var stub = CreateFDSChaincodeStub()

	num_entries := 10
	var entries = initializeEntries(num_entries)
	for _, entry := range entries {
		stub.RegisterFraudEntry(entry)
	}

	expected := num_entries * 4
	actual := stub.GetKVSLength()

	if !(expected == actual) {
		t.Errorf("Expected %v, but returned %v instead", expected, actual)
	}

}

func TestInvalidRegister(t *testing.T) {
	var stub = CreateFDSChaincodeStub()

	longEntry := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}
	shortEntry := []string{"0", "1", "2", "3", "4", "5", "6"}

	expected := false
	actual := stub.RegisterFraudEntry(longEntry)
	if expected != actual {
		t.Errorf("Expected %v, but returned %v instead", expected, actual)
	}

	expected = false
	actual = stub.RegisterFraudEntry(shortEntry)
	if expected != actual {
		t.Errorf("Expected %v, but returned %v instead", expected, actual)
	}
}

func TestRegisterMultipleFraudEntryWithSameCID(t *testing.T) {
	var stub = CreateFDSChaincodeStub()

	num_entries := 10
	var entries = initializeEntries(num_entries)
	for _, entry := range entries {
		stub.RegisterFraudEntry(entry)
	}
}

// ===========================================================
//  FDSChaincodeStub 조회 테스트
// ===========================================================

func TestLookupWithCID(t *testing.T) {
	var stub = CreateFDSChaincodeStub()

	num_entries := 10
	var entries = initializeEntries(num_entries)
	for _, entry := range entries {
		stub.RegisterFraudEntry(entry)
	}

	expected := entries[5:6]
	actual, _ := stub.LookupWithCID(entries[5][IND_CID])

	if !(reflect.DeepEqual(expected, actual)) {
		t.Errorf("Expected %v, but returned %v instead", expected, actual)
	}
}

func TestLookupWithMAC(t *testing.T) {
	var stub = CreateFDSChaincodeStub()

	num_entries := 10
	var entries = initializeEntries(num_entries)
	for _, entry := range entries {
		stub.RegisterFraudEntry(entry)
	}

	expected := entries[5:6]
	actual, _ := stub.LookupWithMAC(entries[5][IND_MAC])

	if !(reflect.DeepEqual(expected, actual)) {
		t.Errorf("Expected %v, but returned %v instead", expected, actual)
	}
}

func TestLookupWithUUID(t *testing.T) {
	var stub = CreateFDSChaincodeStub()

	num_entries := 10
	var entries = initializeEntries(num_entries)
	for _, entry := range entries {
		stub.RegisterFraudEntry(entry)
	}

	expected := entries[5:6]
	actual, _ := stub.LookupWithUUID(entries[5][IND_UUID])

	if !(reflect.DeepEqual(expected, actual)) {
		t.Errorf("Expected %v, but returned %v instead", expected, actual)
	}
}

// ===========================================================
//  FDSChaincodeStub 삭제 테스트
// ===========================================================

func TestRemoveWithCID(t *testing.T) {
	var stub = CreateFDSChaincodeStub()

	num_entries := 10
	var entries = initializeEntries(num_entries)
	for _, entry := range entries {
		stub.RegisterFraudEntry(entry)
	}
	stub.RemoveWithCID(entries[5][IND_CID])

	expected := [][]string{}
	actual, _ := stub.LookupWithCID(entries[5][IND_CID])
	if !(len(expected) == len(actual)) {
		t.Errorf("Expected %d entries, but returned %d entries", len(expected), len(actual))
	}

	expected = entries[2:3]
	actual, _ = stub.LookupWithCID(entries[2][IND_CID])
	if !(reflect.DeepEqual(expected, actual)) {
		t.Errorf("Expected %v, but returned %v instead", expected, actual)
	}
}

func TestRemoveWithMAC(t *testing.T) {
	var stub = CreateFDSChaincodeStub()

	num_entries := 10
	var entries = initializeEntries(num_entries)
	for _, entry := range entries {
		stub.RegisterFraudEntry(entry)
	}
	stub.RemoveWithMAC(entries[5][IND_MAC])

	expected := [][]string{}
	actual, _ := stub.LookupWithMAC(entries[5][IND_MAC])
	if !(len(expected) == len(actual)) {
		t.Errorf("Expected %d entries, but returned %d entries", len(expected), len(actual))
	}

	expected = entries[2:3]
	actual, _ = stub.LookupWithMAC(entries[2][IND_MAC])
	if !(reflect.DeepEqual(expected, actual)) {
		t.Errorf("Expected %v, but returned %v instead", expected, actual)
	}
}

func TestRemoveWithUUID(t *testing.T) {
	var stub = CreateFDSChaincodeStub()

	num_entries := 10
	var entries = initializeEntries(num_entries)
	for _, entry := range entries {
		stub.RegisterFraudEntry(entry)
	}
	stub.RemoveWithUUID(entries[5][IND_UUID])

	expected := [][]string{}
	actual, _ := stub.LookupWithUUID(entries[5][IND_UUID])
	if !(len(expected) == len(actual)) {
		t.Errorf("Expected %d entries, but returned %d entries", len(expected), len(actual))
	}

	expected = entries[2:3]
	actual, _ = stub.LookupWithUUID(entries[2][IND_UUID])
	if !(reflect.DeepEqual(expected, actual)) {
		t.Errorf("Expected %v, but returned %v instead", expected, actual)
	}
}

// ===========================================================
//  SLAChaincodeStub 등록 테스트
// ===========================================================

func TestRegisterContractAndSearchContractDetails(t *testing.T) {
	var c = CreateStub()
	var s SimpleChaincode

	reg_args := []string{"SLA_REG_2007-01-00001", "SLA_REG_2007-01-00001,ITSM 도급계약,도급계약,2017.01-2017.12,신한은행,지은탁,010-1234-5678,2017.01.31,대상"}
	s.registerContract(c, reg_args)

	reg_args = []string{"SLA_REG_2007-01-00002", "SLA_REG_2007-01-00002,ITSM 도급계약,도급계약,2017.01-2017.12,신한은행,지은탁,010-1234-5678,2017.01.31,대상"}
	s.registerContract(c, reg_args)

	reg_args = []string{"SLA_REG_2007-01-00003", "SLA_REG_2007-01-00003,ITSM 도급계약,도급계약,2017.01-2017.12,신한은행,지은탁,010-1234-5678,2017.01.31,대상"}
	s.registerContract(c, reg_args)

	search_args := []string{"SLA_REG_2007-01-00001"}
	expected := "SLA_REG_2007-01-00001,ITSM 도급계약,도급계약,2017.01-2017.12,신한은행,지은탁,010-1234-5678,2017.01.31,대상"
	actual := s.searchContractByID(c, search_args)

	if expected != actual {
		t.Errorf("Expected:%v, but returned:%v", expected, actual)
	}

	search_args = []string{"SLA_REG_2007-01-00002"}
	expected = "SLA_REG_2007-01-00002,ITSM 도급계약,도급계약,2017.01-2017.12,신한은행,지은탁,010-1234-5678,2017.01.31,대상"
	actual = s.searchContractByID(c, search_args)

	if expected != actual {
		t.Errorf("Expected:%v, but returned:%v", expected, actual)
	}

	// 없는 데이터 조회
	search_args = []string{"SLA_REG_2007-01-00004"}
	expected = ""
	actual = s.searchContractByID(c, search_args)

	if expected != actual {
		t.Errorf("Expected:%v, but returned:%v", expected, actual)
	}
}

// ===========================================================
//  SLAChaincodeStub 검색 테스트
// ===========================================================

// SLA 이름검색
func TestSearchContractListByName(t *testing.T) {

	var c = CreateStub()
	var s SimpleChaincode

	reg_args := []string{"SLA_REG_2007-01-00001", "SLA_REG_2007-01-00001,ITSM 도급계약,도급계약,2017.01-2017.12,신한은행,지은탁,010-1234-5678,2017.01.31,대상"}
	s.registerContract(c, reg_args)

	reg_args = []string{"SL_REG_2007-01-00002", "SLA_REG_2007-01-00002,ITSM 도급계약,도급계약,2017.01-2017.12,신한은행,지은탁,010-1234-5678,2017.01.31,대상"}
	s.registerContract(c, reg_args)

	reg_args = []string{"SLA_REG_2007-01-00003", "SLA_REG_2007-01-00003,ITSM 도급계약3,도급계약,2017.01-2017.12,신한은행,지은탁,010-1234-5678,2017.01.31,대상"}
	s.registerContract(c, reg_args)

	args := []string{"ITSM 도급계약"}
	expected := []string{"SLA_REG_2007-01-00001,ITSM 도급계약,도급계약,2017.01-2017.12,신한은행,지은탁,010-1234-5678,2017.01.31,대상",
		"SLA_REG_2007-01-00002,ITSM 도급계약,도급계약,2017.01-2017.12,신한은행,지은탁,010-1234-5678,2017.01.31,대상"}
	actual := s.searchContractListByName(c, args)

	if len(expected) != len(actual) {
		t.Errorf("Expected:%v, but returned:%v", len(expected), len(actual))
	}

	for !(reflect.DeepEqual(expected, actual)) {
		t.Errorf("Expected:%v, but returned:%v", expected, actual)

	}
}

// SLA 고객사검색
func TestSearchContractListByClient(t *testing.T) {

	var c = CreateStub()
	var s SimpleChaincode

	reg_args := []string{"SLA_REG_2007-01-00001", "SLA_REG_2007-01-00001,ITSM 도급계약,도급계약,2017.01-2017.12,신한은행,지은탁,010-1234-5678,2017.01.31,대상"}
	s.registerContract(c, reg_args)

	reg_args = []string{"SL_REG_2007-01-00002", "SLA_REG_2007-01-00002,ITSM 도급계약,도급계약,2017.01-2017.12,신한은행,지은탁,010-1234-5678,2017.01.31,대상"}
	s.registerContract(c, reg_args)

	reg_args = []string{"SLA_REG_2007-01-00003", "SLA_REG_2007-01-00003,ITSM 도급계약,도급계약,2017.01-2017.12,신한생명,지은탁,010-1234-5678,2017.01.31,대상"}
	s.registerContract(c, reg_args)

	args := []string{"신한은행"}
	expected := []string{"SLA_REG_2007-01-00001,ITSM 도급계약,도급계약,2017.01-2017.12,신한은행,지은탁,010-1234-5678,2017.01.31,대상",
		"SLA_REG_2007-01-00002,ITSM 도급계약,도급계약,2017.01-2017.12,신한은행,지은탁,010-1234-5678,2017.01.31,대상"}
	actual := s.searchContractListByClient(c, args)

	if len(expected) != len(actual) {
		t.Errorf("Expected:%v, but returned:%v", len(expected), len(actual))
	}

	for !(reflect.DeepEqual(expected, actual)) {
		t.Errorf("Expected:%v, but returned:%v", expected, actual)

	}
}

// ===========================================================
//  SLAChaincodeStub 업데이트 테스트
// ===========================================================

func TestUpdateContractId(t *testing.T) {

	var c = CreateStub()
	var s SimpleChaincode

	reg_args := []string{"SLA_REG_2007-01-00001", "SLA_REG_2007-01-00001,ITSM 도급계약,도급계약,2017.01-2017.12,신한은행,지은탁,010-1234-5678,2017.01.31,대상"}
	s.registerContract(c, reg_args)

	reg_args = []string{"SL_REG_2007-01-00002", "SLA_REG_2007-01-00002,ITSM 도급계약,도급계약,2017.01-2017.12,신한은행,지은탁,010-1234-5678,2017.01.31,대상"}
	s.registerContract(c, reg_args)

	reg_args = []string{"SLA_REG_2007-01-00003", "SLA_REG_2007-01-00003,ITSM 도급계약,도급계약,2017.01-2017.12,신한생명,지은탁,010-1234-5678,2017.01.31,대상"}
	s.registerContract(c, reg_args)

	args := []string{"SLA_REG_2007-01-00001", "SLA_REG_2007-01-00001,ITSM 도급계약,도급계약_UPDATED,2017.01-2017.12,신한은행,지은탁,010-1234-5678,2017.01.31,대상"}
	expected := "SLA_REG_2007-01-00001,ITSM 도급계약,도급계약_UPDATED,2017.01-2017.12,신한은행,지은탁,010-1234-5678,2017.01.31,대상"

	actual := s.updateContractId(c, args)

	if expected != actual {
		t.Errorf("Expected:%v, but returned:%v", expected, actual)
	}
}

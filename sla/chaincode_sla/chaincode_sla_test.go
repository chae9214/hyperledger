package main

// go test chaincode_sla_test.go chaincode_sla.go

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"testing"
	"time"
)

func checkInit(t *testing.T, stub *shim.MockStub, args []string) {
	_, err := stub.MockInit("1", "init", args)
	if err != nil {
		fmt.Println("Init failed", err)
		t.FailNow()
	}
}

func checkState(t *testing.T, stub *shim.MockStub, name string, value string) {
	bytes := stub.State[name]
	if bytes == nil {
		fmt.Println("State", name, "failed to get value")
		t.FailNow()
	}
	if string(bytes) != value {
		fmt.Printf("State value %v [%v] was not the expected value[%v]\n", name, string(bytes), value)
		t.FailNow()
	}
}

func checkQuery(t *testing.T, stub *shim.MockStub, fnName string, args []string, value string) {
	bytes, err := stub.MockQuery(fnName, args)
	if err != nil {
		fmt.Println("Query", fnName, "failed", err)
		t.FailNow()
	}
	if bytes == nil {
		fmt.Println("Query", fnName, "failed to get value")
		t.FailNow()
	}
	if string(bytes) != value {
		fmt.Printf("State value %v [%v] was not the expected value[%v]\n", fnName, string(bytes), value)
		t.FailNow()
	}
}

func checkInvoke(t *testing.T, stub *shim.MockStub, fnName string, args []string) {
	_, err := stub.MockInvoke("1", fnName, args)
	if err != nil {
		fmt.Println("Invoke", fnName, "failed", err)
		t.FailNow()
	}
}

func TestChaincodeSla_Init(t *testing.T) {
	scc := new(SimpleChaincode)
	stub := shim.NewMockStub("sla_chaincode", scc)

	checkInit(t, stub, []string{})

	checkState(t, stub, SLA_CONTRACT_ID_COUNT_KEY, "1")
	checkState(t, stub, SLA_EVALUATION_ID_COUNT_KEY, "1")
	checkState(t, stub, CURRENT_YEAR_KEY, "2017")
}

func TestChaincodeSla_Query_slaGetContractId(t *testing.T) {
	scc := new(SimpleChaincode)
	stub := shim.NewMockStub("sla_chaincode", scc)

	stub.MockTransactionStart("init")
	checkInit(t, stub, []string{})

	stub.MockTransactionStart("query")
	checkQuery(t, stub, "slaGetContractId", []string{}, "SLA_CONT_2017_00001")

	stub.MockTransactionStart("query")
	checkQuery(t, stub, "slaGetContractId", []string{}, "SLA_CONT_2017_00002")
}

// 최초 생성 + 임시 저장
//				[input]		[expected result]
// Progression: "" 	 	--> "TEMP"
// ----------------------------------------------------------------------------
// In KVS
// SLA_ALL_DATA						: "SLA_CONT_2017_00005|SLA_CONT_2017_00006"
func TestChaincodeSla_Invoke_slaCreateTempContract(t *testing.T) {
	inputContractContentInJson :=
		`{
  "RegId": "SLA_CONT_2017_00005",
  "Name": "홍길동",
  "Kind": "보통계약",
  "StaDate": "2017-02-01",
  "EndDate": "2017-12-01",
  "Client": "신한은행",
  "ClientPerson": "개인",
  "ClientPersonTel": "010-1111-2222",
  "AssessDate": "2017-12-31",
  "Progression": "",
  "AssessYn": "포함",
  "Approvals": [
    {
      "ApprovalUserId": "test",
      "ApprovalCompany": "test",
      "ApprovalDepartment": "test",
      "ApprovalName": "test",
      "ApprovalState": "test",
      "ApprovalDate": "test",
      "ApprovalComment": "test",
      "ApprovalAlarm": "test"
    },
    {
      "ApprovalUserId": "test2",
      "ApprovalCompany": "test2",
      "ApprovalDepartment": "test2",
      "ApprovalName": "test2",
      "ApprovalState": "test2",
      "ApprovalDate": "test2",
      "ApprovalComment": "test2",
      "ApprovalAlarm": "test2"
    }
  ],
  "ServiceItems": [
    {
      "ServiceItem": "test",
      "ScoreItem": "test",
      "MeasurementItem": "test",
      "ExplainItem": "test",
      "DivideScore": "test"
    }
  ]
}`

	expectedContractContentInJson :=
		`{
  "RegId": "SLA_CONT_2017_00005",
  "Name": "홍길동",
  "Kind": "보통계약",
  "StaDate": "2017-02-01",
  "EndDate": "2017-12-01",
  "Client": "신한은행",
  "ClientPerson": "개인",
  "ClientPersonTel": "010-1111-2222",
  "AssessDate": "2017-12-31",
  "Progression": "TEMP",
  "AssessYn": "포함",
  "Approvals": [
    {
      "ApprovalUserId": "test",
      "ApprovalCompany": "test",
      "ApprovalDepartment": "test",
      "ApprovalName": "test",
      "ApprovalState": "test",
      "ApprovalDate": "test",
      "ApprovalComment": "test",
      "ApprovalAlarm": "test"
    },
    {
      "ApprovalUserId": "test2",
      "ApprovalCompany": "test2",
      "ApprovalDepartment": "test2",
      "ApprovalName": "test2",
      "ApprovalState": "test2",
      "ApprovalDate": "test2",
      "ApprovalComment": "test2",
      "ApprovalAlarm": "test2"
    }
  ],
  "ServiceItems": [
    {
      "ServiceItem": "test",
      "ScoreItem": "test",
      "MeasurementItem": "test",
      "ExplainItem": "test",
      "DivideScore": "test"
    }
  ]
}`

	scc := new(SimpleChaincode)
	stub := shim.NewMockStub("sla_chaincode", scc)

	stub.MockTransactionStart("init")
	checkInit(t, stub, []string{})

	stub.MockTransactionStart("invoke")
	checkInvoke(t, stub, "slaCreateTempContract", []string{inputContractContentInJson})

	stub.MockTransactionStart("query")
	checkQuery(t, stub, "slaGetContractWithId", []string{"SLA_CONT_2017_00005"}, expectedContractContentInJson)
	checkState(t, stub, "SLA_ALL_DATA", "SLA_CONT_2017_00005")
}

// 최초 생성 + 바로 결제 요청
//				[input]		[expected result]
// Progression: "" 	 	--> "IN_PROGRESS_INTERNAL_REVIEW_REQUESTED"
// ----------------------------------------------------------------------------
// In KVS
// SLA_ALL_DATA						: "SLA_CONT_2017_00005|SLA_CONT_2017_00006"
func TestChaincodeSla_Invoke_slaSubmitContract(t *testing.T) {
	inputContractContentInJson :=
		`{
  "RegId": "SLA_CONT_2017_00005",
  "Name": "신한은행도급계약_201701",
  "Kind": "보통계약",
  "StaDate": "2017-02-01",
  "EndDate": "2017-12-01",
  "Client": "신한은행",
  "ClientPerson": "개인",
  "ClientPersonTel": "010-1111-2222",
  "AssessDate": "2017-12-31",
  "Progression": "",
  "AssessYn": "포함",
  "Approvals": [
    {
      "ApprovalUserId": "test",
      "ApprovalCompany": "test",
      "ApprovalDepartment": "test",
      "ApprovalName": "test",
      "ApprovalState": "test",
      "ApprovalDate": "test",
      "ApprovalComment": "test",
      "ApprovalAlarm": "test"
    },
    {
      "ApprovalUserId": "test2",
      "ApprovalCompany": "test2",
      "ApprovalDepartment": "test2",
      "ApprovalName": "test2",
      "ApprovalState": "test2",
      "ApprovalDate": "test2",
      "ApprovalComment": "test2",
      "ApprovalAlarm": "test2"
    }
  ],
  "ServiceItems": [
    {
      "ServiceItem": "test",
      "ScoreItem": "test",
      "MeasurementItem": "test",
      "ExplainItem": "test",
      "DivideScore": "test"
    }
  ]
}`

	expectedContractContentInJson :=
		`{
  "RegId": "SLA_CONT_2017_00005",
  "Name": "신한은행도급계약_201701",
  "Kind": "보통계약",
  "StaDate": "2017-02-01",
  "EndDate": "2017-12-01",
  "Client": "신한은행",
  "ClientPerson": "개인",
  "ClientPersonTel": "010-1111-2222",
  "AssessDate": "2017-12-31",
  "Progression": "IN_PROGRESS_INTERNAL_REVIEW_REQUESTED",
  "AssessYn": "포함",
  "Approvals": [
    {
      "ApprovalUserId": "test",
      "ApprovalCompany": "test",
      "ApprovalDepartment": "test",
      "ApprovalName": "test",
      "ApprovalState": "test",
      "ApprovalDate": "test",
      "ApprovalComment": "test",
      "ApprovalAlarm": "test"
    },
    {
      "ApprovalUserId": "test2",
      "ApprovalCompany": "test2",
      "ApprovalDepartment": "test2",
      "ApprovalName": "test2",
      "ApprovalState": "test2",
      "ApprovalDate": "test2",
      "ApprovalComment": "test2",
      "ApprovalAlarm": "test2"
    }
  ],
  "ServiceItems": [
    {
      "ServiceItem": "test",
      "ScoreItem": "test",
      "MeasurementItem": "test",
      "ExplainItem": "test",
      "DivideScore": "test"
    }
  ]
}`

	scc := new(SimpleChaincode)
	stub := shim.NewMockStub("sla_chaincode", scc)

	stub.MockTransactionStart("init")
	checkInit(t, stub, []string{})

	stub.MockTransactionStart("invoke")
	checkInvoke(t, stub, "slaSubmitContract", []string{inputContractContentInJson})

	stub.MockTransactionStart("query")
	checkQuery(t, stub, "slaGetContractWithId", []string{"SLA_CONT_2017_00005"}, expectedContractContentInJson)

	checkState(t, stub, "신한은행도급계약_201701", "SLA_CONT_2017_00005")
	checkState(t, stub, "신한은행", "SLA_CONT_2017_00005")
	checkState(t, stub, "SLA_ALL_DATA", "SLA_CONT_2017_00005")
}

// 최초 생성 + 바로 결제 요청
// Test1: 계약상태 변경 확인 -------------------------------------------------------
//									  [input]		[expected result]
// SLA_CONT_2017_00005	Progression : "" 	 	--> "IN_PROGRESS_INTERNAL_REVIEW_REQUESTED"
// SLA_CONT_2017_00006	Progression : "" 	 	--> "IN_PROGRESS_INTERNAL_REVIEW_REQUESTED"
//
// Test2: 기타 KVS 저장 내용 확인 ------------------------------------------------------
//					[Key]					[Value]
// 계약명 			"신한은행도급계약_201701" 	"SLA_CONT_2017_00005"
// 					"신한은행도급계약_201702" 	"SLA_CONT_2017_00006"
// 고객사명 			"신한은행", 				"SLA_CONT_2017_00005|SLA_CONT_2017_00006"
// 전체				SLA_ALL_DATA			"SLA_CONT_2017_00005|SLA_CONT_2017_00006"
// ----------------------------------------------------------------------------
func TestChaincodeSla_SLA_ALL_DATA_AfterTwoSlaSubmitContract(t *testing.T) {
	inputContractContentInJson_1 :=
		`{
  "RegId": "SLA_CONT_2017_00005",
  "Name": "신한은행도급계약_201701",
  "Kind": "보통계약",
  "StaDate": "2017-02-01",
  "EndDate": "2017-12-01",
  "Client": "신한은행",
  "ClientPerson": "개인",
  "ClientPersonTel": "010-1111-2222",
  "AssessDate": "2017-12-31",
  "Progression": "",
  "AssessYn": "포함",
  "Approvals": [
    {
      "ApprovalUserId": "test",
      "ApprovalCompany": "test",
      "ApprovalDepartment": "test",
      "ApprovalName": "test",
      "ApprovalState": "test",
      "ApprovalDate": "test",
      "ApprovalComment": "test",
      "ApprovalAlarm": "test"
    },
    {
      "ApprovalUserId": "test2",
      "ApprovalCompany": "test2",
      "ApprovalDepartment": "test2",
      "ApprovalName": "test2",
      "ApprovalState": "test2",
      "ApprovalDate": "test2",
      "ApprovalComment": "test2",
      "ApprovalAlarm": "test2"
    }
  ],
  "ServiceItems": [
    {
      "ServiceItem": "test",
      "ScoreItem": "test",
      "MeasurementItem": "test",
      "ExplainItem": "test",
      "DivideScore": "test"
    }
  ]
}`
	inputContractContentInJson_2 :=
		`{
  "RegId": "SLA_CONT_2017_00006",
  "Name": "신한은행도급계약_201702",
  "Kind": "보통계약",
  "StaDate": "2017-02-01",
  "EndDate": "2017-12-01",
  "Client": "신한은행",
  "ClientPerson": "개인",
  "ClientPersonTel": "010-1111-2222",
  "AssessDate": "2017-12-31",
  "Progression": "",
  "AssessYn": "포함",
  "Approvals": [
    {
      "ApprovalUserId": "test",
      "ApprovalCompany": "test",
      "ApprovalDepartment": "test",
      "ApprovalName": "test",
      "ApprovalState": "test",
      "ApprovalDate": "test",
      "ApprovalComment": "test",
      "ApprovalAlarm": "test"
    },
    {
      "ApprovalUserId": "test2",
      "ApprovalCompany": "test2",
      "ApprovalDepartment": "test2",
      "ApprovalName": "test2",
      "ApprovalState": "test2",
      "ApprovalDate": "test2",
      "ApprovalComment": "test2",
      "ApprovalAlarm": "test2"
    }
  ],
  "ServiceItems": [
    {
      "ServiceItem": "test",
      "ScoreItem": "test",
      "MeasurementItem": "test",
      "ExplainItem": "test",
      "DivideScore": "test"
    }
  ]
}`

	expectedContractContentInJson_1 :=
		`{
  "RegId": "SLA_CONT_2017_00005",
  "Name": "신한은행도급계약_201701",
  "Kind": "보통계약",
  "StaDate": "2017-02-01",
  "EndDate": "2017-12-01",
  "Client": "신한은행",
  "ClientPerson": "개인",
  "ClientPersonTel": "010-1111-2222",
  "AssessDate": "2017-12-31",
  "Progression": "IN_PROGRESS_INTERNAL_REVIEW_REQUESTED",
  "AssessYn": "포함",
  "Approvals": [
    {
      "ApprovalUserId": "test",
      "ApprovalCompany": "test",
      "ApprovalDepartment": "test",
      "ApprovalName": "test",
      "ApprovalState": "test",
      "ApprovalDate": "test",
      "ApprovalComment": "test",
      "ApprovalAlarm": "test"
    },
    {
      "ApprovalUserId": "test2",
      "ApprovalCompany": "test2",
      "ApprovalDepartment": "test2",
      "ApprovalName": "test2",
      "ApprovalState": "test2",
      "ApprovalDate": "test2",
      "ApprovalComment": "test2",
      "ApprovalAlarm": "test2"
    }
  ],
  "ServiceItems": [
    {
      "ServiceItem": "test",
      "ScoreItem": "test",
      "MeasurementItem": "test",
      "ExplainItem": "test",
      "DivideScore": "test"
    }
  ]
}`

	expectedContractContentInJson_2 :=
		`{
  "RegId": "SLA_CONT_2017_00006",
  "Name": "신한은행도급계약_201702",
  "Kind": "보통계약",
  "StaDate": "2017-02-01",
  "EndDate": "2017-12-01",
  "Client": "신한은행",
  "ClientPerson": "개인",
  "ClientPersonTel": "010-1111-2222",
  "AssessDate": "2017-12-31",
  "Progression": "IN_PROGRESS_INTERNAL_REVIEW_REQUESTED",
  "AssessYn": "포함",
  "Approvals": [
    {
      "ApprovalUserId": "test",
      "ApprovalCompany": "test",
      "ApprovalDepartment": "test",
      "ApprovalName": "test",
      "ApprovalState": "test",
      "ApprovalDate": "test",
      "ApprovalComment": "test",
      "ApprovalAlarm": "test"
    },
    {
      "ApprovalUserId": "test2",
      "ApprovalCompany": "test2",
      "ApprovalDepartment": "test2",
      "ApprovalName": "test2",
      "ApprovalState": "test2",
      "ApprovalDate": "test2",
      "ApprovalComment": "test2",
      "ApprovalAlarm": "test2"
    }
  ],
  "ServiceItems": [
    {
      "ServiceItem": "test",
      "ScoreItem": "test",
      "MeasurementItem": "test",
      "ExplainItem": "test",
      "DivideScore": "test"
    }
  ]
}`

	scc := new(SimpleChaincode)
	stub := shim.NewMockStub("sla_chaincode", scc)

	stub.MockTransactionStart("init")
	checkInit(t, stub, []string{})

	stub.MockTransactionStart("invoke")
	checkInvoke(t, stub, "slaSubmitContract", []string{inputContractContentInJson_1})

	stub.MockTransactionStart("query")
	checkQuery(t, stub, "slaGetContractWithId", []string{"SLA_CONT_2017_00005"}, expectedContractContentInJson_1)

	stub.MockTransactionStart("invoke")
	checkInvoke(t, stub, "slaSubmitContract", []string{inputContractContentInJson_2})

	stub.MockTransactionStart("query")
	checkQuery(t, stub, "slaGetContractWithId", []string{"SLA_CONT_2017_00006"}, expectedContractContentInJson_2)

	// KVS 확인
	checkState(t, stub, "신한은행도급계약_201701", "SLA_CONT_2017_00005")
	checkState(t, stub, "신한은행도급계약_201702", "SLA_CONT_2017_00006")
	checkState(t, stub, "신한은행", "SLA_CONT_2017_00005|SLA_CONT_2017_00006")
	checkState(t, stub, "SLA_ALL_DATA", "SLA_CONT_2017_00005|SLA_CONT_2017_00006")
}

// 계약을 업데이트 합니다.
// Test: 계약상태 변경 확인 -------------------------------------------------------
//										[input]				[expected result]
// SLA_CONT_2017_00010	StaDate : 		"2017-02-01" 	--> "2018-02-01"
// 	 					EndDate : 		"2017-12-01" 	--> "2018-12-01"
//						ClientPerson  :	"개인"			--> "개개인"
//
func TestChaincodeSla_Invoke_slaUpdateContract(t *testing.T) {

	inputContractContentInJson := `{
  "RegId": "SLA_CONT_2017_00005",
  "Name": "신한은행도급계약_201701",
  "Kind": "보통계약",
  "StaDate": "2017-02-01",
  "EndDate": "2017-12-01",
  "Client": "신한은행",
  "ClientPerson": "개인",
  "ClientPersonTel": "010-1111-2222",
  "AssessDate": "2017-12-31",
  "Progression": "IN_PROGRESS_INTERNAL_REVIEW_REQUESTED",
  "AssessYn": "포함",
  "Approvals": [
    {
      "ApprovalUserId": "test",
      "ApprovalCompany": "test",
      "ApprovalDepartment": "test",
      "ApprovalName": "test",
      "ApprovalState": "test",
      "ApprovalDate": "test",
      "ApprovalComment": "test",
      "ApprovalAlarm": "test"
    },
    {
      "ApprovalUserId": "test2",
      "ApprovalCompany": "test2",
      "ApprovalDepartment": "test2",
      "ApprovalName": "test2",
      "ApprovalState": "test2",
      "ApprovalDate": "test2",
      "ApprovalComment": "test2",
      "ApprovalAlarm": "test2"
    }
  ],
  "ServiceItems": [
    {
      "ServiceItem": "test",
      "ScoreItem": "test",
      "MeasurementItem": "test",
      "ExplainItem": "test",
      "DivideScore": "test"
    }
  ]
}`

	// Update 내용 적용
	updateContractContentInJson := `{
  "RegId": "SLA_CONT_2017_00005",
  "Name": "신한은행도급계약_201701",
  "Kind": "보통계약",
  "StaDate": "2018-02-01",
  "EndDate": "2018-12-01",
  "Client": "신한은행",
  "ClientPerson": "개개인",
  "ClientPersonTel": "010-1111-2222",
  "AssessDate": "2017-12-31",
  "Progression": "IN_PROGRESS_INTERNAL_REVIEW_REQUESTED",
  "AssessYn": "포함",
  "Approvals": [
    {
      "ApprovalUserId": "test",
      "ApprovalCompany": "test",
      "ApprovalDepartment": "test",
      "ApprovalName": "test",
      "ApprovalState": "test",
      "ApprovalDate": "test",
      "ApprovalComment": "test",
      "ApprovalAlarm": "test"
    },
    {
      "ApprovalUserId": "test2",
      "ApprovalCompany": "test2",
      "ApprovalDepartment": "test2",
      "ApprovalName": "test2",
      "ApprovalState": "test2",
      "ApprovalDate": "test2",
      "ApprovalComment": "test2",
      "ApprovalAlarm": "test2"
    }
  ],
  "ServiceItems": [
    {
      "ServiceItem": "test",
      "ScoreItem": "test",
      "MeasurementItem": "test",
      "ExplainItem": "test",
      "DivideScore": "test"
    }
  ]
}`

	scc := new(SimpleChaincode)
	stub := shim.NewMockStub("sla_chaincode", scc)

	checkInit(t, stub, []string{})
	checkInvoke(t, stub, "slaSubmitContract", []string{inputContractContentInJson})
	checkInvoke(t, stub, "slaUpdateContract", []string{updateContractContentInJson}) // todo update --> updated

	checkQuery(t, stub, "slaGetContractWithId", []string{"SLA_CONT_2017_00005"}, updateContractContentInJson) // todo update --> expected
}

//계약을 승인합니다
// Test1.: 계약상태 변경 확인 -------------------------------------------------------
// 계약 데이터: SLA_CONT_2017_00010
//												[input]											[expected result]
// 1. SubmitContract 호출			Progression : 	"TEMP"											"IN_PROGRESS_INTERNAL_REVIEW_REQUESTED"
// 2. ApprovalContract 호출 		Progression :   "IN_PROGRESS_INTERNAL_REVIEW_REQUESTED"			"IN_PROGRESS_CLIENT_REVIEW_REQUESTED"
// 3. ApprovalContract 재호출 	Progression :   "IN_PROGRESS_CLIENT_REVIEW_REQUESTED"			"IN_PROGRESS_CLIENT_MANAGER_REVIEW_REQUESTED"
// 4. ApprovalContract 재호출 	Progression :   "IN_PROGRESS_CLIENT_MANAGER_REVIEW_REQUESTED"	"CLOSED"
// Test2.: 개별 승인 (Approval) 변경 확인 -------------------------------------------
//
//
//
//
func TestChaincodeSla_Invoke_slaApproveContract(t *testing.T) {

	inputContractContentInJson :=
		`{
  "RegId": "SLA_CONT_2017_00010",
  "Name": "홍길동",
  "Kind": "보통계약",
  "StaDate": "2017-02-01",
  "EndDate": "2017-12-01",
  "Client": "신한은행",
  "ClientPerson": "개인",
  "ClientPersonTel": "010-1111-2222",
  "AssessDate": "2017-12-31",
  "Progression": "TEMP",
  "AssessYn": "포함",
  "Approvals": [
     {
      "ApprovalUserId": "기안자_A",
      "ApprovalCompany": "test",
      "ApprovalDepartment": "test",
      "ApprovalName": "test",
      "ApprovalState": "TEMP",
      "ApprovalDate": "test",
      "ApprovalComment": "test",
      "ApprovalAlarm": "test"
    },
    {
      "ApprovalUserId": "내부관리자_A",
      "ApprovalCompany": "test",
      "ApprovalDepartment": "test",
      "ApprovalName": "test",
      "ApprovalState": "TEMP",
      "ApprovalDate": "test",
      "ApprovalComment": "test",
      "ApprovalAlarm": "test"
    },
    {
      "ApprovalUserId": "고객_A",
      "ApprovalCompany": "test2",
      "ApprovalDepartment": "test2",
      "ApprovalName": "test2",
      "ApprovalState": "test2",
      "ApprovalDate": "test2",
      "ApprovalComment": "test2",
      "ApprovalAlarm": "test2"
    },
    {
      "ApprovalUserId": "고객관리자_A",
      "ApprovalCompany": "test2",
      "ApprovalDepartment": "test2",
      "ApprovalName": "test2",
      "ApprovalState": "test2",
      "ApprovalDate": "test2",
      "ApprovalComment": "test2",
      "ApprovalAlarm": "test2"
    }
  ],
  "ServiceItems": [
    {
      "ServiceItem": "test",
      "ScoreItem": "test",
      "MeasurementItem": "test",
      "ExplainItem": "test",
      "DivideScore": "test"
    }
  ]
}`

	approvalDate := time.Now().Format("2006-01-02")
	// 예상 결과갑
	expectedContractContentInJson :=
		`{
  "RegId": "SLA_CONT_2017_00010",
  "Name": "홍길동",
  "Kind": "보통계약",
  "StaDate": "2017-02-01",
  "EndDate": "2017-12-01",
  "Client": "신한은행",
  "ClientPerson": "개인",
  "ClientPersonTel": "010-1111-2222",
  "AssessDate": "2017-12-31",
  "Progression": "IN_PROGRESS_CLIENT_REVIEW_REQUESTED",
  "AssessYn": "포함",
  "Approvals": [
    {
      "ApprovalUserId": "기안자_A",
      "ApprovalCompany": "test",
      "ApprovalDepartment": "test",
      "ApprovalName": "test",
      "ApprovalState": "TEMP",
      "ApprovalDate": "test",
      "ApprovalComment": "test",
      "ApprovalAlarm": "test"
    },
    {
      "ApprovalUserId": "내부관리자_A",
      "ApprovalCompany": "test",
      "ApprovalDepartment": "test",
      "ApprovalName": "test",
      "ApprovalState": "APPROVED",
      "ApprovalDate": "` + approvalDate + `",` + "\n" +
			`      "ApprovalComment": "내용확인 하였음",
      "ApprovalAlarm": "test"
    },
    {
      "ApprovalUserId": "고객_A",
      "ApprovalCompany": "test2",
      "ApprovalDepartment": "test2",
      "ApprovalName": "test2",
      "ApprovalState": "test2",
      "ApprovalDate": "test2",
      "ApprovalComment": "test2",
      "ApprovalAlarm": "test2"
    },
    {
      "ApprovalUserId": "고객관리자_A",
      "ApprovalCompany": "test2",
      "ApprovalDepartment": "test2",
      "ApprovalName": "test2",
      "ApprovalState": "test2",
      "ApprovalDate": "test2",
      "ApprovalComment": "test2",
      "ApprovalAlarm": "test2"
    }
  ],
  "ServiceItems": [
    {
      "ServiceItem": "test",
      "ScoreItem": "test",
      "MeasurementItem": "test",
      "ExplainItem": "test",
      "DivideScore": "test"
    }
  ]
}`
	scc := new(SimpleChaincode)
	stub := shim.NewMockStub("sla_chaincode", scc)

	stub.MockTransactionStart("init")
	checkInit(t, stub, []string{})

	stub.MockTransactionStart("invoke")
	checkInvoke(t, stub, "slaSubmitContract", []string{inputContractContentInJson})

	// 승인 테스트 Input
	SlaContractRegId := "SLA_CONT_2017_00010"
	SlaContractApprovalUserId := "내부관리자_A"
	SlaContractApprovalComment := "내용확인 하였음"

	stub.MockTransactionStart("invoke")
	checkInvoke(t, stub, "slaApproveContract", []string{SlaContractRegId, SlaContractApprovalUserId, SlaContractApprovalComment})

	stub.MockTransactionStart("query")
	checkQuery(t, stub, "slaGetContractWithId", []string{"SLA_CONT_2017_00010"}, expectedContractContentInJson)
}

// 4.계약을 반려합니다.
func TestChaincodeSla_Invoke_slaRejectContract(t *testing.T) {

	inputContractContentInJson :=
		`{
  "RegId": "SLA_CONT_2017_00010",
  "Name": "홍길동",
  "Kind": "보통계약",
  "StaDate": "2017-02-01",
  "EndDate": "2017-12-01",
  "Client": "신한은행",
  "ClientPerson": "개인",
  "ClientPersonTel": "010-1111-2222",
  "AssessDate": "2017-12-31",
  "Progression": "TEMP",
  "AssessYn": "포함",
  "Approvals": [
     {
      "ApprovalUserId": "기안자_A",
      "ApprovalCompany": "test",
      "ApprovalDepartment": "test",
      "ApprovalName": "test",
      "ApprovalState": "TEMP",
      "ApprovalDate": "test",
      "ApprovalComment": "test",
      "ApprovalAlarm": "test"
    },
    {
      "ApprovalUserId": "내부관리자_A",
      "ApprovalCompany": "test",
      "ApprovalDepartment": "test",
      "ApprovalName": "test",
      "ApprovalState": "TEMP",
      "ApprovalDate": "test",
      "ApprovalComment": "test",
      "ApprovalAlarm": "test"
    },
    {
      "ApprovalUserId": "고객_A",
      "ApprovalCompany": "test2",
      "ApprovalDepartment": "test2",
      "ApprovalName": "test2",
      "ApprovalState": "test2",
      "ApprovalDate": "test2",
      "ApprovalComment": "test2",
      "ApprovalAlarm": "test2"
    },
    {
      "ApprovalUserId": "고객관리자_A",
      "ApprovalCompany": "test2",
      "ApprovalDepartment": "test2",
      "ApprovalName": "test2",
      "ApprovalState": "test2",
      "ApprovalDate": "test2",
      "ApprovalComment": "test2",
      "ApprovalAlarm": "test2"
    }
  ],
  "ServiceItems": [
    {
      "ServiceItem": "test",
      "ScoreItem": "test",
      "MeasurementItem": "test",
      "ExplainItem": "test",
      "DivideScore": "test"
    }
  ]
}`

	approvalDate := time.Now().Format("2006-01-02")
	// 예상 결과갑
	expectedContractContentInJson :=
		`{
  "RegId": "SLA_CONT_2017_00010",
  "Name": "홍길동",
  "Kind": "보통계약",
  "StaDate": "2017-02-01",
  "EndDate": "2017-12-01",
  "Client": "신한은행",
  "ClientPerson": "개인",
  "ClientPersonTel": "010-1111-2222",
  "AssessDate": "2017-12-31",
  "Progression": "IN_PROGRESS_INTERNAL_REVIEW_REQUESTED",
  "AssessYn": "포함",
  "Approvals": [
    {
      "ApprovalUserId": "기안자_A",
      "ApprovalCompany": "test",
      "ApprovalDepartment": "test",
      "ApprovalName": "test",
      "ApprovalState": "TEMP",
      "ApprovalDate": "test",
      "ApprovalComment": "test",
      "ApprovalAlarm": "test"
    },
    {
      "ApprovalUserId": "내부관리자_A",
      "ApprovalCompany": "test",
      "ApprovalDepartment": "test",
      "ApprovalName": "test",
      "ApprovalState": "REJECTED",
      "ApprovalDate": "` + approvalDate + `",` + "\n" +
			`      "ApprovalComment": "추가내용 필요함",
      "ApprovalAlarm": "test"
    },
    {
      "ApprovalUserId": "고객_A",
      "ApprovalCompany": "test2",
      "ApprovalDepartment": "test2",
      "ApprovalName": "test2",
      "ApprovalState": "test2",
      "ApprovalDate": "test2",
      "ApprovalComment": "test2",
      "ApprovalAlarm": "test2"
    },
    {
      "ApprovalUserId": "고객관리자_A",
      "ApprovalCompany": "test2",
      "ApprovalDepartment": "test2",
      "ApprovalName": "test2",
      "ApprovalState": "test2",
      "ApprovalDate": "test2",
      "ApprovalComment": "test2",
      "ApprovalAlarm": "test2"
    }
  ],
  "ServiceItems": [
    {
      "ServiceItem": "test",
      "ScoreItem": "test",
      "MeasurementItem": "test",
      "ExplainItem": "test",
      "DivideScore": "test"
    }
  ]
}`

	scc := new(SimpleChaincode)
	stub := shim.NewMockStub("sla_chaincode", scc)

	stub.MockTransactionStart("init")
	checkInit(t, stub, []string{})

	stub.MockTransactionStart("invoke")
	checkInvoke(t, stub, "slaSubmitContract", []string{inputContractContentInJson})

	// 승인 테스트
	SlaContractRegId := "SLA_CONT_2017_00010"
	SlaContractApprovalUserId := "내부관리자_A"
	SlaContractApprovalComment := "추가내용 필요함"

	stub.MockTransactionStart("invoke")
	checkInvoke(t, stub, "slaRejectContract", []string{SlaContractRegId, SlaContractApprovalUserId, SlaContractApprovalComment})

	stub.MockTransactionStart("query")
	checkQuery(t, stub, "slaGetContractWithId", []string{"SLA_CONT_2017_00010"}, expectedContractContentInJson)
}

// 5.계약을 최종 승인합니다.
// func TestChaincodeSla_Invoke_slaCloseContract(t *testing.T) {

// 	contractContentInJson := `{
//                   "RegId": "SLA_CONT_2017_00001",
//                   "Name": "홍길동",
//                   "Kind": "보통계약",
//                   "StaDate": "2017-02-01",
//                   "EndDate": "2017-12-01",
//                   "Client": "신한은행",
//                   "ClientPerson": "개인",
//                   "ClientPersonTel": "010-1111-2222",
//                   "AssessDate": "2017-12-31",
//                   "Progression": "작성",
//                   "AssessYn": "포함",
//                   "Approvals": [
//                     {
//                       "ApprovalUserId": "test",
//                       "ApprovalCompany": "test",
//                       "ApprovalDepartment": "test",
//                       "ApprovalName": "test",
//                       "ApprovalState": "test",
//                       "ApprovalDate": "test",
//                       "ApprovalComment": "test",
//                       "ApprovalAlarm": "test"
//                     },
//                     {
//                       "ApprovalUserId": "test2",
//                       "ApprovalCompany": "test2",
//                       "ApprovalDepartment": "test2",
//                       "ApprovalName": "test2",
//                       "ApprovalState": "test2",
//                       "ApprovalDate": "test2",
//                       "ApprovalComment": "test2",
//                       "ApprovalAlarm": "test2"
//                     }
//                   ],
//                   "ServiceItems": [
//                     {
//                       "ServiceItem": "test",
//                       "ScoreItem": "test",
//                       "MeasurementItem": "test",
//                       "ExplainItem": "test",
//                       "DivideScore": "test"
//                     }
//                   ]
//                 }`

// 	// Update 내용 적용
// 	updateContractContentInJson := `{"RegId":"SLA_CONT_2017_00001","Name":"홍길동","Kind":"보통계약","StaDate":"2017-02-01","EndDate":"2017-12-01","Client":"신한은행","ClientPerson":"개인","ClientPersonTel":"010-1111-2222","AssessDate":"2017-12-31","Progression":"CLOSED","AssessYn":"포함","Approvals":[{"ApprovalUserId":"20170101","ApprovalCompany":"test","ApprovalDepartment":"test","ApprovalName":"test","ApprovalState":"REJECTED","ApprovalDate":"2017-02-17","ApprovalComment":"반려 하였음","ApprovalAlarm":"test"},{"ApprovalUserId":"test2","ApprovalCompany":"test2","ApprovalDepartment":"test2","ApprovalName":"test2","ApprovalState":"test2","ApprovalDate":"test2","ApprovalComment":"test2","ApprovalAlarm":"test2"}],"ServiceItems":[{"ServiceItem":"test","ScoreItem":"test","MeasurementItem":"test","ExplainItem":"test","DivideScore":"test"}]}`

// 	SlaContractRegId := "SLA_CONT_2017_00001"
// 	SlaContractApprovalUserId := "20170101"
// 	SlaContractApprovalComment := "내용확인 하였음"
// 	SlaContractProgression := "IN_PROGRESS_CLIENT_MANAGER_REVIEW_REQUESTED"

// 	scc := new(SimpleChaincode)
// 	stub := shim.NewMockStub("sla_chaincode", scc)

// 	stub.MockTransactionStart("init")
// 	checkInit(t, stub, []string{})

// 	stub.MockTransactionStart("invoke")
// 	checkInvoke(t, stub, "slaCreateContract", []string{contractContentInJson})

// 	stub.MockTransactionStart("invoke")
// 	checkInvoke(t, stub, "slaRejectContract", []string{SlaContractRegId, SlaContractApprovalUserId, SlaContractApprovalComment, SlaContractProgression})

// 	stub.MockTransactionStart("query")
// 	checkQuery(t, stub, "slaGetContractWithId", []string{"SLA_CONT_2017_00001"}, updateContractContentInJson)

// }

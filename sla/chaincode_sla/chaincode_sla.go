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
//chaincode_example05 show's how chaincode ID can be passed in as a paramete1222221`1111``1`1`1``1`1``1`1`1``11r instead of
//hard-coding.

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"strconv"
	"strings"
	"time"
)

// ===========================================================
//  Struct 및 Constant 정의
// ===========================================================
type SimpleChaincode struct {
}

// Sla Contract 구조체를 설정합니다.
type SlaContract struct {
	RegId           string           `json:  "RegId"`           // SLA계약등록번호
	Name            string           `json:  "Name"`            // SLA계약명
	Kind            string           `json:  "Kind"`            // SLA계약종류
	StaDate         string           `json:  "StaDate"`         // SLA계약시작일
	EndDate         string           `json:  "EndDate"`         // SLA계약종료일
	Client          string           `json:  "Client"`          // 고객사명
	ClientPerson    string           `json:  "ClientPerson"`    // 고객담당자명
	ClientPersonTel string           `json:  "ClientPersonTel"` // 고객담당자전화번호
	AssessDate      string           `json:  "AssessDate"`      // 평가예정일
	Progression     string           `json:  "Progression"`     // 진행단계
	AssessYn        string           `json:  "AssessYn"`        // SLA평가 대상여부
	Approvals       []SlaApproval    `json:  "Approvals"`       // SLA결재선정보
	ServiceItems    []SlaServiceItem `json:  "ServiceItems"`    // SLA평가항목
}

// Sla Approval 구조체를 설정합니다.
type SlaApproval struct {
	ApprovalUserId     string `json:  "ApprovalUserId"`     // 결재사용자ID
	ApprovalCompany    string `json:  "ApprovalCompany"`    // 결재회사명
	ApprovalDepartment string `json:  "ApprovalDepartment"` // 결재부서명
	ApprovalName       string `json:  "ApprovalName"`       // 결재자명
	ApprovalState      string `json:  "ApprovalState"`      // 결재상태
	ApprovalDate       string `json:  "ApprovalDate"`       // 결재일자
	ApprovalComment    string `json:  "ApprovalComment"`    // 의견내용
	ApprovalAlram      string `json:  "ApprovalAlram"`      // 알람여부  TODO Alram --> Alarm
}

// Sla ServiceItem 구조체를 설정합니다.
type SlaServiceItem struct {
	ServiceItem     string `json:  "ServiceItem"`     // 서비스항목
	ScoreItem       string `json:  "ScoreItem"`       // 평가항목
	MeasurementItem string `json:  "MeasurementItem"` // 측정기준
	ExplainItem     string `json:  "ExplainItem"`     // 설명
	DivideScore     string `json:  "DivideScore"`     // SLA배분점수
}

// Sla EvaluationRoot구조체를 설정합니다.
type SlaEvalutionRoot struct {
	RegId        string           `json:  "RegId"`        // SLA계약등록번호
	ContractId   string           `json:  "ContractId"`   // SLA계약명
	Status       string           `json:  "Status"`       // SLA계약명
	Evaluations  []SlaEvaluation  `json:  "Evaluations"`  // SLA평가등록번호
	ServiceItems []SlaServiceItem `json:  "ServiceItems"` // SLA평가항목
}

// Sla Evaluation 구조체를 설정합니다.
type SlaEvaluation struct {
	RegId                 string        `json:  "SlaContractRegId"` // SLA계약등록번호
	EvaluationRootId      string        `json:  "SlaContractName"`  // SLA계약명
	ScoresForServiceItems string        `json:  "SlaContractName"`  // SLA평가점수항목
	Approvals             []SlaApproval `json:  "Approvals"`        // SLA결재선정보
}

// key-value store 의 키 구분자
const FIELDSEP = "|"
const ENTRYSEP = ","
const SLA_ALL_DATA = "SLA_ALL_DATA"

const CONTRACT_ID_PREFIX = "SLA_CONT_"
const EVALUATION_ID_PREFIX = "SLA_EVAL_"

const SLA_CONTRACT_ID_COUNT_KEY = "SLA_CONTRACT_ID_COUNT"
const SLA_EVALUATION_ID_COUNT_KEY = "SLA_EVALUATION_ID_COUNT"

const CURRENT_YEAR_KEY = "CURRENT_YEAR"

// ===========================================================
//  function 정의 함수
// ===========================================================

func padLeft(str string) string {

	length := 5
	pad := "0"

	for {
		if len(str) >= length {
			return str
		}
		str = pad + str
	}
}

// ===========================================================
//  Initialization 함수
// ===========================================================

// 초기화를 처리합니다.
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	var err error
	year, _, _ := time.Now().Date()

	err = stub.PutState(SLA_CONTRACT_ID_COUNT_KEY, []byte(strconv.Itoa(1)))
	err = stub.PutState(SLA_EVALUATION_ID_COUNT_KEY, []byte(strconv.Itoa(1)))
	err = stub.PutState(CURRENT_YEAR_KEY, []byte(strconv.Itoa(year))) // 현재 year
	return nil, err
}

// 기능 이벤트를 호출합니다.
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	switch function {

	// 계약 상태 변경
	case "slaCreateContract": // 요청자가 계약 생성 (최초 생성 및 임시저장)
		return t.slaCreateContract(stub, args)

	case "slaUpdateContract": // 요청자가 계약 수정 (임시저장 후 / 승인거절 후)
		return t.slaUpdateContract(stub, args)

	case "slaUpdateContract1": // 요청자가 계약 수정 (임시저장 후 / 승인거절 후)
		return t.slaUpdateContract1(stub, args)

	case "slaAbandonContract": // 요청자가 계약 폐기 (임시저장 후 / 승인거절 후)
		return t.slaUpdateContract(stub, args)

	// 결재 요청 / 승인 / 거절
	case "slaSubmitContract": // 요청자 --> 결재자  (최초 생성 후 결재 요청 / 수정  후 결재 요청 )
		return t.slaSubmitContract(stub, args)

	case "slaApproveContract": // 결재자 --> 다음 결재자
		return t.slaApproveContract(stub, args)

	case "slaRejectContract": // 결재자 --> 요청자
		return t.slaRejectContract(stub, args)

	// 최종 승인
	case "slaCloseContract": // 최종 결재자 승인
		return t.slaCloseContract(stub, args)

	// 전체 평가 생성
	case "slaCreateEvaluationTemplateFromContract": // 최초 평가 생성 (계약등록 최종 승인 후),
		_, err1 := t.slaCreateEvaluationRootFromContract(stub, args)
		_, err2 := t.slaCreateEvaluationsFromContract(stub, args)

		if err1 == nil {
			err1 = err2
		}
		return nil, err1

	// 개별 평가 진행
	case "slaInitEvaluationValues": // 개별 평가의 평가점수 입력
		return t.slaInitEvaluationValues(stub, args)

	case "slaUpdateEvaluationValues": // 평가점수 수정
		return t.slaUpdateEvaluationValues(stub, args)

	// 결재 요청 / 승인 / 거절
	case "slaSubmitEvaluation": // 요청자 --> 결재자  (최초 생성 후 결재 요청 / 수정  후 결재 요청 )
		return t.slaSubmitEvaluation(stub, args)

	case "slaApproveEvaluation": // 결재자 --> 다음 결재자
		return t.slaApproveEvaluation(stub, args)

	case "slaRejectEvaluation": // 결재자 --> 요청자
		return t.slaRejectEvaluation(stub, args)

	// 지급 요청 / 승인 / 거절
	case "slaSubmitPayment": // 요청자 --> 지급자
		return t.slaSubmitPayment(stub, args)

	case "slaClosePayment": // 지급자
		return t.slaClosePayment(stub, args)

	// 개별 평가 마무리
	case "slaCloseEvaluation":
		return t.slaCloseEvaluation(stub, args)

	// 전체 평가 마무리
	case "slaCloseEvaluationRoot": // 마지막 개별 평가가 마무리될 경우, 자동 호출
		return t.slaCloseEvaluationRoot(stub, args)

	}
	return nil, errors.New("Invalid invoke function name. Expecting \"slaCreateContract\" \"slaUpdateContract\" \"slaApproveContract\" \"slaRejectContract\"")
}

// 쿼리 이벤트를 처리합니다.
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	switch function {

	case "slaGetContractId":
		return t.slaGetContractId(stub, args)

	case "slaGetAllContracts":
		return t.slaGetAllContracts(stub, args)

	case "slaGetContractWithId":
		return t.slaGetContractWithId(stub, args)

	case "slaGetContractsWithName":
		return t.slaGetContractsWithName(stub, args)

	case "slaGetContractsWithClient":
		return t.slaGetContractsWithClient(stub, args)

	case "slaGetEvaluationId":
		return t.slaGetEvaluationId(stub, args)

	case "slaGetAllEvaluations":
		return t.slaGetAllContracts(stub, args)

	case "slaGetEvaluationWithId":
		return t.slaGetContractWithId(stub, args)

	case "slaGetEvaluationsWithName":
		return t.slaGetContractsWithName(stub, args)

	case "slaGetEvaluationsWithClient":
		return t.slaGetContractsWithClient(stub, args)

	}
	return nil, errors.New("Invalid Query function name. Expecting \"????slaGetAllContracts\" \"slaGetContractWithId\" \"slaGetContractsWithName\" \"slaGetContractsWithClient\"")
}

// 메인함수를 처리합니다.
func main() {

	// 블록체인 이벤트를 호출합니다.
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Er4ror starting Simple chaincode: %s", err)
	}
}

// ===========================================================
//  SLAChaincodeStub 등록 함수
// ===========================================================

func (t *SimpleChaincode) slaGetContractId(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	var contractId string
	var err error

	// 1.저장된 해당 계약 카운트를 호출
	currentCountInBytes, err := stub.GetState(SLA_CONTRACT_ID_COUNT_KEY)
	if err != nil {
		return nil, errors.New("Failed to SLA_CONTRACT_ID_COUNT_KEY with " + SLA_CONTRACT_ID_COUNT_KEY)
	}

	if currentCountInBytes == nil { // if not initialized
		err = stub.PutState(SLA_CONTRACT_ID_COUNT_KEY, []byte(strconv.Itoa(1)))
	}
	currentCount, _ := strconv.Atoi(string(currentCountInBytes))

	// 2.새로운 연도일 경우, 계약 카운트를 초기화
	kvsCurrentYearInBytes, err := stub.GetState(CURRENT_YEAR_KEY)
	if err != nil {
		return nil, errors.New("Failed to kvsCurrentYearInBytes with " + CURRENT_YEAR_KEY)
	}

	currentYear, _, _ := time.Now().Date()
	if kvsCurrentYearInBytes == nil { // if not initialized
		err = stub.PutState(CURRENT_YEAR_KEY, []byte(strconv.Itoa(currentYear)))
	}

	if string(kvsCurrentYearInBytes) != strconv.Itoa(currentYear) { // new year starts
		err = stub.PutState(CURRENT_YEAR_KEY, []byte(strconv.Itoa(currentYear))) // 새로운 현재 연도
		err = stub.PutState(SLA_CONTRACT_ID_COUNT_KEY, []byte(strconv.Itoa(1)))  // 카운트는 1부터
	}

	// 3. 계약번호 채번을 생성합니다.
	contractId = CONTRACT_ID_PREFIX + strconv.Itoa(currentYear) + "_" + padLeft(strconv.Itoa(currentCount))

	// 4. 다음 계약번호 카운트를 저장
	nextCount := currentCount + 1
	stub.PutState(SLA_CONTRACT_ID_COUNT_KEY, []byte(strconv.Itoa(nextCount)))

	return []byte(contractId), nil
}

// 계약을 등록합니다.
func (t *SimpleChaincode) slaCreateContract(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	var err error
	var data SlaContract

	content := args[0]

	fmt.Printf("slaCreateContract Input Args:%s\n", args[0])

	// JSON 데이터를 디코딩(Unmarshal)합니다.
	err = json.Unmarshal([]byte(content), &data)
	if err != nil {
		return nil, errors.New("Failed to registerContractByIdToJSON with " + content)
	}

	// JSON 데이터를 정렬하여 디코딩(Unmarshal)합니다.
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return nil, errors.New("Failed to registerContractByIdToJSON with " + content)
	}

	fmt.Printf("= jsonData ====================================================\n")
	fmt.Println(data.RegId)
	fmt.Println(data.Name)
	fmt.Println(data.Client)
	fmt.Println(SLA_ALL_DATA)
	fmt.Println("")
	fmt.Println(string(jsonData))
	fmt.Printf("===============================================================\n")

	contractID := data.RegId
	contractName := data.Name
	contractClient := data.Client

	// A01. 계약ID 등록합니다.
	err = stub.PutState(data.RegId, []byte(content))

	if err != nil {
		return nil, errors.New("Failed to put state with" + content)

	} else {
		fmt.Println("SlaContractRegId : ok")
	}

	// A02. 계약명 등록합니다.
	{
		var err error

		// 계약명으로 기존내역를 조회합니다.
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
			err = stub.PutState(contractName, []byte(contractIDsInString+FIELDSEP+contractID))
			if err != nil {
				return nil, err
			}
		}
	}

	// A03. 고객사명 등록합니다.
	{
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
			err = stub.PutState(contractClient, []byte(contractIDsInString+FIELDSEP+contractID))
			if err != nil {
				return nil, err
			}
		}
	}

	// A04. 전체 조회 등록합니다.
	{
		var err error

		// 데이터를 전체 조회합니다.
		contractALLIDsInBytes, err := stub.GetState(SLA_ALL_DATA) // 리턴값 ([]byte, error)

		fmt.Println("= A04 == 01:" + string(contractALLIDsInBytes))

		if err != nil {
			return nil, errors.New("Failed to get state with" + string(contractALLIDsInBytes))
		}

		fmt.Println("= A04 == 02:" + string(contractALLIDsInBytes))

		contractALLIDsInString := string(contractALLIDsInBytes)

		fmt.Println("= A04 == 03:" + contractALLIDsInString)

		// 기존내역이 없을경우 "계약명"-"계약ID목록" 등록
		if contractALLIDsInString == "" {
			err = stub.PutState(SLA_ALL_DATA, []byte(contractID))
			fmt.Println("= A04 == 04:")

			if err != nil {
				return nil, err
			}

		} else {
			err = stub.PutState(SLA_ALL_DATA, []byte(contractALLIDsInString+FIELDSEP+contractID))

			fmt.Println("= A04 == 05:")

			if err != nil {
				return nil, err
			}
		}
	}

	/*
		err = stub.PutState(data[0].SlaContractName, []byte(content))
		if err != nil {
			return nil, errors.New("Failed to put state with" + content)
		} else {
		    fmt.Println("SlaContractName : ok")
		}

		err = stub.PutState(data[0].SlaContractClient, []byte(content))
		if err != nil {
			return nil, errors.New("Failed to put state with" + content)
		} else {
		    fmt.Println("SlaContractClient : ok")
		}

		err = stub.PutState(SLA_ALL_DATA, []byte(content))
		if err != nil {
			return nil, errors.New("Failed to put state with" + content)
		} else {
		    fmt.Println("SLA_ALL_DATA : ok")
		}

	*/
	return nil, nil
}

// // ===========================================================
// //  SLAChaincodeStub 업데이트 함수
// // ===========================================================

// // 계약을 업데이트합니다. (기본)
// func (t *SimpleChaincode) updateContractId(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

// 	var dataInBytes string
// 	var err error

// 	fmt.Printf("updateContractId Input Args:%s\n", args[0])

// 	if len(args) != 1 {
// 		return nil, errors.New("Incorrect number of arguments. Expecting name of the value to slaGetContractsWithClient")
// 	}

// 	dataInBytes = args[0]
// 	contractID := args[0]

// 	// 기존내용 조회
// 	contractIDsInBytes, err := stub.GetState(contractID)
// 	if err != nil {
// 		return nil, errors.New("Failed to get state with " + string(contractIDsInBytes))
// 	}

// 	// UPDATDE 처리
// 	stub.PutState(contractID, []byte(args[1]))

// 	// 변경내용 조회
// 	update_value, err := stub.GetState(contractID)
// 	if err != nil {
// 		return nil, errors.New("Failed to get state with " + dataInBytes)
// 	}

// 	fmt.Printf("slaGetContractsWithClient Response:%s\n", update_value)

// 	return []byte(update_value), nil
// }

// 1.계약을 업데이트 합니다.
func (t *SimpleChaincode) slaUpdateContract(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	var err error
	var data SlaContract

	content := args[0]

	fmt.Printf("slaCreateContract Input Args:%s\n", args[0])

	// JSON 데이터를 디코딩(Unmarshal)합니다.
	err = json.Unmarshal([]byte(content), &data)
	if err != nil {
		return nil, errors.New("Failed to registerContractByIdToJSON with " + content)
	}
	/*
		// 임시저장 데이터를 처리합니다.

		data.Name = "김아무개"
		data.Kind = "특수"
		data.StaDate = "2018-01-01"
		data.EndDate = "2018-01-01"
		data.Client = "신한은행(Shinhan.com)"
		data.ClientPerson = "신한아무개"
		data.ClientPersonTel = "010-9999-9999"
		data.AssessDate = "2018-12-12"
		data.Progression = "진행"
		data.AssessYn = "Y"

		// JSON 데이터를 정렬하여 디코딩(Unmarshal)합니다.
		jsonData, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			return nil, errors.New("Failed to registerContractByIdToJSON with " + content)
		}

		// JSON 데이터를 디코딩(Unmarshal)합니다.
		err = json.Unmarshal([]byte(content), &data)
	*/
	fmt.Printf("= contractListBytes contractListBytes ====================================================\n")
	contractListStirng, _ := json.Marshal(data)

	fmt.Println(string(contractListStirng))

	fmt.Printf("= err ====================================================\n")
	fmt.Println(data.RegId)

	fmt.Printf("= jsonData ====================================================\n")
	fmt.Println(data.RegId)
	fmt.Println(data.Name)
	fmt.Println(data.Client)
	fmt.Println(SLA_ALL_DATA)
	fmt.Println("")

	fmt.Printf("===============================================================\n")

	// 계약ID 통해 데이터를 업데이트 합니다.
	err = stub.PutState(data.RegId, []byte(content))

	if err != nil {
		return nil, errors.New("Failed to put state with" + content)
	}

	return nil, nil
}

// 1.계약을 업데이트 합니다.
func (t *SimpleChaincode) slaUpdateContract1(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	// create 와 유사
	var err error
	var data SlaContract

	content := args[0]

	fmt.Printf("slaCreateContract Input Args:%s\n", args[0])

	// JSON 데이터를 디코딩(Unmarshal)합니다.
	err = json.Unmarshal([]byte(content), &data)
	if err != nil {
		return nil, errors.New("Failed to registerContractByIdToJSON with " + content)
	}

	// 임시저장 데이터를 처리합니다.
	data.Name = "김아무개"
	data.Kind = "특수"
	data.StaDate = "2018-01-01"
	data.EndDate = "2018-01-01"
	data.Client = "신한은행(Shinhan.com)"
	data.ClientPerson = "신한아무개"
	data.ClientPersonTel = "010-9999-9999"
	data.AssessDate = "2018-12-12"
	data.Progression = "진행"
	data.AssessYn = "Y"

	// JSON 데이터를 정렬하여 디코딩(Unmarshal)합니다.
	//	jsonData, err := json.MarshalIndent(data, "", "  ")
	//	if err != nil {
	//		return nil, errors.New("Failed to registerContractByIdToJSON with " + content)
	//	}

	// JSON 데이터를 디코딩(Unmarshal)합니다.
	//err = json.Unmarshal([]byte(content), &data)

	fmt.Printf("= contractListBytes contractListBytes ====================================================\n")
	contractListStirng, _ := json.Marshal(data)

	fmt.Println(string(contractListStirng))

	fmt.Printf("= err ====================================================\n")
	fmt.Println(data.RegId)

	fmt.Printf("= jsonData ====================================================\n")
	fmt.Println(data.RegId)
	fmt.Println(data.Name)
	fmt.Println(data.Client)
	fmt.Println(SLA_ALL_DATA)
	fmt.Println("")

	fmt.Printf("===============================================================\n")

	contractID := data.RegId
	contractName := data.Name
	contractClient := data.Client

	// A01. 계약ID 등록합니다.
	err = stub.PutState(data.RegId, []byte(contractListStirng))

	if err != nil {
		return nil, errors.New("Failed to put state with" + content)

	} else {
		fmt.Println("SlaContractRegId : ok")
	}

	// A02. 계약명 등록합니다.
	{
		var err error

		// 계약명으로 기존내역를 조회합니다.
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
			err = stub.PutState(contractName, []byte(contractIDsInString+FIELDSEP+contractID))
			if err != nil {
				return nil, err
			}
		}
	}

	// A03. 고객사명 등록합니다.
	{
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
			err = stub.PutState(contractClient, []byte(contractIDsInString+FIELDSEP+contractID))
			if err != nil {
				return nil, err
			}
		}
	}

	// A04. 전체 조회 등록합니다.
	{
		var err error

		// 데이터를 전체 조회합니다.
		contractALLIDsInBytes, err := stub.GetState(SLA_ALL_DATA) // 리턴값 ([]byte, error)

		fmt.Println("= A04 == 01:" + string(contractALLIDsInBytes))

		if err != nil {
			return nil, errors.New("Failed to get state with" + string(contractALLIDsInBytes))
		}

		fmt.Println("= A04 == 02:" + string(contractALLIDsInBytes))

		contractALLIDsInString := string(contractALLIDsInBytes)

		fmt.Println("= A04 == 03:" + contractALLIDsInString)

		// 기존내역이 없을경우 "계약명"-"계약ID목록" 등록
		if contractALLIDsInString == "" {
			err = stub.PutState(SLA_ALL_DATA, []byte(contractID))
			fmt.Println("= A04 == 04:")

			if err != nil {
				return nil, err
			}

		} else {
			err = stub.PutState(SLA_ALL_DATA, []byte(contractALLIDsInString+FIELDSEP+contractID))

			fmt.Println("= A04 == 05:")

			if err != nil {
				return nil, err
			}
		}
	}

	return nil, nil
}

// 2.계약을 승인합니다.
func (t *SimpleChaincode) slaApproveContract(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	// 계약번호 > 계약 시퀀스 >  결제 정보

	var err error
	var data SlaContract

	slaTime := time.Now()

	SlaContractRegId := args[0]                             // SLA계약등록번호
	SlaContractApprovalUserId := args[1]                    // 결재사용자ID
	SlaContractApprovalComment := args[2]                   // 의견내용
	SlaContractProgression := args[3]                       // 진행단계
	SlaContractApprovalState := "승인"                        // 승인상태
	SlaContractApprovalDate := slaTime.Format("2006-01-02") // 현재일자

	SlaApprovalNo, _ := strconv.Atoi(SlaContractProgression)
	fmt.Println(SlaContractApprovalDate)

	//content := args[0]

	fmt.Printf("slaApproveContract Input Args:%s\n", args[0])
	fmt.Printf("slaApproveContract Input Args:%s\n", args[1])
	fmt.Printf("slaApproveContract Input Args:%s\n", args[2])
	fmt.Printf("slaApproveContract Input Args:%s\n", args[3])
	fmt.Printf("slaApproveContract Input Args:%s\n", SlaContractProgression)

	// 계약ID으로 목록 조회
	contractIDsInBytes, err := stub.GetState(SlaContractRegId)
	if err != nil {
		return nil, errors.New("Failed to get state with " + string(contractIDsInBytes))
	}
	fmt.Printf("=================[" + string(contractIDsInBytes) + "]=================\n")

	// JSON 데이터를 디코딩(Unmarshal)합니다.
	err = json.Unmarshal([]byte(contractIDsInBytes), &data)
	if err != nil {
		return nil, errors.New("Failed to slaApproveContract with " + string(contractIDsInBytes))
	}

	// 데이터 입력
	data.Approvals[SlaApprovalNo].ApprovalUserId = SlaContractApprovalUserId   // 결재사용자ID
	data.Approvals[SlaApprovalNo].ApprovalState = SlaContractApprovalState     // 결재상태
	data.Approvals[SlaApprovalNo].ApprovalDate = SlaContractApprovalDate       // 결재일자
	data.Approvals[SlaApprovalNo].ApprovalComment = SlaContractApprovalComment // 의견내용

	content, _ := json.Marshal(data)

	// JSON 데이터를 정렬하여 디코딩(Unmarshal)합니다.
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return nil, errors.New("Failed to registerContractByIdToJSON with " + string(contractIDsInBytes))
	}

	fmt.Printf("=================[" + string(jsonData) + "]=================\n")

	// 계약ID 등록합니다.
	err = stub.PutState(SlaContractRegId, []byte(content))

	if err != nil {
		return nil, errors.New("Failed to put state with" + string(content))
	}

	return nil, nil
}

// 3.계약을 반려합니다.
func (t *SimpleChaincode) slaRejectContract(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	// 계약번호 > 계약 시퀀스 >  결제 정보

	return nil, nil
}

func (t *SimpleChaincode) slaAbandonContract(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	// 업데이트를 abandon state  변경
	return nil, nil
}
func (t *SimpleChaincode) slaSubmitContract(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	// 업데이트를 submit state  변경
	return nil, nil
}
func (t *SimpleChaincode) slaCloseContract(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	// 업데이트를 submit state  변경
	return nil, nil
}
func (t *SimpleChaincode) slaGetEvaluationId(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return nil, nil
}

// --------------------

func (t *SimpleChaincode) slaCreateEvaluationRootFromContract(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return nil, nil
}
func (t *SimpleChaincode) slaCreateEvaluationsFromContract(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return nil, nil
}

func (t *SimpleChaincode) slaInitEvaluationValues(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return nil, nil
}
func (t *SimpleChaincode) slaUpdateEvaluationValues(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return nil, nil
}
func (t *SimpleChaincode) slaSubmitEvaluation(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return nil, nil
}
func (t *SimpleChaincode) slaApproveEvaluation(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return nil, nil
}
func (t *SimpleChaincode) slaRejectEvaluation(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return nil, nil
}
func (t *SimpleChaincode) slaSubmitPayment(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return nil, nil
}
func (t *SimpleChaincode) slaClosePayment(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return nil, nil
}
func (t *SimpleChaincode) slaCloseEvaluation(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return nil, nil
}
func (t *SimpleChaincode) slaCloseEvaluationRoot(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return nil, nil
}

// ===========================================================
//  SLAChaincodeStub 검색 함수
// ===========================================================

// SLA 데이터 전체를 조회합니다.  (abandon 제외)
func (t *SimpleChaincode) slaGetAllContracts(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	var dataInBytes string
	var err error

	fmt.Printf("slaGetAllContracts Input Args:%s\n", args[0])

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the value to slaGetContractsWithName")
	}

	dataInBytes = args[0]

	// 계약명으로 계약ID목록 조회
	contractIDsInBytes, err := stub.GetState(SLA_ALL_DATA)
	contractIDsInString := string(contractIDsInBytes)
	if err != nil {
		return nil, errors.New("Failed to get state with " + dataInBytes)
	}

	// 계약ID목록의 형태를 스트링에서 배열로 전환
	contractIDs := strings.Split(contractIDsInString, FIELDSEP)

	// 리턴값 초기화
	contractList := make([]string, len(contractIDs))

	// 계약 전체 ID목록 조회
	for i, _ := range contractIDs {
		contractInBytes, _ := stub.GetState(contractIDs[i])
		contractList[i] = string(contractInBytes)
	}

	contractListBytes := strings.Join(contractList, ENTRYSEP)

	fmt.Printf("slaGetContractsWithName Response:%s\n", contractListBytes)

	return []byte(contractListBytes), nil

}

// ID으로 조회합니다.
func (t *SimpleChaincode) slaGetContractWithId(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	var dataInBytes string
	var err error

	fmt.Printf("slaGetContractWithId Input Args:%s\n", args[0])

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the Value to slaGetContractWithId")
	}

	dataInBytes = args[0]
	Valuebytes, err := stub.GetState(args[0])

	if err != nil {
		return nil, errors.New("Failed to get state with" + dataInBytes)
	}

	fmt.Printf("searchbyid Response:%s\n", Valuebytes)

	return Valuebytes, nil
}

// 계약명으로 조회합니다.
func (t *SimpleChaincode) slaGetContractsWithName(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	var dataInBytes string
	var err error
	var data SlaContract

	fmt.Printf("slaGetContractsWithName Input Args:%s\n", args[0])

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the value to slaGetContractsWithName")
	}

	dataInBytes = args[0]
	contractName := args[0]

	// 계약명으로 계약ID목록 조회
	contractIDsInBytes, err := stub.GetState(contractName)
	contractIDsInString := string(contractIDsInBytes)
	if err != nil {
		return nil, errors.New("Failed to get state with " + dataInBytes)
	}

	// 계약ID목록의 형태를 스트링에서 배열로 전환합니다.
	contractIDs := strings.Split(contractIDsInString, FIELDSEP)

	// 리턴값 초기화
	contractList := make([]SlaContract, len(contractIDs))

	// 계약ID목록으로 계약내용을 추출하여 계약목록 작성
	for i, _ := range contractIDs {
		contractInBytes, _ := stub.GetState(contractIDs[i])

		err = json.Unmarshal(contractInBytes, &data)
		contractList[i] = data
	}

	contractListBytes, _ := json.Marshal(contractList)

	fmt.Printf("slaGetContractsWithName Response:%s\n", contractListBytes)

	return []byte(contractListBytes), nil

}

// 고객사명으로 조회합니다.
func (t *SimpleChaincode) slaGetContractsWithClient(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	var dataInBytes string
	var err error
	var data SlaContract

	fmt.Printf("slaGetContractsWithClient Input Args:%s\n", args[0])

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the value to slaGetContractsWithClient")
	}

	dataInBytes = args[0]
	contractClient := args[0]

	// 계약명으로 계약ID목록 조회
	contractIDsInBytes, err := stub.GetState(contractClient)
	contractIDsInString := string(contractIDsInBytes)
	if err != nil {
		return nil, errors.New("Failed to get state with " + dataInBytes)
	}

	// 계약ID목록의 형태를 스트링에서 배열로 전환합니다.
	contractIDs := strings.Split(contractIDsInString, FIELDSEP)

	// 리턴값 초기화
	contractList := make([]SlaContract, len(contractIDs))

	// 계약ID목록으로 계약내용을 추출하여 계약목록 작성
	for i, _ := range contractIDs {
		contractInBytes, _ := stub.GetState(contractIDs[i])

		err = json.Unmarshal(contractInBytes, &data)
		contractList[i] = data
	}

	contractListBytes, _ := json.Marshal(contractList)

	fmt.Printf("slaGetContractsWithName Response:%s\n", contractListBytes)

	return []byte(contractListBytes), nil

}

func (t *SimpleChaincode) slaGetAllEvaluations(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return nil, nil
}
func (t *SimpleChaincode) slaGetEvaluationWithId(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return nil, nil
}
func (t *SimpleChaincode) slaGetEvaluationsWithName(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return nil, nil
}
func (t *SimpleChaincode) slaGetEvaluationsWithClient(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return nil, nil
}

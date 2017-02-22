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
	Progression     string           `json:  "Progression"`     // 진행단계   const SLA_CONTRACT_PROGRESSION.... 사용할 것
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
	ApprovalState      string `json:  "ApprovalState"`      // 결재상태 const SLA_APPROVAL_STATE.... 사용할 것
	ApprovalDate       string `json:  "ApprovalDate"`       // 결재일자
	ApprovalComment    string `json:  "ApprovalComment"`    // 의견내용
	ApprovalAlarm      string `json:  "ApprovalAlarm"`      // 알람여부
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

// const ENTRYSEP = ","
const SLA_ALL_DATA = "SLA_ALL_DATA"

const CONTRACT_TEMP_ID_PREFIX = "SLA_CONT_TEMP_"
const CONTRACT_ID_PREFIX = "SLA_CONT_"
const EVALUATION_TEMP_ID_PREFIX = "SLA_EVAL_TEMP_"
const EVALUATION_ID_PREFIX = "SLA_EVAL_"

const SLA_CONTRACT_TEMP_ID_COUNT_KEY = "SLA_CONTRACT_TEMP_ID_COUNT"
const SLA_CONTRACT_ID_COUNT_KEY = "SLA_CONTRACT_ID_COUNT"
const SLA_EVALUATION_TEMP_ID_COUNT_KEY = "SLA_EVALUATION_TEMP_ID_COUNT"
const SLA_EVALUATION_ID_COUNT_KEY = "SLA_EVALUATION_ID_COUNT"
const CURRENT_YEAR_KEY = "CURRENT_YEAR"

//  계약 상태: SlaContract -> PROGRESSION
// -------------------------------------------------------------------------------------------------------
//  1. [초기 --> slaCreateTempContract] 임시저장상태: 							SLA_CONTRACT_PROGRESSION_TEMP
//  2. [초기 --> slaSubmitContract
//      slaCreateContract --> slaSubmitContract] Submit되어서 내부 검토 요청상태: 	SLA_CONTRACT_PROGRESSION_IN_PROGRESS_INTERNAL_REVIEW_REQUESTED
//  3. [slaSubmitContract --> slaApproveContract] 고객현업 검토 요청상태:			SLA_CONTRACT_PROGRESSION_IN_PROGRESS_CLIENT_REVIEW_REQUESTED
//  4. [slaApproveContract --> slaApproveContract]고객관리자 검토 요청상태: 		SLA_CONTRACT_PROGRESSION_IN_PROGRESS_CLIENT_MANAGER_REVIEW_REQUESTED
//  5. [slaApproveContract --> slaApproveContract] 계약 등록 완료상태: 			SLA_CONTRACT_PROGRESSION_CLOSED
//  6. [Any state --> slaAbandonContract] 계약 폐기상태:						SLA_CONTRACT_PROGRESSION_ABANDONED
// -------------------------------------------------------------------------------------------------------
//  * 반려될 경우 (slaRejectContract) 해당 상태에 머무른다
// -------------------------------------------------------------------------------------------------------
const SLA_CONTRACT_PROGRESSION_TEMP = "TEMP" // "TEMP": 임시저장
const SLA_CONTRACT_PROGRESSION_IN_PROGRESS_INTERNAL_REVIEW_REQUESTED = "IN_PROGRESS_INTERNAL_REVIEW_REQUESTED"
const SLA_CONTRACT_PROGRESSION_IN_PROGRESS_CLIENT_REVIEW_REQUESTED = "IN_PROGRESS_CLIENT_REVIEW_REQUESTED"
const SLA_CONTRACT_PROGRESSION_IN_PROGRESS_CLIENT_MANAGER_REVIEW_REQUESTED = "IN_PROGRESS_CLIENT_MANAGER_REVIEW_REQUESTED"
const SLA_CONTRACT_PROGRESSION_CLOSED = "CLOSED"
const SLA_CONTRACT_PROGRESSION_ABANDONED = "ABANDONED"

// 결재 상태: SlaContract -> Approvals --> ApprovalState
const SLA_APPROVAL_STATE_TEMP = "TEMP" // "TEMP": 임시저장
const SLA_APPROVAL_STATE_SUBMITTED = "SUBMITTED"
const SLA_APPROVAL_STATE_APPROVED = "APPROVED"
const SLA_APPROVAL_STATE_REJECTED = "REJECTED"

// ===========================================================
// Utility 함수
// ===========================================================

func padLeft(str string, padLength int) string {
	pad := "0"

	for {
		if len(str) >= padLength {
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
	err = stub.PutState(SLA_CONTRACT_TEMP_ID_COUNT_KEY, []byte(strconv.Itoa(1)))
	err = stub.PutState(SLA_EVALUATION_ID_COUNT_KEY, []byte(strconv.Itoa(1)))
	err = stub.PutState(SLA_EVALUATION_TEMP_ID_COUNT_KEY, []byte(strconv.Itoa(1)))
	err = stub.PutState(CURRENT_YEAR_KEY, []byte(strconv.Itoa(year))) // 현재 year
	// err = stub.PutState(SLA_ALL_DATA, []byte(""))
	return nil, err
}

// 기능 이벤트를 호출합니다.
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	switch function {

	case "slaGetContractTempId":
		return t.slaGetTempContractId(stub, args)

	case "slaGetContractId":
		return t.slaGetContractId(stub, args)

	// 최초 생성 + 임시 저장
	case "slaCreateTempContract": // 요청자가 계약 생성 (최초 생성하고 임시저장할 경우)
		return t.slaCreateTempContract(stub, args) //done

	// 최초 생성 + 바로 결제 요청
	case "slaSubmitContract": // 요청자 --> 결재자 (최초 생성하고 바로 결재 요청할 경우 / 수정  후 결재 요청 )
		return t.slaSubmitContract(stub, args) //done

	// 결재 요청 / 승인 / 거절
	case "slaUpdateContract": // 요청자가 계약 수정 (임시저장 후 / 승인거절 후)
		return t.slaUpdateContract(stub, args) //done

	case "slaApproveContract": // 결재자 --> 다음 결재자
		return t.slaApproveContract(stub, args) // done

	case "slaRejectContract": // 결재자 --> 요청자
		return t.slaRejectContract(stub, args) // done

	// 최종 승인
	case "slaCloseContract": // 최종 결재자 승인
		return t.slaCloseContract(stub, args)

	// 최종 폐
	case "slaAbandonContract": // 요청자가 계약 폐기 (임시저장 후 / 승인거절 후)
		return t.slaUpdateContract(stub, args)

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
	return nil, errors.New("Invalid Query function name. Expecting \"slaGetAllContracts\" \"slaGetContractWithId\" \"slaGetContractsWithName\" \"slaGetContractsWithClient\"")
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

func (t *SimpleChaincode) slaGetTempContractId(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var contractTempId string
	var err error

	// 1.저장된 해당 계약 카운트를 호출
	currentCountInBytes, err := stub.GetState(SLA_CONTRACT_TEMP_ID_COUNT_KEY)
	if err != nil {
		return nil, errors.New("Failed to SLA_CONTRACT_TEMP_ID_COUNT_KEY with " + SLA_CONTRACT_TEMP_ID_COUNT_KEY)
	}
	if currentCountInBytes == nil { // if not initialized
		err = stub.PutState(SLA_CONTRACT_TEMP_ID_COUNT_KEY, []byte(strconv.Itoa(1)))
	}
	currentCount, _ := strconv.Atoi(string(currentCountInBytes))

	// 2.카운트가 1000을 넘어가면 초기화
	if currentCount > 100000 { // new year starts
		err = stub.PutState(SLA_CONTRACT_TEMP_ID_COUNT_KEY, []byte(strconv.Itoa(1))) // 카운트는 1부터
	}

	// 3. 계약번호 채번을 생성합니다.
	currentYear, _, _ := time.Now().Date()
	contractTempId = CONTRACT_TEMP_ID_PREFIX + strconv.Itoa(currentYear) + "_" + padLeft(strconv.Itoa(currentCount), 5)

	// 4. 다음 계약번호 카운트를 저장
	nextCount := currentCount + 1
	stub.PutState(SLA_CONTRACT_TEMP_ID_COUNT_KEY, []byte(strconv.Itoa(nextCount)))

	// TODO 5. 혹시 같은 계약번호가 있는지 확인할 것.

	return []byte(contractTempId), nil
}

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
	contractId = CONTRACT_ID_PREFIX + strconv.Itoa(currentYear) + "_" + padLeft(strconv.Itoa(currentCount), 5)

	// 4. 다음 계약번호 카운트를 저장
	nextCount := currentCount + 1
	stub.PutState(SLA_CONTRACT_ID_COUNT_KEY, []byte(strconv.Itoa(nextCount)))

	// TODO 5. 혹시 같은 계약번호가 있는지 확인할 것.

	return []byte(contractId), nil
}

// 최초 생성 + 임시저장 계약을 생성합니다.
// KVS: 계약ID와 전체계약에 대한 KVS는 저장
//		계약명과 고객명에 따른 KVS는 저장하지 않음.
func (t *SimpleChaincode) slaCreateTempContract(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	var targetContract SlaContract

	// JSON 데이터를 디코딩(Unmarshal)합니다: string -> []byte -> golang struct
	err = json.Unmarshal([]byte(args[0]), &targetContract)
	if err != nil {
		return nil, errors.New("Failed to registerContractByIdToJSON with " + args[0])
	}

	// 계약상태를 SLA_CONTRACT_PROGRESSION_TEMP 로 변경
	targetContract.Progression = SLA_CONTRACT_PROGRESSION_TEMP

	// JSON 데이터를 정렬하여 디코딩(Unmarshal)합니다: golang struct --> string
	targetContractInJson, err := json.MarshalIndent(targetContract, "", "  ")
	if err != nil {
		return nil, errors.New("Failed to registerContractByIdToJSON with " + string(targetContractInJson))
	}

	// A01. 계약ID 등록합니다.
	err = stub.PutState(targetContract.RegId, targetContractInJson)
	if err != nil {
		return nil, errors.New("Failed to put state with" + args[0])
	}

	// 전체 조회 등록합니다.
	{
		var err error

		// 데이터를 전체 조회합니다.
		contractALLIdsInBytes, err := stub.GetState(SLA_ALL_DATA) // 리턴값 ([]byte, error)
		if err != nil {
			return nil, errors.New("Failed to get state with" + string(contractALLIdsInBytes))
		}
		contractALLIdsInString := string(contractALLIdsInBytes)

		// 최초이 없을경우 "계약명"-"계약ID목록" 등록 or 기존에 포함되어 있으면
		if contractALLIdsInString == "" {
			err = stub.PutState(SLA_ALL_DATA, []byte(targetContract.RegId))
			if err != nil {
				return nil, err
			}
		} else if !strings.Contains(contractALLIdsInString, targetContract.RegId) { // 계약번호가 존재하지 않을 때만 추가
			err = stub.PutState(SLA_ALL_DATA, []byte(contractALLIdsInString+FIELDSEP+targetContract.RegId))
			if err != nil {
				return nil, err
			}
		}
	}
	return nil, nil

}

// 최초 생성 + 내부결제를 요청합니다.
// KVS: 계약ID, 전체계약, 계약명, 고객명에 따른 KVS 저장
func (t *SimpleChaincode) slaSubmitContract(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	var targetContract SlaContract

	// todo 동일한 ID의 계약이 있을경우 튕겨 내야함.... submit 두번 할 수 없도록 해야함....

	// JSON 데이터를 디코딩(Unmarshal)합니다: string -> []byte -> golang struct
	err = json.Unmarshal([]byte(args[0]), &targetContract)
	if err != nil {
		return nil, errors.New("Failed to slaSubmitContract with " + args[0])
	}

	// 상태 변경: 계약상태 + Approvals[0]의 상태
	// 계약상태를 SLA_CONTRACT_PROGRESSION_TEMP 로 변경
	targetContract.Progression = SLA_CONTRACT_PROGRESSION_IN_PROGRESS_INTERNAL_REVIEW_REQUESTED
	// 첫번째 Approval의 state를 "SUBMITTED"로 변경
	targetContract.Approvals[0].ApprovalState = SLA_APPROVAL_STATE_SUBMITTED

	// JSON 데이터를 정렬하여 디코딩(Unmarshal)합니다: golang struct --> string
	targetContractInJson, err := json.MarshalIndent(targetContract, "", "  ")
	if err != nil {
		return nil, errors.New("Failed to registerContractByIdToJSON with " + string(targetContractInJson))
	}

	// A01. 계약ID 등록합니다.
	err = stub.PutState(targetContract.RegId, targetContractInJson)
	if err != nil {
		return nil, errors.New("Failed to put state with" + args[0])
	}

	// A02. 계약명 등록합니다.
	{
		var err error

		// 계약명으로 기존내역를 조회합니다.
		contractIDsInBytes, err := stub.GetState(targetContract.Name)
		if err != nil {
			return nil, errors.New("Failed to get state with" + targetContract.Name)
		}
		contractIDsInString := string(contractIDsInBytes)

		//기존내역이 없을경우 "계약명"-"계약ID목록" 등록
		if contractIDsInString == "" {
			err = stub.PutState(targetContract.Name, []byte(targetContract.RegId))
		} else {
			err = stub.PutState(targetContract.Name, []byte(contractIDsInString+FIELDSEP+targetContract.RegId))
		}
		if err != nil {
			return nil, err
		}
	}

	// A03. 고객사명 등록합니다.
	{
		var err error

		//계약명으로 기존내역조회
		contractIDsInBytes, _ := stub.GetState(targetContract.Client) // 리턴값 ([]byte, error)
		if err != nil {
			return nil, errors.New("Failed to get state with" + targetContract.Client)
		}
		contractIDsInString := string(contractIDsInBytes)

		//기존내역이 없을경우 "고객사명"-"계약ID목록" 등록
		if contractIDsInString == "" {
			err = stub.PutState(targetContract.Client, []byte(targetContract.RegId))
		} else {
			err = stub.PutState(targetContract.Client, []byte(contractIDsInString+FIELDSEP+targetContract.RegId))
		}
		if err != nil {
			return nil, err
		}
	}

	// A04. 전체 조회 등록합니다.  단 임시저장 때 계약번호를 딴 경우가 있으므로 이미 있는경우에는 건너 뛴다
	{
		var err error

		// 데이터를 전체 조회합니다.
		contractALLIdsInBytes, err := stub.GetState(SLA_ALL_DATA) // 리턴값 ([]byte, error)
		if err != nil {
			return nil, errors.New("Failed to get state with" + string(contractALLIdsInBytes))
		}
		contractALLIdsInString := string(contractALLIdsInBytes)

		// 최초이 없을경우 "계약명"-"계약ID목록" 등록 or
		if contractALLIdsInString == "" {
			err = stub.PutState(SLA_ALL_DATA, []byte(targetContract.RegId))
			if err != nil {
				return nil, err
			}
		} else if !strings.Contains(contractALLIdsInString, targetContract.RegId) { // 계약번호가 존재하지 않을 때만 추가
			err = stub.PutState(SLA_ALL_DATA, []byte(contractALLIdsInString+FIELDSEP+targetContract.RegId))
			if err != nil {
				return nil, err
			}
		}
	}
	return nil, nil
}

// 1.계약을 업데이트 합니다.
func (t *SimpleChaincode) slaUpdateContract(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	var err error
	var targetContract SlaContract

	fmt.Printf("slaUpdateContract Input Args:%s\n", args[0])

	// JSON 데이터를 디코딩(Unmarshal)합니다.
	err = json.Unmarshal([]byte(args[0]), &targetContract)
	if err != nil {
		return nil, errors.New("Failed to slaUpdateContract with " + args[0])
	}

	// 계약ID 통해 데이터를 업데이트 합니다.
	err = stub.PutState(targetContract.RegId, []byte(args[0]))
	if err != nil {
		return nil, errors.New("Failed to put state with" + args[0])
	}

	return nil, nil
}

// 2.계약을 승인합니다.
func (t *SimpleChaincode) slaApproveContract(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	//계약번호 > 계약 시퀀스 >  결제 정보
	var err error
	var targetContract SlaContract

	// 1. 해당 계약 찾기 및 Struct 생성
	slaContractRegId := args[0]
	targetContractInBytes, err := stub.GetState(slaContractRegId)
	if err != nil {
		return nil, errors.New("Failed to get state with " + string(targetContractInBytes))
	}

	// JSON 데이터를 디코딩(Unmarshal)합니다: string -> []byte -> golang struct
	err = json.Unmarshal([]byte(targetContractInBytes), &targetContract)
	if err != nil {
		return nil, errors.New("Failed to slaApproveContract with " + string(targetContractInBytes))
	}

	// 2. 현재의 진행단계(Progression)를 확인하여 현재 진행단계를 찾고 변경할 Approval을 찾음
	var newProgression string       // new progression for target contract
	var targetApproval *SlaApproval // approval to update for target contract

	switch targetContract.Progression {
	case SLA_CONTRACT_PROGRESSION_TEMP:
		return nil, errors.New("Reject cannot have the current progression of \"TEMP\" ")
	case SLA_CONTRACT_PROGRESSION_IN_PROGRESS_INTERNAL_REVIEW_REQUESTED:
		newProgression = SLA_CONTRACT_PROGRESSION_IN_PROGRESS_CLIENT_REVIEW_REQUESTED
		targetApproval = &(targetContract.Approvals[1]) // approval to change
	case SLA_CONTRACT_PROGRESSION_IN_PROGRESS_CLIENT_REVIEW_REQUESTED:
		newProgression = SLA_CONTRACT_PROGRESSION_IN_PROGRESS_CLIENT_MANAGER_REVIEW_REQUESTED
		targetApproval = &(targetContract.Approvals[2]) // approval to change
	case SLA_CONTRACT_PROGRESSION_IN_PROGRESS_CLIENT_MANAGER_REVIEW_REQUESTED:
		newProgression = SLA_CONTRACT_PROGRESSION_CLOSED
		targetApproval = &(targetContract.Approvals[3]) // approval to change
	case SLA_CONTRACT_PROGRESSION_CLOSED:
		return nil, errors.New("Approval cannot have the current progression of \"CLOSED\" ")
	case SLA_CONTRACT_PROGRESSION_ABANDONED:
		return nil, errors.New("Approval cannot have the current progression of \"ABANDONED\" ")
	default:
		return nil, errors.New("Approval cannot have the current progression of " + targetContract.Progression)
	}

	//3. 해당 결재 내용으로 변경
	targetContract.Progression = newProgression                   // 새 진행단계 (Progression)
	targetApproval.ApprovalUserId = args[1]                       // 결재사용자ID
	targetApproval.ApprovalState = SLA_APPROVAL_STATE_APPROVED    // 승인상태
	targetApproval.ApprovalDate = time.Now().Format("2006-01-02") // 현재일자
	targetApproval.ApprovalComment = args[2]                      // 의견내용

	targetContractInJson, _ := json.MarshalIndent(targetContract, "", "  ")

	// 4. 변경된 계약을 KVS에 저장합니다.
	err = stub.PutState(slaContractRegId, []byte(targetContractInJson))
	if err != nil {
		return nil, errors.New("Failed to put state with" + string(targetContractInJson))
	}

	return nil, nil
}

// 3.계약을 반려합니다.
func (t *SimpleChaincode) slaRejectContract(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	// 계약번호 > 계약 시퀀스 >  결제 정보

	// 계약번호 > 계약 시퀀스 >  결제 정보
	var err error
	var targetContract SlaContract

	// 1. 해당 계약 찾기 및 Struct 생성
	SlaContractRegId := args[0]
	targetContractInBytes, err := stub.GetState(SlaContractRegId)
	if err != nil {
		return nil, errors.New("Failed to get state with " + string(targetContractInBytes))
	}
	// 저장된 JSON형태의 데이터를 go의 struct 형태로 디코딩(Unmarshal)합니다.
	err = json.Unmarshal([]byte(targetContractInBytes), &targetContract)
	if err != nil {
		return nil, errors.New("Failed to slaApproveContract with " + string(targetContractInBytes))
	}

	// 2. 현재의 진행단계(Progression)를 확인하여 변경할 Approval을 찾음
	// * 현재의 진행단계(Pregression)은 현재 상태 그대로 유지
	var targetApproval *SlaApproval // approval to update for target contract

	switch targetContract.Progression {
	case SLA_CONTRACT_PROGRESSION_TEMP:
		return nil, errors.New("Reject cannot have the current progression of \"TEMP\" ")
	case SLA_CONTRACT_PROGRESSION_IN_PROGRESS_INTERNAL_REVIEW_REQUESTED:
		targetApproval = &(targetContract.Approvals[1]) // approval to change
	case SLA_CONTRACT_PROGRESSION_IN_PROGRESS_CLIENT_REVIEW_REQUESTED:
		targetApproval = &(targetContract.Approvals[2]) // approval to change
	case SLA_CONTRACT_PROGRESSION_IN_PROGRESS_CLIENT_MANAGER_REVIEW_REQUESTED:
		targetApproval = &(targetContract.Approvals[3]) // approval to change
	case SLA_CONTRACT_PROGRESSION_CLOSED:
		return nil, errors.New("Reject cannot have the current progression of \"CLOSED\" ")
	case SLA_CONTRACT_PROGRESSION_ABANDONED:
		return nil, errors.New("Reject cannot have the current progression of \"ABANDONED\" ")
	default:
		return nil, errors.New("Reject cannot have the current progression of " + targetContract.Progression)
	}

	// 3. 데이터 내용 변경
	targetApproval.ApprovalUserId = args[1]                       // 결재사용자ID
	targetApproval.ApprovalState = SLA_APPROVAL_STATE_REJECTED    // 승인상태
	targetApproval.ApprovalDate = time.Now().Format("2006-01-02") // 현재일자
	targetApproval.ApprovalComment = args[2]                      // 의견내용

	targetContractInJson, _ := json.MarshalIndent(targetContract, "", "  ")

	// 4. 변경된 계약을 KVS에 저장합니다.
	err = stub.PutState(SlaContractRegId, []byte(targetContractInJson))

	if err != nil {
		return nil, errors.New("Failed to put state with" + string(targetContractInJson))
	}

	return nil, nil
}

// 4.계약을 마무리(종료)합니다.
func (t *SimpleChaincode) slaCloseContract(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	//계약번호 > 계약 시퀀스 >  결제 정보
	var err error
	var targetContract SlaContract

	// 1. 해당 계약 찾기 및 Struct 생성
	slaContractRegId := args[0]
	targetContractInBytes, err := stub.GetState(slaContractRegId)
	if err != nil {
		return nil, errors.New("Failed to get state with " + string(targetContractInBytes))
	}

	// JSON 데이터를 디코딩(Unmarshal)합니다: string -> []byte -> golang struct
	err = json.Unmarshal([]byte(targetContractInBytes), &targetContract)
	if err != nil {
		return nil, errors.New("Failed to slaCloseContract with " + string(targetContractInBytes))
	}

	// 2. 현재의 진행단계(Progression)가 고객관리자에게 승인요청하는 최종 직전 단계여야함
	if targetContract.Progression != SLA_CONTRACT_PROGRESSION_IN_PROGRESS_CLIENT_MANAGER_REVIEW_REQUESTED {
		return nil, errors.New("slaCloseContract cannot have the current progression of " + targetContract.Progression)
	}

	// 3. 최종 진행단계(Progression)를 확인하여 현재 진행단계를 찾고 변경할 Approval을 찾음
	targetApproval := &(targetContract.Approvals[3]) // approval to update for target contract

	// 4. 해당 결재 내용으로 변경
	targetContract.Progression = SLA_CONTRACT_PROGRESSION_CLOSED  // 새 진행단계 (Progression)
	targetApproval.ApprovalUserId = args[1]                       // 결재사용자ID
	targetApproval.ApprovalState = SLA_APPROVAL_STATE_APPROVED    // 승인상태
	targetApproval.ApprovalDate = time.Now().Format("2006-01-02") // 현재일자
	targetApproval.ApprovalComment = args[2]                      // 의견내용

	targetContractInJson, _ := json.MarshalIndent(targetContract, "", "  ")

	// 5. 변경된 계약을 KVS에 저장합니다.
	err = stub.PutState(slaContractRegId, []byte(targetContractInJson))
	if err != nil {
		return nil, errors.New("Failed to put state with" + string(targetContractInJson))
	}
	return nil, nil
}

// 5.계약을 폐기합니다.
func (t *SimpleChaincode) slaAbandonContract(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	var targetContract SlaContract

	// 1. 해당 계약 찾기 및 Struct 생성
	slaContractRegId := args[0]
	targetContractInBytes, err := stub.GetState(slaContractRegId)
	if err != nil {
		return nil, errors.New("Failed to get state with " + string(targetContractInBytes))
	}

	// 2. JSON 데이터를 디코딩(Unmarshal)합니다: string -> []byte -> golang struct
	err = json.Unmarshal([]byte(targetContractInBytes), &targetContract)
	if err != nil {
		return nil, errors.New("Failed to slaCloseContract with " + string(targetContractInBytes))
	}

	// 3. 해당 결재 내용으로 변경
	targetContract.Progression = SLA_CONTRACT_PROGRESSION_ABANDONED // 새 진행단계 (Progression)

	// 4. Json 형태로 저장: golang struct -> string
	targetContractInJson, _ := json.MarshalIndent(targetContract, "", "  ")

	// 5. 변경된 계약을 KVS에 저장합니다.
	err = stub.PutState(slaContractRegId, []byte(targetContractInJson))
	if err != nil {
		return nil, errors.New("Failed to put state with" + string(targetContractInJson))
	}
	return nil, nil
}

// --------------------

func (t *SimpleChaincode) slaGetEvaluationId(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return nil, nil
}

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

// SLA 데이터 전체를 조회합니다.  (abandon 포함)
func (t *SimpleChaincode) slaGetAllContracts(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error

	// 전체 계약ID목록 조회
	contractIDsInBytes, err := stub.GetState(SLA_ALL_DATA)
	if err != nil {
		return nil, errors.New("Failed to get state with " + SLA_ALL_DATA)
	}
	contractIDsInString := string(contractIDsInBytes)

	// 계약ID목록의 형태를 스트링에서 배열로 전환
	contractIDs := strings.Split(contractIDsInString, FIELDSEP)

	// 리턴값 초기화
	contractList := make([]string, len(contractIDs))

	// 계약 전체 ID목록 조회
	for i, _ := range contractIDs {
		contractInBytes, _ := stub.GetState(contractIDs[i])
		if err != nil {
			return nil, errors.New("Failed to get state with " + contractIDs[i])
		}
		contractList[i] = string(contractInBytes)
	}
	// 계약 전체신규_20170221 String
	ContractsInJson, err := json.MarshalIndent(contractList, "", "  ")
	if err != nil {
		return nil, errors.New("Failed to json.MarshalIndent with " + strings.Join(contractList, ","))
	}

	return []byte(ContractsInJson), nil

}

// ID으로 조회합니다.
func (t *SimpleChaincode) slaGetContractWithId(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the Value to slaGetContractWithId")
	}

	contractInBytes, err := stub.GetState(args[0])
	if err != nil {
		return nil, errors.New("Failed to get state with" + args[0])
	}

	return contractInBytes, nil
}

// 계약명으로 조회합니다.
func (t *SimpleChaincode) slaGetContractsWithName(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the value to slaGetContractsWithName")
	}
	contractName := args[0]

	// 계약명으로 계약ID목록 조회
	contractIDsInBytes, err := stub.GetState(contractName)
	if err != nil {
		return nil, errors.New("Failed to get state with " + contractName)
	}
	contractIDsInString := string(contractIDsInBytes)

	// 계약ID목록의 형태를 스트링에서 배열로 전환합니다.
	contractIDs := strings.Split(contractIDsInString, FIELDSEP)

	// 리턴값 초기화
	contractList := make([]string, len(contractIDs))

	// 계약ID목록으로 계약내용을 추출하여 계약목록 작성
	for i, _ := range contractIDs {
		contractInBytes, _ := stub.GetState(contractIDs[i])
		if err != nil {
			return nil, errors.New("Failed to get state with " + contractIDs[i])
		}
		contractList[i] = string(contractInBytes)
	}
	// 계약 전체신규_20170221 String
	ContractsInJson, err := json.MarshalIndent(contractList, "", "  ")
	if err != nil {
		return nil, errors.New("Failed to json.MarshalIndent with " + strings.Join(contractList, ","))
	}

	return []byte(ContractsInJson), nil
}

// 고객사명으로 조회합니다.
func (t *SimpleChaincode) slaGetContractsWithClient(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the value to slaGetContractsWithClient")
	}
	clientName := args[0]

	// 계약명으로 계약ID목록 조회
	contractIDsInBytes, err := stub.GetState(clientName)
	if err != nil {
		return nil, errors.New("Failed to get state with " + clientName)
	}
	contractIDsInString := string(contractIDsInBytes)

	// 계약ID목록의 형태를 스트링에서 배열로 전환합니다.
	contractIDs := strings.Split(contractIDsInString, FIELDSEP)

	// 리턴값 초기화
	contractList := make([]string, len(contractIDs))

	// 계약ID목록으로 계약내용을 추출하여 계약목록 작성
	for i, _ := range contractIDs {
		contractInBytes, _ := stub.GetState(contractIDs[i])
		if err != nil {
			return nil, errors.New("Failed to get state with " + contractIDs[i])
		}
		contractList[i] = string(contractInBytes)
	}
	// 계약 전체신규_20170221 String
	ContractsInJson, err := json.MarshalIndent(contractList, "", "  ")
	if err != nil {
		return nil, errors.New("Failed to json.MarshalIndent with " + strings.Join(contractList, ","))
	}

	return []byte(ContractsInJson), nil
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

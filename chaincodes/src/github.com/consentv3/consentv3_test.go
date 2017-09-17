package main

import (
	"testing"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"time"
	"strconv"
	"encoding/json"
	"strings"
	"github.com/op/go-logging"
	"os"
)
const (
	APPID1 = "APP4TESTS1"
	APPID2 = "APP4TESTS2"
	OWNERID1 = "owner1"
	OWNERID2 = "owner2"
	CONSUMERID1 = "consumer1"
	CONSUMERID2 = "consumer2"
	DATATYPE1 = "type1"
	DATATYPE2 = "type2"
	DATAACCESS1 = "access1"
	DATAACCESS2 = "access2"
)

func  InitLogger()  {
	logLevel     := logging.INFO
	//consLogLevel := logging.DEBUG
	f            := os.Stderr
	//loggerModule := ""
	var format    = logging.MustStringFormatter(
		`%{color}%{time:15:04:05.000} [%{module}:%{shortfile}] %{level:.4s} : %{color:reset} %{message}`,
	)
	backend := logging.NewLogBackend(f, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	backendLeveled := logging.AddModuleLevel(backendFormatter)
	backendLeveled.SetLevel(logging.Level(logLevel), "mock")

	//backendConsent := logging.NewLogBackend(f, "", 0)
	//backendConsFormatter := logging.NewBackendFormatter(backendConsent, format)
	//backendConsLeveled := logging.AddModuleLevel(backendConsFormatter)
	//backendConsLeveled.SetLevel(consLogLevel,LOGGERMODULE)
	//logging.SetBackend(backendLeveled, backendConsLeveled)
	logging.SetBackend(backendLeveled)
}


// =====================================================================================================================
// Get version of smart contract (nominal case)
// =====================================================================================================================
func TestConsentV2_GetVersionNominal(t *testing.T) {
	scc := new(ConsentCC)
	stub := shim.NewMockStub("consentv2", scc)
	res := stub.MockInvoke("1", [][]byte{[]byte("getversion")})
	if res.Status != shim.OK {
		t.Log("bad status received, expected: 200 received:"+strconv.FormatInt(int64(res.Status), 10))
		t.Log("response: "+ string(res.Message))
		t.FailNow()
	}
	if res.Payload == nil {
		t.Log("getversion failed to get chaincode version")
		t.FailNow()
	}
	version := string(res.Payload)
	if  version != VERSION {
		t.Log("getversion", version, "was not", VERSION, "as expected")
		t.FailNow()
	}
}

// =====================================================================================================================
// Create a consent (nominal case)
// =====================================================================================================================
func TestConsentV2_CreateConsentNominal(t *testing.T){
	scc := new(ConsentCC)
	stub := shim.NewMockStub("consentv2", scc)
	res := stub.MockInvoke("1", [][]byte{[]byte("postconsent"), []byte(APPID1), []byte(OWNERID1), []byte(CONSUMERID1), []byte(DATATYPE1), []byte(DATAACCESS1), []byte(getStringDateNow(0)), []byte(getStringDateNow(7))})
	if res.Status != shim.OK {
		t.Log("bad status received, expected: 200 received:"+strconv.FormatInt(int64(res.Status), 10))
		t.Log("response: "+ string(res.Message))
		t.FailNow()
	}
	if res.Payload == nil {
		t.Log("postconsent failed to create a consent and get the consent id")
		t.FailNow()
	}
}

// =====================================================================================================================
// Get a consent (nominal case)
// =====================================================================================================================
func TestConsentV2_GetConsentNominal(t *testing.T) {
	scc := new(ConsentCC)
	stub := shim.NewMockStub("consentv2", scc)
	res := stub.MockInvoke("1", [][]byte{[]byte("postconsent"), []byte(APPID1), []byte(OWNERID1), []byte(CONSUMERID1), []byte(DATATYPE1), []byte(DATAACCESS1), []byte(getStringDateNow(0)), []byte(getStringDateNow(7))})
	consentID := string(res.Payload)
	res = stub.MockInvoke("2", [][]byte{[]byte("getconsent"), []byte(APPID1), []byte(consentID)})
	if res.Status != shim.OK {
		t.Log("bad status received, expected: 200 received:"+strconv.FormatInt(int64(res.Status), 10))
		t.Log("response: "+ string(res.Message))
		t.FailNow()
	}
	var consent consent
	err := json.Unmarshal(res.Payload, &consent)
	if err != nil {
		t.Log("getconsent", string(res.Payload))
		t.FailNow()
	}
	if consent.AppID != APPID1 {
		t.Log("getconsent bad appID")
		t.FailNow()
	}
	if consent.ConsentID != consentID {
		t.Log("getconsent bad consentID")
		t.FailNow()
	}
}

// =====================================================================================================================
// Inactivate a consent (nominal case)
// =====================================================================================================================
func TestConsentV2_RemoveConsentNominal(t *testing.T) {
	scc := new(ConsentCC)
	stub := shim.NewMockStub("consentv2", scc)
	res := stub.MockInvoke("1", [][]byte{[]byte("postconsent"), []byte(APPID1), []byte(OWNERID1), []byte(CONSUMERID1), []byte(DATATYPE1), []byte(DATAACCESS1), []byte(getStringDateNow(0)), []byte(getStringDateNow(7))})
	consentID := string(res.Payload)
	res = stub.MockInvoke("2", [][]byte{[]byte("deleteconsent"), []byte(APPID1), []byte(consentID)})
	if res.Status != shim.OK {
		t.Log("bad status received, expected: 200 received:"+strconv.FormatInt(int64(res.Status), 10))
		t.Log("response: "+ string(res.Message))
		t.FailNow()
	}
	res = stub.MockInvoke("3", [][]byte{[]byte("getconsent"), []byte(APPID1), []byte(consentID)})
	if res.Status != shim.ERROR{
		t.Log("bad status received, expected: 500 received:"+strconv.FormatInt(int64(res.Status), 10))
		t.FailNow()
	}
	if !strings.Contains(res.Message, errorConsentNotExist){
		t.Log("Bad return message, expected:"+errorConsentNotExist+" reveived:"+string(res.Message))
		t.FailNow()
	}
}

// =====================================================================================================================
// Get list of all consents for an application (nominal case)
// =====================================================================================================================
func TestConsentV2_GetAllConsents4AppIDNominal(t *testing.T) {
	scc := new(ConsentCC)
	stub := shim.NewMockStub("consentv2", scc)
	stub.MockInvoke("1", [][]byte{[]byte("postconsent"), []byte(APPID1), []byte(OWNERID1), []byte(CONSUMERID1), []byte(DATATYPE1), []byte(DATAACCESS1), []byte(getStringDateNow(0)), []byte(getStringDateNow(7))})
	stub.MockInvoke("2", [][]byte{[]byte("postconsent"), []byte(APPID1), []byte(OWNERID1), []byte(CONSUMERID2), []byte(DATATYPE1), []byte(DATAACCESS1), []byte(getStringDateNow(0)), []byte(getStringDateNow(7))})
	stub.MockInvoke("3", [][]byte{[]byte("postconsent"), []byte(APPID1), []byte(OWNERID2), []byte(CONSUMERID1), []byte(DATATYPE1), []byte(DATAACCESS1), []byte(getStringDateNow(0)), []byte(getStringDateNow(7))})
	stub.MockInvoke("4", [][]byte{[]byte("postconsent"), []byte(APPID1), []byte(OWNERID2), []byte(CONSUMERID2), []byte(DATATYPE1), []byte(DATAACCESS1), []byte(getStringDateNow(0)), []byte(getStringDateNow(7))})
	stub.MockInvoke("5", [][]byte{[]byte("postconsent"), []byte(APPID2), []byte(OWNERID2), []byte(CONSUMERID2), []byte(DATATYPE1), []byte(DATAACCESS1), []byte(getStringDateNow(0)), []byte(getStringDateNow(7))})
	res := stub.MockInvoke("6", [][]byte{[]byte("getconsents"), []byte(APPID1)})
	if res.Status != shim.OK {
		t.Log("bad status received, expected: 200 received:"+strconv.FormatInt(int64(res.Status), 10))
		t.Log("response: "+ string(res.Message))
		t.FailNow()
	}
	consents := make([]consent, 0)
	err := json.Unmarshal(res.Payload, &consents)
	if err != nil {
		t.Log("getconsents", string(res.Payload))
		t.FailNow()
	}
	if len(consents) != 4{
		t.Error("4 expected, but ",strconv.Itoa(len(consents)), "reveived")
		t.FailNow()
	}
}

// =====================================================================================================================
// Delete all consents for an application (nominal case)
// =====================================================================================================================
func TestConsentV2_DeleteConsents4AppIDNominal(t *testing.T) {
	scc := new(ConsentCC)
	stub := shim.NewMockStub("consentv2", scc)
	stub.MockInvoke("1", [][]byte{[]byte("postconsent"), []byte(APPID1), []byte(OWNERID1), []byte(CONSUMERID1), []byte(DATATYPE1), []byte(DATAACCESS1), []byte(getStringDateNow(0)), []byte(getStringDateNow(7))})
	stub.MockInvoke("2", [][]byte{[]byte("postconsent"), []byte(APPID1), []byte(OWNERID1), []byte(CONSUMERID2), []byte(DATATYPE1), []byte(DATAACCESS1), []byte(getStringDateNow(0)), []byte(getStringDateNow(7))})
	stub.MockInvoke("3", [][]byte{[]byte("postconsent"), []byte(APPID1), []byte(OWNERID2), []byte(CONSUMERID1), []byte(DATATYPE1), []byte(DATAACCESS1), []byte(getStringDateNow(0)), []byte(getStringDateNow(7))})
	stub.MockInvoke("4", [][]byte{[]byte("postconsent"), []byte(APPID1), []byte(OWNERID2), []byte(CONSUMERID2), []byte(DATATYPE1), []byte(DATAACCESS1), []byte(getStringDateNow(0)), []byte(getStringDateNow(7))})
	stub.MockInvoke("5", [][]byte{[]byte("postconsent"), []byte(APPID2), []byte(OWNERID2), []byte(CONSUMERID2), []byte(DATATYPE1), []byte(DATAACCESS1), []byte(getStringDateNow(0)), []byte(getStringDateNow(7))})
	res := stub.MockInvoke("6", [][]byte{[]byte("deleteconsents"), []byte(APPID1)})
	if res.Status != shim.OK {
		t.Log("bad status received, expected: 200 received:"+strconv.FormatInt(int64(res.Status), 10))
		t.Log("response: "+ string(res.Message))
		t.FailNow()
	}
	res = stub.MockInvoke("7", [][]byte{[]byte("getconsents"), []byte(APPID1)})
	if res.Status != shim.OK{
		t.Log("bad status received, expected: 200 received:"+strconv.FormatInt(int64(res.Status), 10))
		t.Log("response: "+ string(res.Message))
		t.FailNow()
	}
	consents := make([]consent, 0)
	err := json.Unmarshal(res.Payload, &consents)
	if err != nil {
		t.Log("getconsents", string(res.Payload))
		t.FailNow()
	}
	if len(consents) != 0{
		t.Error("empty list expected, but ",strconv.Itoa(len(consents)), "reveived")
		t.FailNow()
	}
	res = stub.MockInvoke("8", [][]byte{[]byte("getconsents"), []byte(APPID2)})
	if res.Status != shim.OK{
		t.Log("bad status received, expected: 200 received:"+strconv.FormatInt(int64(res.Status), 10))
		t.Log("response: "+ string(res.Message))
		t.FailNow()
	}
	consents = make([]consent, 0)
	err = json.Unmarshal(res.Payload, &consents)
	if err != nil {
		t.Log("getconsents", string(res.Payload))
		t.FailNow()
	}
	if len(consents) != 1{
		t.Error("one consent expected, but ",strconv.Itoa(len(consents)), "reveived")
		t.FailNow()
	}
}

// =====================================================================================================================
// Get list of all consents for an owner (nominal case)
// =====================================================================================================================
func TestConsentV2_GetOwnerConsentsNominal(t *testing.T) {
	scc := new(ConsentCC)
	stub := shim.NewMockStub("consentv2", scc)
	stub.MockInvoke("1", [][]byte{[]byte("postconsent"), []byte(APPID1), []byte(OWNERID1), []byte(CONSUMERID1), []byte(DATATYPE1), []byte(DATAACCESS1), []byte(getStringDateNow(0)), []byte(getStringDateNow(7))})
	stub.MockInvoke("2", [][]byte{[]byte("postconsent"), []byte(APPID1), []byte(OWNERID1), []byte(CONSUMERID2), []byte(DATATYPE1), []byte(DATAACCESS1), []byte(getStringDateNow(0)), []byte(getStringDateNow(7))})
	stub.MockInvoke("3", [][]byte{[]byte("postconsent"), []byte(APPID2), []byte(OWNERID1), []byte(CONSUMERID2), []byte(DATATYPE1), []byte(DATAACCESS1), []byte(getStringDateNow(0)), []byte(getStringDateNow(7))})
	stub.MockInvoke("4", [][]byte{[]byte("postconsent"), []byte(APPID1), []byte(OWNERID2), []byte(CONSUMERID2), []byte(DATATYPE1), []byte(DATAACCESS1), []byte(getStringDateNow(0)), []byte(getStringDateNow(7))})
	res := stub.MockInvoke("5", [][]byte{[]byte("getownerconsents"), []byte(APPID1), []byte(OWNERID1)})
	if res.Status != shim.OK {
		t.Log("bad status received, expected: 200 received:"+strconv.FormatInt(int64(res.Status), 10))
		t.Log("response: "+ string(res.Message))
		t.FailNow()
	}
	consents := make([]consent, 0)
	err := json.Unmarshal(res.Payload, &consents)
	if err != nil {
		t.Log("getconsent", string(res.Payload))
		t.FailNow()
	}
	if len(consents) != 2{
		t.Error("2 expected, but ",strconv.Itoa(len(consents)), "reveived")
		t.FailNow()
	}
}

// =====================================================================================================================
// Get list of all consents for a consumer (nominal case)
// =====================================================================================================================
func TestConsentV2_GetConsumerConsentsNominal(t *testing.T) {
	scc := new(ConsentCC)
	stub := shim.NewMockStub("consentv2", scc)
	stub.MockInvoke("1", [][]byte{[]byte("postconsent"), []byte(APPID1), []byte(OWNERID1), []byte(CONSUMERID1), []byte(DATATYPE1), []byte(DATAACCESS1), []byte(getStringDateNow(0)), []byte(getStringDateNow(7))})
	stub.MockInvoke("2", [][]byte{[]byte("postconsent"), []byte(APPID1), []byte(OWNERID2), []byte(CONSUMERID1), []byte(DATATYPE1), []byte(DATAACCESS1), []byte(getStringDateNow(0)), []byte(getStringDateNow(7))})
	stub.MockInvoke("3", [][]byte{[]byte("postconsent"), []byte(APPID2), []byte(OWNERID1), []byte(CONSUMERID1), []byte(DATATYPE1), []byte(DATAACCESS1), []byte(getStringDateNow(0)), []byte(getStringDateNow(7))})
	stub.MockInvoke("4", [][]byte{[]byte("postconsent"), []byte(APPID1), []byte(OWNERID1), []byte(CONSUMERID2), []byte(DATATYPE1), []byte(DATAACCESS1), []byte(getStringDateNow(0)), []byte(getStringDateNow(7))})

	res := stub.MockInvoke("5", [][]byte{[]byte("getconsumerconsents"), []byte(APPID1), []byte(CONSUMERID1)})
	if res.Status != shim.OK {
		t.Log("bad status received, expected: 200 received:"+strconv.FormatInt(int64(res.Status), 10))
		t.Log("response: "+ string(res.Message))
		t.FailNow()
	}
	consents := make([]consent, 0)
	err := json.Unmarshal(res.Payload, &consents)
	if err != nil {
		t.Log("getconsent", string(res.Payload))
		t.FailNow()
	}
	if len(consents) != 2{
		t.Error("2 expected, but ",strconv.Itoa(len(consents)), "reveived")
		t.FailNow()
	}
}

// =====================================================================================================================
// Get list of all consents for a consumer and a owner (nominal case)
// =====================================================================================================================
func TestConsentV2_GetConsumerOwnerConsentsNominal(t *testing.T) {
	scc := new(ConsentCC)
	stub := shim.NewMockStub("consentv2", scc)
	stub.MockInvoke("1", [][]byte{[]byte("postconsent"), []byte(APPID1), []byte(OWNERID1), []byte(CONSUMERID1), []byte(DATATYPE1), []byte(DATAACCESS1), []byte(getStringDateNow(0)), []byte(getStringDateNow(7))})
	stub.MockInvoke("2", [][]byte{[]byte("postconsent"), []byte(APPID1), []byte(OWNERID2), []byte(CONSUMERID1), []byte(DATATYPE1), []byte(DATAACCESS1), []byte(getStringDateNow(0)), []byte(getStringDateNow(7))})
	stub.MockInvoke("3", [][]byte{[]byte("postconsent"), []byte(APPID1), []byte(OWNERID1), []byte(CONSUMERID1), []byte(DATATYPE2), []byte(DATAACCESS1), []byte(getStringDateNow(0)), []byte(getStringDateNow(7))})
	stub.MockInvoke("4", [][]byte{[]byte("postconsent"), []byte(APPID1), []byte(OWNERID1), []byte(CONSUMERID2), []byte(DATATYPE1), []byte(DATAACCESS1), []byte(getStringDateNow(0)), []byte(getStringDateNow(7))})

	res := stub.MockInvoke("5", [][]byte{[]byte("getconsumerownerconsents"), []byte(APPID1), []byte(CONSUMERID1), []byte(OWNERID1)})
	if res.Status != shim.OK {
		t.Log("bad status received, expected: 200 received:"+strconv.FormatInt(int64(res.Status), 10))
		t.Log("response: "+ string(res.Message))
		t.FailNow()
	}
	consents := make([]consent, 0)
	err := json.Unmarshal(res.Payload, &consents)
	if err != nil {
		t.Log("getconsent", string(res.Payload))
		t.FailNow()
	}
	if len(consents) != 2{
		t.Error("2 expected, but ",strconv.Itoa(len(consents)), "reveived")
		t.FailNow()
	}
}
// =====================================================================================================================
// Is consent exist (nominal case)
// =====================================================================================================================
func TestConsentV2_IsConsentNominal(t *testing.T) {
	scc := new(ConsentCC)
	stub := shim.NewMockStub("consentv2", scc)
	stub.MockInvoke("1", [][]byte{[]byte("postconsent"), []byte(APPID1), []byte(OWNERID1), []byte(CONSUMERID1), []byte(DATATYPE1), []byte(DATAACCESS1), []byte(getStringDateNow(0)), []byte(getStringDateNow(7))})

	res := stub.MockInvoke("1", [][]byte{[]byte("isconsent"), []byte(APPID1), []byte(OWNERID1), []byte(CONSUMERID1), []byte(DATATYPE1), []byte(DATAACCESS1)})
	if res.Status != shim.OK {
		t.Log("bad status received, expected: 200 received:"+strconv.FormatInt(int64(res.Status), 10))
		t.Log("response: "+ string(res.Message))
		t.FailNow()
	}
	response := string(res.Payload)
	if response != AUTHORIZED{
		t.Log(AUTHORIZED, "expected, but ",response, "reveived")
		t.FailNow()
	}
}

// =====================================================================================================================
// Get bad function name --> error
// =====================================================================================================================
func TestConsentV2_GetBadFunction(t *testing.T) {
	scc := new(ConsentCC)
	stub := shim.NewMockStub("consentv2", scc)
	res := stub.MockInvoke("1", [][]byte{[]byte("badFunction")})
	if res.Status != shim.ERROR{
		t.Log("bad status received, expected: 500 received:"+strconv.FormatInt(int64(res.Status), 10))
		t.FailNow()
	}
	if !strings.Contains(res.Message, errorBadFunctionName){
		t.Log("Bad return message, expected:"+errorBadFunctionName+" reveived:"+string(res.Message))
		t.FailNow()
	}
}

// =====================================================================================================================
// Get version of smart contract with a param --> error
// =====================================================================================================================
func TestConsentV2_GetVersionWithParam(t *testing.T) {
	scc := new(ConsentCC)
	stub := shim.NewMockStub("consentv2", scc)
	res := stub.MockInvoke("1", [][]byte{[]byte("getversion"), []byte(APPID1), []byte("bad param")})
	if res.Status != shim.ERROR{
		t.Log("bad status received, expected: 500 received:"+strconv.FormatInt(int64(res.Status), 10))
		t.FailNow()
	}
	if !strings.Contains(res.Message, errorArgs){
		t.Log("Bad return message, expected:"+errorArgs+" reveived:"+string(res.Message))
		t.FailNow()
	}
}

// =====================================================================================================================
// Create a consent with a missing param  --> error
// =====================================================================================================================
func TestConsentV2_CreateConsentWithAMissingParam(t *testing.T){
	scc := new(ConsentCC)
	stub := shim.NewMockStub("consentv2", scc)
	res := stub.MockInvoke("1", [][]byte{[]byte("postconsent"), []byte(APPID1), []byte(OWNERID1), []byte(CONSUMERID1), []byte(DATATYPE1), []byte(getStringDateNow(0)), []byte(getStringDateNow(7))})
	if res.Status != shim.ERROR || !strings.Contains(res.Message, "Incorrect number of arguments"){
		t.Log("postconsent", string(res.Message))
		t.FailNow()
	}
	if res.Status != shim.ERROR{
		t.Log("bad status received, expected: 500 received:"+strconv.FormatInt(int64(res.Status), 10))
		t.FailNow()
	}
	if !strings.Contains(res.Message, errorArgs){
		t.Log("Bad return message, expected:"+errorArgs+" reveived:"+string(res.Message))
		t.FailNow()
	}
}

// =====================================================================================================================
// Create a consent with a bad date format for dt_begin  --> error
// =====================================================================================================================
func TestConsentV2_CreateConsentWithBadStartingDateFormat(t *testing.T){
	scc := new(ConsentCC)
	stub := shim.NewMockStub("consentv2", scc)
	res := stub.MockInvoke("1", [][]byte{[]byte("postconsent"), []byte(APPID1), []byte(OWNERID1), []byte(CONSUMERID1), []byte(DATATYPE1),  []byte(DATAACCESS1), []byte("2017"), []byte(getStringDateNow(7))})
	if res.Status != shim.ERROR{
		t.Log("bad status received, expected: 500 received:"+strconv.FormatInt(int64(res.Status), 10))
		t.FailNow()
	}
	if !strings.Contains(res.Message, errorDateBegin){
		t.Log("Bad return message, expected:"+errorDateBegin+" reveived:"+string(res.Message))
		t.FailNow()
	}
}

// =====================================================================================================================
// Create a consent with a bad date format for dt_end  --> error
// =====================================================================================================================
func TestConsentV2_CreateConsentWithBadEndingDateFormat(t *testing.T){
	scc := new(ConsentCC)
	stub := shim.NewMockStub("consentv2", scc)
	res := stub.MockInvoke("1", [][]byte{[]byte("postconsent"), []byte(APPID1), []byte(OWNERID1), []byte(CONSUMERID1), []byte(DATATYPE1),  []byte(DATAACCESS1), []byte(getStringDateNow(7)), []byte("2017")})
	if res.Status != shim.ERROR{
		t.Log("bad status received, expected: 500 received:"+strconv.FormatInt(int64(res.Status), 10))
		t.FailNow()
	}
	if !strings.Contains(res.Message, errorDateEnd){
		t.Log("Bad return message, expected:"+errorDateEnd+" reveived:"+string(res.Message))
		t.FailNow()
	}
}

// =====================================================================================================================
// Create a consent with a bad period, ending date anterior to begin date  --> error
// =====================================================================================================================
func TestConsentV2_CreateConsentWithBadPeriod(t *testing.T){
	scc := new(ConsentCC)
	stub := shim.NewMockStub("consentv2", scc)
	res := stub.MockInvoke("1", [][]byte{[]byte("postconsent"), []byte(APPID1), []byte(OWNERID1), []byte(CONSUMERID1), []byte(DATATYPE1),  []byte(DATAACCESS1), []byte(getStringDateNow(7)),[]byte(getStringDateNow(0))})
	if res.Status != shim.ERROR{
		t.Log("bad status received, expected: 500 received:"+strconv.FormatInt(int64(res.Status), 10))
		t.FailNow()
	}
	if !strings.Contains(res.Message, errorPeriod){
		t.Log("Bad return message, expected:"+errorPeriod+" reveived:"+string(res.Message))
		t.FailNow()
	}
}

// =====================================================================================================================
// Inactivate a consent with a missing parameter --> error
// =====================================================================================================================
func TestConsentV2_InactivateConsentWithMissingParameter(t *testing.T) {
	scc := new(ConsentCC)
	stub := shim.NewMockStub("consentv2", scc)
	res := stub.MockInvoke("1", [][]byte{[]byte("postconsent"), []byte(APPID1), []byte(OWNERID1), []byte(CONSUMERID1), []byte(DATATYPE1), []byte(DATAACCESS1), []byte(getStringDateNow(0)), []byte(getStringDateNow(7))})
	consentID := string(res.Payload)
	res = stub.MockInvoke("2", [][]byte{[]byte("deleteconsent"), []byte(consentID)})
	if res.Status != shim.ERROR{
		t.Log("bad status received, expected: 500 received:"+strconv.FormatInt(int64(res.Status), 10))
		t.FailNow()
	}
	if !strings.Contains(res.Message, errorArgs){
		t.Log("Bad return message, expected:"+errorArgs+" reveived:"+string(res.Message))
		t.FailNow()
	}
}


// =====================================================================================================================
// Inactivate a consent with a bad consentID --> error
// =====================================================================================================================
func TestConsentV2_InactivateConsentWithBadConsentID(t *testing.T) {
	scc := new(ConsentCC)
	stub := shim.NewMockStub("consentv2", scc)
	res := stub.MockInvoke("1", [][]byte{[]byte("postconsent"), []byte(APPID1), []byte(OWNERID1), []byte(CONSUMERID1), []byte(DATATYPE1), []byte(DATAACCESS1), []byte(getStringDateNow(0)), []byte(getStringDateNow(7))})
	res = stub.MockInvoke("2", [][]byte{[]byte("deleteconsent"), []byte(APPID1), []byte("badconsentid")})
	if res.Status != shim.ERROR{
		t.Log("bad status received, expected: 500 received:"+strconv.FormatInt(int64(res.Status), 10))
		t.FailNow()
	}
	if !strings.Contains(res.Message, errorConsentNotExist){
		t.Log("Bad return message, expected:"+errorConsentNotExist+" reveived:"+string(res.Message))
		t.FailNow()
	}
}

// =====================================================================================================================
// Inactivate a consent with a bad applicationID --> error
// =====================================================================================================================
func TestConsentV2_InactivateConsentWithBadApplicationID(t *testing.T) {
	scc := new(ConsentCC)
	stub := shim.NewMockStub("consentv2", scc)
	res := stub.MockInvoke("1", [][]byte{[]byte("postconsent"), []byte(APPID1), []byte(OWNERID1), []byte(CONSUMERID1), []byte(DATATYPE1), []byte(DATAACCESS1), []byte(getStringDateNow(0)), []byte(getStringDateNow(7))})
	consentID := string(res.Payload)
	res = stub.MockInvoke("2", [][]byte{[]byte("deleteconsent"), []byte("badappid"), []byte(consentID)})
	if res.Status != shim.ERROR{
		t.Log("bad status received, expected: 500 received:"+strconv.FormatInt(int64(res.Status), 10))
		t.FailNow()
	}
	if !strings.Contains(res.Message, errorRemoveConsent){
		t.Log("Bad return message, expected:"+errorRemoveConsent+" reveived:"+string(res.Message))
		t.FailNow()
	}
}

// =====================================================================================================================
// Inactivate a consent with another applicationID --> error
// =====================================================================================================================
func TestConsentV2_InactivateConsentWithAnotherApplicationID(t *testing.T) {
	scc := new(ConsentCC)
	stub := shim.NewMockStub("consentv2", scc)
	res := stub.MockInvoke("1", [][]byte{[]byte("postconsent"), []byte(APPID1), []byte(OWNERID1), []byte(CONSUMERID1), []byte(DATATYPE1), []byte(DATAACCESS1), []byte(getStringDateNow(0)), []byte(getStringDateNow(7))})
	consentID := string(res.Payload)
	stub.MockInvoke("2", [][]byte{[]byte("postconsent"), []byte(APPID2), []byte(OWNERID1), []byte(CONSUMERID1), []byte(DATATYPE1), []byte(DATAACCESS1), []byte(getStringDateNow(0)), []byte(getStringDateNow(7))})
	res = stub.MockInvoke("3", [][]byte{[]byte("deleteconsent"), []byte(APPID2), []byte(consentID)})
	if res.Status != shim.ERROR{
		t.Log("bad status received, expected: 500 received:"+strconv.FormatInt(int64(res.Status), 10))
		t.FailNow()
	}
	if !strings.Contains(res.Message, errorRemoveConsent){
		t.Log("Bad return message, expected:"+errorRemoveConsent+" reveived:"+string(res.Message))
		t.FailNow()
	}
}

// =====================================================================================================================
// Get consent with bad applicationID missing parameter --> error
// =====================================================================================================================
func TestConsentV2_GetConsentWithMissingParameter(t *testing.T) {
	scc := new(ConsentCC)
	stub := shim.NewMockStub("consentv2", scc)
	res := stub.MockInvoke("1", [][]byte{[]byte("postconsent"), []byte(APPID1), []byte(OWNERID1), []byte(CONSUMERID1), []byte(DATATYPE1), []byte(DATAACCESS1), []byte(getStringDateNow(0)), []byte(getStringDateNow(7))})
	consentID := string(res.Payload)
	res = stub.MockInvoke("2", [][]byte{[]byte("getconsent"), []byte(consentID)})
	if res.Status != shim.ERROR{
		t.Log("bad status received, expected: 500 received:"+strconv.FormatInt(int64(res.Status), 10))
		t.FailNow()
	}
	if !strings.Contains(res.Message, errorArgs){
		t.Log("Bad return message, expected:"+errorArgs+" reveived:"+string(res.Message))
		t.FailNow()
	}
}

// =====================================================================================================================
// Get consent with a bad application ID --> error
// =====================================================================================================================
func TestConsentV2_GetConsentWithBadApplicationID(t *testing.T) {
	scc := new(ConsentCC)
	stub := shim.NewMockStub("consentv2", scc)
	res := stub.MockInvoke("1", [][]byte{[]byte("postconsent"), []byte(APPID1), []byte(OWNERID1), []byte(CONSUMERID1), []byte(DATATYPE1), []byte(DATAACCESS1), []byte(getStringDateNow(0)), []byte(getStringDateNow(7))})
	consentID := string(res.Payload)
	res = stub.MockInvoke("2", [][]byte{[]byte("getconsent"), []byte("badappid"), []byte(consentID)})
	if res.Status != shim.ERROR{
		t.Log("bad status received, expected: 500 received:"+strconv.FormatInt(int64(res.Status), 10))
		t.FailNow()
	}
	if !strings.Contains(res.Message, errorConsentNotExist){
		t.Log("Bad return message, expected:"+errorConsentNotExist+" reveived:"+string(res.Message))
		t.FailNow()
	}
}

// =====================================================================================================================
// Get consent with another application ID --> error
// =====================================================================================================================
func TestConsentV2_GetConsentWithAnotherApplicationID(t *testing.T) {
	scc := new(ConsentCC)
	stub := shim.NewMockStub("consentv2", scc)
	res := stub.MockInvoke("1", [][]byte{[]byte("postconsent"), []byte(APPID1), []byte(OWNERID1), []byte(CONSUMERID1), []byte(DATATYPE1), []byte(DATAACCESS1), []byte(getStringDateNow(0)), []byte(getStringDateNow(7))})
	consentID := string(res.Payload)
	stub.MockInvoke("2", [][]byte{[]byte("postconsent"), []byte(APPID2), []byte(OWNERID1), []byte(CONSUMERID1), []byte(DATATYPE1), []byte(DATAACCESS1), []byte(getStringDateNow(0)), []byte(getStringDateNow(7))})
	res = stub.MockInvoke("3", [][]byte{[]byte("getconsent"), []byte(APPID2), []byte(consentID)})
	if res.Status != shim.ERROR{
		t.Log("bad status received, expected: 500 received:"+strconv.FormatInt(int64(res.Status), 10))
		t.FailNow()
	}
	if !strings.Contains(res.Message, errorConsentNotExist){
		t.Log("Bad return message, expected:"+errorConsentNotExist+" reveived:"+string(res.Message))
		t.FailNow()
	}
}
// =====================================================================================================================
// Get inactivate consent --> error
// =====================================================================================================================
func TestConsentV2_GetRemovedConsent(t *testing.T) {
	scc := new(ConsentCC)
	stub := shim.NewMockStub("consentv2", scc)
	res := stub.MockInvoke("1", [][]byte{[]byte("postconsent"), []byte(APPID1), []byte(OWNERID1), []byte(CONSUMERID1), []byte(DATATYPE1), []byte(DATAACCESS1), []byte(getStringDateNow(0)), []byte(getStringDateNow(7))})
	consentID := string(res.Payload)
	res = stub.MockInvoke("2", [][]byte{[]byte("deleteconsent"), []byte(APPID1), []byte(consentID)})
	if res.Status != shim.OK {
		t.Log("deleteconsent", string(res.Message))
		t.FailNow()
	}
	res = stub.MockInvoke("3", [][]byte{[]byte("getconsent"), []byte(APPID1), []byte(consentID)})
	if res.Status == shim.OK {
		t.Log("getconsent: the consent should be inactive")
		t.FailNow()
	}
	res = stub.MockInvoke("4", [][]byte{[]byte("getconsent"), []byte(APPID1), []byte(consentID)})
	if res.Status != shim.ERROR{
		t.Log("bad status received, expected: 500 received:"+strconv.FormatInt(int64(res.Status), 10))
		t.FailNow()
	}
	if !strings.Contains(res.Message, errorConsentNotExist){
		t.Log("Bad return message, expected:"+errorConsentNotExist+" reveived:"+string(res.Message))
		t.FailNow()
	}
}

// =====================================================================================================================
// Get list of all consents empty list
// =====================================================================================================================
func TestConsentV2_GetAllConsentsEmptyList(t *testing.T) {
	scc := new(ConsentCC)
	stub := shim.NewMockStub("consentv2", scc)
	res := stub.MockInvoke("1", [][]byte{[]byte("getconsents"), []byte(APPID1)})
	if res.Status != shim.OK{
		t.Log("bad status received, expected: 200 received:"+strconv.FormatInt(int64(res.Status), 10))
		t.Log("response: "+ string(res.Message))
		t.FailNow()
	}
	consents := make([]consent, 0)
	err := json.Unmarshal(res.Payload, &consents)
	if err != nil {
		t.Log("getconsents", string(res.Payload))
		t.FailNow()
	}
	if len(consents) != 0{
		t.Error("empty list expected, but ",strconv.Itoa(len(consents)), "reveived")
		t.FailNow()
	}
}

// =====================================================================================================================
// Get list of all consents for an application with a missing parameter --> error
// =====================================================================================================================
func TestConsentV2_GetAllConsents4AppIDWithMissingParameter(t *testing.T) {
	scc := new(ConsentCC)
	stub := shim.NewMockStub("consentv2", scc)
	stub.MockInvoke("1", [][]byte{[]byte("postconsent"), []byte(APPID1), []byte(OWNERID1), []byte(CONSUMERID1), []byte(DATATYPE1), []byte(DATAACCESS1), []byte(getStringDateNow(0)), []byte(getStringDateNow(7))})
	stub.MockInvoke("2", [][]byte{[]byte("postconsent"), []byte(APPID1), []byte(OWNERID1), []byte(CONSUMERID2), []byte(DATATYPE1), []byte(DATAACCESS1), []byte(getStringDateNow(0)), []byte(getStringDateNow(7))})
	stub.MockInvoke("3", [][]byte{[]byte("postconsent"), []byte(APPID1), []byte(OWNERID2), []byte(CONSUMERID1), []byte(DATATYPE1), []byte(DATAACCESS1), []byte(getStringDateNow(0)), []byte(getStringDateNow(7))})
	stub.MockInvoke("4", [][]byte{[]byte("postconsent"), []byte(APPID1), []byte(OWNERID2), []byte(CONSUMERID2), []byte(DATATYPE1), []byte(DATAACCESS1), []byte(getStringDateNow(0)), []byte(getStringDateNow(7))})
	res := stub.MockInvoke("6", [][]byte{[]byte("getconsents")})
	if res.Status != shim.ERROR{
		t.Log("bad status received, expected: 500 received:"+strconv.FormatInt(int64(res.Status), 10))
		t.FailNow()
	}
	if !strings.Contains(res.Message, errorArgs){
		t.Log("Bad return message, expected:"+errorArgs+" reveived:"+string(res.Message))
		t.FailNow()
	}
}

// =====================================================================================================================
// Get list of all consents for an application with bad applicationID --> empty list
// =====================================================================================================================
func TestConsentV2_GetAllConsents4AppIDWithBadAppID(t *testing.T) {
	scc := new(ConsentCC)
	stub := shim.NewMockStub("consentv2", scc)
	stub.MockInvoke("1", [][]byte{[]byte("postconsent"), []byte(APPID1), []byte(OWNERID1), []byte(CONSUMERID1), []byte(DATATYPE1), []byte(DATAACCESS1), []byte(getStringDateNow(0)), []byte(getStringDateNow(7))})
	stub.MockInvoke("2", [][]byte{[]byte("postconsent"), []byte(APPID1), []byte(OWNERID1), []byte(CONSUMERID2), []byte(DATATYPE1), []byte(DATAACCESS1), []byte(getStringDateNow(0)), []byte(getStringDateNow(7))})
	stub.MockInvoke("3", [][]byte{[]byte("postconsent"), []byte(APPID1), []byte(OWNERID2), []byte(CONSUMERID1), []byte(DATATYPE1), []byte(DATAACCESS1), []byte(getStringDateNow(0)), []byte(getStringDateNow(7))})
	stub.MockInvoke("4", [][]byte{[]byte("postconsent"), []byte(APPID1), []byte(OWNERID2), []byte(CONSUMERID2), []byte(DATATYPE1), []byte(DATAACCESS1), []byte(getStringDateNow(0)), []byte(getStringDateNow(7))})
	res := stub.MockInvoke("5", [][]byte{[]byte("getconsents"), []byte("badappid")})
	if res.Status != shim.OK{
		t.Log("bad status received, expected: 200 received:"+strconv.FormatInt(int64(res.Status), 10))
		t.Log("response: "+ string(res.Message))
		t.FailNow()
	}
	consents := make([]consent, 0)
	err := json.Unmarshal(res.Payload, &consents)
	if err != nil {
		t.Log("getconsents", string(res.Payload))
		t.FailNow()
	}
	if len(consents) != 0{
		t.Error("empty list expected, but ",strconv.Itoa(len(consents)), "reveived")
		t.FailNow()
	}
}

// =====================================================================================================================
// Get list of all consents for an owner empty list
// =====================================================================================================================
func TestConsentV2_GetOwnerConsentsEmptyList(t *testing.T) {
	scc := new(ConsentCC)
	stub := shim.NewMockStub("consentv2", scc)
	res := stub.MockInvoke("1", [][]byte{[]byte("getownerconsents"), []byte(APPID1), []byte(OWNERID1)})
	if res.Status != shim.OK {
		t.Log("bad status received, expected: 200 received:"+strconv.FormatInt(int64(res.Status), 10))
		t.Log("response: "+ string(res.Message))
		t.FailNow()
	}
	consents := make([]consent, 0)
	err := json.Unmarshal(res.Payload, &consents)
	if err != nil {
		t.Log("getconsent", string(res.Payload))
		t.FailNow()
	}
	if len(consents) != 0{
		t.Error("0 expected, but ",strconv.Itoa(len(consents)), "reveived")
		t.FailNow()
	}
}

// =====================================================================================================================
// Get list of all consents for an owner with missing parameter --> error
// =====================================================================================================================
func TestConsentV2_GetOwnerConsentsWithMissingParameter(t *testing.T) {
	scc := new(ConsentCC)
	stub := shim.NewMockStub("consentv2", scc)
	stub.MockInvoke("1", [][]byte{[]byte("postconsent"), []byte(APPID1), []byte(OWNERID1), []byte(CONSUMERID1), []byte(DATATYPE1), []byte(DATAACCESS1), []byte(getStringDateNow(0)), []byte(getStringDateNow(7))})
	stub.MockInvoke("2", [][]byte{[]byte("postconsent"), []byte(APPID1), []byte(OWNERID1), []byte(CONSUMERID2), []byte(DATATYPE1), []byte(DATAACCESS1), []byte(getStringDateNow(0)), []byte(getStringDateNow(7))})
	res := stub.MockInvoke("3", [][]byte{[]byte("getownerconsents"), []byte(APPID1)})
	if res.Status != shim.ERROR{
		t.Log("bad status received, expected: 500 received:"+strconv.FormatInt(int64(res.Status), 10))
		t.FailNow()
	}
	if !strings.Contains(res.Message, errorArgs){
		t.Log("Bad return message, expected:"+errorArgs+" reveived:"+string(res.Message))
		t.FailNow()
	}
}

// =====================================================================================================================
// Get list of all consents for an owner with bad application ID --> empty list
// =====================================================================================================================
func TestConsentV2_GetOwnerConsentsWithBadApplicationID(t *testing.T) {
	scc := new(ConsentCC)
	stub := shim.NewMockStub("consentv2", scc)
	stub.MockInvoke("1", [][]byte{[]byte("postconsent"), []byte(APPID1), []byte(OWNERID1), []byte(CONSUMERID1), []byte(DATATYPE1), []byte(DATAACCESS1), []byte(getStringDateNow(0)), []byte(getStringDateNow(7))})
	stub.MockInvoke("2", [][]byte{[]byte("postconsent"), []byte(APPID1), []byte(OWNERID1), []byte(CONSUMERID2), []byte(DATATYPE1), []byte(DATAACCESS1), []byte(getStringDateNow(0)), []byte(getStringDateNow(7))})
	res := stub.MockInvoke("3", [][]byte{[]byte("getownerconsents"), []byte("badappid"), []byte(OWNERID1)})
	if res.Status != shim.OK{
		t.Log("bad status received, expected: 200 received:"+strconv.FormatInt(int64(res.Status), 10))
		t.Log("response: "+ string(res.Message))
		t.FailNow()
	}
	consents := make([]consent, 0)
	err := json.Unmarshal(res.Payload, &consents)
	if err != nil {
		t.Log("getconsents", string(res.Payload))
		t.FailNow()
	}
	if len(consents) != 0{
		t.Error("empty list expected, but ",strconv.Itoa(len(consents)), "reveived")
		t.FailNow()
	}
}

// =====================================================================================================================
// Get list of all consents for an owner with other application ID --> empty list
// =====================================================================================================================
func TestConsentV2_GetOwnerConsentsWithOtherApplicationID(t *testing.T) {
	scc := new(ConsentCC)
	stub := shim.NewMockStub("consentv2", scc)
	stub.MockInvoke("1", [][]byte{[]byte("postconsent"), []byte(APPID1), []byte(OWNERID1), []byte(CONSUMERID1), []byte(DATATYPE1), []byte(DATAACCESS1), []byte(getStringDateNow(0)), []byte(getStringDateNow(7))})
	stub.MockInvoke("2", [][]byte{[]byte("postconsent"), []byte(APPID1), []byte(OWNERID1), []byte(CONSUMERID2), []byte(DATATYPE1), []byte(DATAACCESS1), []byte(getStringDateNow(0)), []byte(getStringDateNow(7))})
	stub.MockInvoke("3", [][]byte{[]byte("postconsent"), []byte(APPID2), []byte(OWNERID2), []byte(CONSUMERID1), []byte(DATATYPE1), []byte(DATAACCESS1), []byte(getStringDateNow(0)), []byte(getStringDateNow(7))})
	res := stub.MockInvoke("4", [][]byte{[]byte("getownerconsents"), []byte(APPID2), []byte(OWNERID1)})
	if res.Status != shim.OK{
		t.Log("bad status received, expected: 200 received:"+strconv.FormatInt(int64(res.Status), 10))
		t.Log("response: "+ string(res.Message))
		t.FailNow()
	}
	consents := make([]consent, 0)
	err := json.Unmarshal(res.Payload, &consents)
	if err != nil {
		t.Log("getconsents", string(res.Payload))
		t.FailNow()
	}
	if len(consents) != 0{
		t.Error("empty list expected, but ",strconv.Itoa(len(consents)), "reveived")
		t.FailNow()
	}
}

// =====================================================================================================================
// Get list of all consents for an owner with bad owner ID --> empty list
// =====================================================================================================================
func TestConsentV2_GetOwnerConsentsWithBadOwnerID(t *testing.T) {
	scc := new(ConsentCC)
	stub := shim.NewMockStub("consentv2", scc)
	stub.MockInvoke("1", [][]byte{[]byte("postconsent"), []byte(APPID1), []byte(OWNERID1), []byte(CONSUMERID1), []byte(DATATYPE1), []byte(DATAACCESS1), []byte(getStringDateNow(0)), []byte(getStringDateNow(7))})
	stub.MockInvoke("2", [][]byte{[]byte("postconsent"), []byte(APPID1), []byte(OWNERID1), []byte(CONSUMERID2), []byte(DATATYPE1), []byte(DATAACCESS1), []byte(getStringDateNow(0)), []byte(getStringDateNow(7))})
	stub.MockInvoke("3", [][]byte{[]byte("postconsent"), []byte(APPID2), []byte(OWNERID2), []byte(CONSUMERID1), []byte(DATATYPE1), []byte(DATAACCESS1), []byte(getStringDateNow(0)), []byte(getStringDateNow(7))})
	res := stub.MockInvoke("4", [][]byte{[]byte("getownerconsents"), []byte(APPID2), []byte("badownerid")})
	if res.Status != shim.OK{
		t.Log("bad status received, expected: 200 received:"+strconv.FormatInt(int64(res.Status), 10))
		t.Log("response: "+ string(res.Message))
		t.FailNow()
	}
	consents := make([]consent, 0)
	err := json.Unmarshal(res.Payload, &consents)
	if err != nil {
		t.Log("getconsents", string(res.Payload))
		t.FailNow()
	}
	if len(consents) != 0{
		t.Error("empty list expected, but ",strconv.Itoa(len(consents)), "reveived")
		t.FailNow()
	}
}


// =====================================================================================================================
// Get list of all consents for a consumer empty list
// =====================================================================================================================
func TestConsentV2_GetConsumerEmptyList(t *testing.T) {
	scc := new(ConsentCC)
	stub := shim.NewMockStub("consentv2", scc)
	res := stub.MockInvoke("1", [][]byte{[]byte("getconsumerconsents"), []byte(APPID1), []byte(CONSUMERID1)})
	if res.Status != shim.OK {
		t.Log("bad status received, expected: 200 received:"+strconv.FormatInt(int64(res.Status), 10))
		t.Log("response: "+ string(res.Message))
		t.FailNow()
	}
	consents := make([]consent, 0)
	err := json.Unmarshal(res.Payload, &consents)
	if err != nil {
		t.Log("getconsent", string(res.Payload))
		t.FailNow()
	}
	if len(consents) != 0{
		t.Error("0 expected, but ",strconv.Itoa(len(consents)), "reveived")
		t.FailNow()
	}
}


// =====================================================================================================================
// Get list of all consents for a consumer with a missing parameter
// =====================================================================================================================
func TestConsentV2_GetConsumerConsentsWithMissingParameter(t *testing.T) {
	scc := new(ConsentCC)
	stub := shim.NewMockStub("consentv2", scc)
	stub.MockInvoke("1", [][]byte{[]byte("postconsent"), []byte(APPID1), []byte(OWNERID1), []byte(CONSUMERID1), []byte(DATATYPE1), []byte(DATAACCESS1), []byte(getStringDateNow(0)), []byte(getStringDateNow(7))})
	stub.MockInvoke("2", [][]byte{[]byte("postconsent"), []byte(APPID1), []byte(OWNERID2), []byte(CONSUMERID1), []byte(DATATYPE1), []byte(DATAACCESS1), []byte(getStringDateNow(0)), []byte(getStringDateNow(7))})
	stub.MockInvoke("3", [][]byte{[]byte("postconsent"), []byte(APPID2), []byte(OWNERID1), []byte(CONSUMERID1), []byte(DATATYPE1), []byte(DATAACCESS1), []byte(getStringDateNow(0)), []byte(getStringDateNow(7))})
	stub.MockInvoke("4", [][]byte{[]byte("postconsent"), []byte(APPID1), []byte(OWNERID1), []byte(CONSUMERID2), []byte(DATATYPE1), []byte(DATAACCESS1), []byte(getStringDateNow(0)), []byte(getStringDateNow(7))})

	res := stub.MockInvoke("5", [][]byte{[]byte("getconsumerconsents"), []byte(APPID1)})
	if res.Status != shim.ERROR{
		t.Log("bad status received, expected: 500 received:"+strconv.FormatInt(int64(res.Status), 10))
		t.FailNow()
	}
	if !strings.Contains(res.Message, errorArgs){
		t.Log("Bad return message, expected:"+errorArgs+" reveived:"+string(res.Message))
		t.FailNow()
	}
}

// =====================================================================================================================
// Get list of all consents for a consumer with bad application ID
// =====================================================================================================================
func TestConsentV2_GetConsumerConsentsWithBadApplicationID(t *testing.T) {
	scc := new(ConsentCC)
	stub := shim.NewMockStub("consentv2", scc)
	stub.MockInvoke("1", [][]byte{[]byte("postconsent"), []byte(APPID1), []byte(OWNERID1), []byte(CONSUMERID1), []byte(DATATYPE1), []byte(DATAACCESS1), []byte(getStringDateNow(0)), []byte(getStringDateNow(7))})
	stub.MockInvoke("2", [][]byte{[]byte("postconsent"), []byte(APPID1), []byte(OWNERID2), []byte(CONSUMERID1), []byte(DATATYPE1), []byte(DATAACCESS1), []byte(getStringDateNow(0)), []byte(getStringDateNow(7))})
	stub.MockInvoke("3", [][]byte{[]byte("postconsent"), []byte(APPID2), []byte(OWNERID1), []byte(CONSUMERID1), []byte(DATATYPE1), []byte(DATAACCESS1), []byte(getStringDateNow(0)), []byte(getStringDateNow(7))})
	stub.MockInvoke("4", [][]byte{[]byte("postconsent"), []byte(APPID1), []byte(OWNERID1), []byte(CONSUMERID2), []byte(DATATYPE1), []byte(DATAACCESS1), []byte(getStringDateNow(0)), []byte(getStringDateNow(7))})

	res := stub.MockInvoke("5", [][]byte{[]byte("getconsumerconsents"), []byte("badappid"), []byte(CONSUMERID1)})
	if res.Status != shim.OK{
		t.Log("bad status received, expected: 200 received:"+strconv.FormatInt(int64(res.Status), 10))
		t.Log("response: "+ string(res.Message))
		t.FailNow()
	}
	consents := make([]consent, 0)
	err := json.Unmarshal(res.Payload, &consents)
	if err != nil {
		t.Log("getconsents", string(res.Payload))
		t.FailNow()
	}
	if len(consents) != 0{
		t.Error("empty list expected, but ",strconv.Itoa(len(consents)), "reveived")
		t.FailNow()
	}
}

// =====================================================================================================================
// Get list of all consents for a consumer with other application ID
// =====================================================================================================================
func TestConsentV2_GetConsumerConsentsWithOtherApplicationID(t *testing.T) {
	scc := new(ConsentCC)
	stub := shim.NewMockStub("consentv2", scc)
	stub.MockInvoke("1", [][]byte{[]byte("postconsent"), []byte(APPID1), []byte(OWNERID1), []byte(CONSUMERID1), []byte(DATATYPE1), []byte(DATAACCESS1), []byte(getStringDateNow(0)), []byte(getStringDateNow(7))})
	stub.MockInvoke("2", [][]byte{[]byte("postconsent"), []byte(APPID1), []byte(OWNERID2), []byte(CONSUMERID1), []byte(DATATYPE1), []byte(DATAACCESS1), []byte(getStringDateNow(0)), []byte(getStringDateNow(7))})
	stub.MockInvoke("3", [][]byte{[]byte("postconsent"), []byte(APPID2), []byte(OWNERID1), []byte(CONSUMERID1), []byte(DATATYPE1), []byte(DATAACCESS1), []byte(getStringDateNow(0)), []byte(getStringDateNow(7))})
	stub.MockInvoke("4", [][]byte{[]byte("postconsent"), []byte(APPID1), []byte(OWNERID1), []byte(CONSUMERID2), []byte(DATATYPE1), []byte(DATAACCESS1), []byte(getStringDateNow(0)), []byte(getStringDateNow(7))})

	res := stub.MockInvoke("5", [][]byte{[]byte("getconsumerconsents"), []byte(APPID2), []byte(CONSUMERID2)})
	if res.Status != shim.OK{
		t.Log("bad status received, expected: 200 received:"+strconv.FormatInt(int64(res.Status), 10))
		t.Log("response: "+ string(res.Message))
		t.FailNow()
	}
	consents := make([]consent, 0)
	err := json.Unmarshal(res.Payload, &consents)
	if err != nil {
		t.Log("getconsents", string(res.Payload))
		t.FailNow()
	}
	if len(consents) != 0{
		t.Error("empty list expected, but ",strconv.Itoa(len(consents)), "reveived")
		t.FailNow()
	}
}

// =====================================================================================================================
// Get list of all consents for a consumer with bad owner ID
// =====================================================================================================================
func TestConsentV2_GetConsumerConsentsWithBadOwnerID(t *testing.T) {
	scc := new(ConsentCC)
	stub := shim.NewMockStub("consentv2", scc)
	stub.MockInvoke("1", [][]byte{[]byte("postconsent"), []byte(APPID1), []byte(OWNERID1), []byte(CONSUMERID1), []byte(DATATYPE1), []byte(DATAACCESS1), []byte(getStringDateNow(0)), []byte(getStringDateNow(7))})
	stub.MockInvoke("2", [][]byte{[]byte("postconsent"), []byte(APPID1), []byte(OWNERID2), []byte(CONSUMERID1), []byte(DATATYPE1), []byte(DATAACCESS1), []byte(getStringDateNow(0)), []byte(getStringDateNow(7))})
	stub.MockInvoke("3", [][]byte{[]byte("postconsent"), []byte(APPID2), []byte(OWNERID1), []byte(CONSUMERID1), []byte(DATATYPE1), []byte(DATAACCESS1), []byte(getStringDateNow(0)), []byte(getStringDateNow(7))})
	stub.MockInvoke("4", [][]byte{[]byte("postconsent"), []byte(APPID1), []byte(OWNERID1), []byte(CONSUMERID2), []byte(DATATYPE1), []byte(DATAACCESS1), []byte(getStringDateNow(0)), []byte(getStringDateNow(7))})

	res := stub.MockInvoke("5", [][]byte{[]byte("getconsumerconsents"), []byte(APPID1), []byte("badownerid")})
	if res.Status != shim.OK{
		t.Log("bad status received, expected: 200 received:"+strconv.FormatInt(int64(res.Status), 10))
		t.Log("response: "+ string(res.Message))
		t.FailNow()
	}
	consents := make([]consent, 0)
	err := json.Unmarshal(res.Payload, &consents)
	if err != nil {
		t.Log("getconsents", string(res.Payload))
		t.FailNow()
	}
	if len(consents) != 0{
		t.Error("empty list expected, but ",strconv.Itoa(len(consents)), "reveived")
		t.FailNow()
	}
}


// =====================================================================================================================
// Delete all consents for an application with missing parameter
// =====================================================================================================================
func TestConsentV2_DeleteConsents4AppIDWithMissingParameter(t *testing.T) {
	scc := new(ConsentCC)
	stub := shim.NewMockStub("consentv2", scc)
	stub.MockInvoke("1", [][]byte{[]byte("postconsent"), []byte(APPID1), []byte(OWNERID1), []byte(CONSUMERID1), []byte(DATATYPE1), []byte(DATAACCESS1), []byte(getStringDateNow(0)), []byte(getStringDateNow(7))})
	stub.MockInvoke("2", [][]byte{[]byte("postconsent"), []byte(APPID1), []byte(OWNERID1), []byte(CONSUMERID2), []byte(DATATYPE1), []byte(DATAACCESS1), []byte(getStringDateNow(0)), []byte(getStringDateNow(7))})
	stub.MockInvoke("3", [][]byte{[]byte("postconsent"), []byte(APPID1), []byte(OWNERID2), []byte(CONSUMERID1), []byte(DATATYPE1), []byte(DATAACCESS1), []byte(getStringDateNow(0)), []byte(getStringDateNow(7))})
	stub.MockInvoke("4", [][]byte{[]byte("postconsent"), []byte(APPID1), []byte(OWNERID2), []byte(CONSUMERID2), []byte(DATATYPE1), []byte(DATAACCESS1), []byte(getStringDateNow(0)), []byte(getStringDateNow(7))})
	stub.MockInvoke("5", [][]byte{[]byte("postconsent"), []byte(APPID2), []byte(OWNERID2), []byte(CONSUMERID2), []byte(DATATYPE1), []byte(DATAACCESS1), []byte(getStringDateNow(0)), []byte(getStringDateNow(7))})
	res := stub.MockInvoke("6", [][]byte{[]byte("deleteconsents")})

	if res.Status != shim.ERROR{
		t.Log("bad status received, expected: 500 received:"+strconv.FormatInt(int64(res.Status), 10))
		t.FailNow()
	}
	if !strings.Contains(res.Message, errorArgs){
		t.Log("Bad return message, expected:"+errorArgs+" reveived:"+string(res.Message))
		t.FailNow()
	}
}


// =====================================================================================================================
// Is consent exist with missing parameter
// =====================================================================================================================
func TestConsentV2_IsConsentWithMissingParameter(t *testing.T) {
	scc := new(ConsentCC)
	stub := shim.NewMockStub("consentv2", scc)
	stub.MockInvoke("1", [][]byte{[]byte("postconsent"), []byte(APPID1), []byte(OWNERID1), []byte(CONSUMERID1), []byte(DATATYPE1), []byte(DATAACCESS1), []byte(getStringDateNow(0)), []byte(getStringDateNow(7))})
	res := stub.MockInvoke("1", [][]byte{[]byte("isconsent"), []byte(APPID1), []byte(CONSUMERID1), []byte(DATATYPE1), []byte(DATAACCESS1)})

	if res.Status != shim.ERROR{
		t.Log("bad status received, expected: 500 received:"+strconv.FormatInt(int64(res.Status), 10))
		t.FailNow()
	}
	if !strings.Contains(res.Message, errorArgs){
		t.Log("Bad return message, expected:"+errorArgs+" reveived:"+string(res.Message))
		t.FailNow()
	}
}


// =====================================================================================================================
// Is consent exist with one different parameter
// =====================================================================================================================
func TestConsentV2_IsConsentWithanotherApplicationID(t *testing.T) {
	scc := new(ConsentCC)
	stub := shim.NewMockStub("consentv2", scc)
	stub.MockInvoke("1", [][]byte{[]byte("postconsent"), []byte(APPID1), []byte(OWNERID1), []byte(CONSUMERID1), []byte(DATATYPE1), []byte(DATAACCESS1), []byte(getStringDateNow(0)), []byte(getStringDateNow(7))})

	res := stub.MockInvoke("1", [][]byte{[]byte("isconsent"), []byte(APPID2), []byte(OWNERID1), []byte(CONSUMERID1), []byte(DATATYPE1), []byte(DATAACCESS1)})
	if res.Status != shim.OK {
		t.Log("bad status received, expected: 200 received:"+strconv.FormatInt(int64(res.Status), 10))
		t.Log("response: "+ string(res.Message))
		t.FailNow()
	}
	response := string(res.Payload)
	if response != NOT_AUTHORIZED{
		t.Log(NOT_AUTHORIZED, "expected, but ",response, "reveived")
		t.FailNow()
	}
}

// =====================================================================================================================
// Is consent exist with one different dataType
// =====================================================================================================================
func TestConsentV2_IsConsentWithDifferentDataTypeParameter(t *testing.T) {
	scc := new(ConsentCC)
	stub := shim.NewMockStub("consentv2", scc)
	stub.MockInvoke("1", [][]byte{[]byte("postconsent"), []byte(APPID1), []byte(OWNERID1), []byte(CONSUMERID1), []byte(DATATYPE1), []byte(DATAACCESS1), []byte(getStringDateNow(0)), []byte(getStringDateNow(7))})

	res := stub.MockInvoke("1", [][]byte{[]byte("isconsent"), []byte(APPID1), []byte(OWNERID1), []byte(CONSUMERID1), []byte(DATATYPE2), []byte(DATAACCESS1)})
	if res.Status != shim.OK {
		t.Log("bad status received, expected: 200 received:"+strconv.FormatInt(int64(res.Status), 10))
		t.Log("response: "+ string(res.Message))
		t.FailNow()
	}
	response := string(res.Payload)
	if response != NOT_AUTHORIZED{
		t.Log(NOT_AUTHORIZED, "expected, but ",response, "reveived")
		t.FailNow()
	}
}

// =====================================================================================================================
// Is consent exist with one different dataAccess
// =====================================================================================================================
func TestConsentV2_IsConsentWithOneDifferentDataAccess(t *testing.T) {
	scc := new(ConsentCC)
	stub := shim.NewMockStub("consentv2", scc)
	stub.MockInvoke("1", [][]byte{[]byte("postconsent"), []byte(APPID1), []byte(OWNERID1), []byte(CONSUMERID1), []byte(DATATYPE1), []byte(DATAACCESS1), []byte(getStringDateNow(0)), []byte(getStringDateNow(7))})

	res := stub.MockInvoke("1", [][]byte{[]byte("isconsent"), []byte(APPID1), []byte(OWNERID1), []byte(CONSUMERID1), []byte(DATATYPE1), []byte(DATAACCESS2)})
	if res.Status != shim.OK {
		t.Log("bad status received, expected: 200 received:"+strconv.FormatInt(int64(res.Status), 10))
		t.Log("response: "+ string(res.Message))
		t.FailNow()
	}
	response := string(res.Payload)
	if response != NOT_AUTHORIZED{
		t.Log(NOT_AUTHORIZED, "expected, but ",response, "reveived")
		t.FailNow()
	}
}

// =====================================================================================================================
// Is consent exist with old period
// =====================================================================================================================
func TestConsentV2_IsConsentWithOldPeriod(t *testing.T) {
	scc := new(ConsentCC)
	stub := shim.NewMockStub("consentv2", scc)
	stub.MockInvoke("1", [][]byte{[]byte("postconsent"), []byte(APPID1), []byte(OWNERID1), []byte(CONSUMERID1), []byte(DATATYPE1), []byte(DATAACCESS1), []byte(getStringDateNow(-5)), []byte(getStringDateNow(-2))})

	res := stub.MockInvoke("1", [][]byte{[]byte("isconsent"), []byte(APPID1), []byte(OWNERID1), []byte(CONSUMERID1), []byte(DATATYPE1), []byte(DATAACCESS1)})
	if res.Status != shim.OK {
		t.Log("bad status received, expected: 200 received:"+strconv.FormatInt(int64(res.Status), 10))
		t.Log("response: "+ string(res.Message))
		t.FailNow()
	}
	response := string(res.Payload)
	if response != NOT_AUTHORIZED{
		t.Log(NOT_AUTHORIZED, "expected, but ",response, "reveived")
		t.FailNow()
	}
}

// =====================================================================================================================
// Is consent exist with future period
// =====================================================================================================================
func TestConsentV2_IsConsentWithFuturePeriod(t *testing.T) {
	scc := new(ConsentCC)
	stub := shim.NewMockStub("consentv2", scc)
	stub.MockInvoke("1", [][]byte{[]byte("postconsent"), []byte(APPID1), []byte(OWNERID1), []byte(CONSUMERID1), []byte(DATATYPE1), []byte(DATAACCESS1), []byte(getStringDateNow(1)), []byte(getStringDateNow(7))})

	res := stub.MockInvoke("1", [][]byte{[]byte("isconsent"), []byte(APPID1), []byte(OWNERID1), []byte(CONSUMERID1), []byte(DATATYPE1), []byte(DATAACCESS1)})
	if res.Status != shim.OK {
		t.Log("bad status received, expected: 200 received:"+strconv.FormatInt(int64(res.Status), 10))
		t.Log("response: "+ string(res.Message))
		t.FailNow()
	}
	response := string(res.Payload)
	if response != NOT_AUTHORIZED{
		t.Log(NOT_AUTHORIZED, "expected, but ",response, "reveived")
		t.FailNow()
	}
}
func getStringDateNow(nbdaysafter time.Duration) string{
	t := time.Now().Add(nbdaysafter * 24 * time.Hour)
	return t.Format("2006-01-02")
}

func TestMain(m *testing.M) {
	InitLogger()
	code := m.Run()
	os.Exit(code)
}
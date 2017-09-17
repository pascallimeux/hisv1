/**********************************************************
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
**********************************************************/

// =====================================================================================================================
// Author : Pascal Limeux Orange Labs
// Purpose: Providing consent management interfaces.
// to deploy:
// ./peer chaincode install -n consent -v 1.0 -p github.com/pascallimeux/consent/consentv2 -o 127.0.0.1:7050
// ./peer chaincode instantiate -n consent -v 1.0 -C mch -c '{"Args":[]}' -p github.com/pascallimeux/consent/consentv2
//  -o 127.0.0.1:7050
// =====================================================================================================================

package main

import (
	"fmt"
	"errors"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"time"
	"encoding/json"
	"strconv"
)


// =====================================================================================================================
// Chaincode constantes
// =====================================================================================================================
const (
	// chaincode constants
	LOGGERMODULE   = "consent"
	VERSION        = "Orange Consent Application chaincode ver 3 Dated 2017-05-22"
	AUTHORIZED     = "True"
	NOT_AUTHORIZED = "False"

	//Chainccode index
	indexApp       = "app~id" 			// to get all consents for appID
	indexOwner     = "app~owner~id" 		// to get all consents for appID and ownerID
	indexConsumer  = "app~consumer~id" 		// to get all consents for appID and consumerID
	indexConsumerOwner = "app~consumer~owner~id" 	// to get all consents for appID consumerID and ownerID
	indexIsConsent = "app~isconsent"		// to check if a consent exist

	// Chaincode errors
	errorArgs                 = "Incorrect number of arguments."
	errorBadFunctionName      = "Invalid function, expecting \"postconsent\" \"removeconsent\" " +
				    "\"resetconsents\" \"getconsent\" \"getownerconsents\" \"getconsumerconsents\" " +
				    "\"getconsents\" \"isconsent\" \"getversion\""
	errorCreateConsent        = "Create consent!"
	errorGetConsent           = "Get consent:"
	errorConsentNotExist      = "Consent does not exist:"
	errorRemoveConsent        = "remove a consent:"
	errorRemoveConsent4App    = "Remove all consent for appID:"
	errorGetConsents4Owner    = "Get list of consents for ownerID:"
	errorGetConsents4Consumer = "Get list of consents for consumerID:"
	errorGetConsents4ConsumerOwner = "Get list of consents for consumerID and ownerID:"
	errorGetConsents4AppID    = "Get list of consents for appID:"
	errorGetConsent4Params    = "Get consent for parms:"
	errorDateBegin            = "Dt_begin format error:"
	errorDateEnd              = "Dt_end format error:"
	errorPeriod               = "Period not valid from:"
)
var logger = shim.NewLogger(LOGGERMODULE)
//var logger = logging.MustGetLogger(LOGGERMODULE)

type ConsentCC struct {
}

// =====================================================================================================================
// AppID:      string: id of the client application
// ConsentID:  string: id of the record allow to identify a consent
// ConsumerID: string: id of the data consumer
// OwnerID:    string: id of the data owner
// DataType:   string: type of data ex: ('BC'-->Body composition, 'BM'--> Body measurement, 'BP'-->Bloodpressure,
// 					 'WS'-->Weightscale, 'CGM'-->Continue glucose monitoring, 'HR'-->Heart rate,
// 					 'BGM -->Blood glucose monitoring, 'CAR'-->Cardio vascular and fitness)
// dataAccess: string: type of data access ex: ('C'-->Create, 'R'-->Read, 'U'-->Update, 'D'-->Delete, 'L'-->List )
// Dt_begin:   date:   starting date of the consent  (the date format is: yyyy-mm-dd)
// Dt_end:     date:   ending date of the consent (the date format is: yyyy-mm-dd)
// =====================================================================================================================
type consent struct {
	AppID 		string     `json:"appid"`
	ConsentID      	string     `json:"consentid"`
	OwnerID       	string     `json:"ownerid"`
	ConsumerID      string     `json:"consumerid"`
	DataType      	string     `json:"datatype"`
	DataAccess      string     `json:"dataaccess"`
	Dt_begin      	time.Time  `json:"dtbegin"`
	Dt_end       	time.Time  `json:"dtend"`
}

// =====================================================================================================================
// Init - Initializes chaincode
// =====================================================================================================================
func (c *ConsentCC) Init(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Debug("Init() : calling method -")
	return shim.Success(nil)
}

// =====================================================================================================================
// Invoke - Entry point for Invocations
// =====================================================================================================================
func (c *ConsentCC) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	logger.Debug("Invoke("+function+") : calling method -")
	switch function {
	case "postconsent":
		return c.createConsent(stub, args)
	case "deleteconsent":
		return c.deleteConsent(stub, args)
	case "deleteconsents" :
		return c.deleteConsents4AppID(stub, args)
	case "getconsent" :
		return c.getConsent(stub, args)
	case "getownerconsents" :
		return c.getOwnerConsents(stub, args)
	case "getconsumerconsents" :
		return c.getConsumerConsents(stub, args)
	case "getconsumerownerconsents" :
		return c.getConsumerOwnerConsents(stub, args)
	case "getconsents" :
		return c.getConsents4AppID(stub, args)
	case "isconsent" :
		return c.isConsent(stub, args)
	case "getversion" :
		return c.getVersion(stub, args)
	default:
		return shim.Error(buildError(errorBadFunctionName+function))
	}
}

// =====================================================================================================================
// getVersion - Get the version of the SmartContract
// example:
// ./peer chaincode invoke -C mch -n consent -c '{"Args":["getversion"]}' -o 127.0.0.1:7050
// return the version of this smartcontract
// =====================================================================================================================
func (c *ConsentCC)getVersion(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	logger.Debug("getVersion() : calling method -")
	if len(args) != 0 {
		errStr :=errorArgs+" None argument expecting!"
		return shim.Error(buildError(errStr))
	}
	valAsBytes := []byte(VERSION)
	return shim.Success(valAsBytes)
}

// =====================================================================================================================
// Create a Consent
// example:
// ./peer chaincode invoke -C mch -n consent -c '{"Args":["postconsent","APPID","OWNERID","CONSUMERID","DATATYPE",
// 							"DATAACCESS", "DT_BEGIN", "DT_END"]}' -o 127.0.0.1:7050
// return the consentID
// =====================================================================================================================
func (c *ConsentCC)createConsent(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 7 {
		errStr := errorArgs+" expecting appID, ownerID, consumerID, dataType, dataAccess, dt_begin, dt_end!"
		return shim.Error(buildError(errStr))
	}
	appID      := args[0]
	consentID  := stub.GetTxID()
	ownerID    := args[1]
	consumerID := args[2]
	dataType   := args[3]
	dataAccess := args[4]
	logger.Debug("createConsent(AppID=" + appID + ", consentID=" + consentID + ", Ownerid=" + ownerID + "  Consumerid=" + consumerID + ", Datatype=" + dataType + ", Dataaccess=" +
		dataAccess + ", Dt_begin=" + args[5] + ", Dt_end=" + args[6] + ") : calling method -")
	dt_begin, dt_end, err := checkDates(args[5], args[6])
	if err != nil {
		return shim.Error(buildError(err.Error()))
	}
	consent := &consent{appID, consentID, ownerID, consumerID, dataType, dataAccess, dt_begin, dt_end}
	consentSONasBytes, err := json.Marshal(consent)
	if err != nil {
		return shim.Error(buildError(errorCreateConsent))
	}
	err = stub.PutState(consentID, consentSONasBytes)
	if err != nil {
		return shim.Error(buildError(errorCreateConsent))
	}

	err = createIndex(stub, *consent)
	if err != nil {
		return shim.Error(buildError(errorCreateConsent))
	}
	valAsBytes := []byte(consentID)
	return shim.Success(valAsBytes)
}

// =====================================================================================================================
// Get a Consent from appID and consentID
// example:
// ./peer chaincode invoke -C mch -n consent -c '{"Args":["getconsent","APPID","CONSENTID"]}' -o 127.0.0.1:7050
// =====================================================================================================================
func (c *ConsentCC)getConsent(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		errStr := errorArgs+" Expecting appID, consentID!"
		return shim.Error(buildError(errStr))
	}
	appID     := args[0]
	consentID := args[1]
	logger.Debug("getConsent(Appid=" + appID + ", ConsentID=" + consentID +") : calling method -")

	valAsBytes, err := stub.GetState(consentID) //get the consent from chaincode state
	if err != nil {
		logger.Error("Failed to get consent: " + consentID +" "+err.Error())
		return shim.Error(buildError(errorGetConsent+ consentID))
	} else if valAsBytes == nil {
		return shim.Error(buildError(errorConsentNotExist+ consentID ))
	}
	var consent consent
	err = json.Unmarshal(valAsBytes, &consent)
	if err != nil {
		logger.Error("Failed to unmarshal state for " + consentID +" "+err.Error())
		return shim.Error(buildError(errorGetConsent+ consentID))
	}
	if appID != consent.AppID {
		logger.Error("Consent does not exist: " + consentID + " for this AppID:" + appID)
		return shim.Error(buildError(errorConsentNotExist+ consentID ))
	}
	return shim.Success(valAsBytes)
}

// =====================================================================================================================
// Inactivate a Consent
// example:
// ./peer chaincode invoke -C mch -n consent -c '{"Args":["removeconsent","APPID","OWNERID"]}' -o 127.0.0.1:7050
// =====================================================================================================================
func (c *ConsentCC)deleteConsent(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		errStr := errorArgs+" Expecting appID, consentID!"
		return shim.Error(buildError(errStr))
	}
	appID     := args[0]
	consentID := args[1]
	logger.Debug("deleteConsent(Appid="+ appID+ ", ConsentID="+ consentID+") : calling method -")

	consentAsBytes, err := stub.GetState(consentID)
	if err != nil {
		logger.Error("Failed to get consent:" + err.Error())
		return shim.Error(buildError(errorGetConsent))
	} else if consentAsBytes == nil {
		return shim.Error(buildError(errorConsentNotExist))
	}
	var consent consent
	err = json.Unmarshal(consentAsBytes, &consent)
	if err != nil {
		logger.Error("Failed to unmarshal state for " + consentID +" "+err.Error())
		return shim.Error(buildError(errorGetConsent+ consentID))
	}
	if consent.AppID != appID{
		logger.Error("Failed to remove consent, bad appid:" + appID)
		return shim.Error(buildError(errorRemoveConsent))
	}
	err = deleteConsent(stub, consentID)
	if err != nil {
		return shim.Error(err.Error())
		return shim.Error(buildError(errorRemoveConsent))
	}
	return shim.Success(nil)
}


// =====================================================================================================================
// resetConsents - Remove all consents
// example:
// ./peer chaincode invoke -C mch -n consent -c '{"Args":["resetconsents","APPID"]}' -o 127.0.0.1:7050
// =====================================================================================================================
func (c *ConsentCC)deleteConsents4AppID(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		errStr := errorArgs+" Expecting appID!"
		return shim.Error(buildError(errStr))
	}
	appID := args[0]
	logger.Debug("deleteConsents4AppID(Appid="+ appID + ") : calling method -")
	consents, err := getConsentsByIndex(stub, indexApp, []string{appID})
	if err != nil {
		errStr := err.Error()
		logger.Error(errStr)
		return shim.Error(errStr)
	}
	logger.Debug("delete "+ strconv.Itoa(len(consents))+" consents")
	for i := 0; i < len(consents); i++ {
		err = deleteConsent(stub, consents[i].ConsentID)
	}
	if err != nil {
		return shim.Error(err.Error())
		return shim.Error(buildError(errorRemoveConsent4App+appID))
	}
	return shim.Success(nil)
}

// =====================================================================================================================
// Get a Consents for an appID
// example:
// ./peer chaincode invoke -C mychanel -n consent -c '{"Args":["getconsents","APPID","ALLM"]}' -o 127.0.0.1:7050
// ./peer chaincode invoke -C mychanel -n consent -c '{"Args":["getconsents","APPID"]}' -o 127.0.0.1:7050
// =====================================================================================================================
func (c *ConsentCC)getConsents4AppID(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 && len(args) != 2 {
		errStr := errorArgs+" Expecting appID!"
		return shim.Error(buildError(errStr))
	}
	appID := args[0]
	logger.Debug("getConsents4AppID(Appid="+ appID+") : calling method -")
	consents, err := getConsentsByIndex(stub, indexApp,  []string{appID})
	if err != nil {
		return shim.Error(buildError(errorGetConsents4AppID+appID))
	}
	valAsBytes, err := json.Marshal(consents)
	if err != nil {
		return shim.Success([]byte("[]"))
	}
	return shim.Success(valAsBytes)
}


// =====================================================================================================================
// Get a Consents for an appID and a OwnerID
// example:
// ./peer chaincode invoke -C mch -n consent -c '{"Args":["getownerconsents","APPID","OWNERID"]}' -o 127.0.0.1:7050
// =====================================================================================================================
func (c *ConsentCC)getOwnerConsents(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		errStr := errorArgs+" Expecting appID, ownerID!"
		return shim.Error(buildError(errStr))
	}
	appID   := args[0]
	ownerID := args[1]
	logger.Debug("getOwnerConsents(Appid="+ appID+ ", OwnerID="+ownerID+") : calling method -")

	consents, err := getConsentsByIndex(stub, indexOwner,  []string{appID, ownerID})
	if err != nil {
		return shim.Error(buildError(errorGetConsents4Owner+ownerID+" appID:"+appID))
	}
	valAsBytes, err := json.Marshal(consents)
	if err != nil {
		return shim.Success([]byte("[]"))
	}
	return shim.Success(valAsBytes)
}

// =====================================================================================================================
// Get a Consent for an appID and a consumerID
// example:
// ./peer chaincode invoke -C mch -n consent -c '{"Args":["getconsumerconsents","APPID","CONSUMERID"]}'
// 						-o 127.0.0.1:7050
// =====================================================================================================================
func (c *ConsentCC)getConsumerConsents(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		errStr := errorArgs+" Expecting appID, consumerID!"
		return shim.Error(buildError(errStr))
	}
	appID      := args[0]
	consumerID := args[1]
	logger.Debug("getConsumerConsents(Appid="+ appID+ ", ConsumerID="+ consumerID+") : calling method -")
	consents, err := getConsentsByIndex(stub, indexConsumer,  []string{appID, consumerID})
	if err != nil {
		return shim.Error(buildError(errorGetConsents4Consumer+consumerID+" appID:"+appID))
	}
	valAsBytes, err := json.Marshal(consents)
	if err != nil {
		return shim.Success([]byte("[]"))
	}
	return shim.Success(valAsBytes)
}

// =====================================================================================================================
// Get a Consent for an appID,a consumerID and a ownerID
// example:
// ./peer chaincode invoke -C mch -n consent -c '{"Args":["getconsumerownerconsents","APPID","CONSUMERID", "OWNERID"]}'
// 						-o 127.0.0.1:7050
// =====================================================================================================================
func (c *ConsentCC)getConsumerOwnerConsents(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 3 {
		errStr := errorArgs+" Expecting appID, consumerID, ownerID!"
		return shim.Error(buildError(errStr))
	}
	appID      := args[0]
	consumerID := args[1]
	ownerID    := args[2]
	logger.Debug("getConsumerConsents(Appid="+ appID+ ", ConsumerID="+ consumerID+ ", OwnerID="+ ownerID+") : calling method -")
	fmt.Println(("getConsumerConsents(Appid="+ appID+ ", ConsumerID="+ consumerID+ ", OwnerID="+ ownerID+") : calling method -"))


	consents, err := getConsentsByIndex(stub, indexConsumerOwner,  []string{appID, consumerID, ownerID})
	if err != nil {
		return shim.Error(buildError(errorGetConsents4ConsumerOwner +"consumerID" + consumerID+" ownerID:"+ownerID+" appID:"+appID))
	}
	valAsBytes, err := json.Marshal(consents)
	if err != nil {
		return shim.Success([]byte("[]"))
	}
	return shim.Success(valAsBytes)
}

// =====================================================================================================================
// Verify if a consent exist
// example:
// ./peer chaincode invoke -C mch -n consent -c '{"Args":["isconsent","APPID", "OWNERID", "CONSUMERID", "DATATYPE",
// 							"ACCESSTYPE"]}' -o 127.0.0.1:7050
// =====================================================================================================================
func (c *ConsentCC)isConsent(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 5 {
		errStr := errorArgs+" Expecting AppID, OwnerID, CounsumerID, Datatype, Dataaccess!"
		return shim.Error(buildError(errStr))
	}
	logger.Debug("isConsent(Appid="+ args[0]+ ", Ownerid="+ args[1]+", Consumerid="+ args[2]+ ", Datatype="+ args[3]+
		", Dataaccess=" + args[4] +") : calling method -")

	appID := args[0]
	ownerID := args[1]
	consumerID := args[2]
	dataType := args[3]
	dataAccess := args[4]
	consents, err := getConsentsByIndex(stub, indexIsConsent, []string{appID, ownerID, consumerID,
		dataType, dataAccess})
	if err != nil {
		return shim.Error(buildError(errorGetConsent4Params+"appID:"+appID+" OwnerID:"+ownerID+" ConsumerID:"+
		consumerID+" dataType:"+dataType+" DataAccess:"+dataAccess))
	}
	for i := 0; i < len(consents); i++ {
		isValid := isValidToday(consents[i].Dt_begin, consents[i].Dt_end)
		if isValid {
			return shim.Success([]byte(AUTHORIZED))
		}
	}
	return shim.Success([]byte(NOT_AUTHORIZED))
}

// =====================================================================================================================
// createIndex - Create all index for a consent
// =====================================================================================================================
func createIndex(stub shim.ChaincodeStubInterface, consent consent) error {
	logger.Debug("createIndex(consentID="+ consent.ConsentID+") : calling method -")
	AppIndex, err := stub.CreateCompositeKey(indexApp, []string{consent.AppID, consent.ConsentID})
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	stub.PutState(AppIndex, []byte{0x00})

	OwnerIndex, err := stub.CreateCompositeKey(indexOwner, []string{consent.AppID, consent.OwnerID, consent.ConsentID})
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	stub.PutState(OwnerIndex, []byte{0x00})

	ConsumerIndex, err := stub.CreateCompositeKey(indexConsumer, []string{consent.AppID, consent.ConsumerID,
		consent.ConsentID})
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	stub.PutState(ConsumerIndex, []byte{0x00})

	ConsumerOwnerIndex, err := stub.CreateCompositeKey(indexConsumerOwner, []string{consent.AppID, consent.ConsumerID,
		consent.OwnerID, consent.ConsentID})
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	stub.PutState(ConsumerOwnerIndex, []byte{0x00})

	IsConsentIndex, err := stub.CreateCompositeKey(indexIsConsent, []string{consent.AppID, consent.OwnerID,
		consent.ConsumerID, consent.DataType, consent.DataAccess, consent.ConsentID})
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	stub.PutState(IsConsentIndex, []byte{0x00})
	return nil
}

// =====================================================================================================================
// deleteIndex - delete all index for a consent
// =====================================================================================================================
func deleteIndex(stub shim.ChaincodeStubInterface, consent consent) error {
	logger.Debug("deleteIndex(consentID="+ consent.ConsentID+") : calling method -")
	AppIndex, err := stub.CreateCompositeKey(indexApp, []string{consent.AppID, consent.ConsentID})
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	stub.DelState(AppIndex)

	OwnerIndex, err := stub.CreateCompositeKey(indexOwner, []string{consent.AppID, consent.OwnerID,
		consent.ConsentID})
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	stub.DelState(OwnerIndex)

	ConsumerIndex, err := stub.CreateCompositeKey(indexConsumer, []string{consent.AppID, consent.ConsumerID,
		consent.ConsentID})
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	stub.DelState(ConsumerIndex)

	IsConsentIndex, err := stub.CreateCompositeKey(indexIsConsent, []string{consent.AppID, consent.OwnerID,
		consent.ConsumerID, consent.DataType, consent.DataAccess, consent.ConsentID})
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	stub.DelState(IsConsentIndex)
	return nil
}

// =====================================================================================================================
// deleteConsent - delete a consent
// =====================================================================================================================
func deleteConsent(stub shim.ChaincodeStubInterface, consentID string) error {
	logger.Debug("deleteConsent(ConsentID="+ consentID+ ") : calling method -")
	valAsBytes, err := stub.GetState(consentID)
	if err != nil {
		errStr := "Failed to get consent:" + consentID
		logger.Error(errStr)
		return errors.New(errStr)
	} else if valAsBytes == nil {
		errStr := ("Consent does not exist:" + consentID )
		logger.Error(errStr)
		return errors.New(errStr)
	}
	var consent consent
	err = json.Unmarshal([]byte(valAsBytes), &consent)
	if err != nil {
		errStr := ("Failed to decode JSON for consent:" + consentID )
		logger.Error(errStr)
		return errors.New(errStr)
	}

	err = stub.DelState(consentID)
	if err != nil {
		errStr := ("Failed to delete consent:" + err.Error())
		logger.Error(errStr)
		return errors.New(errStr)
	}
	err = deleteIndex(stub, consent)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	return nil
}


// =====================================================================================================================
// use index to retrieve a list of consents
// =====================================================================================================================
func getConsentsByIndex(stub shim.ChaincodeStubInterface, index string, keys []string) ([]consent, error) {
	logger.Debug("getConsentsByIndex(index="+ index +") : calling method -")
	resultsIterator, err := stub.GetStateByPartialCompositeKey(index, keys)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	defer resultsIterator.Close()

	var i int
	consents := make([]consent, 0)
	for i = 0; resultsIterator.HasNext(); i++ {
		indexKey, err := resultsIterator.Next()
		if err != nil {
			logger.Error(err)
			return nil, err
		}
		_, compositeKeyParts, err := stub.SplitCompositeKey(indexKey.Key)
		if err != nil {
			logger.Error(err)
			return nil, err
		}
		consentID := compositeKeyParts[len(compositeKeyParts) - 1]
		consentAsBytes, err := stub.GetState(consentID)
		if err != nil {
			logger.Error(err)
			return nil, err
		}
		consent := consent{}
		err = json.Unmarshal(consentAsBytes, &consent)
		if err == nil {
			consents = append (consents, consent)
		}
	}
	return consents, nil
}

// =====================================================================================================================
// Build a json error to return
// =====================================================================================================================
func buildError(errorStr string) (jsonResp string){
	jsonResp = "{\"Error\":\""+errorStr+"\"}"
	logger.Error("return: "+ jsonResp)
	return
}

// =====================================================================================================================
// Check if the consent have a valid date for today
// =====================================================================================================================
func isValidToday(dt_start, dt_end time.Time) bool {
	logger.Debug("isValidToday(dt_start="+dt_start.String()+", dt_end="+dt_end.String()+") : calling method -")
	now := time.Now()
	dt_end = dt_end.Add(24 * time.Hour)
	isValid := now.After(dt_start) && now.Before(dt_end)
	return isValid
}

// =====================================================================================================================
// Convert date string to date: format (DD-MM-YYYY)
// =====================================================================================================================
func dateString2Date(datestr string) (time.Time, error) {
	logger.Debug("dateString2Date("+datestr+") : calling method -")
	layout := "2006-01-02"
	date, err := time.Parse(layout, datestr)
	if err != nil {
		logger.Error("dateString2Date error: ", err)
	}
	return date, err
}

// =====================================================================================================================
// Check if the period is valid (begin anterior to end)
// =====================================================================================================================
func checkDates(start, end string) ( dt_start , dt_end time.Time,  err error) {
	logger.Debug("CheckPeriod(start="+start+", end="+end+") : calling method -")

	dt_start, err = dateString2Date(start)
	if err != nil {
		errStr := errorDateBegin+ start
		logger.Error(errStr, err)
		err = errors.New(errStr)
		return
	}

	dt_end, err = dateString2Date(end)
	if err != nil {
		errStr := errorDateEnd+ end
		logger.Error(errStr, err)
		err = errors.New(errStr)
		return
	}
	dt_end = dt_end.Add(24 * time.Hour)

	if dt_start.Before(dt_end) {
		return
	} else {
		errStr := errorPeriod + dt_start.String() + " to:" + dt_end.String()
		logger.Error(errStr, err)
		err = errors.New(errStr)
		return
	}
}

// =====================================================================================================================
// Main
// =====================================================================================================================
func main() {
	err := shim.Start(new(ConsentCC))
	if err != nil {
		fmt.Printf("Error starting Smart contract Consent V2: %s", err)
	}
}


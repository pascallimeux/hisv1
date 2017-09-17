package test

import (
	"fmt"
	"testing"
	"github.com/pascallimeux/hisv1/blockchain"
	"os"
)

func TestChainCodeInvoke(t *testing.T) {

		// Add parameters for the initialization
	testSetup := blockchain.FabricSetup{
		// Channel parameters
		ChannelId:        "mychannel",
		ChannelConfig:    "../fixtures/channel/mychannel.tx",
		// Chaincode parameters
		ChaincodeId:      "heroes-service",
		ChaincodeVersion: "v1.0.0",
		ChaincodeGoPath:  os.Getenv("GOPATH"),
		ChaincodePath:    "github.com/hero_cc",
		// user parameters
		UserLogin:	  "admin",
		UserPwd:	  "adminpw",
		UserOrgID:	  "peerorg1",
		// his parameters
		ConfigFile:       "../config.yaml",
		StatStorePath:    "../keys/enroll_user",
	}

	err := blockchain.Initialize(&testSetup)
	if err != nil {
		t.Fatalf("Unable to initialize the Fabric SDK: %v\n", err)
	}

	testSetup.GetChaincodes()
/*
	// Install and instantiate the chaincode
	err = testSetup.InstallAndInstantiateCC()
	if err != nil {
		t.Fatalf("Unable to install and instantiate the chaincode: %v\n", err)
	}

	// Get Query value before invoke
	value, err := testSetup.QueryHello()
	if err != nil {
		t.Fatalf("getQueryHello return error: %v", err)
	}
	fmt.Printf("*** QueryValue before invoke %s\n", value)

	value, err = testSetup.InvokeHello("250")
	if err != nil {
		t.Fatalf("getInvokeHello return error: %v", err)
	}
	fmt.Printf("*** QueryValue after invoke %s\n", value)
*/
}

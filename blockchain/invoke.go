package blockchain

import (
	"fmt"
	"time"
	"github.com/hyperledger/fabric-sdk-go/api/apitxn"
)

// InvokeHello
func (setup *FabricSetup) InvokeHello(value string) (string, error) {

	// Prepare arguments
	var args []string
	args = append(args, "invoke")
	args = append(args, "hello")
	args = append(args, value)

	// Add data that will be visible in the proposal, like a description of the invoke request
	transientDataMap := make(map[string][]byte)
	transientDataMap["result"] = []byte("Transient data in hello invoke")

	// Make a next transaction proposal and send it
	transactionProposalResponse, txID, err := CreateAndSendTransactionProposal(
		setup.Channel,
		setup.ChaincodeId,
		"invoke",
		args,
		[]apitxn.ProposalProcessor{setup.Channel.PrimaryPeer()},
		transientDataMap,
	)
	if err != nil {
		return "", fmt.Errorf("Create and send transaction proposal in the invoke hello return error: %v", err)
	}

	// Register the Fabric SDK to listen to the event that will come back when the transaction will be send
	done, fail := RegisterTxEvent(txID, setup.EventHub)

	// Send the final transaction signed by endorser
	if _, err := CreateAndSendTransaction(setup.Channel, transactionProposalResponse); err != nil {
		return "", fmt.Errorf("Create and send transaction in the invoke hello return error: %v", err)
	}

	// Wait for the result of the submission
	select {
	// Transaction Ok
	case <-done:
		return txID.ID, nil

	// Transaction failed
	case <-fail:
		return "", fmt.Errorf("Error received from eventhub for txid(%s) error(%v)", txID, fail)

	// Transaction timeout
	case <-time.After(time.Second * 30):
		return "", fmt.Errorf("Didn't receive block event for txid(%s)", txID)
	}
}

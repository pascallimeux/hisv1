package blockchain

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/api/apitxn"
)

// QueryHello query the chaincode to get the state of hello
func (setup *FabricSetup) QueryHello() (string, error) {

	// Prepare arguments
	var args []string
	args = append(args, "query")
	args = append(args, "hello")

	// Make the proposal and submit it to the network (via our primary peer)
	transactionProposalResponses, _, err := CreateAndSendTransactionProposal(
		setup.Channel,
		setup.ChaincodeId,
		"invoke",
		args,
		[]apitxn.ProposalProcessor{setup.Channel.PrimaryPeer()}, // Peer contacted when submitted the proposal
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("Create and send transaction proposal return error in the query hello: %v", err)
	}
	return string(transactionProposalResponses[0].ProposalResponse.GetResponse().Payload), nil
}
package blockchain

import (
	"fmt"
	"github.com/hyperledger/fabric/common/cauthdsl"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabric-txn/admin"
	"github.com/hyperledger/fabric-sdk-go/api/apitxn"
	"github.com/hyperledger/fabric-sdk-go/def/fabapi/opt"
	fab "github.com/hyperledger/fabric-sdk-go/api/apifabclient"
	deffab "github.com/hyperledger/fabric-sdk-go/def/fabapi"
)

// FabricSetup implementation
type FabricSetup struct {
	Client           fab.FabricClient
	Channel          fab.Channel
	EventHub         fab.EventHub
	Initialized      bool
	ChannelId        string
	ChannelConfig    string
	ChaincodeId      string
	ChaincodeVersion string
	ChaincodeGoPath  string
	ChaincodePath    string
	UserLogin	 string
	UserPwd		 string
	UserOrgID	 string
	ConfigFile       string
	StatStorePath    string
}


// Initialize reads the configuration file and sets up the client, chain and event hub
func Initialize(setup *FabricSetup) error {

	sdkOptions := deffab.Options{
		ConfigFile: setup.ConfigFile,
		StateStoreOpts: opt.StateStoreOpts{
			Path: setup.StatStorePath,
		},
	}

	// This will make a user access (here the admin) to interact with the network
	// To do so, it will contact the Fabric CA to check if the user has access
	// and give it to him (enrollment)
	sc, err := GetClient(setup.UserLogin, setup.UserPwd, setup.UserOrgID, sdkOptions)
	if err != nil {
		return fmt.Errorf("Create client failed: %v", err)
	}
	setup.Client = sc

	// Make a new instance of channel pre-configured with the info we have provided,
	// but for now we can't use this channel because we need to create and
	// make some peer join it
	channel, err := GetChannel(setup.Client, setup.ChannelId, []string{setup.UserOrgID})
	if err != nil {
		return fmt.Errorf("Create channel (%s) failed: %v", setup.ChannelId, err)
	}
	setup.Channel = channel

	// Get an orderer user that will validate a proposed order
	// The authentication will be made with local certificates
	ordererUser, err := GetDefaultImplPreEnrolledUser(
		setup.Client,
		"ordererOrganizations/example.com/users/Admin@example.com/msp/keystore",
		"ordererOrganizations/example.com/users/Admin@example.com/msp/signcerts",
		"ordererAdmin",
		"peerorg1",
	)
	if err != nil {
		return fmt.Errorf("Unable to get the orderer user failed: %v", err)
	}

	// Get an organisation user (admin) that will be used to sign the proposal
	// The authentication will be made with local certificates
	admUser, err := GetDefaultImplPreEnrolledUser(
		setup.Client,
		"peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/keystore",
		"peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/signcerts",
		"peerorg1Admin",
		"peerorg1",
	)
	if err != nil {
		return fmt.Errorf("Unable to get the organisation user failed: %v", err)
	}

	// Initialize the channel "mychannel" based on the genesis block by
	//   1. locating in fixtures/channel/mychannel.tx and
	//   2. joining the peer given in the configuration file to this channel
	// Check if primary
	if err := CreateAndJoinChannel(setup.Client, ordererUser, admUser, setup.ChannelConfig, channel); err != nil {
		return fmt.Errorf("CreateAndJoinChannel return error: %v", err)
	}

	// Give the organisation user to the client for next proposal
	setup.Client.SetUserContext(admUser)

	// Setup Event Hub
	// This will allow us to listen for some event from the chaincode
	// and act on it. We won't use it for now.
	eventHub, err := GetEventHub(setup.Client, setup.UserOrgID)
	if err != nil {
		return err
	}
	if err := eventHub.Connect(); err != nil {
		return fmt.Errorf("Failed eventHub.Connect() [%s]", err)
	}
	setup.EventHub = eventHub

	// Tell that the initialization is done
	setup.Initialized = true

	return nil
}


func (setup *FabricSetup) InstallAndInstantiateCC() error {
	if setup.ChaincodeId == "" {
		setup.ChaincodeId = GenerateRandomID()
	}
	if err := setup.InstallCC(setup.ChaincodeId, setup.ChaincodePath, setup.ChaincodeVersion, nil); err != nil {
		return err
	}

	return setup.InstantiateCC(setup.ChaincodeId, setup.ChaincodePath, setup.ChaincodeVersion, []string{"init"})
}

func (setup *FabricSetup) InstantiateCC(chainCodeID string, chainCodePath string, chainCodeVersion string, args []string) error {
	chaincodePolicy := cauthdsl.SignedByMspMember(setup.Client.UserContext().MspID())

	if err := admin.SendInstantiateCC(setup.Channel, chainCodeID, args, chainCodePath, chainCodeVersion, chaincodePolicy, []apitxn.ProposalProcessor{setup.Channel.PrimaryPeer()}, setup.EventHub); err != nil {
		return err
	}
	return nil
}


func (setup *FabricSetup) InstallCC(chainCodeID string, chainCodePath string, chainCodeVersion string, chaincodePackage []byte) error {
	// installCC requires AdminUser privileges so setting user context with Admin User
	//setup.Client.SetUserContext(setup.AdminUser)

	// must reset client user context to normal user once done with Admin privilieges
	//defer setup.Client.SetUserContext(setup.NormalUser)

	if err := admin.SendInstallCC(setup.Client, chainCodeID, chainCodePath, chainCodeVersion, chaincodePackage, setup.Channel.Peers(), GetDeployPath()); err != nil {
		return fmt.Errorf("SendInstallProposal return error: %v", err)
	}

	return nil
}

func (setup *FabricSetup) GetChaincodes ()([]string, error){
	chaincodeQueryResponse, err := setup.Client.QueryInstalledChaincodes(setup.Channel.PrimaryPeer())
	return chaincodeQueryResponse.Chaincodes, err
}
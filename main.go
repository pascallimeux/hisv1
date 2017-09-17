package main

import (
	"fmt"
	"os"
	"runtime"
	"path/filepath"
	"github.com/pascallimeux/hisv1/blockchain"
	"github.com/pascallimeux/hisv1/web"
	"github.com/pascallimeux/hisv1/web/controllers"
)

// Fix empty GOPATH with golang 1.8 (see https://github.com/golang/go/blob/1363eeba6589fca217e155c829b2a7c00bc32a92/src/go/build/build.go#L260-L277)
func defaultGOPATH() string {
	env := "HOME"
	if runtime.GOOS == "windows" {
		env = "USERPROFILE"
	} else if runtime.GOOS == "plan9" {
		env = "home"
	}
	if home := os.Getenv(env); home != "" {
		def := filepath.Join(home, "go")
		if filepath.Clean(def) == filepath.Clean(runtime.GOROOT()) {
			// Don't set the default GOPATH to GOROOT,
			// as that will trigger warnings from the go tool.
			return ""
		}
		return def
	}
	return ""
}

func main() {
	//ENROLUSERREPO := "./keys/enroll_user"

	// Setup correctly the GOPATH in the environment
	if goPath := os.Getenv("GOPATH"); goPath == "" {
		os.Setenv("GOPATH", defaultGOPATH())
	}

	// Add parameters for the initialization
	fabricSdk := blockchain.FabricSetup{
		// Channel parameters
		ChannelId:        "mychannel",
		ChannelConfig:    "fixtures/channel/mychannel.tx",
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
		ConfigFile:       "config.yaml",
		StatStorePath:    "./keys/enroll_user",
	}


	// Initialize the Fabric SDK
	err := blockchain.Initialize(&fabricSdk)
	if err != nil {
		fmt.Printf("Unable to initialize the Fabric SDK: %v\n", err)
	}

	// Install and instantiate the chaincode
	err = fabricSdk.InstallAndInstantiateCC()
	if err != nil {
		fmt.Printf("Unable to install and instantiate the chaincode: %v\n", err)
	}

	// Make the web application listening
	app := &controllers.Application{
		Fabric: &fabricSdk,
	}
	web.Serve(app)
}

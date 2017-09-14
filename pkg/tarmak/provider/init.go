package provider

import (
	"fmt"

	"github.com/jetstack/tarmak/pkg/tarmak/config"
	"github.com/jetstack/tarmak/pkg/tarmak/initialize"
	"github.com/jetstack/tarmak/pkg/tarmak/interfaces"
	"github.com/jetstack/tarmak/pkg/tarmak/utils"
)

func Init(t interfaces.Tarmak) {
	conf, err := config.New(t)
	if err != nil {
		t.Log().Fatal(err)
	}

	// If no providers, need to init (not nil)
	if conf.Providers() == nil {
		conf.InitConfig()
	}

	cloudProvider, credentialsSource, credentials := initialize.AskProvider()
	providerName := initialize.AskProviderName(conf)
	bucketPrefix := initialize.AskBucketPrefix(cloudProvider, conf)
	publicZone := initialize.AskPublicZone(cloudProvider, conf)

	/* This will be generated from the s3 bucket prefix right now. Not too sure but would like to keep it like that. Maybe we call the bucket_prefix resource_prefix for provider wide resources */
	//open = &utils.Open{
	//	Query:    "What is the dynamo DB lock table name?",
	//	Prompt:   "> ",
	//	Required: true,
	//}
	//dynamoDbLockName := open.Ask()

	fmt.Printf("\nCloud Provider ---->%s\n", cloudProvider)
	fmt.Printf("Provider Name ----->%s\n", providerName)
	fmt.Printf("Credentials Source >%s\n", credentialsSource)
	if credentialsSource == "AWS folder" {
		fmt.Printf("Profile Name ------>%s\n", credentials)
	} else {
		fmt.Printf("Vault Prefix ------>%s\n", credentials)
	}
	fmt.Printf("Public Zone ------->%s\n", publicZone)
	fmt.Printf("Resource Prefix --->%s\n", bucketPrefix)

	yn := &utils.YesNo{
		Query:   "Are these input correct?",
		Prompt:  "> ",
		Default: true,
	}
	if yn.Ask() && cloudProvider == "AWS" {
		prov := config.NewAWSProfileProvider(providerName, credentials)
		prov.AWS.PublicZone = publicZone
		prov.AWS.BucketPrefix = bucketPrefix
		conf.AppendProvider(prov)

		fmt.Print("Accepted.\n")
	} else {
		fmt.Print("Aborting.\n")
	}
}

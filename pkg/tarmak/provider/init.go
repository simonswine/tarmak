package provider

import (
	"fmt"

	"github.com/jetstack/tarmak/pkg/tarmak/config"
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
		conf.EmptyConfig()
	}

	sel := &utils.Select{
		Query:   "Select a provider",
		Prompt:  "> ",
		Choice:  &[]string{"AWS"},
		Default: 1,
	}
	cloudProvider := sel.Ask()

	var name string
	open := &utils.Open{
		Query:    "Enter a unique name for this provider",
		Prompt:   "> ",
		Required: true,
	}
	for name == "" {
		resp := open.Ask()
		if err := conf.ValidName(resp, "[a-z0-9-]+"); err != nil {
			fmt.Printf("Name is not valid: %v", err)
		} else if err := conf.UniqueProviderName(resp); err != nil {
			fmt.Printf("Name is not valid: %v", err)
		} else {
			name = resp
		}
	}

	sel = &utils.Select{
		Query:   "Where should the credentials come from?",
		Prompt:  "> ",
		Choice:  &[]string{"AWS folder", "Vault"},
		Default: 1,
	}
	credentialsSource := sel.Ask()

	var profileName string
	var vaultPrefix string
	if credentialsSource == "AWS folder" {
		open := &utils.Open{
			Query:    "What is the profile name?",
			Prompt:   "> ",
			Required: true,
		}
		profileName = open.Ask()

	} else {
		open := &utils.Open{
			Query:   "Which path should be used for AWS credentials?",
			Prompt:  "> ",
			Default: "jetstack/aws/jetstack-dev/sts/admin",
		}
		vaultPrefix = open.Ask()
	}

	query := "Whats is the s3 bucket prefix?"
	if cloudProvider != "AWS" {
		query = "What is the resource prefix?"
	}
	open = &utils.Open{
		Query:    query,
		Prompt:   "> ",
		Required: true,
	}
	resourcePrefix := open.Ask()

	// TODO: Validate if the zone is delegated in through the root servers
	var publicZone string
	if cloudProvider == "AWS" {
		zones := make(map[string]bool)
		var choice []string

		// create map of all zones used in all providers
		for _, p := range conf.Providers() {
			if p.AWS != nil {
				zones[p.AWS.PublicZone] = true
			}
		}

		// put list in []string for input select
		for zone := range zones {
			choice = append(choice, zone)
		}
		choice = append(choice, "enter custom zone")

		sel = &utils.Select{
			Query:   "Select public zone",
			Prompt:  "> ",
			Choice:  &choice,
			Default: 1,
		}
		publicZone = sel.Ask()
	}

	if cloudProvider != "AWS" || publicZone == "enter custom zone" {
		open = &utils.Open{
			Query:    "What public zone should be used?\nPlease make sure you can delegate this zone.",
			Prompt:   "> ",
			Required: true,
		}
		publicZone = open.Ask()
	}

	/* This will be generated from the s3 bucket prefix right now. Not too sure but would like to keep it like that. Maybe we call the bucket_prefix resource_prefix for provider wide resources */
	//open = &utils.Open{
	//	Query:    "What is the dynamo DB lock table name?",
	//	Prompt:   "> ",
	//	Required: true,
	//}
	//dynamoDbLockName := open.Ask()

	fmt.Printf("\nCloud Provider ---->%s\n", cloudProvider)
	fmt.Printf("Provider Name ----->%s\n", name)
	fmt.Printf("Credentials Source >%s\n", credentialsSource)
	if credentialsSource == "AWS folder" {
		fmt.Printf("Profile Name ------>%s\n", profileName)
	} else {
		fmt.Printf("Vault Prefix ------>%s\n", vaultPrefix)
	}
	fmt.Printf("Public Zone ------->%s\n", publicZone)
	fmt.Printf("Resource Prefix --->%s\n", resourcePrefix)

	yn := &utils.YesNo{
		Query:   "Are these input correct?",
		Prompt:  "> ",
		Default: true,
	}
	if yn.Ask() && cloudProvider == "AWS" {
		prov := config.NewAWSProfileProvider(name, profileName)
		prov.AWS.PublicZone = publicZone
		prov.AWS.BucketPrefix = resourcePrefix
		conf.AppendProvider(prov)

		fmt.Print("Accepted.\n")
	} else {
		fmt.Print("Aborting.\n")
	}
}

package provider

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jetstack/tarmak/pkg/tarmak/config"
	"github.com/jetstack/tarmak/pkg/tarmak/utils"
)

func Init(cmd *cobra.Command) {

	sel := &utils.Select{
		Query:   "Select a provider",
		Prompt:  "> ",
		Choice:  &[]string{"AWS"},
		Default: 1,
	}
	cloudProvider := sel.Ask()

	sel = &utils.Select{
		Query:   "Where should the credentials come from?",
		Prompt:  "> ",
		Choice:  &[]string{"AWS folder", "Vault Path"},
		Default: 1,
	}
	credentialsSource := sel.Ask()

	sel = &utils.Select{
		Query:  "Which public zone should be used for DNS records?",
		Prompt: "> ",
		Choice: &[]string{"us-east-1", "us-east-2", "us-west-1", "us-west-2", "ca-central-1", "eu-west-1", "eu-central-1", "eu-west-2", "ap-northeast-1", "ap-northeast-2", "ap-southeast-1", "ap-southeast-2", "ap-south-1", "sa-east-1", "enter custom zone"},
	}
	publicZone := sel.Ask()

	if publicZone == "enter custom zone" {
		open := &utils.Open{
			Query:    "Enter custom zone",
			Prompt:   "> ",
			Required: true,
		}
		publicZone = open.Ask()
	}

	open := &utils.Open{
		Query:    "What should be the s3 bucket prefix",
		Prompt:   "> ",
		Required: true,
	}
	bucketPrefix := open.Ask()

	open = &utils.Open{
		Query:    "What should be the dynamo db lock table name",
		Prompt:   "> ",
		Required: true,
	}
	dynamoDbLockName := open.Ask()

	fmt.Printf("cloudProvider >%s\n", cloudProvider)
	fmt.Printf("credentialsSource >%s\n", credentialsSource)
	fmt.Printf("publicZone >%s\n", publicZone)
	fmt.Printf("bucketPrefix >%s\n", bucketPrefix)
	fmt.Printf("dynamoDbLockName >%s\n", dynamoDbLockName)

	prov := config.NewAWSProfileProvider("name", "profile")
	prov.AWS.PublicZone = publicZone
	prov.AWS.BucketPrefix = bucketPrefix
}

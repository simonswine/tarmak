package initialize

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Sirupsen/logrus"

	"github.com/jetstack/tarmak/pkg/tarmak/config"
	"github.com/jetstack/tarmak/pkg/tarmak/interfaces"
	"github.com/jetstack/tarmak/pkg/tarmak/utils"
)

type Init struct {
	log    *logrus.Entry
	tarmak interfaces.Tarmak
}

func New(t interfaces.Tarmak) *Init {
	return &Init{
		log:    t.Log(),
		tarmak: t,
	}
}

func parseContextName(in string) (environment string, context string, err error) {
	in = strings.ToLower(in)

	splitted := false

	for i, c := range in {
		if !splitted && c == '-' {
			splitted = true
			environment = in[0:i]
			context = in[i+1 : len(in)]
		} else if c < '0' || (c > '9' && c < 'a') || c > 'z' {
			return "", "", fmt.Errorf("invalid char '%c' in string '%s' at pos %d", c, in, i)
		}
	}

	if !splitted {
		return "", "", fmt.Errorf("string '%s' did not contain '-'", in)
	}
	return environment, context, nil
}

func (i *Init) Run() (err error) {

	conf := i.tarmak.Config()
	if conf.ConfigIsEmpty() {
		conf, err = config.New(i.tarmak)
		if err != nil {
			i.tarmak.Log().Fatal(err)
		}
		conf.InitConfig()
	}

	/* TODO: support multiple cluster in one env
	query = "What kind of cluster do you want to initialise?"
	options = []string{"create new single cluster environment", "create new multi cluster environment", "add new cluster to existing multi cluster environment"}
	kind, err := ui.Select(query, options, &input.Options{
		Default: options[0],
		Loop:    true,
		ValidateFunc: func(s string) error {
			if s != "1" {
				return fmt.Errorf(`option "%s" is currently not supported`, s)
			}
			return nil
		},
	})
	if err != nil {
		return err
	}
	*/

	combinedName, environment, context := AskCombindName(conf)
	cloudProvider, credentialsSource, credentials := AskProvider()
	//awsRegion := AskRegion()
	contact, err := AskEmail(conf)
	if err != nil {
		fmt.Print(err)
		return nil
	}
	bucketPrefix := AskBucketPrefix(cloudProvider, conf)
	publicZone := AskPublicZone(cloudProvider, conf)
	privateZone := AskPrivateZone()
	projectName := AskProjectName()

	fmt.Printf("\nCombined Name ->%s\n", combinedName)
	fmt.Printf("Environment --->%s\n", environment)
	fmt.Printf("Context ------->%s\n", context)
	fmt.Printf("Cloud Provider >%s\n", cloudProvider)
	if credentialsSource == "AWS folder" {
		fmt.Printf("Profile Name -->%s\n", credentials)
	} else {
		fmt.Printf("Vault Prefix -->%s\n", credentials)
	}
	fmt.Printf("Contact ------->%s\n", contact)
	fmt.Printf("Bucket Prefix ->%s\n", bucketPrefix)
	fmt.Printf("Public Zone --->%s\n", publicZone)
	fmt.Printf("Private Zone -->%s\n", privateZone)
	fmt.Printf("Project Name -->%s\n", projectName)

	yn := &utils.YesNo{
		Query:   "Are these input correct?",
		Prompt:  "> ",
		Default: true,
	}
	if yn.Ask() {
		fmt.Print("Accepted.\n")
	} else {
		fmt.Print("Aborted.\n")
	}

	//env := config.Environment{
	//	Contact: contact,
	//	Project: project,
	//	AWS: &config.AWSConfig{
	//		VaultPath: vaultPrefix,
	//		Region:    awsRegion,
	//	},
	//	Name: environmentName,
	//	Contexts: []config.Context{
	//		config.Context{
	//			Name:      contextName,
	//			BaseImage: "centos-puppet-agent",
	//			Stacks: []config.Stack{
	//				config.Stack{
	//					State: &config.StackState{
	//						BucketPrefix: bucketPrefix,
	//						PublicZone:   publicZone,
	//					},
	//				},
	//				config.Stack{
	//					Network: &config.StackNetwork{
	//						NetworkCIDR: "10.98.0.0/20",
	//						PrivateZone: privateZone,
	//					},
	//				},
	//				config.Stack{
	//					Tools: &config.StackTools{},
	//				},
	//				config.Stack{
	//					Vault: &config.StackVault{},
	//				},
	//				config.Stack{
	//					Kubernetes: &config.StackKubernetes{},
	//					NodeGroups: config.DefaultKubernetesNodeGroupAWSOneMasterThreeEtcdThreeWorker(),
	//				},
	//			},
	//		},
	//	},
	//}

	//return i.tarmak.MergeEnvironment(env)
	return nil
}

func AskCombindName(conf interfaces.Config) (combinedName, environment, context string) {
	var err error
	open := &utils.Open{
		Query:    "What should be the name of the cluster?\nThe name consists of two parts seperated by a dash. First part is the environment name, second part the cluster name. Both names should be matching [a-z0-9]+",
		Prompt:   "> ",
		Required: true,
	}

	for environment == "" {
		combinedName = open.Ask()
		environment, context, err = parseContextName(combinedName)
		if err != nil {
			fmt.Print(err)
			open.Query = ""
		} else if err := conf.UniqueEnvironmentName(environment); err != nil {
			fmt.Printf("Invalid environment name: %v", err)
			environment = ""
		}
	}
	// TODO ensure max length of both is not longer than 24 chars (verify that limit from AWS)

	return combinedName, environment, context
}

func AskProvider() (provider, credentialsSource, credentials string) {
	sel := &utils.Select{
		Query:   "Select a provider",
		Prompt:  "> ",
		Choice:  &[]string{"AWS"},
		Default: 1,
	}
	provider = sel.Ask()

	sel = &utils.Select{
		Query:   "Where should the credentials come from?",
		Prompt:  "> ",
		Choice:  &[]string{"AWS folder", "Vault Path"},
		Default: 1,
	}
	credentialsSource = sel.Ask()

	if credentialsSource == "AWS folder" {
		open := &utils.Open{
			Query:    "What is the profile name?",
			Prompt:   "> ",
			Required: true,
		}
		credentials = open.Ask()

	} else {
		open := &utils.Open{
			Query:   "Which path should be used for AWS credentials?",
			Prompt:  "> ",
			Default: "jetstack/aws/jetstack-dev/sts/admin",
		}
		credentials = open.Ask()
	}

	return provider, credentialsSource, credentials
}

//func AskRegion() (awsRegion string) {
//
//	sel := &utils.Select{
//		Query:  "Which region should be used for DNS records?",
//		Prompt: "> ",
//		Choice: &[]string{"us-east-1", "us-east-2", "us-west-1", "us-west-2", "ca-central-1", "eu-west-1", "eu-central-1", "eu-west-2", "ap-northeast-1", "ap-northeast-2", "ap-southeast-1", "ap-southeast-2", "ap-south-1", "sa-east-1", "enter custom zone"},
//	}
//	awsRegion = sel.Ask()
//
//	if awsRegion == "enter custom zone" {
//		open := &utils.Open{
//			Query:    "Enter custom zone",
//			Prompt:   "> ",
//			Required: true,
//		}
//		awsRegion = open.Ask()
//	}
//	// TODO: validate region
//
//	return awsRegion
//}

func AskEmail(conf interfaces.Config) (contact string, err error) {
	open := &utils.Open{
		Query:    "What is the mail address of someone responsible?",
		Prompt:   "> ",
		Default:  conf.Contact(),
		Required: true,
	}
	contact = open.Ask()

	var fail error
	if err = utils.ValidateFormat(contact); err != nil {
		fail = err

		// Not sure if this is a good idea bc of privacy concerns
	} else if err = utils.ValidateHost(contact); err != nil {
		fail = err
	}
	if fail != nil {
		yn := &utils.YesNo{
			Query:   fmt.Sprintf("Error verifying email, did you spell it correctly?:\n%v\nUse anyway?", fail),
			Prompt:  "> ",
			Default: true,
		}
		if !yn.Ask() {
			return "", errors.New("Aborting")
		}
	}

	return contact, nil
}

func AskBucketPrefix(cloudProvider string, conf interfaces.Config) (bucketPrefix string) {
	query := "Whats is the s3 bucket prefix?"
	if cloudProvider != "AWS" {
		query = "What is the resource prefix?"
	}

	open := &utils.Open{
		Query:   query,
		Prompt:  "> ",
		Default: "tarmak-",
	}
	for bucketPrefix == "" {
		bucketPrefix = open.Ask()
		if err := conf.ValidName(bucketPrefix, "[a-z0-9-_]+"); err != nil {
			fmt.Printf("Name is not valid: %v", err)
			open.Query = ""
			bucketPrefix = ""
		}
	}

	return bucketPrefix
}

func AskPublicZone(cloudProvider string, conf interfaces.Config) (publicZone string) {

	// TODO: Validate if the zone is delegated in through the root servers

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

		sel := &utils.Select{
			Query:   "Select public zone",
			Prompt:  "> ",
			Choice:  &choice,
			Default: 1,
		}
		publicZone = sel.Ask()
	}

	if cloudProvider != "AWS" || publicZone == "enter custom zone" {
		open := &utils.Open{
			Query:    "What public zone should be used?\nPlease make sure you can delegate this zone.",
			Prompt:   "> ",
			Required: true,
		}
		publicZone = open.Ask()
	}

	return publicZone
}

func AskPrivateZone() (privateZone string) {

	// TODO: verify domain name

	open := &utils.Open{
		Query:   "What private zone should be used?",
		Prompt:  "> ",
		Default: "tarmak.local",
	}
	return open.Ask()
}

func AskProjectName() (projectName string) {

	open := &utils.Open{
		Query:   "What is the project name?",
		Prompt:  "> ",
		Default: "k8s-playground",
	}
	return open.Ask()
}

func AskProviderName(conf interfaces.Config) (providerName string) {
	open := &utils.Open{
		Query:    "Enter a unique name for this provider",
		Prompt:   "> ",
		Required: true,
	}

	for providerName == "" {
		resp := open.Ask()
		if err := conf.ValidName(resp, "[a-z0-9-]+"); err != nil {
			fmt.Printf("Name is not valid: %v", err)
		} else if err := conf.UniqueProviderName(resp); err != nil {
			fmt.Printf("Name is not valid: %v", err)
		} else {
			providerName = resp
		}
	}
	return providerName
}

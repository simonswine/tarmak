package initialize

import (
	"fmt"
	//"os"
	"strings"

	"github.com/Sirupsen/logrus"
	//"github.com/tcnksm/go-input"

	//"github.com/jetstack/tarmak/pkg/tarmak/config"
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

func (i *Init) Run() error {

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

	query := "What should be the name of the cluster?\nThe name consists of two parts seperated by a dash. First part is the environment name, second part the cluster name. Both names should be matching [a-z0-9]+"
	open := &utils.Open{
		Query:    query,
		Prompt:   "> ",
		Required: true,
	}

	combinedName := open.Ask()
	environment, context, err := parseContextName(combinedName)
	if err != nil {
		return err
	}
	// TODO verify environment name not taken yet
	// TODO ensure max length of both is not longer than 24 chars (verify that limit from AWS)

	query = "Select a provider"
	sel := &utils.Select{
		Query:   "Select a provider",
		Prompt:  "> ",
		Choice:  &[]string{"AWS"},
		Default: 1,
	}
	provider := sel.Ask()

	sel = &utils.Select{
		Query:   "Where should the credentials come from?",
		Prompt:  "> ",
		Choice:  &[]string{"AWS folder", "Vault Path"},
		Default: 1,
	}
	credentialsSource := sel.Ask()

	var vaultPrefix string
	if credentialsSource == "AWS folder" {
		open = &utils.Open{
			Query:   "Which path should be used for AWS credentials?",
			Prompt:  "> ",
			Default: "jetstack/aws/jetstack-dev/sts/admin",
		}
		vaultPrefix = open.Ask()
	}

	sel = &utils.Select{
		Query:  "Which region should be used for DNS records?",
		Prompt: "> ",
		Choice: &[]string{"us-east-1", "us-east-2", "us-west-1", "us-west-2", "ca-central-1", "eu-west-1", "eu-central-1", "eu-west-2", "ap-northeast-1", "ap-northeast-2", "ap-southeast-1", "ap-southeast-2", "ap-south-1", "sa-east-1", "enter custom zone"},
	}
	awsRegion := sel.Ask()

	if awsRegion == "enter custom zone" {
		open := &utils.Open{
			Query:    "Enter custom zone",
			Prompt:   "> ",
			Required: true,
		}
		awsRegion = open.Ask()
	}
	// TODO: validate region

	open = &utils.Open{
		Query:    "What is the s3 bucket prefix?",
		Prompt:   "> ",
		Default:  "tarmak-",
		Required: true,
	}
	bucketPrefix := open.Ask()
	// TODO: verify bucket prefix [a-z0-9-_]

	query = "What public zone should be used?\nPlease make sure you can delegate this zone to AWS!"
	open = &utils.Open{
		Query:    query,
		Prompt:   "> ",
		Required: true,
	}
	publicZone := open.Ask()
	// TODO: verify domain name

	open = &utils.Open{
		Query:   "What private zone should be used?",
		Prompt:  "> ",
		Default: "tarmak.local",
	}
	privateZone := open.Ask()
	// TODO: verify domain name

	open = &utils.Open{
		Query:    "What is the mail address of someone responsible?",
		Prompt:   "> ",
		Required: true,
	}
	contact := open.Ask()
	// TODO: use default from existing config
	// TODO: verify mail

	open = &utils.Open{
		Query:   "What is the project name?",
		Prompt:  "> ",
		Default: "k8s-playground",
	}
	projectName := open.Ask()

	fmt.Printf("\nEnvironment--->%s\n", environment)
	fmt.Printf("Context------->%s\n", context)
	fmt.Printf("Cloud Provider>%s\n", provider)
	fmt.Printf("Vault Prefix-->%s\n", vaultPrefix)
	fmt.Printf("Bucket Prefix->%s\n", bucketPrefix)
	fmt.Printf("Public Zone--->%s\n", publicZone)
	fmt.Printf("Private Zone-->%s\n", privateZone)
	fmt.Printf("Contact------->%s\n", contact)
	fmt.Printf("Project Name-->%s\n", projectName)

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

package tarmak

import (
	"fmt"
	"strings"

	"github.com/jetstack/tarmak/pkg/tarmak/config"
	"github.com/jetstack/tarmak/pkg/tarmak/interfaces"
)

const FlagTerraformStacks = "terraform-stacks"
const FlagForceDestroyStateStack = "force-destroy-state"

func (t *Tarmak) Terraform() interfaces.Terraform {
	return t.terraform
}

func (t *Tarmak) stackList() ([]interfaces.Stack, error) {
	var zeroStackList []interfaces.Stack
	allStacks := t.Context().Stacks()

	selectedStackNames, err := t.cmd.Flags().GetStringSlice(FlagTerraformStacks)
	if err != nil {
		return zeroStackList, fmt.Errorf(
			"could not find flag %s: %s", FlagTerraformStacks, err,
		)
	}
	if len(selectedStackNames) == 0 {
		return allStacks, nil
	}

	selectedStackMap := map[string]bool{}
	for _, stack := range allStacks {
		selectedStackMap[stack.Name()] = false
	}
	unrecognisedStackNames := []string{}
	for _, stackName := range selectedStackNames {
		_, stackRecognised := selectedStackMap[stackName]
		if stackRecognised {
			selectedStackMap[stackName] = true
		} else {
			unrecognisedStackNames = append(unrecognisedStackNames, stackName)
		}
	}
	if len(unrecognisedStackNames) > 0 {
		return zeroStackList, fmt.Errorf(
			"unrecognised --%s: %v",
			FlagTerraformStacks,
			unrecognisedStackNames,
		)
	}
	stackList := []interfaces.Stack{}
	for _, stack := range allStacks {
		stackSelected := selectedStackMap[stack.Name()]
		if stackSelected {
			stackList = append(stackList, stack)
		}
	}

	return stackList, nil
}

func (t *Tarmak) CmdTerraformApply(args []string) error {
	t.discoverAMIID()
	stacks, err := t.stackList()
	if err != nil {
		return err
	}
	for _, stack := range stacks {
		t.log.WithField("stack", stack.Name()).Debug("applying stack")
		err := t.terraform.Apply(stack, args)
		if err != nil {
			t.log.Fatal(err)
		}
	}
	return nil
}

func (t *Tarmak) CmdTerraformDestroy(args []string) error {
	forceDestroyStateStack, err := t.cmd.Flags().GetBool(FlagForceDestroyStateStack)
	if err != nil {
		return fmt.Errorf("could not find flag %s: %s", FlagForceDestroyStateStack, err)
	}
	t.discoverAMIID()
	stacks, err := t.stackList()
	if err != nil {
		return err
	}
	for posStack, _ := range stacks {
		stack := stacks[len(stacks)-posStack-1]
		if !forceDestroyStateStack && stack.Name() == config.StackNameState {
			t.log.Debugf("ignoring stack '%s'", stack.Name())
			continue
		}
		err := t.terraform.Destroy(stack, args)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *Tarmak) CmdTerraformShell(args []string) error {
	paramStackName := ""
	if len(args) > 0 {
		paramStackName = strings.ToLower(args[0])
	}

	// find matching stacks
	stacks := t.Context().Stacks()
	stackNames := make([]string, len(stacks))
	for pos, stack := range stacks {
		stackNames[pos] = stack.Name()
		if stack.Name() == paramStackName {
			// prepare stack's shell
			t.discoverAMIID()
			return t.terraform.Shell(stack, args)
		}
	}

	return fmt.Errorf("you have to provide exactly one parameter that contains one of the stack names %s", strings.Join(stackNames, ", "))
}

package utils

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

type Open struct {
	Query    string
	Prompt   string
	Required bool
	Default  string
}

type Select struct {
	Query   string
	Prompt  string
	Choice  *[]string
	Default int
}

type YesNo struct {
	Query   string
	Prompt  string
	Default bool
}

func (y *YesNo) Ask() (responce bool) {
	go Catch()

	reader := bufio.NewReader(os.Stdin)
	var option string
	if y.Default {
		option = " [Y/n]"
	} else {
		option = " [y/N]"
	}

	fmt.Printf("\n%s%s", y.Query, option)

	for {
		fmt.Printf("\n%s", y.Prompt)

		res, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		res = res[0 : len(res)-1]

		res = strings.ToLower(res)
		if res == "y" || res == "yes" {
			return true
		} else if res == "n" || res == "no" {
			return false
		} else if res == "" {
			return y.Default
		} else {
			fmt.Printf("Bad responce. %s", option)
		}
	}
}

func (s *Select) Ask() (responce string) {
	go Catch()

	reader := bufio.NewReader(os.Stdin)
	if s.Query != "" {
		fmt.Printf("\n%s", s.Query)
	}
	if s.Default > 0 {
		fmt.Printf(" (default %d)", s.Default)
	}
	fmt.Print("\n")

	for i, s := range *s.Choice {
		fmt.Printf("%d. %s\n", i+1, s)
	}

	var n int
	for {
		fmt.Print(s.Prompt)

		res, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		res = res[0 : len(res)-1]

		if res == "" {
			if s.Default > 0 {
				//TODO: Make this better
				choices := *s.Choice
				return choices[s.Default-1]
			} else {
				fmt.Printf("Nothing entered.\n")
			}

		} else if n, err = strconv.Atoi(res); err != nil || n < 1 || n > len(*s.Choice) {
			fmt.Printf("Responce must be a number between 1 and %d\n", len(*s.Choice))
		} else {
			//TODO: improve this:
			choices := *s.Choice
			return choices[n-1]
		}
	}
}

func (o *Open) Ask() (responce string) {
	go Catch()

	reader := bufio.NewReader(os.Stdin)
	if o.Query != "" {
		fmt.Printf("\n%s", o.Query)
	}
	if !o.Required {
		fmt.Printf(" (default %s)", o.Default)
	}
	fmt.Print("\n")

	for {
		if o.Prompt != "" {
			fmt.Print(o.Prompt)
		}

		res, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		res = res[0 : len(res)-1]

		if res == "" {
			if o.Required {
				fmt.Print("Nothing entered\n")
			} else {
				return o.Default
			}

		} else {
			return res
		}
	}
}

func Catch() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	fmt.Printf("\nHandle exit\n")

	os.Exit(1)
}

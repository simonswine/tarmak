package utils

import (
	"errors"
	"fmt"
	"net"
	"net/smtp"
	"regexp"
	"strings"
)

func ValidateFormat(email string) error {
	regex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	if !regex.MatchString(email) {
		return errors.New("invalid format")
	}

	return nil
}

func ValidateHost(email string) error {
	host := getHost(email)
	mx, err := net.LookupMX(host)
	if err != nil {
		return errors.New("unresolvable host")
	}

	client, err := smtp.Dial(fmt.Sprintf("%s:%d", mx[0].Host, 25))
	if err != nil {
		return err
	}
	defer client.Close()

	if err = client.Hello("google.com"); err != nil {
		return err
	}
	if err = client.Mail("noreply@jetstack.io"); err != nil {
		return err
	}
	if err := client.Rcpt(email); err != nil {
		return err
	}

	return nil
}

func getHost(email string) string {
	split := strings.Split(email, "@")
	return split[1]
}

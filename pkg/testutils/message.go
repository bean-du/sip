package testutils

import (
	"fmt"
	"strings"

	"github.com/bean-du/sip/pkg/log"
	"github.com/bean-du/sip/pkg/sip"
	"github.com/bean-du/sip/pkg/sip/parser"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func Message(rawMsg []string) sip.Message {
	msg, err := parser.ParseMessage([]byte(strings.Join(rawMsg, "\r\n")), log.NewDefaultLogrusLogger())
	Expect(err).ToNot(HaveOccurred())
	return msg
}

func Request(rawMsg []string) sip.Request {
	msg := Message(rawMsg)
	switch msg := msg.(type) {
	case sip.Request:
		return msg
	case sip.Response:
		Fail(fmt.Sprintf("%s is not a request", msg.Short()))
	default:
		Fail(fmt.Sprintf("%s is not a request", msg.Short()))
	}
	return nil
}

func Response(rawMsg []string) sip.Response {
	msg := Message(rawMsg)
	switch msg := msg.(type) {
	case sip.Response:
		return msg
	case sip.Request:
		Fail(fmt.Sprintf("%s is not a response", msg.Short()))
	default:
		Fail(fmt.Sprintf("%s is not a response", msg.Short()))
	}
	return nil
}

package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bean-du/sip/pkg/account"
	"github.com/bean-du/sip/pkg/log"
	"github.com/bean-du/sip/pkg/media/rtp"
	"github.com/bean-du/sip/pkg/sip/parser"
	"github.com/bean-du/sip/pkg/stack"
	"github.com/bean-du/sip/pkg/ua"
	"github.com/bean-du/sip/pkg/utils"
)

var (
	logger log.Logger
	udp    *rtp.RtpUDPStream
)

func init() {
	logger = utils.NewLogrusLogger(log.DebugLevel, "Register", nil)
}

func main() {

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	stack := stack.NewSipStack(&stack.SipStackConfig{
		UserAgent:  "Go Sip Client/example-register",
		Extensions: []string{"replaces", "outbound"},
		Dns:        "8.8.8.8"})

	if err := stack.Listen("udp", "0.0.0.0:5066"); err != nil {
		logger.Panic(err)
	}

	ua := ua.NewUserAgent(&ua.UserAgentConfig{
		SipStack: stack,
	})

	ua.RegisterStateHandler = func(state account.RegisterState) {
		logger.Infof("RegisterStateHandler: user => %s, state => %v, expires => %v, reason => %v", state.Account.AuthInfo.AuthUser, state.StatusCode, state.Expiration, state.Reason)
	}

	uri, err := parser.ParseUri("sip:100@127.0.0.1") // this acts as an identifier, not connection info
	if err != nil {
		logger.Error(err)
	}

	profile := account.NewProfile(uri.Clone(), "goSIP",
		&account.AuthInfo{
			AuthUser: "100",
			Password: "100",
			Realm:    "b2bua",
		},
		1800,
		stack,
	)

	recipient, err := parser.ParseSipUri("sip:100@127.0.0.1;transport=udp") // this is the remote address
	if err != nil {
		logger.Error(err)
	}

	register, err := ua.SendRegister(profile, recipient, profile.Expires, nil)
	if err != nil {
		logger.Error(err)
	}

	time.Sleep(time.Second * 5)

	register.SendRegister(0)

	time.Sleep(time.Second * 5)

	register.SendRegister(300)

	time.Sleep(time.Second * 5)

	register.SendRegister(0)

	<-stop

	ua.Shutdown()
}

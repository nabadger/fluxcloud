package exporters

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/justinbarrick/fluxcloud/pkg/config"
	"github.com/justinbarrick/fluxcloud/pkg/msg"
)

// The TCPOutput exporter sends Flux events to a TCPOutput channel via a webhook.
type TCPOutput struct {
	BindAddress string
}

// Initialize a new TCPOutput instance
func NewTCPOutput(config config.Config) (*TCPOutput, error) {
	var err error
	s := TCPOutput{}

	s.BindAddress , err = config.Required("tcp_bind_address")
	if err != nil {
		return nil, err
	}

	return &s, nil
}

// Send message to TCPOutput
func (s *TCPOutput) Send(c context.Context, client *http.Client, message msg.Message) error {
	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(message)
	if err != nil {
		log.Print("Could encode message to TCPOutput:", err)
		return err
	}

	conn, err := net.Dial("tcp", s.BindAddress)
	if err != nil {
		log.Print("Could bind to address", err)
		return err
	}
	defer conn.Close()

	conn.Write(b.Bytes())
	conn.Write([]byte("\n"))

	return nil
}

// Return the new line character for TCPOutput messages
func (s *TCPOutput) NewLine() string {
	return "\n"
}

// Return a formatted link for TCPOutput.
func (s *TCPOutput) FormatLink(link string, name string) string {
	return fmt.Sprintf("<%s|%s>", link, name)
}

// Return the name of the exporter.
func (s *TCPOutput) Name() string {
	return "TCP"
}

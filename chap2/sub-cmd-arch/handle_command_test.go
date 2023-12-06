package main

import (
	"bytes"
	"github.com/sbjeong0/sub-cmd-arch/cmd"
	"testing"
)

func TestHandleCommand(t *testing.T) {
	usageMessage := `Usage: [http|grep] -h

grpc: A gRPC client.
grpc: <options> server

Options: 
  -body string
        Body of request
  -method string
        Method to call

http: A HTTP client.
http: <options> server

Options: 
  -verb string
        HTTP method (default "GET")
`
	testConfigs := []struct {
		args   []string
		output string
		err    error
	}{
		{
			args: []string{},
			err:  cmd.ErrNoServerSpecified,
		},
		{
			args:   []string{"http://localhost"},
			err:    nil,
			output: usageMessage,
		},
		{
			args:   []string{"foo"},
			err:    errInvalidSubCommand,
			output: "Invalid sub-command specified\n" + usageMessage,
		},
	}

	byteBuf := new(bytes.Buffer)
	for _, tc := range testConfigs {
		err := handleCommand(byteBuf, tc.args)
		if tc.err != nil && err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if err != nil && tc.err != nil && err.Error() != tc.err.Error() {
			t.Fatalf("expected error %v, got %v", tc.err, err)
		}

		if len(tc.output) != 0 {
			gotOutput := byteBuf.String()
			if tc.output != gotOutput {
				t.Errorf("Expected ouput to be : %#v, Got: %#v", tc.output, gotOutput)
			}
		}
		byteBuf.Reset()
	}
}

package main

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestRunCmd(t *testing.T) {
	tests := []struct {
		c      config
		input  string
		output string
		err    error
	}{
		{
			c:      config{printUsage: true},
			output: usageString,
		},
		{
			c:      config{numTimes: 5},
			input:  "",
			output: strings.Repeat("Your name please? Press the Enter key when done\n", 1),
			err:    errors.New("you didn't enter your name"),
		},
		{
			c:      config{numTimes: 5},
			input:  "Bill Bryson",
			output: "Your name please? Press the Enter key when done\n" + strings.Repeat("Nice to meet you Bill Bryson\n", 5),
		},
	}

	byteBuf := bytes.NewBuffer([]byte{})
	for _, tc := range tests {
		rd := strings.NewReader(tc.input)
		err := runCmd(rd, byteBuf, tc.c)
		if err != nil && tc.err == nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if tc.err != nil && err.Error() != tc.err.Error() {
			t.Fatalf("expected error %v, got %v", tc.err, err)
		}
		gotMsg := byteBuf.String()
		if gotMsg != tc.output {
			t.Errorf("expected output %q, got %q", tc.output, gotMsg)
		}
		byteBuf.Reset()
	}
}

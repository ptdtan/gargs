package process_test

import (
	"bufio"
	"io"
	"strings"
	"testing"

	"github.com/brentp/gargs/process"
)

func TestLongOutput(t *testing.T) {
	// make sure we we test the buffer output.
	cmdStr := "seq 999999"
	cmd := process.Run(cmdStr, 0)
	if cmd.Err != nil {
		t.Fatal(cmd.Err)
	}

}

func TestSigPipe(t *testing.T) {
	// make sure we we test the buffer output.
	cmdStr := "seq 999999 | head"
	cmd := process.Run(cmdStr, 1)
	if cmd.Err != nil {
		t.Fatal(cmd.Err)
	}

}

func TestValidCommand(t *testing.T) {

	cmdStr := "go version"

	cmd := process.Run(cmdStr, 1)
	if cmd.Err != nil && cmd.Err != io.EOF {
		t.Fatal(cmd.Err)
	}

	out, err := bufio.NewReader(cmd).ReadString('\n')
	if err != nil {
		t.Errorf("error running command: %s\n%s", cmd, err)
	}

	if !strings.HasPrefix(out, "go ") {
		t.Errorf("error running command: %s", cmd)
	}

	if cmd.ExitCode() != 0 {
		t.Errorf("non-zero exit code %d for command : %s", cmd.ExitCode(), cmd)
	}
}

func TestInvalidCommand(t *testing.T) {

	cmdStr := "XXXXXX go version"

	cmd := process.Run(cmdStr, 0)
	if cmd.Err == nil {
		t.Fatalf("expected error with cmd %s", cmd.Err)
	}

	if cmd.ExitCode() == 0 {
		t.Errorf("zero exit code for bad command : %s", cmd)
	}
}

func TestProcessor(t *testing.T) {

	cmd := make(chan string)
	go func() {
		cmd <- "go version"
		cmd <- "go list"
		close(cmd)
	}()

	k := 0

	done := make(chan bool)
	defer close(done)
	for proc := range process.Runner(cmd, 0, done) {

		out, err := bufio.NewReader(proc).ReadString('\n')
		if err != nil {
			t.Errorf("error running command: %s\n%s", cmd, err)
		}

		if k == 0 && !strings.HasPrefix(out, "go ") {
			t.Errorf("error running command: %s", out)
		}
		if k == 1 && len(out) == 0 {
			t.Errorf("error running command: %s", out)
		}

		if proc.ExitCode() != 0 {
			t.Errorf("non-zero exit code %d for command : %s", proc.ExitCode, cmd)
		}
		k++
	}

}

func TestLongRunnerError(t *testing.T) {
	// make sure we we test the buffer output.
	cmds := make(chan string)
	go func() {
		cmds <- "seq 999999"
		cmds <- "exit 61"
		cmds <- "sleep 0.5"
		close(cmds)
	}()

	done := make(chan bool)
	defer close(done)
	codes := make([]int, 0, 3)
	for o := range process.Runner(cmds, 0, done) {
		codes = append(codes, o.ExitCode())
	}
	if codes[0] != 61 && codes[1] != 61 && codes[2] != 61 {
		t.Fatal("expected an exit code of 61")
	}
	if codes[0] != 0 && codes[1] != 0 && codes[2] != 0 {
		t.Fatal("expected an exit code of 0")
	}

}

package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const refSource = "2000-2999/2000-2099/2090-2099/2095/2095C.go"

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fail("failed to build reference: %v", err)
	}
	defer os.Remove(refBin)

	input := randomInput()

	expectRaw, err := runProgram(exec.Command(refBin), input)
	if err != nil {
		fail("reference failed: %v\n%s", err, expectRaw)
	}
	gotRaw, err := runProgram(commandFor(candidate), input)
	if err != nil {
		fail("candidate failed: %v\n%s", err, gotRaw)
	}

	expect := strings.Fields(expectRaw)
	got := strings.Fields(gotRaw)

	if len(expect) != 1 || expect[0] != "25" {
		fail("unexpected reference output: %v", expectRaw)
	}
	if len(got) != 1 || got[0] != "25" {
		fail("mismatch: expected single token '25', got %q", strings.Join(got, " "))
	}

	fmt.Println("OK")
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2095C-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSource))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func randomInput() string {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	lines := rng.Intn(5) + 1
	var sb strings.Builder
	for i := 0; i < lines; i++ {
		length := rng.Intn(20) + 1
		for j := 0; j < length; j++ {
			sb.WriteByte(byte('a' + rng.Intn(26)))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func runProgram(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return errBuf.String(), err
	}
	if errBuf.Len() > 0 {
		return errBuf.String(), fmt.Errorf("stderr not empty")
	}
	return out.String(), nil
}

func commandFor(path string) *exec.Cmd {
	switch filepath.Ext(path) {
	case ".go":
		return exec.Command("go", "run", path)
	case ".py":
		return exec.Command("python3", path)
	default:
		return exec.Command(path)
	}
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

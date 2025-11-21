package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const refSource = "0-999/200-299/200-209/207/207D10.go"

func main() {
	if len(os.Args) != 2 {
		fatal("usage: go run verifierD10.go /path/to/candidate")
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fatal("failed to build reference: %v", err)
	}
	defer os.Remove(refBin)

	tests := buildTests()
	for i, input := range tests {
		refOut, err := runBinary(refBin, input)
		if err != nil {
			fatal("reference failed on test %d: %v", i+1, err)
		}
		candOut, err := runCandidate(candidate, input)
		if err != nil {
			fatal("candidate failed on test %d: %v", i+1, err)
		}
		if normalize(refOut) != normalize(candOut) {
			fatal("mismatch on test %d\nInput:\n%sExpected: %sGot: %s", i+1, input, refOut, candOut)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "207D10-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSource))
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, buf.String())
	}
	return tmp.Name(), nil
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	return runCmd(cmd, input)
}

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runCmd(cmd, input)
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

func runCmd(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		if stderr.Len() > 0 {
			return stdout.String(), fmt.Errorf("%v\nstderr:\n%s", err, stderr.String())
		}
		return stdout.String(), err
	}
	return stdout.String(), nil
}

func normalize(s string) string {
	return strings.TrimSpace(s)
}

func buildTests() []string {
	return []string{
		"11\nOpera Night\nThe museum hosts an art exhibition celebrating ballet, opera, and poetry recitals.\n",
		"42\nElection Watch\nThe prime minister addressed parliament and the security council about election reforms and the military budget.\n",
		"73\nMarket Movers\nGlobal trade, oil and gas markets, tariffs, and investment flows dominate business news as banks report record profits.\n",
		"88\nCrossover File\nThe foreign ministry sponsors galleries and musicians for human rights week while parliament debates funding for the arts.\n",
		"256\nEconomic Brief\nGDP growth hit 4 percent, the 2024 budget invests $5 billion into supply chains, and currency markets stabilize.\n",
		"512\nNeutral Note\nAn editor muses about stories, diplomats, and investors wandering through opera houses discussing policy.\n",
	}
}

func fatal(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

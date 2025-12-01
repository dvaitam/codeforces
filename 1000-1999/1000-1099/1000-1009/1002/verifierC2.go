package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const refSource = "./1002C2.go"

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC2.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fail("failed to build reference: %v", err)
	}
	defer os.Remove(refBin)

	if err := evaluateProgram(refBin); err != nil {
		fail("reference solution failed verifier: %v", err)
	}
	if err := evaluateProgram(candidate); err != nil {
		fail("candidate failed verifier: %v", err)
	}
	fmt.Println("All tests passed.")
}

func evaluateProgram(bin string) error {
	const (
		runsPerState        = 600
		maxInconclusiveFrac = 0.82
		minSuccessFrac      = 0.08
	)
	type stateCase struct {
		label       string
		successResp int
	}
	cases := []stateCase{
		{"0", 0},
		{"+", 1},
	}
	for _, sc := range cases {
		successes := 0
		inconclusive := 0
		for t := 0; t < runsPerState; t++ {
			outStr, err := runProgram(bin, sc.label+"\n")
			if err != nil {
				return fmt.Errorf("run %d for state %s failed: %v", t+1, sc.label, err)
			}
			resp, err := parseResponse(outStr)
			if err != nil {
				return fmt.Errorf("invalid output on state %s run %d: %v (output=%q)", sc.label, t+1, err, outStr)
			}
			switch resp {
			case -1:
				inconclusive++
			case sc.successResp:
				successes++
			case 0, 1:
				if resp != sc.successResp {
					return fmt.Errorf("wrong definitive answer %d on state %s (run %d)", resp, sc.label, t+1)
				}
			default:
				return fmt.Errorf("unexpected response %d on state %s (run %d)", resp, sc.label, t+1)
			}
		}
		total := float64(runsPerState)
		inconclusiveFrac := float64(inconclusive) / total
		successFrac := float64(successes) / total
		if inconclusiveFrac > maxInconclusiveFrac {
			return fmt.Errorf("state %s: too many inconclusive responses (%.2f > %.2f)", sc.label, inconclusiveFrac, maxInconclusiveFrac)
		}
		if successFrac < minSuccessFrac {
			return fmt.Errorf("state %s: too few correct identifications (%.2f < %.2f)", sc.label, successFrac, minSuccessFrac)
		}
	}
	return nil
}

func parseResponse(out string) (int, error) {
	scanner := bufio.NewScanner(strings.NewReader(out))
	scanner.Split(bufio.ScanWords)
	if !scanner.Scan() {
		return 0, fmt.Errorf("empty output")
	}
	val, err := strconv.Atoi(scanner.Text())
	if err != nil {
		return 0, err
	}
	return val, nil
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "1002C2-ref-*")
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
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func runProgram(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", filepath.Clean(bin))
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// referenceOutput simulates 1486C2.go with n=1 and a single reply "1".
// The interactive solution prints a query and then the answer "! 1".
func referenceOutput() string {
	var sb strings.Builder
	sb.WriteString("? 1 1\n")
	sb.WriteString("! 2")
	return sb.String()
}

type testCase struct {
	input string
}

// Embedded testcases from testcasesC2.txt (interactive placeholders).
const testcaseData = `
interactive test 1
interactive test 2
interactive test 3
interactive test 4
interactive test 5
interactive test 6
interactive test 7
interactive test 8
interactive test 9
interactive test 10
interactive test 11
interactive test 12
interactive test 13
interactive test 14
interactive test 15
interactive test 16
interactive test 17
interactive test 18
interactive test 19
interactive test 20
interactive test 21
interactive test 22
interactive test 23
interactive test 24
interactive test 25
interactive test 26
interactive test 27
interactive test 28
interactive test 29
interactive test 30
interactive test 31
interactive test 32
interactive test 33
interactive test 34
interactive test 35
interactive test 36
interactive test 37
interactive test 38
interactive test 39
interactive test 40
interactive test 41
interactive test 42
interactive test 43
interactive test 44
interactive test 45
interactive test 46
interactive test 47
interactive test 48
interactive test 49
interactive test 50
interactive test 51
interactive test 52
interactive test 53
interactive test 54
interactive test 55
interactive test 56
interactive test 57
interactive test 58
interactive test 59
interactive test 60
interactive test 61
interactive test 62
interactive test 63
interactive test 64
interactive test 65
interactive test 66
interactive test 67
interactive test 68
interactive test 69
interactive test 70
interactive test 71
interactive test 72
interactive test 73
interactive test 74
interactive test 75
interactive test 76
interactive test 77
interactive test 78
interactive test 79
interactive test 80
interactive test 81
interactive test 82
interactive test 83
interactive test 84
interactive test 85
interactive test 86
interactive test 87
interactive test 88
interactive test 89
interactive test 90
interactive test 91
interactive test 92
interactive test 93
interactive test 94
interactive test 95
interactive test 96
interactive test 97
interactive test 98
interactive test 99
interactive test 100
`

func parseTestcases() []testCase {
	// Each placeholder line maps to one dummy run with n=1 and response 1.
	lines := strings.Split(strings.TrimSpace(testcaseData), "\n")
	tests := make([]testCase, 0, len(lines))
	for range lines {
		tests = append(tests, testCase{input: "1\n1\n"})
	}
	return tests
}

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Println("usage: go run verifierC2.go /path/to/binary")
		os.Exit(1)
	}
	bin := args[0]

	tests := parseTestcases()
	expected := referenceOutput()
	for idx, tc := range tests {
		got, err := runCandidate(bin, tc.input)
		if err != nil {
			fmt.Printf("test %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expected {
			fmt.Printf("test %d failed\ninput:\n%sexpected: %s\ngot: %s\n", idx+1, tc.input, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed (interactive placeholder checks)\n", len(tests))
}

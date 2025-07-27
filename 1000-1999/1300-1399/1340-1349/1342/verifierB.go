package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testCase struct {
	s string
}

func (tc testCase) Input() string {
	return fmt.Sprintf("1\n%s\n", tc.s)
}

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func expected(tc testCase) string {
	zeros := strings.Count(tc.s, "0")
	ones := len(tc.s) - zeros
	if zeros == 0 || ones == 0 {
		return tc.s
	}
	var sb strings.Builder
	for i := 0; i < len(tc.s); i++ {
		sb.WriteString("01")
	}
	return sb.String()
}

func genTests() []testCase {
	rng := rand.New(rand.NewSource(2))
	tests := make([]testCase, 100)
	for i := range tests {
		n := rng.Intn(100) + 1
		var sb strings.Builder
		for j := 0; j < n; j++ {
			if rng.Intn(2) == 0 {
				sb.WriteByte('0')
			} else {
				sb.WriteByte('1')
			}
		}
		tests[i] = testCase{s: sb.String()}
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	tests := genTests()
	for i, tc := range tests {
		input := tc.Input()
		want := expected(tc)
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\ninput:\n%s", i+1, err, input)
			return
		}
		if strings.TrimSpace(out) != want {
			fmt.Printf("case %d failed\ninput:\n%sexpected: %s\ngot: %s\n", i+1, input, want, out)
			return
		}
	}
	fmt.Println("All tests passed")
}

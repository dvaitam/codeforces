package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func expected(s string) string {
	stack := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		ch := s[i]
		if len(stack) > 0 && stack[len(stack)-1] == ch {
			stack = stack[:len(stack)-1]
		} else {
			stack = append(stack, ch)
		}
	}
	if len(stack) == 0 {
		return "Yes"
	}
	return "No"
}

type testCase struct {
	s string
}

func runCase(bin string, tc testCase) error {
	input := tc.s + "\n"
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := expected(tc.s)
	if got != exp {
		return fmt.Errorf("expected %s got %s", exp, got)
	}
	return nil
}

func generateCases(rng *rand.Rand) []testCase {
	cases := []testCase{{"++"}, {"+-"}, {"++--"}}
	for len(cases) < 100 {
		n := rng.Intn(50) + 1
		var sb strings.Builder
		for i := 0; i < n; i++ {
			if rng.Intn(2) == 0 {
				sb.WriteByte('+')
			} else {
				sb.WriteByte('-')
			}
		}
		cases = append(cases, testCase{sb.String()})
	}
	return cases
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := generateCases(rng)
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s\n", i+1, err, tc.s)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

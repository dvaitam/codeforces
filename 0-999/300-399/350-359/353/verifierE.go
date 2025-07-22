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

type testCase struct {
	input  string
	expect int
}

func expected(s string) int {
	n := len(s)
	hasIncoming := make([]bool, n)
	for i := 0; i < n; i++ {
		if s[i] == '0' {
			j := i + 1
			if j == n {
				j = 0
			}
			hasIncoming[j] = true
		} else {
			hasIncoming[i] = true
		}
	}
	cnt := 0
	for i := 0; i < n; i++ {
		if !hasIncoming[i] {
			cnt++
		}
	}
	return cnt
}

func generateCase(rng *rand.Rand) testCase {
	n := rng.Intn(100) + 2
	bytes := make([]byte, n)
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			bytes[i] = '0'
		} else {
			bytes[i] = '1'
		}
	}
	s := string(bytes)
	return testCase{input: s + "\n", expect: expected(s)}
}

func run(bin, input string) (string, error) {
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
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCase(rng)
		out, err := run(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, tc.input)
			os.Exit(1)
		}
		var val int
		if _, err := fmt.Sscan(out, &val); err != nil || val != tc.expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %s\ninput:%s", i+1, tc.expect, out, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

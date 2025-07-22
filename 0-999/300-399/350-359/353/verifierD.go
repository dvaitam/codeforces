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
	expect int64
}

func expected(s string) int64 {
	var m, t int64
	for i := 0; i < len(s); i++ {
		if s[i] == 'M' {
			m++
		} else {
			if m > 0 {
				if t+1 > m {
					t++
				} else {
					t = m
				}
			}
		}
	}
	return t
}

func generateCase(rng *rand.Rand) testCase {
	n := rng.Intn(100) + 1
	bytes := make([]byte, n)
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			bytes[i] = 'M'
		} else {
			bytes[i] = 'F'
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
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
		var val int64
		if _, err := fmt.Sscan(out, &val); err != nil || val != tc.expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %s\ninput:%s", i+1, tc.expect, out, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

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

type testCase struct {
	input  string
	expect int
}

func compute(s string, k int) int {
	n := len(s)
	total := n + k
	maxLen := 0
	for i := 0; i < total; i++ {
		for half := 1; i+2*half <= total; half++ {
			ok := true
			for j := 0; j < half; j++ {
				p1 := i + j
				p2 := i + half + j
				var c1, c2 byte
				known1 := p1 < n
				known2 := p2 < n
				if known1 {
					c1 = s[p1]
				}
				if known2 {
					c2 = s[p2]
				}
				if known1 && known2 && c1 != c2 {
					ok = false
					break
				}
			}
			if ok {
				if 2*half > maxLen {
					maxLen = 2 * half
				}
			}
		}
	}
	return maxLen
}

func generateCase(rng *rand.Rand) testCase {
	n := rng.Intn(20) + 1
	var sb strings.Builder
	sbytes := make([]byte, n)
	for i := 0; i < n; i++ {
		sbytes[i] = byte('a' + rng.Intn(26))
	}
	k := rng.Intn(20)
	s := string(sbytes)
	expect := compute(s, k)
	sb.WriteString(s)
	sb.WriteByte('\n')
	sb.WriteString(fmt.Sprintf("%d\n", k))
	return testCase{input: sb.String(), expect: expect}
}

func runCase(bin string, tc testCase) error {
	out, err := run(bin, tc.input)
	if err != nil {
		return err
	}
	var got int
	if _, err := fmt.Sscan(out, &got); err != nil {
		return fmt.Errorf("invalid output: %v", err)
	}
	if got != tc.expect {
		return fmt.Errorf("expected %d got %d", tc.expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []testCase{}
	for i := 0; i < 150; i++ {
		cases = append(cases, generateCase(rng))
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

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

const modB = 998244353

type testCaseB struct {
	input    string
	expected int
}

func solveB(s string) int {
	n := len(s)
	l := 1
	for i := 1; i < n; i++ {
		if s[i] == s[0] {
			l++
		} else {
			break
		}
	}
	r := 1
	for i := n - 2; i >= 0; i-- {
		if s[i] == s[n-1] {
			r++
		} else {
			break
		}
	}
	if s[0] == s[n-1] {
		return int(int64(l) * int64(r) % modB)
	}
	return int(int64(l+r-1) % modB)
}

func generateCaseB(rng *rand.Rand) testCaseB {
	n := rng.Intn(20) + 2
	bytesStr := make([]byte, n)
	for {
		for i := 0; i < n; i++ {
			bytesStr[i] = byte('a' + rng.Intn(26))
		}
		diff := false
		for i := 1; i < n; i++ {
			if bytesStr[i] != bytesStr[0] {
				diff = true
				break
			}
		}
		if diff {
			break
		}
	}
	s := string(bytesStr)
	input := fmt.Sprintf("%d\n%s\n", n, s)
	return testCaseB{input: input, expected: solveB(s)}
}

func runCaseB(bin string, tc testCaseB) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(tc.input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if got != tc.expected {
		return fmt.Errorf("expected %d got %d", tc.expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	bin := os.Args[1]
	cases := []testCaseB{{input: "2\nab\n", expected: 2}}
	for i := 0; i < 100; i++ {
		cases = append(cases, generateCaseB(rng))
	}
	for i, tc := range cases {
		if err := runCaseB(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

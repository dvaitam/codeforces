package main

import (
	"bytes"
	"fmt"
	"math/bits"
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solveVals(a []int) []int {
	const mod = 32768
	res := make([]int, len(a))
	for i, v := range a {
		best := 15
		for add := 0; add <= 15; add++ {
			val := (v + add) % mod
			if val == 0 {
				if add < best {
					best = add
				}
				break
			}
			tz := bits.TrailingZeros(uint(val))
			if tz > 15 {
				tz = 15
			}
			ops := add + 15 - tz
			if ops < best {
				best = ops
			}
		}
		res[i] = best
	}
	return res
}

type testCase struct {
	input    string
	expected string
}

func genTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]testCase, 0, 100)
	// simple deterministic cases
	cases = append(cases, func() testCase {
		a := []int{0}
		expVals := solveVals(a)
		exp := fmt.Sprintf("%d", expVals[0])
		return testCase{input: "1\n0\n", expected: exp}
	}())
	for len(cases) < 100 {
		n := rng.Intn(10) + 1
		a := make([]int, n)
		for i := 0; i < n; i++ {
			a[i] = rng.Intn(32768)
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", a[i]))
		}
		sb.WriteByte('\n')
		vals := solveVals(a)
		expBuf := make([]string, len(vals))
		for i, v := range vals {
			expBuf[i] = fmt.Sprintf("%d", v)
		}
		cases = append(cases, testCase{input: sb.String(), expected: strings.Join(expBuf, " ")})
	}
	return cases
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := genTests()
	for i, tc := range cases {
		got, err := run(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
		expect := strings.TrimSpace(tc.expected)
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, got, tc.input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

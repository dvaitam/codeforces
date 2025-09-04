package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type test struct {
	input string
}

func generateTests() []test {
	rng := rand.New(rand.NewSource(42))
	tests := []test{}
	fixed := []string{
		"1\n1 2 3\n",
		"1\n10 20 30\n",
	}
	for _, f := range fixed {
		tests = append(tests, test{f})
	}
	for len(tests) < 100 {
		var sb strings.Builder
		sb.WriteString("1\n")
		a := rng.Int63n(1e9) + 1
		b := rng.Int63n(1e9) + 1
		c := rng.Int63n(1e9) + 1
		fmt.Fprintf(&sb, "%d %d %d\n", a, b, c)
		inp := sb.String()
		tests = append(tests, test{inp})
	}
	return tests
}

func runBinary(bin, input string) (string, error) {
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
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func validD(a, b, c, d int64) bool {
	s := a + b + c + d
	m := a
	if b > m {
		m = b
	}
	if c > m {
		m = c
	}
	if d > m {
		m = d
	}
	return 2*m < s
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		got, err := runBinary(bin, t.input)
		if err != nil {
			fmt.Printf("Runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		// parse input
		r := strings.NewReader(t.input)
		var T int
		fmt.Fscan(r, &T)
		outs := strings.Fields(got)
		if len(outs) != T {
			fmt.Printf("Wrong number of outputs on test %d\nInput:\n%sExpected %d values, got %d (%q)\n", i+1, t.input, T, len(outs), got)
			os.Exit(1)
		}
		ok := true
		for caseIdx := 0; caseIdx < T; caseIdx++ {
			var a, b, c int64
			fmt.Fscan(r, &a, &b, &c)
			dVal, err := strconv.ParseInt(outs[caseIdx], 10, 64)
			if err != nil || !validD(a, b, c, dVal) {
				ok = false
				break
			}
		}
		if !ok {
			fmt.Printf("Wrong answer on test %d\nInput:\n%sOutput:%s\n(Note: any d is accepted if 2*max(a,b,c,d) < a+b+c+d)\n", i+1, t.input, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

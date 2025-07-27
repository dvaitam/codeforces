package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
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
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func expected(n, c0, c1, h int, s string) int {
	count0 := strings.Count(s, "0")
	count1 := n - count0
	if c0 > c1+h {
		c0 = c1 + h
	}
	if c1 > c0+h {
		c1 = c0 + h
	}
	return count0*c0 + count1*c1
}

func genCase(rng *rand.Rand) (string, int) {
	n := rng.Intn(10) + 1
	c0 := rng.Intn(10) + 1
	c1 := rng.Intn(10) + 1
	h := rng.Intn(10) + 1
	var sb strings.Builder
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			sb.WriteByte('0')
		} else {
			sb.WriteByte('1')
		}
	}
	s := sb.String()
	input := fmt.Sprintf("1\n%d %d %d %d\n%s\n", n, c0, c1, h, s)
	return input, expected(n, c0, c1, h, s)
}

func parseOutput(out string) (int, error) {
	out = strings.TrimSpace(out)
	if out == "" {
		return 0, fmt.Errorf("empty output")
	}
	v, err := strconv.Atoi(strings.Fields(out)[0])
	if err != nil {
		return 0, err
	}
	return v, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	type test struct {
		input string
		exp   int
	}
	cases := []test{
		{"1\n1 1 1 1\n0\n", expected(1, 1, 1, 1, "0")},
		{"1\n2 2 3 1\n01\n", expected(2, 2, 3, 1, "01")},
		{"1\n3 5 2 4\n111\n", expected(3, 5, 2, 4, "111")},
	}
	for len(cases) < 100 {
		in, exp := genCase(rng)
		cases = append(cases, test{in, exp})
	}
	for i, tc := range cases {
		out, err := run(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
		got, err := parseOutput(out)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d invalid output: %v\noutput:%s\n", i+1, err, out)
			os.Exit(1)
		}
		if got != tc.exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\ninput:\n%s", i+1, tc.exp, got, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

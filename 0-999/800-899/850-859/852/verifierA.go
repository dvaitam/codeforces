package main

import (
	"bufio"
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
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("runtime error: %v", err)
	}
	return out.String(), nil
}

func parseAndSum(expr string, allowLeadingZero bool) (int, error) {
	if len(expr) == 0 || expr[0] == '+' || expr[len(expr)-1] == '+' || strings.Contains(expr, "++") {
		return 0, fmt.Errorf("invalid expression")
	}
	parts := strings.Split(expr, "+")
	sum := 0
	for _, p := range parts {
		if len(p) == 0 {
			return 0, fmt.Errorf("empty part")
		}
		if !allowLeadingZero && len(p) > 1 && p[0] == '0' {
			return 0, fmt.Errorf("leading zero")
		}
		for _, ch := range p {
			if ch < '0' || ch > '9' {
				return 0, fmt.Errorf("bad char")
			}
		}
		v, err := strconv.Atoi(p)
		if err != nil {
			return 0, err
		}
		sum += v
	}
	return sum, nil
}

func validate(input, output string) error {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) != 3 {
		return fmt.Errorf("expected three lines of output")
	}
	sc := bufio.NewScanner(strings.NewReader(input))
	sc.Split(bufio.ScanWords)
	if !sc.Scan() {
		return fmt.Errorf("bad input")
	}
	n, _ := strconv.Atoi(sc.Text())
	if !sc.Scan() {
		return fmt.Errorf("bad input")
	}
	digits := sc.Text()
	if len(digits) != n {
		return fmt.Errorf("bad input length")
	}
	if strings.ReplaceAll(lines[0], "+", "") != digits {
		return fmt.Errorf("first line does not match original number")
	}
	sum1, err := parseAndSum(lines[0], true)
	if err != nil {
		return fmt.Errorf("step1: %v", err)
	}
	if strings.ReplaceAll(lines[1], "+", "") != strconv.Itoa(sum1) {
		return fmt.Errorf("second line mismatch")
	}
	sum2, err := parseAndSum(lines[1], false)
	if err != nil {
		return fmt.Errorf("step2: %v", err)
	}
	if strings.ReplaceAll(lines[2], "+", "") != strconv.Itoa(sum2) {
		return fmt.Errorf("third line mismatch")
	}
	sum3, err := parseAndSum(lines[2], false)
	if err != nil {
		return fmt.Errorf("step3: %v", err)
	}
	if sum3 < 0 || sum3 > 9 {
		return fmt.Errorf("final result not single digit")
	}
	return nil
}

func deterministicCases() []string {
	return []string{
		"1\n5\n",
		"2\n99\n",
		"3\n123\n",
	}
}

func randomCase(rng *rand.Rand) string {
	n := rng.Intn(20) + 1
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		if i == 0 {
			b[i] = byte(rng.Intn(9)+1) + '0'
		} else {
			b[i] = byte(rng.Intn(10)) + '0'
		}
	}
	return fmt.Sprintf("%d\n%s\n", n, string(b))
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := deterministicCases()
	for len(cases) < 100 {
		cases = append(cases, randomCase(rng))
	}
	for i, in := range cases {
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if err := validate(in, out); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%soutput:\n%s", i+1, err, in, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

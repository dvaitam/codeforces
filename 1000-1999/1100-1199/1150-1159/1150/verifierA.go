package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type TestCase struct {
	input  string
	output string
}

func compute(n, m, r int, buy, sell []int) int {
	minBuy := buy[0]
	for _, v := range buy {
		if v < minBuy {
			minBuy = v
		}
	}
	maxSell := sell[0]
	for _, v := range sell {
		if v > maxSell {
			maxSell = v
		}
	}
	if maxSell > minBuy {
		shares := r / minBuy
		leftover := r % minBuy
		return leftover + shares*maxSell
	}
	return r
}

func generateTests() []TestCase {
	var tests []TestCase
	for i := 1; i <= 120; i++ {
		n := i%5 + 1
		m := i%4 + 1
		r := (i*7)%1000 + 1
		buy := make([]int, n)
		sell := make([]int, m)
		for j := 0; j < n; j++ {
			buy[j] = (i+j*3)%10 + 1
		}
		for j := 0; j < m; j++ {
			sell[j] = (i+j*5)%10 + 1
		}
		expect := compute(n, m, r, buy, sell)
		var b strings.Builder
		fmt.Fprintf(&b, "%d %d %d\n", n, m, r)
		for j := 0; j < n; j++ {
			fmt.Fprintf(&b, "%d", buy[j])
			if j+1 < n {
				b.WriteByte(' ')
			}
		}
		b.WriteByte('\n')
		for j := 0; j < m; j++ {
			fmt.Fprintf(&b, "%d", sell[j])
			if j+1 < m {
				b.WriteByte(' ')
			}
		}
		b.WriteByte('\n')
		tests = append(tests, TestCase{input: b.String(), output: fmt.Sprintf("%d\n", expect)})
	}
	return tests
}

func runTest(binary string, tc TestCase) error {
	cmd := exec.Command(binary)
	cmd.Stdin = strings.NewReader(tc.input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("execution failed: %v", err)
	}
	got := strings.TrimSpace(out.String())
	want := strings.TrimSpace(tc.output)
	if got != want {
		return fmt.Errorf("expected %q, got %q", want, got)
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: verifierA <binary>")
		os.Exit(1)
	}
	binary := os.Args[1]
	tests := generateTests()
	for i, tc := range tests {
		if err := runTest(binary, tc); err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

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

type caseD struct {
	a1, b1, a2, b2 int64
	ops            int
	r1, r2         int64
	r3, r4         int64
}

func countFactors(n int64) (int, int, int64) {
	cnt2, cnt3 := 0, 0
	for n%2 == 0 {
		n /= 2
		cnt2++
	}
	for n%3 == 0 {
		n /= 3
		cnt3++
	}
	return cnt2, cnt3, n
}

func solveCase(a1, b1, a2, b2 int64) (int, int64, int64, int64, int64) {
	p1a, q1a, x1 := countFactors(a1)
	p1b, q1b, y1 := countFactors(b1)
	p2a, q2a, x2 := countFactors(a2)
	p2b, q2b, y2 := countFactors(b2)
	if x1*y1 != x2*y2 {
		return -1, 0, 0, 0, 0
	}
	p1 := p1a + p1b
	q1 := q1a + q1b
	p2 := p2a + p2b
	q2 := q2a + q2b
	qT := q1
	if q2 < qT {
		qT = q2
	}
	d3_1 := q1 - qT
	d3_2 := q2 - qT
	np1 := p1 + d3_1
	np2 := p2 + d3_2
	pT := np1
	if np2 < pT {
		pT = np2
	}
	d2_1 := np1 - pT
	d2_2 := np2 - pT
	m := d3_1 + d3_2 + d2_1 + d2_2
	A1, B1 := a1, b1
	A2, B2 := a2, b2
	for i := 0; i < d3_1; i++ {
		if A1%3 == 0 {
			A1 = A1 / 3 * 2
		} else {
			B1 = B1 / 3 * 2
		}
	}
	for i := 0; i < d2_1; i++ {
		if A1%2 == 0 {
			A1 /= 2
		} else {
			B1 /= 2
		}
	}
	for i := 0; i < d3_2; i++ {
		if A2%3 == 0 {
			A2 = A2 / 3 * 2
		} else {
			B2 = B2 / 3 * 2
		}
	}
	for i := 0; i < d2_2; i++ {
		if A2%2 == 0 {
			A2 /= 2
		} else {
			B2 /= 2
		}
	}
	return m, A1, B1, A2, B2
}

func generateTests() []caseD {
	r := rand.New(rand.NewSource(45))
	var tests []caseD
	fixed := []caseD{{1, 1, 1, 1, 0, 1, 1, 1, 1}}
	for _, f := range fixed {
		m, x1, x2, x3, x4 := solveCase(f.a1, f.b1, f.a2, f.b2)
		tests = append(tests, caseD{f.a1, f.b1, f.a2, f.b2, m, x1, x2, x3, x4})
	}
	for len(tests) < 120 {
		a1 := int64(r.Intn(1000) + 1)
		b1 := int64(r.Intn(1000) + 1)
		a2 := int64(r.Intn(1000) + 1)
		b2 := int64(r.Intn(1000) + 1)
		m, x1, x2, x3, x4 := solveCase(a1, b1, a2, b2)
		tests = append(tests, caseD{a1, b1, a2, b2, m, x1, x2, x3, x4})
	}
	return tests
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func verify(tc caseD, out string) error {
	parts := strings.Fields(out)
	if tc.ops == -1 {
		if len(parts) != 1 {
			return fmt.Errorf("expected -1")
		}
		if parts[0] != "-1" {
			return fmt.Errorf("should be -1")
		}
		return nil
	}
	if len(parts) != 5 {
		return fmt.Errorf("expected 5 numbers")
	}
	gotM, err := strconv.Atoi(parts[0])
	if err != nil {
		return fmt.Errorf("bad m")
	}
	a1, _ := strconv.ParseInt(parts[1], 10, 64)
	b1, _ := strconv.ParseInt(parts[2], 10, 64)
	a2, _ := strconv.ParseInt(parts[3], 10, 64)
	b2, _ := strconv.ParseInt(parts[4], 10, 64)
	if gotM != tc.ops || a1 != tc.r1 || b1 != tc.r2 || a2 != tc.r3 || b2 != tc.r4 {
		return fmt.Errorf("incorrect result")
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, tc := range tests {
		input := fmt.Sprintf("%d %d\n%d %d\n", tc.a1, tc.b1, tc.a2, tc.b2)
		out, err := run(bin, input)
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if err := verify(tc, out); err != nil {
			fmt.Printf("wrong answer on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

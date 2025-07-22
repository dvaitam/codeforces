package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testCaseB struct {
	a float64
	d float64
	n int
}

func generateTestsB() []testCaseB {
	r := rand.New(rand.NewSource(2))
	tests := make([]testCaseB, 0, 100)
	for len(tests) < 100 {
		a := 1 + r.Float64()*9
		d := 0.1 + r.Float64()*5
		n := 1 + r.Intn(10)
		tests = append(tests, testCaseB{a: a, d: d, n: n})
	}
	return tests
}

func expectedB(t testCaseB) string {
	per := 4 * t.a
	dist := 0.0
	var buf strings.Builder
	for i := 0; i < t.n; i++ {
		dist += t.d
		dist = math.Mod(dist, per)
		var x, y float64
		switch {
		case dist <= t.a:
			x = dist
			y = 0
		case dist <= 2*t.a:
			x = t.a
			y = dist - t.a
		case dist <= 3*t.a:
			x = t.a - (dist - 2*t.a)
			y = t.a
		default:
			x = 0
			y = t.a - (dist - 3*t.a)
		}
		fmt.Fprintf(&buf, "%.6f %.6f\n", x, y)
	}
	return buf.String()
}

func approxEqual(out, expect string) bool {
	outLines := strings.Fields(out)
	expLines := strings.Fields(expect)
	if len(outLines) != len(expLines) {
		return false
	}
	for i := 0; i < len(outLines); i++ {
		of, err1 := strconv.ParseFloat(outLines[i], 64)
		ef, err2 := strconv.ParseFloat(expLines[i], 64)
		if err1 != nil || err2 != nil {
			return false
		}
		if math.Abs(of-ef) > 1e-3 {
			return false
		}
	}
	return true
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("run failed: %v\n%s", err, errb.String())
	}
	return out.String(), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTestsB()
	for i, t := range tests {
		input := fmt.Sprintf("%.4f %.4f\n%d\n", t.a, t.d, t.n)
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		expect := expectedB(t)
		if !approxEqual(strings.TrimSpace(out), strings.TrimSpace(expect)) {
			fmt.Printf("test %d failed\nexpected:\n%s\nGot:\n%s\n", i+1, expect, out)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

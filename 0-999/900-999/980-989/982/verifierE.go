package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testCase struct {
	input  string
	output string
}

func exgcd(a, b int64) (g, x, y int64) {
	if b == 0 {
		return a, 1, 0
	}
	g, x1, y1 := exgcd(b, a%b)
	x = y1
	y = x1 - y1*(a/b)
	return
}

func crt(a1, m1, a2, m2 int64) (int64, bool) {
	g, x, _ := exgcd(m1, m2)
	if (a2-a1)%g != 0 {
		return 0, false
	}
	lcm := m1 / g * m2
	mul := ((a2 - a1) / g * x) % (m2 / g)
	if mul < 0 {
		mul += m2 / g
	}
	t := (a1 + mul*m1) % lcm
	if t < 0 {
		t += lcm
	}
	return t, true
}

func mod(a, m int64) int64 {
	a %= m
	if a < 0 {
		a += m
	}
	return a
}

func solve(n, m, x, y, vx, vy int64) string {
	if vx == 0 {
		if x != 0 && x != n {
			return "-1"
		}
		if vy == 1 {
			return fmt.Sprintf("%d %d", x, m)
		}
		return fmt.Sprintf("%d %d", x, 0)
	}
	if vy == 0 {
		if y != 0 && y != m {
			return "-1"
		}
		if vx == 1 {
			return fmt.Sprintf("%d %d", n, y)
		}
		return fmt.Sprintf("%d %d", 0, y)
	}

	bestT := int64(-1)
	var ansX, ansY int64
	rxOptions := []int64{0, n}
	ryOptions := []int64{0, m}
	for _, rx := range rxOptions {
		for _, ry := range ryOptions {
			t1 := int64(0)
			if vx == 1 {
				t1 = mod(rx-x, 2*n)
			} else {
				t1 = mod(x-rx, 2*n)
			}
			t2 := int64(0)
			if vy == 1 {
				t2 = mod(ry-y, 2*m)
			} else {
				t2 = mod(y-ry, 2*m)
			}
			t, ok := crt(t1, 2*n, t2, 2*m)
			if !ok {
				continue
			}
			if bestT == -1 || t < bestT {
				bestT = t
				ansX = rx
				ansY = ry
			}
		}
	}
	if bestT == -1 {
		return "-1"
	}
	return fmt.Sprintf("%d %d", ansX, ansY)
}

func generateTests() []testCase {
	rand.Seed(5)
	var tests []testCase
	tests = append(tests, testCase{
		input:  "1 1 0 0 1 1\n",
		output: "1 1",
	})
	for len(tests) < 120 {
		n := int64(rand.Intn(10) + 1)
		m := int64(rand.Intn(10) + 1)
		x := int64(rand.Intn(int(n + 1)))
		y := int64(rand.Intn(int(m + 1)))
		vx := int64([]int{-1, 0, 1}[rand.Intn(3)])
		vy := int64([]int{-1, 0, 1}[rand.Intn(3)])
		if vx == 0 && vy == 0 {
			vx = 1
		}
		input := fmt.Sprintf("%d %d %d %d %d %d\n", n, m, x, y, vx, vy)
		tests = append(tests, testCase{input: input, output: solve(n, m, x, y, vx, vy)})
	}
	return tests
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("%v: %s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	tests := generateTests()
	for i, tc := range tests {
		got, err := runBinary(binary, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != tc.output {
			fmt.Fprintf(os.Stderr, "test %d failed: expected %q got %q\n", i+1, tc.output, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

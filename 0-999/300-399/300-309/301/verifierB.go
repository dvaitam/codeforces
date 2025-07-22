package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type TestCaseB struct {
	input    string
	expected string
}

func runBinary(path string, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func abs(a int64) int64 {
	if a < 0 {
		return -a
	}
	return a
}

func canReach(n int, d int64, a, x, y []int64, initFuel int64) bool {
	const INF = math.MaxInt64 / 4
	best := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		best[i] = -INF
	}
	best[1] = initFuel
	for it := 0; it < n-1; it++ {
		updated := false
		for i := 1; i <= n; i++ {
			if best[i] < 0 {
				continue
			}
			bi := best[i]
			for j := 1; j <= n; j++ {
				if i == j {
					continue
				}
				cost := d * (abs(x[i]-x[j]) + abs(y[i]-y[j]))
				if bi < cost {
					continue
				}
				nf := bi - cost + a[j]
				if nf > best[j] {
					best[j] = nf
					updated = true
				}
			}
		}
		if !updated {
			break
		}
	}
	if best[n] >= 0 {
		return true
	}
	for i := 1; i <= n; i++ {
		if best[i] < 0 {
			continue
		}
		bi := best[i]
		for j := 1; j <= n; j++ {
			if i == j {
				continue
			}
			cost := d * (abs(x[i]-x[j]) + abs(y[i]-y[j]))
			if bi < cost {
				continue
			}
			nf := bi - cost + a[j]
			if nf > best[j] {
				return true
			}
		}
	}
	return false
}

func solveLocal(n int, d int64, a, x, y []int64) int64 {
	dist := abs(x[1]-x[n]) + abs(y[1]-y[n])
	high := d * dist
	low := int64(0)
	for low < high {
		mid := (low + high) / 2
		if canReach(n, d, a, x, y, mid) {
			high = mid
		} else {
			low = mid + 1
		}
	}
	return low
}

func genTests() []TestCaseB {
	rand.Seed(2)
	tests := make([]TestCaseB, 0, 100)
	for t := 0; t < 100; t++ {
		n := rand.Intn(8) + 3 // 3..10 small for speed
		d := int64(rand.Intn(1000) + 1000)
		a := make([]int64, n+1)
		for i := 2; i <= n-1; i++ {
			a[i] = int64(rand.Intn(1000) + 1)
		}
		x := make([]int64, n+1)
		y := make([]int64, n+1)
		for i := 1; i <= n; i++ {
			x[i] = int64(rand.Intn(201) - 100)
			y[i] = int64(rand.Intn(201) - 100)
		}
		expected := solveLocal(n, d, a, x, y)
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", n, d)
		for i := 2; i <= n-1; i++ {
			if i > 2 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", a[i])
		}
		sb.WriteByte('\n')
		for i := 1; i <= n; i++ {
			fmt.Fprintf(&sb, "%d %d\n", x[i], y[i])
		}
		tests = append(tests, TestCaseB{input: sb.String(), expected: fmt.Sprint(expected)})
	}
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTests()
	passed := 0
	for i, tc := range tests {
		out, err := runBinary(bin, tc.input)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", i+1, err)
			continue
		}
		if strings.TrimSpace(out) != tc.expected {
			fmt.Printf("test %d failed: expected %s got %s\n", i+1, tc.expected, out)
		} else {
			passed++
		}
	}
	fmt.Printf("passed %d/%d tests\n", passed, len(tests))
	if passed != len(tests) {
		os.Exit(1)
	}
}

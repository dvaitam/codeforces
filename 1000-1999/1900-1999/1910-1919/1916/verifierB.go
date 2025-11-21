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

type testCase struct {
	input string
	pairs [][2]int64
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTests()
	for idx, tc := range tests {
		out, err := run(bin, tc.input)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\ninput:\n%s\n", idx+1, err, tc.input)
			os.Exit(1)
		}
		if err := check(tc, strings.TrimSpace(out)); err != nil {
			fmt.Printf("test %d failed: %v\ninput:\n%s\noutput:\n%s\n", idx+1, err, tc.input, out)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}

func check(tc testCase, out string) error {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != len(tc.pairs) {
		return fmt.Errorf("expected %d answers but got %d", len(tc.pairs), len(lines))
	}
	for i, line := range lines {
		x, err := strconv.ParseInt(strings.TrimSpace(line), 10, 64)
		if err != nil {
			return fmt.Errorf("test %d: output is not integer: %v", i+1, err)
		}
		a := tc.pairs[i][0]
		b := tc.pairs[i][1]
		if !validX(x, a, b) {
			return fmt.Errorf("test %d: x=%d does not have a=%d and b=%d as two largest divisors", i+1, x, a, b)
		}
	}
	return nil
}

func validX(x, a, b int64) bool {
	if x <= b || x > 1_000_000_000 {
		return false
	}
	divs := make([]int64, 0)
	for d := int64(1); d*d <= x; d++ {
		if x%d == 0 {
			divs = append(divs, d)
			if d != x/d {
				divs = append(divs, x/d)
			}
		}
	}
	var max1, max2, max3 int64
	for _, d := range divs {
		if d == x {
			continue
		}
		if d > max1 {
			max3 = max2
			max2 = max1
			max1 = d
		} else if d > max2 {
			max3 = max2
			max2 = d
		} else if d > max3 {
			max3 = d
		}
	}
	return max2 == b && max1 == a
}

func genTests() []testCase {
	rand.Seed(42)
	tests := []testCase{
		makeCase([][2]int64{{2, 3}, {5, 10}, {3, 11}}),
	}
	for i := 0; i < 100; i++ {
		n := rand.Intn(5) + 1
		pairs := make([][2]int64, n)
		for j := 0; j < n; j++ {
			x := rand.Int63n(1_000_000_000-10) + 10
			divs := divisors(x)
			if len(divs) < 3 {
				j--
				continue
			}
			pairs[j] = [2]int64{divs[len(divs)-3], divs[len(divs)-2]}
		}
		tests = append(tests, makeCase(pairs))
	}
	return tests
}

func divisors(x int64) []int64 {
	res := make([]int64, 0)
	for d := int64(1); d*d <= x; d++ {
		if x%d == 0 {
			res = append(res, d)
			if d != x/d {
				res = append(res, x/d)
			}
		}
	}
	sortSlice(res)
	return res
}

func sortSlice(a []int64) {
	for i := 0; i < len(a); i++ {
		for j := i + 1; j < len(a); j++ {
			if a[i] > a[j] {
				a[i], a[j] = a[j], a[i]
			}
		}
	}
}

func makeCase(pairs [][2]int64) testCase {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(pairs))
	for _, p := range pairs {
		fmt.Fprintf(&sb, "%d %d\n", p[0], p[1])
	}
	return testCase{
		input: sb.String(),
		pairs: pairs,
	}
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

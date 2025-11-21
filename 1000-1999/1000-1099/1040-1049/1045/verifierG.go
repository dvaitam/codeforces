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

type robot struct {
	x  int64
	r  int64
	iq int
}

type testCase struct {
	input  string
	expect string
	robots []robot
	k      int
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
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
		if err := checkOutput(strings.TrimSpace(out), tc); err != nil {
			fmt.Printf("test %d failed: %v\ninput:\n%s\nexpected: %s\nactual: %s\n", idx+1, err, tc.input, tc.expect, strings.TrimSpace(out))
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}

func checkOutput(out string, tc testCase) error {
	if out == "" {
		return fmt.Errorf("empty output")
	}
	val, err := strconv.ParseInt(out, 10, 64)
	if err != nil {
		return fmt.Errorf("output is not an integer: %v", err)
	}
	exp, _ := strconv.ParseInt(tc.expect, 10, 64)
	if val != exp {
		return fmt.Errorf("expected %d but got %d", exp, val)
	}
	return nil
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

func genTests() []testCase {
	rand.Seed(42)
	tests := []testCase{
		makeTest([]robot{{0, 0, 0}}, 0),
		makeTest([]robot{{0, 1, 0}, {1, 1, 0}}, 0),
		makeTest([]robot{{0, 5, 1}, {10, 5, 3}, {5, 5, 2}}, 1),
	}
	for i := 0; i < 200; i++ {
		n := rand.Intn(50) + 1
		k := rand.Intn(5)
		rbs := make([]robot, n)
		for j := 0; j < n; j++ {
			rbs[j] = robot{
				x:  int64(rand.Intn(1000)),
				r:  int64(rand.Intn(100)),
				iq: rand.Intn(30),
			}
		}
		tests = append(tests, makeTest(rbs, k))
	}
	return tests
}

func makeTest(robots []robot, k int) testCase {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", len(robots), k)
	for _, r := range robots {
		fmt.Fprintf(&sb, "%d %d %d\n", r.x, r.r, r.iq)
	}
	input := sb.String()
	expect := solveRef(robots, k)
	return testCase{
		input:  input,
		expect: fmt.Sprintf("%d", expect),
		robots: robots,
		k:      k,
	}
}

func solveRef(robots []robot, k int) int64 {
	var ans int64
	n := len(robots)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if canTalk(robots[i], robots[j], k) {
				ans++
			}
			if canTalk(robots[j], robots[i], k) {
				ans++
			}
		}
	}
	return ans
}

func canTalk(a, b robot, k int) bool {
	if absInt(a.iq-b.iq) > k {
		return false
	}
	return abs64(a.x-b.x) <= a.r
}

func absInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func abs64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

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
)

type testCase struct {
	input string
	arrs  [][]int
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA1.go /path/to/binary")
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
			fmt.Printf("test %d failed: %v\ninput:\n%s\noutput:\n%s\n", idx+1, err, tc.input, out)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}

func checkOutput(out string, tc testCase) error {
	lines := strings.Fields(out)
	if len(lines) == 0 {
		return fmt.Errorf("empty output")
	}
	reader := bufio.NewReader(strings.NewReader(out))
	var firstToken string
	fmt.Fscan(reader, &firstToken)
	if firstToken == "-1" {
		if hasSolution(tc.arrs[0]) {
			return fmt.Errorf("reported -1 but solution exists")
		}
		return nil
	}
	k, err := strconv.Atoi(firstToken)
	if err != nil || k <= 0 {
		return fmt.Errorf("invalid segments count %q", firstToken)
	}
	segs := make([][2]int, k)
	for i := 0; i < k; i++ {
		if _, err := fmt.Fscan(reader, &segs[i][0], &segs[i][1]); err != nil {
			return fmt.Errorf("failed to read segment %d: %v", i+1, err)
		}
	}
	return validateSegments(tc.arrs[0], segs)
}

func validateSegments(arr []int, segs [][2]int) error {
	n := len(arr)
	if len(segs) == 0 {
		return fmt.Errorf("no segments provided")
	}
	if segs[0][0] != 1 || segs[len(segs)-1][1] != n {
		return fmt.Errorf("segments must cover full array")
	}
	prevR := 0
	total := 0
	for _, seg := range segs {
		l, r := seg[0], seg[1]
		if l < 1 || r > n || l > r {
			return fmt.Errorf("invalid segment [%d, %d]", l, r)
		}
		if l != prevR+1 {
			return fmt.Errorf("segments not contiguous: expected %d got %d", prevR+1, l)
		}
		prevR = r
		sum := 0
		sign := 1
		for i := l - 1; i < r; i++ {
			sum += sign * arr[i]
			sign *= -1
		}
		total += sum
	}
	if total != 0 {
		return fmt.Errorf("total alternating sum %d != 0", total)
	}
	return nil
}

func genTests() []testCase {
	rand.Seed(42)
	tests := []testCase{
		makeTest([]int{1, -1, 1, -1}),
		makeTest([]int{-1, -1, 1, 1, -1, 1}),
		makeTest([]int{1}),
	}
	for i := 0; i < 100; i++ {
		n := rand.Intn(10) + 1
		arr := make([]int, n)
		for j := range arr {
			if rand.Intn(2) == 0 {
				arr[j] = 1
			} else {
				arr[j] = -1
			}
		}
		tests = append(tests, makeTest(arr))
	}
	return tests
}

func makeTest(arr []int) testCase {
	var sb strings.Builder
	fmt.Fprintf(&sb, "1\n%d\n", len(arr))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	return testCase{
		input: sb.String(),
		arrs:  [][]int{arr},
	}
}

func hasSolution(arr []int) bool {
	if len(arr)%2 == 1 {
		return false
	}
	return true
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

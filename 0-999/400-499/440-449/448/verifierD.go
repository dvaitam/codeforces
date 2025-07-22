package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type testCase struct {
	input    string
	expected int64
}

func countLE(x, n, m int64) int64 {
	var cnt int64
	for i := int64(1); i <= n; i++ {
		t := x / i
		if t > m {
			t = m
		}
		cnt += t
	}
	return cnt
}

func expectedD(n, m, k int64) int64 {
	if n > m {
		n, m = m, n
	}
	left, right := int64(1), n*m
	var ans int64
	for left <= right {
		mid := (left + right) / 2
		cnt := countLE(mid, n, m)
		if cnt >= k {
			ans = mid
			right = mid - 1
		} else {
			left = mid + 1
		}
	}
	return ans
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]testCase, 100)
	for i := 0; i < 100; i++ {
		n := int64(rng.Intn(20) + 1)
		m := int64(rng.Intn(20) + 1)
		k := int64(rng.Intn(int(n*m)) + 1)
		input := fmt.Sprintf("%d %d %d\n", n, m, k)
		cases[i] = testCase{input: input, expected: expectedD(n, m, k)}
	}
	return cases
}

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
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := generateTests()
	for i, tc := range cases {
		out, err := run(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
		var got int64
		fmt.Sscan(out, &got)
		if got != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\ninput:\n%s", i+1, tc.expected, got, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

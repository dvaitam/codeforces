package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type testCase struct {
	n   int
	arr []int
}

func generateCases() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]testCase, 100)
	for i := range cases {
		n := rng.Intn(5) + 2 // n between 2 and 6
		arr := make([]int, 2*n)
		for j := range arr {
			arr[j] = rng.Intn(1000) + 1
		}
		cases[i] = testCase{n: n, arr: arr}
	}
	return cases
}

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return out.String(), nil
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func check(tc testCase, output string) error {
	fields := strings.Fields(output)
	if len(fields) != 2*(tc.n-1) {
		return fmt.Errorf("expected %d numbers got %d", 2*(tc.n-1), len(fields))
	}
	used := make([]bool, 2*tc.n+1)
	g := 0
	for i := 0; i < tc.n-1; i++ {
		a, err := strconv.Atoi(fields[2*i])
		if err != nil {
			return fmt.Errorf("invalid integer %s", fields[2*i])
		}
		b, err := strconv.Atoi(fields[2*i+1])
		if err != nil {
			return fmt.Errorf("invalid integer %s", fields[2*i+1])
		}
		if a < 1 || a > 2*tc.n || b < 1 || b > 2*tc.n || a == b || used[a] || used[b] {
			return fmt.Errorf("invalid indices")
		}
		used[a] = true
		used[b] = true
		sum := tc.arr[a-1] + tc.arr[b-1]
		if g == 0 {
			g = sum
		} else {
			g = gcd(g, sum)
		}
	}
	if g <= 1 {
		return fmt.Errorf("gcd <= 1")
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := generateCases()
	for i, tc := range cases {
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d\n", tc.n))
		for j, v := range tc.arr {
			sb.WriteString(fmt.Sprintf("%d", v))
			if j+1 < len(tc.arr) {
				sb.WriteByte(' ')
			}
		}
		sb.WriteByte('\n')
		out, err := runCandidate(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if err := check(tc, out); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

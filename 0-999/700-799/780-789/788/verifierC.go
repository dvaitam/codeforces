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

func runCandidate(bin, input string) (string, error) {
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

func gcd(a, b int) int {
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

type testCase struct {
	n     int
	types []int
}

func solveC(tc testCase) int {
	n := tc.n
	a := tc.types
	seen := make([]bool, 1001)
	less := []int{}
	greater := []int{}
	minA, maxA := 1001, -1
	for _, v := range a {
		if seen[v] {
			continue
		}
		seen[v] = true
		if v == n {
			return 1
		}
		if v < n {
			less = append(less, v)
		} else {
			greater = append(greater, v)
		}
		if v < minA {
			minA = v
		}
		if v > maxA {
			maxA = v
		}
	}
	if n < minA || n > maxA || len(less) == 0 || len(greater) == 0 {
		return -1
	}
	ans := int(^uint(0) >> 1)
	for _, ai := range less {
		for _, aj := range greater {
			g := gcd(n-ai, aj-n)
			L := (aj - ai) / g
			if L < ans {
				ans = L
			}
		}
	}
	return ans
}

func genCase(rng *rand.Rand) testCase {
	n := rng.Intn(1001)
	k := rng.Intn(10) + 1
	types := make([]int, k)
	for i := range types {
		types[i] = rng.Intn(1001)
	}
	return testCase{n: n, types: types}
}

func formatInput(tc testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", tc.n, len(tc.types))
	for i, v := range tc.types {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	return sb.String()
}

func expected(tc testCase) string {
	return fmt.Sprintf("%d", solveC(tc))
}

func runCase(bin string, tc testCase) error {
	input := formatInput(tc)
	exp := expected(tc)
	out, err := runCandidate(bin, input)
	if err != nil {
		return err
	}
	if out != strings.TrimSpace(exp) {
		return fmt.Errorf("expected %s got %s\ninput:\n%s", exp, out, input)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var cases []testCase
	cases = append(cases, testCase{n: 0, types: []int{0}})
	cases = append(cases, testCase{n: 500, types: []int{300, 700}})
	for len(cases) < 100 {
		cases = append(cases, genCase(rng))
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

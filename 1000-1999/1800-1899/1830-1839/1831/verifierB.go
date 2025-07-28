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

type testB struct {
	n int
	a []int
	b []int
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
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("runtime error: %v", err)
	}
	return out.String(), nil
}

func solve(tc testB) int {
	cntA := make(map[int]int)
	for i := 0; i < tc.n; {
		j := i
		for j < tc.n && tc.a[j] == tc.a[i] {
			j++
		}
		if l := j - i; l > cntA[tc.a[i]] {
			cntA[tc.a[i]] = l
		}
		i = j
	}
	cntB := make(map[int]int)
	for i := 0; i < tc.n; {
		j := i
		for j < tc.n && tc.b[j] == tc.b[i] {
			j++
		}
		if l := j - i; l > cntB[tc.b[i]] {
			cntB[tc.b[i]] = l
		}
		i = j
	}
	ans := 0
	for v, la := range cntA {
		if la+cntB[v] > ans {
			ans = la + cntB[v]
		}
	}
	for v, lb := range cntB {
		if lb > ans && cntA[v] == 0 {
			ans = lb
		}
	}
	return ans
}

func genTests() []testB {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := []testB{
		{n: 1, a: []int{1}, b: []int{1}},
		{n: 2, a: []int{1, 1}, b: []int{1, 1}},
		{n: 3, a: []int{1, 2, 1}, b: []int{1, 1, 1}},
	}
	for len(tests) < 100 {
		n := rng.Intn(20) + 1
		a := make([]int, n)
		b := make([]int, n)
		for i := 0; i < n; i++ {
			a[i] = rng.Intn(2*n) + 1
		}
		for i := 0; i < n; i++ {
			b[i] = rng.Intn(2*n) + 1
		}
		tests = append(tests, testB{n: n, a: a, b: b})
	}
	return tests
}

func runCase(bin string, tc testB) error {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(strconv.Itoa(tc.n))
	sb.WriteByte('\n')
	for i, v := range tc.a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	for i, v := range tc.b {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')

	expected := solve(tc)
	out, err := runCandidate(bin, sb.String())
	if err != nil {
		return fmt.Errorf("%v\ninput:\n%s", err, sb.String())
	}
	fields := strings.Fields(strings.TrimSpace(out))
	if len(fields) != 1 {
		return fmt.Errorf("expected single integer got %q", out)
	}
	val, err := strconv.Atoi(fields[0])
	if err != nil {
		return fmt.Errorf("cannot parse integer")
	}
	if val != expected {
		return fmt.Errorf("expected %d got %d\ninput:\n%s", expected, val, sb.String())
	}
	return nil
}

func main() {
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := args[0]
	tests := genTests()
	for i, tc := range tests {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

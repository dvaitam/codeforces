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

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return out.String() + errBuf.String(), fmt.Errorf("runtime error: %v", err)
	}
	return out.String(), nil
}

func solveB(n, k int, p []int) int {
	expected := 1
	for _, v := range p {
		if v == expected {
			expected++
		}
	}
	m := expected - 1
	remaining := n - m
	if remaining <= 0 {
		return 0
	}
	return (remaining + k - 1) / k
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	type test struct {
		n, k int
		p    []int
	}

	var cases []test
	// deterministic cases
	cases = append(cases, test{n: 2, k: 1, p: []int{1, 2}})
	cases = append(cases, test{n: 3, k: 1, p: []int{2, 1, 3}})
	cases = append(cases, test{n: 5, k: 2, p: []int{2, 5, 1, 3, 4}})

	for len(cases) < 100 {
		n := rng.Intn(10) + 2
		k := rng.Intn(n) + 1
		perm := rand.Perm(n)
		for i := range perm {
			perm[i]++
		}
		cases = append(cases, test{n: n, k: k, p: perm})
	}

	for i, tc := range cases {
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.k))
		for j, v := range tc.p {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		expected := solveB(tc.n, tc.k, tc.p)
		out, err := runCandidate(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\ninput:\n%s", i+1, err, sb.String())
			os.Exit(1)
		}
		fields := strings.Fields(strings.TrimSpace(out))
		if len(fields) != 1 {
			fmt.Fprintf(os.Stderr, "case %d: expected single integer got %q\n", i+1, out)
			os.Exit(1)
		}
		val, err := strconv.Atoi(fields[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: cannot parse integer\n", i+1)
			os.Exit(1)
		}
		if val != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\ninput:\n%s", i+1, expected, val, sb.String())
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

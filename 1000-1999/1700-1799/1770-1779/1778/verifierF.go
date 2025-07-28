package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testCaseF struct {
	n, k  int
	a     []int
	edges [][2]int
}

func parseCasesF(path string) ([]testCaseF, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	in := bufio.NewReader(f)
	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return nil, err
	}
	cases := make([]testCaseF, t)
	for i := 0; i < t; i++ {
		var n, k int
		if _, err := fmt.Fscan(in, &n, &k); err != nil {
			return nil, err
		}
		a := make([]int, n)
		for j := 0; j < n; j++ {
			if _, err := fmt.Fscan(in, &a[j]); err != nil {
				return nil, err
			}
		}
		edges := make([][2]int, n-1)
		for j := 0; j < n-1; j++ {
			var u, v int
			if _, err := fmt.Fscan(in, &u, &v); err != nil {
				return nil, err
			}
			edges[j] = [2]int{u, v}
		}
		cases[i] = testCaseF{n: n, k: k, a: a, edges: edges}
	}
	return cases, nil
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func solveF(tc testCaseF) int {
	g := tc.a[0]
	for i := 1; i < tc.n; i++ {
		g = gcd(g, tc.a[i])
	}
	ans := tc.a[0]
	if tc.k > 0 {
		ans *= g
	}
	return ans
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
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases, err := parseCasesF("testcasesF.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		sb := strings.Builder{}
		fmt.Fprintf(&sb, "1\n%d %d\n", tc.n, tc.k)
		for i, v := range tc.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		for _, e := range tc.edges {
			fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
		}
		input := sb.String()
		expected := solveF(tc)
		gotStr, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		got, err := strconv.Atoi(strings.TrimSpace(gotStr))
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: cannot parse output: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\n", idx+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	oracle := filepath.Join(dir, "oracleE")
	cmd := exec.Command("go", "build", "-o", oracle, "981E.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
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
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func formatCase(n int, ops [][3]int) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, len(ops)))
	for _, op := range ops {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", op[0], op[1], op[2]))
	}
	return sb.String()
}

func addHandcraftedTests(cases *[]string) {
	*cases = append(*cases,
		formatCase(1, [][3]int{{1, 1, 1}}),
		formatCase(5, [][3]int{{1, 5, 3}}),
		formatCase(5, [][3]int{{1, 5, 2}, {2, 5, 2}, {1, 3, 1}}),
		formatCase(6, [][3]int{{1, 2, 1}, {2, 3, 1}, {3, 4, 1}, {4, 5, 1}, {5, 6, 1}}),
		formatCase(14, [][3]int{{8, 8, 1}, {8, 9, 6}, {3, 9, 6}, {2, 6, 11}, {12, 12, 14}, {2, 7, 7}, {7, 8, 6}, {3, 4, 6}, {4, 4, 3}, {13, 13, 14}, {11, 12, 8}, {12, 12, 1}, {12, 12, 3}, {5, 7, 9}, {7, 8, 4}, {6, 6, 10}, {12, 13, 11}, {14, 14, 13}, {6, 9, 3}, {5, 9, 11}, {9, 13, 5}, {6, 9, 10}, {13, 14, 5}, {10, 10, 4}, {6, 11, 7}, {13, 13, 7}, {10, 14, 9}, {13, 14, 3}, {7, 8, 11}, {6, 13, 10}, {13, 14, 7}, {10, 14, 1}, {3, 7, 5}, {3, 3, 8}, {2, 5, 6}}),
	)
}

func bruteExpected(input string) (string, error) {
	fields := strings.Fields(input)
	if len(fields) < 2 {
		return "", fmt.Errorf("bad input: missing n and q")
	}
	n, err := strconv.Atoi(fields[0])
	if err != nil {
		return "", fmt.Errorf("bad n: %w", err)
	}
	q, err := strconv.Atoi(fields[1])
	if err != nil {
		return "", fmt.Errorf("bad q: %w", err)
	}
	need := 2 + 3*q
	if len(fields) < need {
		return "", fmt.Errorf("bad input: expected %d fields, got %d", need, len(fields))
	}
	type op struct{ l, r, v int }
	ops := make([]op, q)
	idx := 2
	for i := 0; i < q; i++ {
		l, err := strconv.Atoi(fields[idx])
		if err != nil {
			return "", fmt.Errorf("bad l at op %d: %w", i+1, err)
		}
		r, err := strconv.Atoi(fields[idx+1])
		if err != nil {
			return "", fmt.Errorf("bad r at op %d: %w", i+1, err)
		}
		v, err := strconv.Atoi(fields[idx+2])
		if err != nil {
			return "", fmt.Errorf("bad v at op %d: %w", i+1, err)
		}
		ops[i] = op{l, r, v}
		idx += 3
	}

	possible := make([]bool, n+1)
	for pos := 1; pos <= n; pos++ {
		dp := make([]bool, n+1)
		dp[0] = true
		for _, op := range ops {
			if pos < op.l || pos > op.r || op.v > n {
				continue
			}
			for s := n; s >= op.v; s-- {
				if dp[s-op.v] {
					dp[s] = true
				}
			}
		}
		for s := 1; s <= n; s++ {
			if dp[s] {
				possible[s] = true
			}
		}
	}

	ans := make([]int, 0, n)
	for s := 1; s <= n; s++ {
		if possible[s] {
			ans = append(ans, s)
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(ans)))
	for i, v := range ans {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	return strings.TrimSpace(sb.String()), nil
}

func addRandomTests(cases *[]string) {
	rng := rand.New(rand.NewSource(981))
	for t := 0; t < 300; t++ {
		n := rng.Intn(20) + 1
		q := rng.Intn(35) + 1
		ops := make([][3]int, q)
		for i := 0; i < q; i++ {
			l := rng.Intn(n) + 1
			r := rng.Intn(n-l+1) + l
			v := rng.Intn(n) + 1
			ops[i] = [3]int{l, r, v}
		}
		*cases = append(*cases, formatCase(n, ops))
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	cases := make([]string, 0, 304)
	addHandcraftedTests(&cases)
	addRandomTests(&cases)

	for idx, input := range cases {
		exp, err := run(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle error on case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		brute, err := bruteExpected(input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "brute-force error on case %d: %v\ninput:\n%s\n", idx+1, err, input)
			os.Exit(1)
		}
		if exp != brute {
			fmt.Fprintf(os.Stderr, "oracle mismatch on case %d\ninput:\n%sexpected by brute:\n%s\ngoracle output:\n%s\n", idx+1, input, brute, exp)
			os.Exit(1)
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s\n", idx+1, err, input)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", idx+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d generated tests passed\n", len(cases))
}

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

func buildOracle() (string, error) {
	exe := "oracleC"
	cmd := exec.Command("go", "build", "-o", exe, "1814C.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle: %v\n%s", err, out)
	}
	return exe, nil
}

func generateCase(rng *rand.Rand) string {
	t := rng.Intn(3) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", t)
	for i := 0; i < t; i++ {
		n := rng.Intn(6) + 2
		s1 := rng.Intn(5) + 1
		s2 := rng.Intn(5) + 1
		fmt.Fprintf(&sb, "%d %d %d\n", n, s1, s2)
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", rng.Intn(20)+1)
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func runProg(exe, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(exe, ".go") {
		cmd = exec.Command("go", "run", exe)
	} else {
		cmd = exec.Command(exe)
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

type testCase struct {
	n, s1, s2 int
	r         []int // 1-indexed
}

func parseInput(input string) []testCase {
	fields := strings.Fields(input)
	idx := 0
	t, _ := strconv.Atoi(fields[idx])
	idx++
	cases := make([]testCase, t)
	for i := 0; i < t; i++ {
		n, _ := strconv.Atoi(fields[idx])
		idx++
		s1, _ := strconv.Atoi(fields[idx])
		idx++
		s2, _ := strconv.Atoi(fields[idx])
		idx++
		r := make([]int, n+1)
		for j := 1; j <= n; j++ {
			r[j], _ = strconv.Atoi(fields[idx])
			idx++
		}
		cases[i] = testCase{n, s1, s2, r}
	}
	return cases
}

// computeCost returns the total search time for the given assignment.
// listA and listB contain 1-indexed box numbers in list order.
func computeCost(tc testCase, listA, listB []int) int64 {
	var cost int64
	for j, box := range listA {
		cost += int64(j+1) * int64(tc.s1) * int64(tc.r[box])
	}
	for j, box := range listB {
		cost += int64(j+1) * int64(tc.s2) * int64(tc.r[box])
	}
	return cost
}

// parseLists reads 2*len(cases) lines of output and returns (listAs, listBs, error).
func parseLists(output string, cases []testCase) ([][]int, [][]int, error) {
	fields := strings.Fields(output)
	idx := 0
	listAs := make([][]int, len(cases))
	listBs := make([][]int, len(cases))
	for i, tc := range cases {
		if idx >= len(fields) {
			return nil, nil, fmt.Errorf("test %d: unexpected end of output", i+1)
		}
		ka, _ := strconv.Atoi(fields[idx])
		idx++
		listA := make([]int, ka)
		for j := 0; j < ka; j++ {
			if idx >= len(fields) {
				return nil, nil, fmt.Errorf("test %d: unexpected end in list A", i+1)
			}
			listA[j], _ = strconv.Atoi(fields[idx])
			idx++
		}
		if idx >= len(fields) {
			return nil, nil, fmt.Errorf("test %d: unexpected end before list B", i+1)
		}
		kb, _ := strconv.Atoi(fields[idx])
		idx++
		listB := make([]int, kb)
		for j := 0; j < kb; j++ {
			if idx >= len(fields) {
				return nil, nil, fmt.Errorf("test %d: unexpected end in list B", i+1)
			}
			listB[j], _ = strconv.Atoi(fields[idx])
			idx++
		}
		if ka+kb != tc.n {
			return nil, nil, fmt.Errorf("test %d: list sizes %d+%d != n=%d", i+1, ka, kb, tc.n)
		}
		// Validate partition.
		seen := make([]bool, tc.n+1)
		for _, b := range listA {
			if b < 1 || b > tc.n || seen[b] {
				return nil, nil, fmt.Errorf("test %d: invalid or duplicate box %d in list A", i+1, b)
			}
			seen[b] = true
		}
		for _, b := range listB {
			if b < 1 || b > tc.n || seen[b] {
				return nil, nil, fmt.Errorf("test %d: invalid or duplicate box %d in list B", i+1, b)
			}
			seen[b] = true
		}
		listAs[i] = listA
		listBs[i] = listB
	}
	return listAs, listBs, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := generateCase(rng)
		exp, err := runProg("./"+oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failure on case %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		got, err := runProg(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}

		cases := parseInput(input)

		expAs, expBs, err := parseLists(exp, cases)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle invalid output on case %d: %v\ninput:\n%s\noutput:\n%s\n", i+1, err, input, exp)
			os.Exit(1)
		}
		gotAs, gotBs, err := parseLists(got, cases)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: invalid output: %v\ninput:\n%s\ngot:\n%s\n", i+1, err, input, got)
			os.Exit(1)
		}

		for j, tc := range cases {
			expCost := computeCost(tc, expAs[j], expBs[j])
			gotCost := computeCost(tc, gotAs[j], gotBs[j])
			if expCost != gotCost {
				fmt.Fprintf(os.Stderr, "case %d test %d cost mismatch: expected %d, got %d\nexpected:\n%s\ngot:\n%s\ninput:\n%s",
					i+1, j+1, expCost, gotCost, exp, got, input)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}

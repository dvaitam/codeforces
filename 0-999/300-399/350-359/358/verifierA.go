package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesA.txt so the verifier is self-contained.
const testcasesA = `8 6 -18 -4 12 11 5 -1 10
7 17 -7 12 -12 -2 20 -14
6 14 18 -11 -1 -14 -16
7 10 15 -14 2 7 0 -7
10 10 8 13 -4 -17 15 -20 -15 5 14
9 1 -5 0 -16 -8 -6 19 -11 8
3 -15 0 12
9 -14 -1 15 -2 -13 18 1 -7 17
9 -15 18 4 0 16 -5 -2 -9 -8
4 -18 19 -4 10
3 -15 -12 -11
2 -15 14
8 13 -3 20 -5 -7 6 19 8
9 2 -15 0 -13 11 1 -8 -5 -19
6 -13 -6 3 -10 1 7
2 -14 -11
5 -18 16 20 14 18
3 -19 -13 20
5 18 16 -13 5 -15
7 -13 -18 18 -19 -8 -9 20
9 -7 -17 -19 14 7 -14 -4 -16 -6
3 -1 2 7
4 -17 12 9 -18
3 5 -8 -4
7 10 16 -10 -7 -17 18 15
7 13 -4 -13 8 -9 -20 10
8 16 12 -1 2 4 -4 -11 -20
9 -15 1 -18 14 -3 -12 -5 10 2
6 2 17 -12 -1 4 6
3 -20 18 -8
7 -10 -5 -6 8 4 6 -18
8 16 6 -18 -10 8 -16 -4 17
9 13 11 15 -20 -18 19 0 -1 9
2 6 -8
10 20 -15 -12 -20 5 6 0 17 -7 13
2 13 19
3 -8 -13 18
5 -1 -3 -9 -14 10
8 20 -15 -19 -3 8 -13 -4 -12
10 2 -13 -11 -3 -19 -18 15 -7 -4 0
7 16 -18 18 11 9 7 3
10 -9 -7 4 17 -2 -20 -12 -11 -3 1
7 3 -15 1 -18 17 -3 -10
4 17 -2 3 5
10 -12 -2 -13 10 -5 -17 -1 -9 -16 14
8 1 -1 6 -14 17 15 10 14
7 1 -13 10 19 11 7 -18
6 1 -11 -10 16 4 -15
3 -15 -8 -6
2 4 -20
3 5 15 13
6 8 11 17 -7 7 -15
7 -6 -4 17 -10 7 -8 2
3 -16 -19 13
9 -8 -13 11 5 -4 -7 -18 15 -11
3 -8 9 4
7 14 -11 -14 11 19 5 7
10 11 0 20 18 -8 14 -6 -20 1 19
7 -18 13 -11 -4 18 4 -2
9 -16 -15 13 -18 20 -6 -12 17 -1
2 8 1
4 -11 9 3 12
8 13 12 -18 16 -15 20 -16 7
5 -2 14 18 6 10
8 18 17 -6 -19 -20 -9 -1 12
6 1 -16 11 -4 -1 6
8 4 -17 -10 -12 -5 -2 1 19
2 10 6
4 11 18 -15 -11
7 6 -18 9 4 18 -17 -14
9 -11 -19 -18 -12 0 -14 2 -8 4
9 -13 -17 9 1 20 -2 -12 4 15
3 13 -8 -18
8 8 3 -8 9 2 -16 -18 14
9 -4 -19 13 16 17 -7 -6 -15 12
10 6 12 -1 -13 -11 7 15 -15 -14 20
3 -14 6 -11
2 8 7
8 -19 11 0 -4 -15 2 -16 -13
7 -19 2 19 -9 -20 -6 3
3 18 -11 -7
2 -7 -13
2 -2 3
2 18 -6
4 -9 9 -13 10
7 -4 -12 -19 -7 3 1 10
6 -2 15 0 -9 -15 -14
10 17 -1 -10 4 -11 -12 -6 0 12 -5
5 -9 -2 3 6 -18
4 18 -19 5 -16
3 -12 6 -1
10 6 -11 17 7 -1 2 -15 -5 8 3
10 -17 4 6 -20 18 0 8 -7 3 -2
9 -15 -9 -14 -3 -13 15 -11 8 5
4 6 7 -9 -5
9 1 13 -11 2 9 -15 10 -7 -2
2 8 19
9 -20 -7 -1 -13 18 14 -11 7 10
3 11 -6 14
8 -3 -19 -13 20 -18 -20 -4 5`

func runCandidate(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func expected(xs []int) string {
	type seg struct{ a, b int }
	segs := make([]seg, 0, len(xs)-1)
	for i := 0; i+1 < len(xs); i++ {
		a, b := xs[i], xs[i+1]
		if a > b {
			a, b = b, a
		}
		segs = append(segs, seg{a, b})
	}
	for i := 0; i < len(segs); i++ {
		a1, b1 := segs[i].a, segs[i].b
		for j := i + 1; j < len(segs); j++ {
			a2, b2 := segs[j].a, segs[j].b
			if (a1 < a2 && a2 < b1 && b1 < b2) || (a2 < a1 && a1 < b2 && b2 < b1) {
				return "yes"
			}
		}
	}
	return "no"
}

func parseCases() ([]([]int), error) {
	lines := strings.Split(strings.TrimSpace(testcasesA), "\n")
	cases := make([][]int, 0, len(lines))
	for idx, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		fields := strings.Fields(line)
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: bad n", idx+1)
		}
		if len(fields) != n+1 {
			return nil, fmt.Errorf("line %d: expected %d numbers got %d", idx+1, n, len(fields)-1)
		}
		xs := make([]int, n)
		for i := 0; i < n; i++ {
			v, err := strconv.Atoi(fields[i+1])
			if err != nil {
				return nil, fmt.Errorf("line %d: bad value %d", idx+1, i+1)
			}
			xs[i] = v
		}
		cases = append(cases, xs)
	}
	return cases, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: verifierA /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases, err := parseCases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, xs := range cases {
		want := expected(xs)
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", len(xs))
		for i, v := range xs {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		got, err := runCandidate(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

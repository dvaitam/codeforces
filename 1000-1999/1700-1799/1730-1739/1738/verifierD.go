package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testCase struct {
	input     string
	expectedK int
	expectedA []int
}

const testcaseData = `4 4 4 1 2
8 9 1 9 0 7 4 8 3
4 5 3 4 4
8 6 2 3 2 8 6 0 1
3 4 0 2
1 1
8 9 6 6 6 9 7 2 5
2 0 1
8 3 4 6 4 6 8 6 9
6 6 3 5 0 4 2
6 1 3 4 4 1 1
8 7 1 5 1 6 2 0 4
7 6 1 0 0 6 5 8
5 4 1 0 2 0
2 0 0
4 3 2 4 2
3 0 2 2
6 2 6 6 7 6 1
5 3 5 5 5 1
5 3 2 4 2 4
6 0 6 5 0 6 2
1 2
6 7 5 5 4 7 0
1 2
1 1
5 5 3 2 4 4
6 2 5 2 5 5 4
5 6 3 0 6 6
1 2
3 2 4 1
5 1 2 1 5 3
2 0 2
6 3 7 2 1 5 3
8 4 3 1 0 8 3 5 9
3 2 2 0
6 2 6 4 4 7 5
7 4 6 6 0 6 2 3
1 1
7 8 3 0 7 8 4 8
6 3 1 4 1 3 0
1 2
4 3 4 0 0
8 1 2 8 4 3 0 8 8
7 0 1 5 2 4 8 7
1 1
4 1 0 4 0
3 1 2 1
1 1
7 0 4 3 4 8 8 6
1 1
6 0 0 2 0 1 0
2 3 0
2 3 2
3 2 0 2
7 6 4 5 4 3 5 6
2 1 0
7 1 2 0 5 7 8 6
1 2
7 0 5 7 5 6 6 7
1 0
4 4 2 5 4
2 3 1
7 2 0 5 5 8 4 1
8 1 8 6 1 5 9 8 1
1 1
3 1 3 0
2 0 3
3 0 2 0
1 0
8 4 9 4 1 0 9 8 8
4 0 4 5 0
1 2
6 2 1 3 2 3 7
7 4 5 6 5 8 6 1
7 8 3 6 2 6 8 7
3 3 1 1
2 3 3
8 9 2 2 4 3 2 9 8
6 3 4 6 4 3 4
1 1
8 6 3 2 9 5 3 5 7
3 3 3 4
4 3 4 5 4
1 1
2 3 0
8 3 3 1 3 4 3 3 4
3 1 4 0
5 1 6 0 2 1
7 1 1 1 1 4 4 0
6 7 5 0 0 5 5
7 6 7 1 3 7 6 2
6 1 4 1 6 1 7
5 0 4 5 2 5
6 7 4 4 1 5 1
8 8 5 0 4 9 2 2 2
6 7 1 1 2 5 6
5 5 1 3 3 2
3 0 0 1
7 5 1 4 4 6 0 2
1 1`

func solveD(n int, b []int) (int, []int) {
	v := make([][]int, n+2)
	k := 0
	for i := 1; i <= n; i++ {
		if b[i] > i {
			k++
		}
		if b[i] >= 0 && b[i] <= n+1 {
			v[b[i]] = append(v[b[i]], i)
		}
	}
	a := make([]int, 0, n)
	cur := 0
	if len(v[n+1]) > 0 {
		cur = n + 1
	}
	cnt := 0
	for cnt < n {
		cnt += len(v[cur])
		last := len(v[cur]) - 1
		good := -1
		for j := 0; j <= last; j++ {
			nxt := v[cur][j]
			if nxt >= 0 && nxt < len(v) && len(v[nxt]) > 0 {
				good = j
			}
		}
		if good != -1 && good != last {
			v[cur][good], v[cur][last] = v[cur][last], v[cur][good]
		}
		a = append(a, v[cur]...)
		if len(v[cur]) > 0 {
			cur = v[cur][len(v[cur])-1]
		} else {
			break
		}
	}
	return k, a
}

func parseTestcase(line string) (testCase, error) {
	fields := strings.Fields(line)
	if len(fields) < 1 {
		return testCase{}, fmt.Errorf("empty testcase")
	}
	n, err := strconv.Atoi(fields[0])
	if err != nil {
		return testCase{}, fmt.Errorf("bad n: %w", err)
	}
	if len(fields) != n+1 {
		return testCase{}, fmt.Errorf("expected %d numbers, got %d", n, len(fields)-1)
	}
	bvals := make([]int, n+1)
	for i := 0; i < n; i++ {
		val, err := strconv.Atoi(fields[1+i])
		if err != nil {
			return testCase{}, fmt.Errorf("bad b[%d]: %w", i+1, err)
		}
		bvals[i+1] = val
	}

	var input strings.Builder
	input.WriteString("1\n")
	input.WriteString(fmt.Sprintf("%d\n", n))
	for i := 1; i <= n; i++ {
		if i > 1 {
			input.WriteByte(' ')
		}
		input.WriteString(fmt.Sprintf("%d", bvals[i]))
	}
	input.WriteByte('\n')

	k, perm := solveD(n, bvals)
	return testCase{input: input.String(), expectedK: k, expectedA: perm}, nil
}

func loadTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcaseData), "\n")
	cases := make([]testCase, 0, len(lines))
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		tc, err := parseTestcase(line)
		if err != nil {
			return nil, fmt.Errorf("case %d: %w", idx+1, err)
		}
		cases = append(cases, tc)
	}
	return cases, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierD /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases, err := loadTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range cases {
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(tc.input)
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx+1, err, stderr.String())
			os.Exit(1)
		}
		rawLines := strings.Split(out.String(), "\n")
		lineIdx := 0
		for lineIdx < len(rawLines) && strings.TrimSpace(rawLines[lineIdx]) == "" {
			lineIdx++
		}
		if lineIdx == len(rawLines) {
			fmt.Printf("test %d: missing k line\n", idx+1)
			os.Exit(1)
		}
		kLine := strings.TrimSpace(rawLines[lineIdx])
		lineIdx++

		gotK, err := strconv.Atoi(kLine)
		if err != nil {
			fmt.Printf("test %d: cannot parse k: %v\n", idx+1, err)
			os.Exit(1)
		}

		gotA := make([]int, 0, len(tc.expectedA))
		if len(tc.expectedA) > 0 {
			for lineIdx < len(rawLines) && strings.TrimSpace(rawLines[lineIdx]) == "" {
				lineIdx++
			}
			if lineIdx == len(rawLines) {
				fmt.Printf("test %d: missing permutation line\n", idx+1)
				os.Exit(1)
			}
			gotParts := strings.Fields(strings.TrimSpace(rawLines[lineIdx]))
			if len(gotParts) != len(tc.expectedA) {
				fmt.Printf("test %d: expected %d numbers got %d\n", idx+1, len(tc.expectedA), len(gotParts))
				os.Exit(1)
			}
			for i := 0; i < len(gotParts); i++ {
				val, err := strconv.Atoi(gotParts[i])
				if err != nil {
					fmt.Printf("test %d: bad integer in output: %v\n", idx+1, err)
					os.Exit(1)
				}
				gotA = append(gotA, val)
			}
		}

		if gotK != tc.expectedK {
			fmt.Printf("test %d: expected k=%d got %d\n", idx+1, tc.expectedK, gotK)
			os.Exit(1)
		}
		if len(tc.expectedA) != 0 {
			for i := 0; i < len(gotA); i++ {
				if gotA[i] != tc.expectedA[i] {
					fmt.Printf("test %d: permutation mismatch\nexpected %v\ngot %v\n", idx+1, tc.expectedA, gotA)
					os.Exit(1)
				}
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

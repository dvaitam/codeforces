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
	n      int
	k      int
	blocks [][]int
}

// Embedded testcases from testcasesG.txt.
const testcaseData = `
4 1 3 4 2 2
2 1 1 2
5 4 3 4 2 4 4 1 5 4 3 4 5 5 5 5 5 1 5 1 4 4
1 2 1 1 1 1
2 3 2 2 2 2 2 1 2 2 1
5 2 4 4 1 5 1 2 5 5
2 4 2 2 2 1 1 2 1 2 1 2
3 2 3 3 1 1 1 3
4 3 1 4 2 4 4 1 2
3 5 3 1 1 1 1 3 1 2 2 2 3 2 2 3
2 1 1 1
1 5 1 1 1 1 1 1 1 1 1 1
4 4 4 1 1 2 3 2 1 2 2 1 3 1 2
1 5 1 1 1 1 1 1 1 1 1 1
4 5 3 1 1 3 3 3 4 1 4 1 3 4 2 1 3 3 4 2 2
5 1 1 5
5 1 4 5 3 4 5
3 3 3 3 3 2 1 3 3 2 3 3
4 4 2 2 3 2 1 1 1 4 1 2
5 4 4 1 3 3 5 2 1 2 3 3 1 3 3 5 5 3
1 1 1 1
4 5 1 2 2 3 3 4 4 1 2 4 2 4 3 4 4 1 3 3
5 5 1 2 5 3 3 1 2 3 1 3 4 5 2 5 5 4 2 5 1 3
3 1 2 2 1
4 2 3 4 3 2 3 2 1 3
1 2 1 1 1 1
5 2 2 3 4 5 3 3 2 1 5
2 2 1 2 2 1 1
3 1 3 3 1 2
4 1 1 3
1 5 1 1 1 1 1 1 1 1 1 1
3 2 1 2 1 1
5 5 5 1 2 5 3 5 3 3 2 3 5 2 2 1 3 5 3 5 4 1 2 1 4
1 1 1 1
1 4 1 1 1 1 1 1 1 1
4 5 2 1 4 1 4 3 4 4 2 2 2 3 3 4 2 3
3 4 2 1 1 2 1 2 2 2 1 2 2 3
5 2 5 1 4 1 3 5 1 3
5 4 4 5 4 2 5 2 5 2 4 4 1 3 2 5 4 4 5 3 5
2 2 1 2 1 2
3 5 1 2 1 2 2 1 3 1 3 3 2 2 2
1 1 1 1
5 1 2 2 4
5 3 4 5 5 4 2 2 3 3 5 3 4 3 1 5
1 5 1 1 1 1 1 1 1 1 1 1
1 2 1 1 1 1
4 4 4 2 2 4 2 2 4 3 4 2 3 2 2 1 1
1 2 1 1 1 1
2 4 2 2 1 2 2 2 2 2 1 1 2
5 1 5 4 2 5 3 4
5 5 3 1 1 5 4 3 1 5 4 1 2 5 3 4 5 1 5 5 1 2 4 1 5
5 3 5 1 3 5 1 4 4 5 4 5 3 4 3 3 3 2
1 5 1 1 1 1 1 1 1 1 1 1
1 2 1 1 1 1
5 3 3 3 1 1 5 2 5 2 4 3 2 3 2
4 2 1 3 1 2
5 2 4 4 4 2 4 2 2 5
3 1 3 2 1 1
2 4 2 2 2 2 1 2 2 2 2 1 2
1 3 1 1 1 1 1 1
4 5 3 3 3 3 1 3 4 4 2 2 3 2 3 1 1 1
1 4 1 1 1 1 1 1 1 1
1 1 1 1
5 2 5 1 1 1 5 2 4 5 2 1 2
3 3 2 3 1 3 2 2 2 3 3 2 2
2 1 2 1 2
4 3 1 3 3 2 4 1 4 2 4 2 4
1 2 1 1 1 1
3 1 1 1
1 5 1 1 1 1 1 1 1 1 1 1
3 4 3 1 3 3 2 2 2 2 3 1 3 1 1 2
4 1 3 3 3 4
1 1 1 1
1 4 1 1 1 1 1 1 1 1
1 3 1 1 1 1 1 1
2 1 1 2
4 1 3 4 4 1
1 3 1 1 1 1 1 1
5 4 5 2 5 2 1 1 2 4 1 4 5 1 2 2 1 2
1 3 1 1 1 1 1 1
4 5 2 1 2 4 3 4 3 1 2 3 2 1 2 3 1 1 3
2 4 1 2 1 2 1 1 1 1
5 4 2 1 4 1 1 1 2 1 5
3 4 3 2 3 1 1 1 3 2 3 2 2 1 1
5 1 5 2 4 5 1 3
3 4 1 3 1 1 1 3 1 3
3 5 3 1 2 3 3 1 2 1 1 1 1 3 2 2 1
2 1 1 2
3 2 2 2 2 3 2 1 1
1 5 1 1 1 1 1 1 1 1 1 1
2 3 1 1 2 1 1 1 2
2 4 1 1 2 2 1 2 1 1 2 1 2
3 2 2 3 1 3 2 2 1
3 2 1 3 3 1 2 1
1 4 1 1 1 1 1 1 1 1
5 5 2 1 4 5 4 2 5 4 4 3 1 4 2 2 3 4 5 2 5 2 5 1
4 1 1 4
5 1 4 5 4 4 1
4 4 2 1 1 4 2 2 3 3 2 1 1 4 3 2 4 1
4 4 2 4 2 3 4 1 1 4 3 1 3 4 3 2 1 3
`

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcaseData), "\n")
	var cases []testCase
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			return nil, fmt.Errorf("case %d too short", idx+1)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("case %d bad n: %v", idx+1, err)
		}
		k, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil, fmt.Errorf("case %d bad k: %v", idx+1, err)
		}
		pos := 2
		blocks := make([][]int, 0, k)
		for i := 0; i < k; i++ {
			if pos >= len(fields) {
				return nil, fmt.Errorf("case %d missing q for block %d", idx+1, i+1)
			}
			q, err := strconv.Atoi(fields[pos])
			if err != nil {
				return nil, fmt.Errorf("case %d bad q %d: %v", idx+1, i+1, err)
			}
			pos++
			if pos+q > len(fields) {
				return nil, fmt.Errorf("case %d not enough elements for block %d", idx+1, i+1)
			}
			block := make([]int, q)
			for j := 0; j < q; j++ {
				val, err := strconv.Atoi(fields[pos+j])
				if err != nil {
					return nil, fmt.Errorf("case %d bad value %d in block %d: %v", idx+1, j+1, i+1, err)
				}
				block[j] = val
			}
			blocks = append(blocks, block)
			pos += q
		}
		if pos != len(fields) {
			return nil, fmt.Errorf("case %d has leftover data", idx+1)
		}
		cases = append(cases, testCase{n: n, k: k, blocks: blocks})
	}
	return cases, nil
}

// solve mirrors 1552G.go placeholder.
func solve(tc testCase) string {
	return "REJECTED"
}

func runCandidate(bin string, tc testCase) (string, error) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.k))
	for _, block := range tc.blocks {
		sb.WriteString(strconv.Itoa(len(block)))
		for _, v := range block {
			sb.WriteByte(' ')
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
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
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for idx, tc := range tests {
		expect := solve(tc)
		got, err := runCandidate(bin, tc)
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("case %d failed: expected %s got %s\n", idx+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

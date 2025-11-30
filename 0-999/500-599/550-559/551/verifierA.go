package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testcase struct {
	n       int
	ratings []int
}

// Embedded testcases from testcasesA.txt.
const testcasesRaw = `7 14 2 9 17 16 13 10
8 12 19 7 17 5 10 5 4
10 9 18 20 5 10 4 3 11 16 18
2 12 14
6 20 7 18 16 15 17
5 2 18 1 3 13
1 20
8 11 8 11 3 7 19 8 8
3 18 15 3
2 11 17
8 4 10 18 10 4 18 11 18
4 20 18 19 10
8 3 20 13 11 19 8 10 6
4 6 2 20 9
8 3 3 5 5 2 3 18 13
9 9 17 8 7 19 14 19 9 15
8 12 3 11 20 4 16 19 11
4 8 1 9 4
4 12 6 11 14
1 4
3 8 2 19
9 20 3 1 4 7 20 19 4 13
2 12 4
1 20
1 7
3 4 16 7
1 1
9 14 20 4 9 3 8 3 10 12
7 6 2 17 15 2 20 4
7 7 9 12 16 19 6 7
1 6
3 11 17 9
2 20 15
3 1 16 14
10 17 10 12 13 9 5 18 1 15 3
6 2 18 9 5 8 16
6 20 10 12 19 20 5
5 13 14 3 1 20
4 11 6 8 8
8 13 19 14 2 13 19 14 2
3 15 3 9
3 15 17 16
9 20 1 2 16 11 10 15 2 14
4 18 3 5 1
7 14 11 1 7 1 1 17
10 4 7 4 20 7 10 9 6 4 16
7 3 1 9 15 4 9 5
9 12 4 5 9 1 2 2 7 9
9 11 12 19 2 20 16 15 14 12
9 6 7 13 19 10 1 5 5 9
6 11 12 3 11 20 2
1 9
3 5 19 10
6 13 18 5 10 4 16
4 2 10 6 17
2 10 13
6 10 14 4 4 18 16
8 11 11 4 16 4 16 14 2
5 11 5 6 19 13
2 3 3
4 8 2 13 1
2 13 18
9 10 15 16 19 7 14 3 12 8
5 19 6 14 7 12
2 3 1
9 15 7 4 16 13 9 7 2 7
10 5 4 7 15 13 12 18 5 4 20
8 5 19 13 14 17 16 11 16
8 7 18 20 8 1 11 11 11
1 17
3 9 20 5
7 19 10 16 3 3 17 2
2 8 5
1 10
1 15
6 6 5 15 12 17 13
9 17 2 19 3 17 20 3 14 7
5 18 20 14 16 13
10 19 8 1 1 6 10 17 19 9 11
2 16 9
5 14 13 13 2 6
3 8 10 11
1 2
8 14 5 16 20 3 5 12 14
1 20
8 13 15 2 4 16 5 1 2
10 20 5 11 4 18 12 7 13 16 4
1 20
8 20 11 4 20 10 5 13 10
2 17 7
1 13
8 12 7 15 12 3 2 2 16
5 1 17 19 19 7
4 3 17 17 14
9 10 4 5 14 19 14 3 4 14
2 4 14
3 1 15 14
7 1 16 11 9 3 12 3
2 12 1
6 12 6 1 8 12 3`

func parseTestcases() ([]testcase, error) {
	lines := strings.Split(testcasesRaw, "\n")
	var res []testcase
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: invalid n", idx+1)
		}
		if len(fields) != n+1 {
			return nil, fmt.Errorf("line %d: expected %d ratings got %d", idx+1, n, len(fields)-1)
		}
		ratings := make([]int, n)
		for i := 0; i < n; i++ {
			ratings[i], err = strconv.Atoi(fields[i+1])
			if err != nil {
				return nil, fmt.Errorf("line %d: bad rating", idx+1)
			}
		}
		res = append(res, testcase{n: n, ratings: ratings})
	}
	return res, nil
}

// Embedded solver logic from 551A.go.
func solve(ratings []int) []int {
	freq := make([]int, 2001)
	for _, v := range ratings {
		if v >= 0 && v <= 2000 {
			freq[v]++
		}
	}
	suffix := make([]int, 2002)
	for r := 2000; r >= 1; r-- {
		suffix[r] = suffix[r+1] + freq[r+1]
	}
	res := make([]int, len(ratings))
	for i, v := range ratings {
		res[i] = suffix[v] + 1
	}
	return res
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse embedded testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range cases {
		expect := solve(tc.ratings)

		var input strings.Builder
		input.WriteString(strconv.Itoa(tc.n))
		input.WriteByte('\n')
		for i, v := range tc.ratings {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(strconv.Itoa(v))
		}
		input.WriteByte('\n')

		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input.String())
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("case %d: runtime error: %v\nstderr: %s\n", idx+1, err, stderr.String())
			os.Exit(1)
		}
		outParts := strings.Fields(strings.TrimSpace(out.String()))
		if len(outParts) != tc.n {
			fmt.Printf("case %d: expected %d numbers got %d\n", idx+1, tc.n, len(outParts))
			os.Exit(1)
		}
		for i := 0; i < tc.n; i++ {
			got, err := strconv.Atoi(outParts[i])
			if err != nil || got != expect[i] {
				fmt.Printf("case %d failed: expected %v got %v\n", idx+1, expect, outParts)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

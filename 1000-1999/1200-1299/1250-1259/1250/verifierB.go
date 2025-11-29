package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

type testcase struct {
	n, k  int
	teams []int
}

// solve mirrors the 1250B reference solution logic.
func solve(n, k int, teams []int) int64 {
	if len(teams) != k {
		return 0
	}
	sort.Ints(teams)
	best := int64(^uint64(0) >> 1)
	maxPair := 0
	for pairs := 0; pairs <= k/2; pairs++ {
		if pairs > 0 {
			sum := teams[pairs-1] + teams[k-pairs]
			if sum > maxPair {
				maxPair = sum
			}
		}
		rides := k - pairs
		largest := 0
		if pairs <= k-pairs-1 {
			largest = teams[k-pairs-1]
		}
		cap := largest
		if maxPair > cap {
			cap = maxPair
		}
		cost := int64(rides) * int64(cap)
		if cost < best {
			best = cost
		}
	}
	return best
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func loadTestcases() ([]testcase, error) {
	lines := strings.Split(strings.TrimSpace(testcaseData), "\n")
	tests := make([]testcase, 0, len(lines))
	for i, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			return nil, fmt.Errorf("test %d: not enough fields", i+1)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("test %d: bad n: %v", i+1, err)
		}
		k, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil, fmt.Errorf("test %d: bad k: %v", i+1, err)
		}
		if len(fields) != 2+n {
			return nil, fmt.Errorf("test %d: expected %d team entries, got %d", i+1, n, len(fields)-2)
		}
		counts := make([]int, k)
		for j := 0; j < n; j++ {
			val, err := strconv.Atoi(fields[2+j])
			if err != nil {
				return nil, fmt.Errorf("test %d: bad team value: %v", i+1, err)
			}
			if val <= 0 || val > k {
				return nil, fmt.Errorf("test %d: team value out of range: %d", i+1, val)
			}
			counts[val-1]++
		}
		tests = append(tests, testcase{n: n, k: k, teams: counts})
	}
	return tests, nil
}

const testcaseData = `
5 4 1 3 1 3 3
10 2 2 2 2 1 1 2 1 2 1 1
14 3 1 3 1 3 1 2 2 2 3 2 2 2 2 3
9 4 1 2 1 4 2 3 1 4 3
11 5 1 5 4 3 2 5 2 1 2 5 1
25 5 5 5 4 4 3 3 3 1 4 3 3 5 5 5 3 3 2 4 2 1 5 2 5 4 3
5 4 4 2 1 3 2
7 2 2 1 1 1 2 1 2
5 1 1 1 1 1 1
9 2 1 2 2 2 2 2 2 2 1
3 3 2 1 2
8 2 2 2 2 1 2 2 1 2
3 3 2 2 3
2 2 1 2
3 3 1 1 3
5 1 1 1 1 1 1
8 4 3 4 2 2 4 3 1 2
4 2 2 1 1 2
7 2 1 1 1 2 1 2 2
15 3 2 3 2 1 2 3 1 1 3 3 1 3 3 1 2
2 1 1 1
20 4 2 3 2 4 1 1 3 1 1 1 4 3 4 3 4 3 1 2 3 1
4 2 1 2 1 1
5 4 3 4 1 2 3
4 1 1 1 1 1
12 3 2 2 1 1 2 2 1 2 1 2 1 1
3 1 1 1 1
25 5 1 3 4 5 4 4 1 1 5 4 3 1 5 1 1 3 3 1 4 1 2 5 3 1 1
14 4 2 2 2 2 2 2 3 1 4 4 1 2 2 3
5 3 1 2 1 1 2
15 4 2 4 1 2 1 4 2 2 1 1 3 1 1 1 1
24 5 2 2 4 1 4 4 5 3 2 2 3 5 2 5 4 1 2 4 5 3 1 4 4 2
16 4 2 4 2 4 4 1 3 3 3 1 4 2 3 1 3 4
5 1 1 1 1 1 1
17 4 1 4 4 2 2 3 4 4 2 3 1 3 3 2 3 4 1
6 2 1 2 2 2 1 2
18 5 1 4 3 4 2 3 2 2 1 5 3 2 1 4 5 4 4 1
4 1 1 1 1 1
5 2 1 1 2 2 2
13 3 1 2 2 1 2 3 3 2 3 3 1 2 3
14 5 2 4 3 4 4 5 3 4 4 1 2 2 1 5
25 5 5 1 2 5 5 5 1 3 5 2 4 1 1 1 1 3 4 3 1 4 3 5 3 4 2
10 2 1 1 1 1 2 1 2 2 2 1
9 5 4 2 4 5 1 4 1 5 4
6 3 3 2 1 2 3 1
11 3 2 3 3 2 2 2 1 1 3 2 1
12 3 1 2 1 1 2 2 3 2 3 1 2 2
3 2 2 1 2
2 1 1 1
22 5 2 1 1 4 5 5 4 1 4 5 3 3 1 1 2 2 4 3 3 2 1 4
11 3 2 1 1 2 3 2 1 3 2 2 1
6 2 1 2 2 2 2 1
4 1 1 1 1 1
9 3 2 2 2 3 2 1 2 2 1
8 4 3 2 2 3 4 2 1 4
18 4 4 2 4 3 4 2 3 1 1 4 3 4 2 4 1 1 4 1
6 3 3 2 1 1 1 3
3 3 1 1 1
4 4 2 1 2 1
10 4 1 4 1 2 2 2 4 4 4 1
10 4 4 1 3 1 3 3 3 4 2 1
4 3 3 1 2 1
4 2 2 1 2 1
20 4 1 1 4 1 4 4 2 4 3 2 4 3 4 1 3 2 4 2 3 4
4 4 4 1 4 4
8 2 1 2 2 2 1 1 2 2
13 3 2 3 1 3 3 1 1 2 3 3 3 1 1
7 2 1 2 1 2 2 1 2
10 4 2 4 1 3 2 3 2 4 2 4
11 5 5 4 1 2 4 5 3 5 2 1 1
2 1 1 1
12 5 1 5 4 3 3 4 2 2 3 5 3 4
18 4 3 4 1 2 2 1 2 4 4 2 3 1 1 4 2 1 4 4
11 5 1 2 5 4 3 4 5 5 1 1 5
4 4 2 4 3 3
1 1 1
5 5 5 5 3 2 4
4 3 3 2 2 1
16 4 3 2 3 1 1 2 3 4 1 4 4 3 4 1 1 4
9 5 2 1 1 4 2 1 3 3 4
3 1 1 1 1
22 5 1 3 5 5 1 2 4 1 3 3 2 3 1 5 5 4 1 3 5 1 5 1
9 5 4 1 5 2 3 1 5 2 2
8 2 1 2 2 1 2 2 1 1
14 4 2 3 1 1 3 3 3 1 1 1 3 1 2 3
10 2 2 2 1 2 2 1 1 2 2 1
6 3 3 2 3 3 2 1
1 1 1
3 1 1 1 1
10 2 1 2 1 2 2 2 1 2 2 2
7 4 4 3 2 1 2 4 4
2 1 1 1
2 1 1 1
3 1 1 1 1
13 4 3 4 4 4 1 2 3 2 2 2 2 3 2
10 2 2 1 1 1 2 1 1 1 1 1
14 5 1 4 2 1 4 2 4 5 4 2 4 5 5 4
2 1 1 1
1 1 1
18 4 3 4 3 4 3 1 2 2 2 1 1 1 4 2 4 2 3 1
`

func formatInput(tc testcase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(tc.n))
	sb.WriteByte(' ')
	sb.WriteString(strconv.Itoa(tc.k))
	for team, cnt := range tc.teams {
		for i := 0; i < cnt; i++ {
			sb.WriteByte(' ')
			sb.WriteString(strconv.Itoa(team + 1))
		}
	}
	sb.WriteByte('\n')
	return sb.String()
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

	tests, err := loadTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load testcases: %v\n", err)
		os.Exit(1)
	}

	for i, tc := range tests {
		input := formatInput(tc)
		expect := fmt.Sprintf("%d", solve(tc.n, tc.k, append([]int(nil), tc.teams...)))
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, got, input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

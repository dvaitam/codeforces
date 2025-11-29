package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type tile struct {
	num  int
	suit byte
}

func parseTile(s string) tile {
	return tile{num: int(s[0] - '0'), suit: s[1]}
}

// solve mirrors the 1191B solution logic.
func solve(a, b, c string) int {
	t := [3]tile{parseTile(a), parseTile(b), parseTile(c)}
	if t[0] == t[1] && t[1] == t[2] {
		return 0
	}
	if t[0].suit == t[1].suit && t[1].suit == t[2].suit {
		nums := []int{t[0].num, t[1].num, t[2].num}
		sort.Ints(nums)
		if nums[0]+1 == nums[1] && nums[1]+1 == nums[2] {
			return 0
		}
	}
	for i := 0; i < 3; i++ {
		for j := i + 1; j < 3; j++ {
			ti, tj := t[i], t[j]
			if ti == tj {
				return 1
			}
			if ti.suit == tj.suit {
				d := ti.num - tj.num
				if d < 0 {
					d = -d
				}
				if d <= 2 {
					return 1
				}
			}
		}
	}
	return 2
}

const testcaseData = `
1m 1m 1m
1m 1m 1p
1m 1m 1s
1m 1m 2m
1m 1m 2p
1m 1m 2s
1m 1m 3m
1m 1m 3p
1m 1m 3s
1m 1m 4m
1m 1m 4p
1m 1m 4s
1m 1m 5m
1m 1m 5p
1m 1m 5s
1m 1m 6m
1m 1m 6p
1m 1m 6s
1m 1m 7m
1m 1m 7p
1m 1m 7s
1m 1m 8m
1m 1m 8p
1m 1m 8s
1m 1m 9m
1m 1m 9p
1m 1m 9s
1m 1p 1m
1m 1p 1p
1m 1p 1s
1m 1p 2m
1m 1p 2p
1m 1p 2s
1m 1p 3m
1m 1p 3p
1m 1p 3s
1m 1p 4m
1m 1p 4p
1m 1p 4s
1m 1p 5m
1m 1p 5p
1m 1p 5s
1m 1p 6m
1m 1p 6p
1m 1p 6s
1m 1p 7m
1m 1p 7p
1m 1p 7s
1m 1p 8m
1m 1p 8p
1m 1p 8s
1m 1p 9m
1m 1p 9p
1m 1p 9s
1m 1s 1m
1m 1s 1p
1m 1s 1s
1m 1s 2m
1m 1s 2p
1m 1s 2s
1m 1s 3m
1m 1s 3p
1m 1s 3s
1m 1s 4m
1m 1s 4p
1m 1s 4s
1m 1s 5m
1m 1s 5p
1m 1s 5s
1m 1s 6m
1m 1s 6p
1m 1s 6s
1m 1s 7m
1m 1s 7p
1m 1s 7s
1m 1s 8m
1m 1s 8p
1m 1s 8s
1m 1s 9m
1m 1s 9p
1m 1s 9s
1m 2m 1m
1m 2m 1p
1m 2m 1s
1m 2m 2m
1m 2m 2p
1m 2m 2s
1m 2m 3m
1m 2m 3p
1m 2m 3s
1m 2m 4m
1m 2m 4p
1m 2m 4s
1m 2m 5m
1m 2m 5p
1m 2m 5s
1m 2m 6m
1m 2m 6p
1m 2m 6s
1m 2m 7m
`

type testcase struct {
	a, b, c string
}

func loadTestcases() ([]testcase, error) {
	fields := strings.Fields(testcaseData)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no embedded testcases")
	}
	if len(fields)%3 != 0 {
		return nil, fmt.Errorf("malformed testcases: not divisible by 3")
	}
	res := make([]testcase, 0, len(fields)/3)
	for i := 0; i < len(fields); i += 3 {
		res = append(res, testcase{a: fields[i], b: fields[i+1], c: fields[i+2]})
	}
	return res, nil
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
		expect := solve(tc.a, tc.b, tc.c)
		input := fmt.Sprintf("%s %s %s\n", tc.a, tc.b, tc.c)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if got != fmt.Sprintf("%d", expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %s\ninput:\n%s", i+1, expect, got, input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

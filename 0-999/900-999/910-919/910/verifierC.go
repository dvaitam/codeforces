package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

type testCaseC struct {
	words []string
}

func parseCases(path string) ([]testCaseC, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	in := bufio.NewReader(f)
	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return nil, err
	}
	cases := make([]testCaseC, T)
	for i := 0; i < T; i++ {
		var n int
		if _, err := fmt.Fscan(in, &n); err != nil {
			return nil, err
		}
		words := make([]string, n)
		for j := 0; j < n; j++ {
			if _, err := fmt.Fscan(in, &words[j]); err != nil {
				return nil, err
			}
		}
		cases[i] = testCaseC{words: words}
	}
	return cases, nil
}

func solve(tc testCaseC) int64 {
	weights := make([]int64, 10)
	leading := make([]bool, 10)
	pow10 := [7]int64{1}
	for i := 1; i < 7; i++ {
		pow10[i] = pow10[i-1] * 10
	}
	for _, s := range tc.words {
		leading[s[0]-'a'] = true
		L := len(s)
		for j := L - 1; j >= 0; j-- {
			letter := s[j] - 'a'
			pos := L - 1 - j
			weights[letter] += pow10[pos]
		}
	}
	best := int64(^uint64(0) >> 1)
	for zero := 0; zero < 10; zero++ {
		if leading[zero] {
			continue
		}
		type pair struct {
			w   int64
			idx int
		}
		pairs := make([]pair, 0, 9)
		for i := 0; i < 10; i++ {
			if i == zero {
				continue
			}
			pairs = append(pairs, pair{weights[i], i})
		}
		sort.Slice(pairs, func(i, j int) bool {
			if pairs[i].w == pairs[j].w {
				return pairs[i].idx < pairs[j].idx
			}
			return pairs[i].w > pairs[j].w
		})
		digits := make([]int, 10)
		digits[zero] = 0
		for i, p := range pairs {
			digits[p.idx] = i + 1
		}
		var sum int64
		for i := 0; i < 10; i++ {
			sum += weights[i] * int64(digits[i])
		}
		if sum < best {
			best = sum
		}
	}
	return best
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases, err := parseCases("testcasesC.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", len(tc.words)))
		for _, w := range tc.words {
			sb.WriteString(w)
			sb.WriteByte('\n')
		}
		expected := solve(tc)
		gotStr, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		got, err := strconv.ParseInt(strings.TrimSpace(gotStr), 10, 64)
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

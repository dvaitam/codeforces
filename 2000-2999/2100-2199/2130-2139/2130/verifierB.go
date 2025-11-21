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

type testCase struct {
	n   int
	s   int
	arr []int
	cnt [3]int
}

func genCase(rng *rand.Rand) testCase {
	n := rng.Intn(5) + 3 // between 3 and 7
	cnt := [3]int{1, 1, 1}
	for i := 3; i < n; i++ {
		v := rng.Intn(3)
		cnt[v]++
	}
	arr := make([]int, 0, n)
	for v := 0; v < 3; v++ {
		for i := 0; i < cnt[v]; i++ {
			arr = append(arr, v)
		}
	}
	rng.Shuffle(len(arr), func(i, j int) {
		arr[i], arr[j] = arr[j], arr[i]
	})
	s := rng.Intn(25) + 1
	return testCase{n: n, s: s, arr: arr, cnt: cnt}
}

func existsPath(arr []int, s int) bool {
	n := len(arr)
	start := arr[0]
	if start > s {
		return false
	}
	visited := make([][]bool, n)
	for i := range visited {
		visited[i] = make([]bool, s+1)
	}
	type state struct {
		pos int
		sum int
	}
	queue := []state{{0, start}}
	visited[0][start] = true
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		if cur.pos == n-1 && cur.sum == s {
			return true
		}
		for _, nxt := range []int{cur.pos - 1, cur.pos + 1} {
			if nxt < 0 || nxt >= n {
				continue
			}
			ns := cur.sum + arr[nxt]
			if ns > s || visited[nxt][ns] {
				continue
			}
			visited[nxt][ns] = true
			queue = append(queue, state{nxt, ns})
		}
	}
	return false
}

func canBlock(cnt [3]int, s int) (bool, []int) {
	n := cnt[0] + cnt[1] + cnt[2]
	arr := make([]int, n)
	var dfs func(int) (bool, []int)
	dfs = func(idx int) (bool, []int) {
		if idx == n {
			if !existsPath(arr, s) {
				cp := make([]int, n)
				copy(cp, arr)
				return true, cp
			}
			return false, nil
		}
		for v := 0; v < 3; v++ {
			if cnt[v] == 0 {
				continue
			}
			cnt[v]--
			arr[idx] = v
			if ok, res := dfs(idx + 1); ok {
				return true, res
			}
			cnt[v]++
		}
		return false, nil
	}
	return dfs(0)
}

func buildInput(cases []testCase) (string, []bool) {
	var sb strings.Builder
	fmt.Fprintln(&sb, len(cases))
	expectBlock := make([]bool, len(cases))
	for i, tc := range cases {
		fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.s)
		for j, v := range tc.arr {
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprint(&sb, v)
		}
		sb.WriteByte('\n')
		ok, _ := canBlock(tc.cnt, tc.s)
		expectBlock[i] = ok
	}
	return sb.String(), expectBlock
}

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return out.String(), nil
}

func parseCase(tokens []string, ptr *int, n int) (bool, []int, error) {
	if *ptr >= len(tokens) {
		return false, nil, fmt.Errorf("not enough tokens")
	}
	if tokens[*ptr] == "-1" {
		*ptr++
		return true, nil, nil
	}
	if len(tokens)-*ptr < n {
		return false, nil, fmt.Errorf("not enough numbers for arrangement")
	}
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		val, err := strconv.Atoi(tokens[*ptr+i])
		if err != nil {
			return false, nil, fmt.Errorf("invalid integer %q", tokens[*ptr+i])
		}
		if val < 0 || val > 2 {
			return false, nil, fmt.Errorf("value out of range %d", val)
		}
		arr[i] = val
	}
	*ptr += n
	return false, arr, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/2130B_binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var cases []testCase
	for i := 0; i < 200; i++ {
		cases = append(cases, genCase(rng))
	}
	input, expectBlock := buildInput(cases)
	output, err := runCandidate(bin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to run candidate: %v\n", err)
		os.Exit(1)
	}
	tokens := strings.Fields(output)
	ptr := 0
	for idx, tc := range cases {
		isNeg, arr, err := parseCase(tokens, &ptr, tc.n)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d parse error: %v\n", idx+1, err)
			os.Exit(1)
		}
		if isNeg {
			if expectBlock[idx] {
				fmt.Fprintf(os.Stderr, "case %d: arrangement exists but candidate printed -1\nn=%d s=%d arr=%v\n", idx+1, tc.n, tc.s, tc.arr)
				os.Exit(1)
			}
			continue
		}
		// validate multiset match
		cnt := [3]int{}
		for _, v := range arr {
			cnt[v]++
		}
		if cnt != tc.cnt {
			fmt.Fprintf(os.Stderr, "case %d: arrangement uses wrong multiset expected %v got %v\n", idx+1, tc.cnt, cnt)
			os.Exit(1)
		}
		if existsPath(arr, tc.s) {
			fmt.Fprintf(os.Stderr, "case %d: provided arrangement still allows a path\nn=%d s=%d arr=%v\n", idx+1, tc.n, tc.s, arr)
			os.Exit(1)
		}
		if !expectBlock[idx] {
			fmt.Fprintf(os.Stderr, "case %d: candidate claims arrangement blocks but impossible\nn=%d s=%d arr=%v original=%v\n", idx+1, tc.n, tc.s, arr, tc.arr)
			os.Exit(1)
		}
	}
	if ptr != len(tokens) {
		fmt.Fprintf(os.Stderr, "extra output tokens found\n")
		os.Exit(1)
	}
	fmt.Println("All tests passed")
}

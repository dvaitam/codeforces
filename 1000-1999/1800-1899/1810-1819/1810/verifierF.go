package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type testCase struct {
	n   int
	m   int
	q   int
	a   []int
	ops [][2]int
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

// generateSplits returns all m-tuples of positive integers summing to n where
// each value corresponds to a valid subtree leaf count.
func generateSplits(n, m int, cur []int, res *[][]int) {
	if len(cur) == m {
		if n == 0 {
			copyTuple := append([]int(nil), cur...)
			*res = append(*res, copyTuple)
		}
		return
	}
	minRemaining := (m - len(cur) - 1) + 1
	for v := 1; v <= n-minRemaining+1; v++ {
		if (v-1)%(m-1) != 0 {
			continue
		}
		generateSplits(n-v, m, append(cur, v), res)
	}
}

var depthMemo = make(map[[2]int][][]int)

func depthShapes(n, m int) [][]int {
	key := [2]int{n, m}
	if v, ok := depthMemo[key]; ok {
		return v
	}
	if n == 1 {
		depthMemo[key] = [][]int{{0}}
		return depthMemo[key]
	}
	var res [][]int
	var splits [][]int
	generateSplits(n, m, nil, &splits)
	for _, sp := range splits {
		children := make([][][]int, m)
		for i, v := range sp {
			children[i] = depthShapes(v, m)
		}
		var combine func(idx int, cur []int)
		combine = func(idx int, cur []int) {
			if idx == m {
				res = append(res, append([]int(nil), cur...))
				return
			}
			for _, child := range children[idx] {
				tmp := append([]int(nil), cur...)
				for _, d := range child {
					tmp = append(tmp, d+1)
				}
				combine(idx+1, tmp)
			}
		}
		combine(0, nil)
	}
	depthMemo[key] = res
	return res
}

func bestValue(a []int, m int) int {
	vals := append([]int(nil), a...)
	sort.Slice(vals, func(i, j int) bool { return vals[i] > vals[j] })
	ans := 0
	for i, depths := range depthShapes(len(vals), m) {
		sort.Ints(depths)
		cur := 0
		for j, v := range vals {
			if v+depths[j] > cur {
				cur = v + depths[j]
			}
		}
		if i == 0 || cur < ans {
			ans = cur
		}
	}
	return ans
}

func buildCase(rng *rand.Rand) testCase {
	m := rng.Intn(3) + 2 // ensure m >= 2
	k := rng.Intn(2) + 1
	n := k*(m-1) + 1
	q := rng.Intn(5) + 1
	a := make([]int, n)
	for i := range a {
		a[i] = rng.Intn(5) + 1
	}
	ops := make([][2]int, q)
	for i := 0; i < q; i++ {
		ops[i] = [2]int{rng.Intn(n) + 1, rng.Intn(5) + 1}
	}
	return testCase{n: n, m: m, q: q, a: a, ops: ops}
}

func expectedOutputs(tc testCase) []int {
	arr := append([]int(nil), tc.a...)
	res := make([]int, tc.q)
	for i, op := range tc.ops {
		arr[op[0]-1] = op[1]
		res[i] = bestValue(arr, tc.m)
	}
	return res
}

func formatCase(tc testCase) string {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d %d\n", tc.n, tc.m, tc.q))
	for i, v := range tc.a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	for _, op := range tc.ops {
		sb.WriteString(fmt.Sprintf("%d %d\n", op[0], op[1]))
	}
	return sb.String()
}

func parseInts(out string) ([]int, error) {
	fields := strings.Fields(out)
	res := make([]int, len(fields))
	for i, f := range fields {
		var v int
		if _, err := fmt.Sscan(f, &v); err != nil {
			return nil, err
		}
		res[i] = v
	}
	return res, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		tc := buildCase(rng)
		input := formatCase(tc)
		candOut, cErr := runBinary(candidate, input)
		if cErr != nil {
			fmt.Fprintf(os.Stderr, "case %d candidate error: %v\n%s", t+1, cErr, candOut)
			os.Exit(1)
		}
		expected := expectedOutputs(tc)
		got, err := parseInts(candOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed to parse output: %v\n%s", t+1, err, candOut)
			os.Exit(1)
		}
		if len(got) != len(expected) {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%sexpected length %d, got %d (%v)", t+1, input, len(expected), len(got), got)
			os.Exit(1)
		}
		for i, v := range expected {
			if got[i] != v {
				fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%sexpected:%v\nactual:%v\n", t+1, input, expected, got)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}

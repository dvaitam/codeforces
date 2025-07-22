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

const MaxOut = 100000

type testCase struct {
	input    string
	expected string
}

func sequenceE(X, k int64) []int64 {
	vals := make([]int64, 0)
	for i := int64(1); i*i <= X; i++ {
		if X%i == 0 {
			vals = append(vals, i)
			if i*i != X {
				vals = append(vals, X/i)
			}
		}
	}
	sort.Slice(vals, func(i, j int) bool { return vals[i] < vals[j] })
	idx := make(map[int64]int, len(vals))
	for i, v := range vals {
		idx[v] = i
	}
	idx1 := idx[1]
	children := make([][]int, len(vals))
	for i := range vals {
		v := vals[i]
		for j := 0; j <= i; j++ {
			if v%vals[j] == 0 {
				children[i] = append(children[i], j)
			}
		}
	}
	rem := MaxOut
	res := make([]int64, 0, MaxOut)
	var gen func(int, int64)
	gen = func(vIdx int, depth int64) {
		if rem <= 0 {
			return
		}
		if vIdx == idx1 {
			res = append(res, 1)
			rem--
			return
		}
		ch := children[vIdx]
		if len(ch) == 2 && ch[0] == idx1 && ch[1] == vIdx {
			if depth >= int64(rem) {
				for i := 0; i < rem; i++ {
					res = append(res, 1)
				}
				rem = 0
			} else {
				for i := int64(0); i < depth && rem > 0; i++ {
					res = append(res, 1)
					rem--
				}
				if rem > 0 {
					res = append(res, vals[vIdx])
					rem--
				}
			}
			return
		}
		if depth == 0 {
			res = append(res, vals[vIdx])
			rem--
			return
		}
		for _, u := range ch {
			if rem <= 0 {
				return
			}
			gen(u, depth-1)
		}
	}
	gen(idx[X], k)
	return res
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]testCase, 100)
	for i := 0; i < 100; i++ {
		X := int64(rng.Intn(1000) + 1)
		k := int64(rng.Intn(10))
		seq := sequenceE(X, k)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", X, k))
		var out strings.Builder
		for j, v := range seq {
			if j > 0 {
				out.WriteByte(' ')
			}
			out.WriteString(fmt.Sprintf("%d", v))
		}
		out.WriteByte('\n')
		cases[i] = testCase{input: sb.String(), expected: out.String()}
	}
	return cases
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
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := generateTests()
	for i, tc := range cases {
		out, err := run(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
		if out != strings.TrimSpace(tc.expected) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, strings.TrimSpace(tc.expected), out, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

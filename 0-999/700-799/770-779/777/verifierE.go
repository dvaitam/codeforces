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

type ring struct {
	a, b int
	h    int
}

type testCase struct {
	rings []ring
}

// Fenwick tree for prefix maximum
type bit struct {
	n    int
	tree []int64
}

func newBIT(n int) *bit {
	return &bit{n: n, tree: make([]int64, n+2)}
}

func (b *bit) update(i int, val int64) {
	for i <= b.n {
		if val > b.tree[i] {
			b.tree[i] = val
		}
		i += i & -i
	}
}

func (b *bit) query(i int) int64 {
	var res int64
	for i > 0 {
		if b.tree[i] > res {
			res = b.tree[i]
		}
		i -= i & -i
	}
	return res
}

func solve(tc testCase) int64 {
	rings := make([]ring, len(tc.rings))
	copy(rings, tc.rings)
	sort.Slice(rings, func(i, j int) bool {
		if rings[i].b == rings[j].b {
			return rings[i].a > rings[j].a
		}
		return rings[i].b > rings[j].b
	})
	bVals := make([]int, len(rings))
	for i := range rings {
		bVals[i] = rings[i].b
	}
	sort.Ints(bVals)
	uniq := make([]int, 0, len(bVals))
	for i := len(bVals) - 1; i >= 0; i-- {
		if i == len(bVals)-1 || bVals[i] != bVals[i+1] {
			uniq = append(uniq, bVals[i])
		}
	}
	idxMap := make(map[int]int)
	for i, v := range uniq {
		idxMap[v] = i
	}
	bt := newBIT(len(uniq))
	var ans int64
	for _, r := range rings {
		pos := sort.Search(len(uniq), func(i int) bool { return uniq[i] <= r.a })
		var best int64
		if pos > 0 {
			best = bt.query(pos)
		}
		cur := int64(r.h) + best
		bt.update(idxMap[r.b]+1, cur)
		if cur > ans {
			ans = cur
		}
	}
	return ans
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

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(tc.rings)))
	for _, r := range tc.rings {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", r.a, r.b, r.h))
	}
	return sb.String()
}

func genCases() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]testCase, 0, 110)
	for i := 0; i < 110; i++ {
		n := rng.Intn(8) + 1
		arr := make([]ring, n)
		for j := 0; j < n; j++ {
			a := rng.Intn(50) + 1
			b := a + rng.Intn(50) + 1
			h := rng.Intn(50) + 1
			arr[j] = ring{a: a, b: b, h: h}
		}
		cases = append(cases, testCase{rings: arr})
	}
	return cases
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := genCases()
	for i, tc := range cases {
		input := buildInput(tc)
		out, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		exp := fmt.Sprint(solve(tc))
		if strings.TrimSpace(out) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

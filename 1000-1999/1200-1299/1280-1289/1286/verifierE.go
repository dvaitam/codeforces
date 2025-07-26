package main

import (
	"bytes"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type SegTree struct {
	n    int
	size int
	tree []int64
}

func NewSegTree(n int) *SegTree {
	size := 1
	for size < n {
		size <<= 1
	}
	tree := make([]int64, 2*size)
	for i := range tree {
		tree[i] = 1<<63 - 1
	}
	return &SegTree{n: n, size: size, tree: tree}
}

func (st *SegTree) Update(pos int, val int64) {
	i := pos + st.size - 1
	st.tree[i] = val
	for i >>= 1; i > 0; i >>= 1 {
		if st.tree[i<<1] < st.tree[i<<1|1] {
			st.tree[i] = st.tree[i<<1]
		} else {
			st.tree[i] = st.tree[i<<1|1]
		}
	}
}

func (st *SegTree) Query(l, r int) int64 {
	l += st.size - 1
	r += st.size - 1
	res := int64(1<<63 - 1)
	for l <= r {
		if l&1 == 1 {
			if st.tree[l] < res {
				res = st.tree[l]
			}
			l++
		}
		if r&1 == 0 {
			if st.tree[r] < res {
				res = st.tree[r]
			}
			r--
		}
		l >>= 1
		r >>= 1
	}
	return res
}

func solveE(n int, queries []struct {
	c byte
	w int64
}) []string {
	S := make([]byte, 0, n)
	W := make([]int64, 0, n+1)
	pi := make([]int, n+1)
	seg := NewSegTree(n + 2)

	MASK := int64((1 << 30) - 1)
	ans := big.NewInt(0)
	prevAns := big.NewInt(0)
	res := make([]string, 0, n)

	for i := 1; i <= n; i++ {
		ch := queries[i-1].c
		w := queries[i-1].w

		shift := new(big.Int).Mod(prevAns, big.NewInt(26))
		shiftInt := int(shift.Int64())
		c := byte('a' + (int(ch-'a')+shiftInt)%26)
		maskVal := new(big.Int).And(prevAns, big.NewInt(MASK))
		w ^= maskVal.Int64()

		S = append(S, c)
		W = append(W, w)
		seg.Update(i, w)

		j := pi[i-1]
		for j > 0 && S[i-1] != S[j] {
			j = pi[j]
		}
		if S[i-1] == S[j] {
			j++
		}
		pi[i] = j

		cur := i
		for cur > 0 {
			l := cur
			minv := seg.Query(i-l+1, i)
			ans.Add(ans, big.NewInt(minv))
			cur = pi[cur]
		}

		res = append(res, ans.String())
		prevAns.Set(ans)
	}
	return res
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
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func generateTests() []string {
	rng := rand.New(rand.NewSource(6))
	tests := make([]string, 100)
	letters := []byte("abcdefghijklmnopqrstuvwxyz")
	for i := range tests {
		n := rng.Intn(10) + 1
		var sb strings.Builder
		sb.WriteString(fmt.Sprint(n))
		sb.WriteByte('\n')
		for j := 0; j < n; j++ {
			ch := letters[rng.Intn(len(letters))]
			w := rng.Int63n(100)
			fmt.Fprintf(&sb, "%c %d\n", ch, w)
		}
		tests[i] = sb.String()
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, input := range tests {
		lines := strings.Fields(input)
		n := atoi(lines[0])
		queries := make([]struct {
			c byte
			w int64
		}, n)
		idx := 1
		for j := 0; j < n; j++ {
			c := lines[idx][0]
			w, _ := strconv.ParseInt(lines[idx+1], 10, 64)
			queries[j] = struct {
				c byte
				w int64
			}{c: c, w: w}
			idx += 2
		}
		expLines := solveE(n, queries)
		exp := strings.Join(expLines, "\n")
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected\n%s\ngot\n%s\ninput:\n%s", i+1, exp, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

func atoi(s string) int {
	v, _ := strconv.Atoi(s)
	return v
}

package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

const alpha = 26

var trees [alpha][]int
var lazy [alpha][]int

func push(c, node, l, r int) {
	if lazy[c][node] != -1 {
		val := lazy[c][node]
		trees[c][node] = (r - l + 1) * val
		if l != r {
			lazy[c][node*2] = val
			lazy[c][node*2+1] = val
		}
		lazy[c][node] = -1
	}
}

func update(c, node, l, r, ql, qr, val int) {
	push(c, node, l, r)
	if ql > r || qr < l {
		return
	}
	if ql <= l && r <= qr {
		lazy[c][node] = val
		push(c, node, l, r)
		return
	}
	mid := (l + r) / 2
	update(c, node*2, l, mid, ql, qr, val)
	update(c, node*2+1, mid+1, r, ql, qr, val)
	trees[c][node] = trees[c][node*2] + trees[c][node*2+1]
}

func query(c, node, l, r, ql, qr int) int {
	push(c, node, l, r)
	if ql > r || qr < l {
		return 0
	}
	if ql <= l && r <= qr {
		return trees[c][node]
	}
	mid := (l + r) / 2
	return query(c, node*2, l, mid, ql, qr) + query(c, node*2+1, mid+1, r, ql, qr)
}

type testCase struct {
	input    string
	expected string
}

func expectedAnswer(n int, s string, queries [][3]int) string {
	for i := 0; i < alpha; i++ {
		trees[i] = make([]int, 4*(n+2))
		lazy[i] = make([]int, 4*(n+2))
		for j := range lazy[i] {
			lazy[i][j] = -1
		}
	}
	for i, ch := range s {
		c := int(ch - 'a')
		update(c, 1, 1, n, i+1, i+1, 1)
	}
	for _, q := range queries {
		l := q[0]
		r := q[1]
		k := q[2]
		if l > r {
			l, r = r, l
		}
		cnt := make([]int, alpha)
		for c := 0; c < alpha; c++ {
			cnt[c] = query(c, 1, 1, n, l, r)
			if cnt[c] > 0 {
				update(c, 1, 1, n, l, r, 0)
			}
		}
		idx := l
		if k == 1 {
			for c := 0; c < alpha; c++ {
				if cnt[c] > 0 {
					update(c, 1, 1, n, idx, idx+cnt[c]-1, 1)
					idx += cnt[c]
				}
			}
		} else {
			for c := alpha - 1; c >= 0; c-- {
				if cnt[c] > 0 {
					update(c, 1, 1, n, idx, idx+cnt[c]-1, 1)
					idx += cnt[c]
				}
			}
		}
	}
	out := make([]byte, n)
	for i := 1; i <= n; i++ {
		for c := 0; c < alpha; c++ {
			if query(c, 1, 1, n, i, i) > 0 {
				out[i-1] = byte('a' + c)
				break
			}
		}
	}
	return string(out)
}

func generateCase(rng *rand.Rand) testCase {
	n := rng.Intn(10) + 1
	q := rng.Intn(10)
	b := make([]byte, n)
	for i := range b {
		b[i] = byte('a' + rng.Intn(3))
	}
	s := string(b)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, q))
	sb.WriteString(fmt.Sprintf("%s\n", s))
	queries := make([][3]int, q)
	for i := 0; i < q; i++ {
		l := rng.Intn(n) + 1
		r := rng.Intn(n) + 1
		if l > r {
			l, r = r, l
		}
		k := rng.Intn(2)
		sb.WriteString(fmt.Sprintf("%d %d %d\n", l, r, k))
		queries[i] = [3]int{l, r, k}
	}
	expect := expectedAnswer(n, s, queries)
	return testCase{input: sb.String(), expected: expect}
}

func runCase(bin string, tc testCase) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(tc.input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != tc.expected {
		return fmt.Errorf("expected %s got %s", tc.expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCase(rng)
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

package main

import (
	"bufio"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type test struct {
	input    string
	expected string
}

// BIT implements a Fenwick tree for prefix sums and order statistics.
type BIT struct {
	n    int
	tree []int
}

func NewBIT(n int) *BIT { return &BIT{n: n, tree: make([]int, n+2)} }
func (b *BIT) Add(i, delta int) {
	for x := i; x <= b.n; x += x & -x {
		b.tree[x] += delta
	}
}
func (b *BIT) Sum(i int) int {
	s := 0
	for x := i; x > 0; x -= x & -x {
		s += b.tree[x]
	}
	return s
}
func (b *BIT) FindKth(k int) int {
	pos := 0
	bitMask := 1
	for bitMask<<1 <= b.n {
		bitMask <<= 1
	}
	for d := bitMask; d > 0; d >>= 1 {
		nxt := pos + d
		if nxt <= b.n && b.tree[nxt] < k {
			k -= b.tree[nxt]
			pos = nxt
		}
	}
	return pos + 1
}
func (b *BIT) FindNextFrom(l int) int {
	pre := b.Sum(l - 1)
	total := b.Sum(b.n)
	if pre == total {
		return b.n + 1
	}
	return b.FindKth(pre + 1)
}

func charIndex(ch byte) int {
	if ch >= 'a' && ch <= 'z' {
		return int(ch - 'a')
	}
	if ch >= 'A' && ch <= 'Z' {
		return 26 + int(ch-'A')
	}
	return 52 + int(ch-'0')
}

func solve(input string) string {
	reader := strings.NewReader(input)
	in := bufio.NewReader(reader)
	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return ""
	}
	var s string
	fmt.Fscan(in, &s)
	bits := make([]*BIT, 62)
	for i := 0; i < 62; i++ {
		bits[i] = NewBIT(n)
	}
	global := NewBIT(n)
	alive := make([]bool, n+1)
	for i := 1; i <= n; i++ {
		global.Add(i, 1)
		idx := charIndex(s[i-1])
		bits[idx].Add(i, 1)
		alive[i] = true
	}
	for op := 0; op < m; op++ {
		var l, r int
		var cs string
		fmt.Fscan(in, &l, &r, &cs)
		c := cs[0]
		idx := charIndex(c)
		lIdx := global.FindKth(l)
		rIdx := global.FindKth(r)
		limit := rIdx
		pos := bits[idx].FindNextFrom(lIdx)
		for pos <= limit {
			if alive[pos] {
				alive[pos] = false
				global.Add(pos, -1)
				bits[idx].Add(pos, -1)
			}
			pos = bits[idx].FindNextFrom(pos + 1)
		}
	}
	res := make([]byte, 0, n)
	for i := 1; i <= n; i++ {
		if alive[i] {
			res = append(res, s[i-1])
		}
	}
	return string(res)
}

func generateTests() []test {
	rand.Seed(8996)
	var tests []test
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	for len(tests) < 100 {
		n := rand.Intn(20) + 1
		m := rand.Intn(20) + 1
		cur := make([]byte, n)
		for i := 0; i < n; i++ {
			cur[i] = chars[rand.Intn(len(chars))]
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", n, m)
		sb.Write(cur)
		sb.WriteByte('\n')
		for i := 0; i < m; i++ {
			if len(cur) == 0 {
				fmt.Fprintf(&sb, "1 1 a\n")
				cur = cur[:0]
				continue
			}
			l := rand.Intn(len(cur)) + 1
			r := rand.Intn(len(cur)-l+1) + l
			c := chars[rand.Intn(len(chars))]
			fmt.Fprintf(&sb, "%d %d %c\n", l, r, c)
			filtered := make([]byte, 0, len(cur))
			for idx, ch := range cur {
				pos := idx + 1
				if pos >= l && pos <= r && ch == c {
					continue
				}
				filtered = append(filtered, ch)
			}
			cur = filtered
		}
		inp := sb.String()
		tests = append(tests, test{inp, solve(inp)})
	}
	return tests
}

func runBinary(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.CommandContext(ctx, "go", "run", bin)
	} else {
		cmd = exec.CommandContext(ctx, bin)
	}
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("time limit")
	}
	if err != nil {
		return "", fmt.Errorf("%v: %s", err, out)
	}
	return strings.TrimSpace(string(out)), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		out, err := runBinary(bin, t.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if out != strings.TrimSpace(t.expected) {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d\ninput:\n%sexpected:%s\n got:%s\n", i+1, t.input, t.expected, out)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

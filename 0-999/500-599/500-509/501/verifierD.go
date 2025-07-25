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

type Fenwick struct {
	n    int
	tree []int
}

func NewFenwick(n int) *Fenwick { return &Fenwick{n: n, tree: make([]int, n+1)} }

func (f *Fenwick) Add(i, delta int) {
	for ; i <= f.n; i += i & -i {
		f.tree[i] += delta
	}
}
func (f *Fenwick) Sum(i int) int {
	s := 0
	for ; i > 0; i -= i & -i {
		s += f.tree[i]
	}
	return s
}

func (f *Fenwick) FindByOrder(order int) int {
	idx := 0
	bit := 1
	for bit<<1 <= f.n {
		bit <<= 1
	}
	for b := bit; b > 0; b >>= 1 {
		next := idx + b
		if next <= f.n && f.tree[next] < order {
			order -= f.tree[next]
			idx = next
		}
	}
	return idx + 1
}

func solveD(p, q []int) string {
	n := len(p)
	c := make([]int, n)
	d := make([]int, n)
	bitP := NewFenwick(n)
	bitQ := NewFenwick(n)
	for i := 1; i <= n; i++ {
		bitP.Add(i, 1)
		bitQ.Add(i, 1)
	}
	for i := 0; i < n; i++ {
		x := p[i] + 1
		c[i] = bitP.Sum(x - 1)
		bitP.Add(x, -1)
		y := q[i] + 1
		d[i] = bitQ.Sum(y - 1)
		bitQ.Add(y, -1)
	}
	e := make([]int, n)
	carry := 0
	for i := n - 1; i >= 0; i-- {
		base := n - i
		sum := c[i] + d[i] + carry
		e[i] = sum % base
		carry = sum / base
	}
	bitR := NewFenwick(n)
	for i := 1; i <= n; i++ {
		bitR.Add(i, 1)
	}
	res := make([]int, n)
	for i := 0; i < n; i++ {
		idx := bitR.FindByOrder(e[i] + 1)
		res[i] = idx - 1
		bitR.Add(idx, -1)
	}
	var sb strings.Builder
	for i, v := range res {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(8) + 1
	p := rng.Perm(n)
	q := rng.Perm(n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range p {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	for i, v := range q {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	expect := solveD(p, q)
	return sb.String(), expect
}

func runCase(bin, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(expected)
	if got != exp {
		return fmt.Errorf("expected %s got %s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

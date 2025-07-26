package main

import (
	"bytes"
	"fmt"
	"math/bits"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

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

type BIT struct {
	n    int
	tree []int64
}

func NewBIT(n int) *BIT {
	return &BIT{n: n, tree: make([]int64, n+1)}
}

func (b *BIT) Add(i int, v int64) {
	for i <= b.n {
		b.tree[i] += v
		i += i & -i
	}
}

func (b *BIT) Sum(i int) int64 {
	var s int64
	for i > 0 {
		s += b.tree[i]
		i -= i & -i
	}
	return s
}

func (b *BIT) LowerBound(target int64) int {
	idx := 0
	bitMask := 1 << (bits.Len(uint(b.n)) - 1)
	for bitMask > 0 {
		next := idx + bitMask
		if next <= b.n && b.tree[next] <= target {
			target -= b.tree[next]
			idx = next
		}
		bitMask >>= 1
	}
	return idx + 1
}

func solveCase(s []int64) []int {
	n := len(s)
	bit := NewBIT(n)
	for i := 1; i <= n; i++ {
		bit.Add(i, int64(i))
	}
	p := make([]int, n)
	for i := n - 1; i >= 0; i-- {
		idx := bit.LowerBound(s[i])
		p[i] = idx
		bit.Add(idx, -int64(idx))
	}
	return p
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(6) + 1
	bit := NewBIT(n)
	for i := 1; i <= n; i++ {
		bit.Add(i, int64(i))
	}
	s := make([]int64, n)
	for i := n - 1; i >= 0; i-- {
		val := rng.Int63n(bit.Sum(n))
		idx := bit.LowerBound(val)
		s[i] = val
		bit.Add(idx, -int64(idx))
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range s {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	p := solveCase(s)
	var exp strings.Builder
	for i, v := range p {
		if i > 0 {
			exp.WriteByte(' ')
		}
		exp.WriteString(strconv.Itoa(v))
	}
	return sb.String(), exp.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

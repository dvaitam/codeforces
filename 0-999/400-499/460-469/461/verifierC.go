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

type BIT struct {
	n   int
	bit []int64
}

func NewBIT(n int) *BIT {
	return &BIT{n: n, bit: make([]int64, n+1)}
}

func (b *BIT) Add(i int, v int64) {
	for x := i; x <= b.n; x += x & -x {
		b.bit[x] += v
	}
}

func (b *BIT) Sum(i int) int64 {
	var s int64
	for x := i; x > 0; x -= x & -x {
		s += b.bit[x]
	}
	return s
}

func (b *BIT) RangeSum(l, r int) int64 {
	if l > r {
		return 0
	}
	return b.Sum(r) - b.Sum(l-1)
}

func solveCase(n int, queries []query) string {
	bit := NewBIT(n)
	for i := 1; i <= n; i++ {
		bit.Add(i, 1)
	}
	l, r := 0, n
	rev := false
	var out strings.Builder
	for _, q := range queries {
		if q.t == 1 {
			p := q.p
			w := r - l
			if p <= w-p {
				for i := 0; i < p; i++ {
					var x1, x2 int
					if !rev {
						x1 = l + i
						x2 = l + (w - 1 - i)
					} else {
						x1 = r - 1 - i
						x2 = r - 1 - (w - 1 - i)
					}
					v := bit.RangeSum(x1+1, x1+1)
					bit.Add(x2+1, v)
				}
				if !rev {
					l += p
				} else {
					r -= p
				}
			} else {
				sz := w - p
				for i := 0; i < sz; i++ {
					var x1, x2 int
					if !rev {
						x1 = l + (w - 1 - i)
						x2 = l + i
					} else {
						x1 = r - 1 - (w - 1 - i)
						x2 = r - 1 - i
					}
					v := bit.RangeSum(x1+1, x1+1)
					bit.Add(x2+1, v)
				}
				if !rev {
					r -= sz
				} else {
					l += sz
				}
				rev = !rev
			}
		} else {
			li := q.l
			ri := q.r
			var a, b int
			if !rev {
				a = l + li
				b = l + ri - 1
			} else {
				a = r - 1 - li
				b = r - ri
			}
			if a > b {
				a, b = b, a
			}
			ans := bit.RangeSum(a+1, b+1)
			out.WriteString(fmt.Sprintf("%d\n", ans))
		}
	}
	return strings.TrimSpace(out.String())
}

type query struct {
	t int
	p int
	l int
	r int
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 1
	q := rng.Intn(10) + 1
	var queries []query
	width := n
	for i := 0; i < q; i++ {
		if rng.Intn(2) == 0 && width > 1 {
			p := rng.Intn(width-1) + 1
			queries = append(queries, query{t: 1, p: p})
			if p <= width-p {
				width -= p
			} else {
				width = p
			}
		} else {
			li := rng.Intn(width)
			ri := rng.Intn(width-li) + li + 1
			queries = append(queries, query{t: 2, l: li, r: ri})
		}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, q)
	for _, qu := range queries {
		if qu.t == 1 {
			fmt.Fprintf(&sb, "1 %d\n", qu.p)
		} else {
			fmt.Fprintf(&sb, "2 %d %d\n", qu.l, qu.r)
		}
	}
	exp := solveCase(n, queries)
	return sb.String(), exp
}

func runCase(bin, in, exp string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(in)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

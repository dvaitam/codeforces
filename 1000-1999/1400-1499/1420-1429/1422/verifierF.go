package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

const MOD = 1000000007

type Fenwick struct {
	n    int
	tree []int64
}

func NewFenwick(n int) *Fenwick {
	f := &Fenwick{n: n, tree: make([]int64, n+2)}
	for i := range f.tree {
		f.tree[i] = 1
	}
	return f
}

func (f *Fenwick) mul(i int, v int64) {
	for ; i <= f.n; i += i & -i {
		f.tree[i] = f.tree[i] * v % MOD
	}
}

func (f *Fenwick) rangeMul(l, r int, v, invV int64) {
	if l > r {
		return
	}
	f.mul(l, v)
	if r+1 <= f.n {
		f.mul(r+1, invV)
	}
}

func (f *Fenwick) query(x int) int64 {
	res := int64(1)
	for i := x; i > 0; i -= i & -i {
		res = res * f.tree[i] % MOD
	}
	return res
}

func modPow(a int64, e int) int64 {
	res := int64(1)
	for e > 0 {
		if e&1 != 0 {
			res = res * a % MOD
		}
		a = a * a % MOD
		e >>= 1
	}
	return res
}

func modInv(a int64) int64 {
	return modPow(a, MOD-2)
}

type test struct {
	input    string
	expected string
}

func solve(input string) string {
	r := strings.NewReader(input)
	var n int
	fmt.Fscan(r, &n)
	a := make([]int, n+1)
	maxA := 0
	for i := 1; i <= n; i++ {
		fmt.Fscan(r, &a[i])
		if a[i] > maxA {
			maxA = a[i]
		}
	}
	spf := make([]int, maxA+1)
	for i := 2; i <= maxA; i++ {
		if spf[i] == 0 {
			for j := i; j <= maxA; j += i {
				if spf[j] == 0 {
					spf[j] = i
				}
			}
		}
	}
	type pe struct{ pos, exp int }
	stacks := make([][]pe, maxA+1)
	for p := 2; p <= maxA; p++ {
		if spf[p] == p {
			stacks[p] = []pe{{0, 0}}
		}
	}
	bit := NewFenwick(n)
	for i := 1; i <= n; i++ {
		x := a[i]
		for x > 1 {
			p := spf[x]
			cnt := 0
			for x%p == 0 {
				x /= p
				cnt++
			}
			if stacks[p] == nil {
				stacks[p] = []pe{{0, 0}}
			}
			stk := stacks[p]
			for len(stk) > 1 && stk[len(stk)-1].exp <= cnt {
				last := stk[len(stk)-1]
				stk = stk[:len(stk)-1]
				prev := stk[len(stk)-1]
				pe_k := modPow(int64(p), last.exp)
				inv_pe_k := modInv(pe_k)
				l := prev.pos + 1
				r := last.pos
				bit.rangeMul(l, r, inv_pe_k, pe_k)
			}
			prev := stk[len(stk)-1]
			pe_v := modPow(int64(p), cnt)
			inv_pe_v := modInv(pe_v)
			l := prev.pos + 1
			r := i
			bit.rangeMul(l, r, pe_v, inv_pe_v)
			stk = append(stk, pe{i, cnt})
			stacks[p] = stk
		}
	}
	var q int
	fmt.Fscan(r, &q)
	last := int64(0)
	var out strings.Builder
	for qi := 0; qi < q; qi++ {
		var x, y int
		fmt.Fscan(r, &x, &y)
		l := (last+int64(x))%int64(n) + 1
		rpos := (last+int64(y))%int64(n) + 1
		if l > rpos {
			l, rpos = rpos, l
		}
		prer := bit.query(int(rpos))
		prel := bit.query(int(l) - 1)
		invPrel := modInv(prel)
		ans := prer * invPrel % MOD
		out.WriteString(fmt.Sprint(ans))
		if qi+1 < q {
			out.WriteByte(' ')
		}
		last = ans
	}
	return out.String()
}

func generateTests() []test {
	rng := rand.New(rand.NewSource(47))
	tests := []test{}
	for len(tests) < 100 {
		n := rng.Intn(5) + 1
		q := rng.Intn(5) + 1
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", n)
		for i := 1; i <= n; i++ {
			if i > 1 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", rng.Intn(10)+1)
		}
		sb.WriteByte('\n')
		fmt.Fprintf(&sb, "%d\n", q)
		for i := 0; i < q; i++ {
			x := rng.Intn(n)
			y := rng.Intn(n)
			fmt.Fprintf(&sb, "%d %d\n", x, y)
		}
		inp := sb.String()
		tests = append(tests, test{inp, solve(inp)})
	}
	return tests
}

func runBinary(bin, input string) (string, error) {
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
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		got, err := runBinary(bin, t.input)
		if err != nil {
			fmt.Printf("Runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(t.expected) {
			fmt.Printf("Wrong answer on test %d\nInput:\n%sExpected:%s\nGot:%s\n", i+1, t.input, t.expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

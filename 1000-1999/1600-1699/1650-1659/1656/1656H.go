package main

import (
	"bufio"
	"fmt"
	"math/big"
	"os"
)

func gcdBig(a, b *big.Int) *big.Int {
	var g big.Int
	g.GCD(nil, nil, a, b)
	return &g
}

func lcmBig(a, b *big.Int) *big.Int {
	g := gcdBig(a, b)
	if g.Sign() == 0 {
		return big.NewInt(0)
	}
	t := new(big.Int).Div(new(big.Int).Set(a), g)
	t.Mul(t, b)
	if t.Sign() < 0 {
		t.Neg(t)
	}
	return t
}

func lcmSlice(arr []*big.Int) *big.Int {
	res := big.NewInt(1)
	for _, v := range arr {
		res = lcmBig(res, v)
	}
	return res
}

func cloneSlice(src []*big.Int) []*big.Int {
	dst := make([]*big.Int, len(src))
	copy(dst, src)
	return dst
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(in, &n, &m)
		a := make([]*big.Int, n)
		for i := 0; i < n; i++ {
			var s string
			fmt.Fscan(in, &s)
			val, _ := new(big.Int).SetString(s, 10)
			a[i] = val
		}
		b := make([]*big.Int, m)
		for i := 0; i < m; i++ {
			var s string
			fmt.Fscan(in, &s)
			val, _ := new(big.Int).SetString(s, 10)
			b[i] = val
		}

		setA := cloneSlice(a)
		setB := cloneSlice(b)

		for {
			if len(setA) == 0 || len(setB) == 0 {
				break
			}
			lA := lcmSlice(setA)
			lB := lcmSlice(setB)
			changed := false
			// filter A
			tmpA := setA[:0]
			for _, x := range setA {
				if gcdBig(x, lB).Cmp(x) == 0 {
					tmpA = append(tmpA, x)
				} else {
					changed = true
				}
			}
			setA = cloneSlice(tmpA)
			if len(setA) == 0 {
				break
			}
			// filter B
			lA = lcmSlice(setA)
			tmpB := setB[:0]
			for _, x := range setB {
				if gcdBig(x, lA).Cmp(x) == 0 {
					tmpB = append(tmpB, x)
				} else {
					changed = true
				}
			}
			setB = cloneSlice(tmpB)
			if !changed {
				break
			}
		}

		if len(setA) > 0 && len(setB) > 0 {
			lA := lcmSlice(setA)
			lB := lcmSlice(setB)
			if lA.Cmp(lB) == 0 {
				fmt.Fprintln(out, "YES")
				fmt.Fprintf(out, "%d %d\n", len(setA), len(setB))
				for i, x := range setA {
					if i > 0 {
						out.WriteByte(' ')
					}
					out.WriteString(x.String())
				}
				out.WriteByte('\n')
				for i, x := range setB {
					if i > 0 {
						out.WriteByte(' ')
					}
					out.WriteString(x.String())
				}
				out.WriteByte('\n')
				continue
			}
		}
		fmt.Fprintln(out, "NO")
	}
}

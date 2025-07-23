package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1000000007
const M = 1 << 16

var varMask [256]int

func initMasks() {
	for _, ch := range []byte{'A', 'B', 'C', 'D', 'a', 'b', 'c', 'd'} {
		m := 0
		for j := 0; j < 16; j++ {
			A := (j >> 3) & 1
			B := (j >> 2) & 1
			C := (j >> 1) & 1
			D := j & 1
			val := 0
			switch ch {
			case 'A':
				val = A
			case 'B':
				val = B
			case 'C':
				val = C
			case 'D':
				val = D
			case 'a':
				val = 1 - A
			case 'b':
				val = 1 - B
			case 'c':
				val = 1 - C
			case 'd':
				val = 1 - D
			}
			if val == 1 {
				m |= 1 << j
			}
		}
		varMask[ch] = m
	}
}

type Node struct {
	leaf        bool
	ch          byte
	op          byte
	left, right *Node
}

func parseExpr(s []byte, p *int) *Node {
	if s[*p] == '(' {
		(*p)++ // '('
		left := parseExpr(s, p)
		(*p)++ // ')'
		op := s[*p]
		(*p)++
		(*p)++ // '('
		right := parseExpr(s, p)
		(*p)++ // ')'
		return &Node{op: op, left: left, right: right}
	}
	ch := s[*p]
	(*p)++
	return &Node{leaf: true, ch: ch}
}

func fwtAnd(a []int64, inv bool) {
	for bit := 0; bit < 16; bit++ {
		for mask := 0; mask < M; mask++ {
			if mask&(1<<bit) != 0 {
				if !inv {
					a[mask] = (a[mask] + a[mask^(1<<bit)]) % MOD
				} else {
					a[mask] = (a[mask] - a[mask^(1<<bit)] + MOD) % MOD
				}
			}
		}
	}
}

func fwtOr(a []int64, inv bool) {
	for bit := 0; bit < 16; bit++ {
		for mask := 0; mask < M; mask++ {
			if mask&(1<<bit) == 0 {
				if !inv {
					a[mask] = (a[mask] + a[mask|(1<<bit)]) % MOD
				} else {
					a[mask] = (a[mask] - a[mask|(1<<bit)] + MOD) % MOD
				}
			}
		}
	}
}

func eval(node *Node) []int64 {
	if node.leaf {
		arr := make([]int64, M)
		if node.ch == '?' {
			for _, c := range []byte{'A', 'B', 'C', 'D', 'a', 'b', 'c', 'd'} {
				m := varMask[c]
				arr[m] = (arr[m] + 1) % MOD
			}
		} else {
			m := varMask[node.ch]
			arr[m] = 1
		}
		return arr
	}
	left := eval(node.left)
	right := eval(node.right)
	res := make([]int64, M)
	if node.op == '&' || node.op == '?' {
		A := make([]int64, M)
		B := make([]int64, M)
		copy(A, left)
		copy(B, right)
		fwtAnd(A, false)
		fwtAnd(B, false)
		for i := 0; i < M; i++ {
			A[i] = A[i] * B[i] % MOD
		}
		fwtAnd(A, true)
		for i := 0; i < M; i++ {
			res[i] = (res[i] + A[i]) % MOD
		}
	}
	if node.op == '|' || node.op == '?' {
		A := make([]int64, M)
		B := make([]int64, M)
		copy(A, left)
		copy(B, right)
		fwtOr(A, false)
		fwtOr(B, false)
		for i := 0; i < M; i++ {
			A[i] = A[i] * B[i] % MOD
		}
		fwtOr(A, true)
		for i := 0; i < M; i++ {
			res[i] = (res[i] + A[i]) % MOD
		}
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var s string
	if _, err := fmt.Fscan(in, &s); err != nil {
		return
	}
	var n int
	fmt.Fscan(in, &n)
	qs := make([]struct{ idx, val int }, n)
	for i := 0; i < n; i++ {
		var a, b, c, d, e int
		fmt.Fscan(in, &a, &b, &c, &d, &e)
		idx := (a << 3) | (b << 2) | (c << 1) | d
		qs[i] = struct{ idx, val int }{idx, e}
	}
	initMasks()
	b := []byte(s)
	pos := 0
	root := parseExpr(b, &pos)
	dp := eval(root)
	ans := int64(0)
	for mask := 0; mask < M; mask++ {
		if dp[mask] == 0 {
			continue
		}
		ok := true
		for _, q := range qs {
			bit := (mask >> q.idx) & 1
			if int(bit) != q.val {
				ok = false
				break
			}
		}
		if ok {
			ans = (ans + dp[mask]) % MOD
		}
	}
	fmt.Fprintln(out, ans)
}

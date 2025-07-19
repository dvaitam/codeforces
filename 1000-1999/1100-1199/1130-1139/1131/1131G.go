package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

var (
	reader = bufio.NewReader(os.Stdin)
	writer = bufio.NewWriter(os.Stdout)
)

func readInt() (int, error) {
	sign := 1
	var x int
	for {
		ch, _, err := reader.ReadRune()
		if err != nil {
			return 0, err
		}
		if ch == '-' {
			sign = -1
			break
		}
		if ch >= '0' && ch <= '9' {
			x = int(ch - '0')
			break
		}
	}
	for {
		ch, _, err := reader.ReadRune()
		if err != nil && err != io.EOF {
			return 0, err
		}
		if ch < '0' || ch > '9' {
			break
		}
		x = x*10 + int(ch-'0')
	}
	return x * sign, nil
}

func main() {
	defer writer.Flush()
	B, err := readInt()
	if err != nil {
		return
	}
	n, _ := readInt()
	// read blocks
	type block struct {
		hs []int
		cs []int64
	}
	blocks := make([]block, B)
	for i := 0; i < B; i++ {
		L, _ := readInt()
		bs := block{hs: make([]int, L), cs: make([]int64, L)}
		for j := 0; j < L; j++ {
			v, _ := readInt()
			bs.hs[j] = v
		}
		for j := 0; j < L; j++ {
			v, _ := readInt()
			bs.cs[j] = int64(v)
		}
		blocks[i] = bs
	}
	q, _ := readInt()
	h := make([]int, n)
	c := make([]int64, n)
	cnt := 0
	for i := 0; i < q; i++ {
		id, _ := readInt()
		mul, _ := readInt()
		id--
		bs := blocks[id]
		L := len(bs.hs)
		for j := 0; j < L; j++ {
			h[cnt+j] = bs.hs[j]
			c[cnt+j] = bs.cs[j] * int64(mul)
		}
		cnt += L
	}
	// prepare arrays
	lb := make([]int, n)
	rb := make([]int, n)
	stk := make([]int, 0, n)
	// compute lb
	for i := 0; i < n; i++ {
		L := i - h[i] + 1
		if L < 0 {
			L = 0
		}
		for len(stk) > 0 {
			j := stk[len(stk)-1]
			if j-h[j]+1 < L {
				break
			}
			stk = stk[:len(stk)-1]
		}
		if len(stk) > 0 && stk[len(stk)-1] >= L {
			prev := stk[len(stk)-1]
			// lb[i] = min(lb[prev], L)
			x := lb[prev]
			if L < x {
				x = L
			}
			lb[i] = x
		} else {
			lb[i] = L
		}
		stk = append(stk, i)
	}
	// compute rb
	stk = stk[:0]
	for i := n - 1; i >= 0; i-- {
		R := i + h[i] - 1
		if R >= n {
			R = n - 1
		}
		for len(stk) > 0 {
			j := stk[len(stk)-1]
			if j+h[j]-1 > R {
				break
			}
			stk = stk[:len(stk)-1]
		}
		if len(stk) > 0 && stk[len(stk)-1] <= R {
			prev := stk[len(stk)-1]
			x := rb[prev]
			if R < x {
				x = R
			}
			rb[i] = x
		} else {
			rb[i] = R
		}
		stk = append(stk, i)
	}
	// dp
	const INF = int64(9e18)
	dp := make([]int64, n)
	stk = stk[:0]
	for i := 0; i < n; i++ {
		dp[i] = INF
		// pop expired
		for len(stk) > 0 && rb[stk[len(stk)-1]] < i {
			stk = stk[:len(stk)-1]
		}
		var d int
		if len(stk) > 0 {
			d = stk[len(stk)-1]
		} else {
			d = -1
		}
		if d >= 0 && rb[d] >= i {
			var prev int64
			if d > 0 {
				prev = dp[d-1]
			}
			v := prev + c[d]
			if v < dp[i] {
				dp[i] = v
			}
		}
		// cover to left
		var prevL int64
		if lb[i] > 0 {
			prevL = dp[lb[i]-1]
		}
		v2 := prevL + c[i]
		if v2 < dp[i] {
			dp[i] = v2
		}
		// consider extend
		var prevAll int64
		if i > 0 {
			prevAll = dp[i-1]
		}
		C := prevAll + c[i]
		if len(stk) > 0 && rb[d] == rb[i] {
			var alt int64
			if d > 0 {
				alt = dp[d-1]
			}
			alt += c[d]
			if C < alt {
				stk[len(stk)-1] = i
			}
		} else if len(stk) == 0 {
			stk = append(stk, i)
		} else {
			var alt int64
			if d > 0 {
				alt = dp[d-1]
			}
			alt += c[d]
			if C < alt {
				stk = append(stk, i)
			}
		}
	}
	// output
	if n > 0 {
		fmt.Fprintln(writer, dp[n-1])
	}
}

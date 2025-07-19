package main

import (
	"bufio"
	"fmt"
	"os"
)

var rdr = bufio.NewReader(os.Stdin)
var wtr = bufio.NewWriter(os.Stdout)
var m int

func readInt() int {
	b, _ := rdr.ReadByte()
	for (b < '0' || b > '9') && b != '-' {
		b, _ = rdr.ReadByte()
	}
	sign := 1
	if b == '-' {
		sign = -1
		b, _ = rdr.ReadByte()
	}
	n := 0
	for b >= '0' && b <= '9' {
		n = n*10 + int(b-'0')
		b, _ = rdr.ReadByte()
	}
	return sign * n
}

func dis(x, y int) int {
	if x < y {
		return y - x
	}
	return m - x + y
}

func find(pre []int, x int) int {
	if pre[x] != x {
		pre[x] = find(pre, pre[x])
	}
	return pre[x]
}

func main() {
	defer wtr.Flush()
	n := readInt()
	m = readInt()
	a := make([]int, n)
	cnt := make([]int, m)
	p := make([]int, m)
	nxt := make([]int, n)
	for i := 0; i < m; i++ {
		p[i] = -1
	}
	for i := 0; i < n; i++ {
		a[i] = readInt()
		r := a[i] % m
		if r < 0 {
			r += m
		}
		cnt[r]++
		nxt[i] = p[r]
		p[r] = i
	}
	d := n / m
	pre := make([]int, m)
	tmp := 0
	for ; tmp < m; tmp++ {
		if cnt[tmp] > d {
			break
		}
	}
	// initialize pre array
	pre[tmp] = tmp
	for i := tmp + 1; i < m; i++ {
		if cnt[i] > d {
			pre[i] = i
		} else {
			pre[i] = pre[i-1]
		}
	}
	if tmp > 0 {
		if cnt[0] > d {
			pre[0] = 0
		} else {
			pre[0] = pre[m-1]
		}
		for i := 1; i < tmp; i++ {
			if cnt[i] > d {
				pre[i] = i
			} else {
				pre[i] = pre[i-1]
			}
		}
	}
	var ans int64 = 0
	// balance
	for i := 0; i < m; i++ {
		for cnt[i] < d {
			j := find(pre, i)
			// move from j to i
			for cnt[j] > d && cnt[i] < d {
				k := p[j]
				p[j] = nxt[k]
				cnt[j]--
				cnt[i]++
				delta := dis(j, i)
				a[k] += delta
				ans += int64(delta)
			}
			if cnt[j] == d {
				// remove j
				prev := j - 1
				if prev < 0 {
					prev = m - 1
				}
				pre[j] = find(pre, prev)
			}
		}
	}
	fmt.Fprintln(wtr, ans)
	for i := 0; i < n; i++ {
		fmt.Fprintf(wtr, "%d", a[i])
		if i+1 < n {
			wtr.WriteByte(' ')
		}
	}
	wtr.WriteByte('\n')
}

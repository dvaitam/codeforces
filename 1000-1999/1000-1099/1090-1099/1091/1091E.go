package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	fmt.Fscan(reader, &n)
	D := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &D[i])
	}
	// sort descending
	sort.Slice(D[1:], func(i, j int) bool { return D[1+i] > D[1+j] })
	// prefix and suffix sums
	S1 := make([]int64, n+2)
	S2 := make([]int64, n+2)
	for i := 1; i <= n; i++ {
		S1[i] = S1[i-1] + int64(D[i])
	}
	for i := n; i >= 1; i-- {
		S2[i] = S2[i+1] + int64(D[i])
	}
	// T[i]: minimal index j such that D[j] <= i
	T := make([]int, n+2)
	j := n + 1
	for i := 0; i <= n; i++ {
		for j > 1 && D[j-1] <= i {
			j--
		}
		T[i] = j
	}
	// helper get
	get := func(i, k int) int64 {
		t := T[k]
		if i > t {
			t = i
		}
		// positions i..t-1 contribute k, t..n contribute D
		return int64(t-i)*int64(k) + S2[t]
	}
	// P1 and P2
	P1 := make([]int64, n+2)
	P2 := make([]int64, n+2)
	for i := 1; i <= n; i++ {
		pi1 := S1[i] - int64(i)*(int64(i)-1) - get(i+1, i)
		if pi1 > int64(i) {
			P1[i] = int64(n) + 1
		} else {
			P1[i] = pi1
		}
		P2[i] = int64(i+1)*int64(i) + get(i+1, i+1) - S1[i]
	}
	P1[0] = 0
	P2[n+1] = int64(n) + 1
	for i := 1; i <= n; i++ {
		if P1[i-1] > P1[i] {
			P1[i] = P1[i-1]
		}
	}
	for i := n; i >= 1; i-- {
		if P2[i+1] < P2[i] {
			P2[i] = P2[i+1]
		}
	}
	total := S1[n]
	var ok bool
	// collect results
	var res []int
	j = n + 1
	for i := int(total & 1); i <= n; i += 2 {
		for j > 1 && i >= D[j-1] {
			j--
		}
		cond1 := P1[j-1] <= int64(i)
		cond2 := P2[j] >= int64(i)
		cond3 := S1[j-1]+int64(i) <= int64(j)*(int64(j)-1)+get(j, j)
		if cond1 && cond2 && cond3 {
			res = append(res, i)
			ok = true
		}
	}
	if !ok {
		fmt.Fprint(writer, -1)
	} else {
		for idx, v := range res {
			if idx > 0 {
				writer.WriteByte(' ')
			}
			fmt.Fprint(writer, v)
		}
	}
}

package main

import (
	"bufio"
	"fmt"
	"os"
)

type pair struct {
	val int
	cnt int64
}

func canon(arr []int, m int) []pair {
	res := make([]pair, 0, len(arr))
	for _, x := range arr {
		v := x
		c := int64(1)
		for v%m == 0 {
			v /= m
			c *= int64(m)
		}
		if len(res) > 0 && res[len(res)-1].val == v {
			res[len(res)-1].cnt += c
		} else {
			res = append(res, pair{val: v, cnt: c})
		}
	}
	return res
}

func canTransform(aPairs, bPairs []pair) bool {
	i, j := 0, 0
	for i < len(aPairs) && j < len(bPairs) {
		if aPairs[i].val != bPairs[j].val {
			return false
		}
		if aPairs[i].cnt < bPairs[j].cnt {
			bPairs[j].cnt -= aPairs[i].cnt
			i++
		} else if aPairs[i].cnt > bPairs[j].cnt {
			aPairs[i].cnt -= bPairs[j].cnt
			j++
		} else {
			i++
			j++
		}
	}
	return i == len(aPairs) && j == len(bPairs)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	fmt.Fscan(reader, &T)
	for ; T > 0; T-- {
		var n, m int
		fmt.Fscan(reader, &n, &m)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		var k int
		fmt.Fscan(reader, &k)
		b := make([]int, k)
		for i := 0; i < k; i++ {
			fmt.Fscan(reader, &b[i])
		}
		aPairs := canon(a, m)
		bPairs := canon(b, m)
		if canTransform(aPairs, bPairs) {
			fmt.Fprintln(writer, "Yes")
		} else {
			fmt.Fprintln(writer, "No")
		}
	}
}

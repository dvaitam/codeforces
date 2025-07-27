package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func exgcd(a, b int64) (int64, int64, int64) {
	if b == 0 {
		return a, 1, 0
	}
	g, x1, y1 := exgcd(b, a%b)
	return g, y1, x1 - (a/b)*y1
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int
	var k int64
	if _, err := fmt.Fscan(reader, &n, &m, &k); err != nil {
		return
	}
	maxVal := n
	if m > maxVal {
		maxVal = m
	}
	maxVal = 2*maxVal + 5
	posA := make([]int, maxVal)
	posB := make([]int, maxVal)
	for i := range posA {
		posA[i] = -1
		posB[i] = -1
	}
	for i := 0; i < n; i++ {
		var x int
		fmt.Fscan(reader, &x)
		if x < len(posA) {
			posA[x] = i
		}
	}
	for i := 0; i < m; i++ {
		var x int
		fmt.Fscan(reader, &x)
		if x < len(posB) {
			posB[x] = i
		}
	}

	n1 := int64(n)
	m1 := int64(m)
	g := gcd(n1, m1)
	lcm := n1 / g * m1
	inv := int64(0)
	if g == 1 {
		_, invTmp, _ := exgcd(n1, m1)
		inv = (invTmp%m1 + m1) % m1
	} else {
		_, invTmp, _ := exgcd(n1/g, m1/g)
		inv = (invTmp%(m1/g) + m1/g) % (m1 / g)
	}

	matches := make([]int64, 0)
	for c := 0; c < maxVal; c++ {
		ia := posA[c]
		ib := posB[c]
		if ia == -1 || ib == -1 {
			continue
		}
		i64 := int64(ia)
		j64 := int64(ib)
		if (i64-j64)%g != 0 {
			continue
		}
		var t int64
		if g == 1 {
			diff := (j64 - i64) % m1
			if diff < 0 {
				diff += m1
			}
			t = (i64 + n1*((diff*inv)%m1)) % lcm
		} else {
			mg := m1 / g
			diff := ((j64 - i64) / g) % mg
			if diff < 0 {
				diff += mg
			}
			t = (i64 + n1*((diff*inv)%mg)) % lcm
		}
		matches = append(matches, t+1)
	}
	sort.Slice(matches, func(i, j int) bool { return matches[i] < matches[j] })

	diffsPerCycle := lcm - int64(len(matches))
	cycles := int64(0)
	if diffsPerCycle > 0 {
		cycles = (k - 1) / diffsPerCycle
	}
	k -= cycles * diffsPerCycle
	base := cycles * lcm
	// binary search within one cycle
	left, right := int64(1), lcm
	for left < right {
		mid := (left + right) / 2
		eq := int64(sort.Search(len(matches), func(i int) bool { return matches[i] > mid }))
		diff := mid - eq
		if diff >= k {
			right = mid
		} else {
			left = mid + 1
		}
	}
	fmt.Fprintln(writer, base+left)
}

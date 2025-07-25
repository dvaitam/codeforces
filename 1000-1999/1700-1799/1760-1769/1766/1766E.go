package main

import (
	"bufio"
	"fmt"
	"os"
)

const stateCount = 16

var states = [stateCount]int{0, 1, 2, 3, 5, 6, 7, 9, 10, 11, 25, 26, 27, 37, 38, 39}
var indexMap = map[int]int{}
var trans [stateCount][4]int
var inc [stateCount][4]int

func encode(m1, m2, m3 int) int {
	return m1 | (m2 << 2) | (m3 << 4)
}

func appendState(state, x int) (int, int) {
	m1 := state & 3
	m2 := (state >> 2) & 3
	m3 := (state >> 4) & 3
	if x == 0 {
		return state, 1
	}
	if m1 != 0 && (m1&x) != 0 {
		m1 = x
		return encode(m1, m2, m3), 0
	}
	if m2 != 0 && (m2&x) != 0 {
		m2 = x
		return encode(m1, m2, m3), 0
	}
	if m3 != 0 && (m3&x) != 0 {
		m3 = x
		return encode(m1, m2, m3), 0
	}
	if m1 == 0 {
		m1 = x
		return encode(m1, m2, m3), 1
	}
	if m2 == 0 {
		m2 = x
		return encode(m1, m2, m3), 1
	}
	if m3 == 0 {
		m3 = x
		return encode(m1, m2, m3), 1
	}
	// unreachable for given constraints
	return state, 1
}

func init() {
	for i, s := range states {
		indexMap[s] = i
	}
	for i, s := range states {
		for x := 0; x < 4; x++ {
			ns, d := appendState(s, x)
			trans[i][x] = indexMap[ns]
			inc[i][x] = d
		}
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &arr[i])
	}

	var cnt [stateCount]int64
	var sum [stateCount]int64
	var ans int64

	for _, v := range arr {
		var nextCnt [stateCount]int64
		var nextSum [stateCount]int64

		ns := trans[0][v]
		d := inc[0][v]
		nextCnt[ns] += 1
		nextSum[ns] += int64(d)

		for i := 0; i < stateCount; i++ {
			if cnt[i] == 0 {
				continue
			}
			ns := trans[i][v]
			dd := inc[i][v]
			nextCnt[ns] += cnt[i]
			nextSum[ns] += sum[i] + int64(dd)*cnt[i]
		}

		var partial int64
		for i := 0; i < stateCount; i++ {
			partial += nextSum[i]
		}
		ans += partial
		cnt = nextCnt
		sum = nextSum
	}

	fmt.Fprintln(out, ans)
}

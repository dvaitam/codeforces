package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

const mod int64 = 1000000007

var (
	k           int
	opp         = []int{1, 0, 3, 2, 5, 4}
	dpGeneric   [][]int64
	setNodes    map[int64]bool
	presetColor map[int64]int
	memo        map[int64][6]int64
)

func depth(v int64) int {
	return bits.Len64(uint64(v))
}

func dfs(v int64) [6]int64 {
	if val, ok := memo[v]; ok {
		return val
	}
	h := k - depth(v) + 1
	var res [6]int64
	if !setNodes[v] {
		for i := 0; i < 6; i++ {
			res[i] = dpGeneric[h][i]
		}
		return res
	}
	if h == 1 {
		for i := 0; i < 6; i++ {
			if c, ok := presetColor[v]; ok && c != i {
				res[i] = 0
			} else {
				res[i] = 1
			}
		}
		memo[v] = res
		return res
	}
	var leftArr, rightArr [6]int64
	for i := 0; i < 6; i++ {
		leftArr[i] = dpGeneric[h-1][i]
		rightArr[i] = dpGeneric[h-1][i]
	}
	leftID := v * 2
	rightID := v*2 + 1
	if setNodes[leftID] {
		leftArr = dfs(leftID)
	}
	if setNodes[rightID] {
		rightArr = dfs(rightID)
	}
	totalLeft := int64(0)
	totalRight := int64(0)
	for i := 0; i < 6; i++ {
		totalLeft += leftArr[i]
		if totalLeft >= mod {
			totalLeft -= mod
		}
		totalRight += rightArr[i]
		if totalRight >= mod {
			totalRight -= mod
		}
	}
	for i := 0; i < 6; i++ {
		if c, ok := presetColor[v]; ok && c != i {
			res[i] = 0
			continue
		}
		allowedLeft := (totalLeft - leftArr[i] - leftArr[opp[i]]) % mod
		if allowedLeft < 0 {
			allowedLeft += mod
		}
		allowedRight := (totalRight - rightArr[i] - rightArr[opp[i]]) % mod
		if allowedRight < 0 {
			allowedRight += mod
		}
		res[i] = allowedLeft * allowedRight % mod
	}
	memo[v] = res
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	fmt.Fscan(reader, &k)
	fmt.Fscan(reader, &n)

	presetColor = make(map[int64]int)
	setNodes = make(map[int64]bool)
	setNodes[1] = true
	for i := 0; i < n; i++ {
		var v int64
		var s string
		fmt.Fscan(reader, &v, &s)
		col := 0
		switch s {
		case "white":
			col = 0
		case "yellow":
			col = 1
		case "green":
			col = 2
		case "blue":
			col = 3
		case "red":
			col = 4
		case "orange":
			col = 5
		}
		presetColor[v] = col
		for cur := v; cur >= 1; cur /= 2 {
			setNodes[cur] = true
			if cur == 1 {
				break
			}
		}
	}

	dpGeneric = make([][]int64, k+1)
	dpGeneric[1] = make([]int64, 6)
	for i := 0; i < 6; i++ {
		dpGeneric[1][i] = 1
	}
	for h := 2; h <= k; h++ {
		dpGeneric[h] = make([]int64, 6)
		total := int64(0)
		for i := 0; i < 6; i++ {
			total += dpGeneric[h-1][i]
			if total >= mod {
				total -= mod
			}
		}
		for i := 0; i < 6; i++ {
			allowed := (total - dpGeneric[h-1][i] - dpGeneric[h-1][opp[i]]) % mod
			if allowed < 0 {
				allowed += mod
			}
			dpGeneric[h][i] = allowed * allowed % mod
		}
	}

	memo = make(map[int64][6]int64)
	rootArr := dfs(1)
	ans := int64(0)
	for i := 0; i < 6; i++ {
		ans += rootArr[i]
		ans %= mod
	}
	fmt.Fprintln(writer, ans)
}

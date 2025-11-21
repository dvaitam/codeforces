package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var k int
	var t int64
	var s string
	if _, err := fmt.Fscan(in, &k, &t); err != nil {
		return
	}
	if _, err := fmt.Fscan(in, &s); err != nil {
		return
	}
	n := 1 << k
	bytes := []byte(s)
	if k == 0 {
		fmt.Println(s)
		return
	}
	tMod := int(t % int64(k))
	if tMod < 0 {
		tMod += k
	}
	R := 1 << tMod
	C := n / R
	result := solveColumnMajor(bytes, R, C)
	fmt.Println(string(result))
}

type colCandidate struct {
	q int
	r int
}

type columnInfo struct {
	doubled []byte
	prefix  []uint64
}

const hashBase uint64 = 1315423911

func solveColumnMajor(data []byte, R, C int) []byte {
	columns := make([]columnInfo, C)
	pow := make([]uint64, 2*R+1)
	pow[0] = 1
	for i := 1; i < len(pow); i++ {
		pow[i] = pow[i-1] * hashBase
	}
	for col := 0; col < C; col++ {
		colData := make([]byte, R)
		for row := 0; row < R; row++ {
			colData[row] = data[row*C+col]
		}
		dbl := make([]byte, 2*R)
		copy(dbl, colData)
		copy(dbl[R:], colData)
		prefix := make([]uint64, len(dbl)+1)
		for i := 0; i < len(dbl); i++ {
			prefix[i+1] = prefix[i]*hashBase + uint64(dbl[i])
		}
		columns[col] = columnInfo{doubled: dbl, prefix: prefix}
	}
	bestSet := false
	var best colCandidate
	for q := 0; q < R; q++ {
		r := bestStartColumns(columns, pow, R, C, q)
		if !bestSet {
			bestSet = true
			best = colCandidate{q: q, r: r}
			continue
		}
		if compareColumnCandidates(columns, pow, R, C, q, r, best.q, best.r) < 0 {
			best = colCandidate{q: q, r: r}
		}
	}
	res := make([]byte, 0, R*C)
	for col := best.r; col < C; col++ {
		res = append(res, getColumnRotation(columns[col], R, best.q)...)
	}
	nextShift := (best.q + 1) % R
	for col := 0; col < best.r; col++ {
		res = append(res, getColumnRotation(columns[col], R, nextShift)...)
	}
	return res
}

func getColumnRotation(info columnInfo, R, shift int) []byte {
	if R == 0 {
		return nil
	}
	return info.doubled[shift : shift+R]
}

func bestStartColumns(columns []columnInfo, pow []uint64, R, C, q int) int {
	nextShift := (q + 1) % R
	i, j, k := 0, 1, 0
	for i < C && j < C && k < C {
		col1, shift1 := columnToken(C, q, nextShift, i+k)
		col2, shift2 := columnToken(C, q, nextShift, j+k)
		cmp := compareColumns(columns, pow, R, col1, shift1, col2, shift2)
		if cmp == 0 {
			k++
		} else if cmp < 0 {
			j += k + 1
			if j == i {
				j++
			}
			k = 0
		} else {
			i += k + 1
			if i == j {
				i++
			}
			k = 0
		}
	}
	if i >= C {
		return j % C
	}
	if j >= C {
		return i % C
	}
	if i < j {
		return i
	}
	return j
}

func columnToken(C, q, nextShift, pos int) (int, int) {
	if pos < C {
		return pos, q
	}
	return pos - C, nextShift
}

func compareColumns(columns []columnInfo, pow []uint64, R, col1, shift1, col2, shift2 int) int {
	if col1 == col2 && shift1 == shift2 {
		return 0
	}
	l := longestCommonPrefix(columns[col1], columns[col2], shift1, shift2, R, pow)
	if l == R {
		return 0
	}
	c1 := columns[col1].doubled[shift1+l]
	c2 := columns[col2].doubled[shift2+l]
	if c1 < c2 {
		return -1
	}
	return 1
}

func longestCommonPrefix(a, b columnInfo, startA, startB, limit int, pow []uint64) int {
	low, high := 0, limit
	for low < high {
		mid := (low + high + 1) >> 1
		if columnHash(a, startA, mid, pow) == columnHash(b, startB, mid, pow) {
			low = mid
		} else {
			high = mid - 1
		}
	}
	return low
}

func columnHash(info columnInfo, start, length int, pow []uint64) uint64 {
	end := start + length
	return info.prefix[end] - info.prefix[start]*pow[length]
}

func compareColumnCandidates(columns []columnInfo, pow []uint64, R, C, q1, r1, q2, r2 int) int {
	part1 := C - r1
	part2 := C - r2
	shift1Wrap := (q1 + 1) % R
	shift2Wrap := (q2 + 1) % R
	for idx := 0; idx < C; idx++ {
		var col1, shift1 int
		if idx < part1 {
			col1 = r1 + idx
			shift1 = q1
		} else {
			col1 = idx - part1
			shift1 = shift1Wrap
		}
		var col2, shift2 int
		if idx < part2 {
			col2 = r2 + idx
			shift2 = q2
		} else {
			col2 = idx - part2
			shift2 = shift2Wrap
		}
		cmp := compareColumns(columns, pow, R, col1, shift1, col2, shift2)
		if cmp != 0 {
			return cmp
		}
	}
	return 0
}

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

type dpState struct {
	bits   [][]uint64
	active []int
	used   []bool
}

func ensureRow(bitsets [][]uint64, idx int, wordLen int) []uint64 {
	row := bitsets[idx]
	if row == nil {
		row = make([]uint64, wordLen)
		bitsets[idx] = row
	}
	return row
}

func setBit(bitsets [][]uint64, used []bool, active *[]int, prod int, sum int, wordLen int) {
	row := bitsets[prod]
	if row == nil {
		row = make([]uint64, wordLen)
		bitsets[prod] = row
	}
	word := sum >> 6
	mask := uint64(1) << (uint(sum) & 63)
	if row[word]&mask != 0 {
		return
	}
	row[word] |= mask
	if !used[prod] {
		used[prod] = true
		*active = append(*active, prod)
	}
}

func clearRows(bitsets [][]uint64, active []int, used []bool) {
	for _, prod := range active {
		row := bitsets[prod]
		if row == nil {
			continue
		}
		for i := range row {
			row[i] = 0
		}
		used[prod] = false
	}
}

func iterateRow(row []uint64, limit int, fn func(sum int) bool) bool {
	for idx, word := range row {
		w := word
		for w != 0 {
			lsb := w & -w
			bit := bits.TrailingZeros64(w)
			sum := (idx << 6) + bit
			if sum <= limit {
				if fn(sum) {
					return true
				}
			}
			w -= lsb
		}
	}
	return false
}

func capSingle(x int, m int, capVal int) int {
	if x == 0 {
		return 0
	}
	if x > m {
		return capVal
	}
	return x
}

func capMultiply(prod int, x int, m int, capVal int) int {
	if x == 0 || prod == 0 {
		return 0
	}
	if prod == capVal {
		return capVal
	}
	val := int64(prod) * int64(x)
	if val > int64(m) {
		return capVal
	}
	return int(val)
}

func solveCase(a []int, m int) bool {
	n := len(a)
	if n == 0 {
		return false
	}
	prodCap := m + 1
	prodSize := prodCap + 1
	wordLen := ((m + 1) + 63) >> 6

	currBits := make([][]uint64, prodSize)
	nextBits := make([][]uint64, prodSize)
	currActive := make([]int, 0, prodSize)
	nextActive := make([]int, 0, prodSize)
	currUsed := make([]bool, prodSize)
	nextUsed := make([]bool, prodSize)

	firstProd := capSingle(a[0], m, prodCap)
	setBit(currBits, currUsed, &currActive, firstProd, 0, wordLen)

	for i := 1; i < n; i++ {
		x := a[i]
		clearRows(nextBits, nextActive, nextUsed)
		nextActive = nextActive[:0]

		for _, prod := range currActive {
			row := currBits[prod]
			if row == nil {
				continue
			}
			iterateRow(row, m, func(sum int) bool {
				// multiply
				newProd := capMultiply(prod, x, m, prodCap)
				setBit(nextBits, nextUsed, &nextActive, newProd, sum, wordLen)

				// plus
				if prod != prodCap {
					newSum := sum + prod
					if newSum <= m {
						prod2 := capSingle(x, m, prodCap)
						setBit(nextBits, nextUsed, &nextActive, prod2, newSum, wordLen)
					}
				}
				return false
			})
		}

		currBits, nextBits = nextBits, currBits
		currActive, nextActive = nextActive, currActive
		currUsed, nextUsed = nextUsed, currUsed
	}

	target := m
	for _, prod := range currActive {
		if prod == prodCap {
			continue
		}
		row := currBits[prod]
		if row == nil {
			continue
		}
		if iterateRow(row, m, func(sum int) bool {
			if sum+prod == target {
				return true
			}
			return false
		}) {
			return true
		}
	}
	return false
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
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		if solveCase(a, m) {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}

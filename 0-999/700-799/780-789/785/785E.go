package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
)

var (
	n, q      int
	arr       []int
	blocks    [][]int
	blockSize int
)

func build() {
	blockSize = int(math.Sqrt(float64(n))) + 1
	numBlocks := (n + blockSize - 1) / blockSize
	arr = make([]int, n+1)
	blocks = make([][]int, numBlocks)
	for i := 1; i <= n; i++ {
		arr[i] = i
	}
	for b := 0; b < numBlocks; b++ {
		start := b*blockSize + 1
		end := (b + 1) * blockSize
		if end > n {
			end = n
		}
		if start > end {
			blocks[b] = []int{}
			continue
		}
		blocks[b] = make([]int, end-start+1)
		for i := start; i <= end; i++ {
			blocks[b][i-start] = arr[i]
		}
		sort.Ints(blocks[b])
	}
}

func removeVal(idx, val int) {
	b := (idx - 1) / blockSize
	slice := blocks[b]
	pos := sort.Search(len(slice), func(i int) bool { return slice[i] >= val })
	if pos < len(slice) && slice[pos] == val {
		blocks[b] = append(slice[:pos], slice[pos+1:]...)
	}
	arr[idx] = 0
}

func addVal(idx, val int) {
	b := (idx - 1) / blockSize
	slice := blocks[b]
	pos := sort.Search(len(slice), func(i int) bool { return slice[i] >= val })
	slice = append(slice, 0)
	copy(slice[pos+1:], slice[pos:])
	slice[pos] = val
	blocks[b] = slice
	arr[idx] = val
}

func countLess(l, r, val int) int {
	if l > r {
		return 0
	}
	res := 0
	b1 := (l - 1) / blockSize
	b2 := (r - 1) / blockSize
	if b1 == b2 {
		for i := l; i <= r; i++ {
			if arr[i] < val {
				res++
			}
		}
		return res
	}
	end := (b1 + 1) * blockSize
	for i := l; i <= end; i++ {
		if arr[i] < val {
			res++
		}
	}
	for b := b1 + 1; b <= b2-1; b++ {
		slice := blocks[b]
		idx := sort.Search(len(slice), func(i int) bool { return slice[i] >= val })
		res += idx
	}
	start := b2 * blockSize
	for i := start + 1; i <= r; i++ {
		if arr[i] < val {
			res++
		}
	}
	return res
}

func countGreater(l, r, val int) int {
	if l > r {
		return 0
	}
	res := 0
	b1 := (l - 1) / blockSize
	b2 := (r - 1) / blockSize
	if b1 == b2 {
		for i := l; i <= r; i++ {
			if arr[i] > val {
				res++
			}
		}
		return res
	}
	end := (b1 + 1) * blockSize
	for i := l; i <= end; i++ {
		if arr[i] > val {
			res++
		}
	}
	for b := b1 + 1; b <= b2-1; b++ {
		slice := blocks[b]
		idx := sort.Search(len(slice), func(i int) bool { return slice[i] > val })
		res += len(slice) - idx
	}
	start := b2 * blockSize
	for i := start + 1; i <= r; i++ {
		if arr[i] > val {
			res++
		}
	}
	return res
}

func leftGreater(idx, val int) int { return countGreater(1, idx-1, val) }
func rightLess(idx, val int) int   { return countLess(idx+1, n, val) }

func main() {
	in := bufio.NewReader(os.Stdin)
	fmt.Fscan(in, &n, &q)
	build()
	inv := int64(0)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	for t := 0; t < q; t++ {
		var l, r int
		fmt.Fscan(in, &l, &r)
		if l == r {
			fmt.Fprintln(out, inv)
			continue
		}
		if l > r {
			l, r = r, l
		}
		a := arr[l]
		b := arr[r]

		inv -= int64(leftGreater(l, a) + rightLess(l, a))
		inv -= int64(leftGreater(r, b) + rightLess(r, b))
		if a > b {
			inv += 1
		}

		removeVal(l, a)
		removeVal(r, b)
		addVal(l, b)
		addVal(r, a)

		inv += int64(leftGreater(l, b) + rightLess(l, b))
		inv += int64(leftGreater(r, a) + rightLess(r, a))
		if b > a {
			inv -= 1
		}

		fmt.Fprintln(out, inv)
	}
}

package main

import (
	"bufio"
	"fmt"
	"os"
)

const inf int64 = 1 << 60
const blockSize = 120

type matrix [10][10]int64

type solver struct {
	n, k      int
	digits    []int
	blocks    int
	maxPow    int
	blockMats [][]matrix
}

func (sv *solver) applySegment(state [10]int64, l, r int) [10]int64 {
	if l > r {
		return state
	}
	length := r - l + 1
	dp := make([]int64, length+1)
	dp[0] = state[0]
	best := state
	best[0] = inf
	for rel := 1; rel <= length; rel++ {
		p := rel - 1
		if p >= 1 {
			digit := sv.digits[l+p-1]
			val := dp[p-1]
			if val < inf {
				cand := val - int64(digit)*int64(p)
				if cand < best[digit] {
					best[digit] = cand
				}
			}
		}
		q := rel
		val := inf
		for d := 1; d <= 9; d++ {
			if best[d] >= inf {
				continue
			}
			tmp := best[d] + int64(d)*int64(q)
			if tmp < val {
				val = tmp
			}
		}
		dp[rel] = val
	}
	length64 := int64(length)
	best[0] = dp[length]
	for d := 1; d <= 9; d++ {
		if best[d] < inf {
			best[d] += int64(d) * length64
		}
	}
	return best
}

func applyMatrix(state [10]int64, mat matrix) [10]int64 {
	var res [10]int64
	for j := 0; j < 10; j++ {
		val := inf
		for i := 0; i < 10; i++ {
			if state[i] >= inf || mat[i][j] >= inf {
				continue
			}
			tmp := state[i] + mat[i][j]
			if tmp < val {
				val = tmp
			}
		}
		res[j] = val
	}
	return res
}

func multiplyMatrix(a, b matrix) matrix {
	var res matrix
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			res[i][j] = inf
			for k := 0; k < 10; k++ {
				if a[i][k] >= inf || b[k][j] >= inf {
					continue
				}
				tmp := a[i][k] + b[k][j]
				if tmp < res[i][j] {
					res[i][j] = tmp
				}
			}
		}
	}
	return res
}

func (sv *solver) buildMatrix(l, r int) matrix {
	var mat matrix
	for idx := 0; idx < 10; idx++ {
		var state [10]int64
		for i := 0; i < 10; i++ {
			state[i] = inf
		}
		state[idx] = 0
		res := sv.applySegment(state, l, r)
		for j := 0; j < 10; j++ {
			mat[idx][j] = res[j]
		}
	}
	return mat
}

func (sv *solver) buildBlocks() {
	sv.blocks = (sv.n + blockSize - 1) / blockSize
	sv.blockMats = make([][]matrix, 1)
	sv.blockMats[0] = make([]matrix, sv.blocks)
	for i := 0; i < sv.blocks; i++ {
		l := i * blockSize
		r := (i+1)*blockSize - 1
		if r >= sv.n {
			r = sv.n - 1
		}
		sv.blockMats[0][i] = sv.buildMatrix(l, r)
	}
	sv.maxPow = 0
	for (1 << uint(sv.maxPow+1)) <= sv.blocks {
		sv.maxPow++
	}
	for p := 1; p <= sv.maxPow; p++ {
		size := sv.blocks - (1 << uint(p)) + 1
		row := make([]matrix, size)
		prev := sv.blockMats[p-1]
		for i := 0; i < size; i++ {
			row[i] = multiplyMatrix(prev[i], prev[i+(1<<uint(p-1))])
		}
		sv.blockMats = append(sv.blockMats, row)
	}
}

func (sv *solver) answerQueries() []int64 {
	results := make([]int64, sv.n-sv.k+1)
	initial := [10]int64{}
	for d := 1; d <= 9; d++ {
		initial[d] = inf
	}
	for start := 0; start <= sv.n-sv.k; start++ {
		l := start
		r := start + sv.k - 1
		state := initial
		// left partial
		for l <= r && l%blockSize != 0 {
			blockEnd := (l/blockSize+1)*blockSize - 1
			if blockEnd >= sv.n {
				blockEnd = sv.n - 1
			}
			segEnd := r
			if blockEnd < segEnd {
				segEnd = blockEnd
			}
			state = sv.applySegment(state, l, segEnd)
			l = segEnd + 1
		}
		if l <= r {
			blockIdx := l / blockSize
			endBlockIdx := r / blockSize
			for blockIdx < endBlockIdx {
				applied := false
				for p := sv.maxPow; p >= 0; p-- {
					if blockIdx+(1<<uint(p)) <= endBlockIdx {
						state = applyMatrix(state, sv.blockMats[p][blockIdx])
						blockIdx += 1 << uint(p)
						l = blockIdx * blockSize
						applied = true
						break
					}
				}
				if !applied {
					break
				}
			}
			if l <= r {
				state = sv.applySegment(state, l, r)
			}
		}
		results[start] = state[0]
	}
	return results
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, k int
	fmt.Fscan(in, &n, &k)
	var s string
	fmt.Fscan(in, &s)
	digits := make([]int, n)
	for i := 0; i < n; i++ {
		digits[i] = int(s[i] - '0')
	}
	sv := solver{n: n, k: k, digits: digits}
	sv.buildBlocks()
	ans := sv.answerQueries()
	out := bufio.NewWriter(os.Stdout)
	for i, v := range ans {
		if i > 0 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, v)
	}
	fmt.Fprintln(out)
	out.Flush()
}

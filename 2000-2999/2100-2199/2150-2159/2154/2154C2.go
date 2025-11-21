package main

import (
	"bufio"
	"fmt"
	"os"
)

const maxValue = 200000 + 2

type pair struct {
	firstCost  int64
	firstIdx   int
	secondCost int64
	secondIdx  int
}

var spf []int

func buildSPF(limit int) []int {
	spf := make([]int, limit+1)
	for i := 2; i <= limit; i++ {
		if spf[i] == 0 {
			spf[i] = i
			if i*i <= limit {
				for j := i * i; j <= limit; j += i {
					if spf[j] == 0 {
						spf[j] = i
					}
				}
			}
		}
	}
	return spf
}

type fastScanner struct {
	r *bufio.Reader
}

func newFastScanner() *fastScanner {
	return &fastScanner{r: bufio.NewReader(os.Stdin)}
}

func (fs *fastScanner) nextInt() int {
	sign := 1
	val := 0
	c, err := fs.r.ReadByte()
	for (c < '0' || c > '9') && c != '-' {
		if err != nil {
			return 0
		}
		c, err = fs.r.ReadByte()
	}
	if c == '-' {
		sign = -1
		c, err = fs.r.ReadByte()
	}
	for c >= '0' && c <= '9' {
		val = val*10 + int(c-'0')
		c, err = fs.r.ReadByte()
		if err != nil {
			break
		}
	}
	return sign * val
}

func main() {
	spf = buildSPF(maxValue)
	fs := newFastScanner()
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	t := fs.nextInt()
	for ; t > 0; t-- {
		n := fs.nextInt()
		a := make([]int, n)
		b := make([]int, n)
		for i := 0; i < n; i++ {
			a[i] = fs.nextInt()
		}
		for i := 0; i < n; i++ {
			b[i] = fs.nextInt()
		}
		ans := solveCase(a, b)
		fmt.Fprintln(out, ans)
	}
}

func solveCase(a, b []int) int64 {
	const inf int64 = 1 << 60
	n := len(a)

	min1 := int64(1 << 60)
	min2 := int64(1 << 60)
	idxMin := -1
	countMin1 := 0
	for i := 0; i < n; i++ {
		val := int64(b[i])
		if val < min1 {
			min2 = min1
			min1 = val
			idxMin = i
			countMin1 = 1
		} else if val == min1 {
			countMin1++
		} else if val < min2 {
			min2 = val
		}
	}
	if countMin1 >= 2 {
		min2 = min1
	}
	sumBound := min1 + min2

	uniqueMin := countMin1 == 1
	primeMap := make(map[int]*pair)
	primeList := make([]int, 0, 64)

	addCandidate := func(p int, idx int, cost int64) {
		entry, ok := primeMap[p]
		if !ok {
			entry = &pair{
				firstCost:  cost,
				firstIdx:   idx,
				secondCost: inf,
				secondIdx:  -1,
			}
			primeMap[p] = entry
			primeList = append(primeList, p)
			return
		}
		if idx == entry.firstIdx {
			if cost < entry.firstCost {
				entry.firstCost = cost
			}
			return
		}
		if idx == entry.secondIdx {
			if cost < entry.secondCost {
				entry.secondCost = cost
			}
			return
		}
		if cost < entry.firstCost {
			entry.secondCost = entry.firstCost
			entry.secondIdx = entry.firstIdx
			entry.firstCost = cost
			entry.firstIdx = idx
		} else if cost < entry.secondCost {
			entry.secondCost = cost
			entry.secondIdx = idx
		}
	}

	processValue := func(idx int, deltaLimit int) {
		for d := 0; d <= deltaLimit; d++ {
			val := a[idx] + d
			if val <= 1 {
				continue
			}
			tmp := val
			for tmp > 1 {
				p := spf[tmp]
				rem := a[idx] % p
				delta := 0
				if rem != 0 {
					delta = p - rem
				}
				cost := int64(delta) * int64(b[idx])
				addCandidate(p, idx, cost)
				for tmp%p == 0 {
					tmp /= p
				}
			}
		}
	}

	for i := 0; i < n; i++ {
		if uniqueMin && i == idxMin {
			continue
		}
		deltaLimit := int(sumBound / int64(b[i]))
		if deltaLimit > 2 {
			deltaLimit = 2
		}
		if deltaLimit < 0 {
			deltaLimit = 0
		}
		processValue(i, deltaLimit)
	}

	if uniqueMin {
		deltaLimit := int(sumBound / int64(b[idxMin]))
		for _, p := range primeList {
			rem := a[idxMin] % p
			delta := 0
			if rem != 0 {
				delta = p - rem
			}
			if delta > deltaLimit {
				continue
			}
			cost := int64(delta) * int64(b[idxMin])
			addCandidate(p, idxMin, cost)
		}
	}

	answer := inf
	for _, p := range primeList {
		entry := primeMap[p]
		if entry.secondIdx == -1 {
			continue
		}
		total := entry.firstCost + entry.secondCost
		if total < answer {
			answer = total
		}
	}

	if answer == inf {
		return 0
	}
	return answer
}

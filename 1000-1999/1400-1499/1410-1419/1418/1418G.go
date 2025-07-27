package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	a := make([]int, n)
	maxVal := 0
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
		if a[i] > maxVal {
			maxVal = a[i]
		}
	}
	if maxVal < n {
		maxVal = n
	}

	rand.Seed(time.Now().UnixNano())

	h1 := make([]uint64, maxVal+1)
	h2 := make([]uint64, maxVal+1)
	for i := 1; i <= maxVal; i++ {
		h1[i] = rand.Uint64()
		h2[i] = rand.Uint64()
	}

	freq := make([]int, maxVal+1)
	positions := make([][]int, maxVal+1)

	prefixState := make([]uint64, n+1)
	state := uint64(0)
	prefixState[0] = state
	countMap := make(map[uint64]int)
	countMap[state] = 1

	left := 1
	ptr := 0
	var ans int64

	for r := 1; r <= n; r++ {
		x := a[r-1]
		prev := freq[x] % 3
		freq[x]++
		cur := freq[x] % 3
		if prev == 0 && cur == 1 {
			state ^= h1[x]
		} else if prev == 1 && cur == 2 {
			state ^= h1[x]
			state ^= h2[x]
		} else if prev == 2 && cur == 0 {
			state ^= h2[x]
		}

		positions[x] = append(positions[x], r)
		if len(positions[x]) == 4 {
			if positions[x][0]+1 > left {
				left = positions[x][0] + 1
			}
			positions[x] = positions[x][1:]
		}

		for ptr < left-1 {
			st := prefixState[ptr]
			countMap[st]--
			if countMap[st] == 0 {
				delete(countMap, st)
			}
			ptr++
		}

		prefixState[r] = state
		ans += int64(countMap[state])
		countMap[state]++
	}

	fmt.Fprintln(writer, ans)
}

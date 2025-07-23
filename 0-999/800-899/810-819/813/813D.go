package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}

	// map unique values to ids
	valToID := make(map[int]int)
	var uniq []int
	for _, v := range a {
		if _, ok := valToID[v]; !ok {
			valToID[v] = len(uniq)
			uniq = append(uniq, v)
		}
	}
	m := len(uniq)

	bestVal1 := make([]int, m)
	bestMod1 := make([]int, 7)

	afterVal := make([][]int16, n)
	afterMod := make([][]int16, n)
	beforeVal := make([][]int16, n)
	beforeMod := make([][]int16, n)
	for i := 0; i < n; i++ {
		afterVal[i] = make([]int16, m)
		afterMod[i] = make([]int16, 7)
		beforeVal[i] = make([]int16, m)
		beforeMod[i] = make([]int16, 7)
	}

	dp1 := make([]int, n)
	ans := 0

	for j := 0; j < n; j++ {
		val := a[j]
		id := valToID[val]
		r := val % 7

		best := 0
		if idMinus, ok := valToID[val-1]; ok {
			if bestVal1[idMinus] > best {
				best = bestVal1[idMinus]
			}
		}
		if idPlus, ok := valToID[val+1]; ok {
			if bestVal1[idPlus] > best {
				best = bestVal1[idPlus]
			}
		}
		if bestMod1[r] > best {
			best = bestMod1[r]
		}
		dp1[j] = best + 1

		for i := 0; i < j; i++ {
			cand := dp1[i]
			if idMinus, ok := valToID[val-1]; ok {
				if int(afterVal[i][idMinus]) > cand {
					cand = int(afterVal[i][idMinus])
				}
			}
			if idPlus, ok := valToID[val+1]; ok {
				if int(afterVal[i][idPlus]) > cand {
					cand = int(afterVal[i][idPlus])
				}
			}
			if int(afterMod[i][r]) > cand {
				cand = int(afterMod[i][r])
			}
			if idMinus, ok := valToID[val-1]; ok {
				if int(beforeVal[i][idMinus]) > cand {
					cand = int(beforeVal[i][idMinus])
				}
			}
			if idPlus, ok := valToID[val+1]; ok {
				if int(beforeVal[i][idPlus]) > cand {
					cand = int(beforeVal[i][idPlus])
				}
			}
			if int(beforeMod[i][r]) > cand {
				cand = int(beforeMod[i][r])
			}
			dp2 := cand + 1
			if dp2 > int(afterVal[i][id]) {
				afterVal[i][id] = int16(dp2)
			}
			if dp2 > int(afterMod[i][r]) {
				afterMod[i][r] = int16(dp2)
			}
			idi := valToID[a[i]]
			if dp2 > int(beforeVal[j][idi]) {
				beforeVal[j][idi] = int16(dp2)
			}
			if dp2 > int(beforeMod[j][a[i]%7]) {
				beforeMod[j][a[i]%7] = int16(dp2)
			}
			if dp2 > ans {
				ans = dp2
			}
		}

		if dp1[j] > bestVal1[id] {
			bestVal1[id] = dp1[j]
		}
		if dp1[j] > bestMod1[r] {
			bestMod1[r] = dp1[j]
		}
		if dp1[j] > ans {
			ans = dp1[j]
		}
	}

	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	fmt.Fprintln(writer, ans)
}

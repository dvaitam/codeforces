package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	arr := make([]int, n)
	maxV := 0
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &arr[i])
		if arr[i] > maxV {
			maxV = arr[i]
		}
	}

	first1 := make([]int, maxV+1)
	first2 := make([]int, maxV+1)
	last1 := make([]int, maxV+1)
	last2 := make([]int, maxV+1)
	cnt := make([]int, maxV+1)
	for i := 0; i <= maxV; i++ {
		first1[i] = n + 1
		first2[i] = n + 1
	}

	for idx, v := range arr {
		pos := idx + 1
		for d := 1; d*d <= v; d++ {
			if v%d == 0 {
				update := func(div int) {
					cnt[div]++
					if pos < first1[div] {
						first2[div] = first1[div]
						first1[div] = pos
					} else if pos < first2[div] {
						first2[div] = pos
					}
					if pos > last1[div] {
						last2[div] = last1[div]
						last1[div] = pos
					} else if pos > last2[div] {
						last2[div] = pos
					}
				}
				update(d)
				if d*d != v {
					update(v / d)
				}
			}
		}
	}

	total := int64(n) * int64(n+1) / 2
	S := make([]int64, maxV+1)
	for g := 1; g <= maxV; g++ {
		if cnt[g] >= 2 {
			c := int64(first2[g])*(int64(n)-int64(last1[g])+1) + int64(first1[g])*(int64(last1[g])-int64(last2[g]))
			S[g] = total - c
		}
	}

	F := make([]int64, maxV+1)
	var ans int64
	for g := maxV; g >= 1; g-- {
		val := S[g]
		for m := g * 2; m <= maxV; m += g {
			val -= F[m]
		}
		if val < 0 {
			val = 0
		}
		F[g] = val
		ans += val * int64(g)
	}

	fmt.Println(ans)
}

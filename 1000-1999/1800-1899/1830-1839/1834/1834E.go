package main

import (
	"bufio"
	"fmt"
	"os"
)

const limit = 300005

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func lcm(a, b int64) int64 {
	g := gcd(a, b)
	res := a / g * b
	if res > int64(limit) {
		return int64(limit + 1)
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	if _, err := fmt.Fscan(reader, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(reader, &n)
		arr := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &arr[i])
		}
		seen := make([]bool, limit+2)
		prev := []int64{}
		for _, val := range arr {
			curMap := make(map[int64]struct{})
			if val <= limit {
				curMap[val] = struct{}{}
			}
			for _, v := range prev {
				l := lcm(v, val)
				if l <= limit {
					curMap[l] = struct{}{}
				}
			}
			prev = prev[:0]
			for k := range curMap {
				prev = append(prev, k)
				seen[int(k)] = true
			}
		}
		ans := 1
		for ans <= limit && seen[ans] {
			ans++
		}
		fmt.Fprintln(writer, ans)
	}
}

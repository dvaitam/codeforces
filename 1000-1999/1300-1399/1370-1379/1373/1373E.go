package main

import (
	"bufio"
	"fmt"
	"os"
)

type key struct {
	pos  int
	sum  int
	mask int
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n, k int
		fmt.Fscan(reader, &n, &k)
		result, ok := solve(n, k)
		if !ok {
			fmt.Fprintln(writer, -1)
		} else {
			// trim leading zeros
			i := 0
			for i < len(result)-1 && result[i] == '0' {
				i++
			}
			fmt.Fprintln(writer, result[i:])
		}
	}
}

func solve(n, k int) (string, bool) {
	maxLen := 20
	carries := make([]int, k+1)
	memo := make(map[key]*string)
	res, ok := dfs(0, 0, carries, n, k, maxLen, memo)
	if !ok {
		return "", false
	}
	return res, true
}

func encode(carries []int) int {
	res := 0
	for _, v := range carries {
		res = res*3 + v
	}
	return res
}

func dfs(pos, sum int, carries []int, target, k, maxLen int, memo map[key]*string) (string, bool) {
	if pos == maxLen {
		if sum == target && allZero(carries) {
			return "", true
		}
		return "", false
	}
	if sum > target {
		return "", false
	}
	maxRemain := (maxLen - pos) * 9 * (k + 1)
	if sum+maxRemain < target {
		return "", false
	}
	key := key{pos, sum, encode(carries)}
	if v, ok := memo[key]; ok {
		if v == nil {
			return "", false
		}
		return *v, true
	}
	var best string
	found := false
	for d := 0; d < 10; d++ {
		newCarries := make([]int, k+1)
		newSum := sum
		for i := 0; i <= k; i++ {
			val := d + carries[i]
			if pos == 0 {
				val += i
			}
			newSum += val % 10
			newCarries[i] = val / 10
		}
		cand, ok := dfs(pos+1, newSum, newCarries, target, k, maxLen, memo)
		if ok {
			candidate := cand + string('0'+d)
			if !found || less(candidate, best) {
				best = candidate
				found = true
			}
		}
	}
	if !found {
		memo[key] = nil
		return "", false
	}
	memo[key] = &best
	return best, true
}

func allZero(arr []int) bool {
	for _, v := range arr {
		if v != 0 {
			return false
		}
	}
	return true
}

func less(a, b string) bool {
	if len(a) != len(b) {
		return len(a) < len(b)
	}
	return a < b
}

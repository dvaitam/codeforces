package main

import (
	"bufio"
	"fmt"
	"os"
)

const modE = 1000000007

type perm []int

func clone(p perm) perm {
	q := make(perm, len(p))
	copy(q, p)
	return q
}

func permKey(p perm) string {
	b := make([]byte, len(p)*2)
	idx := 0
	for _, v := range p {
		b[idx] = byte(v)
		idx++
		b[idx] = ','
		idx++
	}
	return string(b)
}

func nextStates(p perm) []perm {
	n := len(p)
	used := make([]bool, n)
	resMap := make(map[string]perm)
	var dfs func(int, perm)
	dfs = func(i int, cur perm) {
		for i < n && used[i] {
			i++
		}
		if i == n {
			key := permKey(cur)
			if _, ok := resMap[key]; !ok {
				resMap[key] = clone(cur)
			}
			return
		}
		used[i] = true
		dfs(i+1, cur)
		used[i] = false
		for j := i + 1; j < n; j++ {
			if !used[j] {
				used[i] = true
				used[j] = true
				cur[i], cur[j] = cur[j], cur[i]
				dfs(i+1, cur)
				cur[i], cur[j] = cur[j], cur[i]
				used[i] = false
				used[j] = false
			}
		}
	}
	dfs(0, clone(p))
	res := make([]perm, 0, len(resMap))
	for _, v := range resMap {
		res = append(res, v)
	}
	return res
}

func countWays(start perm) int {
	n := len(start)
	goal := make(perm, n)
	for i := 0; i < n; i++ {
		goal[i] = i + 1
	}
	startKey := permKey(start)
	goalKey := permKey(goal)
	dist := map[string]int{startKey: 0}
	ways := map[string]int{startKey: 1}
	queue := []perm{start}
	best := -1
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		d := dist[permKey(cur)]
		if best != -1 && d >= best {
			continue
		}
		if permKey(cur) == goalKey {
			best = d
			continue
		}
		for _, nxt := range nextStates(cur) {
			k := permKey(nxt)
			if prev, ok := dist[k]; !ok {
				dist[k] = d + 1
				ways[k] = ways[permKey(cur)]
				queue = append(queue, nxt)
			} else if prev == d+1 {
				ways[k] = (ways[k] + ways[permKey(cur)]) % modE
			}
		}
	}
	if best == -1 {
		if startKey == goalKey {
			return 1
		}
		return 0
	}
	return ways[goalKey]
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, k int
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &arr[i])
	}
	used := make([]bool, n+1)
	for _, v := range arr {
		if v != 0 {
			used[v] = true
		}
	}
	missing := make([]int, 0)
	for i := 1; i <= n; i++ {
		if !used[i] {
			missing = append(missing, i)
		}
	}
	pos := make([]int, 0)
	for i, v := range arr {
		if v == 0 {
			pos = append(pos, i)
		}
	}
	total := 0
	var dfs func(int)
	dfs = func(idx int) {
		if idx == len(pos) {
			p := make(perm, n)
			copy(p, arr)
			val := countWays(p)
			total += val
			if total >= modE {
				total %= modE
			}
			return
		}
		for i, num := range missing {
			if num == -1 {
				continue
			}
			arr[pos[idx]] = num
			missing[i] = -1
			dfs(idx + 1)
			missing[i] = num
			arr[pos[idx]] = 0
		}
	}
	dfs(0)
	if len(pos) == 0 {
		total = countWays(perm(arr)) % modE
	}
	fmt.Println(total % modE)
}

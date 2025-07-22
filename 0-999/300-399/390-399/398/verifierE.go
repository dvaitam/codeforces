package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const mod = 1000000007

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
				ways[k] = (ways[k] + ways[permKey(cur)]) % mod
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

func solveE(n, k int, arr []int) int {
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
	positions := make([]int, 0)
	for i, v := range arr {
		if v == 0 {
			positions = append(positions, i)
		}
	}
	total := 0
	var dfs func(int)
	dfs = func(idx int) {
		if idx == len(positions) {
			p := make(perm, n)
			copy(p, arr)
			val := countWays(p)
			total += val
			if total >= mod {
				total %= mod
			}
			return
		}
		for i, num := range missing {
			if num == -1 {
				continue
			}
			arr[positions[idx]] = num
			missing[i] = -1
			dfs(idx + 1)
			missing[i] = num
			arr[positions[idx]] = 0
		}
	}
	dfs(0)
	if len(positions) == 0 {
		return countWays(perm(arr)) % mod
	}
	return total % mod
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesE.txt")
	if err != nil {
		fmt.Println("could not open testcasesE.txt:", err)
		os.Exit(1)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		fields := strings.Fields(line)
		if len(fields) < 2 {
			fmt.Printf("bad test case %d\n", idx)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(fields[0])
		k, _ := strconv.Atoi(fields[1])
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			arr[i], _ = strconv.Atoi(fields[2+i])
		}
		exp := solveE(n, k, arr)
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d %d\n", n, k))
		for i := 0; i < n; i++ {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(strconv.Itoa(arr[i]))
		}
		input.WriteByte('\n')
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input.String())
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Test %d runtime error: %v\n", idx, err)
			os.Exit(1)
		}
		gotStr := strings.TrimSpace(string(out))
		got, err := strconv.Atoi(gotStr)
		if err != nil {
			fmt.Printf("Test %d: output not integer\n", idx)
			os.Exit(1)
		}
		if got%mod != exp%mod {
			fmt.Printf("Test %d failed: expected %d got %d\n", idx, exp%mod, got%mod)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}

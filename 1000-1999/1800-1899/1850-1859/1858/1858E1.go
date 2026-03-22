package main

import (
	"bufio"
	"os"
	"strconv"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)

	buf := make([]byte, 0, 1024*1024)
	scanner.Buffer(buf, 10*1024*1024)

	if !scanner.Scan() {
		return
	}
	q, _ := strconv.Atoi(scanner.Text())

	up := make([][20]int, q+1)
	val := make([]int, q+1)

	stack := make([]int, 0, q+1)
	stack = append(stack, 0)

	queries := make([]int, 0, q)

	id := 0

	for i := 0; i < q; i++ {
		scanner.Scan()
		op := scanner.Text()
		if op == "+" {
			scanner.Scan()
			x, _ := strconv.Atoi(scanner.Text())
			id++
			val[id] = x
			p := stack[len(stack)-1]
			up[id][0] = p
			for j := 1; j < 20; j++ {
				up[id][j] = up[up[id][j-1]][j-1]
			}
			stack = append(stack, id)
		} else if op == "-" {
			scanner.Scan()
			k, _ := strconv.Atoi(scanner.Text())
			curr := stack[len(stack)-1]
			for j := 19; j >= 0; j-- {
				if (k & (1 << j)) != 0 {
					curr = up[curr][j]
				}
			}
			stack = append(stack, curr)
		} else if op == "!" {
			stack = stack[:len(stack)-1]
		} else if op == "?" {
			queries = append(queries, stack[len(stack)-1])
		}
	}

	head := make([]int, id+1)
	for i := range head {
		head[i] = -1
	}
	next := make([]int, id+1)

	for i := 1; i <= id; i++ {
		p := up[i][0]
		next[i] = head[p]
		head[p] = i
	}

	ansAt := make([]int, id+1)
	freq := make([]int, 1000005)
	distinct := 0

	var dfs func(int)
	dfs = func(u int) {
		ansAt[u] = distinct
		for v := head[u]; v != -1; v = next[v] {
			x := val[v]
			if freq[x] == 0 {
				distinct++
			}
			freq[x]++

			dfs(v)

			freq[x]--
			if freq[x] == 0 {
				distinct--
			}
		}
	}

	dfs(0)

	out := bufio.NewWriter(os.Stdout)
	for _, qNode := range queries {
		out.WriteString(strconv.Itoa(ansAt[qNode]))
		out.WriteByte('\n')
	}
	out.Flush()
}

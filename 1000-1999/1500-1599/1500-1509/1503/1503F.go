package main

import (
	"bufio"
	"fmt"
	"os"
)

type pair struct {
	first  int
	second bool
}

var (
	n, m       int
	a, b       []int
	ap, bp     []int
	vis, del   []bool
	L, R       []pair
	lson, rson [][]int
	ans        []int
	reader     = bufio.NewReader(os.Stdin)
	writer     = bufio.NewWriter(os.Stdout)
)

func readInt() int {
	var x int
	var sign = 1
	b, err := reader.ReadByte()
	if err != nil {
		return 0
	}
	for (b < '0' || b > '9') && b != '-' {
		b, _ = reader.ReadByte()
	}
	if b == '-' {
		sign = -1
		b, _ = reader.ReadByte()
	}
	for b >= '0' && b <= '9' {
		x = x*10 + int(b-'0')
		b, _ = reader.ReadByte()
	}
	return x * sign
}

func dfs(id int) {
	for _, v := range lson[id] {
		dfs(v)
	}
	ans = append(ans, id)
	for _, v := range rson[id] {
		dfs(v)
	}
}

func calc(start int) {
	queue := make([]int, 0)
	cur := start
	for {
		next := R[cur].first
		if R[cur].second == R[next].second {
			queue = append(queue, cur)
		}
		cur = next
		if cur == start {
			break
		}
	}
	for L[start].first != R[start].first || L[start].second != R[start].second {
		if R[start].first == L[start].first || len(queue) == 0 {
			fmt.Fprintln(writer, "NO")
			writer.Flush()
			os.Exit(0)
		}
		i := queue[0]
		queue = queue[1:]
		if del[i] {
			continue
		}
		ni := R[i].first
		if R[i].second != R[ni].second {
			continue
		}
		del[ni] = true
		ni2 := R[ni].first
		del[ni2] = true
		if R[i].second {
			rson[i] = append(rson[i], ni)
			rson[i] = append(rson[i], ni2)
		} else {
			lson[i] = append([]int{ni}, lson[i]...)
			lson[i] = append([]int{ni2}, lson[i]...)
		}
		j := R[ni2].first
		L[j].first = i
		R[i].first = j
		R[i].second = !L[j].second
		start = i
		if L[i].second != R[i].second {
			queue = append(queue, L[i].first)
		}
		if R[i].second == R[R[i].first].second {
			queue = append(queue, i)
		}
	}
	if R[start].second {
		dfs(start)
		dfs(R[start].first)
	} else {
		dfs(R[start].first)
		dfs(start)
	}
}

func main() {
	defer writer.Flush()
	n = readInt()
	m = 2 * n
	a = make([]int, m+1)
	b = make([]int, m+1)
	ap = make([]int, 2*n+1)
	bp = make([]int, 2*n+1)
	vis = make([]bool, m+1)
	del = make([]bool, m+1)
	L = make([]pair, m+1)
	R = make([]pair, m+1)
	lson = make([][]int, m+1)
	rson = make([][]int, m+1)
	for i := 1; i <= m; i++ {
		ai := readInt()
		bi := readInt()
		a[i] = ai
		b[i] = bi
		ap[n+ai] = i
		bp[n+bi] = i
	}
	for i := 1; i <= m; i++ {
		if !vis[i] {
			j := i
			for !vis[j] {
				vis[j] = true
				nxa := ap[n-a[j]]
				nxb := bp[n-b[j]]
				if !vis[nxa] {
					R[j] = pair{nxa, a[j] > 0}
					L[nxa] = pair{j, a[j] < 0}
					j = nxa
				} else {
					R[j] = pair{nxb, b[j] > 0}
					L[nxb] = pair{j, b[j] < 0}
					j = nxb
				}
			}
			calc(i)
		}
	}
	fmt.Fprintln(writer, "Yes")
	for _, id := range ans {
		fmt.Fprintf(writer, "%d %d\n", a[id], b[id])
	}
}

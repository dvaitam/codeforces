package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type girl struct {
	node    int
	weight  int64
	invited bool
}

func bfsPath(adj [][]int, start, goal int) []int {
	n := len(adj) - 1
	queue := []int{start}
	parent := make([]int, n+1)
	for i := range parent {
		parent[i] = -1
	}
	parent[start] = start
	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]
		if v == goal {
			break
		}
		for _, u := range adj[v] {
			if parent[u] == -1 {
				parent[u] = v
				queue = append(queue, u)
			}
		}
	}
	if parent[goal] == -1 {
		return nil
	}
	path := []int{}
	v := goal
	for v != start {
		path = append(path, v)
		v = parent[v]
	}
	path = append(path, start)
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return path
}

func subtreeNodes(adj [][]int, start int) []int {
	res := []int{}
	stack := []int{start}
	parent := make(map[int]int)
	parent[start] = -1
	for len(stack) > 0 {
		v := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		res = append(res, v)
		for _, u := range adj[v] {
			if u == parent[v] {
				continue
			}
			parent[u] = v
			stack = append(stack, u)
		}
	}
	return res
}

func solveE(n, m, q int, edges [][2]int, girlsAt []int, events [][]int) string {
	adj := make([][]int, n+1)
	for _, e := range edges {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	girls := make([]girl, m)
	for i := 0; i < m; i++ {
		girls[i] = girl{node: girlsAt[i], weight: int64(i + 1)}
	}
	var sb strings.Builder
	for _, ev := range events {
		if ev[0] == 1 {
			v, u, k := ev[1], ev[2], ev[3]
			nodes := bfsPath(adj, v, u)
			cand := []int{}
			for i := 0; i < m; i++ {
				if girls[i].invited {
					continue
				}
				for _, nd := range nodes {
					if girls[i].node == nd {
						cand = append(cand, i)
						break
					}
				}
			}
			for i := 0; i < len(cand); i++ {
				for j := i + 1; j < len(cand); j++ {
					gi, gj := girls[cand[i]], girls[cand[j]]
					if gi.weight > gj.weight || (gi.weight == gj.weight && gi.node > gj.node) {
						cand[i], cand[j] = cand[j], cand[i]
					}
				}
			}
			if k > len(cand) {
				k = len(cand)
			}
			sb.WriteString(strconv.Itoa(k))
			for i := 0; i < k; i++ {
				id := cand[i]
				sb.WriteString(" ")
				sb.WriteString(strconv.Itoa(id + 1))
				girls[id].invited = true
			}
			sb.WriteByte('\n')
		} else {
			v, add := ev[1], ev[2]
			nodes := subtreeNodes(adj, v)
			for i := 0; i < m; i++ {
				if girls[i].invited {
					continue
				}
				for _, nd := range nodes {
					if girls[i].node == nd {
						girls[i].weight += int64(add)
						break
					}
				}
			}
		}
	}
	return sb.String()
}

func runCase(exe, input, expected string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(expected)
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	data, err := os.ReadFile("testcasesE.txt")
	if err != nil {
		fmt.Println("could not read testcasesE.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		if !scan.Scan() {
			fmt.Println("bad file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		m, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		q, _ := strconv.Atoi(scan.Text())
		edges := make([][2]int, n-1)
		for i := 0; i < n-1; i++ {
			scan.Scan()
			u, _ := strconv.Atoi(scan.Text())
			scan.Scan()
			v, _ := strconv.Atoi(scan.Text())
			edges[i] = [2]int{u, v}
		}
		girlsAt := make([]int, m)
		for i := 0; i < m; i++ {
			scan.Scan()
			girlsAt[i], _ = strconv.Atoi(scan.Text())
		}
		events := make([][]int, q)
		for i := 0; i < q; i++ {
			scan.Scan()
			typ, _ := strconv.Atoi(scan.Text())
			if typ == 1 {
				scan.Scan()
				v, _ := strconv.Atoi(scan.Text())
				scan.Scan()
				u, _ := strconv.Atoi(scan.Text())
				scan.Scan()
				k, _ := strconv.Atoi(scan.Text())
				events[i] = []int{1, v, u, k}
			} else {
				scan.Scan()
				v, _ := strconv.Atoi(scan.Text())
				scan.Scan()
				add, _ := strconv.Atoi(scan.Text())
				events[i] = []int{2, v, add}
			}
		}
		var inSB strings.Builder
		inSB.WriteString(fmt.Sprintf("%d %d %d\n", n, m, q))
		for _, e := range edges {
			inSB.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
		}
		for i, c := range girlsAt {
			if i > 0 {
				inSB.WriteByte(' ')
			}
			inSB.WriteString(strconv.Itoa(c))
		}
		inSB.WriteByte('\n')
		for _, ev := range events {
			if ev[0] == 1 {
				inSB.WriteString(fmt.Sprintf("1 %d %d %d\n", ev[1], ev[2], ev[3]))
			} else {
				inSB.WriteString(fmt.Sprintf("2 %d %d\n", ev[1], ev[2]))
			}
		}
		input := inSB.String()
		expected := solveE(n, m, q, edges, girlsAt, events)
		if err := runCase(exe, input, expected); err != nil {
			fmt.Printf("case %d failed: %v\ninput:\n%s", caseIdx+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

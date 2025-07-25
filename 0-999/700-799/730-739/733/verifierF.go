package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

const Log = 19

type edge struct {
	u, v, w, c, id int
}

func expectedCase(input string) string {
	reader := bufio.NewReader(strings.NewReader(input))
	var n, m, s int
	fmt.Fscan(reader, &n, &m)
	edges := make([]edge, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(reader, &edges[i].w)
	}
	for i := 0; i < m; i++ {
		fmt.Fscan(reader, &edges[i].c)
	}
	for i := 0; i < m; i++ {
		fmt.Fscan(reader, &edges[i].u, &edges[i].v)
		edges[i].u--
		edges[i].v--
		edges[i].id = i
	}
	fmt.Fscan(reader, &s)
	sort.Slice(edges, func(i, j int) bool { return edges[i].w < edges[j].w })
	to := make([]int, m)
	for i := 0; i < m; i++ {
		to[edges[i].id] = i
	}
	parent := make([]int, n)
	for i := 0; i < n; i++ {
		parent[i] = i
	}
	var find func(int) int
	find = func(x int) int {
		for parent[x] != x {
			parent[x] = parent[parent[x]]
			x = parent[x]
		}
		return x
	}
	totw := int64(0)
	minc := int(1e9)
	pc := -1
	vis := make([]bool, m)
	adj := make([][]struct{ to, w int }, n)
	for i := 0; i < m; i++ {
		u := edges[i].u
		v := edges[i].v
		fu := find(u)
		fv := find(v)
		if fu != fv {
			parent[fu] = fv
			totw += int64(edges[i].w)
			if edges[i].c < minc {
				minc = edges[i].c
				pc = i
			}
			vis[i] = true
			adj[u] = append(adj[u], struct{ to, w int }{v, edges[i].w})
			adj[v] = append(adj[v], struct{ to, w int }{u, edges[i].w})
		}
	}
	depth := make([]int, n)
	parentUp := make([][]int, Log)
	maxUp := make([][]int, Log)
	for j := 0; j < Log; j++ {
		parentUp[j] = make([]int, n)
		maxUp[j] = make([]int, n)
	}
	queue := []int{0}
	parentUp[0][0] = -1
	for qi := 0; qi < len(queue); qi++ {
		u := queue[qi]
		for _, e := range adj[u] {
			v := e.to
			if v == parentUp[0][u] {
				continue
			}
			depth[v] = depth[u] + 1
			parentUp[0][v] = u
			maxUp[0][v] = e.w
			queue = append(queue, v)
		}
	}
	for j := 1; j < Log; j++ {
		for i := 0; i < n; i++ {
			p := parentUp[j-1][i]
			if p < 0 {
				parentUp[j][i] = -1
				maxUp[j][i] = maxUp[j-1][i]
			} else {
				parentUp[j][i] = parentUp[j-1][p]
				if parentUp[j][i] >= 0 && maxUp[j-1][p] > maxUp[j-1][i] {
					maxUp[j][i] = maxUp[j-1][p]
				} else {
					maxUp[j][i] = maxUp[j-1][i]
				}
			}
		}
	}
	getMax := func(u, v int) int {
		res := 0
		if depth[u] < depth[v] {
			u, v = v, u
		}
		dd := depth[u] - depth[v]
		for j := 0; j < Log; j++ {
			if dd&(1<<j) != 0 {
				if maxUp[j][u] > res {
					res = maxUp[j][u]
				}
				u = parentUp[j][u]
			}
		}
		if u == v {
			return res
		}
		for j := Log - 1; j >= 0; j-- {
			if parentUp[j][u] != parentUp[j][v] {
				if maxUp[j][u] > res {
					res = maxUp[j][u]
				}
				if maxUp[j][v] > res {
					res = maxUp[j][v]
				}
				u = parentUp[j][u]
				v = parentUp[j][v]
			}
		}
		if maxUp[0][u] > res {
			res = maxUp[0][u]
		}
		if maxUp[0][v] > res {
			res = maxUp[0][v]
		}
		return res
	}
	ans := totw - int64(s/minc)
	pos := -1
	for i := 0; i < m; i++ {
		if !vis[i] && edges[i].c < minc {
			wmx := getMax(edges[i].u, edges[i].v)
			alt := totw - int64(wmx) + int64(edges[i].w) - int64(s/edges[i].c)
			if alt < ans {
				ans = alt
				pos = i
			}
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", ans))
	if pos < 0 {
		for orig := 0; orig < m; orig++ {
			idx := to[orig]
			if vis[idx] {
				if idx == pc {
					sb.WriteString(fmt.Sprintf("%d %d\n", orig+1, edges[idx].w-s/minc))
				} else {
					sb.WriteString(fmt.Sprintf("%d %d\n", orig+1, edges[idx].w))
				}
			}
		}
	} else {
		u := edges[pos].u
		v := edges[pos].v
		U, V := 0, 0
		W := 0
		uu := u
		vv := v
		for uu != vv {
			if depth[uu] < depth[vv] {
				uu, vv = vv, uu
			}
			if maxUp[0][uu] > W {
				W = maxUp[0][uu]
				U = uu
				V = parentUp[0][uu]
			}
			uu = parentUp[0][uu]
		}
		if U < V {
			U, V = V, U
		}
		for orig := 0; orig < m; orig++ {
			idx := to[orig]
			if vis[idx] && !(max(edges[idx].u, edges[idx].v) == U && min(edges[idx].u, edges[idx].v) == V) {
				sb.WriteString(fmt.Sprintf("%d %d\n", orig+1, edges[idx].w))
			} else if idx == pos {
				sb.WriteString(fmt.Sprintf("%d %d\n", orig+1, edges[idx].w-s/edges[idx].c))
			}
		}
	}
	return sb.String()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
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
		fmt.Println("Usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	data, err := os.ReadFile("testcasesF.txt")
	if err != nil {
		fmt.Println("could not read testcasesF.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		m, _ := strconv.Atoi(scan.Text())
		w := make([]int, m)
		c := make([]int, m)
		for j := 0; j < m; j++ {
			scan.Scan()
			w[j], _ = strconv.Atoi(scan.Text())
		}
		for j := 0; j < m; j++ {
			scan.Scan()
			c[j], _ = strconv.Atoi(scan.Text())
		}
		edgesStr := make([][2]int, m)
		for j := 0; j < m; j++ {
			scan.Scan()
			edgesStr[j][0], _ = strconv.Atoi(scan.Text())
			scan.Scan()
			edgesStr[j][1], _ = strconv.Atoi(scan.Text())
		}
		scan.Scan()
		s, _ := strconv.Atoi(scan.Text())
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for j := 0; j < m; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(w[j]))
		}
		sb.WriteString("\n")
		for j := 0; j < m; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(c[j]))
		}
		sb.WriteString("\n")
		for j := 0; j < m; j++ {
			sb.WriteString(fmt.Sprintf("%d %d\n", edgesStr[j][0], edgesStr[j][1]))
		}
		sb.WriteString(fmt.Sprintf("%d\n", s))
		exp := expectedCase(sb.String())
		if err := runCase(exe, sb.String(), exp); err != nil {
			fmt.Printf("case %d failed: %v\ninput:%s", i+1, err, sb.String())
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

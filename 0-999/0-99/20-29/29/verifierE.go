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

const embeddedSolutionSource = `package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

type pair struct{ f, s int }

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   v := make([][]int, n)
   for i := 0; i < m; i++ {
       var f, s int
       fmt.Fscan(reader, &f, &s)
       f--
       s--
       v[f] = append(v[f], s)
       v[s] = append(v[s], f)
   }
   mod := 1 << 15
   d1 := make([]int, n)
   u1 := make([]bool, n)
   for i := range d1 {
       d1[i] = mod
   }
   d1[0], u1[0] = 0, true
   queue1 := []int{0}
   for len(queue1) > 0 {
       p := queue1[0]
       queue1 = queue1[1:]
       for _, nei := range v[p] {
           if !u1[nei] {
               u1[nei] = true
               d1[nei] = d1[p] + 1
               queue1 = append(queue1, nei)
           }
       }
   }
   d2 := make([]int, n)
   u2 := make([]bool, n)
   for i := range d2 {
       d2[i] = mod
   }
   d2[n-1], u2[n-1] = 0, true
   queue2 := []int{n - 1}
   for len(queue2) > 0 {
       p := queue2[0]
       queue2 = queue2[1:]
       for _, nei := range v[p] {
           if !u2[nei] {
               u2[nei] = true
               d2[nei] = d2[p] + 1
               queue2 = append(queue2, nei)
           }
       }
   }
   d := make([][]int, n)
   used := make([][]bool, n)
   pf := make([][]int, n)
   ps := make([][]int, n)
   for i := 0; i < n; i++ {
       d[i] = make([]int, n)
       used[i] = make([]bool, n)
       pf[i] = make([]int, n)
       ps[i] = make([]int, n)
       for j := 0; j < n; j++ {
           d[i][j] = mod
       }
   }
   queue := make([]pair, 0, n*n)
   queue = append(queue, pair{0, n - 1})
   used[0][n-1] = true
   d[0][n-1] = 0
   head := 0
   found := false
   for head < len(queue) && !found {
       t := queue[head]
       head++
       t1, t2 := t.f, t.s
       neigh1 := append([]int(nil), v[t1]...)
       neigh2 := append([]int(nil), v[t2]...)
       sort.Slice(neigh1, func(i, j int) bool { return d2[neigh1[i]] < d2[neigh1[j]] })
       sort.Slice(neigh2, func(i, j int) bool { return d1[neigh2[i]] < d1[neigh2[j]] })
       cnt1, cnt2 := 0, 0
       i := 0
       for cnt1 < 4 && i < len(neigh1) && !found {
           j := 0
           for cnt2 < 4 && j < len(neigh2) {
               from, to := neigh1[i], neigh2[j]
               if !used[from][to] && from != to {
                   used[from][to] = true
                   pf[from][to] = t1
                   ps[from][to] = t2
                   d[from][to] = d[t1][t2] + 1
                   if from == n-1 && to == 0 {
                       found = true
                       break
                   }
                   queue = append(queue, pair{from, to})
               }
               j++
               cnt2++
           }
           i++
           cnt1++
       }
   }
   if d[n-1][0] == mod {
       fmt.Println(-1)
       return
   }
   fmt.Println(d[n-1][0])
   ansf := []int{}
   anss := []int{}
   curf, curs := n-1, 0
   for curf != 0 || curs != 0 {
       ansf = append(ansf, curf+1)
       anss = append(anss, curs+1)
       prevf, prevs := pf[curf][curs], ps[curf][curs]
       curf, curs = prevf, prevs
   }
   for i := len(ansf) - 1; i >= 0; i-- {
       fmt.Printf("%d ", ansf[i])
   }
   fmt.Println()
   for i := len(anss) - 1; i >= 0; i-- {
       fmt.Printf("%d ", anss[i])
   }
   fmt.Println()
}`

const testcasesRaw = `
2 1 1 2
3 3 1 2 1 3 2 3
4 4 2 4 1 2 1 3 3 4
6 7 1 6 2 3 3 5 1 4 2 4 4 5 3 4
4 6 1 3 2 3 1 2 2 4 3 4 1 4
3 3 1 2 1 3 2 3
5 7 1 2 1 3 2 5 4 5 2 4 3 5 1 5
6 9 5 6 1 2 4 6 2 3 1 4 1 3 2 5 4 5 1 5
6 7 1 3 3 4 3 6 1 2 2 4 1 4 1 5
5 8 4 5 2 5 3 4 2 4 1 4 3 5 2 3 1 5
3 3 1 2 1 3 2 3
2 1 1 2
2 1 1 2
4 5 1 2 3 4 1 4 2 4 2 3
2 1 1 2
2 1 1 2
3 3 2 3 1 3 1 2
3 3 2 3 1 2 1 3
3 2 2 3 1 3
6 8 4 5 2 3 1 4 1 3 2 5 1 6 3 6 4 6
4 5 2 3 1 2 3 4 1 3 2 4
3 3 1 2 2 3 1 3
5 8 1 4 3 4 2 5 1 3 1 5 1 2 3 5 2 3
5 8 3 5 1 4 1 2 1 5 2 3 1 3 2 4 2 5
2 1 1 2
5 4 3 5 4 5 1 5 2 4
2 1 1 2
3 3 2 3 1 3 1 2
5 7 2 3 3 4 2 5 1 3 1 5 3 5 4 5
2 1 1 2
6 6 3 6 2 5 5 6 1 3 3 5 1 5
6 5 2 4 2 3 5 6 3 5 1 2
4 3 1 4 2 3 1 2
4 4 2 4 3 4 2 3 1 2
4 3 2 4 1 2 1 3
4 6 1 2 1 3 1 4 2 4 3 4 2 3
6 5 2 5 4 6 1 6 3 6 3 4
6 6 1 2 1 3 1 4 1 5 4 5 2 6
6 9 4 5 1 4 1 3 3 5 2 5 3 4 1 6 1 5 5 6
3 2 1 3 1 2
2 1 1 2
2 1 1 2
3 3 1 2 2 3 1 3
5 4 1 3 3 5 2 3 1 4
5 7 1 3 3 4 1 2 1 5 2 3 4 5 2 4
6 7 3 6 1 2 4 5 2 4 1 3 1 6 1 5
5 8 3 4 1 2 2 4 2 5 4 5 1 3 1 4 3 5
4 5 3 4 2 3 1 2 1 4 1 3
6 5 1 6 1 2 1 5 5 6 2 4
6 6 1 3 3 5 5 6 1 2 1 4 3 4
6 5 1 2 3 5 3 6 5 6 2 4
2 1 1 2
4 6 1 2 2 4 3 4 2 3 1 4 1 3
4 3 3 4 2 3 1 4
6 6 5 6 3 5 1 6 4 6 1 4 2 4
6 8 2 5 3 6 1 2 2 4 2 3 1 3 4 6 4 5
4 4 3 4 2 3 1 2 1 4
5 4 3 4 2 4 2 5 1 3
6 8 1 3 2 6 2 4 2 5 5 6 4 5 2 3 3 5
6 5 1 6 5 6 2 4 3 4 1 3
4 4 1 2 1 3 3 4 2 4
4 3 1 4 1 3 3 4
6 7 1 5 3 4 1 2 4 5 1 6 2 3 3 6
5 4 2 3 3 4 3 5 1 3
4 4 1 4 1 3 1 2 3 4
5 5 2 4 1 4 3 4 3 5 1 5
3 3 2 3 1 3 1 2
2 1 1 2
6 7 1 2 2 3 4 5 3 5 2 4 3 4 1 6
5 8 2 5 2 3 1 5 3 4 3 5 1 3 1 4 2 4
3 2 1 3 2 3
5 4 1 2 4 5 3 5 1 5
6 7 1 4 1 6 1 2 3 6 4 5 2 6 2 4
3 2 1 2 1 3
2 1 1 2
4 5 1 3 2 3 2 4 1 4 1 2
6 9 3 4 2 3 1 6 1 3 2 5 2 6 1 5 1 4 1 2
5 5 4 5 2 5 1 2 1 3 3 4
5 6 1 3 2 4 2 3 1 5 3 4 1 2
2 1 1 2
5 8 4 5 1 4 1 3 2 5 2 3 2 4 3 4 1 2
4 3 1 3 1 4 1 2
6 9 3 5 4 6 5 6 2 4 1 6 2 3 1 3 3 4 1 2
6 5 5 6 3 5 2 3 2 4 1 2
4 5 2 4 1 2 1 3 2 3 3 4
3 2 1 3 2 3
3 3 1 3 1 2 2 3
4 5 1 3 3 4 2 4 2 3 1 4
6 6 3 6 4 5 1 2 4 6 1 4 2 3
5 5 2 3 3 5 1 3 1 5 1 4
2 1 1 2
2 1 1 2
4 4 1 2 1 4 3 4 2 4
6 7 3 4 4 5 1 3 5 6 3 5 2 4 2 6
6 8 3 5 3 4 1 2 1 6 2 6 2 5 5 6 4 5
5 8 3 5 1 4 4 5 3 4 1 5 2 3 1 2 1 3
5 8 1 4 1 3 2 3 3 5 1 2 4 5 2 4 2 5
3 2 1 2 2 3
6 9 2 5 1 6 3 4 2 6 1 5 4 6 2 4 4 5 1 3
5 6 1 3 3 4 2 4 4 5 1 2 2 5
`

var (
	_            = embeddedSolutionSource
	rawTestcases = func() []string {
		scanner := bufio.NewScanner(strings.NewReader(testcasesRaw))
		scanner.Buffer(make([]byte, 0, 1024), 1024*1024)
		var cases []string
		for scanner.Scan() {
			if line := strings.TrimSpace(scanner.Text()); line != "" {
				cases = append(cases, line)
			}
		}
		if err := scanner.Err(); err != nil {
			panic(err)
		}
		return cases
	}()
)

func solveCase(n int, m int, edges [][2]int) string {
	adj := make([][]int, n)
	for _, e := range edges {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	const mod = 1 << 15
	d1 := make([]int, n)
	for i := range d1 {
		d1[i] = mod
	}
	d1[0] = 0
	q1 := []int{0}
	for len(q1) > 0 {
		u := q1[0]
		q1 = q1[1:]
		for _, v := range adj[u] {
			if d1[v] == mod {
				d1[v] = d1[u] + 1
				q1 = append(q1, v)
			}
		}
	}
	d2 := make([]int, n)
	for i := range d2 {
		d2[i] = mod
	}
	d2[n-1] = 0
	q2 := []int{n - 1}
	for len(q2) > 0 {
		u := q2[0]
		q2 = q2[1:]
		for _, v := range adj[u] {
			if d2[v] == mod {
				d2[v] = d2[u] + 1
				q2 = append(q2, v)
			}
		}
	}

	pf := make([][]int, n)
	ps := make([][]int, n)
	d := make([][]int, n)
	used := make([][]bool, n)
	for i := 0; i < n; i++ {
		pf[i] = make([]int, n)
		ps[i] = make([]int, n)
		d[i] = make([]int, n)
		used[i] = make([]bool, n)
		for j := 0; j < n; j++ {
			d[i][j] = mod
		}
	}

	type pair struct{ a, b int }
	queue := []pair{{0, n - 1}}
	d[0][n-1] = 0
	used[0][n-1] = true
	head := 0
	found := false
	for head < len(queue) && !found {
		cur := queue[head]
		head++
		u, v := cur.a, cur.b
		neighU := append([]int(nil), adj[u]...)
		neighV := append([]int(nil), adj[v]...)
		sort.Slice(neighU, func(i, j int) bool { return d2[neighU[i]] < d2[neighU[j]] })
		sort.Slice(neighV, func(i, j int) bool { return d1[neighV[i]] < d1[neighV[j]] })
		cntU := 0
		for i := 0; i < len(neighU) && cntU < 4 && !found; i++ {
			cntV := 0
			for j := 0; j < len(neighV) && cntV < 4; j++ {
				a, b := neighU[i], neighV[j]
				if a == b || used[a][b] {
					cntV++
					continue
				}
				used[a][b] = true
				pf[a][b] = u
				ps[a][b] = v
				d[a][b] = d[u][v] + 1
				if a == n-1 && b == 0 {
					found = true
					break
				}
				queue = append(queue, pair{a, b})
				cntV++
			}
			cntU++
		}
	}

	if d[n-1][0] == mod {
		return "-1"
	}
	var pathA, pathB []int
	curA, curB := n-1, 0
	for curA != 0 || curB != 0 {
		pathA = append(pathA, curA+1)
		pathB = append(pathB, curB+1)
		pa, pb := pf[curA][curB], ps[curA][curB]
		curA, curB = pa, pb
	}
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(d[n-1][0]))
	sb.WriteByte('\n')
	for i := len(pathA) - 1; i >= 0; i-- {
		sb.WriteString(strconv.Itoa(pathA[i]))
		sb.WriteByte(' ')
	}
	sb.WriteByte('\n')
	for i := len(pathB) - 1; i >= 0; i-- {
		sb.WriteString(strconv.Itoa(pathB[i]))
		sb.WriteByte(' ')
	}
	sb.WriteByte('\n')
	return strings.TrimSpace(sb.String())
}

func parseCase(line string) (int, int, [][2]int, error) {
	fields := strings.Fields(strings.TrimSpace(line))
	if len(fields) < 2 {
		return 0, 0, nil, fmt.Errorf("invalid line")
	}
	n, err := strconv.Atoi(fields[0])
	if err != nil {
		return 0, 0, nil, err
	}
	m, err := strconv.Atoi(fields[1])
	if err != nil {
		return 0, 0, nil, err
	}
	if len(fields) != 2+2*m {
		return 0, 0, nil, fmt.Errorf("expected %d numbers got %d", 2+2*m, len(fields))
	}
	edges := make([][2]int, m)
	for i := 0; i < m; i++ {
		a, _ := strconv.Atoi(fields[2+2*i])
		b, _ := strconv.Atoi(fields[2+2*i+1])
		edges[i] = [2]int{a - 1, b - 1}
	}
	return n, m, edges, nil
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v, output: %s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	for idx, line := range rawTestcases {
		n, m, edges, err := parseCase(line)
		if err != nil {
			fmt.Printf("case %d invalid: %v\n", idx+1, err)
			os.Exit(1)
		}
		expected := solveCase(n, m, edges)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for _, e := range edges {
			sb.WriteString(fmt.Sprintf("%d %d\n", e[0]+1, e[1]+1))
		}
		got, err := run(bin, sb.String())
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Printf("case %d failed\nexpected:\n%s\n\ngot:\n%s\n", idx+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(rawTestcases))
}

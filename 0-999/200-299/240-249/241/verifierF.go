package main

import (
	"bytes"
	"container/heap"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesF.txt so the verifier is self-contained.
const testcasesF = `4 3 7 398 5a7 59a 78a 3 2 bcade 1 3
4 3 4 #78 889 98b a66 3 3 dbc 3 1
5 5 1 b72#1 8#21# 54a31 a#3b2 56a89 3 4 e 2 4
5 3 19 #7b #62 319 a25 447 3 3 ceb 4 2
5 4 13 a564 b#3b 7466 bba2 2732 5 2 ebcad 1 3
3 4 8 621# 8#bb 1bb1 3 2 cad 1 2
4 4 5 5983 9399 7941 261b 2 2 cbd 2 3
5 4 18 116# #356 3569 1517 9743 4 1 ac 3 3
4 3 20 a22 b89 5ab 7b1 4 2 d 4 2
3 4 4 2454 9277 #952 2 1 baecd 2 2
4 4 5 b859 2113 137b a3#6 3 3 eabc 3 3
4 3 16 #89 b6b 539 542 2 3 c 3 1
5 4 5 44b9 26#a 1266 11a9 3813 3 4 ba 5 1
4 4 7 8#71 #a#9 327a 3849 1 2 cedab 3 1
3 5 14 51545 14925 #b#a9 2 1 edba 2 2
3 5 19 781a# a#126 692b3 3 5 d 1 1
4 4 6 6#a5 a683 8593 a#67 1 2 eb 3 4
3 3 7 145 5b4 229 3 2 adeb 2 1
3 4 13 28a7 9b5# a983 1 4 aedcb 2 2
3 4 4 29a8 9286 4#52 1 2 cdaeb 3 4
3 4 4 b475 4#4b 1543 2 2 da 1 2
5 4 6 56#8 a#86 a14b 384b 1186 5 3 c 2 2
5 5 10 27438 9b317 35645 127a5 82176 4 3 bcae 5 3
5 3 15 921 163 4a1 6a9 48# 5 1 ecad 4 3
3 5 4 18924 57836 29925 1 1 e 2 2
4 5 8 1bb81 ba789 b7459 13#a6 3 5 cedab 3 2
3 3 4 7b4 547 #5b 2 3 cda 1 2
3 4 19 8469 #815 #621 1 1 cb 3 3
5 3 8 a65 aab a4# 174 a63 2 3 ab 4 1
5 5 2 15#a3 37449 5ba85 ab798 #32b3 2 1 e 2 1
5 3 19 b96 7a7 6#b 9a8 a5b 5 1 ecd 2 3
4 5 17 #7993 648a4 b159b 4b847 4 1 d 3 3
4 5 1 1b591 839aa ###4b 718a# 4 2 cbe 3 5
4 5 15 745#9 99415 b1267 4549a 2 2 dae 4 4
5 3 18 #45 55# 6a5 247 #88 3 2 ace 3 3
4 5 6 594b4 7#723 531a# 676a9 1 3 cbea 1 1
4 3 13 8#4 12# #46 8a9 3 1 cebad 3 2
5 5 8 78b#a 65#83 bb183 #b852 b3978 4 4 bae 2 4
5 4 11 a436 ##a1 2959 2##8 2921 2 2 dabc 5 2
4 5 1 39#55 519aa 55163 59423 2 5 c 2 5
5 5 18 a9192 54747 26651 52425 71554 3 2 bd 4 2
5 4 19 b3#5 2#76 91#6 49b3 6#b9 2 4 edbc 5 4
3 4 6 2518 9729 63b2 2 1 cbda 2 3
4 5 5 435#4 96b3# 39624 ba983 2 4 ebad 4 3
5 4 12 4247 46b9 373b 277# bb63 5 3 cba 4 1
3 4 5 2539 92#1 2494 2 2 cabed 1 3
5 5 18 8a418 8##5b 6975a 93166 a1948 3 2 dc 1 1
3 4 11 8996 8b11 9833 2 4 a 1 1
5 4 7 69ba 1635 11b2 43ab 2366 2 2 bace 1 3
5 4 9 8379 39#9 a146 a1#1 5519 3 1 aebcd 2 3
4 3 5 32# b9# 778 836 2 1 ca 1 2
4 4 7 a747 37a2 9396 5646 2 2 ec 4 1
5 5 9 #4736 3495# 5#6b# 9a8#b ##694 4 3 adbe 1 5
4 3 2 5aa a19 bb# a8a 1 3 edbac 1 1
3 3 5 833 794 63# 3 3 bd 2 1
5 5 15 13312 6a764 938#8 ba62b 7b92# 3 4 edba 5 4
3 4 18 835b 4386 4465 3 3 cba 1 1
3 4 18 2536 7348 5619 2 4 ac 3 1
4 3 13 63# 92# 665 6b6 1 2 eadcb 4 3
5 5 8 2b562 8473a 14293 32#b# 744b8 1 4 ab 4 5
5 3 17 651 5#a 7b8 346 #67 4 1 e 2 3
5 3 17 3ab b86 543 b9# 5a4 1 3 ad 5 3
3 4 12 462# b822 6bbb 3 2 b 2 1
3 5 12 47b94 48552 1a781 3 1 dbc 1 2
4 5 12 a548b #81a6 a19b6 71b1b 4 2 eabc 4 2
4 4 5 1a97 8593 8aaa #37b 4 2 ebd 1 2
3 5 12 bb462 #3943 74##6 2 5 ac 3 1
4 3 9 685 92b 153 a#6 4 2 baedc 3 1
4 3 6 73# 136 66a #29 3 3 e 3 3
3 4 16 619a b51a 7539 3 1 bcea 1 3
3 3 20 637 378 ab6 3 2 e 1 1
4 3 14 ab5 b21 995 793 4 3 ce 4 1
5 4 14 9#a7 174b 7a85 791a 4367 3 4 acd 2 4
5 3 19 2a8 4a1 2b# 89b ab8 4 3 ad 4 2
4 4 14 #b## 36a1 #679 314# 2 4 c 1 1
3 4 6 33a3 8a## 7752 3 1 b 1 1
4 4 1 7418 215b 5a8a 38a9 2 3 ca 2 2
3 4 14 6848 2156 12#6 3 1 cbaed 3 3
3 3 11 1ab ab1 447 1 1 cad 2 1
5 3 12 545 #16 454 #27 533 1 2 b 3 3
3 5 20 3814b 86#a5 37993 2 5 ce 1 4
4 5 20 9#6ab 15665 32566 b84bb 2 2 be 1 4
3 4 13 a119 89#5 2257 2 3 dbe 2 3
5 3 15 268 a55 1#b 2b5 5a3 1 3 aecbd 3 2
3 4 16 1313 3a93 55#4 3 3 eab 3 1
3 4 3 ba68 b2a7 78aa 2 1 d 2 1
4 5 1 38#a1 61627 99496 #8#93 2 5 db 2 4
3 5 2 87764 a7599 8aa12 3 1 cdae 2 4
3 5 11 67399 #8a39 53b24 3 5 eacb 3 2
4 3 3 956 2a6 871 465 3 2 abd 3 2
4 3 3 5#6 639 63a abb 4 2 ca 1 1
3 4 3 3239 b6a5 423a 3 4 eadb 1 2
4 5 9 3882a a5422 56b45 77629 4 2 dabc 3 4
4 5 2 575aa 711b8 4a496 5aa71 3 4 d 4 5
5 4 1 5#97 8963 ab#4 8839 83a1 5 4 cea 1 1
4 5 5 89#9a 99647 49b8a 34233 4 2 bca 2 3
3 4 7 3993 627b 491b 1 3 c 1 1
3 4 10 4764 3a85 891a 2 4 bdcea 2 1
4 5 5 a2121 4#ba3 a85b6 85b1b 4 1 adcb 2 5
5 5 18 31156 84a37 72aaa 99488 751b# 3 3 cadbe 4 2`

const INF = 1000000000

type item struct {
	id   int
	dist int
}

type priorityQueue []item

func (pq priorityQueue) Len() int           { return len(pq) } //nolint:revive
func (pq priorityQueue) Less(i, j int) bool { return pq[i].dist < pq[j].dist }
func (pq priorityQueue) Swap(i, j int)      { pq[i], pq[j] = pq[j], pq[i] }
func (pq *priorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(item))
}
func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	it := old[n-1]
	*pq = old[:n-1]
	return it
}

func dijkstra(m, n int, weight [][]int, passable [][]bool, src, dst int) ([]int, []int) {
	N := m * n
	dist := make([]int, N)
	prev := make([]int, N)
	for i := 0; i < N; i++ {
		dist[i] = INF
		prev[i] = -1
	}
	dist[src] = 0
	pq := &priorityQueue{}
	heap.Init(pq)
	heap.Push(pq, item{src, 0})
	dr := []int{-1, 1, 0, 0}
	dc := []int{0, 0, -1, 1}
	for pq.Len() > 0 {
		it := heap.Pop(pq).(item)
		u := it.id
		d := it.dist
		if d != dist[u] {
			continue
		}
		if u == dst {
			break
		}
		r := u / n
		c := u % n
		w := weight[r][c]
		for k := 0; k < 4; k++ {
			nr := r + dr[k]
			nc := c + dc[k]
			if nr < 0 || nr >= m || nc < 0 || nc >= n {
				continue
			}
			if !passable[nr][nc] {
				continue
			}
			v := nr*n + nc
			nd := d + w
			if nd < dist[v] {
				dist[v] = nd
				prev[v] = u
				heap.Push(pq, item{v, nd})
			}
		}
	}
	return dist, prev
}

type testCase struct {
	m, n, k int
	grid    []string
	rs, cs  int
	path    string
	re, ce  int
}

func filterPath(grid []string, path string) string {
	present := make(map[byte]bool)
	for _, row := range grid {
		for i := 0; i < len(row); i++ {
			ch := row[i]
			if ch != '#' && (ch < '1' || ch > '9') {
				present[ch] = true
			}
		}
	}
	var b strings.Builder
	for i := 0; i < len(path); i++ {
		if present[path[i]] {
			b.WriteByte(path[i])
		}
	}
	return b.String()
}

func parseCases() ([]testCase, error) {
	data := strings.TrimSpace(testcasesF)
	if data == "" {
		return nil, fmt.Errorf("no testcases provided")
	}
	lines := strings.Split(data, "\n")
	cases := make([]testCase, 0, len(lines))
	for i, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 8 {
			return nil, fmt.Errorf("line %d malformed", i+1)
		}
		m, err1 := strconv.Atoi(fields[0])
		n, err2 := strconv.Atoi(fields[1])
		k, err3 := strconv.Atoi(fields[2])
		if err1 != nil || err2 != nil || err3 != nil {
			return nil, fmt.Errorf("line %d parse error", i+1)
		}
		if len(fields) != m+8 {
			return nil, fmt.Errorf("line %d expected %d values, got %d", i+1, m+8, len(fields))
		}
		grid := make([]string, m)
		pos := 3
		for r := 0; r < m; r++ {
			grid[r] = fields[pos]
			pos++
		}
		rs, err4 := strconv.Atoi(fields[pos])
		cs, err5 := strconv.Atoi(fields[pos+1])
		path := filterPath(grid, fields[pos+2])
		if path == "" {
			path = "-"
		}
		re, err6 := strconv.Atoi(fields[pos+3])
		ce, err7 := strconv.Atoi(fields[pos+4])
		if err4 != nil || err5 != nil || err6 != nil || err7 != nil {
			return nil, fmt.Errorf("line %d coordinate parse error", i+1)
		}
		cases = append(cases, testCase{
			m:    m,
			n:    n,
			k:    k,
			grid: grid,
			rs:   rs,
			cs:   cs,
			path: path,
			re:   re,
			ce:   ce,
		})
	}
	return cases, nil
}

// solveOne is the embedded solution from 241F.go.
func solveOne(tc testCase) string {
	m, n, k := tc.m, tc.n, tc.k
	grid := tc.grid
	weight := make([][]int, m)
	passable := make([][]bool, m)
	junctions := make(map[byte]int)
	for i := 0; i < m; i++ {
		weight[i] = make([]int, n)
		passable[i] = make([]bool, n)
		for j := 0; j < n; j++ {
			ch := grid[i][j]
			if ch == '#' {
				passable[i][j] = false
			} else {
				passable[i][j] = true
				if ch >= '1' && ch <= '9' {
					weight[i][j] = int(ch - '0')
				} else {
					weight[i][j] = 1
					junctions[ch] = i*n + j
				}
			}
		}
	}
	rs, cs, re, ce := tc.rs-1, tc.cs-1, tc.re-1, tc.ce-1
	seq := make([]int, 0, len(tc.path)+2)
	start := rs*n + cs
	seq = append(seq, start)
	for i := 0; i < len(tc.path); i++ {
		id, ok := junctions[tc.path[i]]
		if !ok {
			continue
		}
		seq = append(seq, id)
	}
	dest := re*n + ce
	seq = append(seq, dest)

	path := make([]int, 0)
	for i := 0; i+1 < len(seq); i++ {
		u := seq[i]
		v := seq[i+1]
		_, prev := dijkstra(m, n, weight, passable, u, v)
		tmp := make([]int, 0)
		cur := v
		for cur != -1 && cur != u {
			tmp = append(tmp, cur)
			cur = prev[cur]
		}
		tmp = append(tmp, u)
		for l, r := 0, len(tmp)-1; l < r; l, r = l+1, r-1 {
			tmp[l], tmp[r] = tmp[r], tmp[l]
		}
		if i > 0 && len(tmp) > 0 {
			tmp = tmp[1:]
		}
		path = append(path, tmp...)
	}

	time := 0
	pos := path[0]
	for i := 0; i+1 < len(path); i++ {
		u := path[i]
		r := u / n
		c := u % n
		w := weight[r][c]
		if time+w > k {
			break
		}
		time += w
		pos = path[i+1]
	}
	rf := pos/n + 1
	cf := pos%n + 1
	return fmt.Sprintf("%d %d", rf, cf)
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", tc.m, tc.n, tc.k))
	for _, row := range tc.grid {
		sb.WriteString(row)
		sb.WriteByte('\n')
	}
	sb.WriteString(fmt.Sprintf("%d %d %s %d %d\n", tc.rs, tc.cs, tc.path, tc.re, tc.ce))
	return sb.String()
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierF /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases, err := parseCases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range cases {
		input := buildInput(tc)
		expect := solveOne(tc)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("test %d failed\nexpected: %s\n got: %s\n", idx+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

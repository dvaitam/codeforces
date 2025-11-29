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

// Embedded source for the reference solution (was 1250E.go).
const solutionSource = `package main

import (
	"bufio"
	"fmt"
	"os"
)

type DSU struct {
	parent []int
	rank   []int
	parity []int
}

func NewDSU(n int) *DSU {
	d := &DSU{
		parent: make([]int, n),
		rank:   make([]int, n),
		parity: make([]int, n),
	}
	for i := range d.parent {
		d.parent[i] = i
	}
	return d
}

func (d *DSU) find(x int) (int, int) {
	if d.parent[x] == x {
		return x, d.parity[x]
	}
	r, p := d.find(d.parent[x])
	d.parent[x] = r
	d.parity[x] ^= p
	return d.parent[x], d.parity[x]
}

func (d *DSU) union(x, y, rel int) bool {
	rx, px := d.find(x)
	ry, py := d.find(y)
	if rx == ry {
		return (px ^ py) == rel
	}
	if d.rank[rx] < d.rank[ry] {
		rx, ry = ry, rx
		px, py = py, px
	}
	d.parent[ry] = rx
	d.parity[ry] = px ^ py ^ rel
	if d.rank[rx] == d.rank[ry] {
		d.rank[rx]++
	}
	return true
}

func reverse(s string) string {
	b := []byte(s)
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	return string(b)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, m, k int
		if _, err := fmt.Fscan(in, &n, &m, &k); err != nil {
			return
		}
		s := make([]string, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &s[i])
		}
		r := make([]string, n)
		for i := 0; i < n; i++ {
			r[i] = reverse(s[i])
		}
		dsu := NewDSU(n)
		possible := true
		for i := 0; i < n && possible; i++ {
			for j := i + 1; j < n; j++ {
				same, diff := 0, 0
				si := s[i]
				sj := s[j]
				rj := r[j]
				for t2 := 0; t2 < m; t2++ {
					if si[t2] == sj[t2] {
						same++
					}
					if si[t2] == rj[t2] {
						diff++
					}
				}
				if same < k && diff < k {
					possible = false
					break
				} else if same >= k && diff < k {
					if !dsu.union(i, j, 0) {
						possible = false
						break
					}
				} else if same < k && diff >= k {
					if !dsu.union(i, j, 1) {
						possible = false
						break
					}
				}
			}
		}
		if !possible {
			fmt.Println(-1)
			continue
		}
		par := make([]int, n)
		comps := make(map[int][]int)
		for i := 0; i < n; i++ {
			root, p := dsu.find(i)
			par[i] = p
			comps[root] = append(comps[root], i)
		}
		var toReverse []int
		for _, list := range comps {
			count1 := 0
			for _, idx := range list {
				if par[idx] == 1 {
					count1++
				}
			}
			if count1 <= len(list)-count1 {
				for _, idx := range list {
					if par[idx] == 1 {
						toReverse = append(toReverse, idx+1)
					}
				}
			} else {
				for _, idx := range list {
					if par[idx] == 0 {
						toReverse = append(toReverse, idx+1)
					}
				}
			}
		}
		fmt.Println(len(toReverse))
		if len(toReverse) > 0 {
			for i, v := range toReverse {
				if i > 0 {
					fmt.Print(" ")
				}
				fmt.Print(v)
			}
		}
		fmt.Println()
	}
}
`

const testcasesRaw = `100
5 1 1
0
1
1
0
0
4 2 2
11
01
11
11
1 5 5
10100
4 4 3
1100
1111
0000
1111
1 4 3
0110
5 5 3
11101
01011
11010
01000
11101
1 1 1
1
2 2 2
01
01
1 2 1
01
4 5 5
10111
11000
10110
00001
2 2 2
11
01
3 2 1
11
11
01
1 1 1
0
5 3 3
000
000
000
110
111
5 1 1
0
1
1
1
0
5 3 2
010
110
000
010
011
4 3 2
011
101
010
110
3 3 2
101
000
110
1 3 1
011
5 1 1
1
1
0
1
0
5 2 2
11
11
10
10
01
1 5 5
10100
3 1 1
1
1
1
5 2 2
01
01
01
01
10
1 2 2
01
5 5 2
01111
00111
10111
10100
11010
5 3 3
101
100
010
010
000
2 3 3
110
011
1 4 3
1100
1 3 3
111
5 2 1
11
01
10
00
01
3 1 1
1
1
0
1 5 5
00110
1 1 1
1
2 4 1
1010
0010
2 2 2
10
00
5 3 2
001
110
100
011
000
3 3 1
100
100
110
3 2 1
00
00
11
5 5 2
10001
01011
01110
10011
11100
5 5 4
00110
11010
11010
10000
11110
5 4 1
1101
0010
0000
1100
0000
2 4 3
1001
0001
4 1 1
0
0
0
1
4 4 3
0010
1001
1111
1001
5 5 2
11011
11110
01100
01010
11011
4 1 1
1
0
0
0
5 3 2
011
000
011
111
011
5 5 1
00000
11011
00100
00001
11011
5 2 1
00
11
01
01
01
4 1 1
1
1
0
0
1 5 5
00111
2 5 1
11000
10111
4 2 1
00
10
10
00
1 5 3
11111
3 2 1
10
00
01
4 2 2
11
10
01
01
4 1 1
1
0
1
0
4 5 1
11110
01100
01100
01010
5 3 3
110
001
011
001
100
1 3 2
011
5 3 3
001
001
111
000
100
3 2 1
10
00
01
4 4 1
0110
0111
0101
1000
3 4 2
1100
0101
0111
2 2 1
10
11
2 1 1
1
0
4 4 1
0111
1010
1111
1010
4 4 3
1011
0100
1000
1110
1 4 3
1100
1 2 2
00
4 4 2
0100
1000
1010
0000
1 4 4
1000
5 1 1
0
1
1
0
1
3 1 1
0
1
1
5 2 1
10
00
10
01
11
5 3 3
001
100
010
110
000
1 2 2
00
3 3 2
001
011
000
3 3 2
100
111
100
4 3 1
001
100
000
110
3 2 2
01
01
01
5 3 2
010
111
100
111
100
1 4 2
0001
4 1 1
0
1
0
0
3 4 4
1111
1010
1011
3 1 1
0
1
1
2 1 1
1
0
1 4 2
0110
2 1 1
0
0
3 3 1
110
101
010
4 5 2
11101
00110
00010
10111
1 5 1
00010
4 2 1
01
11
11
10
2 5 1
11111
00011
4 2 1
01
00
10
10
3 4 2
0110
1101
1111
3 4 3
0011
0010
1101
2 4 3
0001
0010
2 2 1
11
10`

type testCase struct {
	n, m, k int
	rows    []string
}

type dsu struct {
	parent []int
	rank   []int
	parity []int
}

func newDSU(n int) *dsu {
	d := &dsu{
		parent: make([]int, n),
		rank:   make([]int, n),
		parity: make([]int, n),
	}
	for i := range d.parent {
		d.parent[i] = i
	}
	return d
}

func (d *dsu) find(x int) (int, int) {
	if d.parent[x] == x {
		return x, d.parity[x]
	}
	r, p := d.find(d.parent[x])
	d.parent[x] = r
	d.parity[x] ^= p
	return d.parent[x], d.parity[x]
}

func (d *dsu) union(x, y, rel int) bool {
	rx, px := d.find(x)
	ry, py := d.find(y)
	if rx == ry {
		return (px ^ py) == rel
	}
	if d.rank[rx] < d.rank[ry] {
		rx, ry = ry, rx
		px, py = py, px
	}
	d.parent[ry] = rx
	d.parity[ry] = px ^ py ^ rel
	if d.rank[rx] == d.rank[ry] {
		d.rank[rx]++
	}
	return true
}

func reverse(s string) string {
	b := []byte(s)
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	return string(b)
}

func expected(tc testCase) (string, error) {
	n, m, k := tc.n, tc.m, tc.k
	s := make([]string, n)
	copy(s, tc.rows)
	r := make([]string, n)
	for i := 0; i < n; i++ {
		r[i] = reverse(s[i])
	}
	dsu := newDSU(n)
	possible := true
	for i := 0; i < n && possible; i++ {
		for j := i + 1; j < n; j++ {
			same, diff := 0, 0
			si := s[i]
			sj := s[j]
			rj := r[j]
			for t2 := 0; t2 < m; t2++ {
				if si[t2] == sj[t2] {
					same++
				}
				if si[t2] == rj[t2] {
					diff++
				}
			}
			if same < k && diff < k {
				possible = false
				break
			} else if same >= k && diff < k {
				if !dsu.union(i, j, 0) {
					possible = false
					break
				}
			} else if same < k && diff >= k {
				if !dsu.union(i, j, 1) {
					possible = false
					break
				}
			}
		}
	}
	if !possible {
		return "-1", nil
	}
	par := make([]int, n)
	comps := make(map[int][]int)
	for i := 0; i < n; i++ {
		root, p := dsu.find(i)
		par[i] = p
		comps[root] = append(comps[root], i)
	}
	var toReverse []int
	for _, list := range comps {
		count1 := 0
		for _, idx := range list {
			if par[idx] == 1 {
				count1++
			}
		}
		if count1 <= len(list)-count1 {
			for _, idx := range list {
				if par[idx] == 1 {
					toReverse = append(toReverse, idx+1)
				}
			}
		} else {
			for _, idx := range list {
				if par[idx] == 0 {
					toReverse = append(toReverse, idx+1)
				}
			}
		}
	}
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(toReverse)))
	if len(toReverse) > 0 {
		sb.WriteByte('\n')
		for i, v := range toReverse {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
	}
	return sb.String(), nil
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func parseTestcases() ([]testCase, error) {
	scanner := bufio.NewScanner(strings.NewReader(testcasesRaw))
	scanner.Split(bufio.ScanWords)
	if !scanner.Scan() {
		return nil, fmt.Errorf("invalid test data")
	}
	t, _ := strconv.Atoi(scanner.Text())
	cases := make([]testCase, 0, t)
	for i := 0; i < t; i++ {
		if !scanner.Scan() {
			return nil, fmt.Errorf("unexpected EOF at case %d", i+1)
		}
		n, _ := strconv.Atoi(scanner.Text())
		if !scanner.Scan() {
			return nil, fmt.Errorf("missing m at case %d", i+1)
		}
		m, _ := strconv.Atoi(scanner.Text())
		if !scanner.Scan() {
			return nil, fmt.Errorf("missing k at case %d", i+1)
		}
		k, _ := strconv.Atoi(scanner.Text())
		rows := make([]string, n)
		for j := 0; j < n; j++ {
			if !scanner.Scan() {
				return nil, fmt.Errorf("missing row %d in case %d", j+1, i+1)
			}
			rows[j] = scanner.Text()
		}
		cases = append(cases, testCase{n: n, m: m, k: k, rows: rows})
	}
	return cases, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	for i, tc := range cases {
		var sb strings.Builder
		sb.WriteString("1\n")
		fmt.Fprintf(&sb, "%d %d %d\n", tc.n, tc.m, tc.k)
		for _, row := range tc.rows {
			sb.WriteString(row)
			sb.WriteByte('\n')
		}
		input := sb.String()
		want, err := expected(tc)
		if err != nil {
			fmt.Printf("case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(want) {
			fmt.Printf("case %d failed\nexpected:\n%s\ngot:\n%s\n", i+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

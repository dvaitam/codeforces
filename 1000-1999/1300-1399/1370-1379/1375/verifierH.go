package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type pair struct{ a, b int }

type testcase struct {
	n int
	q int
	v []int
	p [][2]int
}

// Embedded copy of testcasesH.txt so the verifier is self-contained.
const testcasesRaw = `1 2 1 1 1 1 1
4 3 4 1 2 3 3 4 2 4 4 4
2 3 1 2 1 2 2 2 2 2
3 2 2 1 3 1 3 3 3
1 1 1 1 1
4 1 4 2 3 1 2 4
4 2 2 1 4 3 1 2 1 3
3 3 3 2 1 2 2 1 3 1 3
1 1 1 1 1
3 3 1 2 3 2 3 3 3 2 2
4 1 1 3 4 2 1 4
2 1 1 2 2 2
1 3 1 1 1 1 1 1 1
3 2 3 1 2 1 1 3 3
3 3 3 1 2 2 3 2 3 2 3
4 2 1 2 3 4 4 4 4 4
3 3 1 3 2 2 2 3 3 3 3
3 3 1 3 2 3 3 2 2 3 3
1 1 1 1 1
3 2 3 1 2 1 1 2 2
4 3 3 2 1 4 2 2 3 4 2 4
2 3 1 2 1 1 1 1 1 2
4 3 2 4 3 1 2 4 1 2 3 4
4 3 2 4 3 1 2 3 4 4 4 4
4 3 4 3 2 1 3 4 1 1 3 4
3 2 1 3 2 2 2 3 3
2 3 2 1 1 1 2 2 1 1
2 3 1 2 2 2 1 1 2 2
4 3 2 1 3 4 2 2 3 3 2 4
2 1 1 2 1 2
2 1 2 1 2 2
3 1 2 1 3 3 3
2 2 1 2 2 2 2 2
2 3 1 2 2 2 1 1 1 2
1 1 1 1 1
4 1 2 4 3 1 3 3
1 3 1 1 1 1 1 1 1
1 3 1 1 1 1 1 1 1
1 3 1 1 1 1 1 1 1
2 3 2 1 2 2 2 2 1 1
2 2 1 2 2 1 1 1
4 1 2 3 4 2 1 4
3 2 1 3 3 3 3 3 3
4 3 3 3 3 1 3 2 2 2 2 4
4 1 2 3 4 4 2 3
2 2 1 1 1 2 2 2
2 2 2 1 1 1 1 1
2 1 2 1 2 1
4 2 3 2 3 4 3 1 3 1
2 2 2 2 2 2 2 1
2 1 1 2 1 1
3 3 1 1 1 1 2 2 2 1 1
4 1 2 3 4 3 2 2
4 2 3 4 1 1 3 1 2 3 3 4
4 2 1 4 4 3 3 4 3 1 3 4
4 1 1 3 4 2 1 4
4 3 4 3 2 1 2 4 1 2 2 4
1 3 1 1 1 1 1 1 1
4 3 1 1 3 4 2 4 2 3 1 1
2 1 2 2 1 1
3 1 2 2 3 3 3 1 1
4 1 2 3 4 1 2 3
3 2 2 2 1 2 2 2 2
2 1 2 1 1 1
4 3 2 1 3 4 3 3 3 2 1 4
3 1 3 3 1 2 2 3 3
4 2 1 4 3 2 3 4 3 4 3 4
4 1 2 4 4 2 1 2
3 3 2 3 1 2 3 3 3 2 3
2 3 1 2 1 2 2 2 2 2
4 3 3 2 1 2 3 3 3 3 1 2
1 3 1 1 1 1 1 1 1
1 1 1 1 1
2 2 2 1 1 1 2 2
2 3 1 2 1 1 2 1 1 1
3 1 2 2 3 3 3 3 3
4 3 1 4 1 2 3 4 1 2 1 4
4 3 4 1 2 3 3 4 2 4 2 4
4 3 3 2 1 2 2 2 2 2 2 2
4 1 1 1 3 1 2 4
3 3 2 2 3 2 3 3 2 3 3
1 1 1 1 1
4 1 1 2 4 3 4 2
1 3 1 1 1 1 1 1 1
4 3 1 4 2 1 3 2 3 1 1 3
4 3 1 1 2 4 2 1 4 4 2 2
4 3 1 3 4 2 2 2 2 4 1 4
4 1 2 1 3 4 4 2
2 2 1 2 1 2 1 1
2 2 1 2 1 2 2 1
2 1 1 1 2 1
1 3 1 1 1 1 1 1 1
4 2 2 2 3 1 2 4 1 4 3 3
1 1 1 1 1
1 1 1 1 1
2 3 2 1 1 2 1 2 2 1
4 3 3 3 1 3 4 3 2 1 2 4
2 2 2 2 1 1 1 1
2 3 1 2 1 1 1 2 2 2
4 2 2 3 3 1 2 3 1 2 1 2
4 3 4 1 2 3 2 4 3 4 3 4
4 3 2 3 1 4 2 2 2 4 4 4
2 3 2 2 1 1 1 2 1 1
4 3 1 2 4 3 4 2 2 3 1 2
1 3 1 1 1 1 1 1 1
3 1 2 1 3 3 3 3 1
2 1 2 1 1 2
1 1 1 1 1`

// solve implements the logic from 1375H.go for a single testcase.
func solve(tc testcase) string {
	const B = 256
	n, q := tc.n, tc.q
	v := make([]int, n)
	copy(v, tc.v)
	for i := range v {
		v[i]--
	}

	type comp struct{ to []int }
	newComp := func(pos []int, low, up int) *comp {
		c := &comp{to: make([]int, len(pos)+1)}
		for i, p := range pos {
			c.to[i+1] = c.to[i]
			x := v[p]
			if x >= low && x < up {
				c.to[i+1]++
			}
		}
		return c
	}
	down := func(c *comp, l, r int) (int, int) {
		if l >= len(c.to) {
			l = len(c.to) - 1
		}
		if r >= len(c.to) {
			r = len(c.to) - 1
		}
		return c.to[l], c.to[r]
	}

	ansPairs := []pair{}
	cnt := 0
	query := func(a, b int) int {
		if a == -1 {
			return b
		}
		if b == -1 {
			return a
		}
		ansPairs = append(ansPairs, pair{a, b})
		id := cnt
		cnt++
		return id
	}

	var gen func(low, up int, pos []int) [][]int
	gen = func(low, up int, pos []int) [][]int {
		m := len(pos)
		res := make([][]int, m+1)
		for i := 0; i <= m; i++ {
			row := make([]int, m+1)
			for j := range row {
				row[j] = -1
			}
			res[i] = row
		}
		if m == 0 {
			return res
		}
		if m == 1 || up-low <= 1 {
			for i := 0; i < m; i++ {
				res[i][i+1] = pos[i]
			}
			for length := 2; length <= m; length++ {
				for i := 0; i+length <= m; i++ {
					j := i + length
					res[i][j] = query(res[i][j-1], pos[j-1])
				}
			}
			return res
		}
		mid := (low + up) / 2
		lcomp := newComp(pos, low, mid)
		ucomp := newComp(pos, mid, up)
		var lpos, upos []int
		for _, p := range pos {
			if v[p] < mid {
				lpos = append(lpos, p)
			} else {
				upos = append(upos, p)
			}
		}
		lres := gen(low, mid, lpos)
		ures := gen(mid, up, upos)
		for i := 0; i < m; i++ {
			for j := i + 1; j <= m; j++ {
				ll, lr := down(lcomp, i, j)
				ul, ur := down(ucomp, i, j)
				res[i][j] = query(lres[ll][lr], ures[ul][ur])
			}
		}
		return res
	}

	BN := (n + B - 1) / B
	poss := make([][]int, BN)
	for i := 0; i < n; i++ {
		bi := v[i] / B
		poss[bi] = append(poss[bi], i)
	}
	idx := make([]int, n)
	for i := 0; i < n; i++ {
		idx[i] = i
	}
	comps := make([]*comp, BN)
	bks := make([][][]int, BN)
	for i := 0; i < BN; i++ {
		bks[i] = gen(i*B, i*B+B, poss[i])
		comps[i] = newComp(idx, i*B, i*B+B)
	}

	qans := make([]int, q)
	for i := range qans {
		qans[i] = -1
	}
	for ti := 0; ti < q; ti++ {
		l, r := tc.p[ti][0]-1, tc.p[ti][1]
		for i := 0; i < BN; i++ {
			ll, rr := down(comps[i], l, r)
			qans[ti] = query(qans[ti], bks[i][ll][rr])
		}
	}

	var out strings.Builder
	fmt.Fprintln(&out, cnt)
	for _, p := range ansPairs {
		fmt.Fprintf(&out, "%d %d\n", p.a+1, p.b+1)
	}
	for i, x := range qans {
		if i > 0 {
			out.WriteByte(' ')
		}
		out.WriteString(strconv.Itoa(x + 1))
	}
	return strings.TrimSpace(out.String())
}

func buildIfGo(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "solbin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput()
		if err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("build failed: %v\n%s", err, out)
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}

	words := strings.Fields(testcasesRaw)
	pos := 0
	nextInt := func() (int, error) {
		if pos >= len(words) {
			return 0, fmt.Errorf("unexpected EOF at token %d", pos)
		}
		v, err := strconv.Atoi(words[pos])
		pos++
		return v, err
	}
	var cases []testcase
	for pos < len(words) {
		n, err := nextInt()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		q, err := nextInt()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		v := make([]int, n)
		for i := 0; i < n; i++ {
			val, err := nextInt()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			v[i] = val
		}
		p := make([][2]int, q)
		for i := 0; i < q; i++ {
			l, err := nextInt()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			r, err := nextInt()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			p[i] = [2]int{l, r}
		}
		cases = append(cases, testcase{n: n, q: q, v: v, p: p})
	}

	bin, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	for i, tc := range cases {
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.q)
		for j, val := range tc.v {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(val))
		}
		sb.WriteByte('\n')
		for _, pr := range tc.p {
			fmt.Fprintf(&sb, "%d %d\n", pr[0], pr[1])
		}
		input := sb.String()
		expect := solve(tc)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Printf("case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Printf("case %d failed\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", i+1, input, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

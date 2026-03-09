package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Pair struct{ first, second int }

func (a Pair) add(b Pair) Pair { return Pair{a.first + b.first, a.second + b.second} }
func (a Pair) sub(b Pair) Pair { return Pair{a.first - b.first, a.second - b.second} }
func (a Pair) less(b Pair) bool {
	if a.first != b.first {
		return a.first < b.first
	}
	return a.second < b.second
}

var (
	n        int
	fArr, s  []int
	tag, pre []int
	c        [][]int
	dp, cdp  [][]Pair
	cpre     []int
	ans      Pair
	todo     []Pair
)

func gao(v int) {
	if tag[v] == 0 {
		tag[v] = 1
	}
	dp[v][0] = Pair{0, 0}
	dp[v][1] = Pair{0, 0}
	pre[v] = -1
	for _, w := range c[v] {
		if tag[w] == -1 {
			continue
		}
		gao(w)
		dp[v][0] = dp[v][0].add(dp[w][1])
		diff := dp[w][1].sub(dp[w][0])
		base := Pair{1, 0}
		if s[v] != s[w] {
			base.second = 1
		}
		tmp := base.sub(diff)
		if dp[v][1].less(tmp) {
			dp[v][1] = tmp
			pre[v] = w
		}
	}
	dp[v][1] = dp[v][1].add(dp[v][0])
}

func dump(v int, flag bool, ret *[]Pair) {
	for _, w := range c[v] {
		if tag[w] == -1 {
			continue
		}
		if !flag || w != pre[v] {
			dump(w, true, ret)
		} else {
			*ret = append(*ret, Pair{v, w})
			dump(w, false, ret)
		}
	}
}

func solve(v0 int) {
	v := v0
	for tag[v] == 0 {
		tag[v] = 1
		v = fArr[v]
	}
	circ := []int{v}
	for u := fArr[v]; u != v; u = fArr[u] {
		circ = append(circ, u)
	}
	for _, u := range circ {
		tag[u] = -1
	}
	for _, u := range circ {
		gao(u)
	}
	m := len(circ)
	cdp[circ[0]][0] = dp[circ[0]][0]
	cdp[circ[0]][1] = dp[circ[0]][1]
	cpre[circ[0]] = -1
	for i := 1; i < m; i++ {
		u := circ[i]
		p := circ[i-1]
		cdp[u][0] = cdp[p][1].add(dp[u][0])
		cdp[u][1] = cdp[p][1].add(dp[u][1])
		cpre[u] = -1
		base := Pair{1, 0}
		if s[p] != s[u] {
			base.second = 1
		}
		tmp := cdp[p][0].add(base).add(dp[u][0])
		if cdp[u][1].less(tmp) {
			cdp[u][1] = tmp
			cpre[u] = 1
		}
	}
	best := cdp[circ[m-1]][1]
	how := []Pair{}
	for i := m - 1; i >= 0; {
		u := circ[i]
		if cpre[u] == -1 {
			dump(u, true, &how)
			i--
		} else {
			prev := circ[i-1]
			dump(u, false, &how)
			dump(prev, false, &how)
			how = append(how, Pair{u, prev})
			i -= 2
		}
	}
	circ = append(circ[1:], circ[0])
	cdp[circ[0]][0] = dp[circ[0]][0]
	cdp[circ[0]][1] = dp[circ[0]][1]
	cpre[circ[0]] = -1
	for i := 1; i < m; i++ {
		u := circ[i]
		p := circ[i-1]
		cdp[u][0] = cdp[p][1].add(dp[u][0])
		cdp[u][1] = cdp[p][1].add(dp[u][1])
		cpre[u] = -1
		base := Pair{1, 0}
		if s[p] != s[u] {
			base.second = 1
		}
		tmp := cdp[p][0].add(base).add(dp[u][0])
		if cdp[u][1].less(tmp) {
			cdp[u][1] = tmp
			cpre[u] = 1
		}
	}
	cand := cdp[circ[m-1]][1]
	if best.less(cand) {
		best = cand
		how = how[:0]
		for i := m - 1; i >= 0; {
			u := circ[i]
			if cpre[u] == -1 {
				dump(u, true, &how)
				i--
			} else {
				prev := circ[i-1]
				dump(u, false, &how)
				dump(prev, false, &how)
				how = append(how, Pair{u, prev})
				i -= 2
			}
		}
	}
	ans = ans.add(best)
	todo = append(todo, how...)
}

// reference returns the optimal (t, e) for the given input.
func reference(nVal int, fa, sx []int) (int, int) {
	n = nVal
	fArr = fa
	s = sx
	tag = make([]int, n+1)
	pre = make([]int, n+1)
	c = make([][]int, n+1)
	dp = make([][]Pair, n+1)
	cdp = make([][]Pair, n+1)
	cpre = make([]int, n+1)
	for i := 1; i <= n; i++ {
		dp[i] = make([]Pair, 2)
		cdp[i] = make([]Pair, 2)
	}
	for i := 1; i <= n; i++ {
		c[fArr[i]] = append(c[fArr[i]], i)
	}
	ans = Pair{0, 0}
	todo = []Pair{}
	for i := 1; i <= n; i++ {
		if tag[i] == 0 {
			solve(i)
		}
	}
	return ans.first, ans.second
}

func runBin(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
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

func makeInput(nVal int, fa, sx []int) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", nVal)
	for i := 1; i <= nVal; i++ {
		fmt.Fprintf(&sb, "%d %d\n", fa[i], sx[i])
	}
	return sb.String()
}

// validate checks that the candidate output is a valid answer.
func validate(output string, nVal int, fa, sx []int, optT, optE int) error {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) == 0 {
		return fmt.Errorf("empty output")
	}
	parts := strings.Fields(lines[0])
	if len(parts) != 2 {
		return fmt.Errorf("first line must have 2 numbers, got %q", lines[0])
	}
	gotT, err1 := strconv.Atoi(parts[0])
	gotE, err2 := strconv.Atoi(parts[1])
	if err1 != nil || err2 != nil {
		return fmt.Errorf("cannot parse first line: %q", lines[0])
	}
	if gotT != optT {
		return fmt.Errorf("wrong number of pairs: got %d, want %d", gotT, optT)
	}
	if gotE != optE {
		return fmt.Errorf("wrong number of boy-girl pairs: got %d, want %d", gotE, optE)
	}
	if len(lines) != 1+gotT {
		return fmt.Errorf("expected %d pair lines, got %d", gotT, len(lines)-1)
	}
	used := make(map[int]bool)
	boyGirl := 0
	for i := 0; i < gotT; i++ {
		pparts := strings.Fields(lines[1+i])
		if len(pparts) != 2 {
			return fmt.Errorf("pair line %d: expected 2 numbers, got %q", i+1, lines[1+i])
		}
		a, e1 := strconv.Atoi(pparts[0])
		b, e2 := strconv.Atoi(pparts[1])
		if e1 != nil || e2 != nil {
			return fmt.Errorf("pair line %d: cannot parse numbers", i+1)
		}
		if a < 1 || a > nVal || b < 1 || b > nVal {
			return fmt.Errorf("pair line %d: student index out of range (%d, %d)", i+1, a, b)
		}
		if a == b {
			return fmt.Errorf("pair line %d: student paired with themselves (%d)", i+1, a)
		}
		// pair {a,b} is valid if f[a]=b or f[b]=a
		if fa[a] != b && fa[b] != a {
			return fmt.Errorf("pair line %d: invalid pair (%d, %d): neither is the other's best friend", i+1, a, b)
		}
		if used[a] {
			return fmt.Errorf("pair line %d: student %d appears more than once", i+1, a)
		}
		if used[b] {
			return fmt.Errorf("pair line %d: student %d appears more than once", i+1, b)
		}
		used[a] = true
		used[b] = true
		if sx[a] != sx[b] {
			boyGirl++
		}
	}
	if boyGirl != gotE {
		return fmt.Errorf("declared e=%d but counted %d boy-girl pairs", gotE, boyGirl)
	}
	return nil
}

func generateTest(rng *rand.Rand) (int, []int, []int) {
	nVal := rng.Intn(10) + 2
	fa := make([]int, nVal+1)
	sx := make([]int, nVal+1)
	for i := 1; i <= nVal; i++ {
		for {
			fa[i] = rng.Intn(nVal) + 1
			if fa[i] != i {
				break
			}
		}
		sx[i] = rng.Intn(2) + 1
	}
	return nVal, fa, sx
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Fixed test cases from problem statement
	type fixedCase struct {
		nVal int
		fa   []int // 1-indexed
		sx   []int // 1-indexed
	}
	fixed := []fixedCase{
		// Sample: n=4, pairs (1->2,s=1),(2->3,s=2),(3->4,s=1),(4->2,s=1)
		{4, []int{0, 2, 3, 4, 2}, []int{0, 1, 2, 1, 1}},
	}

	total := 0
	for _, fc := range fixed {
		optT, optE := reference(fc.nVal, fc.fa, fc.sx)
		input := makeInput(fc.nVal, fc.fa, fc.sx)
		got, err := runBin(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fixed case %d failed: %v\ninput:\n%s", total+1, err, input)
			os.Exit(1)
		}
		if err := validate(got, fc.nVal, fc.fa, fc.sx, optT, optE); err != nil {
			fmt.Fprintf(os.Stderr, "fixed case %d failed: %v\noutput:\n%s\ninput:\n%s", total+1, err, got, input)
			os.Exit(1)
		}
		total++
	}

	for total < 200 {
		nVal, fa, sx := generateTest(rng)
		optT, optE := reference(nVal, fa, sx)
		input := makeInput(nVal, fa, sx)
		got, err := runBin(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", total+1, err, input)
			os.Exit(1)
		}
		if err := validate(got, nVal, fa, sx, optT, optE); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\noutput:\n%s\ninput:\n%s", total+1, err, got, input)
			os.Exit(1)
		}
		total++
	}
	fmt.Printf("All %d tests passed\n", total)
}

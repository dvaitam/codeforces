package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

const INF int64 = 1 << 60

type Fenwick struct {
	n    int
	tree []int64
}

func NewFenwick(n int) *Fenwick {
	return &Fenwick{n: n, tree: make([]int64, n+2)}
}

func (f *Fenwick) add(idx int, val int64) {
	for idx <= f.n+1 {
		f.tree[idx] += val
		idx += idx & -idx
	}
}

func (f *Fenwick) rangeAdd(l, r int, val int64) {
	if l > r {
		return
	}
	f.add(l+1, val)
	f.add(r+2, -val)
}

func (f *Fenwick) sum(idx int) int64 {
	res := int64(0)
	for idx > 0 {
		res += f.tree[idx]
		idx -= idx & -idx
	}
	return res
}

func (f *Fenwick) pointQuery(i int) int64 {
	return f.sum(i + 1)
}

func expected(a []int, p []int, b []int) string {
	m := len(b)
	dpBase := make([]int64, m+1)
	for i := 1; i <= m; i++ {
		dpBase[i] = INF
	}
	bit := NewFenwick(m + 1)

	query := func(idx int) int64 { return dpBase[idx] + bit.pointQuery(idx) }
	setVal := func(idx int, val int64) { dpBase[idx] = val - bit.pointQuery(idx) }
	rangeAdd := func(l, r int, val int64) { bit.rangeAdd(l, r, val) }

	for i := 0; i < len(a); i++ {
		x := a[i]
		cost := int64(p[i])
		pos := sort.Search(len(b), func(j int) bool { return b[j] >= x })
		if pos == m {
			rangeAdd(0, m, cost)
			continue
		}
		old := query(pos)
		if pos > 0 {
			rangeAdd(0, pos-1, cost)
		}
		addVal := cost
		if addVal > 0 {
			addVal = 0
		}
		rangeAdd(pos, m, addVal)
		if b[pos] == x {
			cur := query(pos + 1)
			if old < cur {
				setVal(pos+1, old)
			}
		}
	}

	ans := query(m)
	if ans >= INF/2 {
		return "NO"
	}
	return fmt.Sprintf("YES\n%d", ans)
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for caseNum := 0; caseNum < 100; caseNum++ {
		n := rng.Intn(5) + 1
		a := make([]int, n)
		pArr := make([]int, n)
		for i := 0; i < n; i++ {
			a[i] = rng.Intn(10) + 1
			pArr[i] = rng.Intn(11) - 5
		}
		m := rng.Intn(n) + 1
		b := make([]int, m)
		used := map[int]bool{}
		for i := 0; i < m; i++ {
			val := rng.Intn(10) + 1
			for used[val] {
				val = rng.Intn(10) + 1
			}
			used[val] = true
			b[i] = val
		}
		sort.Ints(b)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i, v := range a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteString("\n")
		for i, v := range pArr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteString("\n")
		sb.WriteString(fmt.Sprintf("%d\n", m))
		for i, v := range b {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteString("\n")
		input := sb.String()
		exp := expected(append([]int(nil), a...), append([]int(nil), pArr...), append([]int(nil), b...))
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", caseNum+1, err, input)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", caseNum+1, exp, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

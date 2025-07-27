package main

import (
	"bytes"
	"container/heap"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type Test struct {
	n int
	a []int
	b []int
}

func (t Test) Input() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t.n))
	for i := 0; i < t.n; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", t.a[i], t.b[i]))
	}
	return sb.String()
}

type Person struct {
	id  int
	a   int
	b   int
	ans int
}

type Item struct {
	b   int
	idx int
}

type PQ []Item

func (pq PQ) Len() int { return len(pq) }
func (pq PQ) Less(i, j int) bool {
	if pq[i].b == pq[j].b {
		return pq[i].idx < pq[j].idx
	}
	return pq[i].b < pq[j].b
}
func (pq PQ) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PQ) Push(x interface{}) { *pq = append(*pq, x.(Item)) }
func (pq *PQ) Pop() interface{} {
	old := *pq
	n := len(old)
	x := old[n-1]
	*pq = old[:n-1]
	return x
}

type Pair struct {
	b  int
	id int
}

func pairMax(x, y Pair) Pair {
	if x.b > y.b {
		return x
	} else if x.b < y.b {
		return y
	}
	if x.id > y.id {
		return x
	}
	return y
}

func expected(t Test) string {
	n := t.n
	p := make([]Person, n)
	for i := 0; i < n; i++ {
		p[i] = Person{id: i + 1, a: t.a[i], b: t.b[i]}
	}
	sort.Slice(p, func(i, j int) bool { return p[i].a < p[j].a })
	pq := &PQ{}
	heap.Init(pq)
	j := 0
	for i := 1; i <= n; i++ {
		for j < n && p[j].a <= i {
			heap.Push(pq, Item{b: p[j].b, idx: j})
			j++
		}
		item := heap.Pop(pq).(Item)
		p[item.idx].ans = i
	}
	sort.Slice(p, func(i, j int) bool { return p[i].ans < p[j].ans })
	ft := make([]Pair, n+2)
	update := func(u, b, id int) {
		for u > 0 {
			if pairMax(ft[u], Pair{b, id}) != ft[u] {
				ft[u] = pairMax(ft[u], Pair{b, id})
			}
			u -= u & -u
		}
	}
	get := func(u int) Pair {
		res := Pair{0, 0}
		for u <= n {
			res = pairMax(res, ft[u])
			u += u & -u
		}
		return res
	}
	resU, resV := 0, 0
	for i := 0; i < n; i++ {
		mx := get(p[i].a)
		if mx.b >= p[i].ans {
			resU = p[i].id
			resV = mx.id
			break
		}
		update(p[i].ans, p[i].b, p[i].id)
	}
	sort.Slice(p, func(i, j int) bool { return p[i].id < p[j].id })
	var sb strings.Builder
	if resU != 0 {
		sb.WriteString("NO\n")
		for i := 0; i < n; i++ {
			sb.WriteString(fmt.Sprintf("%d", p[i].ans))
			if i+1 == n {
				sb.WriteString("\n")
			} else {
				sb.WriteString(" ")
			}
		}
		for i := range p {
			if p[i].id == resU {
				for j := range p {
					if p[j].id == resV {
						p[i].ans, p[j].ans = p[j].ans, p[i].ans
						break
					}
				}
				break
			}
		}
	} else {
		sb.WriteString("YES\n")
	}
	for i := 0; i < n; i++ {
		sb.WriteString(fmt.Sprintf("%d", p[i].ans))
		if i+1 == n {
			sb.WriteString("\n")
		} else {
			sb.WriteString(" ")
		}
	}
	return strings.TrimSpace(sb.String())
}

func runProg(bin, input string) (string, error) {
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

func genTest(rng *rand.Rand) Test {
	n := rng.Intn(10) + 1
	perm := rng.Perm(n)
	a := make([]int, n)
	b := make([]int, n)
	for i := 0; i < n; i++ {
		val := perm[i] + 1
		a[i] = rng.Intn(val) + 1
		b[i] = val + rng.Intn(n-val+1)
	}
	return Test{n: n, a: a, b: b}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	const cases = 100
	for i := 0; i < cases; i++ {
		tc := genTest(rng)
		expect := expected(tc)
		got, err := runProg(bin, tc.Input())
		if err != nil {
			fmt.Printf("case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Printf("case %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i+1, tc.Input(), expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", cases)
}

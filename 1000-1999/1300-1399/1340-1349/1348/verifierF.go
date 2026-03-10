package main

import (
	"bytes"
	"container/heap"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
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

// refSolve returns (isUnique, uniqueOrGreedyPerm).
// If unique, returns true and the unique permutation (1-indexed by position).
// If not unique, returns false and the greedy permutation.
func refSolve(t Test) (bool, []int) {
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

	// Check uniqueness via Fenwick tree (find overlapping interval pair)
	sort.Slice(p, func(i, j int) bool { return p[i].ans < p[j].ans })
	type Pair struct{ b, id int }
	ft := make([]Pair, n+2)
	update := func(u, b, id int) {
		for u > 0 {
			if b > ft[u].b || (b == ft[u].b && id > ft[u].id) {
				ft[u] = Pair{b, id}
			}
			u -= u & -u
		}
	}
	get := func(u int) Pair {
		res := Pair{0, 0}
		for u <= n {
			if ft[u].b > res.b || (ft[u].b == res.b && ft[u].id > res.id) {
				res = ft[u]
			}
			u += u & -u
		}
		return res
	}
	resU := 0
	for i := 0; i < n; i++ {
		mx := get(p[i].a)
		if mx.b >= p[i].ans {
			resU = p[i].id
			break
		}
		update(p[i].ans, p[i].b, p[i].id)
	}

	sort.Slice(p, func(i, j int) bool { return p[i].id < p[j].id })
	perm := make([]int, n)
	for i := range p {
		perm[i] = p[i].ans
	}
	return resU == 0, perm
}

func parsePerm(s string, n int) ([]int, error) {
	fields := strings.Fields(s)
	if len(fields) != n {
		return nil, fmt.Errorf("expected %d values, got %d", n, len(fields))
	}
	perm := make([]int, n)
	for i, f := range fields {
		v, err := strconv.Atoi(f)
		if err != nil {
			return nil, fmt.Errorf("parse error: %v", err)
		}
		perm[i] = v
	}
	return perm, nil
}

func isValidPerm(tc Test, perm []int) bool {
	n := tc.n
	if len(perm) != n {
		return false
	}
	seen := make([]bool, n+1)
	for i, v := range perm {
		if v < 1 || v > n || seen[v] {
			return false
		}
		seen[v] = true
		if v < tc.a[i] || v > tc.b[i] {
			return false
		}
	}
	return true
}

func validate(tc Test, unique bool, refPerm []int, got string) error {
	lines := strings.Split(strings.TrimSpace(got), "\n")
	if len(lines) == 0 {
		return fmt.Errorf("empty output")
	}
	answer := strings.TrimSpace(lines[0])
	if unique {
		if answer != "YES" {
			return fmt.Errorf("expected YES got %s", answer)
		}
		if len(lines) < 2 {
			return fmt.Errorf("YES but no permutation given")
		}
		perm, err := parsePerm(lines[1], tc.n)
		if err != nil {
			return fmt.Errorf("bad permutation: %v", err)
		}
		if !isValidPerm(tc, perm) {
			return fmt.Errorf("permutation is not valid for the constraints")
		}
		for i, v := range perm {
			if v != refPerm[i] {
				return fmt.Errorf("incorrect answer at position %d (unique answer exists)", i+1)
			}
		}
	} else {
		if answer != "NO" {
			return fmt.Errorf("expected NO got %s", answer)
		}
		if len(lines) < 3 {
			return fmt.Errorf("NO but fewer than 2 permutations given")
		}
		perm1, err := parsePerm(lines[1], tc.n)
		if err != nil {
			return fmt.Errorf("bad first permutation: %v", err)
		}
		perm2, err := parsePerm(lines[2], tc.n)
		if err != nil {
			return fmt.Errorf("bad second permutation: %v", err)
		}
		if !isValidPerm(tc, perm1) {
			return fmt.Errorf("first permutation is not valid for the constraints")
		}
		if !isValidPerm(tc, perm2) {
			return fmt.Errorf("second permutation is not valid for the constraints")
		}
		same := true
		for i := range perm1 {
			if perm1[i] != perm2[i] {
				same = false
				break
			}
		}
		if same {
			return fmt.Errorf("two permutations are identical")
		}
	}
	return nil
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
		unique, refPerm := refSolve(tc)
		got, err := runProg(bin, tc.Input())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if err := validate(tc, unique, refPerm, got); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%serr: %v\ngot:\n%s\n", i+1, tc.Input(), err, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", cases)
}

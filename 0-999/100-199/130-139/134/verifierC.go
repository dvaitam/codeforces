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

type Item struct {
	id  int
	deg int
}

type PQ []*Item

func (pq PQ) Len() int            { return len(pq) }
func (pq PQ) Less(i, j int) bool  { return pq[i].deg > pq[j].deg }
func (pq PQ) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PQ) Push(x interface{}) { *pq = append(*pq, x.(*Item)) }
func (pq *PQ) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	*pq = old[:n-1]
	return item
}

type Edge struct{ u, v int }

func oracleSolve(input string) string {
	words := strings.Fields(input)
	idx := 0
	nextInt := func() int {
		v, _ := strconv.Atoi(words[idx])
		idx++
		return v
	}
	n := nextInt()
	s := nextInt()

	pq := make(PQ, 0, n)
	sum := 0
	for i := 1; i <= n; i++ {
		deg := nextInt()
		if deg > 0 {
			pq = append(pq, &Item{id: i, deg: deg})
			sum += deg
		}
	}

	if sum%2 != 0 || sum != s {
		return "No"
	}

	heap.Init(&pq)

	edges := make([]Edge, 0, sum/2)

	for pq.Len() > 0 {
		u := heap.Pop(&pq).(*Item)
		d := u.deg
		if d == 0 {
			continue
		}
		if pq.Len() < d {
			return "No"
		}

		popped := make([]*Item, d)
		for i := 0; i < d; i++ {
			popped[i] = heap.Pop(&pq).(*Item)
		}

		for i := 0; i < d; i++ {
			edges = append(edges, Edge{u.id, popped[i].id})
			popped[i].deg--
		}

		for i := 0; i < d; i++ {
			if popped[i].deg > 0 {
				heap.Push(&pq, popped[i])
			}
		}
	}

	var sb strings.Builder
	sb.WriteString("Yes\n")
	sb.WriteString(strconv.Itoa(len(edges)) + "\n")
	for _, e := range edges {
		sb.WriteString(strconv.Itoa(e.u) + " " + strconv.Itoa(e.v) + "\n")
	}
	return strings.TrimSpace(sb.String())
}

// Validate candidate output: check it's a valid graph with the right degree sequence.
func validateCandidate(input, candOutput string) error {
	words := strings.Fields(input)
	idx := 0
	nextInt := func() int {
		v, _ := strconv.Atoi(words[idx])
		idx++
		return v
	}
	n := nextInt()
	s := nextInt()
	a := make([]int, n+1)
	sum := 0
	for i := 1; i <= n; i++ {
		a[i] = nextInt()
		sum += a[i]
	}

	lines := strings.Split(strings.TrimSpace(candOutput), "\n")
	if len(lines) == 0 {
		return fmt.Errorf("empty output")
	}

	firstLine := strings.TrimSpace(lines[0])

	// Check if answer should be "No"
	if sum%2 != 0 || sum != s {
		if strings.EqualFold(firstLine, "no") {
			return nil
		}
		return fmt.Errorf("expected No, got %s", firstLine)
	}

	// If candidate says No, check if oracle also says No
	if strings.EqualFold(firstLine, "no") {
		oracleOut := oracleSolve(input)
		if strings.HasPrefix(oracleOut, "No") {
			return nil
		}
		return fmt.Errorf("candidate said No but solution exists")
	}

	if !strings.EqualFold(firstLine, "yes") {
		return fmt.Errorf("expected Yes or No, got %s", firstLine)
	}

	if len(lines) < 2 {
		return fmt.Errorf("missing edge count")
	}
	m, err := strconv.Atoi(strings.TrimSpace(lines[1]))
	if err != nil {
		return fmt.Errorf("bad edge count: %v", err)
	}
	if m != s/2 {
		return fmt.Errorf("edge count %d != s/2 = %d", m, s/2)
	}
	if len(lines) < 2+m {
		return fmt.Errorf("not enough edge lines")
	}

	deg := make([]int, n+1)
	edgeSet := make(map[[2]int]bool)
	for i := 0; i < m; i++ {
		parts := strings.Fields(lines[2+i])
		if len(parts) != 2 {
			return fmt.Errorf("bad edge line %d", i)
		}
		u, _ := strconv.Atoi(parts[0])
		v, _ := strconv.Atoi(parts[1])
		if u < 1 || u > n || v < 1 || v > n {
			return fmt.Errorf("edge (%d,%d) out of range", u, v)
		}
		if u == v {
			return fmt.Errorf("self-loop (%d,%d)", u, v)
		}
		key := [2]int{u, v}
		if u > v {
			key = [2]int{v, u}
		}
		if edgeSet[key] {
			return fmt.Errorf("duplicate edge (%d,%d)", u, v)
		}
		edgeSet[key] = true
		deg[u]++
		deg[v]++
	}

	for i := 1; i <= n; i++ {
		if deg[i] != a[i] {
			return fmt.Errorf("node %d: expected deg %d got %d", i, a[i], deg[i])
		}
	}

	return nil
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func generateCase(r *rand.Rand) string {
	n := r.Intn(10) + 2
	// Generate a valid-ish degree sequence
	mode := r.Intn(3)
	a := make([]int, n)
	sum := 0

	switch mode {
	case 0:
		// Generate a random graph and compute degrees
		edges := make([][2]int, 0)
		for i := 0; i < n; i++ {
			for j := i + 1; j < n; j++ {
				if r.Intn(3) == 0 {
					edges = append(edges, [2]int{i, j})
					a[i]++
					a[j]++
				}
			}
		}
		for _, d := range a {
			sum += d
		}
	case 1:
		// Random degrees (may be invalid)
		for i := 0; i < n; i++ {
			a[i] = r.Intn(n)
			sum += a[i]
		}
	case 2:
		// All zeros
		sum = 0
	}

	// Sort descending for Erdos-Gallai style
	sort.Sort(sort.Reverse(sort.IntSlice(a)))

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, sum))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(a[i]))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in := generateCase(r)
		out, err := runBinary(bin, in)
		if err != nil {
			fmt.Printf("Test %d runtime error: %v\nInput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if verr := validateCandidate(in, out); verr != nil {
			fmt.Printf("Test %d failed: %v\nInput:\n%sOutput:\n%s\n", i+1, verr, in, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

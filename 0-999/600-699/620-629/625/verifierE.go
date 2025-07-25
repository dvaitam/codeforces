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
)

// Node holds position and jump power
type Node struct {
	p, z, n int
}

type Item struct {
	idx  int
	z    int
	orig int
}

type PriorityQueue []Item

func (pq PriorityQueue) Len() int { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool {
	if pq[i].z != pq[j].z {
		return pq[i].z < pq[j].z
	}
	return pq[i].orig < pq[j].orig
}
func (pq PriorityQueue) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PriorityQueue) Push(x interface{}) { *pq = append(*pq, x.(Item)) }
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

func Find(x, k int, a []Node, add []int64, m int64) int {
	if x == k {
		return 0
	}
	l := int64(a[k].p - a[x].p)
	if l <= 0 {
		l += m
	}
	l += add[k] - add[x]
	t := int64(0)
	if a[x].n < a[k].n {
		l -= int64(a[x].z)
		t = 1
	}
	if l <= 0 && add[x] == 0 && add[k] == 0 {
		return 1
	}
	if a[k].z < a[x].z {
		d := int64(a[x].z - a[k].z)
		return int((l-1)/d + t + 1)
	}
	return 0
}

func solveE(n int, m int64, frogs []Node) []int {
	a := make([]Node, n)
	copy(a, frogs)
	sort.Slice(a, func(i, j int) bool { return a[i].p < a[j].p })
	pl := make([]int, n+1)
	for i := 0; i < n; i++ {
		pl[a[i].n] = i
	}
	nex := make([]int, n)
	pre := make([]int, n)
	c := make([]bool, n)
	add := make([]int64, n)
	w := make([]int, n)
	for i := 0; i < n; i++ {
		c[i] = true
	}
	pq := &PriorityQueue{}
	heap.Init(pq)
	for j := 1; j <= n; j++ {
		i := pl[j]
		k := (i + 1) % n
		nex[i] = k
		pre[k] = i
		z := Find(i, k, a, add, m)
		if z > 0 {
			w[i] = z
			heap.Push(pq, Item{idx: i, z: z, orig: a[i].n})
		}
	}
	for pq.Len() > 0 {
		item := heap.Pop(pq).(Item)
		x, z := item.idx, item.z
		if !c[x] || w[x] != z {
			continue
		}
		c[nex[x]] = false
		nex[x] = nex[nex[x]]
		num := 1
		for {
			vl := Find(x, nex[x], a, add, m)
			if vl == z {
				c[nex[x]] = false
				nex[x] = nex[nex[x]]
				num++
			} else {
				break
			}
		}
		a[x].z -= num
		if a[x].z < 0 {
			a[x].z = 0
		}
		pre[nex[x]] = x
		add[x] += int64(num) * int64(z)
		z2 := Find(x, nex[x], a, add, m)
		if z2 > 0 {
			w[x] = z2
			heap.Push(pq, Item{idx: x, z: z2, orig: a[x].n})
		} else {
			w[x] = 0
		}
		px := pre[x]
		z3 := Find(px, x, a, add, m)
		if z3 > 0 {
			w[px] = z3
			heap.Push(pq, Item{idx: px, z: z3, orig: a[px].n})
		} else {
			w[px] = 0
		}
	}
	ans := []int{}
	for i := 0; i < n; i++ {
		if c[i] {
			ans = append(ans, a[i].n)
		}
	}
	sort.Ints(ans)
	return ans
}

type testE struct {
	n     int
	m     int64
	frogs []Node
}

func genTests() []testE {
	rand.Seed(5)
	tests := make([]testE, 100)
	for i := range tests {
		n := rand.Intn(5) + 1
		m := int64(rand.Intn(50) + n)
		perm := rand.Perm(int(m))
		frogs := make([]Node, n)
		for j := 0; j < n; j++ {
			frogs[j].p = perm[j]%int(m) + 1
			frogs[j].z = rand.Intn(5) + 1
			frogs[j].n = j + 1
		}
		tests[i] = testE{n: n, m: m, frogs: frogs}
	}
	// simple small test
	tests = append(tests, testE{n: 1, m: 5, frogs: []Node{{p: 3, z: 2, n: 1}}})
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTests()
	for i, t := range tests {
		var b strings.Builder
		b.WriteString(fmt.Sprintf("%d %d\n", t.n, t.m))
		for j := 0; j < t.n; j++ {
			b.WriteString(fmt.Sprintf("%d %d\n", t.frogs[j].p, t.frogs[j].z))
		}
		input := b.String()
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("Test %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		expected := solveE(t.n, t.m, t.frogs)
		reader := strings.Fields(strings.TrimSpace(out.String()))
		if len(reader) < 1 {
			fmt.Printf("Test %d no output\n", i+1)
			os.Exit(1)
		}
		var cntOut int
		fmt.Sscanf(reader[0], "%d", &cntOut)
		if cntOut != len(expected) {
			fmt.Printf("Test %d count mismatch expected %d got %d\n", i+1, len(expected), cntOut)
			os.Exit(1)
		}
		if len(reader)-1 < cntOut {
			fmt.Printf("Test %d insufficient numbers\n", i+1)
			os.Exit(1)
		}
		for j := 0; j < cntOut; j++ {
			var val int
			fmt.Sscanf(reader[1+j], "%d", &val)
			if val != expected[j] {
				fmt.Printf("Test %d wrong answer expected %v got %v\n", i+1, expected, reader[1:1+cntOut])
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

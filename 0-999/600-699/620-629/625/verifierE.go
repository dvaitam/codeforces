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

// ---- Embedded reference solver for CF 625E (Frog Fights) ----

type refEvent struct {
	T    int64
	Type int
	i    int
	j    int
	v_i  int
	v_j  int
}

type refPQ []refEvent

func (pq refPQ) Len() int { return len(pq) }
func (pq refPQ) Less(i, j int) bool {
	if pq[i].T != pq[j].T {
		return pq[i].T < pq[j].T
	}
	return pq[i].Type < pq[j].Type
}
func (pq refPQ) Swap(i, j int) { pq[i], pq[j] = pq[j], pq[i] }
func (pq *refPQ) Push(x interface{}) {
	*pq = append(*pq, x.(refEvent))
}
func (pq *refPQ) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

type refFrog struct {
	id        int64
	p         int64
	a         int64
	C         int64
	laps      int64
	next      int
	prev      int
	active    bool
	version   int
	knockouts int64
	orig_idx  int
}

func refJumps(id int64, n int64, T int64) int64 {
	return (T - id + n) / n
}

func refFloorDiv(a, b int64) int64 {
	if a >= 0 {
		return a / b
	}
	return (a - b + 1) / b
}

func refCeilDiv(a, b int64) int64 {
	if a >= 0 {
		return (a + b - 1) / b
	}
	return a / b
}

var refFrogs []refFrog
var refN int64
var refM int64

func refCatchTime(i, j int, T_curr int64) int64 {
	if refFrogs[i].a == 0 {
		return -1
	}
	R_min := refFloorDiv(T_curr-refFrogs[i].id, refN) + 2

	delta_j := int64(0)
	if refFrogs[i].id < refFrogs[j].id {
		delta_j = -1
	}

	Delta := refFrogs[i].a - refFrogs[j].a
	K := refFrogs[j].C + refFrogs[i].laps*refM - refFrogs[i].C + delta_j*refFrogs[j].a

	var R int64
	if Delta > 0 {
		R = refCeilDiv(K, Delta)
		if R < R_min {
			R = R_min
		}
	} else {
		if R_min*Delta >= K {
			R = R_min
		} else {
			return -1
		}
	}

	return refFrogs[i].id + (R-1)*refN
}

func refUpdateState(i int, T int64) {
	J := refJumps(refFrogs[i].id, refN, T)
	X := refFrogs[i].C + J*refFrogs[i].a
	refFrogs[i].a -= refFrogs[i].knockouts
	if refFrogs[i].a < 0 {
		refFrogs[i].a = 0
	}
	refFrogs[i].C = X - J*refFrogs[i].a
	refFrogs[i].knockouts = 0
	refFrogs[i].version++
}

// solveE runs the reference solver and returns (activeCount, sorted orig_idx list)
func solveE(nFrogs int, mVal int64, positions []int64, powers []int64) (int, []int) {
	refN = int64(nFrogs)
	refM = mVal

	refFrogs = make([]refFrog, nFrogs)
	for i := 0; i < nFrogs; i++ {
		refFrogs[i] = refFrog{
			id:       int64(i + 1),
			p:        positions[i],
			a:        powers[i],
			active:   true,
			orig_idx: i + 1,
		}
	}

	sort.Slice(refFrogs, func(i, j int) bool {
		return refFrogs[i].p < refFrogs[j].p
	})

	for k := 0; k < nFrogs; k++ {
		refFrogs[k].prev = (k - 1 + nFrogs) % nFrogs
		refFrogs[k].next = (k + 1) % nFrogs
		if k == nFrogs-1 {
			refFrogs[k].laps = 1
		} else {
			refFrogs[k].laps = 0
		}
		refFrogs[k].C = refFrogs[k].p
	}

	pq := &refPQ{}
	heap.Init(pq)

	for k := 0; k < nFrogs; k++ {
		T := refCatchTime(k, refFrogs[k].next, 0)
		if T != -1 {
			heap.Push(pq, refEvent{T, 0, k, refFrogs[k].next, 0, 0})
		}
	}

	activeCount := nFrogs

	for pq.Len() > 0 {
		if activeCount <= 1 {
			break
		}
		ev := heap.Pop(pq).(refEvent)

		if ev.Type == 0 {
			i := ev.i
			j := ev.j
			if !refFrogs[i].active || !refFrogs[j].active || refFrogs[i].next != j {
				continue
			}
			if refFrogs[i].version != ev.v_i || refFrogs[j].version != ev.v_j {
				continue
			}

			refFrogs[j].active = false
			activeCount--
			if activeCount <= 1 {
				break
			}

			nxt := refFrogs[j].next
			refFrogs[i].next = nxt
			refFrogs[nxt].prev = i
			refFrogs[i].laps += refFrogs[j].laps
			refFrogs[i].knockouts++

			T_new := refCatchTime(i, nxt, ev.T-1)
			if T_new != -1 {
				heap.Push(pq, refEvent{T_new, 0, i, nxt, refFrogs[i].version, refFrogs[nxt].version})
			}

			heap.Push(pq, refEvent{ev.T, 1, i, 0, refFrogs[i].version, 0})
		} else {
			i := ev.i
			if !refFrogs[i].active || refFrogs[i].version != ev.v_i {
				continue
			}
			if refFrogs[i].knockouts > 0 {
				refUpdateState(i, ev.T)

				prv := refFrogs[i].prev
				T_prv := refCatchTime(prv, i, ev.T)
				if T_prv != -1 {
					heap.Push(pq, refEvent{T_prv, 0, prv, i, refFrogs[prv].version, refFrogs[i].version})
				}

				nxt := refFrogs[i].next
				T_nxt := refCatchTime(i, nxt, ev.T)
				if T_nxt != -1 {
					heap.Push(pq, refEvent{T_nxt, 0, i, nxt, refFrogs[i].version, refFrogs[nxt].version})
				}
			}
		}
	}

	var result []int
	for _, f := range refFrogs {
		if f.active {
			result = append(result, f.orig_idx)
		}
	}
	sort.Ints(result)
	return activeCount, result
}

// ---- Test generation and verification ----

type testE struct {
	n     int
	m     int64
	pos   []int64
	pow   []int64
}

func genTests() []testE {
	rand.Seed(5)
	tests := make([]testE, 100)
	for i := range tests {
		n := rand.Intn(5) + 1
		m := int64(rand.Intn(50) + n)
		perm := rand.Perm(int(m))
		pos := make([]int64, n)
		pow := make([]int64, n)
		for j := 0; j < n; j++ {
			pos[j] = int64(perm[j]%int(m) + 1)
			pow[j] = int64(rand.Intn(5) + 1)
		}
		tests[i] = testE{n: n, m: m, pos: pos, pow: pow}
	}
	// simple small test
	tests = append(tests, testE{n: 1, m: 5, pos: []int64{3}, pow: []int64{2}})
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
			b.WriteString(fmt.Sprintf("%d %d\n", t.pos[j], t.pow[j]))
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
		expCount, expected := solveE(t.n, t.m, t.pos, t.pow)
		reader := strings.Fields(strings.TrimSpace(out.String()))
		if len(reader) < 1 {
			fmt.Printf("Test %d no output\n", i+1)
			os.Exit(1)
		}
		var cntOut int
		fmt.Sscanf(reader[0], "%d", &cntOut)
		if cntOut != expCount {
			fmt.Printf("Test %d count mismatch expected %d got %d\n", i+1, expCount, cntOut)
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

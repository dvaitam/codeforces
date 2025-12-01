package main

import (
	"bytes"
	"container/heap"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type testCaseD struct {
	s string
}

type state struct {
	cost int
	str  string
}

func generateCaseD(rng *rand.Rand) (string, testCaseD) {
	n := rng.Intn(7) + 1 // length 1..7
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			b[i] = '0'
		} else {
			b[i] = '1'
		}
	}
	s := string(b)
	input := fmt.Sprintf("1\n%s\n", s)
	return input, testCaseD{s: s}
}

func isSorted(s string) bool {
	for i := 1; i < len(s); i++ {
		if s[i-1] > s[i] {
			return false
		}
	}
	return true
}

func minCost(s string) int64 {
	const base = 1_000_000
	pq := &statePQ{}
	heap.Push(pq, state{0, s})
	dist := map[string]int{s: 0}
	for pq.Len() > 0 {
		cur := heap.Pop(pq).(state)
		if cur.cost != dist[cur.str] {
			continue
		}
		if isSorted(cur.str) {
			ops := cur.cost / base
			rem := cur.cost % base
			return int64(ops)*1_000_000_000_000 + int64(rem)
		}
		// swaps
		for i := 0; i+1 < len(cur.str); i++ {
			t := []byte(cur.str)
			t[i], t[i+1] = t[i+1], t[i]
			ns := string(t)
			nc := cur.cost + base
			if old, ok := dist[ns]; !ok || nc < old {
				dist[ns] = nc
				heap.Push(pq, state{nc, ns})
			}
		}
		// removals
		for i := 0; i < len(cur.str); i++ {
			ns := cur.str[:i] + cur.str[i+1:]
			nc := cur.cost + base + 1
			if old, ok := dist[ns]; !ok || nc < old {
				dist[ns] = nc
				heap.Push(pq, state{nc, ns})
			}
		}
	}
	return 0
}

type statePQ []state

func (p statePQ) Len() int            { return len(p) }
func (p statePQ) Less(i, j int) bool  { return p[i].cost < p[j].cost }
func (p statePQ) Swap(i, j int)       { p[i], p[j] = p[j], p[i] }
func (p *statePQ) Push(x interface{}) { *p = append(*p, x.(state)) }
func (p *statePQ) Pop() interface{} {
	old := *p
	v := old[len(old)-1]
	*p = old[:len(old)-1]
	return v
}

func expectedD(tc testCaseD) string {
	res := minCost(tc.s)
	return fmt.Sprintf("%d", res)
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

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= 100; i++ {
		input, tc := generateCaseD(rng)
		expect := expectedD(tc)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i, expect, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

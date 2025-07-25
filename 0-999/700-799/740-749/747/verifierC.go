package main

import (
	"bufio"
	"bytes"
	"container/heap"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type task struct {
	t int
	k int
	d int
}

type event struct {
	release int
	servers []int
}

type pq []*event

func (p pq) Len() int            { return len(p) }
func (p pq) Less(i, j int) bool  { return p[i].release < p[j].release }
func (p pq) Swap(i, j int)       { p[i], p[j] = p[j], p[i] }
func (p *pq) Push(x interface{}) { *p = append(*p, x.(*event)) }
func (p *pq) Pop() interface{}   { old := *p; n := len(old); x := old[n-1]; *p = old[:n-1]; return x }

func expected(n int, tasks []task) []int {
	free := make([]bool, n+1)
	for i := 1; i <= n; i++ {
		free[i] = true
	}
	var queue pq
	heap.Init(&queue)
	res := make([]int, len(tasks))
	for idx, tk := range tasks {
		for queue.Len() > 0 && queue[0].release <= tk.t {
			ev := heap.Pop(&queue).(*event)
			for _, id := range ev.servers {
				free[id] = true
			}
		}
		picked := make([]int, 0, tk.k)
		for id := 1; id <= n && len(picked) < tk.k; id++ {
			if free[id] {
				picked = append(picked, id)
			}
		}
		if len(picked) < tk.k {
			res[idx] = -1
			continue
		}
		sum := 0
		for _, id := range picked {
			free[id] = false
			sum += id
		}
		heap.Push(&queue, &event{release: tk.t + tk.d, servers: picked})
		res[idx] = sum
	}
	return res
}

func runCase(bin string, n int, ts []task) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, len(ts)))
	for _, t := range ts {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", t.t, t.k, t.d))
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("runtime error: %v", err)
	}
	scanner := bufio.NewScanner(bytes.NewReader(out))
	res := expected(n, ts)
	for i, exp := range res {
		if !scanner.Scan() {
			return fmt.Errorf("missing output line %d", i+1)
		}
		got, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
		if err != nil {
			return fmt.Errorf("bad integer on line %d", i+1)
		}
		if got != exp {
			return fmt.Errorf("line %d expected %d got %d", i+1, exp, got)
		}
	}
	if scanner.Scan() {
		return fmt.Errorf("extra output")
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rand.Seed(1)
	for tcase := 1; tcase <= 120; tcase++ {
		n := rand.Intn(5) + 1
		q := rand.Intn(5) + 1
		tasks := make([]task, q)
		cur := rand.Intn(3)
		for i := 0; i < q; i++ {
			cur += rand.Intn(3) + 1
			tasks[i] = task{t: cur, k: rand.Intn(n) + 1, d: rand.Intn(5) + 1}
		}
		if err := runCase(bin, n, tasks); err != nil {
			fmt.Printf("Test %d failed: %v\n", tcase, err)
			return
		}
	}
	fmt.Println("OK")
}

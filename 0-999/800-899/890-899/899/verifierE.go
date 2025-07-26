package main

import (
	"container/heap"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type test struct {
	input    string
	expected string
}

type Segment struct {
	length int
	idx    int
	val    int
	left   int
	right  int
	alive  bool
}

var segments []Segment

type PQ []int

func (pq PQ) Len() int { return len(pq) }
func (pq PQ) Less(i, j int) bool {
	a := segments[pq[i]]
	b := segments[pq[j]]
	if a.length != b.length {
		return a.length > b.length
	}
	return a.idx < b.idx
}
func (pq PQ) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PQ) Push(x interface{}) { *pq = append(*pq, x.(int)) }
func (pq *PQ) Pop() interface{} {
	old := *pq
	n := len(old)
	x := old[n-1]
	*pq = old[:n-1]
	return x
}

func solve(input string) string {
	reader := strings.NewReader(input)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return ""
	}
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &arr[i])
	}
	segments = make([]Segment, 0)
	var pq PQ
	start := 0
	for start < n {
		j := start
		for j < n && arr[j] == arr[start] {
			j++
		}
		segID := len(segments)
		segments = append(segments, Segment{
			length: j - start,
			idx:    start,
			val:    arr[start],
			left:   segID - 1,
			right:  -1,
			alive:  true,
		})
		if segID > 0 {
			segments[segID-1].right = segID
		}
		heap.Push(&pq, segID)
		start = j
	}
	ops := 0
	for pq.Len() > 0 {
		id := heap.Pop(&pq).(int)
		if !segments[id].alive {
			continue
		}
		ops++
		left := segments[id].left
		right := segments[id].right
		segments[id].alive = false
		if left != -1 {
			segments[left].right = right
		}
		if right != -1 {
			segments[right].left = left
		}
		if left != -1 && right != -1 && segments[left].alive && segments[right].alive && segments[left].val == segments[right].val {
			newID := len(segments)
			segments[left].alive = false
			segments[right].alive = false
			segments = append(segments, Segment{
				length: segments[left].length + segments[right].length,
				idx:    segments[left].idx,
				val:    segments[left].val,
				left:   segments[left].left,
				right:  segments[right].right,
				alive:  true,
			})
			if segments[newID].left != -1 {
				segments[segments[newID].left].right = newID
			}
			if segments[newID].right != -1 {
				segments[segments[newID].right].left = newID
			}
			heap.Push(&pq, newID)
		}
	}
	return fmt.Sprintf("%d", ops)
}

func generateTests() []test {
	rand.Seed(8995)
	var tests []test
	fixed := []string{"3\n1 2 3\n", "5\n7 7 7 2 2\n"}
	for _, f := range fixed {
		tests = append(tests, test{f, solve(f)})
	}
	for len(tests) < 100 {
		n := rand.Intn(30) + 1
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", n)
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", rand.Intn(5)))
		}
		sb.WriteByte('\n')
		inp := sb.String()
		tests = append(tests, test{inp, solve(inp)})
	}
	return tests
}

func runBinary(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.CommandContext(ctx, "go", "run", bin)
	} else {
		cmd = exec.CommandContext(ctx, bin)
	}
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("time limit")
	}
	if err != nil {
		return "", fmt.Errorf("%v: %s", err, out)
	}
	return strings.TrimSpace(string(out)), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		out, err := runBinary(bin, t.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if out != strings.TrimSpace(t.expected) {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d\ninput:\n%sexpected:%s\n got:%s\n", i+1, t.input, t.expected, out)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

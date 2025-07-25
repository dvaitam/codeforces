package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type Flight struct {
	cost int64
	idx  int
}

type PQ []Flight

func (pq PQ) Len() int            { return len(pq) }
func (pq PQ) Less(i, j int) bool  { return pq[i].cost > pq[j].cost }
func (pq PQ) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PQ) Push(x interface{}) { *pq = append(*pq, x.(Flight)) }
func (pq *PQ) Pop() interface{} {
	old := *pq
	n := len(old)
	x := old[n-1]
	*pq = old[:n-1]
	return x
}

func solveCase(line string) string {
	fields := strings.Fields(line)
	if len(fields) < 2 {
		return ""
	}
	idx := 0
	n, _ := strconv.Atoi(fields[idx])
	idx++
	k, _ := strconv.Atoi(fields[idx])
	idx++
	costs := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		if idx >= len(fields) {
			return ""
		}
		v, _ := strconv.ParseInt(fields[idx], 10, 64)
		idx++
		costs[i] = v
	}
	pq := &PQ{}
	heap.Init(pq)
	ans := make([]int, n+1)
	var total int64
	cur := 1
	for t := k + 1; t <= k+n; t++ {
		for cur <= n && cur <= t {
			heap.Push(pq, Flight{cost: costs[cur], idx: cur})
			cur++
		}
		if pq.Len() == 0 {
			continue
		}
		f := heap.Pop(pq).(Flight)
		ans[f.idx] = t
		total += int64(t-f.idx) * f.cost
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", total))
	for i := 1; i <= n; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(ans[i]))
	}
	return sb.String()
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("runtime error: %v", err)
	}
	return strings.TrimSpace(string(out)), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesA.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		expected := solveCase(line)
		got, err := run(bin, line+"\n")
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx, expected, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}

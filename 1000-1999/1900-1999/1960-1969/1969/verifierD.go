package main

import (
	"bufio"
	"bytes"
	"container/heap"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

type item struct{ a, b int64 }

type maxHeap []int64

func (h maxHeap) Len() int            { return len(h) }
func (h maxHeap) Less(i, j int) bool  { return h[i] > h[j] }
func (h maxHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *maxHeap) Push(x interface{}) { *h = append(*h, x.(int64)) }
func (h *maxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func runCandidate(bin, input string) (string, error) {
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

func expected(n, k int, a, b []int64) int64 {
	items := make([]item, n)
	for i := 0; i < n; i++ {
		items[i] = item{a[i], b[i]}
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].b == items[j].b {
			return items[i].a < items[j].a
		}
		return items[i].b > items[j].b
	})
	suff := make([]int64, n+1)
	for i := n - 1; i >= 0; i-- {
		suff[i] = suff[i+1]
		diff := items[i].b - items[i].a
		if diff > 0 {
			suff[i] += diff
		}
	}
	h := &maxHeap{}
	sum := int64(0)
	ans := int64(0)
	for i := 0; i < n; i++ {
		if len(*h) >= k {
			profit := suff[i] - sum
			if profit > ans {
				ans = profit
			}
		}
		heap.Push(h, items[i].a)
		sum += items[i].a
		if len(*h) > k {
			sum -= heap.Pop(h).(int64)
		}
	}
	return ans
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesD.txt")
	if err != nil {
		fmt.Println("failed to open testcasesD.txt:", err)
		os.Exit(1)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		fields := strings.Fields(line)
		n, _ := strconv.Atoi(fields[0])
		k, _ := strconv.Atoi(fields[1])
		a := make([]int64, n)
		b := make([]int64, n)
		for i := 0; i < n; i++ {
			v, _ := strconv.ParseInt(fields[2+i], 10, 64)
			a[i] = v
		}
		for i := 0; i < n; i++ {
			v, _ := strconv.ParseInt(fields[2+n+i], 10, 64)
			b[i] = v
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "1\n%d %d\n", n, k)
		for i, v := range a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')
		for i, v := range b {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')
		want := strconv.FormatInt(expected(n, k, a, b), 10)
		got, err := runCandidate(bin, sb.String())
		if err != nil {
			fmt.Printf("test %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Printf("test %d failed: expected %s got %s\ninput:\n%s", idx, want, strings.TrimSpace(got), sb.String())
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}

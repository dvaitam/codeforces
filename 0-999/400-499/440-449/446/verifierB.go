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

type MaxHeap []int64

func (h MaxHeap) Len() int            { return len(h) }
func (h MaxHeap) Less(i, j int) bool  { return h[i] > h[j] }
func (h MaxHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *MaxHeap) Push(x interface{}) { *h = append(*h, x.(int64)) }
func (h *MaxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func expected(n, m, k int, p int64, mat [][]int64) int64 {
	rowSum := make([]int64, n)
	colSum := make([]int64, m)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			v := mat[i][j]
			rowSum[i] += v
			colSum[j] += v
		}
	}
	bestR := make([]int64, k+1)
	bestC := make([]int64, k+1)
	hR := make(MaxHeap, len(rowSum))
	copy(hR, rowSum)
	heap.Init(&hR)
	for i := 1; i <= k; i++ {
		x := heap.Pop(&hR).(int64)
		bestR[i] = bestR[i-1] + x
		heap.Push(&hR, x-int64(m)*p)
	}
	hC := make(MaxHeap, len(colSum))
	copy(hC, colSum)
	heap.Init(&hC)
	for i := 1; i <= k; i++ {
		x := heap.Pop(&hC).(int64)
		bestC[i] = bestC[i-1] + x
		heap.Push(&hC, x-int64(n)*p)
	}
	ans := int64(-1 << 60)
	for i := 0; i <= k; i++ {
		j := k - i
		val := bestR[i] + bestC[j] - int64(i)*int64(j)*p
		if val > ans {
			ans = val
		}
	}
	return ans
}

func generateCase(rng *rand.Rand) (string, int64) {
	n := rng.Intn(3) + 1
	m := rng.Intn(3) + 1
	k := rng.Intn(3) + 1
	p := int64(rng.Intn(3) + 1)
	mat := make([][]int64, n)
	for i := 0; i < n; i++ {
		mat[i] = make([]int64, m)
		for j := 0; j < m; j++ {
			mat[i][j] = int64(rng.Intn(5))
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d %d\n", n, m, k, p))
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(mat[i][j]))
		}
		sb.WriteByte('\n')
	}
	return sb.String(), expected(n, m, k, p, mat)
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		var got int64
		if _, err := fmt.Sscan(out, &got); err != nil {
			fmt.Fprintf(os.Stderr, "case %d: cannot parse output: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d: expected %d got %d\ninput:\n%s", i+1, exp, got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

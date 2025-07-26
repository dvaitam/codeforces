package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// BIT for prefix sums
type BIT struct {
	n    int
	tree []int64
}

func NewBIT(n int) *BIT {
	return &BIT{n: n, tree: make([]int64, n+1)}
}

func (b *BIT) Add(i int, v int64) {
	for x := i; x <= b.n; x += x & -x {
		b.tree[x] += v
	}
}

func (b *BIT) Sum(i int) int64 {
	var s int64
	for x := i; x > 0; x -= x & -x {
		s += b.tree[x]
	}
	return s
}

func (b *BIT) RangeSum(l, r int) int64 {
	if l > r {
		return 0
	}
	return b.Sum(r) - b.Sum(l-1)
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

func solveCase(pairs [][2]int) []int64 {
	n := len(pairs)
	bitCnt := NewBIT(n)
	bitSum := NewBIT(n)
	var totalB, S, sumAb int64
	low := int64(1<<63 - 1)
	res := make([]int64, n)
	for i := 0; i < n; i++ {
		a := int64(pairs[i][0])
		b := pairs[i][1]
		if a < low {
			low = a
		}
		cntGt := bitCnt.RangeSum(b+1, n)
		sumGt := bitSum.RangeSum(b+1, n)
		sumLower := totalB - sumGt
		S += int64(b)*cntGt + sumLower
		totalB += int64(b)
		bitCnt.Add(b, 1)
		bitSum.Add(b, int64(b))
		sumAb += int64(b) * a
		res[i] = low*totalB + S - sumAb
	}
	return res
}

func genCase(rng *rand.Rand) [][2]int {
	n := rng.Intn(6) + 1
	used := make(map[int]bool)
	pairs := make([][2]int, n)
	for i := 0; i < n; i++ {
		a := rng.Intn(20) + 1
		var b int
		for {
			b = rng.Intn(n*2) + 1
			if !used[b] {
				used[b] = true
				break
			}
		}
		pairs[i] = [2]int{a, b}
	}
	return pairs
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		pairs := genCase(rng)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", len(pairs)))
		for _, p := range pairs {
			sb.WriteString(fmt.Sprintf("%d %d\n", p[0], p[1]))
		}
		expect := solveCase(pairs)
		out, err := runCandidate(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		fields := strings.Fields(out)
		if len(fields) != len(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d numbers got %d\n", i+1, len(expect), len(fields))
			os.Exit(1)
		}
		for j, f := range fields {
			val, err := strconv.ParseInt(f, 10, 64)
			if err != nil || val != expect[j] {
				fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %s at index %d\n", i+1, expect[j], f, j)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}

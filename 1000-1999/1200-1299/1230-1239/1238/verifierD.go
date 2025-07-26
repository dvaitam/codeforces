package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Fenwick struct {
	n   int
	bit []int
}

func NewFenwick(n int) *Fenwick { return &Fenwick{n, make([]int, n+2)} }
func (f *Fenwick) Add(idx, val int) {
	for idx <= f.n {
		f.bit[idx] += val
		idx += idx & -idx
	}
}
func (f *Fenwick) Sum(idx int) int {
	res := 0
	for idx > 0 {
		res += f.bit[idx]
		idx -= idx & -idx
	}
	return res
}

func countGood(s string) int64 {
	n := len(s)
	next := make([]int, n)
	prev := make([]int, n)
	last := map[byte]int{'A': -1, 'B': -1}
	for i := 0; i < n; i++ {
		prev[i] = last[s[i]]
		last[s[i]] = i
	}
	last['A'] = n
	last['B'] = n
	for i := n - 1; i >= 0; i-- {
		next[i] = last[s[i]]
		last[s[i]] = i
	}
	buckets := make([][]int, n+1)
	for l := 0; l < n; l++ {
		if next[l] < n {
			buckets[next[l]] = append(buckets[next[l]], l+1)
		}
	}
	ft := NewFenwick(n)
	var ans int64
	for r := 0; r < n; r++ {
		for _, idx := range buckets[r] {
			ft.Add(idx, 1)
		}
		if prev[r] >= 0 {
			ans += int64(ft.Sum(prev[r] + 1))
		}
	}
	return ans
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
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func generateTests() []string {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]string, 0, 100)
	for i := 0; i < 100; i++ {
		n := rng.Intn(15) + 1
		b := make([]byte, n)
		for j := 0; j < n; j++ {
			if rng.Intn(2) == 0 {
				b[j] = 'A'
			} else {
				b[j] = 'B'
			}
		}
		tests = append(tests, string(b))
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, s := range tests {
		input := fmt.Sprintf("%d %s\n", len(s), s)
		want := fmt.Sprintf("%d", countGood(s))
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Printf("case %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Printf("case %d failed: expected %s got %s\ninput:%s", i+1, want, got, input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

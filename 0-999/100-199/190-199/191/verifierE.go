package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

func runCandidate(bin string, input string) (string, error) {
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

// Fenwick tree
type BIT struct {
	n    int
	tree []int
}

func newBIT(n int) *BIT {
	return &BIT{n: n, tree: make([]int, n+1)}
}

func (b *BIT) update(i, v int) {
	for ; i <= b.n; i += i & -i {
		b.tree[i] += v
	}
}

func (b *BIT) query(i int) int {
	s := 0
	for ; i > 0; i &= i - 1 {
		s += b.tree[i]
	}
	return s
}

func unique64(a []int64) []int64 {
	if len(a) == 0 {
		return a
	}
	j := 0
	for i := 1; i < len(a); i++ {
		if a[i] != a[j] {
			j++
			a[j] = a[i]
		}
	}
	return a[:j+1]
}

func solve(n int, k int64, arr []int64) string {
	S := make([]int64, n+1)
	for i := 0; i < n; i++ {
		S[i+1] = S[i] + arr[i]
	}
	P := append([]int64(nil), S...)
	sort.Slice(P, func(i, j int) bool { return P[i] < P[j] })
	P = unique64(P)
	count := func(x int64) int64 {
		bit := newBIT(len(P))
		var cnt int64
		for _, sj := range S {
			v := sj - x
			idx := sort.Search(len(P), func(i int) bool { return P[i] > v })
			cnt += int64(bit.query(idx))
			pos := sort.Search(len(P), func(i int) bool { return P[i] >= sj })
			bit.update(pos+1, 1)
		}
		return cnt
	}
	low, high := int64(-1e14), int64(1e14)
	for low < high {
		mid := (low + high + 1) >> 1
		if count(mid) >= k {
			low = mid
		} else {
			high = mid - 1
		}
	}
	return fmt.Sprintf("%d", low)
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(8) + 1
	maxK := int64(n*(n+1)) / 2
	k := rng.Int63n(maxK) + 1
	arr := make([]int64, n)
	for i := range arr {
		arr[i] = rng.Int63n(101) - 50
	}
	var in strings.Builder
	fmt.Fprintf(&in, "%d %d\n", n, k)
	for i, v := range arr {
		if i+1 == n {
			fmt.Fprintf(&in, "%d\n", v)
		} else {
			fmt.Fprintf(&in, "%d ", v)
		}
	}
	exp := solve(n, k, arr)
	return in.String(), exp
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 100; i++ {
		input, expect := generateCase(rng)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

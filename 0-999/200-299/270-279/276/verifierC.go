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

type testCaseC struct {
	n int
	q int
	a []int
	L []int
	R []int
}

func generateCaseC(rng *rand.Rand) testCaseC {
	n := rng.Intn(10) + 1
	q := rng.Intn(10) + 1
	a := make([]int, n)
	for i := range a {
		a[i] = rng.Intn(50)
	}
	L := make([]int, q)
	R := make([]int, q)
	for i := 0; i < q; i++ {
		l := rng.Intn(n) + 1
		r := rng.Intn(n-l+1) + l
		L[i] = l
		R[i] = r
	}
	return testCaseC{n: n, q: q, a: a, L: L, R: R}
}

func buildInputC(tc testCaseC) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.q)
	for i := 0; i < tc.n; i++ {
		fmt.Fprintf(&sb, "%d ", tc.a[i])
	}
	sb.WriteByte('\n')
	for i := 0; i < tc.q; i++ {
		fmt.Fprintf(&sb, "%d %d\n", tc.L[i], tc.R[i])
	}
	return sb.String()
}

func expectedC(tc testCaseC) int64 {
	freqDiff := make([]int, tc.n+1)
	for i := 0; i < tc.q; i++ {
		freqDiff[tc.L[i]-1]++
		freqDiff[tc.R[i]]--
	}
	freq := make([]int, tc.n)
	cur := 0
	for i := 0; i < tc.n; i++ {
		cur += freqDiff[i]
		freq[i] = cur
	}
	arr := append([]int(nil), tc.a...)
	sort.Ints(arr)
	sort.Ints(freq)
	var ans int64
	for i := 0; i < tc.n; i++ {
		ans += int64(arr[i]) * int64(freq[i])
	}
	return ans
}

func runCaseC(bin string, tc testCaseC) error {
	input := buildInputC(tc)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int64
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	want := expectedC(tc)
	if got != want {
		return fmt.Errorf("expected %d got %d", want, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCaseC(rng)
		if err := runCaseC(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, buildInputC(tc))
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

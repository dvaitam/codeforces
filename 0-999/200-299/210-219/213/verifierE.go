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

type caseE struct {
	n, m int
	a    []int
	b    []int
}

func genPerm(rng *rand.Rand, n int) []int {
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = i + 1
	}
	rng.Shuffle(n, func(i, j int) { arr[i], arr[j] = arr[j], arr[i] })
	return arr
}

func genCaseE(rng *rand.Rand) caseE {
	n := rng.Intn(5) + 1
	m := rng.Intn(10) + n
	a := genPerm(rng, n)
	b := genPerm(rng, m)
	return caseE{n, m, a, b}
}

// solveE counts distinct d in [0, m-n] such that (a[0]+d, ..., a[n-1]+d)
// is a subsequence of b. Brute-force is fine for the small cases generated here.
func solveE(tc caseE) int {
	n := tc.n
	m := tc.m
	a := tc.a
	b := tc.b
	// pos[v] = 1-based position of value v in b
	pos := make([]int, m+1)
	for i, v := range b {
		pos[v] = i + 1
	}
	count := 0
	for d := 0; d <= m-n; d++ {
		// verify a[i]+d in [1,m] and positions are strictly increasing
		ok := true
		prev := 0
		for i := 0; i < n; i++ {
			v := a[i] + d
			if v < 1 || v > m {
				ok = false
				break
			}
			if pos[v] <= prev {
				ok = false
				break
			}
			prev = pos[v]
		}
		if ok {
			count++
		}
	}
	return count
}

func runE(bin string, tc caseE) error {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.m)
	for i, v := range tc.a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprint(&sb, v)
	}
	sb.WriteByte('\n')
	for i, v := range tc.b {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprint(&sb, v)
	}
	sb.WriteByte('\n')
	input := sb.String()
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	exp := solveE(tc)
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		tc := genCaseE(rng)
		if err := runE(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

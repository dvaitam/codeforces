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

func serialize(arr []int) string {
	var sb strings.Builder
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	return sb.String()
}

func lexLess(a, b []int) bool {
	for i := range a {
		if a[i] < b[i] {
			return true
		}
		if a[i] > b[i] {
			return false
		}
	}
	return false
}

func bfs(a []int) []int {
	start := append([]int{}, a...)
	best := append([]int{}, a...)
	type node struct{ arr []int }
	q := []node{{start}}
	seen := map[string]bool{serialize(start): true}
	for len(q) > 0 {
		cur := q[0].arr
		q = q[1:]
		if lexLess(cur, best) {
			best = append([]int{}, cur...)
		}
		n := len(cur)
		for i := 0; i < n; i++ {
			for j := i + 1; j < n; j++ {
				if cur[i]^cur[j] < 4 {
					next := append([]int{}, cur...)
					next[i], next[j] = next[j], next[i]
					key := serialize(next)
					if !seen[key] {
						seen[key] = true
						q = append(q, node{next})
					}
				}
			}
		}
	}
	return best
}

type caseG struct{ arr []int }

func genCase(rng *rand.Rand) caseG {
	n := rng.Intn(5) + 2
	arr := make([]int, n)
	for i := range arr {
		arr[i] = rng.Intn(20)
	}
	return caseG{arr}
}

func runCase(bin string, tc caseG) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("1\n%d\n", len(tc.arr)))
	for i, v := range tc.arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
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
	fields := strings.Fields(strings.TrimSpace(out.String()))
	if len(fields) != len(tc.arr) {
		return fmt.Errorf("expected %d numbers got %d", len(tc.arr), len(fields))
	}
	got := make([]int, len(fields))
	for i, f := range fields {
		fmt.Sscan(f, &got[i])
	}
	exp := bfs(tc.arr)
	for i := range exp {
		if got[i] != exp[i] {
			return fmt.Errorf("mismatch at %d exp %d got %d", i, exp[i], got[i])
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		tc := genCase(rng)
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

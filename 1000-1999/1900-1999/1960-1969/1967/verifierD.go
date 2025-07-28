package main

import (
	"bytes"
	"container/list"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type TestCaseD struct {
	n int
	m int
	a []int
	b []int
}

func genCaseD(rng *rand.Rand) TestCaseD {
	n := rng.Intn(3) + 1
	m := rng.Intn(3) + 1
	a := make([]int, n)
	for i := range a {
		a[i] = rng.Intn(m) + 1
	}
	b := make([]int, m)
	for i := range b {
		b[i] = rng.Intn(m) + 1
	}
	return TestCaseD{n: n, m: m, a: a, b: b}
}

func isNonDecreasing(arr []int) bool {
	for i := 1; i < len(arr); i++ {
		if arr[i] < arr[i-1] {
			return false
		}
	}
	return true
}

func stateKey(arr []int) string {
	var sb strings.Builder
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(',')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	return sb.String()
}

func expectedD(tc TestCaseD) int {
	start := append([]int(nil), tc.a...)
	if isNonDecreasing(start) {
		return 0
	}
	visited := map[string]bool{stateKey(start): true}
	type item struct {
		arr  []int
		step int
	}
	q := list.New()
	q.PushBack(item{start, 0})
	for q.Len() > 0 {
		front := q.Remove(q.Front()).(item)
		if front.step > 8 {
			continue
		}
		for mask := 1; mask < (1 << tc.n); mask++ {
			next := append([]int(nil), front.arr...)
			for i := 0; i < tc.n; i++ {
				if mask>>i&1 == 1 {
					next[i] = tc.b[next[i]-1]
				}
			}
			key := stateKey(next)
			if visited[key] {
				continue
			}
			if isNonDecreasing(next) {
				return front.step + 1
			}
			visited[key] = true
			q.PushBack(item{next, front.step + 1})
		}
	}
	return -1
}

func runCaseD(bin string, tc TestCaseD, expect int) error {
	var sb strings.Builder
	sb.WriteString("1\n")
	fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.m)
	for i, v := range tc.a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	for i, v := range tc.b {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')

	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("invalid output: %v", err)
	}
	if got != expect {
		return fmt.Errorf("expected %d got %d", expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := genCaseD(rng)
		exp := expectedD(tc)
		if err := runCaseD(bin, tc, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput: n=%d m=%d a=%v b=%v\n", i+1, err, tc.n, tc.m, tc.a, tc.b)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

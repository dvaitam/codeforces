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

type TestCaseA struct {
	n int
	k int64
	a []int64
}

func genCaseA(rng *rand.Rand) TestCaseA {
	n := rng.Intn(5) + 1
	k := rng.Int63n(50)
	a := make([]int64, n)
	for i := range a {
		a[i] = rng.Int63n(50) + 1
	}
	return TestCaseA{n: n, k: k, a: a}
}

func canMake(t int64, a []int64, k int64) bool {
	n := int64(len(a))
	q := t / n
	r := int(t % n)
	need := int64(0)
	for i, v := range a {
		req := q
		if i < r {
			req++
		}
		if v < req {
			need += req - v
			if need > k {
				return false
			}
		}
	}
	return need <= k
}

func expectedA(tc TestCaseA) int64 {
	a := append([]int64(nil), tc.a...)
	sort.Slice(a, func(i, j int) bool { return a[i] > a[j] })
	sum := int64(0)
	for _, v := range a {
		sum += v
	}
	low, high := int64(0), sum+tc.k
	ans := int64(0)
	for low <= high {
		mid := (low + high) / 2
		if canMake(mid, a, tc.k) {
			ans = mid
			low = mid + 1
		} else {
			high = mid - 1
		}
	}
	if ans >= int64(tc.n) {
		return ans - int64(tc.n) + 1
	}
	return 0
}

func runCaseA(bin string, tc TestCaseA, expect int64) error {
	var sb strings.Builder
	sb.WriteString("1\n")
	fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.k)
	for i, v := range tc.a {
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
	var got int64
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := genCaseA(rng)
		exp := expectedA(tc)
		if err := runCaseA(bin, tc, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput: n=%d k=%d a=%v\n", i+1, err, tc.n, tc.k, tc.a)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

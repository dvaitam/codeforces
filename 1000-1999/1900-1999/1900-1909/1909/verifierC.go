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

func solve(l, r, c []int) int64 {
	n := len(l)
	sort.Ints(l)
	sort.Ints(r)
	sort.Slice(c, func(i, j int) bool { return c[i] > c[j] })
	diff := make([]int, n)
	pool := make([]int, 0, n)
	idx := 0
	for i, rv := range r {
		for idx < n && l[idx] < rv {
			pool = append(pool, l[idx])
			idx++
		}
		lv := pool[len(pool)-1]
		pool = pool[:len(pool)-1]
		diff[i] = rv - lv
	}
	sort.Ints(diff)
	var ans int64
	for i := 0; i < n; i++ {
		ans += int64(diff[i]) * int64(c[i])
	}
	return ans
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(4) + 1
	l := make([]int, n)
	r := make([]int, n)
	c := make([]int, n)
	used := map[int]bool{}
	for i := 0; i < n; i++ {
		for {
			lv := rng.Intn(100)
			rv := lv + rng.Intn(20) + 1
			if !used[lv] && !used[rv] {
				used[lv] = true
				used[rv] = true
				l[i] = lv
				r[i] = rv
				break
			}
		}
		c[i] = rng.Intn(100) + 1
	}
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		sb.WriteString(fmt.Sprintf("%d ", l[i]))
	}
	sb.WriteString("\n")
	for i := 0; i < n; i++ {
		sb.WriteString(fmt.Sprintf("%d ", r[i]))
	}
	sb.WriteString("\n")
	for i := 0; i < n; i++ {
		sb.WriteString(fmt.Sprintf("%d ", c[i]))
	}
	sb.WriteString("\n")
	expect := fmt.Sprintf("%d", solve(append([]int(nil), l...), append([]int(nil), r...), append([]int(nil), c...)))
	return sb.String(), expect
}

func runCase(bin, input, expect string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expect {
		return fmt.Errorf("expected %s got %s", expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, expect := generateCase(rng)
		if err := runCase(bin, in, expect); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

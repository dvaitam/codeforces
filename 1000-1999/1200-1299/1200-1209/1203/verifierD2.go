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

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func solve(s, t string) string {
	n := len(s)
	m := len(t)
	pre := make([]int, m)
	idx := 0
	for i := 0; i < n && idx < m; i++ {
		if s[i] == t[idx] {
			pre[idx] = i
			idx++
		}
	}
	suf := make([]int, m)
	idx = m - 1
	for i := n - 1; i >= 0 && idx >= 0; i-- {
		if s[i] == t[idx] {
			suf[idx] = i
			idx--
		}
	}
	ans := max(suf[0], n-1-pre[m-1])
	for i := 0; i < m-1; i++ {
		if gap := suf[i+1] - pre[i] - 1; gap > ans {
			ans = gap
		}
	}
	return fmt.Sprintf("%d", ans)
}

func randomString(rng *rand.Rand, n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = byte('a' + rng.Intn(26))
	}
	return string(b)
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(20) + 1
	m := rng.Intn(n) + 1
	s := randomString(rng, n)
	idxs := rng.Perm(n)[:m]
	sort.Ints(idxs)
	tb := make([]byte, m)
	for i, idx := range idxs {
		tb[i] = s[idx]
	}
	t := string(tb)
	input := fmt.Sprintf("%s\n%s\n", s, t)
	expect := solve(s, t)
	return input, expect
}

func runCase(bin, input, expect string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(strings.Split(out.String(), "\n")[0])
	if got != expect {
		return fmt.Errorf("expected %s got %s", expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD2.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, expect := generateCase(rng)
		if err := runCase(bin, input, expect); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

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

func solve(s, t string) string {
	n := len(s)
	m := len(t)
	pre := make([]int, m)
	idx := 0
	for i := 0; i < m; i++ {
		for idx < n && s[idx] != t[i] {
			idx++
		}
		if idx == n {
			pre[i] = n
		} else {
			pre[i] = idx
			idx++
		}
	}

	suf := make([]int, m)
	idx = n - 1
	for i := m - 1; i >= 0; i-- {
		for idx >= 0 && s[idx] != t[i] {
			idx--
		}
		if idx < 0 {
			suf[i] = -1
		} else {
			suf[i] = idx
			idx--
		}
	}

	maxDel := 0
	for i := 0; i <= m; i++ {
		var l, r int
		if i == 0 {
			l = 0
		} else {
			l = pre[i-1] + 1
		}
		if i == m {
			r = n - 1
		} else {
			r = suf[i] - 1
		}
		if r >= l && r-l+1 > maxDel {
			maxDel = r - l + 1
		}
	}
	return fmt.Sprintf("%d", maxDel)
}

func randomString(rng *rand.Rand, n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = byte('a' + rng.Intn(26))
	}
	return string(b)
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 1
	m := rng.Intn(n) + 1
	s := randomString(rng, n)
	// ensure t is subsequence: pick indices in order
	idxs := rng.Perm(n)[:m]
	sort.Ints(idxs)
	tbytes := make([]byte, m)
	for i, idx := range idxs {
		tbytes[i] = s[idx]
	}
	t := string(tbytes)
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierD1.go /path/to/binary")
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

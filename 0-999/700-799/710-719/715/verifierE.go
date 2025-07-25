package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Case struct {
	n int
	p []int
	q []int
}

func generatePermutation(n int, rng *rand.Rand) []int {
	perm := rng.Perm(n)
	for i := range perm {
		perm[i]++
	}
	return perm
}

func generateCase(rng *rand.Rand) Case {
	n := rng.Intn(5) + 2 //2..6
	p := generatePermutation(n, rng)
	q := generatePermutation(n, rng)
	// randomly replace with zeros
	for i := 0; i < n; i++ {
		if rng.Intn(3) == 0 {
			p[i] = 0
		}
		if rng.Intn(3) == 0 {
			q[i] = 0
		}
	}
	return Case{n, p, q}
}

func completions(base []int) [][]int {
	n := len(base)
	used := make([]bool, n+1)
	for _, v := range base {
		if v != 0 {
			used[v] = true
		}
	}
	res := [][]int{}
	var cur []int = make([]int, n)
	var dfs func(int)
	dfs = func(i int) {
		if i == n {
			tmp := make([]int, n)
			copy(tmp, cur)
			res = append(res, tmp)
			return
		}
		if base[i] != 0 {
			cur[i] = base[i]
			dfs(i + 1)
			return
		}
		for v := 1; v <= n; v++ {
			if !used[v] {
				used[v] = true
				cur[i] = v
				dfs(i + 1)
				used[v] = false
			}
		}
	}
	dfs(0)
	return res
}

func distance(p, q []int) int {
	n := len(p)
	mp := make([]int, n+1)
	for i := 0; i < n; i++ {
		mp[p[i]] = q[i]
	}
	vis := make([]bool, n+1)
	cycles := 0
	for i := 1; i <= n; i++ {
		if !vis[i] {
			cycles++
			j := i
			for !vis[j] {
				vis[j] = true
				j = mp[j]
			}
		}
	}
	return n - cycles
}

func expectedCounts(c Case) []int {
	cp := completions(c.p)
	cq := completions(c.q)
	counts := make([]int, c.n)
	for _, p := range cp {
		for _, q := range cq {
			d := distance(p, q)
			counts[d]++
		}
	}
	return counts
}

func runCase(bin string, c Case) error {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", c.n)
	for i, v := range c.p {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	for i, v := range c.q {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
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
	outFields := strings.Fields(strings.TrimSpace(out.String()))
	if len(outFields) != c.n {
		return fmt.Errorf("expected %d numbers", c.n)
	}
	expect := expectedCounts(c)
	for i := 0; i < c.n; i++ {
		val, err := strconv.Atoi(outFields[i])
		if err != nil {
			return fmt.Errorf("bad int")
		}
		if val != expect[i] {
			return fmt.Errorf("pos %d expected %d got %d", i, expect[i], val)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		c := generateCase(rng)
		if err := runCase(bin, c); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

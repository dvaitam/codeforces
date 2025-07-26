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

const mod int64 = 1000000007

func solveC(n, m int, gyms [][]int) int64 {
	typeCount := make(map[int]map[int]int)
	for i := 0; i < n; i++ {
		for _, t := range gyms[i] {
			mp := typeCount[t]
			if mp == nil {
				mp = make(map[int]int)
				typeCount[t] = mp
			}
			mp[i]++
		}
	}

	groups := make(map[string]int)
	buf := &bytes.Buffer{}
	idx := make([]int, 0)
	for t := 1; t <= m; t++ {
		mp := typeCount[t]
		if mp == nil {
			groups[""]++
			continue
		}
		buf.Reset()
		idx = idx[:0]
		for gym := range mp {
			idx = append(idx, gym)
		}
		sort.Ints(idx)
		for _, gym := range idx {
			fmt.Fprintf(buf, "%d:%d|", gym, mp[gym])
		}
		key := buf.String()
		groups[key]++
	}

	fact := make([]int64, m+1)
	fact[0] = 1
	for i := 1; i <= m; i++ {
		fact[i] = fact[i-1] * int64(i) % mod
	}

	ans := int64(1)
	for _, s := range groups {
		ans = ans * fact[s] % mod
	}
	return ans
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(3) + 1
	m := rng.Intn(4) + 1
	gyms := make([][]int, n)
	for i := 0; i < n; i++ {
		g := rng.Intn(3)
		gyms[i] = make([]int, g)
		for j := 0; j < g; j++ {
			gyms[i][j] = rng.Intn(m) + 1
		}
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	for i := 0; i < n; i++ {
		fmt.Fprintf(&sb, "%d", len(gyms[i]))
		for _, t := range gyms[i] {
			fmt.Fprintf(&sb, " %d", t)
		}
		sb.WriteByte('\n')
	}
	expected := solveC(n, m, gyms)
	return sb.String(), fmt.Sprintf("%d", expected)
}

func runCase(exe, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(exe, ".go") {
		cmd = exec.Command("go", "run", exe)
	} else {
		cmd = exec.Command(exe)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		got, err := runCase(exe, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

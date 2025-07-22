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

func solveE(n, k int, skills []int, flips [][2]int) int64 {
	sorted := make([]int, n)
	copy(sorted, skills)
	sort.Ints(sorted)
	dt := make([]int64, n+3)
	dsum := make([]int64, n+3)
	for _, f := range flips {
		a := f[0]
		b := f[1]
		l0 := sort.Search(n, func(i int) bool { return sorted[i] >= a })
		r0 := sort.Search(n, func(i int) bool { return sorted[i] > b }) - 1
		if l0 <= r0 {
			l := int64(l0 + 1)
			r := int64(r0 + 1)
			dt[l]++
			dt[r+1]--
			dsum[l] += l + r
			dsum[r+1] -= l + r
		}
	}
	t := make([]int64, n+2)
	sumlr := make([]int64, n+2)
	var ct, cs int64
	for u := 1; u <= n; u++ {
		ct += dt[u]
		cs += dsum[u]
		t[u] = ct
		sumlr[u] = cs
	}
	var trans int64
	for u := 1; u <= n; u++ {
		out := int64(u-1) + sumlr[u] - 2*int64(u)*t[u]
		if out > 1 {
			trans += out * (out - 1) / 2
		}
	}
	nn := int64(n)
	total := nn * (nn - 1) * (nn - 2) / 6
	return total - trans
}

func generateCase(rng *rand.Rand) (string, int64) {
	n := rng.Intn(10) + 3
	k := rng.Intn(10)
	skills := make([]int, n)
	used := make(map[int]bool)
	for i := 0; i < n; i++ {
		val := rng.Intn(100) + 1
		for used[val] {
			val = rng.Intn(100) + 1
		}
		used[val] = true
		skills[i] = val
	}
	flips := make([][2]int, k)
	for i := 0; i < k; i++ {
		a := rng.Intn(100) + 1
		b := a + rng.Intn(100)
		if a > b {
			a, b = b, a
		}
		flips[i] = [2]int{a, b}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
	for i, v := range skills {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteString("\n")
	for _, f := range flips {
		sb.WriteString(fmt.Sprintf("%d %d\n", f[0], f[1]))
	}
	ans := solveE(n, k, skills, flips)
	return sb.String(), ans
}

func runCase(bin, input string, exp int64) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var val int64
	if _, err := fmt.Sscan(strings.TrimSpace(out.String()), &val); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if val != exp {
		return fmt.Errorf("expected %d got %d", exp, val)
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
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

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

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func gcdex(a, b int) (g, x, y int) {
	if a == 0 {
		return b, 0, 1
	}
	g1, x1, y1 := gcdex(b%a, a)
	x = y1 - (b/a)*x1
	y = x1
	return g1, x, y
}

func inv(a, m int) int {
	_, x, _ := gcdex(a, m)
	x %= m
	if x < 0 {
		x += m
	}
	return x
}

func solveCase(n, m int, bad []int) string {
	forb := make([]bool, m)
	for _, v := range bad {
		forb[v] = true
	}
	mp := make(map[int][]int)
	for j := 0; j < m; j++ {
		if !forb[j] {
			d := gcd(m, j)
			mp[d] = append(mp[d], j)
		}
	}
	gcds := make([]int, 0, len(mp))
	for d := range mp {
		gcds = append(gcds, d)
	}
	sort.Ints(gcds)
	cnt := len(gcds)
	ans := make([]int, cnt)
	best := make([]int, cnt)
	for i := range best {
		best[i] = -1
	}
	for i := cnt - 1; i >= 0; i-- {
		for j := i + 1; j < cnt; j++ {
			if gcds[j]%gcds[i] == 0 && ans[j] > ans[i] {
				ans[i] = ans[j]
				best[i] = j
			}
		}
		ans[i] += len(mp[gcds[i]])
	}
	bestI := 0
	for i := 1; i < cnt; i++ {
		if ans[i] > ans[bestI] {
			bestI = i
		}
	}
	var out strings.Builder
	out.WriteString(fmt.Sprintf("%d\n", ans[bestI]))
	prevK, prevG := 1, 1
	for idx := bestI; idx != -1; idx = best[idx] {
		g := gcds[idx]
		for _, k := range mp[g] {
			modSeg := m / prevG
			a := k / prevG
			b := prevK / prevG
			x := (a * inv(b, modSeg)) % modSeg
			fmt.Fprintf(&out, "%d ", x)
			prevK = k
			prevG = g
		}
	}
	out.WriteString("\n")
	return out.String()
}

func genCase(rng *rand.Rand) (string, string) {
	m := rng.Intn(30) + 2
	n := rng.Intn(m)
	badSet := make(map[int]bool)
	for len(badSet) < n {
		v := rng.Intn(m)
		badSet[v] = true
	}
	bad := make([]int, 0, n)
	for v := range badSet {
		bad = append(bad, v)
	}
	sort.Ints(bad)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	if n > 0 {
		for i, v := range bad {
			if i > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", v)
		}
		sb.WriteByte('\n')
	}
	exp := solveCase(n, m, bad)
	return sb.String(), exp
}

func runCase(bin, in, exp string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(in)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, buf.String())
	}
	got := strings.TrimSpace(buf.String())
	if strings.TrimSpace(exp) != strings.TrimSpace(got) {
		return fmt.Errorf("expected \n%s\ngot \n%s", exp, got)
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
		in, exp := genCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

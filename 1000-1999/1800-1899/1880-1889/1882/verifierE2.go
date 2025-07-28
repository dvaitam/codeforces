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

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func calc(n int, a []int, w int) int {
	p := make([]int, n+1)
	vis := make([]bool, n+1)
	for i := 0; i <= n; i++ {
		p[i] = (a[i] + w) % (n + 1)
	}
	ans := n
	work(p, vis, 0, true)
	for i := 1; i <= n; i++ {
		if p[i] == i {
			ans--
		} else if !vis[i] {
			work(p, vis, i, false)
			ans++
		}
	}
	return ans
}

func work(p []int, vis []bool, x int, reset bool) {
	if reset {
		for i := range vis {
			vis[i] = false
		}
	}
	vis[x] = true
	for y := p[x]; y != x; y = p[y] {
		vis[y] = true
	}
}

func dfs(p []int, v *[]int, x, y int) {
	if x == y {
		return
	}
	dfs(p, v, p[x], y)
	*v = append(*v, x-p[x])
}

func solvePerm(n int, a []int, w int) []int {
	p := make([]int, n+1)
	vis := make([]bool, n+1)
	v := make([]int, 0)
	for i := 0; i <= n; i++ {
		p[i] = (a[i] + w) % (n + 1)
	}
	cur := 0
	work(p, vis, 0, true)
	for i := 1; i <= n; i++ {
		if p[i] != i && !vis[i] {
			v = append(v, i-cur)
			p[cur], p[i] = p[i], p[cur]
			cur = i
			work(p, vis, cur, true)
		}
	}
	dfs(p, &v, p[cur], cur)
	return v
}

func F(n, x int) int {
	if x > 0 {
		return x
	}
	return n + 1 + x
}

func expectedOutput(n, m int, a, b []int) string {
	f1 := make([]int, n+1)
	f2 := make([]int, m+1)
	for i := 0; i <= n; i++ {
		f1[i] = calc(n, a, i)
	}
	for j := 0; j <= m; j++ {
		f2[j] = calc(m, b, j)
	}
	ans := 1<<31 - 1
	bi, bj := -1, -1
	for i := 0; i <= n; i++ {
		for j := 0; j <= m; j++ {
			if (f1[i]^f2[j])%2 == 0 {
				mx := f1[i]
				if f2[j] > mx {
					mx = f2[j]
				}
				if mx < ans {
					ans = mx
					bi, bj = i, j
				}
			}
		}
	}
	if bi < 0 {
		return "-1"
	}
	v1 := solvePerm(n, append([]int(nil), a...), bi)
	v2 := solvePerm(m, append([]int(nil), b...), bj)
	for len(v1) < len(v2) {
		v1 = append(v1, 1, n)
	}
	for len(v2) < len(v1) {
		v2 = append(v2, 1, m)
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(v1)))
	for k := range v1 {
		sb.WriteString(fmt.Sprintf("%d %d\n", F(n, v1[k]), F(m, v2[k])))
	}
	return strings.TrimSpace(sb.String())
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(4) + 1
	m := rng.Intn(4) + 1
	p := rng.Perm(n)
	q := rng.Perm(m)
	a := make([]int, n+1)
	b := make([]int, m+1)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i, v := range p {
		a[i+1] = v + 1
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v+1))
	}
	sb.WriteByte('\n')
	for i, v := range q {
		b[i+1] = v + 1
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v+1))
	}
	sb.WriteByte('\n')
	expect := expectedOutput(n, m, a, b)
	return sb.String(), expect
}

func fixedCases() [][2]string {
	return [][2]string{
		{"1 1\n1\n1\n", "0"},
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE2.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	for idx, tc := range fixedCases() {
		out, err := runCandidate(bin, tc[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "fixed case %d failed: %v\ninput:\n%s", idx+1, err, tc[0])
			os.Exit(1)
		}
		if strings.TrimSpace(out) != tc[1] {
			fmt.Fprintf(os.Stderr, "fixed case %d failed: expected %s got %s\ninput:\n%s", idx+1, tc[1], out, tc[0])
			os.Exit(1)
		}
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

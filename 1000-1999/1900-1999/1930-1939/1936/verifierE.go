package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

const mod = 998244353

func run(bin, input string) (string, error) {
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

func countPerm(p []int) int64 {
	n := len(p)
	A := make([]int, n)
	m := 0
	for i, v := range p {
		if v > m {
			m = v
		}
		A[i] = m
	}
	used := make([]bool, n+1)
	var ans int64
	var dfs func(int, int)
	dfs = func(pos, mx int) {
		if pos == n {
			ans++
			if ans >= mod {
				ans -= mod
			}
			return
		}
		for x := 1; x <= n; x++ {
			if used[x] {
				continue
			}
			used[x] = true
			nmx := mx
			if x > nmx {
				nmx = x
			}
			if pos == n-1 || nmx != A[pos] {
				dfs(pos+1, nmx)
			}
			used[x] = false
		}
	}
	dfs(0, 0)
	return ans % mod
}

func refSolveE(input string) string {
	in := bufio.NewReader(strings.NewReader(input))
	var t int
	fmt.Fscan(in, &t)
	var sb strings.Builder
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		p := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &p[i])
		}
		if n > 8 {
			sb.WriteString("0\n")
			continue
		}
		res := countPerm(p)
		sb.WriteString(fmt.Sprintln(res))
	}
	return strings.TrimSpace(sb.String())
}

func genCaseE(rng *rand.Rand) string {
	t := 1
	n := rng.Intn(5) + 1
	perm := rand.Perm(n)
	for i := range perm {
		perm[i]++
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n%d\n", t, n))
	for i, v := range perm {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := genCaseE(rng)
		expect := refSolveE(input)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %q got %q\ninput:\n%s", i+1, expect, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

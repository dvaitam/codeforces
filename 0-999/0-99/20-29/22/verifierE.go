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

func solveCase(a []int) string {
	n := len(a) - 1
	d := make([]int, n+1)
	vis := make([]int, n+1)
	f := make([]int, n+1)
	for i := 1; i <= n; i++ {
		d[a[i]]++
	}
	var dfs func(int) int
	dfs = func(x int) int {
		if f[x] != 0 {
			return f[x]
		}
		if vis[x] != 0 {
			return 0
		}
		vis[x] = 1
		res := dfs(a[x])
		f[x] = res
		if f[x] != 0 {
			return f[x]
		}
		f[x] = x
		return x
	}
	for i := 1; i <= n; i++ {
		if d[i] == 0 || vis[i] == 0 {
			dfs(i)
		}
	}
	ans1 := make([]int, 0)
	ans2 := make([]int, 0)
	for i := 1; i <= n; i++ {
		if d[i] == 0 {
			ans1 = append(ans1, i)
			ans2 = append(ans2, f[i])
			vis[f[i]] = 2
		}
	}
	for i := 1; i <= n; i++ {
		if f[i] == i && vis[i] < 2 {
			ans1 = append(ans1, i)
			ans2 = append(ans2, i)
		}
	}
	m := len(ans1)
	if m == 1 && ans1[0] == ans2[0] {
		m = 0
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", m))
	for i := 0; i < m; i++ {
		j := (i + 1) % m
		sb.WriteString(fmt.Sprintf("%d %d\n", ans2[i], ans1[j]))
	}
	return strings.TrimSpace(sb.String())
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		n := rng.Intn(8) + 2
		a := make([]int, n+1)
		for i := 1; i <= n; i++ {
			v := rng.Intn(n) + 1
			if v == i {
				if v < n {
					v++
				} else {
					v--
				}
			}
			a[i] = v
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i := 1; i <= n; i++ {
			if i > 1 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", a[i]))
		}
		sb.WriteByte('\n')
		input := sb.String()
		expected := solveCase(a)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", t+1, err, input)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected:\n%s\ngot:\n%s\ninput:\n%s", t+1, expected, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

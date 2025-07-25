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

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func lcm(a, b int64) int64 {
	return a / gcd(a, b) * b
}

func expected(n int, f []int) string {
	for i := range f {
		f[i]--
	}
	indeg := make([]int, n)
	for i := 0; i < n; i++ {
		indeg[f[i]]++
	}
	inCycle := make([]bool, n)
	for i := range inCycle {
		inCycle[i] = true
	}
	queue := make([]int, 0)
	for i := 0; i < n; i++ {
		if indeg[i] == 0 {
			queue = append(queue, i)
		}
	}
	for idx := 0; idx < len(queue); idx++ {
		u := queue[idx]
		inCycle[u] = false
		v := f[u]
		indeg[v]--
		if indeg[v] == 0 {
			queue = append(queue, v)
		}
	}
	visited := make([]bool, n)
	L := int64(1)
	for i := 0; i < n; i++ {
		if inCycle[i] && !visited[i] {
			v := i
			cnt := int64(0)
			for {
				visited[v] = true
				cnt++
				v = f[v]
				if v == i {
					break
				}
			}
			L = lcm(L, cnt)
		}
	}
	depth := make([]int, n)
	var calc func(int) int
	calc = func(u int) int {
		if inCycle[u] {
			return 0
		}
		if depth[u] != 0 {
			return depth[u]
		}
		depth[u] = calc(f[u]) + 1
		return depth[u]
	}
	maxd := 0
	for i := 0; i < n; i++ {
		d := calc(i)
		if d > maxd {
			maxd = d
		}
	}
	var k int64
	if maxd == 0 {
		k = L
	} else {
		t := (int64(maxd) + L - 1) / L
		if t < 1 {
			t = 1
		}
		k = t * L
	}
	return fmt.Sprintf("%d", k)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(6) + 1
		f := make([]int, n)
		for j := 0; j < n; j++ {
			f[j] = rng.Intn(n) + 1
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", f[j]))
		}
		sb.WriteByte('\n')
		input := sb.String()
		exp := expected(n, append([]int(nil), f...))
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

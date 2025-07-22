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

func solveD(b [][]int) []int {
	n := len(b)
	a := make([]int, n)
	for k := 0; k <= 30; k++ {
		deg := make([]int, n)
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				if i != j && ((b[i][j]>>k)&1) == 1 {
					deg[i]++
				}
			}
		}
		visited := make([]bool, n)
		for i := 0; i < n; i++ {
			if visited[i] || deg[i] == 0 {
				continue
			}
			queue := []int{i}
			comp := []int{i}
			visited[i] = true
			for head := 0; head < len(queue); head++ {
				u := queue[head]
				for v := 0; v < n; v++ {
					if !visited[v] && ((b[u][v]>>k)&1) == 1 {
						visited[v] = true
						queue = append(queue, v)
						comp = append(comp, v)
					}
				}
			}
			for _, u := range comp {
				a[u] |= 1 << k
			}
		}
	}
	return a
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(4) + 1
	a := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = rng.Intn(1000)
	}
	b := make([][]int, n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		b[i] = make([]int, n)
		for j := 0; j < n; j++ {
			if i == j {
				b[i][j] = -1
			} else {
				b[i][j] = a[i] & a[j]
			}
			sb.WriteString(fmt.Sprintf("%d", b[i][j]))
			if j+1 < n {
				sb.WriteByte(' ')
			}
		}
		sb.WriteByte('\n')
	}
	ansArr := solveD(b)
	var exp strings.Builder
	for i, v := range ansArr {
		if i > 0 {
			exp.WriteByte(' ')
		}
		exp.WriteString(fmt.Sprint(v))
	}
	return sb.String(), exp.String()
}

func runCase(bin, input, expected string) error {
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
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

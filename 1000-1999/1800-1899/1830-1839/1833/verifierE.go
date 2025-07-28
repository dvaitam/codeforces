package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func solveCase(a []int) (int, int) {
	n := len(a) - 1
	adj := make([][]int, n+1)
	for i := 1; i <= n; i++ {
		j := a[i]
		exists := false
		for _, v := range adj[i] {
			if v == j {
				exists = true
				break
			}
		}
		if !exists {
			adj[i] = append(adj[i], j)
		}
		exists = false
		for _, v := range adj[j] {
			if v == i {
				exists = true
				break
			}
		}
		if !exists {
			adj[j] = append(adj[j], i)
		}
	}
	visited := make([]bool, n+1)
	components, cycleComp := 0, 0
	for i := 1; i <= n; i++ {
		if visited[i] {
			continue
		}
		components++
		queue := []int{i}
		visited[i] = true
		isCycle := true
		for len(queue) > 0 {
			v := queue[0]
			queue = queue[1:]
			if len(adj[v]) != 2 {
				isCycle = false
			}
			for _, to := range adj[v] {
				if !visited[to] {
					visited[to] = true
					queue = append(queue, to)
				}
			}
		}
		if isCycle {
			cycleComp++
		}
	}
	pathComp := components - cycleComp
	maxCycles := components
	minCycles := cycleComp
	if pathComp > 0 {
		minCycles++
	}
	return minCycles, maxCycles
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 2
	a := make([]int, n+1)
	for i := 1; i <= n; i++ {
		j := rng.Intn(n) + 1
		if j == i {
			j = (j % n) + 1
		}
		a[i] = j
	}
	input := fmt.Sprintf("1\n%d\n", n)
	for i := 1; i <= n; i++ {
		if i > 1 {
			input += " "
		}
		input += fmt.Sprintf("%d", a[i])
	}
	input += "\n"
	minC, maxC := solveCase(a)
	exp := fmt.Sprintf("%d %d\n", minC, maxC)
	return input, exp
}

func runCase(bin, input, exp string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	if strings.TrimSpace(out.String()) != strings.TrimSpace(exp) {
		return fmt.Errorf("expected %q got %q", strings.TrimSpace(exp), strings.TrimSpace(out.String()))
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(42))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

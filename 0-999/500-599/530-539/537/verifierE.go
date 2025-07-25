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

func expectedAnswerE(n int, edges [][2]int) string {
	children := make([][]int, n+1)
	for _, e := range edges {
		u, v := e[0], e[1]
		children[u] = append(children[u], v)
	}
	depth := make([]int, n+1)
	queue := []int{1}
	for i := 0; i < len(queue); i++ {
		u := queue[i]
		for _, v := range children[u] {
			depth[v] = depth[u] + 1
			queue = append(queue, v)
		}
	}
	stack := []int{1}
	order := make([]int, 0, n)
	for len(stack) > 0 {
		u := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		order = append(order, u)
		for _, v := range children[u] {
			stack = append(stack, v)
		}
	}
	need1 := make([]int, n+1)
	need2 := make([]int, n+1)
	m := 0
	for i := len(order) - 1; i >= 0; i-- {
		u := order[i]
		if len(children[u]) == 0 {
			need1[u] = 1
			need2[u] = 1
			m++
		} else if depth[u]%2 == 0 {
			mn1 := n + 5
			sum2 := 0
			for _, v := range children[u] {
				if need1[v] < mn1 {
					mn1 = need1[v]
				}
				sum2 += need2[v]
			}
			need1[u] = mn1
			need2[u] = sum2
		} else {
			sum1 := 0
			mn2 := n + 5
			for _, v := range children[u] {
				sum1 += need1[v]
				if need2[v] < mn2 {
					mn2 = need2[v]
				}
			}
			need1[u] = sum1
			need2[u] = mn2
		}
	}
	maxRes := m - need1[1] + 1
	minRes := need2[1]
	return fmt.Sprintf("%d %d", maxRes, minRes)
}

func generateCaseE(rng *rand.Rand) (int, [][2]int) {
	n := rng.Intn(10) + 1
	edges := make([][2]int, 0, n-1)
	for v := 2; v <= n; v++ {
		p := rng.Intn(v-1) + 1
		edges = append(edges, [2]int{p, v})
	}
	return n, edges
}

func runCaseE(bin string, n int, edges [][2]int) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprint(n, "\n"))
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	expected := expectedAnswerE(n, edges)
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, edges := generateCaseE(rng)
		if err := runCaseE(bin, n, edges); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

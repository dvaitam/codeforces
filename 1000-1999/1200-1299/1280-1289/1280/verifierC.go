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

type Edge struct {
	to int
	w  int64
}

func solveC(k int, edges [][3]int) (int64, int64) {
	n := 2 * k
	g := make([][]Edge, n+1)
	for _, e := range edges {
		a, b := e[0], e[1]
		w := int64(e[2])
		g[a] = append(g[a], Edge{b, w})
		g[b] = append(g[b], Edge{a, w})
	}
	parent := make([]int, n+1)
	weight := make([]int64, n+1)
	order := make([]int, 0, n)
	stack := []int{1}
	parent[1] = -1
	for len(stack) > 0 {
		u := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		order = append(order, u)
		for _, e := range g[u] {
			if e.to == parent[u] {
				continue
			}
			parent[e.to] = u
			weight[e.to] = e.w
			stack = append(stack, e.to)
		}
	}
	size := make([]int, n+1)
	var G, B int64
	for i := len(order) - 1; i >= 0; i-- {
		u := order[i]
		size[u]++
		if parent[u] != -1 {
			s := size[u]
			if s%2 == 1 {
				G += weight[u]
			}
			if s < n-s {
				B += int64(s) * weight[u]
			} else {
				B += int64(n-s) * weight[u]
			}
			size[parent[u]] += s
		}
	}
	return G, B
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	t := 100
	var input bytes.Buffer
	fmt.Fprintln(&input, t)
	expected := make([]string, t)
	for i := 0; i < t; i++ {
		k := rand.Intn(3) + 1
		n := 2 * k
		edges := make([][3]int, n-1)
		for j := 2; j <= n; j++ {
			p := rand.Intn(j-1) + 1
			w := rand.Intn(5) + 1
			edges[j-2] = [3]int{p, j, w}
		}
		fmt.Fprintln(&input, k)
		for _, e := range edges {
			fmt.Fprintln(&input, e[0], e[1], e[2])
		}
		g, b := solveC(k, edges)
		expected[i] = fmt.Sprintf("%d %d", g, b)
	}
	cmd := exec.Command(bin)
	cmd.Stdin = &input
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("binary error:", err)
		fmt.Print(string(out))
		return
	}
	outputs := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(outputs) != t {
		fmt.Printf("expected %d lines, got %d\n", t, len(outputs))
		fmt.Print(string(out))
		return
	}
	for i := 0; i < t; i++ {
		if strings.TrimSpace(outputs[i]) != expected[i] {
			fmt.Printf("Test %d failed: expected %s got %s\n", i+1, expected[i], strings.TrimSpace(outputs[i]))
			return
		}
	}
	fmt.Println("All tests passed!")
}

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func solveE(n int, dom [][2]int) string {
	adj := make([][]int, n+1)
	deg := make([]int, n+1)
	ok := true
	for _, d := range dom {
		a, b := d[0], d[1]
		if a == b {
			ok = false
		}
		deg[a]++
		deg[b]++
		adj[a] = append(adj[a], b)
		adj[b] = append(adj[b], a)
	}
	if ok {
		for i := 1; i <= n; i++ {
			if deg[i] > 2 {
				ok = false
				break
			}
		}
	}
	if ok {
		vis := make([]bool, n+1)
		for i := 1; i <= n && ok; i++ {
			if vis[i] || deg[i] == 0 {
				continue
			}
			stack := []int{i}
			vis[i] = true
			nodes, edges := 0, 0
			for len(stack) > 0 {
				v := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				nodes++
				edges += len(adj[v])
				for _, to := range adj[v] {
					if !vis[to] {
						vis[to] = true
						stack = append(stack, to)
					}
				}
			}
			edges /= 2
			if edges == nodes && nodes%2 == 1 {
				ok = false
				break
			}
		}
	}
	if ok {
		return "YES"
	}
	return "NO"
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]

	rand.Seed(1)
	const T = 100
	var input bytes.Buffer
	fmt.Fprintln(&input, T)
	expected := make([]string, T)
	for i := 0; i < T; i++ {
		n := rand.Intn(10)/2*2 + 2 // even and at least 2
		fmt.Fprintln(&input, n)
		dom := make([][2]int, n)
		for j := 0; j < n; j++ {
			a := rand.Intn(n) + 1
			b := rand.Intn(n) + 1
			dom[j] = [2]int{a, b}
			fmt.Fprintf(&input, "%d %d\n", a, b)
		}
		expected[i] = solveE(n, dom)
	}

	cmd := exec.Command(binary)
	cmd.Stdin = bytes.NewReader(input.Bytes())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("failed to run binary:", err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(bytes.NewReader(out.Bytes()))
	i := 0
	for scanner.Scan() {
		if i >= T {
			fmt.Println("binary produced extra output")
			os.Exit(1)
		}
		got := strings.TrimSpace(scanner.Text())
		if got != expected[i] {
			fmt.Printf("test %d failed: expected %s got %s\n", i+1, expected[i], got)
			os.Exit(1)
		}
		i++
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("reading output failed:", err)
		os.Exit(1)
	}
	if i < T {
		fmt.Println("binary produced insufficient output")
		os.Exit(1)
	}

	fmt.Println("all tests passed")
}

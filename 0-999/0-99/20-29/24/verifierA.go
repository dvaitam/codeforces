package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func expected(n int, edges [][3]int) int {
	dirCost := make(map[[2]int]int)
	adj := make(map[int][]int)
	
	for _, edge := range edges {
		a, b, c := edge[0], edge[1], edge[2]
		dirCost[[2]int{a, b}] = c
		adj[a] = append(adj[a], b)
		adj[b] = append(adj[b], a)
	}
	
	order := make([]int, 0, n)
	start := 0
	for u := range adj {
		start = u
		break
	}
	
	cur := start
	prev := -1
	for len(order) < n {
		order = append(order, cur)
		for _, nei := range adj[cur] {
			if nei != prev {
				prev, cur = cur, nei
				break
			}
		}
	}
	
	costCW := 0
	costCCW := 0
	for i := 0; i < n; i++ {
		u := order[i]
		v := order[(i+1)%n]
		
		if _, ok := dirCost[[2]int{u, v}]; ok {
		} else if c, ok2 := dirCost[[2]int{v, u}]; ok2 {
			costCW += c
		}
		
		if _, ok := dirCost[[2]int{v, u}]; ok {
		} else if c, ok2 := dirCost[[2]int{u, v}]; ok2 {
			costCCW += c
		}
	}
	
	if costCW < costCCW {
		return costCW
	}
	return costCCW
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	
	for i := 0; i < 100; i++ {
		n := rng.Intn(98) + 3
		
		nodes := make([]int, n)
		for j := 0; j < n; j++ {
			nodes[j] = j + 1
		}
		
		for j := 0; j < n-1; j++ {
			k := j + rng.Intn(n-j)
			nodes[j], nodes[k] = nodes[k], nodes[j]
		}
		
		edges := make([][3]int, n)
		for j := 0; j < n; j++ {
			u := nodes[j]
			v := nodes[(j+1)%n]
			c := rng.Intn(100) + 1
			
			if rng.Intn(2) == 0 {
				edges[j] = [3]int{u, v, c}
			} else {
				edges[j] = [3]int{v, u, c}
			}
		}
		
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d\n", n))
		for _, edge := range edges {
			input.WriteString(fmt.Sprintf("%d %d %d\n", edge[0], edge[1], edge[2]))
		}
		
		expectedOut := expected(n, edges)
		got, err := run(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input.String())
			os.Exit(1)
		}
		
		gotInt, parseErr := strconv.Atoi(got)
		if parseErr != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: cannot parse output: %v\ninput:\n%s", i+1, parseErr, input.String())
			os.Exit(1)
		}
		
		if gotInt != expectedOut {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\ninput:\n%s", i+1, expectedOut, gotInt, input.String())
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
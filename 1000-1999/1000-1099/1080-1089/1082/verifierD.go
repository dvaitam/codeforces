package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func runCmd(name string, args []string, input string) (string, error) {
	cmd := exec.Command(name, args...)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func genCase(seed int64) (int, []int, string) {
	rand.Seed(seed)
	n := rand.Intn(500-3+1) + 3
	a := make([]int, n)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for j := 0; j < n; j++ {
		a[j] = rand.Intn(n-1) + 1
		fmt.Fprintf(&sb, "%d ", a[j])
	}
	sb.WriteByte('\n')
	return n, a, sb.String()
}

func computeDiam(n int, adj [][]int) int {
	maxD := 0
	for i := 1; i <= n; i++ {
		dist := make([]int, n+1)
		for j := 1; j <= n; j++ {
			dist[j] = -1
		}
		dist[i] = 0
		q := []int{i}
		for len(q) > 0 {
			u := q[0]
			q = q[1:]
			if dist[u] > maxD {
				maxD = dist[u]
			}
			for _, v := range adj[u] {
				if dist[v] == -1 {
					dist[v] = dist[u] + 1
					q = append(q, v)
				}
			}
		}
		for j := 1; j <= n; j++ {
			if dist[j] == -1 {
				return -1 // not connected
			}
		}
	}
	return maxD
}

func checkCandidate(n int, a []int, out string) error {
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return fmt.Errorf("empty output")
	}
	
	sum := 0
	cntGe2 := 0
	for _, x := range a {
		sum += x
		if x >= 2 {
			cntGe2++
		}
	}
	
	if strings.ToUpper(fields[0]) == "NO" {
		if sum >= 2*n - 2 {
			return fmt.Errorf("candidate says NO but solution exists")
		}
		return nil
	}
	
	if strings.ToUpper(fields[0]) != "YES" {
		return fmt.Errorf("expected YES or NO, got %s", fields[0])
	}
	
	if sum < 2*n - 2 {
		return fmt.Errorf("candidate says YES but no solution exists")
	}
	
	if len(fields) < 2 {
		return fmt.Errorf("missing diameter")
	}
	
	diam, err := strconv.Atoi(fields[1])
	if err != nil {
		return fmt.Errorf("invalid diameter: %s", fields[1])
	}
	
	expectedDiam := cntGe2 + 1
	if expectedDiam > n - 1 {
		expectedDiam = n - 1
	}
	
	if diam != expectedDiam {
		return fmt.Errorf("expected diameter %d, got %d", expectedDiam, diam)
	}
	
	if len(fields) < 3 {
		return fmt.Errorf("missing number of edges")
	}
	
	m, err := strconv.Atoi(fields[2])
	if err != nil {
		return fmt.Errorf("invalid number of edges: %s", fields[2])
	}
	
	if len(fields) < 3 + 2*m {
		return fmt.Errorf("not enough edges in output")
	}
	
	deg := make([]int, n+1)
	adj := make([][]int, n+1)
	edges := make(map[string]bool)
	
	idx := 3
	for i := 0; i < m; i++ {
		u, err1 := strconv.Atoi(fields[idx])
		v, err2 := strconv.Atoi(fields[idx+1])
		idx += 2
		if err1 != nil || err2 != nil || u < 1 || u > n || v < 1 || v > n {
			return fmt.Errorf("invalid edge: %s %s", fields[idx-2], fields[idx-1])
		}
		if u == v {
			return fmt.Errorf("self loop: %d %d", u, v)
		}
		if u > v {
			u, v = v, u
		}
		edgeStr := fmt.Sprintf("%d-%d", u, v)
		if edges[edgeStr] {
			return fmt.Errorf("multiple edges: %d %d", u, v)
		}
		edges[edgeStr] = true
		deg[u]++
		deg[v]++
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	
	for i := 1; i <= n; i++ {
		if deg[i] > a[i-1] {
			return fmt.Errorf("degree of %d is %d > %d", i, deg[i], a[i-1])
		}
	}
	
	actualDiam := computeDiam(n, adj)
	if actualDiam == -1 {
		return fmt.Errorf("graph is not connected")
	}
	
	if actualDiam != diam {
		return fmt.Errorf("actual diameter %d != printed diameter %d", actualDiam, diam)
	}
	
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	
	bin := os.Args[1]
	
	for i := int64(0); i < 100; i++ {
		n, a, input := genCase(i + 12345)
		
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
		err := cmd.Run()
		
		if err != nil {
			fmt.Printf("candidate error on case %d: %v\n%s\n", i, err, out.String())
			os.Exit(1)
		}
		
		err = checkCandidate(n, a, out.String())
		if err != nil {
			fmt.Printf("case %d failed: %v\ninput:\n%s\ngot:\n%s\n", i, err, input, out.String())
			os.Exit(1)
		}
	}
	
	fmt.Println("All tests passed")
}
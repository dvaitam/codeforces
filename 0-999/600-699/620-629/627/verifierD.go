package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func checkD(n int, k int, a []int, edges [][2]int, th int) bool {
	g := make([][]int, n)
	for _, e := range edges {
		u, v := e[0], e[1]
		g[u] = append(g[u], v)
		g[v] = append(g[v], u)
	}
	vis := make([]bool, n)
	stack := make([]int, 0, n)
	for i := 0; i < n; i++ {
		if !vis[i] && a[i] >= th {
			size := 0
			stack = append(stack[:0], i)
			vis[i] = true
			for len(stack) > 0 {
				v := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				size++
				for _, to := range g[v] {
					if !vis[to] && a[to] >= th {
						vis[to] = true
						stack = append(stack, to)
					}
				}
			}
			if size >= k {
				return true
			}
		}
	}
	return false
}

func expectedD(n, k int, a []int, edges [][2]int) int {
	maxA := 0
	for _, v := range a {
		if v > maxA {
			maxA = v
		}
	}
	low, high := 1, maxA
	for low < high {
		mid := (low + high + 1) / 2
		if checkD(n, k, a, edges, mid) {
			low = mid
		} else {
			high = mid - 1
		}
	}
	return low
}

func generateD(rng *rand.Rand) (string, int) {
	n := rng.Intn(10) + 1
	k := rng.Intn(n) + 1
	a := make([]int, n)
	for i := range a {
		a[i] = rng.Intn(20) + 1
	}
	edges := make([][2]int, n-1)
	for i := 1; i < n; i++ {
		p := rng.Intn(i)
		edges[i-1] = [2]int{i, p}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, k)
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprint(&sb, v)
	}
	sb.WriteByte('\n')
	for _, e := range edges {
		fmt.Fprintf(&sb, "%d %d\n", e[0]+1, e[1]+1)
	}
	exp := expectedD(n, k, a, edges)
	return sb.String(), exp
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
		return
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(45))
	for i := 0; i < 100; i++ {
		input, exp := generateD(rng)
		cmd := exec.Command(exe)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		if err := cmd.Run(); err != nil {
			fmt.Printf("case %d runtime error: %v\n%s", i+1, err, out.String())
			return
		}
		got := strings.TrimSpace(out.String())
		if got != fmt.Sprint(exp) {
			fmt.Printf("case %d failed: expected %d got %s\ninput:\n%s", i+1, exp, got, input)
			return
		}
	}
	fmt.Println("All tests passed")
}

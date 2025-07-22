package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func run(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func runRef(input string) (string, error) {
	cmd := exec.Command("go", "run", "317C.go")
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rand.Seed(3)
	for i := 0; i < 100; i++ {
		n := rand.Intn(4) + 2
		v := rand.Intn(10) + 1
		maxEdges := n * (n - 1) / 2
		m := rand.Intn(maxEdges + 1)
		a := make([]int, n)
		b := make([]int, n)
		for j := 0; j < n; j++ {
			a[j] = rand.Intn(v + 1)
		}
		for j := 0; j < n; j++ {
			b[j] = rand.Intn(v + 1)
		}
		edges := make([][2]int, 0, m)
		used := make(map[[2]int]bool)
		for len(edges) < m {
			x := rand.Intn(n) + 1
			y := rand.Intn(n) + 1
			if x == y {
				continue
			}
			pair := [2]int{x, y}
			if x > y {
				pair = [2]int{y, x}
			}
			if used[pair] {
				continue
			}
			used[pair] = true
			edges = append(edges, pair)
		}
		input := fmt.Sprintf("%d %d %d\n", n, v, m)
		for j := 0; j < n; j++ {
			input += fmt.Sprintf("%d ", a[j])
		}
		input = strings.TrimSpace(input) + "\n"
		for j := 0; j < n; j++ {
			input += fmt.Sprintf("%d ", b[j])
		}
		input = strings.TrimSpace(input) + "\n"
		for _, e := range edges {
			input += fmt.Sprintf("%d %d\n", e[0], e[1])
		}
		exp, err := runRef(input)
		if err != nil {
			fmt.Println("reference run error:", err)
			return
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("binary run error on test %d: %v\n", i+1, err)
			return
		}
		if exp != got {
			fmt.Printf("mismatch on test %d\ninput:\n%sexpected:\n%s\n got:\n%s\n", i+1, input, exp, got)
			return
		}
	}
	fmt.Println("all tests passed")
}

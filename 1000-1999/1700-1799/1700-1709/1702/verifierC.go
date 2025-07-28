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

func solveC(n, k int, arr []int, queries [][2]int) []string {
	first := make(map[int]int)
	last := make(map[int]int)
	for i, u := range arr {
		if _, ok := first[u]; !ok {
			first[u] = i
		}
		last[u] = i
	}
	res := make([]string, len(queries))
	for i, q := range queries {
		a, b := q[0], q[1]
		fa, okA := first[a]
		lb, okB := last[b]
		if okA && okB && fa <= lb {
			res[i] = "YES"
		} else {
			res[i] = "NO"
		}
	}
	return res
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]

	rand.Seed(1)
	const T = 100
	var input bytes.Buffer
	fmt.Fprintln(&input, T)
	expected := []string{}

	for t := 0; t < T; t++ {
		n := rand.Intn(10) + 1
		k := rand.Intn(10) + 1
		fmt.Fprintf(&input, "%d %d\n", n, k)
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			arr[i] = rand.Intn(10) + 1
			fmt.Fprintf(&input, "%d ", arr[i])
		}
		fmt.Fprintln(&input)
		queries := make([][2]int, k)
		for i := 0; i < k; i++ {
			a := rand.Intn(10) + 1
			b := rand.Intn(10) + 1
			queries[i] = [2]int{a, b}
			fmt.Fprintf(&input, "%d %d\n", a, b)
		}
		out := solveC(n, k, arr, queries)
		expected = append(expected, out...)
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
		if i >= len(expected) {
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
	if i < len(expected) {
		fmt.Println("binary produced insufficient output")
		os.Exit(1)
	}

	fmt.Println("all tests passed")
}

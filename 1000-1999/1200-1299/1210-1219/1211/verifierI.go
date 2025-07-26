package main

import (
	"bytes"
	"fmt"
	"math/bits"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func runProgram(bin string, input []byte) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func genInput() (string, []int, [][2]int) {
	n := rand.Intn(10) + 1
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = rand.Intn(16)
	}
	edges := make([][2]int, 0)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if bits.OnesCount(uint(arr[i]^arr[j])) == 1 {
				edges = append(edges, [2]int{i + 1, j + 1})
			}
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, len(edges)))
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	return sb.String(), arr, edges
}

func checkOutput(out string, n int, edges [][2]int) bool {
	fields := strings.Fields(out)
	if len(fields) != n {
		return false
	}
	val := make([]int, n)
	for i, f := range fields {
		var v int
		if _, err := fmt.Sscan(f, &v); err != nil {
			return false
		}
		if v < 0 || v > 15 {
			return false
		}
		val[i] = v
	}
	adj := make(map[[2]int]bool)
	for _, e := range edges {
		if e[0] > e[1] {
			e[0], e[1] = e[1], e[0]
		}
		adj[e] = true
	}
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			diff := bits.OnesCount(uint(val[i]^val[j])) == 1
			_, has := adj[[2]int{i + 1, j + 1}]
			if diff != has {
				return false
			}
		}
	}
	return true
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierI.go /path/to/binary")
		return
	}
	rand.Seed(time.Now().UnixNano())
	bin := os.Args[1]

	for t := 0; t < 100; t++ {
		input, arr, edges := genInput()
		out, err := runProgram(bin, []byte(input))
		if err != nil || !checkOutput(out, len(arr), edges) {
			fmt.Println("Test", t+1, "failed")
			fmt.Println("Input:\n", input)
			fmt.Println("Output:\n", out)
			if err != nil {
				fmt.Println("Error:", err)
			}
			return
		}
	}
	fmt.Println("All tests passed")
}

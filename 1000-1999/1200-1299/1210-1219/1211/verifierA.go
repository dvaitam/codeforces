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

func runProgram(bin string, input []byte) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func hasTriple(r []int) bool {
	n := len(r)
	sufMaxIdx := make([]int, n)
	maxIdx := n - 1
	sufMaxIdx[n-1] = maxIdx
	for i := n - 2; i >= 0; i-- {
		if r[i] >= r[maxIdx] {
			maxIdx = i
		}
		sufMaxIdx[i] = maxIdx
	}
	minIdx := 0
	for j := 1; j < n-1; j++ {
		if r[minIdx] < r[j] && r[j] < r[sufMaxIdx[j+1]] {
			return true
		}
		if r[j] < r[minIdx] {
			minIdx = j
		}
	}
	return false
}

func checkOutput(out string, r []int) bool {
	parts := strings.Fields(out)
	if len(parts) != 3 {
		return false
	}
	a := make([]int, 3)
	for i := 0; i < 3; i++ {
		var v int
		_, err := fmt.Sscan(parts[i], &v)
		if err != nil {
			return false
		}
		a[i] = v
	}
	if a[0] == -1 && a[1] == -1 && a[2] == -1 {
		return !hasTriple(r)
	}
	n := len(r)
	if a[0] < 1 || a[0] > n || a[1] < 1 || a[1] > n || a[2] < 1 || a[2] > n {
		return false
	}
	if !(a[0] != a[1] && a[1] != a[2] && a[0] != a[2]) {
		return false
	}
	if !(r[a[0]-1] < r[a[1]-1] && r[a[1]-1] < r[a[2]-1]) {
		return false
	}
	return true
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		return
	}
	rand.Seed(time.Now().UnixNano())
	bin := os.Args[1]

	for t := 0; t < 100; t++ {
		n := rand.Intn(18) + 3
		rVals := make([]int, n)
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			rVals[i] = rand.Intn(100) + 1
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(fmt.Sprintf("%d", rVals[i]))
		}
		input.WriteByte('\n')
		inBytes := []byte(input.String())
		out, err := runProgram(bin, inBytes)
		if err != nil || !checkOutput(out, rVals) {
			fmt.Println("Test", t+1, "failed")
			fmt.Println("Input:\n", input.String())
			fmt.Println("Output:\n", out)
			if err != nil {
				fmt.Println("Error:", err)
			}
			return
		}
	}
	fmt.Println("All tests passed")
}

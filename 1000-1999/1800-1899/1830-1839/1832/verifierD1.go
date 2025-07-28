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

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solve(arr []int, k int) int {
	n := len(arr)
	a := append([]int(nil), arr...)
	colors := make([]bool, n)
	maxOps := k
	if maxOps > 2000 {
		maxOps = 2000
	}
	for i := 1; i <= maxOps; i++ {
		idx := 0
		for j := 1; j < n; j++ {
			if a[j] < a[idx] {
				idx = j
			}
		}
		if !colors[idx] {
			a[idx] += i
			colors[idx] = true
		} else {
			a[idx] -= i
			colors[idx] = false
		}
	}
	mn := a[0]
	for _, v := range a {
		if v < mn {
			mn = v
		}
	}
	return mn
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(4) + 1
	k := rng.Intn(10) + 1
	arr := make([]int, n)
	for i := range arr {
		arr[i] = rng.Intn(50)
	}
	input := fmt.Sprintf("%d 1\n", n)
	for i := 0; i < n; i++ {
		if i > 0 {
			input += " "
		}
		input += fmt.Sprintf("%d", arr[i])
	}
	input += "\n"
	input += fmt.Sprintf("%d\n", k)
	return input, fmt.Sprintf("%d", solve(arr, k))
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

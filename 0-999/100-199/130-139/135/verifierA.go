package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solveCase(a []int) []int {
	sort.Ints(a)
	n := len(a)
	b := make([]int, n)
	k1 := 0
	for k1 < n && a[k1] == 1 {
		k1++
	}
	if k1 < n {
		for i := 0; i <= k1 && i < n; i++ {
			b[i] = 1
		}
		for i := k1 + 1; i < n; i++ {
			b[i] = a[i-1]
		}
	} else {
		for i := 0; i < n-1; i++ {
			b[i] = 1
		}
		if n > 0 {
			b[n-1] = 2
		}
	}
	return b
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(20) + 1
	arr := make([]int, n)
	for i := range arr {
		arr[i] = rng.Intn(9) + 1
	}
	expectArr := solveCase(append([]int(nil), arr...))
	input := fmt.Sprintf("%d\n", n)
	for i, v := range arr {
		if i > 0 {
			input += " "
		}
		input += fmt.Sprint(v)
	}
	input += "\n"
	var exp strings.Builder
	for i, v := range expectArr {
		if i > 0 {
			exp.WriteByte(' ')
		}
		exp.WriteString(fmt.Sprint(v))
	}
	return input, exp.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
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
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

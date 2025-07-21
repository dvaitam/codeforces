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

func runBinary(bin, input string) (string, error) {
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
	return strings.TrimSpace(out.String()), err
}

func solve(a []int) string {
	maxVal := 0
	for _, v := range a {
		if v > maxVal {
			maxVal = v
		}
	}
	b := make([]int, maxVal+1)
	for _, v := range a {
		b[v]++
	}
	for i := 2; i < len(b); i++ {
		if b[i] > b[i-1] {
			return "-1"
		}
	}
	var sb strings.Builder
	if len(b) > 1 {
		sb.WriteString(fmt.Sprintf("%d\n", b[1]))
	} else {
		sb.WriteString("0\n")
	}
	for _, v := range a {
		sb.WriteString(fmt.Sprintf("%d ", b[v]))
		b[v]--
	}
	sb.WriteByte('\n')
	return strings.TrimSpace(sb.String())
}

func generateCase(r *rand.Rand) (string, string) {
	n := r.Intn(10) + 1
	a := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = r.Intn(5) + 1
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(a[i]))
	}
	sb.WriteByte('\n')
	return sb.String(), solve(a)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(r)
		out, err := runBinary(bin, in)
		if err != nil {
			fmt.Printf("Test %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != exp {
			fmt.Printf("Test %d failed.\nInput:\n%sExpected:\n%s\nGot:\n%s\n", i+1, in, exp, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

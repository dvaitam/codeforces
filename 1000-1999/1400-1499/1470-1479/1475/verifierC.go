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

func solveC(a, b, k int, boys, girls []int) string {
	cntBoy := make([]int, a+1)
	cntGirl := make([]int, b+1)
	for i := 0; i < k; i++ {
		cntBoy[boys[i]]++
		cntGirl[girls[i]]++
	}
	total := int64(k) * int64(k-1) / 2
	for i := 1; i <= a; i++ {
		x := cntBoy[i]
		total -= int64(x) * int64(x-1) / 2
	}
	for i := 1; i <= b; i++ {
		x := cntGirl[i]
		total -= int64(x) * int64(x-1) / 2
	}
	return fmt.Sprintf("%d", total)
}

func generateC(rng *rand.Rand) (string, string) {
	a := rng.Intn(10) + 1
	b := rng.Intn(10) + 1
	k := rng.Intn(a * b)
	if k < 2 {
		k = 2
	}
	boys := make([]int, k)
	girls := make([]int, k)
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d %d\n", a, b, k))
	for i := 0; i < k; i++ {
		boys[i] = rng.Intn(a) + 1
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", boys[i]))
	}
	sb.WriteByte('\n')
	for i := 0; i < k; i++ {
		girls[i] = rng.Intn(b) + 1
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", girls[i]))
	}
	sb.WriteByte('\n')
	out := solveC(a, b, k, boys, girls)
	return sb.String(), out
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
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateC(rng)
		got, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

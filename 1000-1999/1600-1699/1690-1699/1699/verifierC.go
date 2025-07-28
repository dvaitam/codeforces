package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

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

const mod int64 = 1000000007

func solveC(n int, arr []int) int64 {
	pos := make([]int, n)
	for i, v := range arr {
		pos[v] = i
	}
	L, R := pos[0], pos[0]
	ans := int64(1)
	for x := 1; x < n; x++ {
		p := pos[x]
		if p < L {
			L = p
		} else if p > R {
			R = p
		} else {
			choices := int64(R - L + 1 - x)
			ans = (ans * choices) % mod
		}
	}
	return ans
}

func generateCase(rng *rand.Rand) (string, int64) {
	n := rng.Intn(10) + 1
	perm := rng.Perm(n)
	input := fmt.Sprintf("1\n%d\n", n)
	for i, v := range perm {
		if i > 0 {
			input += " "
		}
		input += strconv.Itoa(v)
	}
	input += "\n"
	exp := solveC(n, perm)
	return input, exp
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		val, err := strconv.ParseInt(strings.TrimSpace(out), 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: bad output\ninput:\n%soutput:\n%s", i+1, in, out)
			os.Exit(1)
		}
		if val != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\ninput:\n%s", i+1, exp, val, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

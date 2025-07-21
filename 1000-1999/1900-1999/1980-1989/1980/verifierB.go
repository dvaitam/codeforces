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

func solveCase(n, f, k int, a []int) string {
	fav := a[f-1]
	gt, eq := 0, 0
	for _, v := range a {
		if v > fav {
			gt++
		} else if v == fav {
			eq++
		}
	}
	earliest := gt + 1
	latest := gt + eq
	if k < earliest {
		return "NO"
	} else if k >= latest {
		return "YES"
	}
	return "MAYBE"
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(100) + 1
	f := rng.Intn(n) + 1
	k := rng.Intn(n) + 1
	a := make([]int, n)
	var sb strings.Builder
	for i := range a {
		a[i] = rng.Intn(100) + 1
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(a[i]))
	}
	input := fmt.Sprintf("1\n%d %d %d\n%s\n", n, f, k, sb.String())
	expect := solveCase(n, f, k, a)
	return input, expect
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
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

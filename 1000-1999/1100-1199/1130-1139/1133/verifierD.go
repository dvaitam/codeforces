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

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func solveD(a, b []int64) int {
	counts := make(map[[2]int64]int)
	alwaysZero := 0
	for i := range a {
		ai, bi := a[i], b[i]
		if ai == 0 {
			if bi == 0 {
				alwaysZero++
			}
			continue
		}
		num := -bi
		den := ai
		if den < 0 {
			den = -den
			num = -num
		}
		if num == 0 {
			den = 1
		} else {
			g := gcd(num, den)
			num /= g
			den /= g
		}
		key := [2]int64{num, den}
		counts[key]++
	}
	maxCount := 0
	for _, v := range counts {
		if v > maxCount {
			maxCount = v
		}
	}
	return alwaysZero + maxCount
}

func generateCase(rng *rand.Rand) (string, int) {
	n := rng.Intn(7) + 1
	a := make([]int64, n)
	b := make([]int64, n)
	for i := 0; i < n; i++ {
		a[i] = int64(rng.Intn(21) - 10)
		b[i] = int64(rng.Intn(21) - 10)
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(v, 10))
	}
	sb.WriteByte('\n')
	for i, v := range b {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(v, 10))
	}
	sb.WriteByte('\n')
	return sb.String(), solveD(a, b)
}

func run(bin, input string) (int, error) {
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
		return 0, fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	val, err := strconv.Atoi(strings.TrimSpace(out.String()))
	if err != nil {
		return 0, fmt.Errorf("invalid integer output: %v", err)
	}
	return val, nil
}

func main() {
	if len(os.Args) != 2 && !(len(os.Args) == 3 && os.Args[1] == "--") {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[len(os.Args)-1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d mismatch: expected %d got %d\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

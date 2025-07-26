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

func solveB(n, k int, arr []int) int {
	buckets := make([]int, k)
	for _, v := range arr {
		r := v % k
		if r < 0 {
			r += k
		}
		buckets[r]++
	}
	ans := 0
	for i := 1; i <= k/2; i++ {
		if i == k-i {
			ans += (buckets[i] / 2) * 2
		} else {
			if buckets[i] < buckets[k-i] {
				ans += buckets[i] * 2
			} else {
				ans += buckets[k-i] * 2
			}
		}
	}
	ans += (buckets[0] / 2) * 2
	return ans
}

func generateCase(rng *rand.Rand) (string, int) {
	n := rng.Intn(8) + 2
	k := rng.Intn(9) + 2
	arr := make([]int, n)
	for i := range arr {
		arr[i] = rng.Intn(1000) - 500
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	return sb.String(), solveB(n, k, arr)
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
		fmt.Println("usage: go run verifierB.go /path/to/binary")
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

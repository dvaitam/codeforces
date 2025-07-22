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

const MOD = 1000000007

func solve(n int, k int, a []int) []int {
	res := append([]int(nil), a...)
	for iter := 0; iter < k; iter++ {
		for i := 1; i < n; i++ {
			res[i] = (res[i] + res[i-1]) % MOD
		}
	}
	return res
}

func genCase(rng *rand.Rand) (int, int, []int) {
	n := rng.Intn(5) + 1
	k := rng.Intn(5)
	arr := make([]int, n)
	for i := range arr {
		arr[i] = rng.Intn(1000)
	}
	return n, k, arr
}

func runCase(bin string, n, k int, arr []int, expected []int) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	input := sb.String()
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
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	fields := strings.Fields(out.String())
	if len(fields) != n {
		return fmt.Errorf("expected %d numbers got %d", n, len(fields))
	}
	for i := 0; i < n; i++ {
		var x int
		fmt.Sscan(fields[i], &x)
		if x != expected[i]%MOD {
			return fmt.Errorf("index %d expected %d got %d", i, expected[i]%MOD, x)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, k, arr := genCase(rng)
		exp := solve(n, k, arr)
		if err := runCase(bin, n, k, arr, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

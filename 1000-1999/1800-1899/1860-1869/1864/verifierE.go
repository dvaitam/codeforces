package main

import (
	"bytes"
	"fmt"
	"math/bits"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

const modE = 998244353

func msb(x int) int {
	if x == 0 {
		return -1
	}
	return bits.Len(uint(x)) - 1
}

func turns(a, b int) int {
	if a == 0 {
		if b == 0 {
			return 1
		}
		return 1
	}
	ha := msb(a)
	hb := msb(b)
	if ha < hb {
		return 1
	}
	if ha > hb {
		return 2
	}
	if a < b {
		return 3
	}
	if a == b {
		pc := bits.OnesCount(uint(a)) + 1
		if pc < 4 {
			return pc
		}
		return 4
	}
	thresh := 3 << (ha - 1)
	if hb == ha && b >= thresh {
		return 4
	}
	return 2
}

func modInv(a int) int {
	return modPow(a, modE-2)
}

func modPow(a, e int) int {
	res := 1
	for e > 0 {
		if e&1 == 1 {
			res = res * a % modE
		}
		a = a * a % modE
		e >>= 1
	}
	return res
}

func expectedE(arr []int) string {
	n := len(arr)
	sum := 0
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			sum += turns(arr[i], arr[j])
		}
	}
	inv := modInv(n * n % modE)
	ans := sum % modE * inv % modE
	return fmt.Sprintf("%d\n", ans)
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(4) + 1
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = rng.Intn(1 << 10)
	}
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", arr[i]))
	}
	sb.WriteByte('\n')
	expect := expectedE(arr)
	return sb.String(), expect
}

func runCase(bin, input, expected string) error {
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
	outStr := strings.TrimSpace(out.String())
	expStr := strings.TrimSpace(expected)
	if outStr != expStr {
		return fmt.Errorf("expected %q got %q", expStr, outStr)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

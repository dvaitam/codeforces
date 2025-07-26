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

const mod = 998244353

func expectedE(arr []int64) string {
	n := len(arr) - 1
	var ans, pow int64 = 0, 1
	for i := n - 1; i >= 1; i-- {
		term := arr[i] * int64(n-i+2) % mod * pow % mod
		ans += term
		if ans >= mod {
			ans -= mod
		}
		pow = pow * 2 % mod
	}
	ans += arr[n]
	ans %= mod
	return fmt.Sprint(ans)
}

func genCaseE(rng *rand.Rand) []int64 {
	n := rng.Intn(10) + 1
	arr := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		arr[i] = int64(rng.Intn(1000))
	}
	arr[0] = int64(n)
	return arr
}

func runCaseE(bin string, arr []int64) error {
	n := int(arr[0])
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(n))
	sb.WriteByte('\n')
	for i := 1; i <= n; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(arr[i], 10))
	}
	sb.WriteByte('\n')
	input := sb.String()
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	got := strings.TrimSpace(out.String())
	exp := expectedE(arr)
	if got != exp {
		return fmt.Errorf("expected %s got %s", exp, got)
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
		arr := genCaseE(rng)
		if err := runCaseE(bin, arr); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

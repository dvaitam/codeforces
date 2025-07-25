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

const mod int64 = 1000000007

func expectedA(xs []int64) int64 {
	n := len(xs)
	sort.Slice(xs, func(i, j int) bool { return xs[i] < xs[j] })
	pow2 := make([]int64, n+1)
	pow2[0] = 1
	for i := 1; i <= n; i++ {
		pow2[i] = (pow2[i-1] * 2) % mod
	}
	var ans int64
	for j := 1; j < n; j++ {
		ans = (ans + xs[j]*(pow2[j]-1)) % mod
	}
	for i := 0; i < n-1; i++ {
		ans = (ans - xs[i]*(pow2[n-1-i]-1)) % mod
	}
	if ans < 0 {
		ans += mod
	}
	return ans
}

func generateCaseA(rng *rand.Rand) []int64 {
	n := rng.Intn(8) + 2 // 2..9
	xs := make([]int64, n)
	for i := range xs {
		xs[i] = int64(rng.Intn(1000) + 1)
	}
	return xs
}

func runCase(bin string, xs []int64) error {
	var input strings.Builder
	input.WriteString(fmt.Sprintf("%d\n", len(xs)))
	for i, v := range xs {
		if i > 0 {
			input.WriteByte(' ')
		}
		input.WriteString(fmt.Sprint(v))
	}
	input.WriteByte('\n')
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int64
	if _, err := fmt.Fscan(bytes.NewReader(out.Bytes()), &got); err != nil {
		return fmt.Errorf("invalid output: %v", err)
	}
	expect := expectedA(append([]int64(nil), xs...))
	if got != expect {
		return fmt.Errorf("expected %d got %d", expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		xs := generateCaseA(rng)
		if err := runCase(bin, xs); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%v\n", i+1, err, xs)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

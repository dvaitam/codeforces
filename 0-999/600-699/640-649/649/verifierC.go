package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(8) + 1
	x := rng.Intn(20)
	y := rng.Intn(20)
	arr := make([]int, n)
	for i := range arr {
		arr[i] = rng.Intn(10) + 1
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, x, y))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')

	b := append([]int(nil), arr...)
	sort.Ints(b)
	prefixDouble := make([]int64, n+1)
	prefixOdd := make([]int64, n+1)
	for i := 0; i < n; i++ {
		prefixDouble[i+1] = prefixDouble[i] + int64((b[i]+1)/2)
		if b[i]%2 == 1 {
			prefixOdd[i+1] = prefixOdd[i] + 1
		} else {
			prefixOdd[i+1] = prefixOdd[i]
		}
	}
	feasible := func(k int) bool {
		totalDouble := prefixDouble[k]
		need := totalDouble - int64(x)
		if need <= 0 {
			return true
		}
		odd := prefixOdd[k]
		use1 := odd
		if use1 > int64(y) {
			use1 = int64(y)
		}
		if use1 > need {
			use1 = need
		}
		need -= use1
		yLeft := int64(y) - use1
		if yLeft >= need*2 {
			return true
		}
		return false
	}
	l, r := 0, n
	for l < r {
		mid := (l + r + 1) / 2
		if feasible(mid) {
			l = mid
		} else {
			r = mid - 1
		}
	}
	exp := fmt.Sprintf("%d\n", l)
	return sb.String(), exp
}

func runCase(exe, input, expected string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(expected)
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(exe, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

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

func solve(n, k, q int, arr []int) int64 {
	var ans int64
	cur := 0
	for i := 0; i < n; i++ {
		if arr[i] <= q {
			cur++
		} else {
			if cur >= k {
				l := int64(cur - k + 1)
				ans += l * (l + 1) / 2
			}
			cur = 0
		}
	}
	if cur >= k {
		l := int64(cur - k + 1)
		ans += l * (l + 1) / 2
	}
	return ans
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(20) + 1
	k := rng.Intn(n) + 1
	qv := rng.Intn(21) - 10
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = rng.Intn(41) - 20
	}
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, k, qv))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", arr[i]))
	}
	sb.WriteByte('\n')
	expected := fmt.Sprintf("%d", solve(n, k, qv, arr))
	return sb.String(), expected
}

func runCase(bin, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
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
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

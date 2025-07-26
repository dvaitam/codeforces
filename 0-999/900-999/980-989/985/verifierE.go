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

func expectedE(n, k, d int, arr []int) string {
	b := append([]int(nil), arr...)
	sort.Ints(b)
	dp := make([]bool, n+1)
	pref := make([]int, n+1)
	dp[0] = true
	pref[0] = 1
	l := 0
	for i := 1; i <= n; i++ {
		for l < i && b[i-1]-b[l] > d {
			l++
		}
		if i >= k {
			left := l
			right := i - k
			if left <= right {
				sum := pref[right]
				if left > 0 {
					sum -= pref[left-1]
				}
				if sum > 0 {
					dp[i] = true
				}
			}
		}
		pref[i] = pref[i-1]
		if dp[i] {
			pref[i]++
		}
	}
	if dp[n] {
		return "YES"
	}
	return "NO"
}

func generateCaseE(rng *rand.Rand) (int, int, int, []int) {
	n := rng.Intn(20) + 1
	k := rng.Intn(n) + 1
	d := rng.Intn(10)
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = rng.Intn(50)
	}
	return n, k, d, arr
}

func runCaseE(bin string, n, k, d int, arr []int) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, k, d))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	expected := expectedE(n, k, d, arr)
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
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
		n, k, d, arr := generateCaseE(rng)
		if err := runCaseE(bin, n, k, d, arr); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

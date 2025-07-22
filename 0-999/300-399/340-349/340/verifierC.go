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

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func expectedAnswer(a []int64) (int64, int64) {
	n := len(a)
	sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
	var sumA int64
	for _, v := range a {
		sumA += v
	}
	var prefix int64
	var sumPairs int64
	for i, v := range a {
		sumPairs += v*int64(i) - prefix
		prefix += v
	}
	num := sumA + 2*sumPairs
	den := int64(n)
	g := gcd(num, den)
	return num / g, den / g
}

func generateCase(rng *rand.Rand) []int64 {
	n := rng.Intn(50) + 1
	vals := make([]int64, 0, n)
	used := make(map[int64]bool)
	for len(vals) < n {
		v := int64(rng.Intn(2000))
		if !used[v] {
			used[v] = true
			vals = append(vals, v)
		}
	}
	return vals
}

func runCase(bin string, arr []int64) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(arr)))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(v, 10))
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
	num, den := expectedAnswer(append([]int64(nil), arr...))
	expected := fmt.Sprintf("%d %d", num, den)
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
		arr := generateCase(rng)
		if err := runCase(bin, arr); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

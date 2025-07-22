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

func expected(pies []int) int64 {
	sort.Ints(pies)
	n := len(pies)
	k := 0
	i, j := 0, 0
	for j < n && i < n {
		if pies[j] > pies[i] {
			k++
			i++
			j++
		} else {
			j++
		}
	}
	if k > n-k {
		k = n - k
	}
	var total, free int64
	for _, v := range pies {
		total += int64(v)
	}
	for idx := 0; idx < k; idx++ {
		free += int64(pies[idx])
	}
	return total - free
}

func runCase(bin string, pies []int) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(pies)))
	for i, p := range pies {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", p))
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
	var val int64
	if _, err := fmt.Sscan(strings.TrimSpace(out.String()), &val); err != nil {
		return fmt.Errorf("cannot parse output: %v", err)
	}
	exp := expected(append([]int(nil), pies...))
	if val != exp {
		return fmt.Errorf("expected %d got %d", exp, val)
	}
	return nil
}

func generateCase(rng *rand.Rand) []int {
	n := rng.Intn(10) + 1
	pies := make([]int, n)
	for i := range pies {
		pies[i] = rng.Intn(1000) + 1
	}
	return pies
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 100; i++ {
		pies := generateCase(rng)
		if err := runCase(bin, pies); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

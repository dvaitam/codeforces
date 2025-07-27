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

func run(bin string, input string) (string, error) {
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
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solveE(n, k int, a []int) int64 {
	sort.Slice(a, func(i, j int) bool { return a[i] > a[j] })
	cntNeg := 0
	for _, v := range a {
		if v < 0 {
			cntNeg++
		}
	}
	p := k
	if p > cntNeg {
		p = cntNeg
	}
	mainSize := n - p
	var ans int64
	for i := 0; i < mainSize; i++ {
		mult := mainSize - i - 1
		ans += int64(a[i]) * int64(mult)
	}
	return ans
}

func runCase(bin string, n, k int, a []int) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", a[i]))
	}
	sb.WriteByte('\n')
	expect := fmt.Sprintf("%d", solveE(n, k, append([]int(nil), a...)))
	out, err := run(bin, sb.String())
	if err != nil {
		return err
	}
	if out != expect {
		return fmt.Errorf("expected %s got %s", expect, out)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	total := 0
	// edge cases
	if err := runCase(bin, 1, 0, []int{5}); err != nil {
		fmt.Fprintf(os.Stderr, "case %d failed: %v\n", total+1, err)
		os.Exit(1)
	}
	total++
	if err := runCase(bin, 3, 1, []int{-1, 2, -3}); err != nil {
		fmt.Fprintf(os.Stderr, "case %d failed: %v\n", total+1, err)
		os.Exit(1)
	}
	total++
	for total < 100 {
		n := rng.Intn(20) + 1
		k := rng.Intn(n + 1)
		a := make([]int, n)
		for i := range a {
			a[i] = rng.Intn(21) - 10
		}
		if err := runCase(bin, n, k, a); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", total+1, err)
			os.Exit(1)
		}
		total++
	}
	fmt.Printf("All %d tests passed\n", total)
}

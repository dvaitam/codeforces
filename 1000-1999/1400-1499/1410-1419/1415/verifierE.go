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

func solveE(n, k int, c []int) int64 {
	sort.Ints(c)

	s := make([]int64, n+1)
	s[0] = 0
	for i := 0; i < n; i++ {
		s[i+1] = s[i] + int64(c[i])
	}

	sn := s[n]
	w := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		w[i] = sn - s[i]
	}

	suffixSum := make([]int64, n+2)
	suffixSum[n+1] = 0
	for i := n; i >= 1; i-- {
		suffixSum[i] = suffixSum[i+1] + w[i]
	}

	r := k + 1
	preCalc := make([]int64, n/r+2)
	preCalc[0] = 0
	for q := 1; q*r <= n; q++ {
		preCalc[q] = preCalc[q-1] + w[q*r]
	}

	var maxScore int64 = -9223372036854775808

	for i := 0; i <= n; i++ {
		var currentScore int64
		if i > 0 {
			q := i / r
			currentScore = preCalc[q]
			if i%r != 0 {
				currentScore += w[i]
			}
		}
		currentScore += suffixSum[i+1]
		if currentScore > maxScore {
			maxScore = currentScore
		}
	}
	return maxScore
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

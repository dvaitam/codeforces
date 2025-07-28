package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func lexGreater(a, b []int) bool {
	for i := 0; i < len(a) && i < len(b); i++ {
		if a[i] != b[i] {
			return a[i] > b[i]
		}
	}
	return false
}

func solveCase(p []int) []int {
	n := len(p)
	if n == 1 {
		return p
	}
	posMax := 1
	for i := 1; i < n; i++ {
		if p[i] > p[posMax] {
			posMax = i
		}
	}
	posPrevMax := 0
	for i := 0; i < n-1; i++ {
		if p[i] > p[posPrevMax] {
			posPrevMax = i
		}
	}
	rsetMap := map[int]bool{}
	cand := []int{posMax - 1, posMax, posPrevMax - 1, posPrevMax, n - 1}
	for _, r := range cand {
		if r >= 0 && r < n {
			rsetMap[r] = true
		}
	}
	best := []int{}
	for r := range rsetMap {
		suffix := make([]int, n-r-1)
		copy(suffix, p[r+1:])
		for l := 0; l <= r; l++ {
			candperm := append([]int{}, suffix...)
			for i := r; i >= l; i-- {
				candperm = append(candperm, p[i])
			}
			candperm = append(candperm, p[:l]...)
			if len(best) == 0 || lexGreater(candperm, best) {
				best = candperm
			}
		}
	}
	return best
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 1
	perm := rng.Perm(n)
	for i := 0; i < n; i++ {
		perm[i]++
	}
	input := fmt.Sprintf("1\n%d\n", n)
	for i, v := range perm {
		if i > 0 {
			input += " "
		}
		input += fmt.Sprintf("%d", v)
	}
	input += "\n"
	ans := solveCase(perm)
	var sb strings.Builder
	for i, v := range ans {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return input, sb.String()
}

func runCase(bin, input, exp string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	if strings.TrimSpace(out.String()) != strings.TrimSpace(exp) {
		return fmt.Errorf("expected %q got %q", strings.TrimSpace(exp), strings.TrimSpace(out.String()))
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(42))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

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

func visible(h []int, i, j int) bool {
	n := len(h)
	minH := h[i]
	if h[j] < minH {
		minH = h[j]
	}
	// clockwise i->j
	k := (i + 1) % n
	ok := true
	for k != j {
		if h[k] > minH {
			ok = false
			break
		}
		k = (k + 1) % n
	}
	if ok {
		return true
	}
	// counter-clockwise
	k = (i - 1 + n) % n
	ok = true
	for k != j {
		if h[k] > minH {
			ok = false
			break
		}
		k = (k - 1 + n) % n
	}
	return ok
}

func countPairs(h []int) int {
	n := len(h)
	res := 0
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if visible(h, i, j) {
				res++
			}
		}
	}
	return res
}

func generateCase(rng *rand.Rand) (string, int) {
	n := rng.Intn(6) + 3
	h := make([]int, n)
	for i := range h {
		h[i] = rng.Intn(10) + 1
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range h {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return sb.String(), countPairs(h)
}

func runCase(bin, input string, expected int) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int
	if _, err := fmt.Sscan(strings.TrimSpace(out.String()), &got); err != nil {
		return fmt.Errorf("failed to parse output: %v", err)
	}
	if got != expected {
		return fmt.Errorf("expected %d got %d", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

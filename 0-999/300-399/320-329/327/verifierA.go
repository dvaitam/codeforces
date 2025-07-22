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

func expected(arr []int) int {
	ones := 0
	b := make([]int, len(arr))
	for i, x := range arr {
		if x == 1 {
			ones++
			b[i] = -1
		} else {
			b[i] = 1
		}
	}
	best := b[0]
	curr := b[0]
	for i := 1; i < len(b); i++ {
		if curr < 0 {
			curr = b[i]
		} else {
			curr += b[i]
		}
		if curr > best {
			best = curr
		}
	}
	return ones + best
}

func generateCase(rng *rand.Rand) []int {
	n := rng.Intn(10) + 1
	arr := make([]int, n)
	for i := range arr {
		arr[i] = rng.Intn(2)
	}
	return arr
}

func runCase(bin string, arr []int) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(arr)))
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
	expect := fmt.Sprint(expected(arr))
	if got != expect {
		return fmt.Errorf("expected %s got %s", expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		arr := generateCase(rng)
		if err := runCase(bin, arr); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput: %v\n", i+1, err, arr)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

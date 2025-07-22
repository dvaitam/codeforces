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

func expectedAnswer(k int, a []int) int {
	if k == 0 {
		return 0
	}
	sort.Slice(a, func(i, j int) bool { return a[i] > a[j] })
	sum := 0
	for i, v := range a {
		sum += v
		if sum >= k {
			return i + 1
		}
	}
	return -1
}

func generateCase(rng *rand.Rand) (int, []int) {
	k := rng.Intn(101)
	a := make([]int, 12)
	for i := range a {
		a[i] = rng.Intn(101)
	}
	return k, a
}

func runCase(bin string, k int, a []int) error {
	var input strings.Builder
	input.WriteString(fmt.Sprintf("%d\n", k))
	for i, v := range a {
		input.WriteString(fmt.Sprintf("%d", v))
		if i+1 < len(a) {
			input.WriteByte(' ')
		}
	}
	input.WriteByte('\n')
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	expect := expectedAnswer(k, append([]int(nil), a...))
	if got != expect {
		return fmt.Errorf("expected %d got %d", expect, got)
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

	// edge cases
	for _, k := range []int{0, 1, 50, 100} {
		a := make([]int, 12)
		for i := range a {
			a[i] = 0
		}
		if err := runCase(bin, k, a); err != nil {
			fmt.Fprintf(os.Stderr, "edge case k=%d failed: %v\n", k, err)
			os.Exit(1)
		}
	}

	for i := 0; i < 100; i++ {
		k, a := generateCase(rng)
		if err := runCase(bin, k, a); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput: k=%d a=%v\n", i+1, err, k, a)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

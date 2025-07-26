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

func solve(weights []int) string {
	sort.Ints(weights)
	const maxW = 150001
	used := make(map[int]bool)
	ans := 0
	for _, w := range weights {
		if w-1 > 0 && !used[w-1] {
			used[w-1] = true
			ans++
		} else if !used[w] {
			used[w] = true
			ans++
		} else if !used[w+1] {
			used[w+1] = true
			ans++
		}
	}
	return fmt.Sprintf("%d", ans)
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(20) + 1
	w := make([]int, n)
	for i := range w {
		w[i] = rng.Intn(100) + 1
	}
	parts := make([]string, n)
	for i, v := range w {
		parts[i] = fmt.Sprintf("%d", v)
	}
	input := fmt.Sprintf("%d\n%s\n", n, strings.Join(parts, " "))
	expect := solve(w)
	return input, expect
}

func runCase(bin, input, expect string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(strings.Split(out.String(), "\n")[0])
	if got != expect {
		return fmt.Errorf("expected %s got %s", expect, got)
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
		input, expect := generateCase(rng)
		if err := runCase(bin, input, expect); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

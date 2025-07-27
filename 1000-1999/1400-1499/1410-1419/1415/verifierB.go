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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func solveB(n, k int, c []int) int {
	present := make(map[int]struct{})
	for _, v := range c {
		present[v] = struct{}{}
	}
	ans := n + 1
	for color := range present {
		days := 0
		for i := 0; i < n; {
			if c[i] == color {
				i++
			} else {
				days++
				i += k
			}
		}
		if days < ans {
			ans = days
		}
	}
	return ans
}

func runCase(bin string, n, k int, c []int) error {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", c[i]))
	}
	sb.WriteByte('\n')
	expect := fmt.Sprintf("%d", solveB(n, k, c))
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
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	total := 0
	// some edge cases
	edgeCases := []struct {
		n, k int
		c    []int
	}{
		{1, 1, []int{1}},
		{3, 2, []int{1, 1, 1}},
		{5, 1, []int{1, 2, 3, 4, 5}},
		{4, 4, []int{2, 2, 2, 2}},
	}
	for _, ec := range edgeCases {
		if err := runCase(bin, ec.n, ec.k, ec.c); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", total+1, err)
			os.Exit(1)
		}
		total++
	}
	for total < 100 {
		n := rng.Intn(20) + 1
		k := rng.Intn(n) + 1
		c := make([]int, n)
		for i := 0; i < n; i++ {
			c[i] = rng.Intn(5) + 1
		}
		if err := runCase(bin, n, k, c); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", total+1, err)
			os.Exit(1)
		}
		total++
	}
	fmt.Printf("All %d tests passed\n", total)
}

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

func run(bin, input string) (string, error) {
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

func expected(c []int, X int) int {
	n := len(c)
	w := make([]int, n)
	for i := 0; i < n; i++ {
		w[i] = c[i] * (n - i)
	}
	sort.Ints(w)
	cnt, sum := 0, 0
	for _, wi := range w {
		if sum+wi <= X {
			sum += wi
			cnt++
		} else {
			break
		}
	}
	return cnt
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(6) + 1
	X := rng.Intn(20) + 1
	c := make([]int, n)
	for i := 0; i < n; i++ {
		c[i] = rng.Intn(10) + 1
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, X)
	for i, v := range c {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	input := sb.String()
	ans := expected(c, X)
	return input, fmt.Sprintf("%d", ans)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var cases [][2]string
	cases = append(cases, [2]string{"1 1\n1\n", "1"})
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		cases = append(cases, [2]string{in, exp})
	}
	for i, tc := range cases {
		out, err := run(bin, tc[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc[0])
			os.Exit(1)
		}
		if strings.TrimSpace(out) != tc[1] {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, tc[1], out, tc[0])
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

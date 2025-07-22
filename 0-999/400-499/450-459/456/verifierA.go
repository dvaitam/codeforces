package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type laptop struct {
	price   int
	quality int
}

func solveA(r io.Reader) string {
	in := bufio.NewReader(r)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return ""
	}
	laptops := make([]laptop, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &laptops[i].price, &laptops[i].quality)
	}
	sort.Slice(laptops, func(i, j int) bool { return laptops[i].price < laptops[j].price })
	for i := 1; i < n; i++ {
		if laptops[i-1].quality > laptops[i].quality {
			return "Happy Alex\n"
		}
	}
	return "Poor Alex\n"
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(50) + 1
	prices := rng.Perm(n)
	qualities := rng.Perm(n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", prices[i]+1, qualities[i]+1))
	}
	input := sb.String()
	expected := solveA(strings.NewReader(input))
	return input, expected
}

func runCase(bin, input, expected string) error {
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
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(expected)
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
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
	cases := []struct{ in, exp string }{
		{"1\n1 1\n", "Poor Alex\n"},
		{"2\n1 1\n2 2\n", "Poor Alex\n"},
		{"2\n1 2\n2 1\n", "Happy Alex\n"},
	}
	for i := len(cases); i < 100; i++ {
		in, exp := generateCase(rng)
		cases = append(cases, struct{ in, exp string }{in, exp})
	}
	for i, c := range cases {
		if err := runCase(bin, c.in, c.exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, c.in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

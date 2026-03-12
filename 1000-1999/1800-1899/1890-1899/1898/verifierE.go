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

func expectedE(s, t string) string {
	if len(t) == 0 {
		return "YES"
	}
	var queues [26][]int
	for i := 0; i < len(s); i++ {
		c := int(s[i] - 'a')
		queues[c] = append(queues[c], i)
	}
	var front [26]int
	for i := 0; i < len(t); i++ {
		c := int(t[i] - 'a')
		if front[c] >= len(queues[c]) {
			return "NO"
		}
		pos := queues[c][front[c]]
		front[c]++
		for d := 0; d < c; d++ {
			for front[d] < len(queues[d]) && queues[d][front[d]] <= pos {
				front[d]++
			}
		}
	}
	return "YES"
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(20) + 1
	m := rng.Intn(n) + 1
	var sb strings.Builder
	for i := 0; i < n; i++ {
		sb.WriteByte(byte('a' + rng.Intn(26)))
	}
	s := sb.String()
	var tb strings.Builder
	for i := 0; i < m; i++ {
		tb.WriteByte(byte('a' + rng.Intn(26)))
	}
	t := tb.String()
	input := fmt.Sprintf("1\n%d %d\n%s\n%s\n", n, m, s, t)
	expect := expectedE(s, t)
	return input, expect
}

func runCase(bin, input, exp string) error {
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
	if !strings.EqualFold(got, exp) {
		return fmt.Errorf("expected %s got %s", exp, got)
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
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

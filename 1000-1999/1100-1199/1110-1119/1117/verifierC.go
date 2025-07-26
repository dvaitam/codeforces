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

func expectedC(x1, y1, x2, y2 int, s string) int {
	dx := x2 - x1
	dy := y2 - y1
	for t := 0; t <= 2000; t++ {
		if dx == 0 && dy == 0 {
			return t
		}
		c := s[t%len(s)]
		switch c {
		case 'U':
			dy--
		case 'D':
			dy++
		case 'L':
			dx++
		case 'R':
			dx--
		}
	}
	return -1
}

func generateCase(rng *rand.Rand) (string, int) {
	x1 := rng.Intn(21) - 10
	y1 := rng.Intn(21) - 10
	x2 := rng.Intn(21) - 10
	y2 := rng.Intn(21) - 10
	n := rng.Intn(10) + 1
	var b strings.Builder
	fmt.Fprintf(&b, "%d %d %d %d %d\n", x1, y1, x2, y2, n)
	var s strings.Builder
	moves := []byte{'U', 'D', 'L', 'R'}
	for i := 0; i < n; i++ {
		s.WriteByte(moves[rng.Intn(4)])
	}
	b.WriteString(s.String())
	b.WriteByte('\n')
	return b.String(), expectedC(x1, y1, x2, y2, s.String())
}

func runCase(bin, input string, exp int) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	var got int
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("failed to parse output: %v", err)
	}
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

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

func solveB(a1, a2, a3, a4 int) string {
	var start, end byte
	switch {
	case a3 == a4:
		start, end = '4', '4'
	case a3 == a4+1:
		start, end = '4', '7'
	case a4 == a3+1:
		start, end = '7', '4'
	default:
		return "-1"
	}
	var seg4, seg7, k int
	switch {
	case start == '4' && end == '4':
		k = a3
		seg4, seg7 = k+1, k
	case start == '4' && end == '7':
		k = a4
		seg4, seg7 = k+1, k+1
	case start == '7' && end == '4':
		k = a3
		seg4, seg7 = k+1, k+1
	case start == '7' && end == '7':
		k = a3
		seg4, seg7 = k, k+1
	}
	if a1 < seg4 || a2 < seg7 {
		return "-1"
	}
	size4First := a1 - (seg4 - 1)
	size7Last := a2 - (seg7 - 1)
	totalLen := a1 + a2
	var b strings.Builder
	b.Grow(totalLen)
	idx4, idx7 := 0, 0
	current := start
	for i := 0; i < seg4+seg7; i++ {
		if current == '4' {
			idx4++
			cnt := 1
			if idx4 == 1 {
				cnt = size4First
			}
			b.WriteString(strings.Repeat("4", cnt))
			current = '7'
		} else {
			idx7++
			cnt := 1
			if idx7 == seg7 {
				cnt = size7Last
			}
			b.WriteString(strings.Repeat("7", cnt))
			current = '4'
		}
	}
	return b.String()
}

func run(bin, input string) (string, error) {
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
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func generateCase(rng *rand.Rand) (string, string) {
	a1 := rng.Intn(6)
	a2 := rng.Intn(6)
	a3 := rng.Intn(6)
	a4 := rng.Intn(6)
	input := fmt.Sprintf("%d %d %d %d\n", a1, a2, a3, a4)
	return input, solveB(a1, a2, a3, a4)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

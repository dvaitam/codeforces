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

func solve(a, b string) int {
	cnt47, cnt74 := 0, 0
	for i := range a {
		if a[i] == '4' && b[i] == '7' {
			cnt47++
		} else if a[i] == '7' && b[i] == '4' {
			cnt74++
		}
	}
	if cnt47 > cnt74 {
		return cnt47
	}
	return cnt74
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

func generateCase(rng *rand.Rand) (string, int) {
	n := rng.Intn(50) + 1
	var a, b strings.Builder
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			a.WriteByte('4')
		} else {
			a.WriteByte('7')
		}
		if rng.Intn(2) == 0 {
			b.WriteByte('4')
		} else {
			b.WriteByte('7')
		}
	}
	input := fmt.Sprintf("%s\n%s\n", a.String(), b.String())
	return input, solve(a.String(), b.String())
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
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
		var got int
		fmt.Sscan(out, &got)
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

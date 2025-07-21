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

func expected(a, b, c, d, e, f int64) string {
	if c == 0 && d > 0 {
		return "Ron"
	}
	if a == 0 && b > 0 && ((c > 0 && d > 0) || (c == 0 && d > 0)) {
		return "Ron"
	}
	if e == 0 && f > 0 && (a > 0 && b > 0) && ((c > 0 && d > 0) || (c == 0 && d > 0)) {
		return "Ron"
	}
	if a > 0 && b > 0 && c > 0 && d > 0 && e > 0 && f > 0 {
		if b*d*f > a*c*e {
			return "Ron"
		}
	}
	return "Hermione"
}

func genCase(rng *rand.Rand) (string, string) {
	a := rng.Int63n(1001)
	b := rng.Int63n(1001)
	c := rng.Int63n(1001)
	d := rng.Int63n(1001)
	e := rng.Int63n(1001)
	f := rng.Int63n(1001)
	input := fmt.Sprintf("%d %d %d %d %d %d\n", a, b, c, d, e, f)
	exp := expected(a, b, c, d, e, f)
	return input, exp
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

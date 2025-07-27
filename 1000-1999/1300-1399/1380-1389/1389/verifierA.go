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

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func lcm(a, b int) int {
	return a / gcd(a, b) * b
}

func solveCase(l, r int) (int, int) {
	for x := l; x <= r; x++ {
		for y := x + 1; y <= r; y++ {
			v := lcm(x, y)
			if v >= l && v <= r {
				return x, y
			}
		}
	}
	return -1, -1
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

func genTest(rng *rand.Rand) (string, string) {
	l := rng.Intn(20) + 1
	r := l + rng.Intn(20) + 1
	input := fmt.Sprintf("1\n%d %d\n", l, r)
	x, y := solveCase(l, r)
	expected := fmt.Sprintf("%d %d", x, y)
	return input, expected
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= 100; i++ {
		input, exp := genTest(rng)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", i, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "test %d failed\ninput:\n%sexpected:%s\ngot:%s\n", i, input, exp, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

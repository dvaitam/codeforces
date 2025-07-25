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

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func expected(x1, y1, x2, y2, x, y int) string {
	dx := abs(x2 - x1)
	dy := abs(y2 - y1)
	if dx%x != 0 || dy%y != 0 {
		return "NO"
	}
	kx := dx / x
	ky := dy / y
	if kx%2 == ky%2 {
		return "YES"
	}
	return "NO"
}

func runBin(bin, input string) (string, error) {
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
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func genCase(rng *rand.Rand) (string, string) {
	x1 := rng.Intn(200001) - 100000
	y1 := rng.Intn(200001) - 100000
	x2 := rng.Intn(200001) - 100000
	y2 := rng.Intn(200001) - 100000
	x := rng.Intn(100000) + 1
	y := rng.Intn(100000) + 1
	input := fmt.Sprintf("%d %d %d %d\n%d %d\n", x1, y1, x2, y2, x, y)
	exp := expected(x1, y1, x2, y2, x, y)
	return input, exp
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, exp := genCase(rng)
		out, err := runBin(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\ninput:\n%s\n", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed:\ninput:\n%s\nexpected:%s\ngot:%s\n", i+1, input, exp, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

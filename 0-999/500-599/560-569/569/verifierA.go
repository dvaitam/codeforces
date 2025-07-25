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

func expected(T, S, q int) string {
	cur := S
	count := 0
	for cur < T {
		cur *= q
		count++
	}
	return fmt.Sprint(count)
}

func buildInput(T, S, q int) string {
	return fmt.Sprintf("%d %d %d\n", T, S, q)
}

func randomCase(rng *rand.Rand) (int, int, int) {
	T := rng.Intn(100000-2) + 2
	S := rng.Intn(T-1) + 1
	q := rng.Intn(10000-2) + 2
	return T, S, q
}

func runCase(bin string, T, S, q int) error {
	input := buildInput(T, S, q)
	out, err := run(bin, input)
	if err != nil {
		return err
	}
	exp := expected(T, S, q)
	if out != exp {
		return fmt.Errorf("expected %s got %s", exp, out)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	cases := [][3]int{
		{5, 1, 2},
		{5, 4, 2},
		{100000, 1, 10000},
		{100000, 99999, 2},
		{7, 4, 3},
	}
	for len(cases) < 105 {
		T, S, q := randomCase(rng)
		cases = append(cases, [3]int{T, S, q})
	}

	for i, tc := range cases {
		if err := runCase(bin, tc[0], tc[1], tc[2]); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, buildInput(tc[0], tc[1], tc[2]))
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

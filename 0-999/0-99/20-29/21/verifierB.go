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

func intersectionSize(A1, B1, C1, A2, B2, C2 int) int {
	state := func(A, B, C int) int {
		if A == 0 && B == 0 {
			if C == 0 {
				return 1
			}
			return 0
		}
		return 2
	}
	s1 := state(A1, B1, C1)
	s2 := state(A2, B2, C2)
	if s1 == 0 || s2 == 0 {
		return 0
	}
	if s1 != 2 || s2 != 2 {
		return -1
	}
	det := A1*B2 - A2*B1
	if det != 0 {
		return 1
	}
	if A1*C2 == A2*C1 && B1*C2 == B2*C1 {
		return -1
	}
	return 0
}

func generateCase(rng *rand.Rand) (string, string) {
	A1 := rng.Intn(201) - 100
	B1 := rng.Intn(201) - 100
	C1 := rng.Intn(201) - 100
	A2 := rng.Intn(201) - 100
	B2 := rng.Intn(201) - 100
	C2 := rng.Intn(201) - 100
	expected := fmt.Sprintf("%d", intersectionSize(A1, B1, C1, A2, B2, C2))
	input := fmt.Sprintf("%d %d %d\n%d %d %d\n", A1, B1, C1, A2, B2, C2)
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
	outStr := strings.TrimSpace(out.String())
	if outStr != expected {
		return fmt.Errorf("expected %s got %s", expected, outStr)
	}
	return nil
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
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

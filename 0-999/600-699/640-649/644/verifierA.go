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

type Case struct {
	n, a, b int
	expect  string
}

func solveA(n, a, b int) string {
	if n > a*b {
		return "-1\n"
	}
	M := make([][]int, a)
	for i := 0; i < a; i++ {
		M[i] = make([]int, b)
	}
	cnt := 1
	for i := 0; i < a; i++ {
		for j := 0; j < b; j++ {
			M[i][j] = cnt
			if b%2 == 0 && i%2 == 1 {
				if j%2 == 0 {
					M[i][j]++
				} else {
					M[i][j]--
				}
			}
			if M[i][j] > n {
				M[i][j] = 0
			}
			cnt++
			if cnt > n+1 {
				break
			}
		}
		if cnt > n+1 {
			break
		}
	}
	var sb strings.Builder
	for i := 0; i < a; i++ {
		for j := 0; j < b; j++ {
			if j == b-1 {
				fmt.Fprintf(&sb, "%d\n", M[i][j])
			} else {
				fmt.Fprintf(&sb, "%d ", M[i][j])
			}
		}
	}
	return sb.String()
}

func genCase(rng *rand.Rand) (string, string) {
	a := rng.Intn(20) + 1
	b := rng.Intn(20) + 1
	n := rng.Intn(2*a*b) + 1
	if n > 10000 {
		n = 10000
	}
	input := fmt.Sprintf("%d %d %d\n", n, a, b)
	expected := solveA(n, a, b)
	return input, expected
}

func runCase(bin, input, expected string) error {
	cmd := exec.Command(bin)
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
		return fmt.Errorf("expected:\n%s\ngot:\n%s", exp, got)
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
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

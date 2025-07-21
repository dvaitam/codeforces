package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func solveEquation(a, b, c float64) (int, []float64) {
	if a == 0 && b == 0 && c != 0 {
		return 0, nil
	}
	if a == 0 && b == 0 && c == 0 {
		return -1, nil
	}
	if a == 0 {
		return 1, []float64{-c / b}
	}
	disc := b*b - 4*a*c
	if disc < 0 {
		return 0, nil
	}
	if disc == 0 {
		return 1, []float64{-b / (2 * a)}
	}
	d := math.Sqrt(disc)
	r1 := (-b - d) / (2 * a)
	r2 := (-b + d) / (2 * a)
	if r1 > r2 {
		r1, r2 = r2, r1
	}
	return 2, []float64{r1, r2}
}

func generateCase(rng *rand.Rand) (string, int, []float64) {
	a := float64(rng.Intn(200001) - 100000)
	b := float64(rng.Intn(200001) - 100000)
	c := float64(rng.Intn(200001) - 100000)
	input := fmt.Sprintf("%d %d %d\n", int(a), int(b), int(c))
	cnt, roots := solveEquation(a, b, c)
	return input, cnt, roots
}

func runCase(bin, input string, expCnt int, expRoots []float64) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	tokens := strings.Fields(strings.TrimSpace(out.String()))
	if expCnt == -1 {
		if len(tokens) != 1 || tokens[0] != "-1" {
			return fmt.Errorf("expected -1 got %q", strings.Join(tokens, " "))
		}
		return nil
	}
	if expCnt == 0 {
		if len(tokens) != 1 || tokens[0] != "0" {
			return fmt.Errorf("expected 0 got %q", strings.Join(tokens, " "))
		}
		return nil
	}
	if len(tokens) != expCnt+1 {
		return fmt.Errorf("expected %d values got %d", expCnt+1, len(tokens))
	}
	if tokens[0] != strconv.Itoa(expCnt) {
		return fmt.Errorf("expected count %d got %s", expCnt, tokens[0])
	}
	const eps = 1e-4
	for i := 0; i < expCnt; i++ {
		got, err := strconv.ParseFloat(tokens[i+1], 64)
		if err != nil {
			return fmt.Errorf("failed to parse float '%s'", tokens[i+1])
		}
		if math.Abs(got-expRoots[i]) > eps {
			return fmt.Errorf("root %d: expected %.5f got %.5f", i+1, expRoots[i], got)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, cnt, roots := generateCase(rng)
		if err := runCase(bin, in, cnt, roots); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

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

func buildOracle() (string, error) {
	oracle := "oracleG.bin"
	cmd := exec.Command("go", "build", "-o", oracle, "1297G.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build oracle: %v\n%s", err, out)
	}
	return oracle, nil
}

func runProg(prog, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(prog, ".go") {
		cmd = exec.Command("go", "run", prog)
	} else {
		cmd = exec.Command(prog)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("%v", err)
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(rng *rand.Rand) string {
	// m must be product of 2,3,5,7 powers
	exp2 := rng.Intn(4)
	exp3 := rng.Intn(4)
	exp5 := rng.Intn(4)
	exp7 := rng.Intn(4)
	m := int64(1)
	for i := 0; i < exp2; i++ {
		m *= 2
	}
	for i := 0; i < exp3; i++ {
		m *= 3
	}
	for i := 0; i < exp5; i++ {
		m *= 5
	}
	for i := 0; i < exp7; i++ {
		m *= 7
	}
	if m == 0 {
		m = 1
	}
	k := rng.Int63n(100) + 1
	return fmt.Sprintf("%d %d\n", m, k)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	cand := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := genCase(rng)
		exp, err := runProg(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on case %d: %v\n%s", i+1, err, input)
			os.Exit(1)
		}
		out, err := runProg(cand, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n%s", i+1, err, out)
			os.Exit(1)
		}
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:%s", i+1, exp, out, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

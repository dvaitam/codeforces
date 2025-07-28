package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func buildOracle() (string, error) {
	exe := "oracleC"
	dir, _ := os.Getwd()
	exe = filepath.Join(dir, exe)
	cmd := exec.Command("go", "build", "-o", exe, "1737C.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return exe, nil
}

func runProg(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return out.String() + errb.String(), fmt.Errorf("%v", err)
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(7) + 4
	r0 := rng.Intn(n-1) + 1
	c0 := rng.Intn(n-1) + 1
	r1, c1 := r0, c0
	r2, c2 := r0, c0+1
	r3, c3 := r0+1, c0
	orientation := rng.Intn(4)
	switch orientation {
	case 1:
		r1, c1 = r0, c0+1
		r2, c2 = r0+1, c0+1
		r3, c3 = r0, c0
	case 2:
		r1, c1 = r0+1, c0
		r2, c2 = r0+1, c0+1
		r3, c3 = r0, c0
	case 3:
		r1, c1 = r0, c0
		r2, c2 = r0+1, c0
		r3, c3 = r0+1, c0+1
	}
	x := rng.Intn(n) + 1
	y := rng.Intn(n) + 1
	return fmt.Sprintf("1\n%d\n%d %d %d %d %d %d\n%d %d\n", n, r1, c1, r2, c2, r3, c3, x, y)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := genCase(rng)
		expect, err := runProg(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runProg(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i+1, input, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}

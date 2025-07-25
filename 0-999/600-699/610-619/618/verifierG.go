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

func run(prog, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(prog, ".go") {
		cmd = exec.Command("go", "run", prog)
	} else {
		cmd = exec.Command(prog)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 120; t++ {
		n := rng.Int63n(1_000_000_000) + 1
		p := rng.Int63n(1_000_000_000-1) + 1
		input := fmt.Sprintf("%d %d\n", n, p)
		expectedStr, err := run("618G.go", input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "internal error on test %d: %v\n", t+1, err)
			os.Exit(1)
		}
		gotStr, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", t+1, err, input)
			os.Exit(1)
		}
		exVal, err := strconv.ParseFloat(strings.TrimSpace(expectedStr), 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "internal parse error: %v\n", err)
			os.Exit(1)
		}
		gotVal, err := strconv.ParseFloat(strings.TrimSpace(gotStr), 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: non float output %s\n", t+1, gotStr)
			os.Exit(1)
		}
		if math.Abs(gotVal-exVal) > 1e-4*math.Max(1, math.Abs(exVal)) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %.6f got %.6f\ninput:\n%s", t+1, exVal, gotVal, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

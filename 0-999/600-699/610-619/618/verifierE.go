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

func genCase(rng *rand.Rand, n, m int) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < m; i++ {
		typ := rng.Intn(2) + 1
		seg := rng.Intn(n) + 1
		val := rng.Intn(359) + 1
		sb.WriteString(fmt.Sprintf("%d %d %d\n", typ, seg, val))
	}
	return sb.String()
}

func parseFloat(s string) (float64, float64, error) {
	parts := strings.Fields(s)
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("expected two floats")
	}
	a, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return 0, 0, err
	}
	b, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return 0, 0, err
	}
	return a, b, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 120; t++ {
		n := rng.Intn(10) + 1
		m := rng.Intn(10) + 1
		input := genCase(rng, n, m)
		expectedStr, err := run("618E.go", input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "internal error on test %d: %v\n", t+1, err)
			os.Exit(1)
		}
		gotStr, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", t+1, err, input)
			os.Exit(1)
		}
		exX, exY, err := parseFloat(expectedStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "internal parse error: %v\n", err)
			os.Exit(1)
		}
		gX, gY, err := parseFloat(gotStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: expected two numbers got %s\n", t+1, gotStr)
			os.Exit(1)
		}
		if math.Abs(gX-exX) > 1e-4*math.Max(1, math.Abs(exX)) || math.Abs(gY-exY) > 1e-4*math.Max(1, math.Abs(exY)) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %.5f %.5f got %.5f %.5f\ninput:\n%s", t+1, exX, exY, gX, gY, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

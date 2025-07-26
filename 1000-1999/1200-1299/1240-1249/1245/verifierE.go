package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	oracle := filepath.Join(dir, "oracleE")
	cmd := exec.Command("go", "build", "-o", oracle, "1245E.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
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

func generateCase(rng *rand.Rand) string {
	var sb strings.Builder
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			val := 0
			if i > 0 {
				val = rng.Intn(i)
			}
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", val))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func parseFloat(out string) (float64, error) {
	scanner := bufio.NewScanner(strings.NewReader(out))
	if !scanner.Scan() {
		return 0, fmt.Errorf("no output")
	}
	return strconv.ParseFloat(strings.Fields(scanner.Text())[0], 64)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
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
		input := generateCase(rng)
		wantOut, err := run(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle error on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		want, err := parseFloat(wantOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle parse error on case %d: %v\noutput:%s", i+1, err, wantOut)
			os.Exit(1)
		}
		gotOut, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
		got, err := parseFloat(gotOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: cannot parse output: %v\noutput:%s", i+1, err, gotOut)
			os.Exit(1)
		}
		diff := got - want
		if diff < 0 {
			diff = -diff
		}
		if diff > 1e-6 {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %f got %f\ninput:%s", i+1, want, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

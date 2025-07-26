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
	oracle := filepath.Join(dir, "oracleD")
	cmd := exec.Command("go", "build", "-o", oracle, "1245D.go")
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
	n := rng.Intn(5) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		x := rng.Intn(20)
		y := rng.Intn(20)
		sb.WriteString(fmt.Sprintf("%d %d\n", x, y))
	}
	for i := 0; i < n; i++ {
		c := rng.Intn(100) + 1
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", c))
	}
	sb.WriteByte('\n')
	for i := 0; i < n; i++ {
		k := rng.Intn(10) + 1
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", k))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func parseCost(out string) (int64, error) {
	scanner := bufio.NewScanner(strings.NewReader(out))
	if !scanner.Scan() {
		return 0, fmt.Errorf("no output")
	}
	line := strings.Fields(scanner.Text())
	if len(line) == 0 {
		return 0, fmt.Errorf("empty output")
	}
	return strconv.ParseInt(line[0], 10, 64)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
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
		want, err := parseCost(wantOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle parse error on case %d: %v\noutput:%s", i+1, err, wantOut)
			os.Exit(1)
		}
		gotOut, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
		got, err := parseCost(gotOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: cannot parse output: %v\noutput:%s", i+1, err, gotOut)
			os.Exit(1)
		}
		if got != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected cost %d got %d\ninput:%s", i+1, want, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

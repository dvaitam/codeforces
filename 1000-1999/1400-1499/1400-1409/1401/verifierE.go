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
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	oracle := filepath.Join(dir, "oracleE")
	cmd := exec.Command("go", "build", "-o", oracle, "1401E.go")
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
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(5)
	m := rng.Intn(5)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	usedY := map[int]bool{}
	for i := 0; i < n; i++ {
		var y int
		for {
			y = rng.Intn(1_000_000-1) + 1
			if !usedY[y] {
				usedY[y] = true
				break
			}
		}
		if rng.Intn(2) == 0 {
			lx := 0
			rx := rng.Intn(1_000_000) + 1
			sb.WriteString(fmt.Sprintf("%d %d %d\n", y, lx, rx))
		} else {
			rx := 1_000_000
			lx := rng.Intn(1_000_000)
			sb.WriteString(fmt.Sprintf("%d %d %d\n", y, lx, rx))
		}
	}
	usedX := map[int]bool{}
	for i := 0; i < m; i++ {
		var x int
		for {
			x = rng.Intn(1_000_000-1) + 1
			if !usedX[x] {
				usedX[x] = true
				break
			}
		}
		if rng.Intn(2) == 0 {
			ly := 0
			ry := rng.Intn(1_000_000) + 1
			sb.WriteString(fmt.Sprintf("%d %d %d\n", x, ly, ry))
		} else {
			ry := 1_000_000
			ly := rng.Intn(1_000_000)
			sb.WriteString(fmt.Sprintf("%d %d %d\n", x, ly, ry))
		}
	}
	return sb.String()
}

func runCase(bin, oracle, input string) error {
	exp, err := run(oracle, input)
	if err != nil {
		return fmt.Errorf("oracle error: %v", err)
	}
	got, err := run(bin, input)
	if err != nil {
		return err
	}
	if strings.TrimSpace(got) != strings.TrimSpace(exp) {
		return fmt.Errorf("expected %s got %s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
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
		if err := runCase(bin, oracle, input); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

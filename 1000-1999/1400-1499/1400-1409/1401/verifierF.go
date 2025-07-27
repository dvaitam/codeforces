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
	oracle := filepath.Join(dir, "oracleF")
	cmd := exec.Command("go", "build", "-o", oracle, "1401F.go")
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
	N := 1 << n
	q := rng.Intn(15) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, q))
	for i := 0; i < N; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(rng.Intn(1000)))
	}
	sb.WriteByte('\n')
	sumCnt := 0
	for i := 0; i < q; i++ {
		typ := rng.Intn(4) + 1
		if i == q-1 && sumCnt == 0 {
			typ = 4
		}
		switch typ {
		case 1:
			x := rng.Intn(N) + 1
			k := rng.Intn(1000)
			sb.WriteString(fmt.Sprintf("1 %d %d\n", x, k))
		case 2:
			k := 0
			if n > 0 {
				k = rng.Intn(n + 1)
			}
			sb.WriteString(fmt.Sprintf("2 %d\n", k))
		case 3:
			k := 0
			if n > 0 {
				k = rng.Intn(n)
			}
			sb.WriteString(fmt.Sprintf("3 %d\n", k))
		case 4:
			l := rng.Intn(N) + 1
			r := rng.Intn(N-l+1) + l
			sb.WriteString(fmt.Sprintf("4 %d %d\n", l, r))
			sumCnt++
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
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

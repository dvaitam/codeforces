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
	cmd := exec.Command("go", "build", "-o", oracle, "163E.go")
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func randString(rng *rand.Rand, n int) string {
	letters := []byte{"a", "b", "c"}
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = letters[rng.Intn(len(letters))]
	}
	return string(b)
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(10) + 1
	k := rng.Intn(5) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
	patterns := make([]string, k)
	for i := 0; i < k; i++ {
		patterns[i] = randString(rng, rng.Intn(3)+1)
		sb.WriteString(patterns[i] + "\n")
	}
	for i := 0; i < n; i++ {
		typ := rng.Intn(3)
		switch typ {
		case 0:
			id := rng.Intn(k) + 1
			sb.WriteString(fmt.Sprintf("+%d\n", id))
		case 1:
			id := rng.Intn(k) + 1
			sb.WriteString(fmt.Sprintf("-%d\n", id))
		default:
			text := randString(rng, rng.Intn(4)+1)
			sb.WriteString("?" + text + "\n")
		}
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := generateCase(rng)
		expected, err := run(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on case %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expected) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %q got %q\ninput:\n%s", i+1, expected, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

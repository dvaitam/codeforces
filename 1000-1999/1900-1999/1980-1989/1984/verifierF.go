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
	cmd := exec.Command("go", "build", "-o", oracle, "1984F.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("%v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func randomString(rng *rand.Rand, n int) string {
	b := make([]byte, n)
	chars := []byte{'P', 'S', '?'}
	for i := range b {
		b[i] = chars[rng.Intn(3)]
	}
	return string(b)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
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

	for i := 1; i <= 100; i++ {
		n := rng.Intn(4) + 2
		m := rng.Intn(9) + 2
		s := randomString(rng, n)
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", n, m)
		fmt.Fprintf(&sb, "%s\n", s)
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			val := rng.Int63n(int64(m*n*2+1)) - int64(m*n)
			fmt.Fprintf(&sb, "%d", val)
		}
		sb.WriteByte('\n')
		input := fmt.Sprintf("1\n%s", sb.String())

		exp, err := run(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle error on case %d: %v\n", i, err)
			os.Exit(1)
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n", i, err)
			os.Exit(1)
		}
		if got != exp {
			fmt.Printf("case %d failed\ninput:\n%sexpected: %s\ngot: %s\n", i, input, exp, got)
			os.Exit(1)
		}
	}

	fmt.Println("All tests passed")
}

package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func buildOracle() (string, error) {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	src := filepath.Join(dir, "1463F.go")
	bin := filepath.Join(os.TempDir(), "oracle1463F.bin")
	cmd := exec.Command("go", "build", "-o", bin, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return bin, nil
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

func randomCase(rng *rand.Rand) string {
	n := rng.Uint64()%1000000 + 1
	x := rng.Intn(22) + 1
	y := rng.Intn(22) + 1
	return fmt.Sprintf("%d %d %d\n", n, x, y)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	userBin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]string, 0, 100)
	for len(cases) < 100 {
		cases = append(cases, randomCase(rng))
	}
	for i, in := range cases {
		want, err := run(oracle, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on case %d: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		got, err := run(userBin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if want != got {
			fmt.Printf("test %d failed\ninput:\n%sexpected: %s\ngot: %s\n", i+1, in, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

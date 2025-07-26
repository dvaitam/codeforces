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
)

func buildOracle() (string, error) {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	src := filepath.Join(dir, "1156G.go")
	bin := filepath.Join(os.TempDir(), "oracle1156G.bin")
	cmd := exec.Command("go", "build", "-o", bin, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return bin, nil
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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

var names = []string{"a", "b", "c", "d", "e", "f", "g", "h", "res"}
var ops = []byte{'$', '^', '#', '&'}

func randName(r *rand.Rand) string { return names[r.Intn(len(names))] }

func genCase(r *rand.Rand) string {
	n := r.Intn(5) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 0; i < n; i++ {
		l := randName(r)
		if r.Intn(2) == 0 {
			r1 := randName(r)
			fmt.Fprintf(&sb, "%s=%s\n", l, r1)
		} else {
			r1 := randName(r)
			r2 := randName(r)
			op := ops[r.Intn(len(ops))]
			fmt.Fprintf(&sb, "%s=%s%c%s\n", l, r1, op, r2)
		}
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	userBin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)
	r := rand.New(rand.NewSource(1))
	for i := 0; i < 100; i++ {
		input := genCase(r)
		want, err := run(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on test %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		got, err := run(userBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if want != got {
			fmt.Printf("test %d failed\ninput:\n%sexpected: %s\ngot: %s\n", i+1, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", 100)
}

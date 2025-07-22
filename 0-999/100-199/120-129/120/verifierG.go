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

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func generateWord(rng *rand.Rand) string {
	l := rng.Intn(5) + 1
	b := make([]byte, l)
	for i := range b {
		b[i] = byte('a' + rng.Intn(26))
	}
	return string(b)
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(2) + 1
	tturn := rng.Intn(5) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, tturn)
	for i := 0; i < n; i++ {
		fmt.Fprintf(&sb, "%d %d %d %d\n", rng.Intn(5)+1, rng.Intn(5)+1, rng.Intn(5)+1, rng.Intn(5)+1)
	}
	m := rng.Intn(3) + 1
	fmt.Fprintf(&sb, "%d\n", m)
	for i := 0; i < m; i++ {
		w := generateWord(rng)
		fmt.Fprintf(&sb, "%s\n%d\n", w, rng.Intn(5)+1)
	}
	input := sb.String()
	return input, runRef(input)
}

func runRef(input string) string {
	ref := filepath.Join(filepath.Dir(os.Args[0]), "refG")
	cmd := exec.Command(ref)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Run()
	return out.String()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: go run verifierG.go /path/to/binary\n")
		os.Exit(1)
	}
	candidate := os.Args[1]

	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	ref := filepath.Join(dir, "refG")
	cmd := exec.Command("go", "build", "-o", ref, filepath.Join(dir, "120G.go"))
	if out, err := cmd.CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference solution: %v\n%s", err, out)
		os.Exit(1)
	}
	defer os.Remove(ref)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		input, expect := generateCase(rng)
		candOut, cErr := runBinary(candidate, input)
		if cErr != nil {
			fmt.Fprintf(os.Stderr, "test %d: candidate error: %v\n", t+1, cErr)
			os.Exit(1)
		}
		if strings.TrimSpace(candOut) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "test %d failed\ninput:\n%sexpected:%sactual:%s\n", t+1, input, expect, candOut)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

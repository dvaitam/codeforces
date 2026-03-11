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

func buildRef() (string, error) {
	refSrc := os.Getenv("REFERENCE_SOURCE_PATH")
	if refSrc == "" {
		return "", fmt.Errorf("REFERENCE_SOURCE_PATH not set")
	}
	ref := filepath.Join(os.TempDir(), fmt.Sprintf("ref1252J_%d", time.Now().UnixNano()))
	cmd := exec.Command("go", "build", "-o", ref, refSrc)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v\n%s", err, out)
	}
	return ref, nil
}

func runExe(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("%v\n%s", err, stderr.String())
	}
	return out.String(), nil
}

func genTests(rng *rand.Rand) []string {
	tests := make([]string, 0, 100)
	for i := 0; i < 100; i++ {
		N := rng.Intn(20) + 1
		K := rng.Intn(N + 1) // 0 <= K <= N
		G1 := rng.Intn(100)
		G2 := rng.Intn(100)
		G3 := rng.Intn(100)
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d %d %d %d\n", N, K, G1, G2, G3)
		rocks := 0
		for j := 0; j < N; j++ {
			if rng.Intn(5) == 0 && rocks < 50 {
				sb.WriteByte('#')
				rocks++
			} else {
				sb.WriteByte('.')
			}
		}
		sb.WriteByte('\n')
		tests = append(tests, sb.String())
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierJ.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	rng := rand.New(rand.NewSource(42))
	tests := genTests(rng)
	for i, input := range tests {
		exp, err := runExe(ref, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runExe(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(exp) != strings.TrimSpace(got) {
			fmt.Fprintf(os.Stderr, "Test %d failed\nInput:\n%sExpected: %sGot: %s\n", i+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}

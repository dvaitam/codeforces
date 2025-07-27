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
	src := filepath.Join(dir, "1326F1.go")
	bin := filepath.Join(os.TempDir(), "oracle1326F1.bin")
	cmd := exec.Command("go", "build", "-o", bin, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return bin, nil
}

func runCandidate(bin, input string) (string, error) {
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

func deterministicCases() []string {
	return []string{
		"2\n01\n10\n",
		"3\n011\n101\n110\n",
	}
}

func genMatrix(rng *rand.Rand, n int) []string {
	mat := make([]string, n)
	for i := 0; i < n; i++ {
		b := make([]byte, n)
		for j := 0; j < n; j++ {
			if i == j {
				b[j] = '0'
			} else {
				if rng.Intn(2) == 0 {
					b[j] = '0'
				} else {
					b[j] = '1'
				}
			}
		}
		mat[i] = string(b)
	}
	// ensure symmetry
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if mat[i][j] != mat[j][i] {
				if rng.Intn(2) == 0 {
					b := []byte(mat[i])
					b[j] = mat[j][i]
					mat[i] = string(b)
				} else {
					b := []byte(mat[j])
					b[i] = mat[i][j]
					mat[j] = string(b)
				}
			}
		}
	}
	return mat
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(5) + 2
	mat := genMatrix(rng, n)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 0; i < n; i++ {
		sb.WriteString(mat[i])
		sb.WriteByte('\n')
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF1.go /path/to/binary")
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
	cases := deterministicCases()
	for len(cases) < 100 {
		cases = append(cases, genCase(rng))
	}

	for i, in := range cases {
		want, err := runCandidate(oracle, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on case %d: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		got, err := runCandidate(userBin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if got != want {
			fmt.Printf("case %d failed\ninput:\n%sexpected: %s\n got: %s\n", i+1, in, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

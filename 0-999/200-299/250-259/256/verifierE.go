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

const tests = 40

func buildOracle() (string, error) {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	src := filepath.Join(dir, "256E.go")
	tmp, err := os.CreateTemp("", "oracle256E")
	if err != nil {
		return "", err
	}
	path := tmp.Name()
	tmp.Close()
	os.Remove(path)
	cmd := exec.Command("go", "build", "-o", path, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build oracle: %v\n%s", err, out)
	}
	return path, nil
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

func genMatrix(r *rand.Rand) [3][3]int {
	var w [3][3]int
	ones := 0
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if r.Intn(3) != 0 {
				w[i][j] = 1
				ones++
			}
		}
	}
	if ones == 0 {
		w[r.Intn(3)][r.Intn(3)] = 1
	}
	return w
}

func genCase(r *rand.Rand) string {
	n := r.Intn(25) + 1
	m := r.Intn(40) + 1
	w := genMatrix(r)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteByte(byte('0' + w[i][j]))
		}
		sb.WriteByte('\n')
	}
	forcedIdx := r.Intn(n) + 1
	forcedVal := r.Intn(3) + 1
	for i := 0; i < m; i++ {
		v := r.Intn(n) + 1
		t := r.Intn(4)
		if i == 0 {
			v = forcedIdx
			t = forcedVal
		} else if i == 1 {
			v = forcedIdx
			t = 0
		}
		fmt.Fprintf(&sb, "%d %d\n", v, t)
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
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
	for i := 0; i < tests; i++ {
		input := genCase(r)
		want, err := run(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := run(userBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if want != got {
			fmt.Printf("test %d failed\ninput:\n%s\nexpected: %s\ngot: %s\n", i+1, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", tests)
}

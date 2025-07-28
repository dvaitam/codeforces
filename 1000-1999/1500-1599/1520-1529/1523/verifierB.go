package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const numTestsB = 100

func prepareBinary(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp := filepath.Join(os.TempDir(), "candB")
		cmd := exec.Command("go", "build", "-o", tmp, path)
		if out, err := cmd.CombinedOutput(); err != nil {
			return "", nil, fmt.Errorf("go build failed: %v: %s", err, out)
		}
		return tmp, func() { os.Remove(tmp) }, nil
	}
	return path, nil, nil
}

func buildOracle() (string, func(), error) {
	tmp := filepath.Join(os.TempDir(), "oracleB")
	cmd := exec.Command("go", "build", "-o", tmp, "1523B.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", nil, fmt.Errorf("build oracle failed: %v: %s", err, out)
	}
	return tmp, func() { os.Remove(tmp) }, nil
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	err := cmd.Run()
	return strings.TrimSpace(buf.String()), err
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(10)*2 + 2 // even 2..20
	arr := make([]int, n)
	for i := range arr {
		arr[i] = rng.Intn(21) - 10
	}
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		return
	}
	bin, clean, err := prepareBinary(os.Args[1])
	if err != nil {
		fmt.Println("compile error:", err)
		return
	}
	if clean != nil {
		defer clean()
	}
	oracle, c2, err := buildOracle()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer c2()

	rng := rand.New(rand.NewSource(1))
	for i := 0; i < numTestsB; i++ {
		input := genCase(rng)
		want, err := run(oracle, input)
		if err != nil {
			fmt.Printf("oracle runtime error on case %d: %v\n", i+1, err)
			return
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("candidate runtime error on case %d: %v\ninput:\n%s", i+1, err, input)
			return
		}
		if want != got {
			fmt.Printf("case %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i+1, input, want, got)
			return
		}
	}
	fmt.Println("All tests passed")
}

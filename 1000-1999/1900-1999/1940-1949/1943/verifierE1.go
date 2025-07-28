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
	src := filepath.Join(dir, "1943E1.go")
	bin := filepath.Join(os.TempDir(), "oracle1943E1.bin")
	cmd := exec.Command("go", "build", "-o", bin, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return bin, nil
}

func run(prog string, input []byte) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(prog, ".go") {
		cmd = exec.Command("go", "run", prog)
	} else {
		cmd = exec.Command(prog)
	}
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genTest() []byte {
	m := rand.Intn(5) + 1
	k := rand.Intn(5) + 1
	freq := make([]int, m+1)
	for i := range freq {
		freq[i] = rand.Intn(5)
	}
	var sb strings.Builder
	sb.WriteString("1\n")
	fmt.Fprintf(&sb, "%d %d\n", m, k)
	for i := 0; i <= m; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", freq[i])
	}
	sb.WriteByte('\n')
	return []byte(sb.String())
}

func main() {
	rand.Seed(time.Now().UnixNano())
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE1.go /path/to/binary")
		return
	}
	cand := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	for i := 1; i <= 100; i++ {
		in := genTest()
		expect, err := run(oracle, in)
		if err != nil {
			fmt.Printf("oracle failed on test %d: %v\n", i, err)
			os.Exit(1)
		}
		got, err := run(cand, in)
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\n", i, err)
			os.Exit(1)
		}
		if expect != got {
			fmt.Printf("wrong answer on test %d\ninput:\n%sexpected:%s\ngot:%s\n", i, string(in), expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

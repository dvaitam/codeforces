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
	src := filepath.Join(dir, "1943D2.go")
	bin := filepath.Join(os.TempDir(), "oracle1943D2.bin")
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
	n := rand.Intn(5) + 3
	k := rand.Intn(n) + 1
	p := 1000000007
	return []byte(fmt.Sprintf("1\n%d %d %d\n", n, k, p))
}

func main() {
	rand.Seed(time.Now().UnixNano())
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD2.go /path/to/binary")
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

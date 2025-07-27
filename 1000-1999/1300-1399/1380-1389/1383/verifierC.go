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

func run(bin string, input []byte) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
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

func buildRef() (string, error) {
	_, cur, _, _ := runtime.Caller(0)
	dir := filepath.Dir(cur)
	src := filepath.Join(dir, "1383C.go")
	bin := filepath.Join(os.TempDir(), "1383C_ref.bin")
	cmd := exec.Command("go", "build", "-o", bin, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v\n%s", err, string(out))
	}
	return bin, nil
}

func genCase(rng *rand.Rand) []byte {
	n := rng.Intn(20) + 1
	a := make([]byte, n)
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		a[i] = byte('a' + rng.Intn(20))
		b[i] = byte('a' + rng.Intn(20))
	}
	return []byte(fmt.Sprintf("1\n%d\n%s\n%s\n", n, string(a), string(b)))
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	cand := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= 100; i++ {
		in := genCase(rng)
		exp, err := run(ref, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\n", i, err)
			os.Exit(1)
		}
		got, err := run(cand, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\n", i, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d\ninput:\n%sexpected:%s\ngot:%s\n", i, string(in), exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

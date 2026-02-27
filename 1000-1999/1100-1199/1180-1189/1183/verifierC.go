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

func buildOracle() (string, func(), error) {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	src := filepath.Join(dir, "1183C.go")
	tmp, err := os.CreateTemp("", "oracle1183C*")
	if err != nil {
		return "", nil, err
	}
	tmp.Close()
	cmd := exec.Command("go", "build", "-o", tmp.Name(), src)
	if out, err := cmd.CombinedOutput(); err != nil {
		os.Remove(tmp.Name())
		return "", nil, fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
}

func runProg(prog, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(prog, ".go") {
		cmd = exec.Command("go", "run", prog)
	} else {
		cmd = exec.Command(prog)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return out.String() + stderr.String(), fmt.Errorf("%v", err)
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(rng *rand.Rand) string {
	q := rng.Intn(5) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", q))
	for i := 0; i < q; i++ {
		k := rng.Int63n(100000) + 1
		n := rng.Int63n(100000) + 1
		a := rng.Int63n(1000) + 2
		b := rng.Int63n(a-1) + 1
		fmt.Fprintf(&sb, "%d %d %d %d\n", k, n, a, b)
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, cleanup, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := genCase(rng)
		exp, err := runProg(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle runtime error on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		out, err := runProg(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n%s", i+1, err, out)
			os.Exit(1)
		}
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d failed:\nexpected: %s\ngot: %s\ninput:\n%s", i+1, exp, out, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

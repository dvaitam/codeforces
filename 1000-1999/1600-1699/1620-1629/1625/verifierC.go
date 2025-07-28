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
	_, cur, _, _ := runtime.Caller(0)
	dir := filepath.Dir(cur)
	src := filepath.Join(dir, "1625C.go")
	bin := filepath.Join(os.TempDir(), "oracle1625C.bin")
	cmd := exec.Command("go", "build", "-o", bin, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return bin, nil
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(r *rand.Rand) string {
	n := r.Intn(8) + 1
	l := n + r.Intn(20) + 1
	k := r.Intn(n) + 1
	d := make([]int, n)
	cur := 0
	for i := 0; i < n; i++ {
		cur += r.Intn(5) + 1
		d[i] = cur
	}
	if cur >= l {
		l = cur + r.Intn(5) + 1
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = r.Intn(10) + 1
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", n, l, k)
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", d[i])
	}
	sb.WriteByte('\n')
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", a[i])
	}
	sb.WriteByte('\n')
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	userBin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
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

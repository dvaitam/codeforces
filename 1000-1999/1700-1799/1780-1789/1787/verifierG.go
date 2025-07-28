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
	src := filepath.Join(dir, "1787G.go")
	bin := filepath.Join(os.TempDir(), "oracle1787G.bin")
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
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(r *rand.Rand) string {
	n := r.Intn(5) + 2
	q := r.Intn(5) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, q))
	for i := 2; i <= n; i++ {
		u := i
		v := r.Intn(i-1) + 1
		w := r.Intn(n) + 1
		c := r.Intn(n) + 1
		sb.WriteString(fmt.Sprintf("%d %d %d %d\n", u, v, w, c))
	}
	blocked := make([]bool, n+1)
	for i := 0; i < q; i++ {
		op := 0
		anyBlocked := false
		for j := 1; j <= n; j++ {
			if blocked[j] {
				anyBlocked = true
				break
			}
		}
		if anyBlocked && r.Intn(2) == 0 {
			op = 1
		}
		var x int
		if op == 0 {
			for {
				x = r.Intn(n) + 1
				if !blocked[x] {
					break
				}
			}
			blocked[x] = true
		} else {
			for {
				x = r.Intn(n) + 1
				if blocked[x] {
					break
				}
			}
			blocked[x] = false
		}
		sb.WriteString(fmt.Sprintf("%d %d\n", op, x))
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= 100; i++ {
		input := genCase(r)
		want, err := run(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle error on case %d: %v\n", i, err)
			os.Exit(1)
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i, err, input)
			os.Exit(1)
		}
		if got != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i, want, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All 100 tests passed")
}

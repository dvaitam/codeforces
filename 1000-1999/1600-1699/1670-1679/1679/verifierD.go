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

const numTests = 100

func buildRef() (string, error) {
	_, cur, _, _ := runtime.Caller(0)
	dir := filepath.Dir(cur)
	src := filepath.Join(dir, "1679D.go")
	bin := filepath.Join(os.TempDir(), "1679D_ref.bin")
	cmd := exec.Command("go", "build", "-o", bin, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v\n%s", err, out)
	}
	return bin, nil
}

func run(path string, in []byte) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = bytes.NewReader(in)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genTest(rng *rand.Rand) []byte {
	n := rng.Intn(6) + 1
	maxEdges := n * (n - 1)
	m := 0
	if maxEdges > 0 {
		m = rng.Intn(maxEdges) + 1
	}
	k := rng.Intn(n) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", n, m, k)
	for i := 0; i < n; i++ {
		fmt.Fprintf(&sb, "%d ", rng.Intn(20)+1)
	}
	sb.WriteByte('\n')
	for i := 0; i < m; i++ {
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		for u == v {
			v = rng.Intn(n) + 1
		}
		fmt.Fprintf(&sb, "%d %d\n", u, v)
	}
	return []byte(sb.String())
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	bin := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	for i := 0; i < numTests; i++ {
		input := genTest(rng)
		expect, err := run(ref, input)
		if err != nil {
			fmt.Printf("reference failed on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(expect) != strings.TrimSpace(got) {
			fmt.Printf("wrong answer on test %d\ninput:\n%s\nexpected:%s\nGot:%s\n", i+1, string(input), expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

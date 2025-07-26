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

func buildRef() (string, error) {
	_, cur, _, _ := runtime.Caller(0)
	dir := filepath.Dir(cur)
	src := filepath.Join(dir, "1070I.go")
	bin := filepath.Join(os.TempDir(), "1070I_ref.bin")
	cmd := exec.Command("go", "build", "-o", bin, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v\n%s", err, out)
	}
	return bin, nil
}

func runProg(bin string, input []byte) (string, error) {
	cmd := exec.Command(bin)
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

func genCase(rng *rand.Rand) []byte {
	t := rng.Intn(2) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	for ; t > 0; t-- {
		n := rng.Intn(4) + 1
		m := rng.Intn(n*(n-1)/2 + 1)
		k := rng.Intn(n) + 1
		sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, k))
		edges := make(map[[2]int]struct{})
		for i := 0; i < m; i++ {
			var u, v int
			for {
				u = rng.Intn(n) + 1
				v = rng.Intn(n) + 1
				if u != v {
					if u > v {
						u, v = v, u
					}
					if _, ok := edges[[2]int{u, v}]; !ok {
						edges[[2]int{u, v}] = struct{}{}
						break
					}
				}
			}
			sb.WriteString(fmt.Sprintf("%d %d\n", u, v))
		}
	}
	return []byte(sb.String())
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: go run verifierI.go /path/to/binary")
		os.Exit(1)
	}
	cand := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= 100; i++ {
		input := genCase(rng)
		want, err := runProg(ref, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\n", i, err)
			os.Exit(1)
		}
		got, err := runProg(cand, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\ninput:\n%s", i, err, string(input))
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(want) {
			fmt.Printf("wrong answer on test %d\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i, string(input), want, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

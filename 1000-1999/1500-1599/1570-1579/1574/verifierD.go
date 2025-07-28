package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func buildRef() (string, error) {
	ref := filepath.Join(os.TempDir(), "1574D_ref.bin")
	cmd := exec.Command("go", "build", "-o", ref, "1574D.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v\n%s", err, out)
	}
	return ref, nil
}

func runProg(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return out.String() + errBuf.String(), fmt.Errorf("runtime error: %v", err)
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(3) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	counts := make([]int, n)
	for i := 0; i < n; i++ {
		c := rng.Intn(3) + 1
		counts[i] = c
		fmt.Fprintf(&sb, "%d", c)
		for j := 0; j < c; j++ {
			val := rng.Intn(20) + 1
			fmt.Fprintf(&sb, " %d", val)
		}
		sb.WriteByte('\n')
	}
	m := rng.Intn(3)
	fmt.Fprintf(&sb, "%d\n", m)
	seen := make(map[string]bool)
	for len(seen) < m {
		var build strings.Builder
		for i := 0; i < n; i++ {
			idx := rng.Intn(counts[i]) + 1
			if i > 0 {
				build.WriteByte(' ')
			}
			fmt.Fprintf(&build, "%d", idx)
		}
		s := build.String()
		if seen[s] {
			continue
		}
		seen[s] = true
		sb.WriteString(s)
		sb.WriteByte('\n')
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
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
	for i := 0; i < 100; i++ {
		input := genCase(rng)
		exp, err := runProg(ref, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		out, err := runProg(cand, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on case %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i+1, input, exp, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

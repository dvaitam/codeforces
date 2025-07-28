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

func runBinary(bin string, input []byte) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func buildRef() (string, error) {
	_, cur, _, _ := runtime.Caller(0)
	dir := filepath.Dir(cur)
	src := filepath.Join(dir, "1490C.go")
	ref := filepath.Join(os.TempDir(), "1490C_ref.bin")
	cmd := exec.Command("go", "build", "-o", ref, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v\n%s", err, string(out))
	}
	return ref, nil
}

func genTest(rng *rand.Rand) []byte {
	// x up to 1e12
	x := rng.Int63n(1_000_000_000_000) + 1
	return []byte(fmt.Sprintf("1\n%d\n", x))
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// some fixed edge cases
	edges := []int64{1, 2, 4, 35, 34, 16, 703657519796}
	for i, x := range edges {
		input := []byte(fmt.Sprintf("1\n%d\n", x))
		exp, err := runBinary(ref, input)
		if err != nil {
			fmt.Printf("reference runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(exp) != strings.TrimSpace(got) {
			fmt.Printf("wrong answer on fixed test %d\ninput:\n%sexpected:%s\ngot:%s\n", i+1, string(input), exp, got)
			os.Exit(1)
		}
	}

	for i := 1; i <= 100; i++ {
		input := genTest(rng)
		exp, err := runBinary(ref, input)
		if err != nil {
			fmt.Printf("reference runtime error on random test %d: %v\n", i, err)
			os.Exit(1)
		}
		got, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("candidate runtime error on random test %d: %v\n", i, err)
			os.Exit(1)
		}
		if strings.TrimSpace(exp) != strings.TrimSpace(got) {
			fmt.Printf("wrong answer on random test %d\ninput:\n%sexpected:%s\ngot:%s\n", i, string(input), exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

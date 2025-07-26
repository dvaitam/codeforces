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

func buildRef(dir string) (string, error) {
	ref := filepath.Join(dir, "refA_bin")
	cmd := exec.Command("go", "build", "-o", ref, "896A.go")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("build reference failed: %v\n%s", err, string(out))
	}
	return ref, nil
}

func runBinary(bin string, input []byte) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)

	ref, err := buildRef(dir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	for t := 1; t <= numTests; t++ {
		q := rand.Intn(5) + 1
		var input bytes.Buffer
		input.WriteString(fmt.Sprintf("%d\n", q))
		for i := 0; i < q; i++ {
			n := rand.Intn(100000)
			k := rand.Int63n(1e18) + 1
			input.WriteString(fmt.Sprintf("%d %d\n", n, k))
		}
		inBytes := input.Bytes()

		out, err := runBinary(target, inBytes)
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\n", t, err)
			os.Exit(1)
		}

		refOut, err := runBinary(ref, inBytes)
		if err != nil {
			fmt.Printf("Test %d: reference runtime error: %v\n", t, err)
			os.Exit(1)
		}

		if strings.TrimSpace(out) != strings.TrimSpace(refOut) {
			fmt.Printf("Test %d failed\nInput:\n%sExpected:\n%sGot:\n%s\n", t, string(inBytes), refOut, out)
			os.Exit(1)
		}
	}

	fmt.Printf("Passed %d tests\n", numTests)
}

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

func buildReference() (string, error) {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	refSource := filepath.Join(dir, "961E.go")
	out := filepath.Join(os.TempDir(), "refE.bin")
	cmd := exec.Command("go", "build", "-o", out, refSource)
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return out, nil
}

func runBinary(path string, input []byte) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func genTest() string {
	n := rand.Intn(40) + 1
	a := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = rand.Intn(50) + 1
	}
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range a {
		if i > 0 {
			b.WriteByte(' ')
		}
		b.WriteString(fmt.Sprintf("%d", v))
	}
	b.WriteByte('\n')
	return b.String()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	cand := os.Args[1]
	rand.Seed(time.Now().UnixNano())

	ref, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}

	for i := 0; i < 100; i++ {
		test := genTest()
		expect, err1 := runBinary(ref, []byte(test))
		got, err2 := runBinary(cand, []byte(test))
		if err2 != nil {
			fmt.Printf("Test %d: candidate runtime error: %v\n", i+1, err2)
			fmt.Println("Input:\n" + test)
			os.Exit(1)
		}
		if err1 != nil {
			fmt.Println("reference runtime error:", err1)
			os.Exit(1)
		}
		if expect != got {
			fmt.Printf("Test %d failed\n", i+1)
			fmt.Println("Input:\n" + test)
			fmt.Println("Expected:", expect)
			fmt.Println("Got:", got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

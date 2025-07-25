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
)

func buildReference() (string, error) {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	src := filepath.Join(dir, "803C.go")
	bin := filepath.Join(os.TempDir(), "ref803C.bin")
	cmd := exec.Command("go", "build", "-o", bin, src)
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return bin, nil
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &bytes.Buffer{}
	err := cmd.Run()
	return out.String(), err
}

func genCase() string {
	n := rand.Int63n(1e6) + 1
	k := rand.Int63n(10) + 1
	return fmt.Sprintf("%d %d\n", n, k)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: verifierC.go <path-to-binary>")
		return
	}
	userBin := os.Args[1]
	rand.Seed(1)
	refBin, err := buildReference()
	if err != nil {
		fmt.Println("failed to build reference:", err)
		return
	}
	defer os.Remove(refBin)

	for i := 0; i < 100; i++ {
		input := genCase()
		want, err1 := runBinary(refBin, input)
		if err1 != nil {
			fmt.Println("reference solution failed:", err1)
			return
		}
		got, err2 := runBinary(userBin, input)
		if err2 != nil {
			fmt.Printf("test %d: runtime error: %v\n", i+1, err2)
			fmt.Println("input:\n" + input)
			return
		}
		if strings.TrimSpace(want) != strings.TrimSpace(got) {
			fmt.Printf("test %d failed\ninput:\n%sexpected:\n%s\nactual:\n%s\n", i+1, input, want, got)
			return
		}
	}
	fmt.Println("All tests passed")
}

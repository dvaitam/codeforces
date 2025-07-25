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
	src := filepath.Join(dir, "803A.go")
	bin := filepath.Join(os.TempDir(), "ref803A.bin")
	cmd := exec.Command("go", "build", "-o", bin, src)
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return bin, nil
}

func runBinary(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &bytes.Buffer{}
	err := cmd.Run()
	return out.String(), err
}

func genCase() string {
	// mix of deterministic and random cases
	switch rand.Intn(5) {
	case 0:
		return "1 0\n"
	case 1:
		return "1 1\n"
	case 2:
		return "2 3\n"
	case 3:
		return fmt.Sprintf("%d %d\n", 5, 10)
	}
	n := rand.Intn(10) + 1
	maxK := n * n
	if rand.Intn(10) == 0 {
		n = rand.Intn(100) + 1
		maxK = n * n
		if maxK > 1000000 {
			maxK = 1000000
		}
	}
	k := rand.Intn(maxK + 1)
	return fmt.Sprintf("%d %d\n", n, k)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: verifierA.go <path-to-binary>")
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

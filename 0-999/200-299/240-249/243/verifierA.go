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

func compileRef(src string) (string, error) {
	tmp, err := os.CreateTemp("", "refA-")
	if err != nil {
		return "", err
	}
	tmp.Close()
	os.Remove(tmp.Name())
	cmd := exec.Command("go", "build", "-o", tmp.Name(), src)
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return tmp.Name(), nil
}

func runBinary(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func genTest() string {
	n := rand.Intn(50) + 1
	arr := make([]int, n)
	for i := range arr {
		arr[i] = rand.Intn(1000)
	}
	var b strings.Builder
	fmt.Fprintln(&b, n)
	for i, v := range arr {
		if i > 0 {
			b.WriteByte(' ')
		}
		fmt.Fprint(&b, v)
	}
	b.WriteByte('\n')
	return b.String()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	candidate, _ := filepath.Abs(os.Args[1])

	rand.Seed(time.Now().UnixNano())

	refBin, err := compileRef("243A.go")
	if err != nil {
		fmt.Println("failed to compile reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	for i := 0; i < 100; i++ {
		input := genTest()
		expectedOut, err := runBinary(refBin, input)
		if err != nil {
			fmt.Println("reference solution failed:", err)
			os.Exit(1)
		}
		expected := strings.TrimSpace(expectedOut)

		actualOut, err := runBinary(candidate, input)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		actual := strings.TrimSpace(actualOut)
		if actual != expected {
			fmt.Printf("test %d failed:\ninput:\n%s\nexpected: %s\ngot: %s\n", i+1, input, expected, actual)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}

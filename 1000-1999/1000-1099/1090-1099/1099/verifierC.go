package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func buildRef() (string, error) {
	ref := "./refC_bin"
	cmd := exec.Command("go", "build", "-o", ref, "1099C.go")
	cmd.Stdout = new(bytes.Buffer)
	cmd.Stderr = new(bytes.Buffer)
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return ref, nil
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = new(bytes.Buffer)
	err := cmd.Run()
	return out.String(), err
}

func genTests() []string {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]string, 100)
	letters := []rune("abcdefghijklmnopqrstuvwxyz")
	for i := 0; i < 100; i++ {
		l := rng.Intn(15) + 1
		var sb strings.Builder
		for j := 0; j < l; j++ {
			ch := letters[rng.Intn(len(letters))]
			sb.WriteRune(ch)
			r := rng.Intn(3)
			if r == 0 {
				sb.WriteByte('?')
			} else if r == 1 {
				sb.WriteByte('*')
			}
		}
		k := rng.Intn(20) + 1
		tests[i] = fmt.Sprintf("%s\n%d\n", sb.String(), k)
	}
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	cand := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Println("failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	tests := genTests()
	for idx, input := range tests {
		exp, err1 := runBinary(ref, input)
		out, err2 := runBinary(cand, input)
		if err1 != nil || err2 != nil {
			fmt.Printf("runtime error on test %d\n", idx+1)
			os.Exit(1)
		}
		if strings.TrimSpace(exp) != strings.TrimSpace(out) {
			fmt.Printf("wrong answer on test %d\ninput:\n%sexpected:\n%s\ngot:\n%s", idx+1, input, exp, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

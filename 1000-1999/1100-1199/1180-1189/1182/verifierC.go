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

func run(bin string, input []byte) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func buildRef() (string, error) {
	ref := "./refC.bin"
	cmd := exec.Command("go", "build", "-o", ref, "1182C.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v: %s", err, string(out))
	}
	return ref, nil
}

func randomWord() string {
	l := rand.Intn(8) + 1
	vowels := "aeiou"
	var sb strings.Builder
	ensure := rand.Intn(l)
	for i := 0; i < l; i++ {
		if i == ensure {
			sb.WriteByte(vowels[rand.Intn(len(vowels))])
		} else {
			sb.WriteByte(byte('a' + rand.Intn(26)))
		}
	}
	return sb.String()
}

func genTest() []byte {
	n := rand.Intn(20) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		sb.WriteString(randomWord())
		sb.WriteByte('\n')
	}
	return []byte(sb.String())
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	cand := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 100; i++ {
		input := genTest()
		want, err := run(ref, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error: %v\n", err)
			os.Exit(1)
		}
		got, err := run(cand, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\ninput:\n%s", i+1, err, string(input))
			os.Exit(1)
		}
		if strings.TrimSpace(want) != strings.TrimSpace(got) {
			fmt.Printf("wrong answer on test %d\ninput:\n%sexpected:\n%s\ngot:\n%s", i+1, string(input), want, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}

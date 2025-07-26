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
	ref := "./refD.bin"
	cmd := exec.Command("go", "build", "-o", ref, "1182D.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v: %s", err, string(out))
	}
	return ref, nil
}

func genTest() []byte {
	n := rand.Intn(20) + 1
	if n == 1 {
		return []byte("1\n")
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for v := 2; v <= n; v++ {
		p := rand.Intn(v-1) + 1
		sb.WriteString(fmt.Sprintf("%d %d\n", p, v))
	}
	return []byte(sb.String())
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
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

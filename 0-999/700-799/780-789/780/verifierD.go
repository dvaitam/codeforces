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

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func randName() string {
	l := rand.Intn(3) + 3 // length 3..5
	b := make([]byte, l)
	for i := range b {
		b[i] = byte('A' + rand.Intn(26))
	}
	return string(b)
}

func generateCase() string {
	n := rand.Intn(7) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		sb.WriteString(randName())
		sb.WriteByte(' ')
		sb.WriteString(randName())
		sb.WriteByte('\n')
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]
	rand.Seed(time.Now().UnixNano())

	ref := "refD_bin"
	if err := exec.Command("go", "build", "-o", ref, "780D.go").Run(); err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference solution:", err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	for i := 0; i < 100; i++ {
		input := generateCase()
		want, err := runBinary("./"+ref, input)
		if err != nil {
			fmt.Fprintln(os.Stderr, "failed to run reference solution:", err)
			os.Exit(1)
		}
		got, err := runBinary(target, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: candidate error: %v\n", i+1, err)
			os.Exit(1)
		}
		if want != got {
			fmt.Printf("test %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i+1, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

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

func generateCase() string {
	n := rand.Intn(5) + 1
	m := rand.Intn(n*n + 1)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < m; i++ {
		v := rand.Intn(n) + 1
		u := rand.Intn(n) + 1
		t := rand.Intn(2)
		sb.WriteString(fmt.Sprintf("%d %d %d\n", v, u, t))
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]
	rand.Seed(time.Now().UnixNano())

	ref := "refF_bin"
	if err := exec.Command("go", "build", "-o", ref, "780F.go").Run(); err != nil {
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
			fmt.Printf("test %d failed\ninput:\n%sexpected:%s\ngot:%s\n", i+1, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

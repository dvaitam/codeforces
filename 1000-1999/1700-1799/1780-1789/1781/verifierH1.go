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
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func generateCase() string {
	h := rand.Intn(4) + 1
	w := rand.Intn(4) + 1
	maxK := h * w
	if maxK > 4 {
		maxK = 4
	}
	k := rand.Intn(maxK + 1)
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d %d\n", h, w, k))
	for i := 0; i < k; i++ {
		r := rand.Intn(h) + 1
		c := rand.Intn(w) + 1
		sb.WriteString(fmt.Sprintf("%d %d\n", r, c))
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierH1.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]
	rand.Seed(time.Now().UnixNano())

	ref := "refH1_bin"
	if err := exec.Command("go", "build", "-o", ref, "1781H1.go").Run(); err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference solution:", err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	for i := 0; i < 100; i++ {
		input := generateCase()
		want, err := runBinary("./"+ref, input)
		if err != nil {
			fmt.Fprintln(os.Stderr, "reference runtime error:", err)
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

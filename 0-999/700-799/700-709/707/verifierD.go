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

func compileRef() (string, error) {
	ref := "refD"
	cmd := exec.Command("go", "build", "-o", ref, "707D.go")
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return ref, nil
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func generateTest() string {
	n := rand.Intn(3) + 1
	m := rand.Intn(3) + 1
	q := rand.Intn(20) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, q))
	for i := 1; i <= q; i++ {
		t := rand.Intn(4) + 1
		if i == 1 && t == 4 {
			t = rand.Intn(3) + 1
		}
		switch t {
		case 1:
			sb.WriteString(fmt.Sprintf("1 %d %d\n", rand.Intn(n)+1, rand.Intn(m)+1))
		case 2:
			sb.WriteString(fmt.Sprintf("2 %d %d\n", rand.Intn(n)+1, rand.Intn(m)+1))
		case 3:
			sb.WriteString(fmt.Sprintf("3 %d\n", rand.Intn(n)+1))
		case 4:
			sb.WriteString(fmt.Sprintf("4 %d\n", rand.Intn(i)))
		}
	}
	return sb.String()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	rand.Seed(time.Now().UnixNano())
	ref, err := compileRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, "could not build reference solution:", err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	bin := os.Args[1]
	for t := 0; t < 100; t++ {
		input := generateTest()
		want, err := run(ref, input)
		if err != nil {
			fmt.Fprintln(os.Stderr, "reference failed:", err)
			os.Exit(1)
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintln(os.Stderr, "test", t+1, "error running binary:", err)
			os.Exit(1)
		}
		if got != want {
			fmt.Fprintf(os.Stderr, "test %d failed\ninput:\n%s\nexpected: %s\nactual: %s\n", t+1, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

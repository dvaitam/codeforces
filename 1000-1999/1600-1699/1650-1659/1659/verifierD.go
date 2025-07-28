package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

func run(bin, in string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(in)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func runRef(in string) (string, error) {
	cmd := exec.Command("go", "run", "1659D.go")
	cmd.Stdin = strings.NewReader(in)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func genCase() string {
	n := rand.Intn(10) + 1
	A := make([]int, n)
	for i := 0; i < n; i++ {
		A[i] = rand.Intn(2)
	}
	C := make([]int, n)
	for k := 1; k <= n; k++ {
		tmp := append([]int(nil), A...)
		sort.Ints(tmp[:k])
		for i := 0; i < n; i++ {
			C[i] += tmp[i]
		}
	}
	var sb strings.Builder
	sb.WriteString("1\n")
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", C[i])
	}
	sb.WriteByte('\n')
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(1)
	for t := 1; t <= 100; t++ {
		input := genCase()
		exp, err := runRef(input)
		if err != nil {
			fmt.Println("reference failed:", err)
			os.Exit(1)
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("test %d exec failed: %v\n", t, err)
			os.Exit(1)
		}
		if got != exp {
			fmt.Printf("test %d failed\nexpected: %s\ngot: %s\n", t, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("ok")
}

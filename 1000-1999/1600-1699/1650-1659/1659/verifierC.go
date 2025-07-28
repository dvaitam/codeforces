package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
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
	cmd := exec.Command("go", "run", "1659C.go")
	cmd.Stdin = strings.NewReader(in)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(1)
	for t := 1; t <= 100; t++ {
		n := rand.Intn(10) + 1
		a := rand.Intn(20) + 1
		b := rand.Intn(20) + 1
		var sb strings.Builder
		sb.WriteString("1\n")
		fmt.Fprintf(&sb, "%d %d %d\n", n, a, b)
		cur := 0
		for i := 0; i < n; i++ {
			cur += rand.Intn(10) + 1
			fmt.Fprintf(&sb, "%d ", cur)
		}
		sb.WriteByte('\n')
		input := sb.String()
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
			fmt.Printf("test %d failed: expected %s got %s\ninput:%s\n", t, exp, got, strings.ReplaceAll(input, "\n", " "))
			os.Exit(1)
		}
	}
	fmt.Println("ok")
}

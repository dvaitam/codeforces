package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func runBinary(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func buildOfficial() (string, error) {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	bin := filepath.Join(dir, "officialC.bin")
	cmd := exec.Command("go", "build", "-o", bin, filepath.Join(dir, "1056C.go"))
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return bin, nil
}

func genTest() string {
	n := rand.Intn(3) + 1
	m := rand.Intn(n + 1)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < 2*n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", rand.Intn(100)+1))
	}
	sb.WriteByte('\n')
	used := make(map[int]bool)
	for i := 0; i < m; i++ {
		var a, b int
		for {
			a = rand.Intn(2*n) + 1
			if !used[a] {
				used[a] = true
				break
			}
		}
		for {
			b = rand.Intn(2*n) + 1
			if b != a && !used[b] {
				used[b] = true
				break
			}
		}
		sb.WriteString(fmt.Sprintf("%d %d\n", a, b))
	}
	t := rand.Intn(2) + 1
	sb.WriteString(fmt.Sprintf("%d\n", t))
	perm := rand.Perm(2 * n)
	for i, v := range perm {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v+1))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func main() {
	rand.Seed(44)
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
		return
	}
	candidate := os.Args[1]
	off, err := buildOfficial()
	if err != nil {
		fmt.Println("failed to build official solution:", err)
		return
	}
	defer os.Remove(off)

	for i := 0; i < 100; i++ {
		input := genTest()
		exp, err1 := runBinary(off, input)
		out, err2 := runBinary(candidate, input)
		if err1 != nil || err2 != nil {
			fmt.Printf("Runtime error on test %d\n", i+1)
			fmt.Println("input:\n" + input)
			if err1 != nil {
				fmt.Println("official:", err1)
			}
			if err2 != nil {
				fmt.Println("candidate:", err2)
			}
			return
		}
		if exp != out {
			fmt.Printf("Wrong answer on test %d\nInput:\n%sExpected:%s\nGot:%s\n", i+1, input, exp, out)
			return
		}
	}
	fmt.Println("All tests passed!")
}

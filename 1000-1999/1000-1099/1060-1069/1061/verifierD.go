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
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
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
	bin := filepath.Join(dir, "officialD.bin")
	cmd := exec.Command("go", "build", "-o", bin, filepath.Join(dir, "1061D.go"))
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return bin, nil
}

func genCase(rng *rand.Rand) (int, int64, int64, [][2]int) {
	n := rng.Intn(8) + 1
	x := int64(rng.Intn(20) + 2)
	y := int64(rng.Intn(int(x-1)) + 1)
	segs := make([][2]int, n)
	for i := 0; i < n; i++ {
		l := rng.Intn(20) + 1
		r := l + rng.Intn(20)
		segs[i] = [2]int{l, r}
	}
	return n, x, y, segs
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		return
	}
	candidate := os.Args[1]
	off, err := buildOfficial()
	if err != nil {
		fmt.Println("failed to build official solution:", err)
		return
	}
	defer os.Remove(off)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, x, y, segs := genCase(rng)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d %d\n", n, x, y))
		for _, s := range segs {
			sb.WriteString(fmt.Sprintf("%d %d\n", s[0], s[1]))
		}
		input := sb.String()
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

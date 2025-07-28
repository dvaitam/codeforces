package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func buildRef() (string, error) {
	bin := "refE.bin"
	cmd := exec.Command("go", "build", "-o", bin, "1954E.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference: %v\n%s", err, out)
	}
	return bin, nil
}

func runBin(path, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, path)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		return string(out), fmt.Errorf("time limit exceeded")
	}
	return string(out), err
}

func generateTests() []string {
	rand.Seed(5)
	res := make([]string, 0, 100)
	for len(res) < 100 {
		n := rand.Intn(5) + 1
		arr := make([]string, n)
		maxV := 0
		for i := range arr {
			v := rand.Intn(10) + 1
			arr[i] = fmt.Sprint(v)
			if v > maxV {
				maxV = v
			}
		}
		res = append(res, fmt.Sprintf("%d\n%s\n", n, strings.Join(arr, " ")))
	}
	return res
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	cand := os.Args[1]

	refBin, err := buildRef()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()

	for idx, tcase := range tests {
		expOut, err := runBin("./"+refBin, tcase)
		if err != nil {
			fmt.Println("reference error:", err)
			os.Exit(1)
		}
		candOut, err := runBin(cand, tcase)
		if err != nil {
			fmt.Printf("candidate error on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if !bytes.Equal(bytes.TrimSpace([]byte(expOut)), bytes.TrimSpace([]byte(candOut))) {
			fmt.Printf("test %d failed\n", idx+1)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed!")
}

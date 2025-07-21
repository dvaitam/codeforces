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
	out := "refE.bin"
	cmd := exec.Command("go", "build", "-o", out, "48E.go")
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return "./" + out, nil
}

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func runRef(ref, input string) (string, error) {
	cmd := exec.Command(ref)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func generateCase(r *rand.Rand) string {
	R := r.Intn(10) + 1
	h0 := r.Intn(R + 1)
	t0 := r.Intn(R + 1 - h0)
	n := r.Intn(3) + 1
	m := r.Intn(3) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", h0, t0, R))
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", r.Intn(R+1), r.Intn(R+1)))
	}
	sb.WriteString(fmt.Sprintf("%d\n", m))
	for i := 0; i < m; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", r.Intn(R+1), r.Intn(R+1)))
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	ref, err := compileRef()
	if err != nil {
		fmt.Println("failed to compile reference:", err)
		os.Exit(1)
	}
	defer os.Remove(ref)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in := generateCase(r)
		expect, err := runRef(ref, in)
		if err != nil {
			fmt.Printf("reference runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		out, err := runBinary(bin, in)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(expect) {
			fmt.Printf("test %d failed.\nInput:\n%sExpected:\n%s\nGot:\n%s\n", i+1, in, expect, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

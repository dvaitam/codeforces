package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func buildRef() (string, error) {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	src := filepath.Join(dir, "1872A.go")
	out := filepath.Join(os.TempDir(), "refA.bin")
	cmd := exec.Command("go", "build", "-o", out, src)
	if outb, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v\n%s", err, outb)
	}
	return out, nil
}

func prepareBin(path, tag string) (string, error) {
	if strings.HasSuffix(path, ".go") {
		out := filepath.Join(os.TempDir(), tag+fmt.Sprint(time.Now().UnixNano()))
		cmd := exec.Command("go", "build", "-o", out, path)
		if outb, err := cmd.CombinedOutput(); err != nil {
			return "", fmt.Errorf("build %s failed: %v\n%s", path, err, outb)
		}
		return out, nil
	}
	return path, nil
}

func runBin(path string, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func genCase() string {
	t := rand.Intn(5) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	for i := 0; i < t; i++ {
		a := rand.Intn(100) + 1
		b := rand.Intn(100) + 1
		c := rand.Intn(100) + 1
		sb.WriteString(fmt.Sprintf("%d %d %d\n", a, b, c))
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierA.go /path/to/binary")
		return
	}
	cand, err := prepareBin(os.Args[1], "candA")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 100; i++ {
		input := genCase()
		exp, err := runBin(ref, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference error on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runBin(cand, input)
		if err != nil {
			fmt.Printf("candidate runtime error on case %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if exp != got {
			fmt.Printf("case %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s", i+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

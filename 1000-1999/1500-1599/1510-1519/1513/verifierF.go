package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func compileRef(src string) (string, error) {
	tmp, err := os.CreateTemp("", "refF-")
	if err != nil {
		return "", err
	}
	tmp.Close()
	os.Remove(tmp.Name())
	cmd := exec.Command("go", "build", "-o", tmp.Name(), src)
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return tmp.Name(), nil
}

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func genTest(rng *rand.Rand) string {
	n := rng.Intn(8) + 2
	var sb strings.Builder
	fmt.Fprintln(&sb, n)
	a := make([]string, n)
	b := make([]string, n)
	for i := 0; i < n; i++ {
		a[i] = fmt.Sprintf("%d", rng.Intn(1000)+1)
	}
	for i := 0; i < n; i++ {
		b[i] = fmt.Sprintf("%d", rng.Intn(1000)+1)
	}
	fmt.Fprintln(&sb, strings.Join(a, " "))
	fmt.Fprintln(&sb, strings.Join(b, " "))
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	candidate, _ := filepath.Abs(os.Args[1])
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	ref, err := compileRef("1513F.go")
	if err != nil {
		fmt.Println("failed to compile reference:", err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	for i := 0; i < 100; i++ {
		input := genTest(rng)
		expOut, err := runBinary(ref, input)
		if err != nil {
			fmt.Println("reference failed:", err)
			os.Exit(1)
		}
		gotOut, err := runBinary(candidate, input)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(expOut) != strings.TrimSpace(gotOut) {
			fmt.Printf("test %d failed:\ninput:\n%s\nexpected:%s\ngot:%s\n", i+1, input, strings.TrimSpace(expOut), strings.TrimSpace(gotOut))
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}

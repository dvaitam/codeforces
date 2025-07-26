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
)

func buildReference() (string, error) {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	src := filepath.Join(dir, "803G.go")
	bin := filepath.Join(os.TempDir(), "ref803G.bin")
	cmd := exec.Command("go", "build", "-o", bin, src)
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return bin, nil
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &bytes.Buffer{}
	err := cmd.Run()
	return out.String(), err
}

func genCase() string {
	n := rand.Intn(5) + 1
	k := rand.Intn(5) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", rand.Intn(20)+1))
	}
	sb.WriteByte('\n')
	q := rand.Intn(5) + 1
	sb.WriteString(fmt.Sprintf("%d\n", q))
	for i := 0; i < q; i++ {
		typ := rand.Intn(2) + 1
		if typ == 1 {
			l := rand.Intn(n*k) + 1
			r := l + rand.Intn(n*k-l+1)
			x := rand.Intn(20) + 1
			sb.WriteString(fmt.Sprintf("1 %d %d %d\n", l, r, x))
		} else {
			l := rand.Intn(n*k) + 1
			r := l + rand.Intn(n*k-l+1)
			sb.WriteString(fmt.Sprintf("2 %d %d\n", l, r))
		}
	}
	return sb.String()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: verifierG.go <path-to-binary>")
		return
	}
	userBin := os.Args[1]
	rand.Seed(1)
	refBin, err := buildReference()
	if err != nil {
		fmt.Println("failed to build reference:", err)
		return
	}
	defer os.Remove(refBin)

	for i := 0; i < 100; i++ {
		input := genCase()
		want, err1 := runBinary(refBin, input)
		if err1 != nil {
			fmt.Println("reference solution failed:", err1)
			return
		}
		got, err2 := runBinary(userBin, input)
		if err2 != nil {
			fmt.Printf("test %d: runtime error: %v\n", i+1, err2)
			fmt.Println("input:\n" + input)
			return
		}
		if strings.TrimSpace(want) != strings.TrimSpace(got) {
			fmt.Printf("test %d failed\ninput:\n%sexpected:\n%s\nactual:\n%s\n", i+1, input, want, got)
			return
		}
	}
	fmt.Println("All tests passed")
}

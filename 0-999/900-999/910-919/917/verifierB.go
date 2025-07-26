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

type edge struct {
	v int
	u int
	c byte
}

func buildRef() (string, error) {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	src := filepath.Join(dir, "917B.go")
	bin := filepath.Join(os.TempDir(), "ref917B.bin")
	cmd := exec.Command("go", "build", "-o", bin, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("compile reference: %v\n%s", err, out)
	}
	return bin, nil
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(4) + 2 // 2..5
	edges := make([]edge, 0)
	for i := 1; i <= n; i++ {
		for j := i + 1; j <= n; j++ {
			if rng.Intn(2) == 0 {
				edges = append(edges, edge{i, j, byte('a' + rng.Intn(26))})
			}
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, len(edges)))
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d %c\n", e.v, e.u, e.c))
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		return
	}
	userBin := os.Args[1]

	rand.Seed(time.Now().UnixNano())

	ref, err := buildRef()
	if err != nil {
		fmt.Println("reference compile failed:", err)
		return
	}
	defer os.Remove(ref)

	for i := 0; i < 100; i++ {
		input := genCase(rand.New(rand.NewSource(int64(i) + time.Now().UnixNano())))
		want, err1 := runBinary(ref, input)
		if err1 != nil {
			fmt.Printf("reference failed on test %d: %v\n", i+1, err1)
			fmt.Println("input:\n" + input)
			return
		}
		got, err2 := runBinary(userBin, input)
		if err2 != nil {
			fmt.Printf("test %d runtime error: %v\n", i+1, err2)
			fmt.Println("input:\n" + input)
			return
		}
		if want != got {
			fmt.Printf("test %d failed\ninput:\n%sexpected:\n%s\nactual:\n%s\n", i+1, input, want, got)
			return
		}
	}
	fmt.Println("All tests passed")
}

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
	src := filepath.Join(dir, "917E.go")
	bin := filepath.Join(os.TempDir(), "ref917E.bin")
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

func genTree(rng *rand.Rand, n int) []string {
	edges := make([]string, n-1)
	for v := 2; v <= n; v++ {
		p := rng.Intn(v-1) + 1
		c := byte('a' + rng.Intn(26))
		edges[v-2] = fmt.Sprintf("%d %d %c", p, v, c)
	}
	return edges
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(6) + 2 // 2..7
	m := rng.Intn(3) + 1
	q := rng.Intn(4) + 1
	tree := genTree(rng, n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, q))
	for _, e := range tree {
		sb.WriteString(e)
		sb.WriteByte('\n')
	}
	for i := 0; i < m; i++ {
		l := rng.Intn(5) + 1
		word := make([]byte, l)
		for j := range word {
			word[j] = byte('a' + rng.Intn(26))
		}
		sb.WriteString(string(word))
		sb.WriteByte('\n')
	}
	for i := 0; i < q; i++ {
		a := rng.Intn(n) + 1
		b := rng.Intn(n) + 1
		for b == a {
			b = rng.Intn(n) + 1
		}
		k := rng.Intn(m) + 1
		sb.WriteString(fmt.Sprintf("%d %d %d\n", a, b, k))
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
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

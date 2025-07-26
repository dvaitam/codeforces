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
	src := filepath.Join(dir, "917C.go")
	bin := filepath.Join(os.TempDir(), "ref917C.bin")
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
	k := rng.Intn(4) + 1 // 1..4
	x := rng.Intn(k) + 1
	n := rng.Intn(40) + k
	q := 0
	if n-x > 0 {
		q = rng.Intn(min(5, n-x) + 1)
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d %d\n", x, k, n, q))
	for i := 0; i < k; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", rng.Intn(10)+1))
	}
	sb.WriteByte('\n')
	used := make(map[int]bool)
	for i := 0; i < q; i++ {
		p := rng.Intn(n-x) + x + 1
		for used[p] {
			p = rng.Intn(n-x) + x + 1
		}
		used[p] = true
		w := rng.Intn(21) - 10
		sb.WriteString(fmt.Sprintf("%d %d\n", p, w))
	}
	return sb.String()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
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

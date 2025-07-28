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

func runBinary(bin string, input []byte) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func buildRef() (string, error) {
	_, cur, _, _ := runtime.Caller(0)
	dir := filepath.Dir(cur)
	src := filepath.Join(dir, "1929F.go")
	refBin := filepath.Join(os.TempDir(), "1929F_ref.bin")
	cmd := exec.Command("go", "build", "-o", refBin, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v\n%s", err, string(out))
	}
	return refBin, nil
}

func genTree(n int) ([]int, []int) {
	left := make([]int, n+1)
	right := make([]int, n+1)
	for i := 1; i <= n; i++ {
		left[i] = -1
		right[i] = -1
	}
	avail := []int{1}
	for v := 2; v <= n; v++ {
		pIdx := rand.Intn(len(avail))
		p := avail[pIdx]
		if left[p] == -1 && (rand.Intn(2) == 0 || right[p] != -1) {
			left[p] = v
		} else {
			right[p] = v
		}
		if left[p] != -1 && right[p] != -1 {
			avail[pIdx] = avail[len(avail)-1]
			avail = avail[:len(avail)-1]
		}
		avail = append(avail, v)
	}
	return left, right
}

func genTest() []byte {
	n := rand.Intn(5) + 2
	C := rand.Int63n(20) + 1
	left, right := genTree(n)
	vals := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		if rand.Intn(2) == 0 {
			vals[i] = -1
		} else {
			vals[i] = rand.Int63n(C) + 1
		}
	}
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", n, C))
	for i := 1; i <= n; i++ {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", left[i], right[i], vals[i]))
	}
	return []byte(sb.String())
}

func main() {
	rand.Seed(time.Now().UnixNano())
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		return
	}
	cand := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	for i := 1; i <= 100; i++ {
		in := genTest()
		expected, err := runBinary(ref, in)
		if err != nil {
			fmt.Printf("reference failed on test %d: %v\n", i, err)
			os.Exit(1)
		}
		got, err := runBinary(cand, in)
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\n", i, err)
			os.Exit(1)
		}
		if strings.TrimSpace(expected) != strings.TrimSpace(got) {
			fmt.Printf("wrong answer on test %d\ninput:\n%sexpected:%s\ngot:%s\n", i, string(in), expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

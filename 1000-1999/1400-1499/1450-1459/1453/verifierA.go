package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
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
	src := filepath.Join(dir, "1453A.go")
	refBin := filepath.Join(os.TempDir(), "1453A_ref.bin")
	cmd := exec.Command("go", "build", "-o", refBin, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v\n%s", err, out)
	}
	return refBin, nil
}

func genTest() []byte {
	n := rand.Intn(10) + 1
	m := rand.Intn(10) + 1
	arr1 := rand.Perm(100)[:n]
	arr2 := rand.Perm(100)[:m]
	sort.Ints(arr1)
	sort.Ints(arr2)
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i, v := range arr1 {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v+1))
	}
	sb.WriteByte('\n')
	for i, v := range arr2 {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v+1))
	}
	sb.WriteByte('\n')
	return []byte(sb.String())
}

func main() {
	rand.Seed(time.Now().UnixNano())
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
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

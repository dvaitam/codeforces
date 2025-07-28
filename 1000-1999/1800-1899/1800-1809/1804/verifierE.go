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
	src := filepath.Join(dir, "1804E.go")
	refBin := filepath.Join(os.TempDir(), "1804E_ref.bin")
	cmd := exec.Command("go", "build", "-o", refBin, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v\n%s", err, string(out))
	}
	return refBin, nil
}

func genTest() []byte {
	n := rand.Intn(5) + 2
	maxM := n * (n - 1) / 2
	m := rand.Intn(maxM-(n-1)+1) + (n - 1)
	edges := make(map[[2]int]struct{})
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for len(edges) < m {
		u := rand.Intn(n) + 1
		v := rand.Intn(n) + 1
		if u == v {
			continue
		}
		if u > v {
			u, v = v, u
		}
		pair := [2]int{u, v}
		if _, ok := edges[pair]; ok {
			continue
		}
		edges[pair] = struct{}{}
		sb.WriteString(fmt.Sprintf("%d %d\n", u, v))
	}
	return []byte(sb.String())
}

func main() {
	rand.Seed(time.Now().UnixNano())
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
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

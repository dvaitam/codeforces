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
	src := filepath.Join(dir, "1167D.go")
	refBin := filepath.Join(os.TempDir(), "1167D_ref.bin")
	cmd := exec.Command("go", "build", "-o", refBin, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v\n%s", err, string(out))
	}
	return refBin, nil
}

func isValidRBS(s []byte) bool {
	bal := 0
	for _, ch := range s {
		if ch == '(' {
			bal++
		} else {
			bal--
			if bal < 0 {
				return false
			}
		}
	}
	return bal == 0
}

func randomRBS(n int) string {
	for {
		arr := make([]byte, n)
		for i := 0; i < n/2; i++ {
			arr[i] = '('
		}
		for i := n / 2; i < n; i++ {
			arr[i] = ')'
		}
		rand.Shuffle(n, func(i, j int) { arr[i], arr[j] = arr[j], arr[i] })
		if isValidRBS(arr) {
			return string(arr)
		}
	}
}

func genTest() []byte {
	n := (rand.Intn(10) + 1) * 2
	s := randomRBS(n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n%s\n", n, s))
	return []byte(sb.String())
}

func main() {
	rand.Seed(time.Now().UnixNano())
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
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

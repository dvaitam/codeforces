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

func compileRef() (string, error) {
	dir := filepath.Dir(os.Args[0])
	src := filepath.Join(dir, "1037G.go")
	out := filepath.Join(os.TempDir(), "ref1037G.bin")
	cmd := exec.Command("go", "build", "-o", out, src)
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return out, nil
}

func runCmd(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func genCase(r *rand.Rand) (string, string) {
	n := r.Intn(5) + 1
	letters := []byte{'a', 'b', 'c'}
	s := make([]byte, n)
	for i := range s {
		s[i] = letters[r.Intn(len(letters))]
	}
	m := r.Intn(5) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s\n%d\n", s, m))
	for i := 0; i < m; i++ {
		l := r.Intn(n) + 1
		rpos := r.Intn(n-l+1) + l
		sb.WriteString(fmt.Sprintf("%d %d\n", l, rpos))
	}
	input := sb.String()
	refBin, err := compileRef()
	if err != nil {
		panic(err)
	}
	exp, err := runCmd(refBin, input)
	if err != nil {
		panic(err)
	}
	os.Remove(refBin)
	return input, exp
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		out, err := runCmd(bin, in)
		if err != nil {
			fmt.Printf("Test %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(exp) {
			fmt.Printf("Test %d failed.\nInput:\n%sExpected:\n%s\nGot:\n%s\n", i+1, in, exp, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

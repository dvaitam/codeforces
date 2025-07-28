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

func buildIfGo(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "solbin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		if out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput(); err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("build failed: %v\n%s", err, string(out))
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func runProg(exe string, input []byte) (string, error) {
	cmd := exec.Command(exe)
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return out.String() + errBuf.String(), err
	}
	return strings.TrimSpace(out.String()), nil
}

func buildRef() (string, error) {
	_, cur, _, _ := runtime.Caller(0)
	dir := filepath.Dir(cur)
	src := filepath.Join(dir, "1608E.go")
	refBin := filepath.Join(os.TempDir(), "1608E_ref.bin")
	cmd := exec.Command("go", "build", "-o", refBin, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v\n%s", err, string(out))
	}
	return refBin, nil
}

func genTest() []byte {
	m := rand.Intn(5) + 1
	n := m * 3
	counts := []int{m, m, m}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	cells := make([][3]int, 0, n)
	for c := 0; c < 3; c++ {
		for i := 0; i < counts[c]; i++ {
			x := rand.Intn(21) - 10
			y := rand.Intn(21) - 10
			cells = append(cells, [3]int{x, y, c + 1})
		}
	}
	rand.Shuffle(len(cells), func(i, j int) { cells[i], cells[j] = cells[j], cells[i] })
	for _, v := range cells {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", v[0], v[1], v[2]))
	}
	return []byte(sb.String())
}

func main() {
	rand.Seed(time.Now().UnixNano())
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		return
	}
	exe, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer cleanup()
	ref, err := buildRef()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer os.Remove(ref)
	for i := 1; i <= 100; i++ {
		in := genTest()
		expected, err := runProg(ref, in)
		if err != nil {
			fmt.Printf("reference failed on test %d: %v\n%s", i, err, expected)
			os.Exit(1)
		}
		got, err := runProg(exe, in)
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\n%s", i, err, got)
			os.Exit(1)
		}
		if expected != got {
			fmt.Printf("wrong answer on test %d\ninput:\n%sexpected:%s\ngot:%s\n", i, string(in), expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

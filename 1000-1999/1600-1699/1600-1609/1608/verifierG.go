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
	src := filepath.Join(dir, "1608G.go")
	refBin := filepath.Join(os.TempDir(), "1608G_ref.bin")
	cmd := exec.Command("go", "build", "-o", refBin, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v\n%s", err, string(out))
	}
	return refBin, nil
}

func genTree(n int) [][3]interface{} {
	edges := make([][3]interface{}, 0, n-1)
	for i := 2; i <= n; i++ {
		p := rand.Intn(i-1) + 1
		c := byte('a' + rand.Intn(3))
		edges = append(edges, [3]interface{}{p, i, c})
	}
	return edges
}

func genTest() []byte {
	n := rand.Intn(4) + 2
	m := rand.Intn(4) + 1
	q := rand.Intn(4) + 1
	edges := genTree(n)
	strs := make([]string, m)
	for i := 0; i < m; i++ {
		l := rand.Intn(3) + 1
		b := make([]byte, l)
		for j := 0; j < l; j++ {
			b[j] = byte('a' + rand.Intn(3))
		}
		strs[i] = string(b)
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, q))
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d %c\n", e[0], e[1], e[2]))
	}
	for i := 0; i < m; i++ {
		sb.WriteString(strs[i] + "\n")
	}
	for i := 0; i < q; i++ {
		u := rand.Intn(n) + 1
		v := rand.Intn(n) + 1
		for v == u {
			v = rand.Intn(n) + 1
		}
		l := rand.Intn(m) + 1
		r := rand.Intn(m-l+1) + l
		sb.WriteString(fmt.Sprintf("%d %d %d %d\n", u, v, l, r))
	}
	return []byte(sb.String())
}

func main() {
	rand.Seed(time.Now().UnixNano())
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
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

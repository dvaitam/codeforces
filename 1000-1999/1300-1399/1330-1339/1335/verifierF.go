package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func buildRef() (string, error) {
	ref := "./refF.bin"
	cmd := exec.Command("go", "build", "-o", ref, "1335F.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v: %s", err, string(out))
	}
	return ref, nil
}

func runExe(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func genTest() string {
	t := 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	n := rand.Intn(5) + 1
	m := rand.Intn(5) + 1
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < n; i++ {
		line := make([]byte, m)
		for j := 0; j < m; j++ {
			if rand.Intn(2) == 0 {
				line[j] = '0'
			} else {
				line[j] = '1'
			}
		}
		sb.WriteString(string(line))
		sb.WriteByte('\n')
	}
	for i := 0; i < n; i++ {
		line := make([]byte, m)
		for j := 0; j < m; j++ {
			dirs := []byte{'U', 'R', 'D', 'L'}
			if i == 0 {
				dirs = remove(dirs, 'U')
			}
			if i == n-1 {
				dirs = remove(dirs, 'D')
			}
			if j == 0 {
				dirs = remove(dirs, 'L')
			}
			if j == m-1 {
				dirs = remove(dirs, 'R')
			}
			line[j] = dirs[rand.Intn(len(dirs))]
		}
		sb.WriteString(string(line))
		sb.WriteByte('\n')
	}
	return sb.String()
}

func remove(a []byte, x byte) []byte {
	out := a[:0]
	for _, v := range a {
		if v != x {
			out = append(out, v)
		}
	}
	return out
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	rand.Seed(7)
	for i := 0; i < 100; i++ {
		input := genTest()
		exp, err := runExe(ref, input)
		if err != nil {
			fmt.Printf("reference runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runExe(bin, input)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", i+1, err)
			fmt.Println("input:\n", input)
			os.Exit(1)
		}
		if strings.TrimSpace(exp) != strings.TrimSpace(got) {
			fmt.Printf("test %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}

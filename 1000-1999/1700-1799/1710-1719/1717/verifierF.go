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

func baseDir() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Dir(file)
}

func prepareBinary(path, tag string) (string, error) {
	if strings.HasSuffix(path, ".go") {
		bin := filepath.Join(os.TempDir(), tag+"_"+fmt.Sprint(time.Now().UnixNano()))
		cmd := exec.Command("go", "build", "-o", bin, path)
		if out, err := cmd.CombinedOutput(); err != nil {
			return "", fmt.Errorf("build %s: %v\n%s", path, err, out)
		}
		return bin, nil
	}
	return path, nil
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(5) + 2 // 2..6
	m := rng.Intn(5) + 1
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	s := make([]int, n)
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		if rng.Intn(2) == 0 {
			sb.WriteByte('0')
			s[i] = 0
		} else {
			sb.WriteByte('1')
			s[i] = 1
		}
	}
	sb.WriteByte('\n')
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		if s[i] == 0 {
			sb.WriteByte('0')
		} else {
			sb.WriteString(fmt.Sprintf("%d", rng.Intn(2*m+1)-m))
		}
	}
	sb.WriteByte('\n')
	used := make(map[[2]int]bool)
	for i := 0; i < m; i++ {
		for {
			u := rng.Intn(n) + 1
			v := rng.Intn(n) + 1
			if u == v {
				continue
			}
			p1 := [2]int{u, v}
			p2 := [2]int{v, u}
			if used[p1] || used[p2] {
				continue
			}
			used[p1] = true
			sb.WriteString(fmt.Sprintf("%d %d\n", u, v))
			break
		}
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	candPath, err := prepareBinary(os.Args[1], "candF")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	refSrc := filepath.Join(baseDir(), "1717F.go")
	refPath, err := prepareBinary(refSrc, "refF")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := genCase(rng)
		exp, err := runBinary(refPath, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference error on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runBinary(candPath, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on case %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(exp) != strings.TrimSpace(got) {
			fmt.Printf("case %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s", i+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

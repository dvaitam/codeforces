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
	src := filepath.Join(dir, "1872E.go")
	out := filepath.Join(os.TempDir(), "refE.bin")
	cmd := exec.Command("go", "build", "-o", out, src)
	if outb, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v\n%s", err, outb)
	}
	return out, nil
}

func prepareBin(path, tag string) (string, error) {
	if strings.HasSuffix(path, ".go") {
		out := filepath.Join(os.TempDir(), tag+fmt.Sprint(time.Now().UnixNano()))
		cmd := exec.Command("go", "build", "-o", out, path)
		if outb, err := cmd.CombinedOutput(); err != nil {
			return "", fmt.Errorf("build %s failed: %v\n%s", path, err, outb)
		}
		return out, nil
	}
	return path, nil
}

func runBin(path string, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func genCase() string {
	t := rand.Intn(3) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	for ; t > 0; t-- {
		n := rand.Intn(6) + 1
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", rand.Intn(100)))
		}
		sb.WriteByte('\n')
		var s strings.Builder
		for i := 0; i < n; i++ {
			if rand.Intn(2) == 0 {
				s.WriteByte('0')
			} else {
				s.WriteByte('1')
			}
		}
		sb.WriteString(s.String())
		sb.WriteByte('\n')
		q := rand.Intn(6) + 1
		sb.WriteString(fmt.Sprintf("%d\n", q))
		has := false
		for i := 0; i < q; i++ {
			tp := rand.Intn(2) + 1
			if i == q-1 && !has {
				tp = 2
			}
			if tp == 1 {
				l := rand.Intn(n) + 1
				r := rand.Intn(n-l+1) + l
				sb.WriteString(fmt.Sprintf("1 %d %d\n", l, r))
			} else {
				g := rand.Intn(2)
				sb.WriteString(fmt.Sprintf("2 %d\n", g))
				has = true
			}
		}
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierE.go /path/to/binary")
		return
	}
	cand, err := prepareBin(os.Args[1], "candE")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 100; i++ {
		input := genCase()
		exp, err := runBin(ref, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference error on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runBin(cand, input)
		if err != nil {
			fmt.Printf("candidate runtime error on case %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if exp != got {
			fmt.Printf("case %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s", i+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

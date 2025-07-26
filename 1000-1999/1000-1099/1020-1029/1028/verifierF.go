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
	_, cur, _, _ := runtime.Caller(0)
	dir := filepath.Dir(cur)
	src := filepath.Join(dir, "solF.cpp")
	bin := filepath.Join(os.TempDir(), "ref1028F.bin")
	cmd := exec.Command("g++", "-std=c++17", "-O2", "-pipe", "-march=native", "-o", bin, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v\n%s", err, string(out))
	}
	return bin, nil
}

func run(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(r *rand.Rand) string {
	q := r.Intn(20) + 1
	type pt struct{ x, y int }
	points := make([]pt, 0)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", q))
	for i := 0; i < q; i++ {
		t := r.Intn(3) + 1
		var x, y int
		if t == 1 || t == 2 {
			x = r.Intn(11) - 5
			y = r.Intn(11) - 5
			sb.WriteString(fmt.Sprintf("%d %d %d\n", t, x, y))
			if t == 1 {
				points = append(points, pt{x, y})
			}
			if t == 2 && len(points) > 0 {
				points = points[1:]
			}
		} else {
			x = r.Intn(11) - 5
			y = r.Intn(11) - 5
			sb.WriteString(fmt.Sprintf("3 %d %d\n", x, y))
		}
	}
	return sb.String()
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
		in := genCase(rand.New(rand.NewSource(int64(i))))
		want, err := run(ref, in)
		if err != nil {
			fmt.Printf("reference failed on test %d: %v\n", i, err)
			os.Exit(1)
		}
		got, err := run(cand, in)
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\n", i, err)
			os.Exit(1)
		}
		if strings.TrimSpace(want) != strings.TrimSpace(got) {
			fmt.Printf("wrong answer on test %d\ninput:\n%sexpected:%s\ngot:%s\n", i, in, want, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

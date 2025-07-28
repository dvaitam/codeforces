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

func buildBinary(src, tag string) (string, error) {
	if strings.HasSuffix(src, ".go") {
		out := filepath.Join(os.TempDir(), tag)
		cmd := exec.Command("go", "build", "-o", out, src)
		if outb, err := cmd.CombinedOutput(); err != nil {
			return "", fmt.Errorf("build %s: %v\n%s", src, err, string(outb))
		}
		return out, nil
	}
	return src, nil
}

func runCase(exe, input, expected string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(expected)
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func solveCase(arr []int) (int, int) {
	n := len(arr)
	p, l, r := 0, 0, 0
	j := -1
	for i := 0; i <= n; i++ {
		if i == n || arr[i] == 0 {
			mn0x, mn0y := 0, j+1
			mn1x, mn1y := n, -1
			pw, sign := 0, 0
			for k := j + 1; k < i; k++ {
				if arr[k] < 0 {
					sign ^= 1
				}
				if abs(arr[k]) == 2 {
					pw++
				}
				if sign == 0 {
					if pw-mn0x > p {
						p = pw - mn0x
						l = mn0y
						r = k + 1
					}
				} else {
					if pw-mn1x > p {
						p = pw - mn1x
						l = mn1y
						r = k + 1
					}
				}
				if sign == 0 {
					if pw < mn0x {
						mn0x = pw
						mn0y = k + 1
					}
				} else {
					if pw < mn1x {
						mn1x = pw
						mn1y = k + 1
					}
				}
			}
			j = i
		}
	}
	return l, n - r
}

func generateCase(r *rand.Rand) (string, string) {
	n := r.Intn(20) + 1
	arr := make([]int, n)
	for i := range arr {
		arr[i] = r.Intn(5) - 2 // -2..2
	}
	var b strings.Builder
	fmt.Fprintf(&b, "1\n%d\n", n)
	for i := 0; i < n; i++ {
		if i > 0 {
			b.WriteByte(' ')
		}
		fmt.Fprintf(&b, "%d", arr[i])
	}
	b.WriteByte('\n')
	l, y := solveCase(arr)
	expect := fmt.Sprintf("%d %d\n", l, y)
	return b.String(), expect
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		return
	}
	candSrc := os.Args[1]
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	refSrc := filepath.Join(dir, "1660D.go")

	cand, err := buildBinary(candSrc, "candD.bin")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	ref, err := buildBinary(refSrc, "refD.bin")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(ref, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "reference failed: %v\n", err)
			os.Exit(1)
		}
		if err := runCase(cand, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func buildRef() (string, error) {
	_, cur, _, _ := runtime.Caller(0)
	dir := filepath.Dir(cur)
	src := filepath.Join(dir, "1028E.go")
	bin := filepath.Join(os.TempDir(), "ref1028E.bin")
	cmd := exec.Command("go", "build", "-o", bin, src)
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
	n := r.Intn(5) + 2
	a := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = r.Intn(20) + 1
	}
	b := make([]int, n)
	for i := 0; i < n; i++ {
		b[i] = a[i] % a[(i+1)%n]
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range b {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func check(input, output string) error {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) < 2 || strings.ToUpper(strings.TrimSpace(lines[0])) != "YES" {
		return fmt.Errorf("expected YES")
	}
	parts := strings.Fields(lines[1])
	nfields := len(strings.Fields(strings.Split(input, "\n")[1]))
	if len(parts) != nfields {
		return fmt.Errorf("wrong number of elements")
	}
	bStr := strings.Fields(strings.Split(input, "\n")[1])
	b := make([]int, len(parts))
	for i, s := range bStr {
		v, _ := strconv.Atoi(s)
		b[i] = v
	}
	a := make([]int64, len(parts))
	for i, s := range parts {
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		}
		if v < 1 {
			return fmt.Errorf("a_i < 1")
		}
		a[i] = v
	}
	for i := 0; i < len(a); i++ {
		if int(a[i]%a[(i+1)%len(a)]) != b[i] {
			return fmt.Errorf("output does not satisfy condition")
		}
	}
	return nil
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
		if strings.ToUpper(strings.TrimSpace(want)) == "NO" {
			if strings.ToUpper(strings.TrimSpace(got)) != "NO" {
				fmt.Printf("wrong answer on test %d expected NO\n", i)
				os.Exit(1)
			}
			continue
		}
		if err := check(in, got); err != nil {
			fmt.Printf("wrong answer on test %d: %v\ninput:%soutput:%s\n", i, err, in, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

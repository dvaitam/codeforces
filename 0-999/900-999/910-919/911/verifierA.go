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

func runProg(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func runRef(input string) (string, error) {
	_, self, _, _ := runtime.Caller(0)
	dir := filepath.Dir(self)
	ref := filepath.Join(dir, "911A.go")
	return runProg(ref, input)
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(20) + 2 // 2..21
	minVal := rng.Intn(50) + 1
	arr := make([]int, n)
	for i := range arr {
		arr[i] = minVal + rng.Intn(50)
	}
	// ensure at least two minimums
	i := rng.Intn(n)
	j := rng.Intn(n)
	for j == i {
		j = rng.Intn(n)
	}
	arr[i] = minVal
	arr[j] = minVal
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(n))
	sb.WriteByte('\n')
	for idx, v := range arr {
		if idx > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in := genCase(rng)
		expect, err := runRef(in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runProg(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on case %d: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Printf("case %d failed\ninput:\n%sexpected:%s\nactual:%s\n", i+1, in, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}

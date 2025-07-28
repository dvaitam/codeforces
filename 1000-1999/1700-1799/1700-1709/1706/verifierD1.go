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
	return strings.TrimSpace(out.String()), err
}

func runRef(input string) (string, error) {
	_, self, _, _ := runtime.Caller(0)
	dir := filepath.Dir(self)
	ref := filepath.Join(dir, "1706D1.go")
	return runProg(ref, input)
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(5) + 1
	k := rng.Intn(5) + 1
	arr := make([]int, n)
	val := rng.Intn(5) + 1
	for i := 0; i < n; i++ {
		arr[i] = val
		val += rng.Intn(5)
		if val > 3000 {
			val = 3000
		}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, k)
	for i, v := range arr {
		if i+1 == n {
			fmt.Fprintf(&sb, "%d\n", v)
		} else {
			fmt.Fprintf(&sb, "%d ", v)
		}
	}
	return sb.String()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	const T = 100
	var input strings.Builder
	fmt.Fprintf(&input, "%d\n", T)
	cases := make([]string, T)
	for i := 0; i < T; i++ {
		c := genCase(rng)
		cases[i] = c
		input.WriteString(c)
	}
	expect, err := runRef(input.String())
	if err != nil {
		fmt.Fprintln(os.Stderr, "reference solver failed:", err)
		os.Exit(1)
	}
	got, err := runProg(bin, input.String())
	if err != nil {
		fmt.Fprintln(os.Stderr, "candidate failed:", err)
		os.Exit(1)
	}
	if strings.TrimSpace(got) != strings.TrimSpace(expect) {
		fmt.Printf("output mismatch\ninput:\n%s\nexpected:\n%s\nactual:\n%s\n", input.String(), expect, got)
		os.Exit(1)
	}
	fmt.Println("all tests passed")
}

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
	ref := filepath.Join(dir, "911G.go")
	return runProg(ref, input)
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(10) + 1
	arr := make([]int, n)
	for i := range arr {
		arr[i] = rng.Intn(100) + 1
	}
	q := rng.Intn(10) + 1
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(n))
	sb.WriteByte('\n')
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	sb.WriteString(strconv.Itoa(q))
	sb.WriteByte('\n')
	for i := 0; i < q; i++ {
		l := rng.Intn(n) + 1
		r := rng.Intn(n) + 1
		if l > r {
			l, r = r, l
		}
		x := rng.Intn(100) + 1
		y := rng.Intn(100) + 1
		sb.WriteString(fmt.Sprintf("%d %d %d %d\n", l, r, x, y))
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
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

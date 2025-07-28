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

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	ref := filepath.Join(dir, "refF")
	cmd := exec.Command("go", "build", "-o", ref, filepath.Join(dir, "1810F.go"))
	if out, err := cmd.CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n%s", err, out)
		os.Exit(1)
	}
	defer os.Remove(ref)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		n := rng.Intn(5) + 1
		m := rng.Intn(3) + 1
		q := rng.Intn(5) + 1
		a := make([]int, n)
		for i := range a {
			a[i] = rng.Intn(5) + 1
		}
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, q))
		for i, v := range a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(v))
		}
		sb.WriteByte('\n')
		for i := 0; i < q; i++ {
			x := rng.Intn(n) + 1
			y := rng.Intn(5) + 1
			sb.WriteString(fmt.Sprintf("%d %d\n", x, y))
		}
		input := sb.String()
		candOut, cErr := runBinary(candidate, input)
		refOut, rErr := runBinary(ref, input)
		if cErr != nil {
			fmt.Fprintf(os.Stderr, "case %d candidate error: %v\n%s", t+1, cErr, candOut)
			os.Exit(1)
		}
		if rErr != nil {
			fmt.Fprintf(os.Stderr, "case %d reference error: %v\n%s", t+1, rErr, refOut)
			os.Exit(1)
		}
		if strings.TrimSpace(candOut) != strings.TrimSpace(refOut) {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%sexpected:%sactual:%s", t+1, input, refOut, candOut)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

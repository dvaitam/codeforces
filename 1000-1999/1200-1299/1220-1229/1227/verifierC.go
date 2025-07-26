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

func genCase(rng *rand.Rand) string {
	n := 2 * (rng.Intn(4) + 2) // length 4..10 even
	k := rng.Intn(n/2) + 1
	arr := make([]byte, n)
	for i := 0; i < n/2; i++ {
		arr[i] = '('
	}
	for i := n / 2; i < n; i++ {
		arr[i] = ')'
	}
	rng.Shuffle(n, func(i, j int) { arr[i], arr[j] = arr[j], arr[i] })
	s := string(arr)
	var b bytes.Buffer
	fmt.Fprintln(&b, 1)
	fmt.Fprintf(&b, "%d %d\n", n, k)
	fmt.Fprintln(&b, s)
	return b.String()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: go run verifierC.go /path/to/binary\n")
		os.Exit(1)
	}
	candidate := os.Args[1]

	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	ref := filepath.Join(dir, "refC")
	if out, err := exec.Command("go", "build", "-o", ref, filepath.Join(dir, "1227C.go")).CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference solution: %v\n%s", err, out)
		os.Exit(1)
	}
	defer os.Remove(ref)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		input := genCase(rng)
		candOut, cErr := runBinary(candidate, input)
		refOut, rErr := runBinary(ref, input)
		if cErr != nil {
			fmt.Fprintf(os.Stderr, "case %d: candidate error: %v\ninput:\n%s", t+1, cErr, input)
			os.Exit(1)
		}
		if rErr != nil {
			fmt.Fprintf(os.Stderr, "case %d: reference error: %v\ninput:\n%s", t+1, rErr, input)
			os.Exit(1)
		}
		if candOut != refOut {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%s\nexpected:%s\nactual:%s\n", t+1, input, refOut, candOut)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

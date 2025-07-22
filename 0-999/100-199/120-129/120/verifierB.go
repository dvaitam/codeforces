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
	return out.String(), err
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(20) + 1
	k := rng.Intn(n) + 1
	a := make([]int, n)
	one := false
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			a[i] = 1
			one = true
		}
	}
	if !one {
		idx := rng.Intn(n)
		a[idx] = 1
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, k)
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", a[i])
	}
	sb.WriteByte('\n')
	// compute expected
	pos := k - 1
	var ans int
	for i := pos; i < n; i++ {
		if a[i] == 1 {
			ans = i + 1
			break
		}
	}
	if ans == 0 {
		for i := 0; i < pos; i++ {
			if a[i] == 1 {
				ans = i + 1
				break
			}
		}
	}
	return sb.String(), fmt.Sprintf("%d\n", ans)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: go run verifierB.go /path/to/binary\n")
		os.Exit(1)
	}
	candidate := os.Args[1]

	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	ref := filepath.Join(dir, "refB")
	cmd := exec.Command("go", "build", "-o", ref, filepath.Join(dir, "120B.go"))
	if out, err := cmd.CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference solution: %v\n%s", err, out)
		os.Exit(1)
	}
	defer os.Remove(ref)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		input, expect := generateCase(rng)
		candOut, cErr := runBinary(candidate, input)
		refOut, rErr := runBinary(ref, input)
		if cErr != nil {
			fmt.Fprintf(os.Stderr, "test %d: candidate error: %v\n", t+1, cErr)
			os.Exit(1)
		}
		if rErr != nil {
			fmt.Fprintf(os.Stderr, "test %d: reference error: %v\n", t+1, rErr)
			os.Exit(1)
		}
		if strings.TrimSpace(candOut) != strings.TrimSpace(refOut) {
			fmt.Fprintf(os.Stderr, "test %d failed\ninput:\n%sexpected:%sactual:%s\n", t+1, input, refOut, candOut)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

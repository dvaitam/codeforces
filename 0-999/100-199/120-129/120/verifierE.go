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

func expected(nums []int) string {
	var sb strings.Builder
	for i, n := range nums {
		v := 1
		if n%2 == 1 {
			v = 0
		}
		if i > 0 {
			sb.WriteByte('\n')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	return sb.String()
}

func generateCase(rng *rand.Rand) (string, string) {
	T := 100
	nums := make([]int, T)
	for i := 0; i < T; i++ {
		nums[i] = rng.Intn(10000) + 1
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", T)
	for i := 0; i < T; i++ {
		fmt.Fprintf(&sb, "%d\n", nums[i])
	}
	return sb.String(), expected(nums)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: go run verifierE.go /path/to/binary\n")
		os.Exit(1)
	}
	candidate := os.Args[1]

	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	ref := filepath.Join(dir, "refE")
	cmd := exec.Command("go", "build", "-o", ref, filepath.Join(dir, "120E.go"))
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

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

type jar struct {
	rem   int
	eaten int
}

func expected(n, k int, arr []int) string {
	jars := make([]jar, n)
	for i := 0; i < n; i++ {
		jars[i] = jar{rem: arr[i], eaten: 0}
	}
	total := 0
	for len(jars) > 0 {
		maxIdx := 0
		for i := 1; i < len(jars); i++ {
			if jars[i].rem > jars[maxIdx].rem {
				maxIdx = i
			}
		}
		j := &jars[maxIdx]
		if j.rem < k || j.eaten >= 3 {
			total += j.rem
			jars = append(jars[:maxIdx], jars[maxIdx+1:]...)
		} else {
			j.rem -= k
			j.eaten++
		}
	}
	return fmt.Sprintf("%d\n", total)
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 1
	k := rng.Intn(10) + 1
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = rng.Intn(20) + 1
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, k)
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", arr[i])
	}
	sb.WriteByte('\n')
	return sb.String(), expected(n, k, arr)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: go run verifierC.go /path/to/binary\n")
		os.Exit(1)
	}
	candidate := os.Args[1]

	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	ref := filepath.Join(dir, "refC")
	cmd := exec.Command("go", "build", "-o", ref, filepath.Join(dir, "120C.go"))
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

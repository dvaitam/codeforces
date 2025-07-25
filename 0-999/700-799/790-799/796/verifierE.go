package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
)

func runBinary(binPath, input string) (string, error) {
	cmd := exec.Command(binPath)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	refSrc := filepath.Join(dir, "796E.go")
	refBin := filepath.Join(os.TempDir(), "796E_ref_bin")
	cmd := exec.Command("go", "build", "-o", refBin, refSrc)
	if out, err := cmd.CombinedOutput(); err != nil {
		fmt.Println("Failed to build reference:", err)
		fmt.Println(string(out))
		os.Exit(1)
	}
	defer os.Remove(refBin)

	rng := rand.New(rand.NewSource(46))

	for tc := 1; tc <= 100; tc++ {
		n := rng.Intn(20) + 1
		p := rng.Intn(20) + 1
		k := rng.Intn(min(n, 5)) + 1
		r := rng.Intn(n + 1)
		a := randIndices(n, r, rng)
		s := rng.Intn(n + 1)
		b := randIndices(n, s, rng)

		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d %d\n", n, p, k))
		sb.WriteString(fmt.Sprintf("%d", r))
		for _, v := range a {
			sb.WriteString(fmt.Sprintf(" %d", v))
		}
		sb.WriteByte('\n')
		sb.WriteString(fmt.Sprintf("%d", s))
		for _, v := range b {
			sb.WriteString(fmt.Sprintf(" %d", v))
		}
		sb.WriteByte('\n')
		input := sb.String()

		exp, err := runBinary(refBin, input)
		if err != nil {
			fmt.Println("Reference solution error on test", tc)
			fmt.Println(err)
			os.Exit(1)
		}
		got, err := runBinary(candidate, input)
		if err != nil {
			fmt.Printf("Candidate runtime error on test %d: %v\n", tc, err)
			os.Exit(1)
		}
		if strings.TrimSpace(exp) != strings.TrimSpace(got) {
			fmt.Printf("Wrong answer on test %d\nInput:\n%sExpected:%sGot:%s\n", tc, input, exp, got)
			os.Exit(1)
		}
	}

	fmt.Println("All tests passed.")
}

func randIndices(n, cnt int, rng *rand.Rand) []int {
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = i + 1
	}
	rng.Shuffle(n, func(i, j int) { arr[i], arr[j] = arr[j], arr[i] })
	res := arr[:cnt]
	sort.Ints(res)
	return res
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

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
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	refSrc := filepath.Join(dir, "796B.go")
	refBin := filepath.Join(os.TempDir(), "796B_ref_bin")
	cmd := exec.Command("go", "build", "-o", refBin, refSrc)
	if out, err := cmd.CombinedOutput(); err != nil {
		fmt.Println("Failed to build reference:", err)
		fmt.Println(string(out))
		os.Exit(1)
	}
	defer os.Remove(refBin)

	rng := rand.New(rand.NewSource(43))

	for tc := 1; tc <= 100; tc++ {
		n := rng.Intn(20) + 2
		m := rng.Intn(n-1) + 1
		k := rng.Intn(20) + 1
		holesSet := make(map[int]bool)
		for len(holesSet) < m {
			holesSet[rng.Intn(n)+1] = true
		}
		holes := make([]int, 0, m)
		for h := range holesSet {
			holes = append(holes, h)
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, k))
		for i, h := range holes {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", h))
		}
		sb.WriteByte('\n')
		for i := 0; i < k; i++ {
			u := rng.Intn(n) + 1
			v := rng.Intn(n) + 1
			if u == v {
				v = (v % n) + 1
			}
			sb.WriteString(fmt.Sprintf("%d %d\n", u, v))
		}
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

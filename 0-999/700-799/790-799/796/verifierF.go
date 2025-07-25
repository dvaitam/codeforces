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
		fmt.Println("Usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	refSrc := filepath.Join(dir, "796F.go")
	refBin := filepath.Join(os.TempDir(), "796F_ref_bin")
	cmd := exec.Command("go", "build", "-o", refBin, refSrc)
	if out, err := cmd.CombinedOutput(); err != nil {
		fmt.Println("Failed to build reference:", err)
		fmt.Println(string(out))
		os.Exit(1)
	}
	defer os.Remove(refBin)

	rng := rand.New(rand.NewSource(47))

	for tc := 1; tc <= 100; tc++ {
		n := rng.Intn(6) + 1
		m := rng.Intn(8) + 1
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		usedX := make(map[int]bool)
		for i := 0; i < m; i++ {
			t := rng.Intn(2) + 1
			if t == 1 {
				l := rng.Intn(n) + 1
				r := rng.Intn(n-l+1) + l
				x := rng.Intn(100)
				for usedX[x] {
					x = rng.Intn(100)
				}
				usedX[x] = true
				sb.WriteString(fmt.Sprintf("1 %d %d %d\n", l, r, x))
			} else {
				k := rng.Intn(n) + 1
				d := rng.Intn(100)
				sb.WriteString(fmt.Sprintf("2 %d %d\n", k, d))
			}
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

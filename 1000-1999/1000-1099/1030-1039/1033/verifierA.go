package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func buildIfGo(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "solbin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		if out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput(); err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("build failed: %v\n%s", err, string(out))
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func run(bin string, input []byte) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(bytes.TrimSpace(out)), err
}

func region(x, y, qx, qy int) int {
	ans := 0
	if x < qx {
		ans = 1
	}
	ans <<= 1
	if y < qy {
		ans |= 1
	}
	return ans
}

func solveCase(n, qx, qy, bx, by, cx, cy int) string {
	if region(bx, by, qx, qy) == region(cx, cy, qx, qy) {
		return "YES"
	}
	return "NO"
}

func genCase(rng *rand.Rand) []byte {
	n := rng.Intn(998) + 3 // 3..1000
	qx := rng.Intn(n) + 1
	qy := rng.Intn(n) + 1
	var bx, by, cx, cy int
	for {
		bx = rng.Intn(n) + 1
		by = rng.Intn(n) + 1
		if bx != qx && by != qy && abs(bx-qx) != abs(by-qy) {
			break
		}
	}
	for {
		cx = rng.Intn(n) + 1
		cy = rng.Intn(n) + 1
		if cx != qx && cy != qy && abs(cx-qx) != abs(cy-qy) && (cx != bx || cy != by) {
			break
		}
	}
	return []byte(fmt.Sprintf("%d\n%d %d\n%d %d %d %d\n", n, qx, qy, bx, by, cx, cy))
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	cand, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()
	ref := "./refA.bin"
	if err := exec.Command("go", "build", "-o", ref, "1033A.go").Run(); err != nil {
		fmt.Println("failed to build reference solution:", err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := genCase(rng)
		want, err := run(ref, input)
		if err != nil {
			fmt.Printf("reference failed on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := run(cand, input)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\ninput:\n%s", i+1, err, string(input))
			os.Exit(1)
		}
		if strings.TrimSpace(want) != strings.TrimSpace(got) {
			fmt.Printf("wrong answer on test %d\ninput:\n%sexpected:%s\ngot:%s\n", i+1, string(input), want, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}

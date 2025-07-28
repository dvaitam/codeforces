package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func prepareBinary(path, tag string) (string, error) {
	if strings.HasSuffix(path, ".go") {
		bin := filepath.Join(os.TempDir(), fmt.Sprintf("%s_%d", tag, time.Now().UnixNano()))
		cmd := exec.Command("go", "build", "-o", bin, path)
		if out, err := cmd.CombinedOutput(); err != nil {
			return "", fmt.Errorf("build %s: %v\n%s", path, err, out)
		}
		return bin, nil
	}
	return path, nil
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

type testCaseC struct {
	n int
	s string
}

func genCase(rng *rand.Rand) testCaseC {
	n := rng.Intn(10) + 1
	bytes := make([]byte, n)
	for i := 0; i < n; i++ {
		bytes[i] = byte('a' + rng.Intn(3))
	}
	return testCaseC{n: n, s: string(bytes)}
}

const INF = int(1e9)

func cost(s string, ch byte) int {
	l := 0
	r := len(s) - 1
	cnt := 0
	for l < r {
		if s[l] == s[r] {
			l++
			r--
		} else if s[l] == ch {
			cnt++
			l++
		} else if s[r] == ch {
			cnt++
			r--
		} else {
			return INF
		}
	}
	return cnt
}

func solveCase(tc testCaseC) string {
	best := INF
	for c := byte('a'); c <= byte('z'); c++ {
		v := cost(tc.s, c)
		if v < best {
			best = v
		}
	}
	if best == INF {
		return "-1\n"
	}
	return fmt.Sprintf("%d\n", best)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	candPath, err := prepareBinary(os.Args[1], "candC")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(3))
	for i := 0; i < 100; i++ {
		tc := genCase(rng)
		input := fmt.Sprintf("1\n%d\n%s\n", tc.n, tc.s)
		expected := solveCase(tc)
		got, err := runBinary(candPath, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on case %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(expected) != strings.TrimSpace(got) {
			fmt.Printf("case %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s", i+1, input, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

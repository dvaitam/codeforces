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

type testCaseG struct {
	n   int
	arr []int
	ops string
}

func genCase(rng *rand.Rand) testCaseG {
	n := rng.Intn(6) + 2
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = rng.Intn(20) + 1
	}
	opsBytes := make([]byte, n)
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			opsBytes[i] = '*'
		} else {
			opsBytes[i] = '/'
		}
	}
	return testCaseG{n: n, arr: arr, ops: string(opsBytes)}
}

func factorize(x int) map[int]int {
	res := make(map[int]int)
	d := 2
	for d*d <= x {
		for x%d == 0 {
			res[d]++
			x /= d
		}
		d++
	}
	if x > 1 {
		res[x]++
	}
	return res
}

func solveG(tc testCaseG) string {
	n := tc.n
	factors := make([]map[int]int, n)
	for i := 0; i < n; i++ {
		factors[i] = factorize(tc.arr[i])
	}
	ans := 0
	for l := 0; l < n; l++ {
		counts := make(map[int]int)
		valid := true
		for r := l; r < n; r++ {
			for p, c := range factors[r] {
				if tc.ops[r] == '*' {
					counts[p] += c
				} else {
					counts[p] -= c
					if counts[p] < 0 {
						valid = false
					}
				}
			}
			if !valid {
				break
			}
			ans++
		}
	}
	return fmt.Sprintf("%d\n", ans)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	candPath, err := prepareBinary(os.Args[1], "candG")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(8))
	for i := 0; i < 100; i++ {
		tc := genCase(rng)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", tc.n))
		for j := 0; j < tc.n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", tc.arr[j]))
		}
		sb.WriteByte('\n')
		sb.WriteString(tc.ops)
		sb.WriteByte('\n')
		input := sb.String()
		expected := solveG(tc)
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

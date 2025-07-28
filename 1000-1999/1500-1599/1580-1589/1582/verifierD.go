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

type testCaseD struct {
	n   int
	arr []int64
}

func genCase(rng *rand.Rand) testCaseD {
	n := rng.Intn(9) + 2 // 2..10
	arr := make([]int64, n)
	for i := 0; i < n; i++ {
		v := int64(rng.Intn(20001) - 10000)
		if v == 0 {
			v = 1
		}
		arr[i] = v
	}
	return testCaseD{n: n, arr: arr}
}

func solveCase(tc testCaseD) string {
	n := tc.n
	a := tc.arr
	b := make([]int64, n)
	if n%2 == 0 {
		for i := 0; i < n; i += 2 {
			b[i] = -a[i+1]
			b[i+1] = a[i]
		}
	} else {
		for i := 0; i < n-3; i += 2 {
			b[i] = -a[i+1]
			b[i+1] = a[i]
		}
		i0 := n - 3
		i1 := n - 2
		i2 := n - 1
		if a[i0]+a[i1] != 0 {
			b[i0] = -a[i2]
			b[i1] = -a[i2]
			b[i2] = a[i0] + a[i1]
		} else if a[i0]+a[i2] != 0 {
			b[i0] = -a[i1]
			b[i1] = a[i0] + a[i2]
			b[i2] = -a[i1]
		} else {
			b[i0] = a[i1] + a[i2]
			b[i1] = -a[i0]
			b[i2] = -a[i0]
		}
	}
	var sb strings.Builder
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", b[i]))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	candPath, err := prepareBinary(os.Args[1], "candD")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(4))
	for i := 0; i < 100; i++ {
		tc := genCase(rng)
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d\n", tc.n))
		for j := 0; j < tc.n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", tc.arr[j]))
		}
		sb.WriteByte('\n')
		input := sb.String()
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

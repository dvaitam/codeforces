package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type TestA struct {
	n, k, x int
	arr     []int
}

func (t TestA) Input() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", t.n, t.k, t.x))
	for i, v := range t.arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func expectedA(t TestA) string {
	sum := 0
	for i := 0; i < t.n-t.k; i++ {
		sum += t.arr[i]
	}
	sum += t.k * t.x
	return strconv.Itoa(sum)
}

func run(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		return errBuf.String(), fmt.Errorf("runtime error: %v", err)
	}
	return out.String(), nil
}

func genTests() []TestA {
	rand.Seed(1)
	tests := make([]TestA, 0, 100)
	for i := 0; i < 100; i++ {
		n := rand.Intn(100) + 1
		k := rand.Intn(n) + 1
		arr := make([]int, n)
		arr[0] = rand.Intn(99) + 2
		for j := 1; j < n; j++ {
			val := arr[j-1] + rand.Intn(3)
			if val > 100 {
				val = 100
			}
			arr[j] = val
		}
		x := rand.Intn(arr[0]-1) + 1
		tests = append(tests, TestA{n: n, k: k, x: x, arr: arr})
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierA.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	tests := genTests()
	for i, tc := range tests {
		exp := strings.TrimSpace(expectedA(tc))
		gotRaw, err := run(bin, tc.Input())
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\n%s", i+1, err, gotRaw)
			os.Exit(1)
		}
		got := strings.TrimSpace(gotRaw)
		if got != exp {
			fmt.Printf("test %d failed\ninput:\n%sexpected: %s\ngot: %s\n", i+1, tc.Input(), exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

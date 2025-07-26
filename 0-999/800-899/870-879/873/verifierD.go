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

type TestD struct {
	n, k int
}

func (t TestD) Input() string {
	return fmt.Sprintf("%d %d\n", t.n, t.k)
}

func solveRec(l, r, k, L, R int, ans *[]int) {
	if k == 1 {
		for i := L; i <= R; i++ {
			*ans = append(*ans, i)
		}
		return
	}
	mid := (l + r - 1) >> 1
	Mid := (L + R + 2) >> 1
	leftCnt := mid - l + 1
	if 2*leftCnt-1 >= k-2 {
		solveRec(l, mid, k-2, Mid, R, ans)
		solveRec(mid+1, r, 1, L, Mid-1, ans)
	} else {
		solveRec(l, mid, 2*leftCnt-1, Mid, R, ans)
		solveRec(mid+1, r, k-2*leftCnt, L, Mid-1, ans)
	}
}

func expectedD(t TestD) string {
	n, k := t.n, t.k
	if k%2 == 0 || n*2 <= k {
		return "-1"
	}
	ans := make([]int, 0, n)
	solveRec(1, n, k, 1, n, &ans)
	var sb strings.Builder
	for i, v := range ans {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	return sb.String()
}

func run(bin, input string) (string, error) {
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

func genTests() []TestD {
	rand.Seed(4)
	tests := make([]TestD, 0, 100)
	for i := 0; i < 100; i++ {
		n := rand.Intn(20) + 1
		k := rand.Intn(2*n - 1)
		if k%2 == 0 {
			k++
			if k >= 2*n {
				k -= 2
			}
		}
		if k <= 0 {
			k = 1
		}
		tests = append(tests, TestD{n: n, k: k})
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	tests := genTests()
	for i, tc := range tests {
		exp := strings.TrimSpace(expectedD(tc))
		gotRaw, err := run(bin, tc.Input())
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\n%s", i+1, err, gotRaw)
			os.Exit(1)
		}
		got := strings.TrimSpace(gotRaw)
		if got != exp {
			fmt.Printf("test %d failed\ninput:%sexpected: %s\ngot: %s\n", i+1, tc.Input(), exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

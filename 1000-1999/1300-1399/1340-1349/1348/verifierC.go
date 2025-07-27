package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type Test struct {
	n int
	k int
	s string
}

func (t Test) Input() string {
	return fmt.Sprintf("1\n%d %d\n%s\n", t.n, t.k, t.s)
}

func expected(t Test) string {
	bytesArr := []byte(t.s)
	sort.Slice(bytesArr, func(i, j int) bool { return bytesArr[i] < bytesArr[j] })
	if bytesArr[0] != bytesArr[t.k-1] {
		return string(bytesArr[t.k-1])
	}
	res := []byte{bytesArr[0]}
	if t.k < t.n {
		if bytesArr[t.k] != bytesArr[t.n-1] {
			res = append(res, bytesArr[t.k:]...)
		} else {
			cnt := (t.n - t.k + t.k - 1) / t.k
			res = append(res, []byte(strings.Repeat(string(bytesArr[t.k]), cnt))...)
		}
	}
	return string(res)
}

func runProg(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genTest(rng *rand.Rand) Test {
	n := rng.Intn(20) + 1
	k := rng.Intn(n) + 1
	letters := make([]byte, n)
	for i := range letters {
		letters[i] = byte('a' + rng.Intn(26))
	}
	return Test{n: n, k: k, s: string(letters)}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	const cases = 100
	for i := 0; i < cases; i++ {
		tc := genTest(rng)
		expect := expected(tc)
		got, err := runProg(bin, tc.Input())
		if err != nil {
			fmt.Printf("case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Printf("case %d failed\ninput:\n%sexpected: %s\ngot: %s\n", i+1, tc.Input(), expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", cases)
}

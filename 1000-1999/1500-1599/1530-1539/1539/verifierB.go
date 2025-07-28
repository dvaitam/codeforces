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

type Test struct {
	n int
	q int
	s string
	L []int
	R []int
}

func (tc Test) Input() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.q))
	sb.WriteString(tc.s)
	sb.WriteByte('\n')
	for i := 0; i < tc.q; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.L[i], tc.R[i]))
	}
	return sb.String()
}

func expected(tc Test) string {
	pref := make([]int, tc.n+1)
	for i := 1; i <= tc.n; i++ {
		pref[i] = pref[i-1] + int(tc.s[i-1]-'a'+1)
	}
	var out strings.Builder
	for i := 0; i < tc.q; i++ {
		ans := pref[tc.R[i]] - pref[tc.L[i]-1]
		out.WriteString(fmt.Sprintf("%d\n", ans))
	}
	return strings.TrimSpace(out.String())
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
	n := rng.Intn(30) + 1
	q := rng.Intn(30) + 1
	bytes := make([]byte, n)
	for i := 0; i < n; i++ {
		bytes[i] = byte('a' + rng.Intn(26))
	}
	s := string(bytes)
	L := make([]int, q)
	R := make([]int, q)
	for i := 0; i < q; i++ {
		l := rng.Intn(n) + 1
		r := rng.Intn(n-l+1) + l
		L[i] = l
		R[i] = r
	}
	return Test{n, q, s, L, R}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
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
		if strings.TrimSpace(got) != expect {
			fmt.Printf("case %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i+1, tc.Input(), expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", cases)
}

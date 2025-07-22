package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func comb(n, k int) int {
	if k < 0 || k > n {
		return 0
	}
	res := 1
	for i := 0; i < k; i++ {
		res = res * (n - i) / (i + 1)
	}
	return res
}

func expected(s1, s2 string) float64 {
	target := 0
	for _, c := range s1 {
		if c == '+' {
			target++
		} else {
			target--
		}
	}
	curr := 0
	q := 0
	for _, c := range s2 {
		switch c {
		case '+':
			curr++
		case '-':
			curr--
		case '?':
			q++
		}
	}
	diff := target - curr
	if (diff+q)%2 != 0 || int(math.Abs(float64(diff))) > q {
		return 0
	}
	up := (q + diff) / 2
	ways := comb(q, up)
	return float64(ways) / math.Pow(2, float64(q))
}

type testCase struct {
	s1 string
	s2 string
}

func randStr(rng *rand.Rand, n int, alphabet string) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = alphabet[rng.Intn(len(alphabet))]
	}
	return string(b)
}

func generateRandomCase(rng *rand.Rand) testCase {
	l := rng.Intn(10) + 1
	s1 := randStr(rng, l, "+-")
	s2 := randStr(rng, l, "+-?")
	return testCase{s1, s2}
}

func runCase(bin string, tc testCase) error {
	input := fmt.Sprintf("%s\n%s\n", tc.s1, tc.s2)
	exp := expected(tc.s1, tc.s2)
	out, err := runCandidate(bin, input)
	if err != nil {
		return err
	}
	var got float64
	if _, err := fmt.Sscan(out, &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if math.Abs(got-exp) > 1e-7 {
		return fmt.Errorf("expected %.9f got %.9f", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []testCase{
		{"++", "+?"},
		{"+-", "??"},
	}
	for i := 0; i < 100; i++ {
		cases = append(cases, generateRandomCase(rng))
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Test struct {
	n   int
	k   int
	arr []int
}

func (t Test) Input() string {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", t.n, t.k))
	for i, v := range t.arr {
		if i+1 == len(t.arr) {
			sb.WriteString(fmt.Sprintf("%d\n", v))
		} else {
			sb.WriteString(fmt.Sprintf("%d ", v))
		}
	}
	return sb.String()
}

func expected(t Test) string {
	distinct := make(map[int]bool)
	beauty := []int{}
	for _, v := range t.arr {
		if !distinct[v] {
			distinct[v] = true
			beauty = append(beauty, v)
		}
	}
	if len(beauty) > t.k {
		return "-1"
	}
	for i := 1; i <= t.n && len(beauty) < t.k; i++ {
		if !distinct[i] {
			distinct[i] = true
			beauty = append(beauty, i)
		}
	}
	total := t.n * len(beauty)
	seq := make([]string, 0, total)
	for i := 0; i < t.n; i++ {
		for _, v := range beauty {
			seq = append(seq, strconv.Itoa(v))
		}
	}
	return fmt.Sprintf("%d\n%s", total, strings.Join(seq, " "))
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
	n := rng.Intn(100) + 1
	k := rng.Intn(n) + 1
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = rng.Intn(n) + 1
	}
	return Test{n: n, k: k, arr: arr}
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
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Printf("case %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i+1, tc.Input(), expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", cases)
}

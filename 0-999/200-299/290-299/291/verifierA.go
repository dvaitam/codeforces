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

func solveA(ids []int) string {
	freq := map[int]int{}
	for _, id := range ids {
		if id != 0 {
			freq[id]++
		}
	}
	pairs := 0
	for _, c := range freq {
		if c > 2 {
			return "-1\n"
		}
		if c == 2 {
			pairs++
		}
	}
	return fmt.Sprintf("%d\n", pairs)
}

func genCaseA(rng *rand.Rand) (string, string) {
	n := rng.Intn(15) + 1
	ids := make([]int, n)
	for i := range ids {
		if rng.Intn(4) == 0 {
			ids[i] = 0
		} else {
			ids[i] = rng.Intn(6)
		}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i, v := range ids {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	return sb.String(), solveA(ids)
}

func run(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func main() {
	rand.Seed(time.Now().UnixNano())
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		in, expect := genCaseA(rand.New(rand.NewSource(time.Now().UnixNano() + int64(i))))
		got, err := run(bin, in)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\ninput:\n%soutput:\n%s", i+1, err, in, got)
			return
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Printf("test %d failed\ninput:\n%sexpected:\n%sbut got:\n%s", i+1, in, expect, got)
			return
		}
	}
	fmt.Println("All tests passed")
}

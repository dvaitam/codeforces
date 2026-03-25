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

// ---------- embedded solver (brute-force, correct for small m) ----------

func solveG(m, x int64) int64 {
	if m == 1 {
		return 1 // only 0
	}
	seen := make(map[int64]bool)
	cur := int64(1) % m
	for {
		if seen[cur] {
			break
		}
		seen[cur] = true
		cur = (cur * x) % m
	}
	return int64(len(seen))
}

func solveInput(input []byte) string {
	scanner := strings.NewReader(string(input))
	var sb strings.Builder
	for {
		var m, x int64
		_, err := fmt.Fscan(scanner, &m, &x)
		if err != nil {
			break
		}
		fmt.Fprintln(&sb, solveG(m, x))
	}
	return sb.String()
}

// ---------- verifier logic ----------

func runBin(bin string, input []byte) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func genTest(rng *rand.Rand) []byte {
	t := rng.Intn(5) + 1
	var sb strings.Builder
	for i := 0; i < t; i++ {
		m := rng.Intn(20) + 1
		x := rng.Intn(20) + 1
		sb.WriteString(fmt.Sprintf("%d %d\n", m, x))
	}
	return []byte(sb.String())
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	cand := os.Args[1]

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= 100; i++ {
		input := genTest(rng)
		want := solveInput(input)
		got, err := runBin(cand, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\ninput:\n%s", i, err, string(input))
			os.Exit(1)
		}
		if strings.TrimSpace(want) != strings.TrimSpace(got) {
			fmt.Printf("wrong answer on test %d\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i, string(input), want, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

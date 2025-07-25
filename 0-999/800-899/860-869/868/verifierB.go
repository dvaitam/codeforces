package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type Test struct {
	in  string
	out string
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solveOracle(input string) string {
	reader := bufio.NewReader(strings.NewReader(input))
	var h, m, s, t1, t2 int
	if _, err := fmt.Fscan(reader, &h, &m, &s, &t1, &t2); err != nil {
		return ""
	}
	hPos := float64(h%12) + float64(m)/60.0 + float64(s)/3600.0
	mPos := float64(m)/5.0 + float64(s)/300.0
	sPos := float64(s) / 5.0
	t1Pos := float64(t1 % 12)
	t2Pos := float64(t2 % 12)
	type item struct {
		pos float64
		id  int
	}
	arr := []item{{hPos, 0}, {mPos, 1}, {sPos, 2}, {t1Pos, 3}, {t2Pos, 4}}
	sort.Slice(arr, func(i, j int) bool { return arr[i].pos < arr[j].pos })
	for i := 0; i < 5; i++ {
		j := (i + 1) % 5
		if arr[i].id >= 3 && arr[j].id >= 3 {
			return "YES"
		}
	}
	return "NO"
}

func genCase(rng *rand.Rand) Test {
	h := rng.Intn(12) + 1
	m := rng.Intn(59) + 1 // avoid zero for simplicity
	s := rng.Intn(59) + 1
	t1 := rng.Intn(12) + 1
	t2 := rng.Intn(12) + 1
	for t2 == t1 {
		t2 = rng.Intn(12) + 1
	}
	input := fmt.Sprintf("%d %d %d %d %d\n", h, m, s, t1, t2)
	out := solveOracle(input)
	return Test{input, out}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(2))
	for i := 0; i < 100; i++ {
		tc := genCase(rng)
		got, err := run(bin, tc.in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.in)
			os.Exit(1)
		}
		if got != tc.out {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, tc.out, got, tc.in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

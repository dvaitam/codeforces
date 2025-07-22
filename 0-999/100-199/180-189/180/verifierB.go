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

func primeFactors(n int) map[int]int {
	m := make(map[int]int)
	d := 2
	for d*d <= n {
		for n%d == 0 {
			m[d]++
			n /= d
		}
		d++
	}
	if n > 1 {
		m[n]++
	}
	return m
}

func expected(b, d int) (string, int) {
	fb := primeFactors(b)
	fd := primeFactors(d)
	ok2 := true
	k := 0
	for p, ed := range fd {
		eb := fb[p]
		if eb == 0 {
			ok2 = false
			break
		}
		need := (ed + eb - 1) / eb
		if need > k {
			k = need
		}
	}
	if ok2 {
		return "2-type", k
	}
	if (b-1)%d == 0 {
		return "3-type", 0
	}
	if (b+1)%d == 0 {
		return "11-type", 0
	}
	ok6 := true
	for p := range fd {
		if fb[p] > 0 {
			continue
		}
		if (b-1)%p == 0 || (b+1)%p == 0 {
			continue
		}
		ok6 = false
		break
	}
	if ok6 {
		return "6-type", 0
	}
	return "7-type", 0
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
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

func runCase(bin string, b, d int) error {
	input := fmt.Sprintf("%d %d\n", b, d)
	expType, expK := expected(b, d)
	out, err := run(bin, input)
	if err != nil {
		return err
	}
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return fmt.Errorf("empty output")
	}
	if fields[0] != expType {
		return fmt.Errorf("expected %s got %s", expType, fields[0])
	}
	if expType == "2-type" {
		if len(fields) != 2 {
			return fmt.Errorf("expected 2 lines for 2-type")
		}
		kGot, err := strconv.Atoi(fields[1])
		if err != nil {
			return fmt.Errorf("invalid k: %v", err)
		}
		if kGot != expK {
			return fmt.Errorf("expected %d got %d", expK, kGot)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 && !(len(os.Args) == 3 && os.Args[1] == "--") {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[len(os.Args)-1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	cases := [][2]int{
		{2, 2},
		{10, 3},
		{10, 11},
		{10, 6},
		{2, 7},
	}

	for i := 0; i < 100; i++ {
		b := rng.Intn(99) + 2
		d := rng.Intn(99) + 2
		cases = append(cases, [2]int{b, d})
	}

	for idx, c := range cases {
		if err := runCase(bin, c[0], c[1]); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%d %d\n", idx+1, err, c[0], c[1])
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

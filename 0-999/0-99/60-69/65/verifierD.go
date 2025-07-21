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

func run(bin, input string) (string, error) {
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

func expected(n int, s string) string {
	c := [4]int{}
	k := 0
	for i := 0; i < n && i < len(s); i++ {
		switch s[i] {
		case 'G':
			c[0]++
		case 'H':
			c[1]++
		case 'R':
			c[2]++
		case 'S':
			c[3]++
		case '?':
			k++
		}
	}
	type pair struct{ val, idx int }
	p := []pair{{c[0], 0}, {c[1], 1}, {c[2], 2}, {c[3], 3}}
	for i := 0; i < 4; i++ {
		for j := i + 1; j < 4; j++ {
			if p[j].val < p[i].val {
				p[i], p[j] = p[j], p[i]
			}
		}
	}
	kLeft := k
	level := p[0].val
	idxGroup := 1
	for idxGroup < 4 {
		need := idxGroup * (p[idxGroup].val - level)
		if kLeft >= need {
			kLeft -= need
			level = p[idxGroup].val
			idxGroup++
		} else {
			break
		}
	}
	if idxGroup > 0 {
		f := kLeft / idxGroup
		level += f
		kLeft -= f * idxGroup
	}
	names := []string{"Gryffindor", "Hufflepuff", "Ravenclaw", "Slytherin"}
	res := make([]string, 0, idxGroup)
	for i := 0; i < idxGroup; i++ {
		res = append(res, names[p[i].idx])
	}
	for i := 0; i < len(res); i++ {
		for j := i + 1; j < len(res); j++ {
			if res[j] < res[i] {
				res[i], res[j] = res[j], res[i]
			}
		}
	}
	return strings.Join(res, "\n")
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(20) + 1
	letters := []byte{'G', 'H', 'R', 'S', '?'}
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = letters[rng.Intn(len(letters))]
	}
	s := string(b)
	input := fmt.Sprintf("%d\n%s\n", n, s)
	exp := expected(n, s)
	return input, exp
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected\n%s\ngot\n%s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

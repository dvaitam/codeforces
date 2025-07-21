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

func candidates(y int) []int {
	ystr := fmt.Sprintf("%04d", y)
	cands := make([]int, 0, 50)
	if y >= 1000 && y <= 2011 {
		cands = append(cands, y)
	}
	for j := 0; j < 4; j++ {
		for d := byte('0'); d <= '9'; d++ {
			if d == ystr[j] {
				continue
			}
			if j == 0 && d == '0' {
				continue
			}
			zbytes := []byte(ystr)
			zbytes[j] = d
			z := int(zbytes[0]-'0')*1000 + int(zbytes[1]-'0')*100 + int(zbytes[2]-'0')*10 + int(zbytes[3]-'0')
			if z < 1000 || z > 2011 {
				continue
			}
			cands = append(cands, z)
		}
	}
	sort.Ints(cands)
	uniq := make([]int, 0, len(cands))
	prev := -1
	for _, z := range cands {
		if z != prev {
			uniq = append(uniq, z)
			prev = z
		}
	}
	return uniq
}

func expected(n int, years []int) (string, bool) {
	prev := 1000
	res := make([]int, n)
	for i, y := range years {
		cands := candidates(y)
		pick := -1
		for _, z := range cands {
			if z >= prev {
				pick = z
				break
			}
		}
		if pick < 0 {
			return "No solution", false
		}
		res[i] = pick
		prev = pick
	}
	var sb strings.Builder
	for i, v := range res {
		if i > 0 {
			sb.WriteByte('\n')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	return sb.String(), true
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(8) + 1
	years := make([]int, n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		years[i] = rng.Intn(9000) + 1000
		sb.WriteString(fmt.Sprintf("%d\n", years[i]))
	}
	input := sb.String()
	exp, ok := expected(n, years)
	if !ok {
		exp = "No solution"
	}
	return input, exp
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
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
		if out != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

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

func runBinary(bin, input string) (string, error) {
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

func solveB(sum, lim int) string {
	const maxb = 31
	s := make([]int, maxb)
	for i := 1; i <= lim; i++ {
		for j := 0; j < maxb; j++ {
			if i&(1<<j) != 0 {
				s[j]++
				break
			}
		}
	}
	num := make([]int, maxb)
	fail := false
	for i := 0; i < maxb; i++ {
		if (sum>>i)&1 == 0 {
			continue
		}
		need := 1
		for j := i; j >= 0; j-- {
			if s[j] >= need {
				num[j] += need
				s[j] -= need
				break
			}
			need -= s[j]
			num[j] += s[j]
			need <<= 1
			s[j] = 0
			if j == 0 {
				fail = true
				break
			}
		}
		if fail {
			break
		}
	}
	if fail {
		return "-1"
	}
	total := 0
	for _, v := range num {
		total += v
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", total)
	for i := 1; i <= lim; i++ {
		for j := 0; j < maxb; j++ {
			if i&(1<<j) != 0 && num[j] > 0 {
				fmt.Fprintf(&sb, "%d ", i)
				num[j]--
				break
			}
		}
	}
	sb.WriteByte('\n')
	return strings.TrimSpace(sb.String())
}

func generateCase(rng *rand.Rand) (string, string) {
	sum := rng.Intn(1000) + 1
	lim := rng.Intn(1000) + 1
	input := fmt.Sprintf("%d %d\n", sum, lim)
	exp := solveB(sum, lim)
	return input, exp
}

func manualCase(sum, lim int) (string, string) {
	input := fmt.Sprintf("%d %d\n", sum, lim)
	exp := solveB(sum, lim)
	return input, exp
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var cases [][2]string
	in1, ex1 := manualCase(5, 5)
	cases = append(cases, [2]string{in1, ex1})
	in2, ex2 := manualCase(1, 1)
	cases = append(cases, [2]string{in2, ex2})
	in3, ex3 := manualCase(100, 50)
	cases = append(cases, [2]string{in3, ex3})
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		cases = append(cases, [2]string{in, exp})
	}
	for idx, tc := range cases {
		out, err := runBinary(bin, tc[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", idx+1, err, tc[0])
			os.Exit(1)
		}
		if strings.TrimSpace(out) != tc[1] {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", idx+1, tc[1], out, tc[0])
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

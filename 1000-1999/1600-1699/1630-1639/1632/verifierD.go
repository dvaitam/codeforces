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

type pair struct {
	g int
	l int
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func solveD(arr []int) []int {
	const BIG = 1000000007
	segs := []pair{}
	ans := 0
	res := make([]int, len(arr))
	for i, x := range arr {
		newSegs := []pair{{x, i}}
		for _, p := range segs {
			g := gcd(p.g, x)
			last := &newSegs[len(newSegs)-1]
			if last.g == g {
				if p.l < last.l {
					last.l = p.l
				}
			} else {
				newSegs = append(newSegs, pair{g, p.l})
			}
		}
		bad := false
		for _, p := range newSegs {
			if p.g == i-p.l+1 {
				bad = true
				break
			}
		}
		if bad {
			ans++
			segs = []pair{{BIG, i}}
		} else {
			segs = newSegs
		}
		res[i] = ans
	}
	return res
}

func generateCase(rng *rand.Rand) (string, []int) {
	n := rng.Intn(20) + 1
	arr := make([]int, n)
	for i := range arr {
		arr[i] = rng.Intn(100) + 1
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteString("\n")
	return sb.String(), solveD(arr)
}

func runCase(bin, input string, exp []int) error {
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
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	parts := strings.Fields(strings.TrimSpace(out.String()))
	if len(parts) != len(exp) {
		return fmt.Errorf("expected %d numbers got %d", len(exp), len(parts))
	}
	for i, p := range parts {
		var v int
		if _, err := fmt.Sscan(p, &v); err != nil {
			return fmt.Errorf("bad int at pos %d: %v", i+1, err)
		}
		if v != exp[i] {
			return fmt.Errorf("pos %d expected %d got %d", i+1, exp[i], v)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

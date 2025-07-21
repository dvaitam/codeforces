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

func isGP(a []int, skip int) bool {
	n := len(a)
	k := 0
	var pp, p int
	for i := 0; i < n; i++ {
		if i == skip {
			continue
		}
		x := a[i]
		if k == 0 {
			pp = x
		} else if k == 1 {
			p = x
			if pp == 0 && p != 0 {
				return false
			}
		} else {
			if pp == 0 {
				if x != 0 {
					return false
				}
			} else {
				if int64(p)*int64(p) != int64(pp)*int64(x) {
					return false
				}
			}
			pp = p
			p = x
		}
		k++
	}
	return true
}

func solveD(a []int) string {
	if isGP(a, -1) {
		return "0"
	}
	n := len(a)
	cand := map[int]bool{}
	have := false
	if n >= 2 && a[0] == 0 && a[1] != 0 {
		cand[0] = true
		cand[1] = true
		have = true
	}
	for i := 2; i < n; i++ {
		x0, x1, x2 := a[i-2], a[i-1], a[i]
		fail := false
		if x0 == 0 {
			if x1 != 0 {
				fail = true
			}
		} else {
			if int64(x1)*int64(x1) != int64(x0)*int64(x2) {
				fail = true
			}
		}
		if fail {
			tri := []int{i - 2, i - 1, i}
			if !have {
				for _, j := range tri {
					if j >= 0 && j < n {
						cand[j] = true
					}
				}
				have = true
			} else {
				newC := map[int]bool{}
				for _, j := range tri {
					if cand[j] {
						newC[j] = true
					}
				}
				cand = newC
			}
			if len(cand) == 0 {
				return "2"
			}
		}
	}
	for j := range cand {
		if isGP(a, j) {
			return "1"
		}
	}
	return "2"
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 1
	arr := make([]int, n)
	for i := range arr {
		arr[i] = rng.Intn(10) - 5
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	expected := solveD(arr)
	return sb.String(), expected
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

func main() {
	if len(os.Args) != 2 && !(len(os.Args) == 3 && os.Args[1] == "--") {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[len(os.Args)-1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, expect := generateCase(rng)
		got, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

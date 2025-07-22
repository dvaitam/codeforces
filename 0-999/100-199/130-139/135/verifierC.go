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

func runCandidate(bin, input string) (string, error) {
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

func solveCase(s string) []string {
	n := len(s)
	pos0 := []int{}
	pos1 := []int{}
	qpos := []int{}
	for i, c := range s {
		switch c {
		case '0':
			pos0 = append(pos0, i)
		case '1':
			pos1 = append(pos1, i)
		case '?':
			qpos = append(qpos, i)
		}
	}
	C0 := len(pos0)
	C1 := len(pos1)
	Q := len(qpos)
	D := n - 2
	if D < 0 {
		D = 0
	}
	m := D / 2
	M := (D + 1) / 2
	can00, can01, can10, can11 := false, false, false, false
	if C0+Q >= m+2 {
		can00 = true
	}
	if C1+Q >= M+2 {
		can11 = true
	}
	if C0+Q >= m+1 && C1+Q >= M+1 {
		L := m + 1 - C0
		if L < 0 {
			L = 0
		}
		R := C1 + Q - (M + 1)
		if R > Q {
			R = Q
		}
		if R >= L {
			var Zmin int
			if m < C0 {
				Zmin = pos0[m]
			} else {
				Zmin = qpos[m-C0]
			}
			var Omax int
			if M < C1 {
				Omax = pos1[M]
			} else {
				idx := R + (M - C1)
				if idx < 0 {
					idx = 0
				}
				if idx >= len(qpos) {
					idx = len(qpos) - 1
				}
				Omax = qpos[idx]
			}
			if Zmin < Omax {
				can01 = true
			}
			var Zmax int
			if m < C0 {
				Zmax = pos0[m]
			} else {
				Zmax = qpos[Q-1]
			}
			var Omin int
			if M < C1 {
				Omin = pos1[M]
			} else {
				idx := M - C1
				if idx < 0 {
					idx = 0
				}
				if idx >= len(qpos) {
					idx = len(qpos) - 1
				}
				Omin = qpos[idx]
			}
			if Zmax > Omin {
				can10 = true
			}
		}
	}
	res := []string{}
	if can00 {
		res = append(res, "00")
	}
	if can01 {
		res = append(res, "01")
	}
	if can10 {
		res = append(res, "10")
	}
	if can11 {
		res = append(res, "11")
	}
	return res
}

func generateCase(rng *rand.Rand) (string, []string) {
	n := rng.Intn(20) + 2
	bytes := make([]byte, n)
	chars := []byte{'0', '1', '?'}
	for i := range bytes {
		bytes[i] = chars[rng.Intn(3)]
	}
	s := string(bytes)
	expect := solveCase(s)
	return s + "\n", expect
}

func compare(out string, expect []string) bool {
	fields := strings.Fields(out)
	if len(fields) != len(expect) {
		return false
	}
	for i := range fields {
		if fields[i] != expect[i] {
			return false
		}
	}
	return true
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if !compare(out, exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %v got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

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
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solveOne(n int, p int, a []int) int {
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = a[n-1-i]
	}
	arr0 := arr[0]
	union := map[int]bool{}
	earliest := map[int]int{}
	for _, d := range arr {
		union[d] = true
		earliest[d] = 0
	}
	prefix := p - arr0
	if arr0 != 0 && prefix <= p-1 {
		carry := 1
		i := 1
		for {
			if i < n {
				newd := (arr[i] + carry) % p
				t := prefix
				jlsd := (newd - arr0) % p
				if jlsd < 0 {
					jlsd += p
				}
				if jlsd < t {
					t = jlsd
				}
				if old, ok := earliest[newd]; !ok || t < old {
					earliest[newd] = t
				}
				union[newd] = true
				if arr[i]+carry >= p {
					carry = 1
					i++
					continue
				} else {
					break
				}
			} else {
				newd := carry
				t := prefix
				jlsd := (newd - arr0) % p
				if jlsd < 0 {
					jlsd += p
				}
				if jlsd < t {
					t = jlsd
				}
				if old, ok := earliest[newd]; !ok || t < old {
					earliest[newd] = t
				}
				union[newd] = true
				break
			}
		}
	}
	jSet := map[int]bool{}
	for d := range union {
		j := (d - arr0) % p
		if j < 0 {
			j += p
		}
		jSet[j] = true
	}
	j := p - 1
	for ; j >= 0; j-- {
		if !jSet[j] {
			break
		}
	}
	tMissing := 0
	if j >= 0 {
		tMissing = j
	}
	tUnion := 0
	for _, v := range earliest {
		if v > tUnion {
			tUnion = v
		}
	}
	if tUnion > tMissing {
		return tUnion
	}
	return tMissing
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 1
	p := rng.Intn(20) + 2
	arr := make([]int, n)
	for i := range arr {
		arr[i] = rng.Intn(p)
	}
	ans := solveOne(n, p, arr)
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", n, p))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteString("\n")
	return sb.String(), fmt.Sprint(ans)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
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
		if strings.TrimSpace(out) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

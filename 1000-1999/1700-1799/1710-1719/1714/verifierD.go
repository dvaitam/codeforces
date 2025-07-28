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

func run(bin, input string) (string, error) {
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

func minOps(t string, arr []string) int {
	const INF = int(1e9)
	L := len(t)
	ss := make([][2]int, L)
	for i := 0; i < L; i++ {
		ss[i][0], ss[i][1] = INF, INF
	}
	for i, pat := range arr {
		plen := len(pat)
		for j := 0; j+plen <= L; j++ {
			if t[j:j+plen] == pat {
				end := j + plen - 1
				if ss[end][1] > j {
					ss[end][0] = i
					ss[end][1] = j
				}
			}
		}
	}
	rpos := L - 1
	ans := 0
	for rpos >= 0 {
		mv := rpos
		for i := rpos; i < L; i++ {
			if ss[i][1] <= rpos && ss[i][1]-1 < mv {
				mv = ss[i][1] - 1
			}
		}
		if mv >= rpos {
			return -1
		}
		ans++
		rpos = mv
	}
	return ans
}

func verifyOutput(out string, t string, arr []string, exp int) error {
	fields := strings.Fields(out)
	if exp == -1 {
		if len(fields) != 1 || fields[0] != "-1" {
			return fmt.Errorf("expected -1, got %s", out)
		}
		return nil
	}
	if len(fields) != 1+2*exp {
		return fmt.Errorf("expected %d numbers, got %d", 1+2*exp, len(fields))
	}
	m, err := strconv.Atoi(fields[0])
	if err != nil || m != exp {
		return fmt.Errorf("expected first number %d", exp)
	}
	colored := make([]bool, len(t))
	for i := 0; i < m; i++ {
		idx, err1 := strconv.Atoi(fields[1+2*i])
		pos, err2 := strconv.Atoi(fields[2+2*i])
		if err1 != nil || err2 != nil {
			return fmt.Errorf("bad integers")
		}
		if idx < 1 || idx > len(arr) {
			return fmt.Errorf("bad index")
		}
		pos--
		pat := arr[idx-1]
		if pos < 0 || pos+len(pat) > len(t) {
			return fmt.Errorf("bad position")
		}
		if t[pos:pos+len(pat)] != pat {
			return fmt.Errorf("substring mismatch")
		}
		for j := 0; j < len(pat); j++ {
			colored[pos+j] = true
		}
	}
	for _, c := range colored {
		if !c {
			return fmt.Errorf("not fully covered")
		}
	}
	return nil
}

func genCase(rng *rand.Rand) (string, string, []string, int) {
	L := rng.Intn(8) + 1
	alphabet := []rune("abc")
	var sb strings.Builder
	for i := 0; i < L; i++ {
		sb.WriteRune(alphabet[rng.Intn(len(alphabet))])
	}
	t := sb.String()
	n := rng.Intn(3) + 1
	arr := make([]string, n)
	for i := 0; i < n; i++ {
		ln := rng.Intn(3) + 1
		var tb strings.Builder
		for j := 0; j < ln; j++ {
			tb.WriteRune(alphabet[rng.Intn(len(alphabet))])
		}
		arr[i] = tb.String()
	}
	var input strings.Builder
	input.WriteString("1\n")
	input.WriteString(fmt.Sprintf("%s\n", t))
	input.WriteString(fmt.Sprintf("%d\n", n))
	for _, s := range arr {
		input.WriteString(fmt.Sprintf("%s\n", s))
	}
	exp := minOps(t, arr)
	return input.String(), t, arr, exp
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, t, arr, exp := genCase(rng)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if err := verifyOutput(out, t, arr, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%soutput:\n%s\n", i+1, err, input, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

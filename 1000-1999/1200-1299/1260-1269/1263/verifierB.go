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

type resultB struct {
	cnt   int
	codes []string
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
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solveB(codes []string) resultB {
	mp := make(map[string]int)
	n := len(codes)
	for _, s := range codes {
		mp[s]++
	}
	res := make([]string, n)
	copy(res, codes)
	cnt := 0
	for i := 0; i < n; i++ {
		if mp[res[i]] > 1 {
			cnt++
			orig := res[i]
			tmp := []byte(orig)
			ok := false
			for j := byte('0'); j <= '9'; j++ {
				tmp[3] = j
				s := string(tmp)
				if mp[s] == 0 {
					ok = true
					res[i] = s
					break
				}
			}
			if !ok {
				tmp = []byte(orig)
				for j := byte('0'); j <= '9'; j++ {
					tmp[2] = j
					s := string(tmp)
					if mp[s] == 0 {
						ok = true
						res[i] = s
						break
					}
				}
			}
			mp[orig]--
			mp[res[i]]++
		}
	}
	return resultB{cnt, res}
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(9) + 2
	codes := make([]string, n)
	for i := 0; i < n; i++ {
		var b [4]byte
		for j := 0; j < 4; j++ {
			b[j] = byte('0' + rng.Intn(10))
		}
		codes[i] = string(b[:])
	}
	inputBuilder := strings.Builder{}
	inputBuilder.WriteString("1\n")
	inputBuilder.WriteString(fmt.Sprintf("%d\n", n))
	for _, s := range codes {
		inputBuilder.WriteString(s)
		inputBuilder.WriteByte('\n')
	}
	res := solveB(codes)
	outBuilder := strings.Builder{}
	outBuilder.WriteString(fmt.Sprintf("%d\n", res.cnt))
	for i, s := range res.codes {
		outBuilder.WriteString(s)
		if i+1 < len(res.codes) {
			outBuilder.WriteByte('\n')
		}
	}
	return inputBuilder.String(), strings.TrimSpace(outBuilder.String())
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected\n%s\ngot\n%s\n", i+1, exp, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

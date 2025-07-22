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

func expectedA(n, k int, s string) string {
	p := k - 1
	flag1, flag2 := 0, 0
	if p == 0 {
		flag1, flag2 = 1, 0
	} else if p == n-1 {
		flag1, flag2 = 0, 0
	} else if p < n-1-p {
		flag1, flag2 = 0, 1
	} else {
		flag1, flag2 = 1, 1
	}
	var cmds []string
	cmds = append(cmds, fmt.Sprintf("PRINT %c", s[p]))
	if flag1 == 1 {
		cnt := p + 1
		for cnt < n-1 {
			cmds = append(cmds, "RIGHT")
			cmds = append(cmds, fmt.Sprintf("PRINT %c", s[cnt]))
			cnt++
		}
		if cnt < n {
			cmds = append(cmds, "RIGHT")
			cmds = append(cmds, fmt.Sprintf("PRINT %c", s[cnt]))
		}
		cnt = n - 1
		if flag2 == 1 {
			for cnt >= p {
				cmds = append(cmds, "LEFT")
				cnt--
			}
			for cnt > 0 {
				cmds = append(cmds, fmt.Sprintf("PRINT %c", s[cnt]))
				cmds = append(cmds, "LEFT")
				cnt--
			}
			if cnt >= 0 {
				cmds = append(cmds, fmt.Sprintf("PRINT %c", s[cnt]))
			}
		}
	} else {
		cnt := p - 1
		for cnt > 0 {
			cmds = append(cmds, "LEFT")
			cmds = append(cmds, fmt.Sprintf("PRINT %c", s[cnt]))
			cnt--
		}
		if cnt >= 0 {
			cmds = append(cmds, "LEFT")
			cmds = append(cmds, fmt.Sprintf("PRINT %c", s[cnt]))
		}
		cnt = 0
		if flag2 == 1 {
			for cnt <= p {
				cmds = append(cmds, "RIGHT")
				cnt++
			}
			for cnt < n-1 {
				cmds = append(cmds, fmt.Sprintf("PRINT %c", s[cnt]))
				cmds = append(cmds, "RIGHT")
				cnt++
			}
			if cnt < n {
				cmds = append(cmds, fmt.Sprintf("PRINT %c", s[cnt]))
			}
		}
	}
	return strings.Join(cmds, "\n")
}

func runCase(bin string, n, k int, s string) error {
	input := fmt.Sprintf("%d %d\n%s\n", n, k, s)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	expect := strings.TrimSpace(expectedA(n, k, s))
	if got != expect {
		return fmt.Errorf("expected:\n%s\ngot:\n%s", expect, got)
	}
	return nil
}

func genCase(rng *rand.Rand) (int, int, string) {
	n := rng.Intn(100) + 1
	k := rng.Intn(n) + 1
	alphabet := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789.!?,"
	b := make([]byte, n)
	for i := range b {
		b[i] = alphabet[rng.Intn(len(alphabet))]
	}
	return n, k, string(b)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, k, s := genCase(rng)
		if err := runCase(bin, n, k, s); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%d %d\n%s\n", i+1, err, n, k, s)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

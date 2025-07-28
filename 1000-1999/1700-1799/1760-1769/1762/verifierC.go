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

const mod int64 = 998244353

func expectedC(s string) int64 {
	same, large := int64(1), int64(0)
	ans := int64(1)
	for i := 1; i < len(s); i++ {
		if s[i] == s[i-1] {
			newSame := same
			newLarge := (same + 2*large) % mod
			same = newSame % mod
			large = newLarge
		} else {
			same = same % mod
			large = 0
		}
		ans = (ans + same + large) % mod
	}
	return ans % mod
}

func genTestC(rng *rand.Rand) (string, string) {
	n := rng.Intn(20) + 1
	var sb strings.Builder
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			sb.WriteByte('0')
		} else {
			sb.WriteByte('1')
		}
	}
	s := sb.String()
	exp := expectedC(s)
	input := fmt.Sprintf("1\n%d\n%s\n", n, s)
	return input, fmt.Sprintf("%d", exp)
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
	err := cmd.Run()
	if err != nil {
		return out.String() + errBuf.String(), err
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		input, expect := genTestC(rng)
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n%s", t+1, err, got)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("test %d failed: expected %s got %s\n", t+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", 100)
}

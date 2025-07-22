package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return out.String(), fmt.Errorf("%v: %s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

const MOD int64 = 1000000007

func modpow(a, e int64) int64 {
	res := int64(1)
	a %= MOD
	for e > 0 {
		if e&1 == 1 {
			res = res * a % MOD
		}
		a = a * a % MOD
		e >>= 1
	}
	return res
}

func expected(start, end string, k int64) int64 {
	n := len(start)
	doubled := start + start
	cnt := 0
	cnt0 := 0
	for s := 0; s < n; s++ {
		if doubled[s:s+n] == end {
			cnt++
			if s == 0 {
				cnt0 = 1
			}
		}
	}
	if cnt == 0 {
		return 0
	}
	t := modpow(int64(n-1), k)
	p := int64(1)
	if k&1 == 1 {
		p = MOD - 1
	}
	invN := modpow(int64(n), MOD-2)
	C0 := (t + p*int64(n-1)) % MOD * invN % MOD
	C1 := (t - p + MOD) % MOD * invN % MOD
	ways := (int64(cnt0)*C0 + int64(cnt-cnt0)*C1) % MOD
	return ways
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	if bin == "--" && len(os.Args) >= 3 {
		bin = os.Args[2]
	}
	rand.Seed(1)
	for tcase := 0; tcase < 100; tcase++ {
		n := rand.Intn(5) + 1
		letters := []rune("abcde")
		sb := strings.Builder{}
		for i := 0; i < n; i++ {
			sb.WriteRune(letters[rand.Intn(len(letters))])
		}
		start := sb.String()
		shift := rand.Intn(n)
		end := start[shift:] + start[:shift]
		// occasionally produce non-rotation
		if rand.Intn(4) == 0 {
			endBytes := []rune(end)
			pos := rand.Intn(n)
			endBytes[pos] = letters[rand.Intn(len(letters))]
			end = string(endBytes)
		}
		k := int64(rand.Intn(10))
		input := fmt.Sprintf("%s\n%s\n%d\n", start, end, k)
		exp := expected(start, end, k)
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", tcase+1, err)
			return
		}
		var got int64
		fmt.Sscan(out, &got)
		if got != exp {
			fmt.Printf("test %d failed: expected %d got %d\ninput:\n%s", tcase+1, exp, got, input)
			return
		}
	}
	fmt.Println("All tests passed.")
}

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const mod int64 = 1000000007

var (
	fact    []int64
	invfact []int64
	inv     []int64
)

func powmod(a, b int64) int64 {
	res := int64(1)
	for b > 0 {
		if b&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		b >>= 1
	}
	return res
}

func precompute(n int) {
	fact = make([]int64, n+1)
	invfact = make([]int64, n+1)
	inv = make([]int64, n+1)
	fact[0] = 1
	for i := 1; i <= n; i++ {
		fact[i] = fact[i-1] * int64(i) % mod
	}
	invfact[n] = powmod(fact[n], mod-2)
	for i := n; i > 0; i-- {
		invfact[i-1] = invfact[i] * int64(i) % mod
	}
	for i := 1; i <= n; i++ {
		inv[i] = fact[i-1] * invfact[i] % mod
	}
}

func countLess(s string, freqOrig []int) int64 {
	n := len(s)
	freq := make([]int, 26)
	copy(freq, freqOrig)
	remaining := n
	tot := fact[remaining]
	for i := 0; i < 26; i++ {
		tot = tot * invfact[freq[i]] % mod
	}
	ans := int64(0)
	for i := 0; i < n; i++ {
		ch := int(s[i] - 'a')
		for c := 0; c < ch; c++ {
			if freq[c] > 0 {
				add := tot * int64(freq[c]) % mod * inv[remaining] % mod
				ans = (ans + add) % mod
			}
		}
		if freq[ch] == 0 {
			break
		}
		tot = tot * int64(freq[ch]) % mod * inv[remaining] % mod
		freq[ch]--
		remaining--
	}
	return ans
}

func solveCase(a, b string) string {
	n := len(a)
	if len(b) != n {
		return "0"
	}
	freq := make([]int, 26)
	for i := 0; i < n; i++ {
		freq[a[i]-'a']++
	}
	precompute(n)
	lessB := countLess(b, freq)
	lessA := countLess(a, freq)
	ans := (lessB - lessA - 1) % mod
	if ans < 0 {
		ans += mod
	}
	return fmt.Sprintf("%d", ans)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	const testcasesRaw = `cc de
bccac deefd
cb ff
bcaa fdfe
c d
ca ed
bcc efe
ccba eddd
abcb feef
cbcc efde
b f
cb ff
accac feedd
cbab dedd
bba dff
b f
bcbca deddd
caabb fedfd
bba eee
bcccc dffee
bb ef
cba efe
b f
caacc eeefe
cbcba fdfde
cbb ffe
ba ee
bbbaa fffde
acbab dfeff
a f
bca edd
cca fee
aa fd
cab efd
bcabb feeef
bbcb dedd
b f
bccaa fefff
cba dfe
a d
c f
bc dd
caac edfd
cbaca edefe
b d
ac dd
ba de
babab fffed
baaa dddd
acac feed
abb fef
bba eed
ca ff
acaa eeff
bcacb defef
bcb edd
cb ff
b d
aabb fede
c f
bcacb ffdff
b d
ba fd
acbaa edddf
cbcb ddff
caacc dfdfe
aaaac deffe
bcb efe
b f
bc de
ccccb dfedd
b f
ccbc fdde
ac fe
cc ef
cccb dede
baac edee
bc ef
abccc defde
b d
cc fd
ba de
ac ff
b d
b d
acaa deed
bcc fed
b e
bbaa fffe
acba edfe
b f
acc efe
bccc fede
ccacb fedfe
caccc fddef
aaca efff
bcbca eeedf
a f
cc ff
babb eddd`

	scanner := bufio.NewScanner(strings.NewReader(testcasesRaw))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		fields := strings.Fields(line)
		if len(fields) != 2 {
			fmt.Println("bad test case")
			os.Exit(1)
		}
		a := fields[0]
		b := fields[1]
		expected := solveCase(a, b)
		input := fmt.Sprintf("%s\n%s\n", a, b)
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, stderr.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if got != expected {
			fmt.Printf("test %d failed: expected %s got %s\n", idx, expected, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}

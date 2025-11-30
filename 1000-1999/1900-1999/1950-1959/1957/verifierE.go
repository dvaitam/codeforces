package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const (
	N   = 1000000
	MOD = 1000000007
)

var prefix []int64

const testcasesE = `422
509
197
172
401
49
806
971
822
791
985
348
909
779
699
850
534
193
742
336
188
532
946
556
730
872
609
597
445
156
649
706
633
506
921
680
585
952
219
613
452
34
249
500
927
619
864
326
570
749
202
1
45
64
138
249
454
896
691
625
230
677
701
112
883
989
776
529
436
947
662
650
294
907
700
347
506
702
197
168
382
571
864
366
735
449
401
450
680
326
867
480
69
151
226
773
117
770
160
977`

func sieve(n int) []int {
	isPrime := make([]bool, n+1)
	for i := 2; i <= n; i++ {
		isPrime[i] = true
	}
	for i := 2; i*i <= n; i++ {
		if isPrime[i] {
			for j := i * i; j <= n; j += i {
				isPrime[j] = false
			}
		}
	}
	primes := []int{}
	for i := 2; i <= n; i++ {
		if isPrime[i] {
			primes = append(primes, i)
		}
	}
	return primes
}

func init() {
	primes := sieve(N)
	diff := make([]int64, N+2)

	for _, p := range primes {
		for q := 1; q <= N/p; q++ {
			val := int64((p - (q % p)) % p)
			l := q * p
			r := (q+1)*p - 1
			if r > N {
				r = N
			}
			diff[l] += val
			diff[r+1] -= val
		}
	}

	for q := 1; q <= N/4; q++ {
		if q%2 == 1 {
			l := 4 * q
			r := l + 3
			if r > N {
				r = N
			}
			diff[l] += 2
			diff[r+1] -= 2
		}
	}

	prefix = make([]int64, N+1)
	var cur int64
	for i := 1; i <= N; i++ {
		cur += diff[i]
		curMod := cur % MOD
		if curMod < 0 {
			curMod += MOD
		}
		prefix[i] = (prefix[i-1] + curMod) % MOD
	}
}

func expected(n int) string {
	return fmt.Sprintf("%d", prefix[n]%MOD)
}

func buildInput(n int) string {
	return fmt.Sprintf("1\n%d\n", n)
}

func buildBatchInput(nums []int) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%d\n", len(nums)))
	for _, n := range nums {
		b.WriteString(fmt.Sprintf("%d\n", n))
	}
	return b.String()
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: verifierE /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	var nums []int
	lines := strings.Split(testcasesE, "\n")
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		n, err := strconv.Atoi(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d parse error: %v\n", i+1, err)
			os.Exit(1)
		}
		nums = append(nums, n)
	}

	input := buildBatchInput(nums)
	gotRaw, err := runCandidate(bin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "runtime error: %v\n%s\n", err, gotRaw)
		os.Exit(1)
	}
	gotLines := strings.Split(strings.TrimSpace(gotRaw), "\n")
	if len(gotLines) != len(nums) {
		fmt.Fprintf(os.Stderr, "expected %d lines of output, got %d\n", len(nums), len(gotLines))
		os.Exit(1)
	}
	for i, n := range nums {
		want := expected(n)
		if gotLines[i] != want {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%s\nexpected: %s\ngot: %s\n", i+1, buildInput(n), want, gotLines[i])
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(nums))
}

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const mod = 1000000007

// Embedded source for the reference solution (was 1174E.go).
const solutionSource = `package main

import (
   "bufio"
   "fmt"
   "os"
)

const mod = 1000000007

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(in, &n)
   // spf sieve
   spf := make([]int, n+1)
   spf[1] = 1
   primes := make([]int, 0, n/10)
   for i := 2; i <= n; i++ {
       if spf[i] == 0 {
           spf[i] = i
           primes = append(primes, i)
       }
       for _, p := range primes {
           if p > spf[i] || i*p > n {
               break
           }
           spf[i*p] = p
       }
   }
   // omega and find max
   omega := make([]int, n+1)
   omega[1] = 0
   M := 0
   cands := make([]int, 0)
   for i := 2; i <= n; i++ {
       omega[i] = omega[i/spf[i]] + 1
       if omega[i] > M {
           M = omega[i]
           cands = cands[:0]
           cands = append(cands, i)
       } else if omega[i] == M {
           cands = append(cands, i)
       }
   }
   // factorials and inverses
   fact := make([]int, n+1)
  fact[0] = 1
  for i := 1; i <= n; i++ {
      fact[i] = int(int64(fact[i-1]) * int64(i) % mod)
  }
   // inverses up to n
   inv := make([]int, n+1)
   inv[1] = 1
   for i := 2; i <= n; i++ {
       inv[i] = mod - int(int64(mod/i)*int64(inv[mod%i])%mod)
   }
   fullFact := fact[n-1]
   res := 0
   // process each candidate a0
   for _, a0 := range cands {
       // factor a0
       tmp := a0
       // map primes
       pm := make([]int, 0)
       ex := make([]int, 0)
       for tmp > 1 {
           p := spf[tmp]
           cnt := 0
           for tmp%p == 0 {
               tmp /= p
               cnt++
           }
           pm = append(pm, p)
           ex = append(ex, cnt)
       }
       k := len(pm)
       // dims and strides
       dims := make([]int, k)
       strides := make([]int, k)
       total := 1
       for i := 0; i < k; i++ {
           dims[i] = ex[i] + 1
       }
       for i := 0; i < k; i++ {
           strides[i] = total
           total *= dims[i]
       }
       // DP array stores sum of products of b_j/S_j
       dp := make([]int, total)
       dp[0] = 1
       // precompute base count of multiples of a0
       baseCnt := n / a0
       // iterate states
       for idx := 0; idx < total; idx++ {
           curVal := dp[idx]
           if curVal == 0 {
               continue
           }
           // decode r vector and compute sumRemoved and gPrev
           rem := idx
           sumRemoved := 0
           r := make([]int, k)
           gPrev := a0
           for i := 0; i < k; i++ {
               ri := rem % dims[i]
               rem /= dims[i]
               r[i] = ri
               sumRemoved += ri
               for t := 0; t < ri; t++ {
                   gPrev /= pm[i]
               }
           }
           if sumRemoved == M {
               continue
           }
           // compute S_j = (n-1) - (total used so far) = (n-1) - (floor(n/gPrev) - baseCnt)
           used := n/gPrev - baseCnt
           S := (n - 1) - used
           invS := inv[S]
           // transitions: remove one prime pm[i]
           for i := 0; i < k; i++ {
               if r[i] < ex[i] {
                   gNew := gPrev / pm[i]
                   b := n/gNew - n/gPrev
                   mult := int(int64(b) * int64(invS) % mod)
                   nextIdx := idx + strides[i]
                   dp[nextIdx] = (dp[nextIdx] + int(int64(curVal)*int64(mult)%mod)) % mod
               }
           }
       }
       // full removal index
       sumFrac := dp[total-1]
       // add count = sumFrac * (n-1)! 
       res = (res + int(int64(sumFrac)*int64(fullFact)%mod)) % mod
   }
   fmt.Println(res)
}
`

const testcasesRaw = `785
652
283
418
688
147
611
153
415
889
315
525
63
168
130
916
139
934
494
724
659
735
784
742
49
749
845
534
46
845
570
709
764
716
648
398
974
186
354
821
601
834
86
84
575
180
834
273
208
821
269
337
719
727
261
803
267
531
959
469
933
160
776
923
461
567
158
41
649
601
183
658
526
36
919
777
325
839
971
75
199
665
812
469
627
246
833
471
536
165
728
343
936
672
140
489
792
571
60
558`

var _ = solutionSource

func expected(n int) int {
	spf := make([]int, n+1)
	spf[1] = 1
	primes := make([]int, 0, n/10)
	for i := 2; i <= n; i++ {
		if spf[i] == 0 {
			spf[i] = i
			primes = append(primes, i)
		}
		for _, p := range primes {
			if p > spf[i] || i*p > n {
				break
			}
			spf[i*p] = p
		}
	}
	omega := make([]int, n+1)
	omega[1] = 0
	M := 0
	cands := make([]int, 0)
	for i := 2; i <= n; i++ {
		omega[i] = omega[i/spf[i]] + 1
		if omega[i] > M {
			M = omega[i]
			cands = cands[:0]
			cands = append(cands, i)
		} else if omega[i] == M {
			cands = append(cands, i)
		}
	}
	fact := make([]int, n+1)
	fact[0] = 1
	for i := 1; i <= n; i++ {
		fact[i] = int(int64(fact[i-1]) * int64(i) % mod)
	}
	inv := make([]int, n+1)
	inv[1] = 1
	for i := 2; i <= n; i++ {
		inv[i] = mod - int(int64(mod/i)*int64(inv[mod%i])%mod)
	}
	fullFact := fact[n-1]
	res := 0
	for _, a0 := range cands {
		tmp := a0
		pm := make([]int, 0)
		ex := make([]int, 0)
		for tmp > 1 {
			p := spf[tmp]
			cnt := 0
			for tmp%p == 0 {
				tmp /= p
				cnt++
			}
			pm = append(pm, p)
			ex = append(ex, cnt)
		}
		k := len(pm)
		dims := make([]int, k)
		strides := make([]int, k)
		total := 1
		for i := 0; i < k; i++ {
			dims[i] = ex[i] + 1
		}
		for i := 0; i < k; i++ {
			strides[i] = total
			total *= dims[i]
		}
		dp := make([]int, total)
		dp[0] = 1
		baseCnt := n / a0
		for idx := 0; idx < total; idx++ {
			curVal := dp[idx]
			if curVal == 0 {
				continue
			}
			rem := idx
			sumRemoved := 0
			r := make([]int, k)
			gPrev := a0
			for i := 0; i < k; i++ {
				ri := rem % dims[i]
				rem /= dims[i]
				r[i] = ri
				sumRemoved += ri
				for t := 0; t < ri; t++ {
					gPrev /= pm[i]
				}
			}
			if sumRemoved == M {
				continue
			}
			used := n/gPrev - baseCnt
			S := (n - 1) - used
			invS := inv[S]
			for i := 0; i < k; i++ {
				if r[i] < ex[i] {
					gNew := gPrev / pm[i]
					b := n/gNew - n/gPrev
					mult := int(int64(b) * int64(invS) % mod)
					nextIdx := idx + strides[i]
					dp[nextIdx] = (dp[nextIdx] + int(int64(curVal)*int64(mult)%mod)) % mod
				}
			}
		}
		sumFrac := dp[total-1]
		res = (res + int(int64(sumFrac)*int64(fullFact)%mod)) % mod
	}
	return res
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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	scanner := bufio.NewScanner(strings.NewReader(testcasesRaw))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		n, err := strconv.Atoi(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d invalid integer\n", idx)
			os.Exit(1)
		}
		want := fmt.Sprintf("%d", expected(n))
		input := fmt.Sprintf("%d\n", n)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", idx, err)
			os.Exit(1)
		}
		if got != want {
			fmt.Fprintf(os.Stderr, "test %d failed\nexpected: %s\ngot: %s\n", idx, want, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}

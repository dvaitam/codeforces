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

// Embedded source for the reference solution (was 1139D.go).
const solutionSource = `package main

import (
   "bufio"
   "fmt"
   "os"
)

const mod = 1000000007

func modPow(a, e int64) int64 {
   res := int64(1)
   a %= mod
   for e > 0 {
       if e&1 == 1 {
           res = res * a % mod
       }
       a = a * a % mod
       e >>= 1
   }
   return res
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var m int
   if _, err := fmt.Fscan(in, &m); err != nil {
       return
   }
   // Compute Möbius function with sieve
   mu := make([]int, m+1)
   p := make([]bool, m+1)
   primes := make([]int, m+1)
   mu[1] = 1
   pris := 0
   for i := 2; i <= m; i++ {
       if !p[i] {
           pris++
           primes[pris] = i
           mu[i] = 1
       }
       for j := 1; j <= pris; j++ {
           v := primes[j]
           s := i * v
           if s > m {
               break
           }
           p[s] = true
           if i%v != 0 {
               mu[s] = -mu[i]
           } else {
               mu[s] = 0
               break
           }
       }
   }
   var ansA, ansB int64 = 1, 1
   // accumulate contributions for i = 2..m (skip i=1 as per original logic)
   for i := 2; i <= m; i++ {
       if mu[i] == 0 {
           continue
       }
       // count multiples and non-multiples of i
       cntA := int64(m / i)
       cntB := int64(m - m/i)
       if mu[i] == -1 {
           cntA = (mod - cntA) % mod
       }
       // update numerator and denominator
       ansA = (ansA*cntB + ansB*cntA) % mod
       ansB = ansB * cntB % mod
   }
   // result = ansA / ansB mod
   res := ansA * modPow(ansB, mod-2) % mod
   fmt.Println(res)
}
`

const testcasesRaw = `244
607
558
134
379
938
619
486
641
595
68
621
14
931
858
481
266
565
240
197
735
482
554
857
563
488
407
655
882
155
238
651
156
889
949
536
400
760
16
688
796
66
164
777
981
606
44
309
799
32
844
887
276
485
610
737
943
900
397
732
808
944
438
405
746
821
591
456
988
959
138
900
375
100
37
140
507
223
265
989
689
447
798
642
876
309
432
520
854
396
588
360
547
600
418
599
238
926
345
699`

func modPow(a, e int64) int64 {
	res := int64(1)
	a %= mod
	for e > 0 {
		if e&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		e >>= 1
	}
	return res
}

func expected(m int) int64 {
	mu := make([]int, m+1)
	p := make([]bool, m+1)
	primes := make([]int, m+1)
	mu[1] = 1
	pris := 0
	for i := 2; i <= m; i++ {
		if !p[i] {
			pris++
			primes[pris] = i
			mu[i] = 1
		}
		for j := 1; j <= pris; j++ {
			v := primes[j]
			s := i * v
			if s > m {
				break
			}
			p[s] = true
			if i%v != 0 {
				mu[s] = -mu[i]
			} else {
				mu[s] = 0
				break
			}
		}
	}
	var ansA, ansB int64 = 1, 1
	for i := 2; i <= m; i++ {
		if mu[i] == 0 {
			continue
		}
		cntA := int64(m / i)
		cntB := int64(m - m/i)
		if mu[i] == -1 {
			cntA = (mod - cntA) % mod
		}
		ansA = (ansA*cntB + ansB*cntA) % mod
		ansB = ansB * cntB % mod
	}
	return ansA * modPow(ansB, mod-2) % mod
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
	var _ = solutionSource
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
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
		m, err := strconv.Atoi(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: invalid integer\n", idx)
			os.Exit(1)
		}
		input := fmt.Sprintf("%d\n", m)
		want := fmt.Sprintf("%d", expected(m))
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Printf("test %d failed\nexpected: %s\ngot: %s\n", idx, want, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}

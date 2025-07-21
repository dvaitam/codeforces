package main

import (
   "bufio"
   "fmt"
   "os"
)

const MOD = 1000000007

// fast doubling Fibonacci: returns (F(n), F(n+1)) mod MOD
func fibPair(n int64) (int64, int64) {
   if n == 0 {
       return 0, 1
   }
   a, b := fibPair(n >> 1)
   c := (a * ((b*2%MOD - a + MOD) % MOD)) % MOD
   d := (a*a%MOD + b*b%MOD) % MOD
   if n&1 == 0 {
       return c, d
   }
   return d, (c + d) % MOD
}

// build KMP prefix function
func buildPi(pat string) []int {
   m := len(pat)
   pi := make([]int, m)
   j := 0
   for i := 1; i < m; i++ {
       for j > 0 && pat[i] != pat[j] {
           j = pi[j-1]
       }
       if pat[i] == pat[j] {
           j++
       }
       pi[i] = j
   }
   return pi
}

// count occurrences of pat in text using KMP
func countOcc(text, pat string, pi []int) int64 {
   n, m := len(text), len(pat)
   j := 0
   cnt := int64(0)
   for i := 0; i < n; i++ {
       for j > 0 && text[i] != pat[j] {
           j = pi[j-1]
       }
       if text[i] == pat[j] {
           j++
       }
       if j == m {
           cnt++
           j = pi[j-1]
       }
   }
   return cnt
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var k int64
   var m int
   fmt.Fscan(in, &k, &m)
   patterns := make([]string, m)
   maxL := 0
   for i := 0; i < m; i++ {
       fmt.Fscan(in, &patterns[i])
       if len(patterns[i]) > maxL {
           maxL = len(patterns[i])
       }
   }
   // precompute Fibonacci lengths until >= maxL
   fibLen := []uint64{0, 1, 1}
   Nmax := 2
   for fibLen[Nmax] < uint64(maxL) {
       fibLen = append(fibLen, fibLen[Nmax]+fibLen[Nmax-1])
       if fibLen[Nmax+1] > (1<<63-1) {
           fibLen[Nmax+1] = 1<<63 - 1
       }
       Nmax++
   }
   // need up to Nmax+2
   upto := Nmax + 2
   // build Fibonacci strings up to upto
   fstr := make([]string, upto+1)
   fstr[1] = "a"
   fstr[2] = "b"
   for i := 3; i <= upto; i++ {
       fstr[i] = fstr[i-1] + fstr[i-2]
   }
   // process each pattern
   for _, pat := range patterns {
       L := len(pat)
       // find N0: first index with fibLen >= L
       N0 := 1
       for uint64(L) > fibLen[N0] {
           N0++
       }
       t0 := N0 + 2
       // build pi for pattern
       pi := buildPi(pat)
       // compute C[n] for n < t0
       C := make([]int64, t0)
       for i := 1; i < t0 && i < len(fstr); i++ {
           C[i] = countOcc(fstr[i], pat, pi)
       }
       // compute cross term K for n >= t0
       K := int64(0)
       if t0 < len(fstr) {
           suf := fstr[t0-1]
           if len(suf) > L-1 {
               suf = suf[len(suf)-(L-1):]
           }
           pre := fstr[t0-2]
           if len(pre) > L-1 {
               pre = pre[:L-1]
           }
           // count occurrences of pat in suf+pre that cross boundary
           T := suf + pre
           n := len(suf)
           j := 0
           for i, ch := range T {
               for j > 0 && byte(ch) != pat[j] {
                   j = pi[j-1]
               }
               if byte(ch) == pat[j] {
                   j++
               }
               if j == L {
                   start := i - L + 1
                   if start < n && i >= n {
                       K++
                   }
                   j = pi[j-1]
               }
           }
       }
       var ans int64
       if k < int64(t0) {
           ans = C[k] % MOD
       } else {
           // D(n) = C(n) + K satisfy D(n)=D(n-1)+D(n-2)
           D0 := (C[t0-1] + K) % MOD
           D1 := (C[t0-2] + K) % MOD
           // offset = k-(t0-1)
           offset := k - int64(t0-1)
           Foff, Foff1 := fibPair(offset)   // F(offset), F(offset+1)
           // Dk = D0*Foff1 + D1*Foff
           Dk := (D0*Foff1 + D1*Foff) % MOD
           ans = (Dk - K%MOD + MOD) % MOD
       }
       fmt.Fprintln(out, ans)
   }
}

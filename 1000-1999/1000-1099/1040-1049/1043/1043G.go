package main

import (
   "bufio"
   "fmt"
   "os"
)

// double rolling hash moduli
const (
   mod1 = 1000000007
   mod2 = 1000000009
   base = 91138233
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n int
   fmt.Fscan(in, &n)
   var s string
   fmt.Fscan(in, &s)
   // 1-based
   h1 := make([]int, n+1)
   h2 := make([]int, n+1)
   pow1 := make([]int, n+1)
   pow2 := make([]int, n+1)
   pow1[0], pow2[0] = 1, 1
   for i := 1; i <= n; i++ {
       pow1[i] = int((int64(pow1[i-1]) * base) % mod1)
       pow2[i] = int((int64(pow2[i-1]) * base) % mod2)
   }
   for i := 1; i <= n; i++ {
       c := int(s[i-1])
       h1[i] = int((int64(h1[i-1])*base + int64(c)) % mod1)
       h2[i] = int((int64(h2[i-1])*base + int64(c)) % mod2)
   }
   // prefix char counts
   cnt := make([][26]int, n+1)
   for i := 1; i <= n; i++ {
       for c := 0; c < 26; c++ {
           cnt[i][c] = cnt[i-1][c]
       }
       cnt[i][s[i-1]-'a']++
   }

   q := 0
   fmt.Fscan(in, &q)
   for qi := 0; qi < q; qi++ {
       var l, r int
       fmt.Fscan(in, &l, &r)
       k := r - l + 1
       if k < 2 {
           fmt.Fprintln(out, -1)
           continue
       }
       // compute distinct count and any repeat
       distinct := 0
       hasRepeat := false
       for c := 0; c < 26; c++ {
           f := cnt[r][c] - cnt[l-1][c]
           if f > 0 {
               distinct++
               if f > 1 {
                   hasRepeat = true
               }
           }
       }
       if !hasRepeat {
           fmt.Fprintln(out, -1)
           continue
       }
       // periodic check: uniform char
       periodic := false
       if distinct == 1 {
           periodic = true
       } else {
           // helper to get hash of s[l..r]
           getHash := func(lm, rm, d int) bool {
               // compare s[lm..rm-d] == s[lm+d..rm]
               // here d < (rm-lm+1)
               // compute hash1
               h1a := (h1[rm-d] - int(int64(h1[lm-1])*int64(pow1[rm-d-(lm-1)])%mod1) + mod1) % mod1
               h1b := (h1[rm] - int(int64(h1[lm+d-1])*int64(pow1[rm-(lm+d-1)])%mod1) + mod1) % mod1
               if h1a != h1b {
                   return false
               }
               h2a := (h2[rm-d] - int(int64(h2[lm-1])*int64(pow2[rm-d-(lm-1)])%mod2) + mod2) % mod2
               h2b := (h2[rm] - int(int64(h2[lm+d-1])*int64(pow2[rm-(lm+d-1)])%mod2) + mod2) % mod2
               return h2a == h2b
           }
           // enumerate divisors
           for di := 1; di*di <= k; di++ {
               if k%di != 0 {
                   continue
               }
               // check smallest period di
               if di > 1 && di < k {
                   if getHash(l, r, di) {
                       periodic = true
                       break
                   }
               }
               d2 := k / di
               if d2 > 1 && d2 < k && d2 != di {
                   if getHash(l, r, d2) {
                       periodic = true
                       break
                   }
               }
           }
       }
       if periodic {
           fmt.Fprintln(out, 1)
           continue
       }
       // check for answer 2
       // boundary repeats
       if s[l-1] == s[r-1] || (l+1 <= r && s[l-1] == s[l]) || (r-1 >= l && s[r-1] == s[r-2]) {
           fmt.Fprintln(out, 2)
       } else {
           fmt.Fprintln(out, 3)
       }
   }
}

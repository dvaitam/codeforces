package main

import (
   "bufio"
   "fmt"
   "os"
)

const md = 2006091501

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

func exkmp(S, T string) []int {
   n, m := len(S), len(T)
   nxt := make([]int, m)
   // Z-array for T
   nxt[0] = m
   l, r := 0, 0
   for i := 1; i < m; i++ {
       if i <= r {
           nxt[i] = min(r-i+1, nxt[i-l])
       }
       for i+nxt[i] < m && T[nxt[i]] == T[i+nxt[i]] {
           nxt[i]++
       }
       if i+nxt[i]-1 > r {
           l = i
           r = i + nxt[i] - 1
       }
   }
   // match lengths for S
   le := make([]int, n)
   l, r = 0, 0
   for i := 0; i < n; i++ {
       if i <= r {
           le[i] = min(r-i+1, nxt[i-l])
       }
       for i+le[i] < n && le[i] < m && S[i+le[i]] == T[le[i]] {
           le[i]++
       }
       if i+le[i]-1 > r {
           l = i
           r = i + le[i] - 1
       }
   }
   return le
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var S, T string
   if _, err := fmt.Fscan(in, &S, &T); err != nil {
       return
   }
   n, m := len(S), len(T)
   le := exkmp(S, T)
   // prefix hash and powers
   ha := make([]int64, n+1)
   mi := make([]int64, n+1)
   mi[0] = 1
   var ht int64
   for i := 0; i < m; i++ {
       ht = (ht*10 + int64(T[i]-'0')) % md
   }
   for i := 0; i < n; i++ {
       mi[i+1] = (mi[i] * 10) % md
       ha[i+1] = (ha[i]*10 + int64(S[i]-'0')) % md
   }
   // check function
   check := func(a, b, c int) {
       if a < 0 || c >= n {
           return
       }
       h1 := (ha[b+1] - ha[a]*mi[b-a+1]%md + md) % md
       h2 := (ha[c+1] - ha[b+1]*mi[c-b]%md + md) % md
       if (h1+h2)%md != ht {
           return
       }
       fmt.Printf("%d %d\n%d %d", a+1, b+1, b+2, c+1)
       os.Exit(0)
   }
   // main logic
   for l := 0; l+m <= n; l++ {
       r := l + m - 1
       c := le[l]
       if c == m || S[l+c] > T[c] {
           continue
       }
       z := m - c
       check(l, r, r+z)
       check(l-z, l-1, r)
       z--
       if z > 0 {
           check(l, r, r+z)
           check(l-z, l-1, r)
       }
   }
   if m > 1 {
       for l := 0; l+2*(m-1) <= n; l++ {
           check(l, l+m-2, l+2*m-3)
       }
   }
}

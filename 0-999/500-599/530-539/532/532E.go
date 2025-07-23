package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(reader, &n)
   s := make([]byte, n)
   t := make([]byte, n)
   fmt.Fscan(reader, &s)
   fmt.Fscan(reader, &t)
   // find first mismatch from left
   preMM := n
   for i := 0; i < n; i++ {
       if s[i] != t[i] {
           preMM = i
           break
       }
   }
   // find last mismatch from right
   mmSuffix := -1
   for i := n - 1; i >= 0; i-- {
       if s[i] != t[i] {
           mmSuffix = i
           break
       }
   }
   // prepare eq arrays and f arrays
   f1 := make([]int, n)
   f2 := make([]int, n)
   // eq1: s[i]==t[i+1], eq2: t[i]==s[i+1]
   // f[i]: length of consecutive eq starting at i
   // last position
   if n > 1 {
       f1[n-1] = 0
       f2[n-1] = 0
   }
   for i := n - 2; i >= 0; i-- {
       if s[i] == t[i+1] {
           f1[i] = f1[i+1] + 1
       } else {
           f1[i] = 0
       }
       if t[i] == s[i+1] {
           f2[i] = f2[i+1] + 1
       } else {
           f2[i] = 0
       }
   }
   // sufTrue for j from 0..n
   sufTrue := make([]int, n+1)
   for j := 0; j <= n; j++ {
       if j > mmSuffix {
           sufTrue[j] = 1
       } else {
           sufTrue[j] = 0
       }
   }
   // prefix sum on sufTrue
   pref := make([]int, len(sufTrue)+1)
   for i := 0; i < len(sufTrue); i++ {
       pref[i+1] = pref[i] + sufTrue[i]
   }
   // use hash set to count distinct W
   const B = 1315423911
   // precompute powers
   pow := make([]uint64, n+2)
   pow[0] = 1
   for i := 1; i <= n+1; i++ {
       pow[i] = pow[i-1] * B
   }
   // prefix hashes
   hS := make([]uint64, n+1)
   hT := make([]uint64, n+1)
   for i := 0; i < n; i++ {
       hS[i+1] = hS[i]*B + uint64(s[i])
       hT[i+1] = hT[i]*B + uint64(t[i])
   }
   seen := make(map[uint64]struct{})
   var count int
   // first pass: s->t
   for i := 0; i < n; i++ {
       if i > preMM {
           break
       }
       lo, hi := i+1, i+1+f1[i]
       if hi > n {
           hi = n
       }
       if lo <= hi && pref[hi+1]-pref[lo] > 0 {
           // build hash of W = S[0..i-1] + T[i] + S[i..n-1]
           // part1: S[0..i-1]
           h := hS[i] * pow[n+1-i]
           // part2: T[i]
           h += uint64(t[i]) * pow[n-i]
           // part3: S[i..n-1]
           suf := hS[n] - hS[i]*pow[n-i]
           h += suf
           if _, ok := seen[h]; !ok {
               seen[h] = struct{}{}
               count++
           }
       }
   }
   // second pass: t->s
   for i := 0; i < n; i++ {
       if i > preMM {
           break
       }
       lo, hi := i+1, i+1+f2[i]
       if hi > n {
           hi = n
       }
       if lo <= hi && pref[hi+1]-pref[lo] > 0 {
           // build hash of W = T[0..i-1] + S[i] + T[i..n-1]
           h := hT[i] * pow[n+1-i]
           h += uint64(s[i]) * pow[n-i]
           suf := hT[n] - hT[i]*pow[n-i]
           h += suf
           if _, ok := seen[h]; !ok {
               seen[h] = struct{}{}
               count++
           }
       }
   }
   fmt.Println(count)
}

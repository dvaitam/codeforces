package main

import (
   "bufio"
   "bytes"
   "fmt"
   "io"
   "os"
   "sort"
)

const MOD = 1000000007

// Rolling hash base
const base = 91138233

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(reader, &n)
   words := make([][]byte, n)
   maxLen := 0
   for i := 0; i < n; i++ {
       line, _ := reader.ReadBytes('\n')
       line = bytes.TrimSpace(line)
       words[i] = line
       if len(line) > maxLen {
           maxLen = len(line)
       }
   }
   // precompute powers
   pow := make([]uint64, maxLen+2)
   pow[0] = 1
   for i := 1; i < len(pow); i++ {
       pow[i] = pow[i-1] * base
   }
   // dp for previous
   var dpPrev []int
   // initial dpPrev for first word
   L0 := len(words[0])
   dpPrev = make([]int, L0+1)
   for p := 0; p <= L0; p++ {
       dpPrev[p] = 1
   }
   // iterate pairs
   for idx := 1; idx < n; idx++ {
       s := words[idx-1]
       t := words[idx]
       hs := buildHash(s)
       ht := buildHash(t)
       ls := len(s)
       lt := len(t)
       // prev variants positions: 0..ls (where pos ls means no deletion)
       ps := make([]int, ls+1)
       for i := range ps {
           ps[i] = i
       }
       // current variants: 0..lt
       qs := make([]int, lt+1)
       for i := range qs {
           qs[i] = i
       }
       // sort prev variants
       sort.Slice(ps, func(i, j int) bool {
           return cmpVariant(s, hs, ps[i], s, hs, ps[j], pow) < 0
       })
       // sort curr variants
       sort.Slice(qs, func(i, j int) bool {
           return cmpVariant(t, ht, qs[i], t, ht, qs[j], pow) < 0
       })
       // dpCurr
       dpCurr := make([]int, lt+1)
       sum := 0
       pi := 0
       for _, q := range qs {
           // advance ps while s^ps <= t^q
           for pi < len(ps) && cmpVariant(s, hs, ps[pi], t, ht, q, pow) <= 0 {
               sum = (sum + dpPrev[ps[pi]]) % MOD
               pi++
           }
           dpCurr[q] = sum
       }
       dpPrev = dpCurr
   }
   // sum dpPrev
   ans := 0
   for _, v := range dpPrev {
       ans = (ans + v) % MOD
   }
   fmt.Println(ans)
}

// buildHash builds prefix hash array for s
func buildHash(s []byte) []uint64 {
   n := len(s)
   h := make([]uint64, n+1)
   for i := 0; i < n; i++ {
       h[i+1] = h[i]*base + uint64(s[i]) + 1
   }
   return h
}

// getHash returns hash of substring s[l:r] (0-indexed, l<=r)
func getHash(h []uint64, l, r int) uint64 {
   return h[r] - h[l]*powGlobal[r-l]
}

var powGlobal []uint64

// cmpVariant compares s1 with deletion p1 and s2 with deletion p2
// returns -1 if s1^p1 < s2^p2, 0 if equal, 1 if >
func cmpVariant(s1 []byte, h1 []uint64, p1 int, s2 []byte, h2 []uint64, p2 int, pow []uint64) int {
   // set global pow
   powGlobal = pow
   n1 := len(s1); n2 := len(s2)
   l1 := n1; if p1 < n1 { l1 = n1-1 }
   l2 := n2; if p2 < n2 { l2 = n2-1 }
   // binary search lcp
   lo, hi := 0, min(l1, l2)+1
   for lo+1 < hi {
       mid := (lo + hi) >> 1
       if hashPrefix(s1, h1, p1, mid, pow) == hashPrefix(s2, h2, p2, mid, pow) {
           lo = mid
       } else {
           hi = mid
       }
   }
   lcp := lo
   if lcp == l1 && lcp == l2 {
       return 0
   }
   if lcp == l1 {
       return -1
   }
   if lcp == l2 {
       return 1
   }
   // get chars at lcp
   c1 := getChar(s1, p1, lcp)
   c2 := getChar(s2, p2, lcp)
   if c1 < c2 {
       return -1
   } else if c1 > c2 {
       return 1
   }
   return 0
}

func min(a, b int) int {
   if a < b { return a }
   return b
}

// hashPrefix returns hash of first l chars of s^p
func hashPrefix(s []byte, h []uint64, p, l int, pow []uint64) uint64 {
   // use global powGlobal
   if l <= p {
       // s[0:l]
       return h[l]
   }
   // prefix s[0:p]
   hp := h[p]
   // suffix part: need l-p chars from s[p+1:]
   // get hash of s[p+1 : p+1+(l-p)]
   h2 := h[p+1+(l-p)] - h[p+1]*powGlobal[l-p]
   return hp*powGlobal[l-p] + h2
}

// getChar returns char at index i of s^p
func getChar(s []byte, p, i int) byte {
   if i < p {
       return s[i]
   }
   return s[i+1]
}

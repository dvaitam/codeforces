package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var s1, s2 string
   if _, err := fmt.Fscanln(reader, &s1); err != nil {
       return
   }
   if _, err := fmt.Fscanln(reader, &s2); err != nil {
       return
   }
   n1, n2 := len(s1), len(s2)
   // choose base and precompute powers
   const base = 91138233
   maxn := n1
   if n2 > maxn {
       maxn = n2
   }
   pow := make([]uint64, maxn+2)
   pow[0] = 1
   for i := 1; i <= maxn; i++ {
       pow[i] = pow[i-1] * base
   }
   // prefix hashes
   h1 := make([]uint64, n1+1)
   for i := 0; i < n1; i++ {
       h1[i+1] = h1[i]*base + uint64(s1[i])
   }
   h2 := make([]uint64, n2+1)
   for i := 0; i < n2; i++ {
       h2[i+1] = h2[i]*base + uint64(s2[i])
   }
   // helper to get hash
   get := func(h []uint64, l, r int) uint64 {
       return h[r] - h[l-1]*pow[r-l+1]
   }
   // scan lengths
   minLen := -1
   maxL := n1
   if n2 < maxL {
       maxL = n2
   }
   for L := 1; L <= maxL; L++ {
       cnt1 := make(map[uint64]int)
       for i := 1; i+L-1 <= n1; i++ {
           h := get(h1, i, i+L-1)
           cnt1[h]++
       }
       cnt2 := make(map[uint64]int)
       for i := 1; i+L-1 <= n2; i++ {
           h := get(h2, i, i+L-1)
           cnt2[h]++
       }
       found := false
       for h, c1 := range cnt1 {
           if c1 == 1 {
               if c2, ok := cnt2[h]; ok && c2 == 1 {
                   found = true
                   break
               }
           }
       }
       if found {
           minLen = L
           break
       }
   }
   if minLen < 0 {
       fmt.Println(-1)
   } else {
       fmt.Println(minLen)
   }
}

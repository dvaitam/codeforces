package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(in, &n)
   A := make([]string, n)
   maxLen := 0
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &A[i])
       if l := len(A[i]); l > maxLen {
           maxLen = l
       }
   }
   // prepare rolling hash powers
   maxH := maxLen * 2
   pow := make([]uint64, maxH+1)
   const B = 1315423911
   pow[0] = 1
   for i := 1; i <= maxH; i++ {
       pow[i] = pow[i-1] * B
   }
   ans := 0
   // for each subarray
   for l := 0; l < n; l++ {
       // build word set incrementally
       wordSet := make(map[string]struct{})
       lengths := make(map[int]struct{})
       // track max word length in this subarray for possible pow limit
       localMaxLen := 0
       for r := l; r < n; r++ {
           w := A[r]
           wordSet[w] = struct{}{}
           wl := len(w)
           lengths[wl] = struct{}{}
           if wl > localMaxLen {
               localMaxLen = wl
           }
           // prepare length list
           Ls := make([]int, 0, len(lengths))
           for L := range lengths {
               Ls = append(Ls, L)
           }
           // build hash maps by length
           hashMap := make(map[int]map[uint64]struct{})
           for s := range wordSet {
               h := uint64(0)
               for i := 0; i < len(s); i++ {
                   h = h*B + uint64(s[i])
               }
               l0 := len(s)
               m0, ok := hashMap[l0]
               if !ok {
                   m0 = make(map[uint64]struct{})
                   hashMap[l0] = m0
               }
               m0[h] = struct{}{}
           }
           stable := true
           // test all pairs
           for x := range wordSet {
               if !stable {
                   break
               }
               for y := range wordSet {
                   if !stable {
                       break
                   }
                   S := x + y
                   tot := len(S)
                   // compute prefix hash
                   hpre := make([]uint64, tot+1)
                   hpre[0] = 0
                   for i := 0; i < tot; i++ {
                       hpre[i+1] = hpre[i]*B + uint64(S[i])
                   }
                   // dp prefix and suffix
                   dpP := make([]bool, tot+1)
                   dpE := make([]bool, tot+1)
                   dpP[0] = true
                   for i := 0; i < tot; i++ {
                       if !dpP[i] {
                           continue
                       }
                       for _, L := range Ls {
                           j := i + L
                           if j <= tot {
                               // hash of S[i:j]
                               h := hpre[j] - hpre[i]*pow[L]
                               if m0, ok := hashMap[L]; ok {
                                   if _, ok2 := m0[h]; ok2 {
                                       dpP[j] = true
                                   }
                               }
                           }
                       }
                   }
                   dpE[tot] = true
                   for i := tot - 1; i >= 0; i-- {
                       for _, L := range Ls {
                           j := i + L
                           if j <= tot && dpE[j] {
                               h := hpre[j] - hpre[i]*pow[L]
                               if m0, ok := hashMap[L]; ok {
                                   if _, ok2 := m0[h]; ok2 {
                                       dpE[i] = true
                                       break
                                   }
                               }
                           }
                       }
                   }
                   // check rotations
                   lx := len(x)
                   for sft := 1; sft < tot; sft++ {
                       if sft == lx {
                           continue
                       }
                       if dpP[sft] && dpE[sft] {
                           stable = false
                           break
                       }
                   }
               }
           }
           if stable {
               ans++
           }
       }
   }
   fmt.Println(ans)
}

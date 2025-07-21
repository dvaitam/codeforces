package main

import (
   "bufio"
   "fmt"
   "math/bits"
   "os"
   "strings"
)

// bitset LCS based on Myers' algorithm
func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(in, &n)
   bs := make([]string, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &bs[i])
   }
   var m int
   fmt.Fscan(in, &m)
   idx := make([]int, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(in, &idx[i])
       idx[i]--
   }
   var s string
   fmt.Fscan(in, &s)
   // build Hyper String t (may be large)
   var b strings.Builder
   for _, id := range idx {
       b.WriteString(bs[id])
   }
   t := b.String()
   // prepare bitsets
   k := len(s)
   if k == 0 || len(t) == 0 {
       fmt.Println(0)
       return
   }
   words := (k + 63) >> 6
   // char masks for s
   charMask := make([][]uint64, 26)
   for c := 0; c < 26; c++ {
       charMask[c] = make([]uint64, words)
   }
   for i := 0; i < k; i++ {
       c := s[i] - 'a'
       w := i >> 6
       off := uint(i & 63)
       charMask[c][w] |= 1 << off
   }
   D := make([]uint64, words)
   // temp slices
   X := make([]uint64, words)
   tmp := make([]uint64, words)
   // process t
   for i := 0; i < len(t); i++ {
       c := t[i] - 'a'
       M := charMask[c]
       // X = M | D
       for j := 0; j < words; j++ {
           X[j] = M[j] | D[j]
       }
       // D = (D << 1) | 1
       var carry uint64 = 0
       for j := 0; j < words; j++ {
           newCarry := D[j] >> 63
           D[j] = (D[j] << 1) | carry
           carry = newCarry
       }
       D[0] |= 1
       // tmp = X - D
       var borrow uint64 = 0
       for j := 0; j < words; j++ {
           dj := D[j]
           xj := X[j]
           sub := xj - dj - borrow
           // borrow if xj < dj+borrow
           if xj < dj+borrow {
               borrow = 1
           } else {
               borrow = 0
           }
           tmp[j] = sub
       }
       // D = X & ^tmp
       for j := 0; j < words; j++ {
           D[j] = X[j] &^ tmp[j]
       }
   }
   // count bits in D
   var res int
   for j := 0; j < words; j++ {
       res += bitsOnesCount64(D[j])
   }
   fmt.Println(res)
}

// bitsOnesCount64 counts bits set in x
func bitsOnesCount64(x uint64) int {
   // use builtin
   return bits.OnesCount64(x)
}

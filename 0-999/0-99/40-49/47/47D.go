package main

import (
   "bufio"
   "fmt"
   "math/bits"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   // Read attempts
   s1Masks := make([]uint64, m)
   s2Masks := make([]uint64, m)
   c := make([]int, m)
   // split size
   k := n / 2
   r := n - k
   for i := 0; i < m; i++ {
       // read si: could be one token of length n, or n tokens of "0"/"1"
       var tok string
       fmt.Fscan(reader, &tok)
       bitsRow := make([]byte, n)
       if len(tok) == n {
           for j := 0; j < n; j++ {
               bitsRow[j] = tok[j]
           }
       } else {
           // first token is one bit
           bitsRow[0] = tok[0]
           for j := 1; j < n; j++ {
               fmt.Fscan(reader, &tok)
               bitsRow[j] = tok[0]
           }
       }
       fmt.Fscan(reader, &c[i])
       // build masks
       var m1, m2 uint64
       for j := 0; j < k; j++ {
           if bitsRow[j] == '1' {
               m1 |= 1 << j
           }
       }
       for j := k; j < n; j++ {
           if bitsRow[j] == '1' {
               m2 |= 1 << (j - k)
           }
       }
       s1Masks[i] = m1
       s2Masks[i] = m2
   }
   // precompute base multipliers
   pow6 := make([]int, m)
   pow6[0] = 1
   for i := 1; i < m; i++ {
       pow6[i] = pow6[i-1] * 6
   }
   // map for first half vectors
   counts := make(map[int]int)
   lim1 := 1 << k
   for mask := 0; mask < lim1; mask++ {
       key := 0
       ok := true
       for i := 0; i < m; i++ {
           // matches in first half = k - hamming(mask, s1Masks[i])
           diff := bits.OnesCount64(uint64(mask) ^ s1Masks[i])
           match := k - diff
           if match > c[i] {
               ok = false
               break
           }
           key += match * pow6[i]
       }
       if ok {
           counts[key]++
       }
   }
   // process second half
   var ans int64
   lim2 := 1 << r
   for mask := 0; mask < lim2; mask++ {
       key := 0
       ok := true
       for i := 0; i < m; i++ {
           diff := bits.OnesCount64(uint64(mask) ^ s2Masks[i])
           match := r - diff
           if match > c[i] {
               ok = false
               break
           }
           need := c[i] - match
           key += need * pow6[i]
       }
       if ok {
           if v, found := counts[key]; found {
               ans += int64(v)
           }
       }
   }
   fmt.Println(ans)
}

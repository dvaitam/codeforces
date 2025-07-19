package main

import (
   "bufio"
   "fmt"
   "os"
)

const (
   P  = int64(13331)
   M1 = int64(1000000007)
   M2 = int64(998244353)
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var T int
   fmt.Fscan(reader, &T)
   for T > 0 {
       T--
       var n, m int
       fmt.Fscan(reader, &n, &m)
       s := make([][]byte, n)
       for i := 0; i < n; i++ {
           var str string
           fmt.Fscan(reader, &str)
           s[i] = []byte(str)
       }
       // precompute powers
       p1 := make([]int64, n)
       p2 := make([]int64, n)
       p1[0], p2[0] = 1, 1
       for i := 1; i < n; i++ {
           p1[i] = p1[i-1] * P % M1
           p2[i] = p2[i-1] * P % M2
       }
       // map of hash pair to count
       mp := make(map[uint64]int)
       // count all flip scenarios
       for j := 0; j < m; j++ {
           var h1, h2 int64
           for i := 0; i < n; i++ {
               h1 = (h1*P + int64(s[i][j])) % M1
               h2 = (h2*P + int64(s[i][j])) % M2
           }
           for i := 0; i < n; i++ {
               // delta when flipping bit
               delta1 := int64(int(s[i][j]^1) - int(s[i][j]) + int(M1))
               t1 := (h1 + delta1*p1[n-1-i]) % M1
               if t1 < 0 {
                   t1 += M1
               }
               delta2 := int64(int(s[i][j]^1) - int(s[i][j]) + int(M2))
               t2 := (h2 + delta2*p2[n-1-i]) % M2
               if t2 < 0 {
                   t2 += M2
               }
               key := uint64(t1)<<32 | uint64(t2)
               mp[key]++
           }
       }
       // find best count
       best := 0
       for _, cnt := range mp {
           if cnt > best {
               best = cnt
           }
       }
       // find and output one solution
   FOUND:
       for j := 0; j < m; j++ {
           var h1, h2 int64
           for i := 0; i < n; i++ {
               h1 = (h1*P + int64(s[i][j])) % M1
               h2 = (h2*P + int64(s[i][j])) % M2
           }
           for i := 0; i < n; i++ {
               delta1 := int64(int(s[i][j]^1) - int(s[i][j]) + int(M1))
               t1 := (h1 + delta1*p1[n-1-i]) % M1
               if t1 < 0 {
                   t1 += M1
               }
               delta2 := int64(int(s[i][j]^1) - int(s[i][j]) + int(M2))
               t2 := (h2 + delta2*p2[n-1-i]) % M2
               if t2 < 0 {
                   t2 += M2
               }
               key := uint64(t1)<<32 | uint64(t2)
               if mp[key] == best {
                   // apply flip and output column
                   s[i][j] ^= 1
                   // build output
                   for k := 0; k < n; k++ {
                       writer.WriteByte(s[k][j])
                   }
                   writer.WriteByte('\n')
                   break FOUND
               }
           }
       }
   }
}

package main

import (
   "bufio"
   "fmt"
   "os"
)

const mod = 1000000007

// fast exponentiation
func modPow(a, e int64) int64 {
   res := int64(1)
   a %= mod
   for e > 0 {
       if e&1 == 1 {
           res = res * a % mod
       }
       a = a * a % mod
       e >>= 1
   }
   return res
}

// compute Z-array for string s
func makeZ(s string) []int {
   n := len(s)
   Z := make([]int, n)
   Z[0] = n
   l, r := 0, 0
   for i := 1; i < n; i++ {
       if i <= r {
           k := i - l
           if Z[k] < r-i+1 {
               Z[i] = Z[k]
           } else {
               j := r + 1
               for j < n && s[j] == s[j-i] {
                   j++
               }
               Z[i] = j - i
               l, r = i, j-1
           }
       } else {
           j := 0
           for i+j < n && s[j] == s[i+j] {
               j++
           }
           Z[i] = j
           if j > 0 {
               l, r = i, i+j-1
           }
       }
   }
   return Z
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   fmt.Fscan(reader, &n, &m)
   var p string
   fmt.Fscan(reader, &p)
   lenP := len(p)
   ys := make([]int, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(reader, &ys[i])
   }
   // compute Z-array of p for overlap checks
   Z := makeZ(p)
   // validate overlaps and compute covered length U
   var U int64 = 0
   lastEnd := 0
   for i, y := range ys {
       // positions are 1-indexed
       start := y
       end := y + lenP - 1
       if i > 0 {
           prev := ys[i-1]
           d := start - prev
           if d < lenP {
               // need prefix of length lenP-d to match suffix
               if Z[d] < lenP-d {
                   fmt.Println(0)
                   return
               }
           }
       }
       // merge interval [start, end]
       if start > lastEnd {
           U += int64(lenP)
       } else {
           // overlap = lastEnd - start + 1
           ov := lastEnd - start + 1
           U += int64(lenP - ov)
       }
       if end > lastEnd {
           lastEnd = end
       }
   }
   free := int64(n) - U
   if free < 0 {
       free = 0
   }
   ans := modPow(26, free)
   fmt.Println(ans)
}

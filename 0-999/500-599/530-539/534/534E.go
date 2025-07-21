package main

import (
   "bufio"
   "encoding/binary"
   "fmt"
   "io"
   "math"
   "os"
   "time"
)

// splitmix64 is a deterministic hash function
func splitmix64(x uint64) uint64 {
   x += 0x9e3779b97f4a7c15
   z := x
   z = (z ^ (z >> 30)) * 0xbf58476d1ce4e5b9
   z = (z ^ (z >> 27)) * 0x94d049bb133111eb
   return z ^ (z >> 31)
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   a := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   var m int
   fmt.Fscan(reader, &m)
   b := make([]int, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(reader, &b[i])
   }
   // prepare random values
   seed := uint64(time.Now().UnixNano())
   rnd := make([]uint64, n+1)
   for i := 1; i <= n; i++ {
       rnd[i] = splitmix64(seed + uint64(i))
   }
   // hash of query multiset
   var hashB uint64
   for _, v := range b {
       hashB += rnd[v]
   }
   // build P cycle
   // length T = 2*n-2 (for n>=2)
   T := 2*n - 2
   P := make([]int, T)
   idx := 0
   for i := 1; i <= n; i++ {
       if i == n {
           P[idx] = i
           idx++
           break
       }
       P[idx] = i
       idx++
   }
   for i := n - 1; i >= 2; i-- {
       P[idx] = i
       idx++
   }
   // double P to cover windows
   L := T + m - 1
   Pd := make([]int, L)
   for i := 0; i < L; i++ {
       Pd[i] = P[i%T]
   }
   // prefix hash for Pd
   H := make([]uint64, L+1)
   for i := 0; i < L; i++ {
       H[i+1] = H[i] + rnd[Pd[i]]
   }
   // prefix distance
   D := make([]int64, L)
   for i := 1; i < L; i++ {
       D[i] = D[i-1] + llabs(a[Pd[i]] - a[Pd[i-1]])
   }
   // collect distances for matching windows
   var resultVal int64
   var uniq bool
   seen := make(map[int64]struct{})
   for i := 0; i < T; i++ {
       // window [i, i+m)
       if i+m <= L && H[i+m]-H[i] == hashB {
           dist := D[i+m-1] - D[i]
           if _, ok := seen[dist]; !ok {
               seen[dist] = struct{}{}
               if !uniq && len(seen) == 1 {
                   resultVal = dist
               }
           }
           if len(seen) > 1 {
               fmt.Println(-1)
               return
           }
       }
   }
   // output
   if len(seen) == 1 {
       fmt.Println(resultVal)
   } else {
       fmt.Println(-1)
   }
}

func llabs(x int64) int64 {
   if x < 0 {
       return -x
   }
   return x
}

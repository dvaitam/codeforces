package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, k int
   fmt.Fscan(reader, &n, &k)
   p := make([]int, n+2)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &p[i])
   }
   // Precompute suffix arrays for shooters at positions >= j
   sufPos := make([]bool, n+3)
   sufSafe := make([]bool, n+3)
   sufPos[n+1] = false
   sufSafe[n+1] = true
   for i := n; i >= 1; i-- {
       sufPos[i] = sufPos[i+1] || (p[i] > 0)
       sufSafe[i] = sufSafe[i+1] && (p[i] < 100)
   }
   // visited[i][d] where i in [0..n+1], d=j-i in [0..k+1]
   stride := k + 2
   size := (n + 2) * stride
   visited := make([]byte, size)
   encode := func(i, d int) int32 {
       return int32(i<<12 | d)
   }
   decode := func(code int32) (i, d int) {
       i = int(code >> 12)
       d = int(code & ((1<<12)-1))
       return
   }
   // initial state: i=1, j=min(2,n+1)
   j0 := 2
   if j0 > n+1 {
       j0 = n + 1
   }
   i0 := 1
   d0 := j0 - i0
   start := encode(i0, d0)
   curr := make([]int32, 0, 1024)
   next := make([]int32, 0, 1024)
   curr = append(curr, start)
   visited[i0*stride+d0] = 1
   count := 1
   for round := 0; round < k && len(curr) > 0; round++ {
       next = next[:0]
       for _, code := range curr {
           i, d := decode(code)
           j := i + d
           if i > n || j > n {
               continue // empty or single, no further changes
           }
           kill1 := sufPos[j]
           surv1 := sufSafe[j]
           kill2 := p[i] > 0
           surv2 := p[i] < 100
           // both die
           if kill1 && kill2 {
               ni := j + 1
               nd := 1
               idx := ni*stride + nd
               if visited[idx] == 0 {
                   visited[idx] = 1
                   next = append(next, encode(ni, nd))
                   count++
               }
           }
           // only s1 dies
           if kill1 && surv2 {
               ni := j
               nd := 1
               idx := ni*stride + nd
               if visited[idx] == 0 {
                   visited[idx] = 1
                   next = append(next, encode(ni, nd))
                   count++
               }
           }
           // only s2 dies
           if surv1 && kill2 {
               ni := i
               nd := d + 1
               idx := ni*stride + nd
               if visited[idx] == 0 {
                   visited[idx] = 1
                   next = append(next, encode(ni, nd))
                   count++
               }
           }
           // neither die => same state, skip
       }
       curr, next = next, curr
   }
   // Output the count of distinct situations
   fmt.Println(count)
}

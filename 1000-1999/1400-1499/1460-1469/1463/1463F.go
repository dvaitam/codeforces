package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n uint64
   var x, y int
   if _, err := fmt.Fscan(reader, &n, &x, &y); err != nil {
       return
   }
   // Greedy simulation parameters
   maxXY := x
   if y > maxXY {
       maxXY = y
   }
   M := 5000
   size := M + maxXY + 5
   forbidden := make([]bool, size)
   picks := make([]bool, M+2)
   // Greedy simulate
   for i := 1; i <= M; i++ {
       if !forbidden[i] {
           picks[i] = true
           if i+x < size {
               forbidden[i+x] = true
           }
           if i+y < size {
               forbidden[i+y] = true
           }
       }
   }
   start := maxXY + 1
   // find period P
   P := 0
   for p := 1; p <= M-start; p++ {
       ok := true
       for i := start; i <= M-p; i++ {
           if picks[i] != picks[i+p] {
               ok = false
               break
           }
       }
       if ok {
           P = p
           break
       }
   }
   if P == 0 {
       P = M - start + 1
   }
   // prefix and period counts
   prefixLen := start - 1
   var prefixCnt uint64
   for i := 1; i <= prefixLen && i <= M; i++ {
       if picks[i] {
           prefixCnt++
       }
   }
   var periodCnt uint64
   for i := start; i < start+P && i <= M; i++ {
       if picks[i] {
           periodCnt++
       }
   }
   // compute result
   var res uint64
   if n <= uint64(M) {
       for i := 1; uint64(i) <= n; i++ {
           if picks[i] {
               res++
           }
       }
   } else {
       res = prefixCnt
       remN := n - uint64(prefixLen)
       res += (remN / uint64(P)) * periodCnt
       rem := remN % uint64(P)
       for i := 0; uint64(i) < rem; i++ {
           if picks[start+i] {
               res++
           }
       }
   }
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprint(writer, res)
}

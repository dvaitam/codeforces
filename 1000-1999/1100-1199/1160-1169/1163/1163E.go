package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var N int
   if _, err := fmt.Fscan(in, &N); err != nil {
       return
   }
   const maxLog = 18
   size := 1 << maxLog
   buk := make([]int, size)
   for i := 0; i < N; i++ {
       var x int
       fmt.Fscan(in, &x)
       if x >= 0 && x < size {
           buk[x] = 1
       }
   }
   B := make([]int, maxLog)
   D := make([]int, maxLog)
   var C int

   // Insert x into basis
   ins := func(x0 int) {
       x := x0
       for i := maxLog - 1; i >= 0; i-- {
           if (x>>i)&1 == 0 {
               continue
           }
           if B[i] == 0 {
               B[i] = x
               D[i] = x0
               C++
               break
           }
           x ^= B[i]
       }
   }

   // build sequence A of size 2^j
   var solveSeq func(A []int, j int)
   solveSeq = func(A []int, j int) {
       if j == 0 {
           return
       }
       j--
       step := 1 << j
       A[step] = D[j]
       solveSeq(A[:step], j)
       solveSeq(A[step:], j)
   }

   // try from maxLog down to 1
   for i := maxLog; i >= 1; i-- {
       // reset basis
       for k := 0; k < i; k++ {
           B[k] = 0
           D[k] = 0
       }
       C = 0
       lim := 1 << i
       for k := 0; k < lim; k++ {
           if buk[k] != 0 {
               ins(k)
           }
       }
       if C == i {
           // found
           fmt.Fprintln(out, i)
           A := make([]int, lim)
           solveSeq(A, i)
           x := 0
           for k := 0; k < lim; k++ {
               x ^= A[k]
               fmt.Fprintf(out, "%d ", x)
           }
           out.WriteByte('\n')
           return
       }
   }
   // none found
   fmt.Fprintln(out, 0)
   fmt.Fprintln(out, "0 ")
}

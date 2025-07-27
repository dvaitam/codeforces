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

   var N, Q int
   fmt.Fscan(in, &N, &Q)
   A := make([]int, N)
   for i := 0; i < N; i++ {
       fmt.Fscan(in, &A[i])
   }
   for qi := 0; qi < Q; qi++ {
       var t int
       fmt.Fscan(in, &t)
       if t == 1 {
           var L, R, X int
           fmt.Fscan(in, &L, &R, &X)
           for i := L - 1; i < R; i++ {
               A[i] = X
           }
       } else if t == 2 {
           var K int
           fmt.Fscan(in, &K)
           // sliding window distinct count
           cnt := make(map[int]int)
           distinct := 0
           B := 0
           // initial window
           for i := 0; i < K && i < N; i++ {
               v := A[i]
               if cnt[v] == 0 {
                   distinct++
               }
               cnt[v]++
           }
           B += distinct
           for i := K; i < N; i++ {
               // remove A[i-K]
               v := A[i-K]
               cnt[v]--
               if cnt[v] == 0 {
                   distinct--
               }
               // add A[i]
               v2 := A[i]
               if cnt[v2] == 0 {
                   distinct++
               }
               cnt[v2]++
               B += distinct
           }
           fmt.Fprintln(out, B)
       }
   }
}

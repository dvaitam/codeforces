package main

import (
   "bufio"
   "fmt"
   "os"
)

func nextPerm(a []int) bool {
   n := len(a)
   // find largest i such that a[i] < a[i+1]
   i := n - 2
   for i >= 0 && a[i] >= a[i+1] {
       i--
   }
   if i < 0 {
       return false
   }
   // find largest j > i such that a[j] > a[i]
   j := n - 1
   for j > i && a[j] <= a[i] {
       j--
   }
   // swap
   a[i], a[j] = a[j], a[i]
   // reverse a[i+1:]
   for l, r := i+1, n-1; l < r; l, r = l+1, r-1 {
       a[l], a[r] = a[r], a[l]
   }
   return true
}

func computeF(a []int) int64 {
   var sum int64
   n := len(a)
   for i := 0; i < n; i++ {
       minv := a[i]
       for j := i; j < n; j++ {
           if a[j] < minv {
               minv = a[j]
           }
           sum += int64(minv)
       }
   }
   return sum
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   var m int64
   if _, err := fmt.Fscan(in, &n, &m); err != nil {
       return
   }
   // initial permutation
   perm := make([]int, n)
   for i := 0; i < n; i++ {
       perm[i] = i + 1
   }
   // first pass: find max f
   var maxF int64
   first := true
   tmp := make([]int, n)
   copy(tmp, perm)
   for {
       f := computeF(tmp)
       if first || f > maxF {
           maxF = f
           first = false
       }
       if !nextPerm(tmp) {
           break
       }
   }
   // second pass: find m-th perm with f == maxF
   copy(tmp, perm)
   var cnt int64
   for {
       if computeF(tmp) == maxF {
           cnt++
           if cnt == m {
               // output this perm
               out := bufio.NewWriter(os.Stdout)
               for i, v := range tmp {
                   if i > 0 {
                       out.WriteByte(' ')
                   }
                   fmt.Fprintf(out, "%d", v)
               }
               out.WriteByte('\n')
               out.Flush()
               return
           }
       }
       if !nextPerm(tmp) {
           break
       }
   }
}

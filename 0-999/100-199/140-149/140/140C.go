package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// Try checks if it's possible to form q triples with distinct elements
func Try(A []int, q int) bool {
   n := len(A)
   if n < 3*q {
       return false
   }
   k := 0
   if A[q-1] == A[n-q] {
       return false
   }
   if A[q-1] == A[q] {
       j := 0
       for ; j < q; j++ {
           if A[j] == A[q] {
               break
           }
       }
       i := q
       for ; i < n; i++ {
           if A[i] != A[q] {
               break
           }
       }
       i -= q
       if i > j {
           k += i - j
       }
   }
   if n-k < 3*q {
       return false
   }
   if A[n-q-1] == A[n-q] {
       j := q
       for ; j < n-q; j++ {
           if A[j] == A[n-q] {
               break
           }
       }
       j -= k + q
       i := n - q
       for ; i < n; i++ {
           if A[i] != A[n-q] {
               break
           }
       }
       i -= (n - q)
       if i > j {
           return false
       }
   }
   return true
}

// Out prints the result triples
func Out(A []int, q int) {
   n := len(A)
   i := 0
   j := q
   k := n - q
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprintln(writer, q)
   for cnt := q; cnt > 0; cnt-- {
       for A[i] == A[j] {
           j++
       }
       fmt.Fprintf(writer, "%d %d %d\n", A[k], A[j], A[i])
       i++
       j++
       k++
   }
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(reader, &n)
   A := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &A[i])
   }
   sort.Ints(A)
   l, r := 0, n
   for r > l+1 {
       m := (l + r) / 2
       if Try(A, m) {
           l = m
       } else {
           r = m
       }
   }
   Out(A, l)
}

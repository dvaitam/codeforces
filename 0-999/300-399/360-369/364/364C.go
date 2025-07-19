package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
   "strconv"
)

var primes = []int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47}

func test(n int) []int {
   t := 2 * n * n
   A := make([]int, n)
   AA := make([]int, n)
   A[0] = 1
   m := 1
   for k := 0; m < n; k++ {
       if k >= len(primes) {
           break
       }
       p := primes[k]
       mm := 0
       d := p
       l := m
       cnt := 0
       for mm < n && d <= t {
           var i int
           for i = 0; i < l && mm < n; i++ {
               if d*A[i] <= t {
                   AA[mm] = A[i] * d
                   mm++
                   if d > 1 {
                       cnt++
                   }
               } else {
                   break
               }
           }
           if d == p {
               // record limit of base elements for next rounds
               l = i
           }
           // update d: sequence p, then 1, then p^2, p^3, ...
           if d == p {
               d = 1
           } else if d == 1 {
               d = p * p
           } else {
               d *= p
           }
       }
       // fill from A if needed
       for cnt >= (mm+2)/2 && l < m && mm < n {
           AA[mm] = A[l]
           mm++
           l++
       }
       if mm == m {
           break
       }
       m = mm
       // copy back
       for i := 0; i < m; i++ {
           A[i] = AA[i]
       }
       sort.Ints(A[:m])
   }
   return A[:n]
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   for {
       var n int
       if _, err := fmt.Fscan(reader, &n); err != nil {
           break
       }
       seq := test(n)
       for i, v := range seq {
           if i > 0 {
               writer.WriteByte(' ')
           }
           writer.WriteString(strconv.Itoa(v))
       }
       writer.WriteByte('\n')
   }
}

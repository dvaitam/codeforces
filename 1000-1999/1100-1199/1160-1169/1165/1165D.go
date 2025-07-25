package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func countDivisors(x int64) int {
   cnt := 1
   n := x
   for i := int64(2); i*i <= n; i++ {
       if n%i == 0 {
           e := 0
           for n%i == 0 {
               n /= i
               e++
           }
           cnt *= (e + 1)
       }
   }
   if n > 1 {
       cnt *= 2
   }
   return cnt
}

func main() {
   rd := bufio.NewReader(os.Stdin)
   wr := bufio.NewWriter(os.Stdout)
   defer wr.Flush()

   var t int
   if _, err := fmt.Fscan(rd, &t); err != nil {
       return
   }
   for ; t > 0; t-- {
       var n int
       fmt.Fscan(rd, &n)
       d := make([]int64, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(rd, &d[i])
       }
       sort.Slice(d, func(i, j int) bool { return d[i] < d[j] })
       if n == 0 {
           fmt.Fprintln(wr, -1)
           continue
       }
       candidate := d[0] * d[n-1]
       ok := true
       for i := 0; i < n; i++ {
           if d[i]*d[n-1-i] != candidate {
               ok = false
               break
           }
       }
       if !ok {
           fmt.Fprintln(wr, -1)
           continue
       }
       totalDiv := countDivisors(candidate)
       // subtract 2 for 1 and x itself
       if totalDiv-2 != n {
           fmt.Fprintln(wr, -1)
           continue
       }
       fmt.Fprintln(wr, candidate)
   }
}

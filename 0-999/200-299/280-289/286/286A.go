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

   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   // A lucky permutation exists only if floor(n/2) is even, i.e., n%4==0 or n%4==1
   if n%4 == 2 || n%4 == 3 {
       fmt.Fprintln(out, -1)
       return
   }
   p := make([]int, n+1)
   half := n / 2
   // pair the 2-cycles of reverse(n) in groups of two to form 4-cycles
   for i := 1; i <= half; i += 2 {
       a := i
       b := i + 1
       c := n - a + 1
       d := n - b + 1
       p[a] = b
       p[b] = c
       p[c] = d
       p[d] = a
   }
   // if n is odd, fix the middle element
   if n%2 == 1 {
       m := (n + 1) / 2
       p[m] = m
   }
   // output permutation p[1..n]
   for i := 1; i <= n; i++ {
       if i > 1 {
           out.WriteByte(' ')
       }
       fmt.Fprint(out, p[i])
   }
   out.WriteByte('\n')
}

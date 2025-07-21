package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(reader, &n)
   t := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &t[i])
   }
   // collect divisors of n (possible polygon vertex counts)
   var ks []int
   for i := 1; i*i <= n; i++ {
       if n%i == 0 {
           ks = append(ks, i)
           if i != n/i {
               ks = append(ks, n/i)
           }
       }
   }
   var ans int64 = -1e18
   // for each valid k (number of sculptures to keep)
   for _, k := range ks {
       if k < 3 {
           continue
       }
       s := n / k // spacing between kept sculptures
       // try each starting offset
       for o := 0; o < s; o++ {
           var sum int64
           for j := 0; j < k; j++ {
               sum += t[o+j*s]
           }
           if sum > ans {
               ans = sum
           }
       }
   }
   fmt.Println(ans)
}

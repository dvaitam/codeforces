package main

import (
   "bufio"
   "fmt"
   "os"
)

const mod = 1000000007

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, x, pos int
   fmt.Fscan(reader, &n, &x, &pos)

   // simulate binary search steps to count required smaller and greater placements
   l, r := 0, n
   smaller, bigger := 0, 0
   for l < r {
       mid := (l + r) / 2
       if mid <= pos {
           if mid < pos {
               smaller++
           }
           l = mid + 1
       } else {
           bigger++
           r = mid
       }
   }

   // numbers less than x: x-1, greater than x: n-x
   less := x - 1
   greater := n - x
   // if not enough numbers to assign, answer is 0
   if less < smaller || greater < bigger {
       fmt.Fprintln(writer, 0)
       return
   }

   // compute result
   result := 1
   // choose and arrange smaller numbers
   for i := 0; i < smaller; i++ {
       result = result * (less - i) % mod
   }
   // choose and arrange greater numbers
   for i := 0; i < bigger; i++ {
       result = result * (greater - i) % mod
   }
   // arrange remaining numbers freely
   remain := n - 1 - smaller - bigger
   for i := 1; i <= remain; i++ {
       result = result * i % mod
   }
   fmt.Fprintln(writer, result)
}

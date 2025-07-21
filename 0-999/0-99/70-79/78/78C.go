package main

import (
   "bufio"
   "fmt"
   "os"
)

var (
   k int64
   memo map[int64]int
)

// grundy returns the Grundy number for a pile of size x
func grundy(x int64) int {
   if x < 2*k {
       return 0
   }
   if v, ok := memo[x]; ok {
       return v
   }
   // enumerate divisors d of x such that d >= 2 and x/d >= k
   // if any move leads to Grundy 0, then this is winning (Grundy=1)
   var win bool
   // we can stop early if 0 reachable
   lim := x/k
   for i := int64(1); i*i <= x; i++ {
       if x%i != 0 {
           continue
       }
       // divisor i
       d1 := i
       if d1 >= 2 && d1 <= lim {
           // parts size = x/d1
           if d1%2 == 0 {
               win = true
           } else if grundy(x/d1) == 0 {
               win = true
           }
       }
       if win {
           memo[x] = 1
           return 1
       }
       // paired divisor
       d2 := x / i
       if d2 != d1 && d2 >= 2 && d2 <= lim {
           if d2%2 == 0 {
               win = true
           } else if grundy(x/d2) == 0 {
               win = true
           }
       }
       if win {
           memo[x] = 1
           return 1
       }
   }
   memo[x] = 0
   return 0
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int64
   if _, err := fmt.Fscan(reader, &n, &m, &k); err != nil {
       return
   }
   // if even number of piles, xor cancels out
   if n%2 == 0 {
       fmt.Println("Marsel")
       return
   }
   memo = make(map[int64]int)
   // compute Grundy of a single pile of size m
   g := grundy(m)
   if g != 0 {
       fmt.Println("Timur")
   } else {
       fmt.Println("Marsel")
   }
}

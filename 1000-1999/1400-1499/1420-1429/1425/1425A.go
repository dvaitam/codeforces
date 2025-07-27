package main

import (
   "bufio"
   "fmt"
   "os"
)

// f computes the maximum number of coins the first player can obtain
// in the game where players alternately take either one coin or half the coins (if even).
func f(n uint64) uint64 {
   // Base cases for small n
   if n <= 5 {
       switch n {
       case 0:
           return 0
       case 1, 2:
           return 1
       case 3:
           return 2
       case 4:
           return 3
       case 5:
           return 2
       }
   }
   if n&1 == 1 {
       // n is odd
       if n&3 == 1 {
           // n mod 4 == 1
           return 2 + f((n-3)/2)
       }
       // n mod 4 == 3
       return 1 + f((n-1)/2)
   }
   // n is even
   if n&3 == 2 {
       // n mod 4 == 2
       return n - f(n/2)
   }
   // n mod 4 == 0
   return n - 1 - f(n/2-1)
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var T int
   if _, err := fmt.Fscan(reader, &T); err != nil {
       return
   }
   for i := 0; i < T; i++ {
       var n uint64
       fmt.Fscan(reader, &n)
       fmt.Fprintln(writer, f(n))
   }
}

package main

import (
   "bufio"
   "fmt"
   "os"
)

func min(a, b int64) int64 {
   if a < b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int64
   var a [3]int64
   var b [3]int64
   fmt.Fscan(reader, &n)
   fmt.Fscan(reader, &a[0], &a[1], &a[2])
   fmt.Fscan(reader, &b[0], &b[1], &b[2])
   // maximum Alice wins
   // Alice rock (0) beats Bob scissors (1)
   // Alice scissors (1) beats Bob paper (2)
   // Alice paper (2) beats Bob rock (0)
   maxWins := min(a[0], b[1]) + min(a[1], b[2]) + min(a[2], b[0])

   // minimum Alice wins: Bob maximizes his wins and draws to minimize Alice wins
   // copy counts
   aa := a
   bb := b
   // Bob wins: bb[i] vs aa[(i+1)%3]
   for i := 0; i < 3; i++ {
       j := (i + 1) % 3
       t := min(bb[i], aa[j])
       bb[i] -= t
       aa[j] -= t
   }
   // Bob draws: bb[i] vs aa[i]
   for i := 0; i < 3; i++ {
       t := min(bb[i], aa[i])
       bb[i] -= t
       aa[i] -= t
   }
   // Remaining aa[*] matched with bb[*] => Alice wins
   var minWins int64
   for i := 0; i < 3; i++ {
       minWins += aa[i]
   }
   fmt.Printf("%d %d\n", minWins, maxWins)
}

package main

import (
   "fmt"
   "os"
)

func toggle(b byte) byte {
   if b == 'T' {
       return 'B'
   }
   return 'T'
}

func main() {
   var n, r int
   if _, err := fmt.Fscan(os.Stdin, &n, &r); err != nil {
       return
   }
   if n == 1 {
       if r == 1 {
           fmt.Println(0)
           fmt.Println("T")
           return
       }
       fmt.Println("IMPOSSIBLE")
       return
   }
   if n == 2 {
       if r == 2 {
           fmt.Println(0)
           fmt.Println("TB")
           return
       }
       fmt.Println("IMPOSSIBLE")
       return
   }
   bi := -1
   const INF = 1<<60
   bval := INF
   for i := 1; i+i <= r; i++ {
       x := i
       y := r - i
       ans := 0
       errc := 0
       for x > 0 && y > 0 {
           if x < y {
               x, y = y, x
           }
           dd := x / y
           ans += dd
           errc += dd - 1
           x %= y
       }
       x += y
       errc--
       if x == 1 && ans+1 == n {
           if errc < bval {
               bval = errc
               bi = i
           }
       }
   }
   if bi == -1 {
       fmt.Println("IMPOSSIBLE")
       return
   }
   u := bi
   d := r - bi
   moves := make([]byte, 0, n+2)
   for u > 1 || d > 1 {
       if u > d {
           moves = append(moves, 'T')
           u -= d
       } else {
           moves = append(moves, 'B')
           d -= u
       }
   }
   // reverse moves
   for i, j := 0, len(moves)-1; i < j; i, j = i+1, j-1 {
       moves[i], moves[j] = moves[j], moves[i]
   }
   if len(moves) == 0 {
       fmt.Println("IMPOSSIBLE")
       return
   }
   first := toggle(moves[0])
   last := toggle(moves[len(moves)-1])
   seq := make([]byte, 0, len(moves)+2)
   seq = append(seq, first)
   seq = append(seq, moves...)
   seq = append(seq, last)
   if seq[0] == 'B' {
       for i := range seq {
           seq[i] = toggle(seq[i])
       }
   }
   fmt.Println(bval)
   fmt.Println(string(seq))
}

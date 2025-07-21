package main

import (
   "bufio"
   "fmt"
   "math/big"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var s string
   if _, err := fmt.Fscan(reader, &s); err != nil {
       return
   }
   n := len(s)
   // map A->0, B->1
   sp := make([]int, n)
   for i, c := range s {
       if c == 'A' {
           sp[i] = 0
       } else {
           sp[i] = 1
       }
   }
   total := big.NewInt(0)
   // dp with states x[i-1], x[i]
   for x0 := 0; x0 <= 1; x0++ {
       for x1 := 0; x1 <= 1; x1++ {
           // dp at pos 1: prevDP[x0][x1]
           prevDP := [2][2]*big.Int{}
           // initialize
           for a := 0; a < 2; a++ {
               for b := 0; b < 2; b++ {
                   prevDP[a][b] = big.NewInt(0)
               }
           }
           prevDP[x0][x1].SetInt64(1)
           // DP from pos=1 to n-2
           for pos := 1; pos <= n-2; pos++ {
               currDP := [2][2]*big.Int{}
               for a := 0; a < 2; a++ {
                   for b := 0; b < 2; b++ {
                       currDP[a][b] = big.NewInt(0)
                   }
               }
               for prev := 0; prev < 2; prev++ {
                   for curr := 0; curr < 2; curr++ {
                       cnt := prevDP[prev][curr]
                       if cnt.Sign() == 0 {
                           continue
                       }
                       // try next bit
                       for next := 0; next < 2; next++ {
                           // t_prev: flip from (prev,curr)
                           tPrev := 0
                           if prev == 0 && curr == 1 {
                               tPrev = 1
                           }
                           // t_cur: flip from (curr,next)
                           tCur := 0
                           if curr == 0 && next == 1 {
                               tCur = 1
                           }
                           // check constraint at pos
                           if curr^tPrev^tCur != sp[pos] {
                               continue
                           }
                           // add to dp at pos+1 for state (curr,next)
                           currDP[curr][next].Add(currDP[curr][next], cnt)
                       }
                   }
               }
               prevDP = currDP
           }
           // wrap around: check positions n-1 and 0
           for prev := 0; prev < 2; prev++ {
               for curr := 0; curr < 2; curr++ {
                   cnt := prevDP[prev][curr]
                   if cnt.Sign() == 0 {
                       continue
                   }
                   // compute flips
                   tPrev := 0
                   if prev == 0 && curr == 1 {
                       tPrev = 1
                   }
                   tLast := 0
                   if curr == 0 && x0 == 1 {
                       tLast = 1
                   }
                   t0 := 0
                   if x0 == 0 && x1 == 1 {
                       t0 = 1
                   }
                   // check at n-1
                   if curr^tPrev^tLast != sp[n-1] {
                       continue
                   }
                   // check at 0
                   if x0^tLast^t0 != sp[0] {
                       continue
                   }
                   total.Add(total, cnt)
               }
           }
       }
   }
   // output result
   fmt.Println(total.String())
}

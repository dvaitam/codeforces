package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(in, &n, &m); err != nil {
       return
   }
   N := n + 1
   low := m + 1
   high := n - m
   if low > high {
       fmt.Println(0)
       return
   }
   var sumP, sumQ, sumR int64
   // angle at P obtuse: sum over a,b count c
   for a := low; a <= high; a++ {
       for b := low; b <= high; b++ {
           C := int64(2*N-2*a-b) * int64(N-2*a) + int64(3*b*N)
           denom := int64(2*N-2*a+5*b)
           if denom <= 0 {
               continue
           }
           // dot_P < 0 => c > C/denom
           threshold := C / denom
           kmin := threshold + 1
           if int(kmin) < low {
               kmin = int64(low)
           }
           if int(kmin) <= high {
               sumP += int64(high) - kmin + 1
           }
       }
   }
   // angle at Q obtuse: sum over b,c count a
   for b := low; b <= high; b++ {
       for c := low; c <= high; c++ {
           C := int64(2*N-2*b-c) * int64(N-2*b) + int64(3*c*N)
           denom := int64(2*N-2*b+5*c)
           if denom <= 0 {
               continue
           }
           threshold := C / denom
           amin := threshold + 1
           if int(amin) < low {
               amin = int64(low)
           }
           if int(amin) <= high {
               sumQ += int64(high) - amin + 1
           }
       }
   }
   // angle at R obtuse: sum over c,a count b
   for c := low; c <= high; c++ {
       for a := low; a <= high; a++ {
           C := int64(2*N-2*c-a) * int64(N-2*c) + int64(3*a*N)
           denom := int64(2*N-2*c+5*a)
           if denom <= 0 {
               continue
           }
           threshold := C / denom
           bmin := threshold + 1
           if int(bmin) < low {
               bmin = int64(low)
           }
           if int(bmin) <= high {
               sumR += int64(high) - bmin + 1
           }
       }
   }
   fmt.Println(sumP + sumQ + sumR)
}

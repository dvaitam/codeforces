package main

import (
   "bufio"
   "fmt"
   "os"
)

func gcd(a, b int64) int64 {
   for b != 0 {
       a, b = b, a%b
   }
   return a
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var a, b int64
   if _, err := fmt.Fscan(reader, &a, &b); err != nil {
       return
   }
   g := gcd(a, b)
   L := (a / g) * b
   nextA, nextB := a, b
   prev := int64(0)
   var sumA, sumB int64
   for prev < L {
       // next arrival time
       t := nextA
       if nextB < t {
           t = nextB
       }
       delta := t - prev
       if nextA < nextB {
           sumA += delta
       } else if nextB < nextA {
           sumB += delta
       } else {
           // tie: choose direction with less frequent trains (larger period)
           if a > b {
               sumA += delta
           } else {
               sumB += delta
           }
       }
       prev = t
       if nextA == t {
           nextA += a
       }
       if nextB == t {
           nextB += b
       }
   }
   // Compare total measures
   if sumA > sumB {
       fmt.Println("Dasha")
   } else if sumB > sumA {
       fmt.Println("Masha")
   } else {
       fmt.Println("Equal")
   }
}

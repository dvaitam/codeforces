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
   var n int64
   fmt.Fscan(reader, &n)
   var A, B string
   fmt.Fscan(reader, &A)
   fmt.Fscan(reader, &B)
   m := int64(len(A))
   k := int64(len(B))
   g := gcd(m, k)
   l := m / g * k
   var winA, winB int64
   for i := int64(0); i < l; i++ {
       a := A[i% m]
       b := B[i% k]
       if (a == 'R' && b == 'S') || (a == 'S' && b == 'P') || (a == 'P' && b == 'R') {
           winA++
       } else if (b == 'R' && a == 'S') || (b == 'S' && a == 'P') || (b == 'P' && a == 'R') {
           winB++
       }
   }
   full := n / l
   rem := n % l
   totalA := winA * full
   totalB := winB * full
   for i := int64(0); i < rem; i++ {
       a := A[i% m]
       b := B[i% k]
       if (a == 'R' && b == 'S') || (a == 'S' && b == 'P') || (a == 'P' && b == 'R') {
           totalA++
       } else if (b == 'R' && a == 'S') || (b == 'S' && a == 'P') || (b == 'P' && a == 'R') {
           totalB++
       }
   }
   // Nikephoros loses winB times, Polycarpus loses winA times
   fmt.Printf("%d %d\n", totalB, totalA)
}

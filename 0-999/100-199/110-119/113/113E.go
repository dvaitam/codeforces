package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var H, M, K int64
   if _, err := fmt.Fscan(in, &H, &M, &K); err != nil {
       return
   }
   var x1, y1, x2, y2 int64
   fmt.Fscan(in, &x1, &y1)
   fmt.Fscan(in, &x2, &y2)
   // number of digits in H-1 and M-1
   N1 := digits(H - 1)
   N2 := digits(M - 1)
   if K > N1+N2 {
       fmt.Println(0)
       return
   }
   // precompute powers of 10
   maxP := N1 + N2
   if maxP < K {
       maxP = K
   }
   // safe to cap at 18
   if maxP > 18 {
       maxP = 18
   }
   TEN := make([]int64, maxP+1)
   TEN[0] = 1
   for i := 1; i <= maxP; i++ {
       TEN[i] = TEN[i-1] * 10
   }
   // constant change counts
   C2 := change(N2, M-1, 0)
   C1 := change(N1, H-1, 0)
   // F function
   F := func(x, y int64) int64 {
       // if K <= N2
       if K <= N2 {
           // count full minute-wrap contributions
           d := TEN[K-1]
           part := ( (M-1) / d ) * x
           part += y / d
           need := K - C2
           if need <= 0 {
               part += x
           } else {
               d2 := TEN[need-1]
               part += x / d2
           }
           return part
       }
       // K > N2
       need := K - C2
       d2 := TEN[need-1]
       return x / d2
   }
   a := F(x1, y1)
   b := F(x2, y2)
   var ans int64
   if x1 < x2 || (x1 == x2 && y1 <= y2) {
       ans = b - a
   } else {
       total := F(H-1, M-1)
       ans = total - (a - b)
       if C1 + C2 >= K {
           ans++
       }
   }
   fmt.Println(ans)
}

// digits returns the decimal digit count of a non-negative integer (at least 1)
func digits(x int64) int64 {
   if x < 0 {
       x = -x
   }
   cnt := int64(0)
   for {
       cnt++
       x /= 10
       if x == 0 {
           break
       }
   }
   return cnt
}

// change returns number of differing digits in length-N decimal rep of a and b
func change(N int64, a, b int64) int64 {
   var cnt int64
   for i := int64(0); i < N; i++ {
       if a%10 != b%10 {
           cnt++
       }
       a /= 10
       b /= 10
   }
   return cnt
}

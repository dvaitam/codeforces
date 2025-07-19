package main

import (
   "bufio"
   "fmt"
   "os"
)

var md int

func add(a, b int) int {
   a += b
   if a >= md {
       a -= md
   }
   return a
}

func sub(a, b int) int {
   a -= b
   if a < 0 {
       a += md
   }
   return a
}

func mul(a, b int) int {
   return int((int64(a) * int64(b)) % int64(md))
}

func inv(a int) int {
   a %= md
   if a < 0 {
       a += md
   }
   b := md
   u, v := 0, 1
   for a != 0 {
       t := b / a
       b -= t * a
       a, b = b, a
       u -= t * v
       u, v = v, u
   }
   // b == 1
   if u < 0 {
       u += md
   }
   return u
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, k int
   fmt.Fscan(reader, &n, &k, &md)
   if k >= 20 || n <= (1 << (k - 1)) {
       fmt.Println(0)
       return
   }
   bc := 1 << (k - 1)
   smallSize := n / bc
   bigSize := smallSize + 1
   bigCnt := n % bc
   smallCnt := bc - bigCnt

   blocks := make([]int, bc)
   for i := 0; i < n; i++ {
       blocks[i%bc]++
   }

   // precompute inv of 2 and 4
   inv2 := inv(2)
   inv4 := inv(4)

   ans := 0
   for _, b := range blocks {
       ans = add(ans, mul(mul(b, b-1), inv4))
   }

   // sum of inverses
   sumInv := make([]int, n+2)
   for i := 0; i < n; i++ {
       sumInv[i+1] = sumInv[i]
       sumInv[i+1] = add(sumInv[i+1], inv(i+1))
   }

   // second part
   currSmall, currBig := smallCnt, bigCnt
   for _, b1 := range blocks {
       if b1 == smallSize {
           currSmall--
       } else {
           currBig--
       }
       for x := 2; x <= b1; x++ {
           if currSmall > 0 {
               aux := sub(sumInv[x+smallSize], sumInv[x])
               prob := mul(x-1, aux)
               ans = add(ans, mul(currSmall, mul(prob, inv2)))
           }
           if currBig > 0 {
               aux := sub(sumInv[x+bigSize], sumInv[x])
               prob := mul(x-1, aux)
               ans = add(ans, mul(currBig, mul(prob, inv2)))
           }
       }
       if b1 == smallSize {
           currSmall++
       } else {
           currBig++
       }
   }

   fmt.Println(ans)
}

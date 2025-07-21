package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   data, _ := in.ReadString('\n')
   data = data[:len(data)-1]
   n, err := strconv.ParseInt(data, 10, 64)
   if err != nil {
       return
   }
   // max base to consider for a^2 <= n
   maxA := int64(0)
   for i := int64(1); (i+1)*(i+1) <= n; i++ {
       maxA = i + 1
   }
   if maxA*(maxA) > n {
       // ensure maxA*maxA <= n < (maxA+1)^2
       for maxA*maxA > n {
           maxA--
       }
   }
   // mark perfect powers up to maxA
   perfectSmall := make([]bool, maxA+1)
   // collect all perfect powers <= n
   perfectAll := make(map[int64]struct{})
   // generate perfect powers
   baseMax := int64(1)
   for i := int64(1); (i+1)*(i+1) <= n; i++ {
       baseMax = i + 1
   }
   for b := int64(2); b <= baseMax; b++ {
       p := b * b
       for p <= n {
           // mark small if within small range
           if p <= maxA {
               perfectSmall[p] = true
           }
           perfectAll[p] = struct{}{}
           // multiply, check overflow
           if p > n/b {
               break
           }
           p *= b
       }
   }
   var xorSum int64 = 0
   // include a = 1, chain length = 1
   xorSum ^= 1
   // small primitive bases a=2..maxA
   for a := int64(2); a <= maxA; a++ {
       if perfectSmall[a] {
           continue
       }
       // compute chain length L: max k such that a^k <= n
       // start with p=a, k=1
       var k int64 = 1
       p := a
       for p <= n/a {
           p *= a
           k++
       }
       xorSum ^= k
   }
   // large bases a > maxA up to n: chain length =1, count only primitives => total minus perfect powers > maxA
   totalLarge := n - maxA
   // count perfect powers > maxA
   cntPP := int64(0)
   for pp := range perfectAll {
       if pp > maxA {
           cntPP++
       }
   }
   // primitive large count
   primLarge := totalLarge - cntPP
   if primLarge&1 == 1 {
       xorSum ^= 1
   }
   if xorSum != 0 {
       fmt.Println("Vasya")
   } else {
       fmt.Println("Petya")
   }
}

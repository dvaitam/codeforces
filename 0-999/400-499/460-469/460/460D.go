package main

import (
   "fmt"
)

func main() {
   var l, r int64
   var k int
   if _, err := fmt.Scan(&l, &r, &k); err != nil {
       return
   }
   // Initialize with single element
   result := l
   S := []int64{l}
   // Try small sequences up to length min(4, k)
   maxK := 4
   if k < maxK {
       maxK = k
   }
   for K := 2; K <= maxK; K++ {
       for i := 0; i < 2; i++ {
           start := l + int64(i)
           if start+int64(K)-1 > r {
               continue
           }
           temp := int64(0)
           V := make([]int64, 0, K)
           for j := start; j < start+int64(K); j++ {
               V = append(V, j)
               temp ^= j
           }
           if temp < result {
               result = temp
               S = V
           }
       }
   }
   // Special case for k>=3: attempt a triple with XOR zero
   if k >= 3 {
       // Find most significant bit of l
       msb := 0
       for b := 62; b >= 0; b-- {
           if l&(int64(1)<<b) != 0 {
               msb = b
               break
           }
       }
       A := (int64(1) << (msb + 1)) | (int64(1) << msb)
       B := (int64(1) << (msb + 1)) | (l ^ (int64(1) << msb))
       if A <= r {
           result = 0
           S = []int64{l, A, B}
       }
   }
   // Output result and sequence
   fmt.Println(result)
   fmt.Println(len(S))
   for i, v := range S {
       if i+1 == len(S) {
           fmt.Printf("%d\n", v)
       } else {
           fmt.Printf("%d ", v)
       }
   }
}

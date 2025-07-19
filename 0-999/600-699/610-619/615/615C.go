package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var s1, s2 string
   if _, err := fmt.Fscan(reader, &s1, &s2); err != nil {
       return
   }
   // prepend dummy to make 1-based indexing
   b1 := append([]byte{0}, []byte(s1)...)  // length l1+1, indices 1..l1
   b2 := append([]byte{0}, []byte(s2)...)  // length l2+1, indices 1..l2
   l1 := len(s1)
   l2 := len(s2)
   // store answer pairs
   ansL := make([]int, 0, l2)
   ansR := make([]int, 0, l2)

   i := 1
   for i <= l2 {
       bestLen := 0
       bestL, bestR := 0, 0
       // try all start positions in s1
       for j := 1; j <= l1; j++ {
           if b1[j] != b2[i] {
               continue
           }
           // forward match
           k := j
           for k+1 <= l1 && (k+1-j+i) <= l2 && b1[k+1] == b2[k+1-j+i] {
               k++
           }
           len1 := k - j + 1
           if len1 > bestLen {
               bestLen = len1
               bestL = j
               bestR = k
           }
           // backward match
           k2 := j - 1
           for k2 >= 1 && (j-k2+i) <= l2 && b1[k2] == b2[j-k2+i] {
               k2--
           }
           len2 := j - k2
           if len2 > bestLen {
               bestLen = len2
               bestL = j
               bestR = k2 + 1
           }
       }
       if bestLen == 0 {
           fmt.Fprintln(writer, -1)
           return
       }
       ansL = append(ansL, bestL)
       ansR = append(ansR, bestR)
       i += bestLen
   }
   // output
   fmt.Fprintln(writer, len(ansL))
   for idx := range ansL {
       fmt.Fprintf(writer, "%d %d\n", ansL[idx], ansR[idx])
   }
}

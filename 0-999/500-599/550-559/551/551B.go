package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var a, b, c string
   if _, err := fmt.Fscan(reader, &a); err != nil {
       return
   }
   fmt.Fscan(reader, &b)
   fmt.Fscan(reader, &c)
   // count frequencies
   var freqA, freqB, freqC [26]int
   for i := 0; i < len(a); i++ {
       freqA[a[i]-'a']++
   }
   for i := 0; i < len(b); i++ {
       freqB[b[i]-'a']++
   }
   for i := 0; i < len(c); i++ {
       freqC[c[i]-'a']++
   }
   // maximum times we can use b alone
   maxB := len(a) / len(b)
   for i := 0; i < 26; i++ {
       if freqB[i] > 0 {
           if v := freqA[i] / freqB[i]; v < maxB {
               maxB = v
           }
       }
   }
   bestX, bestY, bestSum := 0, 0, 0
   // try x copies of b
   for x := 0; x <= maxB; x++ {
       // remaining after x*b
       remOK := true
       var rem [26]int
       for i := 0; i < 26; i++ {
           rem[i] = freqA[i] - x*freqB[i]
           if rem[i] < 0 {
               remOK = false
               break
           }
       }
       if !remOK {
           break
       }
       // compute max y for c
       y := len(a) / len(c)
       for i := 0; i < 26; i++ {
           if freqC[i] > 0 {
               if v := rem[i] / freqC[i]; v < y {
                   y = v
               }
           }
       }
       if x+y > bestSum {
           bestSum = x + y
           bestX = x
           bestY = y
       }
   }
   // build result
   res := make([]byte, 0, len(a))
   bBytes := []byte(b)
   cBytes := []byte(c)
   // append b bestX times, c bestY times
   for i := 0; i < bestX; i++ {
       res = append(res, bBytes...)
   }
   for i := 0; i < bestY; i++ {
       res = append(res, cBytes...)
   }
   // append leftover letters
   for i := 0; i < 26; i++ {
       used := bestX*freqB[i] + bestY*freqC[i]
       left := freqA[i] - used
       for j := 0; j < left; j++ {
           res = append(res, byte('a'+i))
       }
   }
   // output
   writer := bufio.NewWriter(os.Stdout)
   writer.Write(res)
   writer.WriteByte('\n')
   writer.Flush()
}

package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var b, d int
   var a, c string
   if _, err := fmt.Fscan(reader, &b, &d); err != nil {
       return
   }
   fmt.Fscan(reader, &a)
   fmt.Fscan(reader, &c)
   aLen := len(a)
   cLen := len(c)
   // Precompute matches and next positions for each c starting index
   count := make([]int, cLen)
   nextIdx := make([]int, cLen)
   for j0 := 0; j0 < cLen; j0++ {
       j := j0
       cnt := 0
       for i := 0; i < aLen; i++ {
           if a[i] == c[j] {
               j++
               if j == cLen {
                   j = 0
                   cnt++
               }
           }
       }
       count[j0] = cnt
       nextIdx[j0] = j
   }
   // Detect cycle over b repetitions
   firstIdx := make([]int, cLen)
   firstCount := make([]int64, cLen)
   curJ := 0
   var totalCount int64
   i := 1
   for i <= b {
       if firstIdx[curJ] != 0 {
           prevI := firstIdx[curJ]
           prevCount := firstCount[curJ]
           cycleLen := i - prevI
           cycleCount := totalCount - prevCount
           remaining := b - i + 1
           fullCycles := remaining / cycleLen
           totalCount += int64(fullCycles) * cycleCount
           leftover := remaining % cycleLen
           for k := 0; k < leftover; k++ {
               totalCount += int64(count[curJ])
               curJ = nextIdx[curJ]
           }
           break
       }
       firstIdx[curJ] = i
       firstCount[curJ] = totalCount
       totalCount += int64(count[curJ])
       curJ = nextIdx[curJ]
       i++
   }
   // Answer is number of full q strings: each q has d copies of c
   // totalCount is number of c strings obtained; divide by d
   fmt.Println(totalCount / int64(d))
}

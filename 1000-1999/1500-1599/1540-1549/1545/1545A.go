package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var t int
   fmt.Fscan(reader, &t)
   for t > 0 {
       t--
       var n int
       fmt.Fscan(reader, &n)
       a := make([]int, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &a[i])
       }
       b := make([]int, n)
       copy(b, a)
       sort.Ints(b)
       // count values by index parity groups for a and sorted b
       cntA := make(map[int][2]int)
       cntB := make(map[int][2]int)
       for i, v := range a {
           p := i % 2
           arr := cntA[v]
           arr[p]++
           cntA[v] = arr
       }
       for i, v := range b {
           p := i % 2
           arr := cntB[v]
           arr[p]++
           cntB[v] = arr
       }
       ok := true
       // check all values have same counts in each parity group
       for v, arrA := range cntA {
           if arrB, found := cntB[v]; !found || arrA != arrB {
               ok = false
               break
           }
       }
       // also ensure no extra in cntB
       if ok {
           for v, arrB := range cntB {
               if arrA, found := cntA[v]; !found || arrA != arrB {
                   ok = false
                   break
               }
           }
       }
       if ok {
           fmt.Fprintln(writer, "YES")
       } else {
           fmt.Fprintln(writer, "NO")
       }
   }
}

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

   var m int
   if _, err := fmt.Fscan(reader, &m); err != nil {
       return
   }
   tp := make([]int, m)
   x := make([]int64, m)
   l := make([]int64, m)
   c := make([]int64, m)
   lenArr := make([]int64, m)
   var curLen int64
   for i := 0; i < m; i++ {
       fmt.Fscan(reader, &tp[i])
       if tp[i] == 1 {
           fmt.Fscan(reader, &x[i])
           curLen++
       } else {
           fmt.Fscan(reader, &l[i], &c[i])
           curLen += l[i] * c[i]
       }
       lenArr[i] = curLen
   }
   var q int
   fmt.Fscan(reader, &q)
   queries := make([]int64, q)
   for i := 0; i < q; i++ {
       fmt.Fscan(reader, &queries[i])
   }
   for i, pos := range queries {
       ans := findValue(pos, tp, x, l, lenArr)
       if i > 0 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, ans)
   }
   writer.WriteByte('\n')
}

// findValue returns the element at position pos in the built sequence
func findValue(pos int64, tp []int, x, l, lenArr []int64) int64 {
   for {
       idx := sort.Search(len(lenArr), func(i int) bool {
           return lenArr[i] >= pos
       })
       if tp[idx] == 1 {
           return x[idx]
       }
       prevLen := int64(0)
       if idx > 0 {
           prevLen = lenArr[idx-1]
       }
       // map pos into the prefix of length l[idx]
       pos = (pos - prevLen - 1) % l[idx] + 1
   }
}

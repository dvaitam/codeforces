package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func gcd(a, b int64) int64 {
   for b != 0 {
       a, b = b, a%b
   }
   return a
}

func solve(reader *bufio.Reader, writer *bufio.Writer) {
   var n int
   fmt.Fscan(reader, &n)
   v := make([]int64, n)
   for i := range v {
       fmt.Fscan(reader, &v[i])
   }
   sort.Slice(v, func(i, j int) bool { return v[i] < v[j] })
   used := make([]bool, n)
   result := make([]int64, 0, n)
   leftovers := make([]int64, 0, n)
   currentG := v[n-1]
   used[n-1] = true
   result = append(result, currentG)
   for {
       nextMaxG := int64(1)
       for i := 0; i < n; i++ {
           if used[i] {
               continue
           }
           d := gcd(v[i], currentG)
           if d == 1 {
               leftovers = append(leftovers, v[i])
               used[i] = true
           } else if d == currentG {
               result = append(result, v[i])
               used[i] = true
           } else if d > nextMaxG {
               nextMaxG = d
           }
       }
       if nextMaxG == 1 {
           break
       }
       currentG = nextMaxG
   }
   for _, x := range leftovers {
       result = append(result, x)
   }
   for i, x := range result {
       if i > 0 {
           writer.WriteByte(' ')
       }
       writer.WriteString(fmt.Sprintf("%d", x))
   }
   writer.WriteByte('\n')
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var t int
   fmt.Fscan(reader, &t)
   for i := 0; i < t; i++ {
       solve(reader, writer)
   }
}

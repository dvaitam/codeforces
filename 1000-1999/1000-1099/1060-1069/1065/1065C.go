package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   n64 := readInt(reader)
   m := readInt(reader)
   n := int(n64)
   a := make([]int64, n)
   for i := 0; i < n; i++ {
       a[i] = readInt(reader)
   }
   sort.Slice(a, func(i, j int) bool { return a[i] > a[j] })
   ans := 0
   i := 0
   for a[0] != a[n-1] {
       t := i
       var tot int64
       // try to level down from current prefix
       for t+1 < n {
           diff := a[t] - a[t+1]
           cost := diff * int64(t+1)
           if tot+cost > m {
               break
           }
           tot += cost
           t++
       }
       // if no progress on t, break to avoid infinite loop
       if t == i {
           // we can only reduce a[i] by remaining m
           i = t
       }
       ans++
       i = t
       // reduce all first t+1 elements by as much as possible
       dec := (m - tot) / int64(i+1)
       a[i] -= dec
   }
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprint(writer, ans)
}

// readInt reads next integer from bufio.Reader
func readInt(r *bufio.Reader) int64 {
   var x int64
   var neg bool
   for {
       b, err := r.ReadByte()
       if err != nil {
           return 0
       }
       if b == '-' {
           neg = true
           break
       }
       if b >= '0' && b <= '9' {
           x = int64(b - '0')
           break
       }
   }
   for {
       b, err := r.ReadByte()
       if err != nil || b < '0' || b > '9' {
           break
       }
       x = x*10 + int64(b-'0')
   }
   if neg {
       x = -x
   }
   return x
}

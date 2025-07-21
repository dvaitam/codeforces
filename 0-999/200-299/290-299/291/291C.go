package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, k int
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   ips := make([]uint32, n)
   for i := 0; i < n; i++ {
       var a, b, c, d uint32
       fmt.Fscan(reader, &a)
       reader.ReadByte() // dot
       fmt.Fscan(reader, &b)
       reader.ReadByte()
       fmt.Fscan(reader, &c)
       reader.ReadByte()
       fmt.Fscan(reader, &d)
       ips[i] = (a<<24 | b<<16 | c<<8 | d)
   }
   sort.Slice(ips, func(i, j int) bool { return ips[i] < ips[j] })
   // try mask lengths from 1 to 31
   for t := 1; t <= 31; t++ {
       shift := 32 - t
       var cnt int
       var last uint32 = ^uint32(0)
       for i, ip := range ips {
           v := ip >> shift
           if i == 0 || v != last {
               cnt++
               last = v
           }
           if cnt > k {
               break
           }
       }
       if cnt == k {
           // build mask
           mask := ^uint32(0) << shift
           m1 := (mask >> 24) & 0xFF
           m2 := (mask >> 16) & 0xFF
           m3 := (mask >> 8) & 0xFF
           m4 := mask & 0xFF
           fmt.Printf("%d.%d.%d.%d\n", m1, m2, m3, m4)
           return
       }
   }
   fmt.Println(-1)
}

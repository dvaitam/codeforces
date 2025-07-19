package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var a00, a01, a10, a11 int64
   if _, err := fmt.Fscan(reader, &a00, &a01, &a10, &a11); err != nil {
       return
   }
   zero, one := int64(-1), int64(-1)
   maxv := a00
   if a11 > maxv {
       maxv = a11
   }
   for i := int64(1); ; i++ {
       v := i * (i - 1) / 2
       if v > maxv {
           break
       }
       if v == a00 {
           zero = i
       }
       if v == a11 {
           one = i
       }
   }
   if zero == -1 || one == -1 {
       fmt.Println("Impossible")
       return
   }
   // special cases
   if a00 == 0 && a11 == 0 && a01 == 0 && a10 == 0 {
       fmt.Println("0")
       return
   }
   if a00 == 0 && a01 == 0 && a10 == 0 {
       for i := int64(0); i < one; i++ {
           fmt.Print("1")
       }
       fmt.Println()
       return
   }
   if a11 == 0 && a01 == 0 && a10 == 0 {
       for i := int64(0); i < zero; i++ {
           fmt.Print("0")
       }
       fmt.Println()
       return
   }
   if zero*one != a01 + a10 {
       fmt.Println("Impossible")
       return
   }
   b01 := zero * one
   lol := int64(0)
   for b01 - a01 >= one {
       lol++
       b01 -= one
   }
   extra := b01 - a01
   // build result
   var sb []byte
   preZeros := zero - lol
   if extra > 0 {
       preZeros--
   }
   for i := int64(0); i < preZeros; i++ {
       sb = append(sb, '0')
   }
   for i := int64(0); i < extra; i++ {
       sb = append(sb, '1')
   }
   if extra > 0 {
       sb = append(sb, '0')
   }
   for i := int64(0); i < one-extra; i++ {
       sb = append(sb, '1')
   }
   for i := int64(0); i < lol; i++ {
       sb = append(sb, '0')
   }
   os.Stdout.Write(sb)
}

package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
)

func sumDigits(x int) int {
   s := 0
   for x > 0 {
       s += x % 10
       x /= 10
   }
   return s
}

func check(x int) bool {
   return sumDigits(sumDigits(x)) < 10
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var n int
   var str string
   fmt.Fscan(reader, &n, &str)
   ds := make([]int, n+2)
   for i := 1; i <= n; i++ {
       ds[i] = int(str[i-1] - '0')
   }
   p := make([]bool, n+2)
   x := 0
   for i := 1; i <= n; i++ {
       x += ds[i]
   }
   xx := x
   if !check(x) {
       for i := 1; i <= n; i++ {
           p[i] = false
       }
       x = xx
       for i := 1; i <= n; i++ {
           if (i&1) == 1 && i+1 <= n {
               p[i] = true
               x += ds[i] * 9
               if check(x) {
                   break
               }
           }
       }
   }
   if !check(x) {
       for i := 1; i <= n; i++ {
           p[i] = false
       }
       x = xx
       for i := 1; i <= n; i++ {
           if (i&1) == 0 && i+1 <= n {
               p[i] = true
               x += ds[i] * 9
               if check(x) {
                   break
               }
           }
       }
   }
   if !check(x) {
       for i := 1; i <= n; i++ {
           p[i] = false
       }
       x = xx
       for i := n; i >= 1; i-- {
           if (i&1) == 1 && i+1 <= n {
               p[i] = true
               x += ds[i] * 9
               if check(x) {
                   break
               }
           }
       }
   }
   if !check(x) {
       for i := 1; i <= n; i++ {
           p[i] = false
       }
       x = xx
       for i := n; i >= 1; i-- {
           if (i&1) == 0 && i+1 <= n {
               p[i] = true
               x += ds[i] * 9
               if check(x) {
                   break
               }
           }
       }
   }
   if !check(x) {
       for i := 1; i <= n; i++ {
           p[i] = false
       }
       x = xx
       for i := 1; i <= n; i++ {
           if i%3 == 1 && i+2 <= n {
               p[i], p[i+1] = true, true
               x += ds[i]*99 + ds[i+1]*9
               if check(x) {
                   break
               }
           }
       }
   }
   if !check(x) {
       for i := 1; i <= n; i++ {
           p[i] = false
       }
       x = xx
       for i := 1; i <= n; i++ {
           if i%3 == 2 && i+2 <= n {
               p[i], p[i+1] = true, true
               x += ds[i]*99 + ds[i+1]*9
               if check(x) {
                   break
               }
           }
       }
   }
   if !check(x) {
       for i := 1; i <= n; i++ {
           p[i] = false
       }
       x = xx
       for i := n; i >= 1; i-- {
           if i%3 == 1 && i+2 <= n {
               p[i], p[i+1] = true, true
               x += ds[i]*99 + ds[i+1]*9
               if check(x) {
                   break
               }
           }
       }
   }
   if !check(x) {
       for i := 1; i <= n; i++ {
           p[i] = false
       }
       x = xx
       for i := n; i >= 1; i-- {
           if i%3 == 2 && i+2 <= n {
               p[i], p[i+1] = true, true
               x += ds[i]*99 + ds[i+1]*9
               if check(x) {
                   break
               }
           }
       }
   }
   // output original expression
   for i := 1; i <= n; i++ {
       writer.WriteByte(byte(ds[i] + '0'))
       if !p[i] && i < n {
           writer.WriteByte('+')
       }
   }
   writer.WriteByte('\n')
   // first reduction
   s2 := strconv.Itoa(x)
   sum1 := 0
   m := len(s2)
   for i := 0; i < m; i++ {
       writer.WriteByte(s2[i])
       sum1 += int(s2[i] - '0')
       if i < m-1 {
           writer.WriteByte('+')
       }
   }
   writer.WriteByte('\n')
   // second reduction
   s3 := strconv.Itoa(sum1)
   m2 := len(s3)
   for i := 0; i < m2; i++ {
       writer.WriteByte(s3[i])
       if i < m2-1 {
           writer.WriteByte('+')
       }
   }
   writer.WriteByte('\n')
}

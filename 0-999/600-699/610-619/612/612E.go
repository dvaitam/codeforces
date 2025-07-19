package main

import (
   "bufio"
   "os"
)

var (
   reader = bufio.NewReader(os.Stdin)
   writer = bufio.NewWriter(os.Stdout)
)

func readInt() int {
   var c byte
   var err error
   // skip non-digits
   c, err = reader.ReadByte()
   for err == nil && (c < '0' || c > '9') {
       c, err = reader.ReadByte()
   }
   var x int
   for err == nil && c >= '0' && c <= '9' {
       x = x*10 + int(c-'0')
       c, err = reader.ReadByte()
   }
   return x
}

func writeInt(x int) {
   if x == 0 {
       writer.WriteByte('0')
       return
   }
   var buf [20]byte
   i := len(buf)
   for x > 0 {
       i--
       buf[i] = byte('0' + x%10)
       x /= 10
   }
   writer.Write(buf[i:])
}

func main() {
   defer writer.Flush()
   n := readInt()
   p := make([]int, n+1)
   ans := make([]int, n+1)
   v := make([]bool, n+1)
   tmp := make([]int, n+1)
   // cycle stores first node index for even-length cycles
   cycle := make([]int, n+1)
   // t2 may need up to 2*n entries
   t2 := make([]int, 2*n+1)

   for i := 1; i <= n; i++ {
       p[i] = readInt()
   }
   for i := 1; i <= n; i++ {
       if !v[i] {
           l := 0
           // collect cycle starting at i
           j := i
           for !v[j] {
               l++
               tmp[l] = j
               v[j] = true
               j = p[j]
           }
           if l&1 == 1 {
               // odd length: rotate by 2
               for j := 1; j <= l; j++ {
                   pos := j*2 - 1
                   if pos > l {
                       pos -= l
                   }
                   t2[pos] = tmp[j]
               }
           } else {
               if cycle[l] != 0 {
                   jidx := cycle[l]
                   r := i
                   // merge two cycles of same even length
                   for k := 1; k <= l*2; k++ {
                       if k&1 == 1 {
                           t2[k] = r
                           r = p[r]
                       } else {
                           t2[k] = jidx
                           jidx = p[jidx]
                       }
                   }
                   cycle[l] = 0
                   l *= 2
               } else {
                   cycle[l] = i
                   continue
               }
           }
           // assign next pointers
           for k := 1; k < l; k++ {
               ans[t2[k]] = t2[k+1]
           }
           ans[t2[l]] = t2[1]
       }
   }
   // check if any not assigned
   for i := 1; i <= n; i++ {
       if ans[i] == 0 {
           writer.WriteString("-1")
           writer.WriteByte('\n')
           return
       }
   }
   // output ans
   for i := 1; i <= n; i++ {
       writeInt(ans[i])
       if i < n {
           writer.WriteByte(' ')
       }
   }
}

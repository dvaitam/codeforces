package main

import (
   "bufio"
   "os"
   "strconv"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   // read int from reader
   readInt := func() int {
       b, _ := reader.ReadByte()
       // skip non-digit, non-minus
       for (b < '0' || b > '9') && b != '-' {
           b, _ = reader.ReadByte()
       }
       neg := false
       if b == '-' {
           neg = true
           b, _ = reader.ReadByte()
       }
       x := 0
       for b >= '0' && b <= '9' {
           x = x*10 + int(b-'0')
           b, _ = reader.ReadByte()
       }
       if neg {
           return -x
       }
       return x
   }

   T := readInt()
   for t := 0; t < T; t++ {
       n := readInt()
       s := make([]byte, n)
       // read n-1 chars '<' or '>'
       for i := 0; i < n-1; {
           b, err := reader.ReadByte()
           if err != nil {
               break
           }
           if b == '<' || b == '>' {
               s[i] = b
               i++
           }
       }
       // sentinel
       s[n-1] = '>'

       // first permutation
       ans1 := make([]int, n)
       now := n
       for i := 0; i < n-1; i++ {
           if s[i] == '<' && s[i+1] == '>' {
               x := i
               for x >= 0 && s[x] == '<' {
                   ans1[x+1] = now
                   now--
                   x--
               }
           }
       }
       for i := 0; i < n; i++ {
           if ans1[i] == 0 {
               ans1[i] = now
               now--
           }
       }
       // output ans1
       for i := 0; i < n; i++ {
           writer.WriteString(strconv.Itoa(ans1[i]))
           writer.WriteByte(' ')
       }
       writer.WriteByte('\n')

       // second permutation
       ans2 := make([]int, n)
       now = n
       for i := n - 2; i >= 0; i-- {
           if s[i] == '<' {
               ans2[i+1] = now
               now--
           }
       }
       for i := 0; i < n; i++ {
           if ans2[i] == 0 {
               ans2[i] = now
               now--
           }
       }
       for i := 0; i < n; i++ {
           writer.WriteString(strconv.Itoa(ans2[i]))
           writer.WriteByte(' ')
       }
       writer.WriteByte('\n')
   }
}

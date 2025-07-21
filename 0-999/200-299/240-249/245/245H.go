package main

import (
   "bufio"
   "bytes"
   "fmt"
   "os"
)

// Fast input reader
type FastReader struct {
   buf *bufio.Reader
}

func NewFastReader(r *os.File) *FastReader {
   return &FastReader{buf: bufio.NewReader(r)}
}

func (fr *FastReader) ReadBytes() []byte {
   line, _ := fr.buf.ReadBytes('\n')
   return bytes.TrimSpace(line)
}

func (fr *FastReader) ReadInt() int {
   var c byte
   var sign, x int = 1, 0
   // skip non-digit
   for {
       b, err := fr.buf.ReadByte()
       if err != nil {
           return x
       }
       c = b
       if (c >= '0' && c <= '9') || c == '-' {
           break
       }
   }
   if c == '-' {
       sign = -1
       c, _ = fr.buf.ReadByte()
   }
   for ; c >= '0' && c <= '9'; c, _ = fr.buf.ReadByte() {
       x = x*10 + int(c - '0')
   }
   return x * sign
}

func main() {
   in := NewFastReader(os.Stdin)
   s := in.ReadBytes()
   n := len(s)
   q := in.ReadInt()
   // queries grouped by right endpoint
   queriesByR := make([][][2]int, n+2)
   cnt := make([]int, n+2)
   ls := make([]int, q)
   rs := make([]int, q)
   for i := 0; i < q; i++ {
       l := in.ReadInt()
       r := in.ReadInt()
       ls[i], rs[i] = l, r
       cnt[r]++
   }
   for r := 1; r <= n; r++ {
       if cnt[r] > 0 {
           queriesByR[r] = make([][2]int, 0, cnt[r])
       }
   }
   for i := 0; i < q; i++ {
       r := rs[i]
       queriesByR[r] = append(queriesByR[r], [2]int{ls[i], i})
   }
   // endsAt: for each r, list of start positions of palindromes ending at r
   endsAt := make([][]int, n+2)
   for r := 1; r <= n; r++ {
       endsAt[r] = make([]int, 0, r)
   }
   // expand palindromes
   for center := 1; center <= n; center++ {
       // odd
       l, r1 := center, center
       for l >= 1 && r1 <= n && s[l-1] == s[r1-1] {
           endsAt[r1] = append(endsAt[r1], l)
           l--
           r1++
       }
       // even
       l, r1 = center, center+1
       for l >= 1 && r1 <= n && s[l-1] == s[r1-1] {
           endsAt[r1] = append(endsAt[r1], l)
           l--
           r1++
       }
   }
   // diff array for range updates
   diff := make([]int, n+3)
   pre := make([]int, n+2)
   ans := make([]int, q)
   // process by right endpoint
   for r := 1; r <= n; r++ {
       // apply updates for palindromes ending at r
       for _, st := range endsAt[r] {
           diff[1]++
           diff[st+1]--
       }
       // build prefix pre[i] = #pal substrings ending <=r with start >= i
       cur := 0
       for i := 1; i <= n; i++ {
           cur += diff[i]
           pre[i] = cur
       }
       // answer queries ending at r
       for _, qr := range queriesByR[r] {
           l, idx := qr[0], qr[1]
           ans[idx] = pre[l]
       }
   }
   // output
   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()
   for i, v := range ans {
       if i > 0 {
           w.WriteByte(' ')
       }
       w.WriteString(fmt.Sprint(v))
   }
}

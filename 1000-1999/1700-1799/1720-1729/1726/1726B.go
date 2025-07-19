package main

import (
   "bufio"
   "fmt"
   "os"
)

var reader = bufio.NewReader(os.Stdin)
var writer = bufio.NewWriter(os.Stdout)

// readInt reads the next integer from standard input.
func readInt() (int64, error) {
   var x int64
   var sign int64 = 1
   // read first non-space byte
   b, err := reader.ReadByte()
   if err != nil {
       return 0, err
   }
   for ; (b < '0' || b > '9') && b != '-'; b, err = reader.ReadByte() {
       if err != nil {
           return 0, err
       }
   }
   if b == '-' {
       sign = -1
       b, err = reader.ReadByte()
       if err != nil {
           return 0, err
       }
   }
   for ; b >= '0' && b <= '9'; b, err = reader.ReadByte() {
       if err != nil {
           break
       }
       x = x*10 + int64(b-'0')
   }
   return x * sign, nil
}

func main() {
   defer writer.Flush()
   t, err := readInt()
   if err != nil {
       return
   }
   for ; t > 0; t-- {
       n, _ := readInt()
       m, _ := readInt()
       // impossible if sum too small
       if m <= n-1 {
           fmt.Fprintln(writer, "No")
           continue
       }
       // all equal
       if m % n == 0 {
           fmt.Fprintln(writer, "Yes")
           v := m / n
           for i := int64(0); i < n; i++ {
               if i > 0 {
                   writer.WriteByte(' ')
               }
               fmt.Fprint(writer, v)
           }
           writer.WriteByte('\n')
           continue
       }
       // even count cannot sum to odd with all same parity
       if n % 2 == 0 && m % 2 != 0 {
           fmt.Fprintln(writer, "No")
           continue
       }
       // otherwise possible
       fmt.Fprintln(writer, "Yes")
       if n % 2 != 0 {
           // use n-1 ones, last is remainder
           for i := int64(0); i < n-1; i++ {
               fmt.Fprint(writer, "1 ")
           }
           fmt.Fprintln(writer, m-(n-1))
       } else {
           // even n: use n-2 ones, last two equal
           for i := int64(0); i < n-2; i++ {
               fmt.Fprint(writer, "1 ")
           }
           rem := m - (n - 2)
           half := rem / 2
           fmt.Fprintln(writer, half, half)
       }
   }
}

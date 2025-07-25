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

   n := readInt(reader)
   k := readInt(reader)
   g := 0
   for i := 0; i < n; i++ {
       a := readInt(reader)
       r := a % k
       g = gcd(g, r)
   }
   g = gcd(g, k)
   // collect reachable residues: multiples of g
   count := k / g
   // output count
   writer.WriteString(strconv.Itoa(count))
   writer.WriteByte('\n')
   // output residues
   for i := 0; i < k; i += g {
       if i > 0 {
           writer.WriteByte(' ')
       }
       writer.WriteString(strconv.Itoa(i))
   }
   writer.WriteByte('\n')
}

// gcd returns the greatest common divisor of a and b.
func gcd(a, b int) int {
   for b != 0 {
       a, b = b, a%b
   }
   return a
}

// readInt reads an integer from bufio.Reader (skips non-digits).
func readInt(r *bufio.Reader) int {
   c, err := r.ReadByte()
   for err == nil && (c < '0' || c > '9') {
       c, err = r.ReadByte()
   }
   if err != nil {
       return 0
   }
   x := 0
   for err == nil && c >= '0' && c <= '9' {
       x = x*10 + int(c-'0')
       c, err = r.ReadByte()
   }
   return x
}

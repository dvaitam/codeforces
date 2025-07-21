package main

import (
   "bufio"
   "fmt"
   "os"
   "strings"
)

// reverseBits8 returns the byte result of reversing the order of bits in b.
func reverseBits8(b byte) byte {
   var r byte
   for i := 0; i < 8; i++ {
       if (b>>i)&1 == 1 {
           r |= 1 << (7 - i)
       }
   }
   return r
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   // Read the entire line including spaces
   s, err := reader.ReadString('\n')
   if err != nil && err.Error() != "EOF" {
       fmt.Fprintln(os.Stderr, "error reading input:", err)
       os.Exit(1)
   }
   s = strings.TrimRight(s, "\r\n")
   // Process each character
   var prevRev byte = 0
   for i := 0; i < len(s); i++ {
       c := s[i]
       // reverse bits of current char to get y
       y := reverseBits8(c)
       // a = (prevRev - y) mod 256
       a := (int(prevRev) - int(y) + 256) % 256
       fmt.Println(a)
       // update prevRev to reverse bits of current char
       prevRev = reverseBits8(c)
   }
}

package main

import (
   "bufio"
   "os"
   "sort"
   "fmt"
)

func readInt(r *bufio.Reader) (int, error) {
   var x, sign int = 0, 1
   for {
       b, err := r.ReadByte()
       if err != nil {
           return 0, err
       }
       if b == '-' {
           sign = -1
       }
       if b >= '0' && b <= '9' {
           x = int(b - '0')
           break
       }
   }
   for {
       b, err := r.ReadByte()
       if err != nil {
           break
       }
       if b < '0' || b > '9' {
           break
       }
       x = x*10 + int(b - '0')
   }
   return x * sign, nil
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   n, err := readInt(reader)
   if err != nil {
       return
   }
   sizes := make([]int, n)
   for i := 0; i < n; i++ {
       v, err := readInt(reader)
       if err != nil {
           return
       }
       sizes[i] = v
   }
   sort.Ints(sizes)
   mid := n / 2
   i, j, count := 0, mid, 0
   for i < mid && j < n {
       if sizes[j] >= 2*sizes[i] {
           count++
           i++
           j++
       } else {
           j++
       }
   }
   // Visible kangaroos = total - held
   res := n - count
   fmt.Fprintln(writer, res)
}

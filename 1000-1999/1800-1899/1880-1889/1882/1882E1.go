package main

import (
   "bufio"
   "fmt"
   "os"
)

func readInt(r *bufio.Reader) (int, error) {
   var x int
   var sign = 1
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
       x = x*10 + int(b-'0')
   }
   return x * sign, nil
}

func solve(n int, a []int) []int {
   var ops []int
   for i := 1; i <= n; i++ {
       for j := i; j <= n; j++ {
           if a[j] == i {
               if j != i {
                   ops = append(ops, i, j-i, n-j+1)
               }
               a[i], a[j] = a[j], a[i]
               break
           }
       }
   }
   return ops
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   n, err := readInt(reader)
   if err != nil {
       return
   }
   m, _ := readInt(reader)
   a := make([]int, n+1)
   b := make([]int, m+1)
   for i := 1; i <= n; i++ {
       a[i], _ = readInt(reader)
   }
   for i := 1; i <= m; i++ {
       b[i], _ = readInt(reader)
   }
   opa := solve(n, a)
   opb := solve(m, b)

   la := len(opa)
   lb := len(opb)
   if (la+lb)&1 == 1 {
       if n&1 == 1 {
           for i := 0; i < n; i++ {
               opa = append(opa, n)
           }
       } else if m&1 == 1 {
           for i := 0; i < m; i++ {
               opb = append(opb, m)
           }
       } else {
           fmt.Fprintln(writer, -1)
           return
       }
       la = len(opa)
       lb = len(opb)
   }
   // balance lengths
   for la < lb {
       opa = append(opa, 1, n)
       la += 2
   }
   for lb < la {
       opb = append(opb, 1, m)
       lb += 2
   }
   // output
   fmt.Fprintln(writer, la)
   for i := 0; i < la; i++ {
       fmt.Fprintln(writer, opa[i], opb[i])
   }
}

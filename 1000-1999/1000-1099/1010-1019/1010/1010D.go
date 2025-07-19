package main

import (
   "bufio"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   n := readInt(reader)
   ch0 := make([]int, n+1)
   ch1 := make([]int, n+1)
   typ := make([]int, n+1) // -1: input, 0: AND,1:OR,2:XOR,3:NOT
   g := make([]int, n+1)
   for i := 1; i <= n; i++ {
       op := readOp(reader)
       switch op {
       case 'I':
           typ[i] = -1
           g[i] = readInt(reader)
       case 'N':
           typ[i] = 3
           c := readInt(reader)
           ch0[i] = c
       case 'A':
           typ[i] = 0
           c0 := readInt(reader)
           c1 := readInt(reader)
           ch0[i], ch1[i] = c0, c1
       case 'O':
           typ[i] = 1
           c0 := readInt(reader)
           c1 := readInt(reader)
           ch0[i], ch1[i] = c0, c1
       case 'X':
           typ[i] = 2
           c0 := readInt(reader)
           c1 := readInt(reader)
           ch0[i], ch1[i] = c0, c1
       }
   }
   // compute values in post-order
   order := make([]int, 0, n)
   stack := make([]int, 0, n)
   stack = append(stack, 1)
   for len(stack) > 0 {
       u := stack[len(stack)-1]
       stack = stack[:len(stack)-1]
       order = append(order, u)
       if typ[u] >= 0 {
           if typ[u] == 3 {
               stack = append(stack, ch0[u])
           } else {
               stack = append(stack, ch1[u])
               stack = append(stack, ch0[u])
           }
       }
   }
   for i := len(order) - 1; i >= 0; i-- {
       u := order[i]
       switch typ[u] {
       case -1:
           // input, value already set
       case 3:
           g[u] = 1 - g[ch0[u]]
       default:
           a, b := g[ch0[u]], g[ch1[u]]
           if typ[u] == 0 {
               g[u] = a & b
           } else if typ[u] == 1 {
               g[u] = a | b
           } else {
               g[u] = a ^ b
           }
       }
   }
   // propagate influence (need)
   need := make([]bool, n+1)
   need[1] = true
   queue := make([]int, 0, n)
   queue = append(queue, 1)
   for qi := 0; qi < len(queue); qi++ {
       u := queue[qi]
       switch typ[u] {
       case -1:
           continue
       case 3:
           v := ch0[u]
           if need[u] && !need[v] {
               need[v] = true
               queue = append(queue, v)
           }
       case 0:
           l, r := ch0[u], ch1[u]
           if need[u] {
               if g[r] == 1 && !need[l] {
                   need[l] = true
                   queue = append(queue, l)
               }
               if g[l] == 1 && !need[r] {
                   need[r] = true
                   queue = append(queue, r)
               }
           }
       case 1:
           l, r := ch0[u], ch1[u]
           if need[u] {
               if g[r] == 0 && !need[l] {
                   need[l] = true
                   queue = append(queue, l)
               }
               if g[l] == 0 && !need[r] {
                   need[r] = true
                   queue = append(queue, r)
               }
           }
       case 2:
           l, r := ch0[u], ch1[u]
           if need[u] {
               if !need[l] {
                   need[l] = true
                   queue = append(queue, l)
               }
               if !need[r] {
                   need[r] = true
                   queue = append(queue, r)
               }
           }
       }
   }
   rootVal := g[1]
   // output results for input leaves
   for i := 1; i <= n; i++ {
       if typ[i] == -1 {
           var out int
           if need[i] {
               out = 1 - rootVal
           } else {
               out = rootVal
           }
           writer.WriteByte(byte(out + '0'))
       }
   }
   writer.WriteByte('\n')
}

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

func readOp(r *bufio.Reader) byte {
   c, err := r.ReadByte()
   for err == nil && (c < 'A' || c > 'Z') {
       c, err = r.ReadByte()
   }
   if err != nil {
       return 0
   }
   return c
}

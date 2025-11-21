package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"time"
)

type fastScanner struct {
	r *bufio.Reader
}

func newFastScanner() *fastScanner {
	return &fastScanner{r: bufio.NewReader(os.Stdin)}
}

func (fs *fastScanner) nextInt64() int64 {
	var sign int64 = 1
	var val int64
	c, err := fs.r.ReadByte()
	for (c < '0' || c > '9') && c != '-' {
		c, err = fs.r.ReadByte()
		if err != nil {
			return 0
		}
	}
	if c == '-' {
		sign = -1
		c, err = fs.r.ReadByte()
		if err != nil {
			return 0
		}
	}
	for c >= '0' && c <= '9' {
		val = val*10 + int64(c-'0')
		c, err = fs.r.ReadByte()
		if err != nil {
			break
		}
	}
	return val * sign
}

type node struct {
	val      int64
	priority int
	cnt      int
	size     int
	left     *node
	right    *node
}

func nodeSize(n *node) int {
	if n == nil {
		return 0
	}
	return n.size
}

func (n *node) maintain() {
	if n != nil {
		n.size = n.cnt + nodeSize(n.left) + nodeSize(n.right)
	}
}

func rotateRight(n *node) *node {
	l := n.left
	n.left = l.right
	l.right = n
	n.maintain()
	l.maintain()
	return l
}

func rotateLeft(n *node) *node {
	r := n.right
	n.right = r.left
	r.left = n
	n.maintain()
	r.maintain()
	return r
}

type treap struct {
	root *node
	rnd  *rand.Rand
}

func newTreap() *treap {
	return &treap{
		rnd: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (t *treap) insert(n *node, val int64) *node {
	if n == nil {
		return &node{
			val:      val,
			priority: t.rnd.Int(),
			cnt:      1,
			size:     1,
		}
	}
	if val == n.val {
		n.cnt++
	} else if val < n.val {
		n.left = t.insert(n.left, val)
		if n.left.priority > n.priority {
			n = rotateRight(n)
		}
	} else {
		n.right = t.insert(n.right, val)
		if n.right.priority > n.priority {
			n = rotateLeft(n)
		}
	}
	n.maintain()
	return n
}

func (t *treap) delete(n *node, val int64) *node {
	if n == nil {
		return nil
	}
	if val == n.val {
		if n.cnt > 1 {
			n.cnt--
		} else {
			if n.left == nil {
				return n.right
			}
			if n.right == nil {
				return n.left
			}
			if n.left.priority > n.right.priority {
				n = rotateRight(n)
				n.right = t.delete(n.right, val)
			} else {
				n = rotateLeft(n)
				n.left = t.delete(n.left, val)
			}
		}
	} else if val < n.val {
		n.left = t.delete(n.left, val)
	} else {
		n.right = t.delete(n.right, val)
	}
	if n != nil {
		n.maintain()
	}
	return n
}

func (t *treap) Insert(val int64) {
	t.root = t.insert(t.root, val)
}

func (t *treap) Delete(val int64) {
	t.root = t.delete(t.root, val)
}

func (t *treap) lowerBound(val int64) (int64, bool) {
	cur := t.root
	var res *node
	for cur != nil {
		if cur.val >= val {
			res = cur
			cur = cur.left
		} else {
			cur = cur.right
		}
	}
	if res == nil {
		return 0, false
	}
	return res.val, true
}

func (t *treap) popLowerBound(val int64) (int64, bool) {
	v, ok := t.lowerBound(val)
	if !ok {
		return 0, false
	}
	t.Delete(v)
	return v, true
}

type monster struct {
	need   int64
	reward int64
}

func main() {
	fs := newFastScanner()
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	T := int(fs.nextInt64())
	for ; T > 0; T-- {
		n := int(fs.nextInt64())
		m := int(fs.nextInt64())
		swords := make([]int64, n)
		for i := 0; i < n; i++ {
			swords[i] = fs.nextInt64()
		}
		needs := make([]int64, m)
		for i := 0; i < m; i++ {
			needs[i] = fs.nextInt64()
		}
		withReward := make([]monster, 0)
		withoutReward := make([]int64, 0)
		for i := 0; i < m; i++ {
			c := fs.nextInt64()
			if c > 0 {
				withReward = append(withReward, monster{need: needs[i], reward: c})
			} else {
				withoutReward = append(withoutReward, needs[i])
			}
		}
		tr := newTreap()
		for _, s := range swords {
			tr.Insert(s)
		}
		var ans int64

		sort.Slice(withReward, func(i, j int) bool {
			if withReward[i].need == withReward[j].need {
				return withReward[i].reward > withReward[j].reward
			}
			return withReward[i].need < withReward[j].need
		})
		for _, mon := range withReward {
			val, ok := tr.popLowerBound(mon.need)
			if !ok {
				break
			}
			if mon.reward > val {
				val = mon.reward
			}
			tr.Insert(val)
			ans++
		}

		sort.Slice(withoutReward, func(i, j int) bool {
			return withoutReward[i] < withoutReward[j]
		})
		for _, need := range withoutReward {
			_, ok := tr.popLowerBound(need)
			if !ok {
				break
			}
			ans++
		}
		fmt.Fprintln(out, ans)
	}
}

#include <cassert>
#include <cctype>
#include <cstdio>
#include <cstdlib>
#include <cstring>
#include <algorithm>
#include <bitset>
#include <iomanip>
#include <iostream>
#include <map>
#include <queue>
#include <set>
#include <sstream>
#include <string>
#include <utility>
#include <vector>
typedef long long LL;
typedef unsigned long long ULL;
typedef long double LD;
using std::cin; using std::cout;
using std::endl;
using std::bitset; using std::map;
using std::queue; using std::priority_queue;
using std::set; using std::string;
using std::stringstream; using std::vector;
using std::pair; using std::make_pair;
typedef pair<int, int> pii;
typedef pair<LL, LL> pll;
typedef pair<ULL, ULL> puu;
#ifdef DEBUG
using std::cerr;
#define pass cerr << "[" << __FUNCTION__ << "] : line = " << __LINE__ << endl;
#define display(x) cerr << #x << " = " << x << endl;
#define displaya(a, st, n)                                                     \
  {                                                                            \
    cerr << #a << " = {";                                                      \
    for (int qwq = (st); qwq <= (n); ++qwq)                                    \
      if (qwq == (st))                                                         \
        cerr << a[qwq];                                                        \
      else                                                                     \
        cerr << ", " << a[qwq];                                                \
    cerr << "}" << endl;                                                       \
  }
#define displayv(a) displaya(a, 0, (int)(a.size() - 1))
#include <ctime>
class MyTimer {
    clock_t st;
public:
    MyTimer() { cerr << std::fixed << std::setprecision(0); reset(); }
    ~MyTimer() { report(); }
    void reset() { st = clock_t(); }
    void report() {  cerr << "Time consumed: " << (clock() - st) * 1e3 / CLOCKS_PER_SEC << "ms" << endl; }
} myTimer;
#else
#define cerr if(false) std::cout
#define pass ;
#define display(x) ;
#define displaya(a, st, n) ;
#define displayv(a) ;
class MyTimer {
public: void reset() {} void report() {}
} myTimer;
#endif

template<typename A, typename B>
std::ostream& operator << (std::ostream &cout, const pair<A, B> &x) {
    return cout << "(" << x.first << ", " << x.second << ")";
}
template<typename T1, typename T2>
inline bool chmin(T1 &a, const T2 &b) { return a > b ? a = b, true : false; }
template<typename T1, typename T2>
inline bool chmax(T1 &a, const T2 &b) { return a < b ? a = b, true : false; }
#ifdef QUICK_READ
char pool[1<<15|1],*it=pool+32768;
#define getchar() (it>=pool+32768?(pool[fread(pool,sizeof(char),1<<15,stdin)]=EOF,*((it=pool)++)):*(it++))
#endif
inline int readint() {
    int a = 0; char c = getchar(), p = 0;
    while(isspace(c)) c = getchar();
    if(c == '-') p = 1, c = getchar();
    while(isdigit(c)) a = a*10 + c - '0', c = getchar();
    return p ? -a : a;
}

const int maxN = 200000 + 233;
int n, m, a[maxN], b[maxN];
int ans[maxN];
struct Item {
    int type, coef;
    // type == 0 : val; type == ? : query
    int time, pos, val;
    Item() {}
    Item(int t, int c, int tm, int p, int v)
        : type(t), coef(c), time(tm), pos(p), val(v) {}
} q[maxN * 8], aux[maxN * 8];
int len = 0;
int C[maxN];
void add(int p, int x) {
    for(int i = p; i <= n; i += i & -i) C[i] += x;
}
int sum(int p) {
    int r = 0;
    for(int i = p; i > 0; i -= i & -i) r += C[i];
    return r;
}

void solve(int l, int r) {
    if(l == r) return;
    int m = (l + r) >> 1;
    solve(l, m); solve(m + 1, r);
    int i = l, j = m + 1;
    int iter = l;
    auto processL = [&](int pos) {
        Item &I = q[pos];
        if(I.type == 0) {
            add(I.val, I.coef);
        }
        aux[iter++] = q[pos];
    };
    auto processR = [&](int pos) {
        Item &I = q[pos];
        if(I.type) {
            ans[I.type] += I.coef * sum(I.val);
        }
        aux[iter++] = q[pos];
    };
    while(i <= m && j <= r) {
        if(q[i].pos <= q[j].pos) processL(i), i++;
        else processR(j), j++;
    }
    while(i <= m) processL(i), i++;
    while(j <= r) processR(j), j++;
    for(int k = l; k <= m; ++k) {
        if(q[k].type == 0) add(q[k].val, -q[k].coef);
    }
    for(int k = l; k <= r; ++k) q[k] = aux[k];
}

/*
4 1
1 2 3 4
1 2 3 4
1 1 3 2 4
*/

int main() {
    n = readint(); m = readint();
    for(int i = 1; i <= n; ++i) a[readint()] = i;
    for(int i = 1; i <= n; ++i) b[i] = a[readint()];
    for(int i = 1; i <= n; ++i) q[++len] = Item(0, 1, 0, i, b[i]);
    for(int i = 1; i <= m; ++i) {
        int op = readint();
        if(op == 1) {
            int vl = readint(), vr = readint(),
                pl = readint(), pr = readint();
            q[++len] = Item(i, 1, i, pr, vr);
            if(pl) q[++len] = Item(i, -1, i, pl - 1, vr);
            if(vl) q[++len] = Item(i, -1, i, pr, vl - 1);
            if(pl && vl) q[++len] = Item(i, 1, i, pl - 1, vl - 1);
        } else {
            int x = readint(), y = readint();
            q[++len] = Item(0, 1, i, x, b[y]);
            q[++len] = Item(0, 1, i, y, b[x]);
            q[++len] = Item(0, -1, i, x, b[x]);
            q[++len] = Item(0, -1, i, y, b[y]);
            std::swap(b[x], b[y]);
            ans[i] = -maxN;
        }
    }
    std::sort(q + 1, q + len + 1, [&](const Item &x, const Item &y) -> bool {
        return x.time < y.time;
    });
    solve(1, len);
//    for(int i = 1; i <= len; ++i) if(q[i].type == 0) {
//        for(int j = 1; j <= len; ++j) if(q[j].type != 0) {
//            ans[q[j].type] += (q[i].time <= q[j].time && q[i].pos <= q[j].pos && q[i].val <= q[j].val) * q[i].coef * q[j].coef;
//        }
//    }
    for(int i = 1; i <= m; ++i) if(ans[i] != -maxN) printf("%d\n", ans[i]);
    return 0;
}
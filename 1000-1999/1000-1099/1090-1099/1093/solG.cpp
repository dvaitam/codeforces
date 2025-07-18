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
const int INF = 0x3f3f3f3f;
const int maxK = 5;
int n, k, q;
int a[maxN][maxK];

#define lson (o << 1)
#define rson (lson | 1)
int rv[maxK], nv[1 << maxK];
int max[maxN * 4][1 << maxK];
void maintain(int o) {
    for(int S = 0; S < (1 << k); ++S)
        max[o][S] = std::max(max[lson][S], max[rson][S]);
}
void build(int o, int L, int R) {
    if(L == R) {
        for(int S = 0; S < (1 << k); ++S) max[o][S] = 0;
        for(int S = 0; S < (1 << k); ++S)
            for(int j = 0; j < k; ++j)
                if(S >> j & 1) max[o][S] += a[L][j]; else max[o][S] -= a[R][j];
    } else {
        int M = (L + R) >> 1;
        build(lson, L, M);
        build(rson, M + 1, R);
        maintain(o);
    }
}
int ql, qr;
int qmax[1 << maxK];
void query(int o, int L, int R) {
    if(ql <= L && R <= qr) {
        for(int S = 0; S < (1 << k); ++S)
            chmax(qmax[S], max[o][S]);
    } else {
        int M = (L + R) >> 1;
        if(ql <= M) query(lson, L, M);
        if(qr > M) query(rson, M + 1, R);
    }
}
void modify(int o, int L, int R) {
    if(L == R) {
        for(int S = 0; S < (1 << k); ++S) max[o][S] = 0;
        for(int S = 0; S < (1 << k); ++S)
            for(int j = 0; j < k; ++j)
                if(S >> j & 1) max[o][S] += a[L][j]; else max[o][S] -= a[R][j];
    } else {
        int M = (L + R) >> 1;
        if(ql <= M) modify(lson, L, M);
        else modify(rson, M + 1, R);
        maintain(o);
    }
}

int main() {
    n = readint(); k = readint();
    for(int i = 1; i <= n; ++i) {
        for(int j = 0; j < k; ++j) a[i][j] = readint();
    }
    build(1, 1, n);
    q = readint();
    while(q--) {
        int op = readint();
        if(op == 1) {
            int i = readint();
            for(int j = 0; j < k; ++j) a[i][j] = readint();
            ql = i;
            modify(1, 1, n);
        } else {
            ql = readint(); qr = readint();
            for(int S = 0; S < (1 << k); ++S) qmax[S] = ~INF;
            query(1, 1, n);
            int ans = ~INF;
            for(int S = 0; S < (1 << k); ++S)
                chmax(ans, qmax[S] + qmax[((1 << k) - 1) ^ S]);
            printf("%d\n", ans);
        }
    }
    return 0;
}
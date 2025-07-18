/*
 * Author: SolitaryWayfarer
 * https://github.com/TheLonelyHeracles/Cplusplus-Improvement-Library
 */
#include <algorithm>
#include <cassert>
#include <cmath>
#include <cstdio>
#include <cstring>
#include <cstdlib>
#include <functional>
#include <deque>
#include <iostream>
#include <list>
#include <map>
#include <memory>
#include <set>
#include <unordered_map>
#include <unordered_set>
#include <utility>
#include <vector>
using namespace std;
#ifndef ONLINE_JUDGE
#include "Cplusplus-Improvement-Library/debug_kit.h"
#else
#define watch(...)
#define OUT(...)
#define _OUT(...)
#define OUT_(...)
#endif

// read integer types
template<typename T> inline bool RD(T& v) { char c; bool n; while (c = getchar(), c != '-' && (c < '0' || c > '9')) if (c == EOF) return false; if (c == '-') { n = true; v = 0; } else { n = false; v = c - '0'; } while (c = getchar(), c >= '0' && c <= '9') v = (v << 3) + (v << 1) + c - '0'; if (n) v *= -1; return true; }
template<typename A, typename ...Args> inline bool RD(A& a, Args&... rest) { return RD(a) && RD(rest...); }
inline int RD() { int v; RD(v); return v; }
#define RDn(a, l, r)          { for (int _i_ = l; _i_ < r; ++_i_) RD((a)[_i_]); }
#define RDnm(a, b, c, d, e)   { for (int __i__ = b; __i__ < c; ++__i__) RDn(a[__i__], d, e) }

// read integer types with buffer
inline char getchar_buffered() { static char _BUF_[1 << 15], *_HEAD_ = _BUF_, *_TAIL_ = _BUF_; return _HEAD_ == _TAIL_ && (_TAIL_ = (_HEAD_ = _BUF_) + fread(_BUF_, 1, 1 << 15, stdin), _HEAD_ == _TAIL_) ? EOF : *_HEAD_++; }
template<typename T> inline bool RDB(T& v) { char c; bool n; while (c = getchar_buffered(), c != '-' && (c < '0' || c > '9')) if (c == EOF) return false; if (c == '-') { n = true; v = 0; } else { n = false; v = c - '0'; } while (c = getchar_buffered() - '0', c >= 0 && c <= 9) v = (v << 3) + (v << 1) + c; if (n) v *= -1; return true; }
template<typename A, typename ...Args> inline bool RDB(A& a, Args&... rest) { return RDB(a) && RDB(rest...); }
inline int RDB() { int v; RDB(v); return v; }
#define RDBn(a, l, r)         { for (int _i_ = l; _i_ < r; ++_i_) RDB((a)[_i_]); }
#define RDBnm(a, b, c, d, e)  { for (int __i__ = b; __i__ < c; ++__i__) RDBn(a[__i__], d, e) }

// write integer types
template<typename T> inline void _WR(T a) { if (a < 0) { putchar('-'); a *= -1; } T t = a / 10; if (t) _WR(t); putchar(a - (t << 1) - (t << 3) + '0'); }
template<typename T> inline void WR_(const T&a) { _WR(a); putchar(' '); }
template<typename T> inline void WR(const T&a) { _WR(a); putchar('\n'); }
template<typename A, typename ...Args> inline void _WR(const A& a, const Args&... rest) { WR_(a); _WR(rest...); }
template<typename A, typename ...Args> inline void WR_(const A& a, const Args&... rest) { WR_(a); WR_(rest...); }
template<typename A, typename ...Args> inline void WR (const A& a, const Args&... rest) { WR_(a); WR(rest...); }
#define WRn(a, l, r)          { for (int _i_ = l; _i_ < r - 1; ++_i_) WR_((a)[_i_]); if(r > 0) WR((a)[r - 1]); }
#define WRnm(a, b, c, d, e)   { for (int __i__ = b; __i__ < c; ++__i__) WRn(a[__i__], d, e) }
///////


template<typename T, typename Operator = std::plus<T>>
class binary_indexed_tree {
    std::vector<T> __c;
    Operator __op;
    int __n;
public:
    binary_indexed_tree() : __n(0) {}
    binary_indexed_tree(int n, T v = 0) { __c.resize(n, v); __n = n; }
    void clear() { __c.clear(); __n = 0; }
    void fill(T v) { for (int i = 0; i < __n; ++i) __c[i] = v; }
    void resize(int n, T v = 0) { __c.clear(); __c.resize(n, v); __n = n; }
    int size() { return __n; }
    void add(int x, T d) {
        for (int i = x; i < __n; i |= i + 1)
            __c[i] = __op(__c[i], d);
    }
    T accu(int x) {
        T ret;
        if (__op(1, 0) == 1) ret = 0;
        else if (__op(1, 1) == 1) ret = 1;
        else if (__op(1, T(-1)) == 1) ret = T(-1);
        else assert(0);
        for (int i = x; i >= 0; i = (i & i + 1) - 1)
            ret = __op(ret, __c[i]);
        return ret;
    }
    T accu(int l, int r) { return accu(r) - accu(l - 1); }
};

////

#define FOR(i, a, b) for (int i = a; i < (int)(b); ++i)
#define ROF(i, a, b) for (int i = (int)(b) - 1; i >= a; --i)
typedef vector<int> VI;
typedef pair<int, int> PII;
typedef long long ll;


const int inf = 0x3f3f3f3f;
const int maxn = 2e5 + 10;

int n, d;
int s[maxn], p[maxn];

int main() {
    RDB(n, d);
    RDBn(s, 0, n);
    RDBn(p, 0, n);

    int cnt = 0;
    int score = s[d - 1];
    score += p[0];

    int pl = 1, pr = n - 1;
    FOR(i, 0, d - 1) {
        if (s[i] + p[pr] > score) {
            ++cnt;
            pl++;
            continue;
        }
        pr--;
    }

    WR(cnt + 1);

    return 0;
}
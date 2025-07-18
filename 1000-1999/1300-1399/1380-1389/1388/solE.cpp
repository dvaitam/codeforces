#ifndef _GLIBCXX_NO_ASSERT

#include <cassert>

#endif

#include <cctype>

#include <cerrno>

#include <cfloat>

#include <ciso646>

#include <climits>

#include <clocale>

#include <cmath>

#include <csetjmp>

#include <csignal>

#include <cstdarg>

#include <cstddef>

#include <cstdio>

#include <cstdlib>

#include <cstring>

#include <ctime>

 

#if __cplusplus >= 201103L

#include <ccomplex>

#include <cfenv>

#include <cinttypes>

//#include <cstdalign>

#include <cstdbool>

#include <cstdint>

#include <ctgmath>

#include <cwchar>

#include <cwctype>

#endif

 

// C++

#include <algorithm>

#include <bitset>

#include <complex>

#include <deque>

#include <exception>

#include <fstream>

#include <functional>

#include <iomanip>

#include <ios>

#include <iosfwd>

#include <iostream>

#include <istream>

#include <iterator>

#include <limits>

#include <list>

#include <locale>

#include <map>

#include <memory>

#include <new>

#include <numeric>

#include <ostream>

#include <queue>

#include <set>

#include <sstream>

#include <stack>

#include <stdexcept>

#include <streambuf>

#include <string>

#include <typeinfo>

#include <utility>

#include <valarray>

#include <vector>

 

#if __cplusplus >= 201103L

#include <array>

#include <atomic>

#include <chrono>

#include <condition_variable>

#include <forward_list>

#include <future>

#include <initializer_list>

#include <mutex>

#include <random>

#include <ratio>

#include <regex>

#include <scoped_allocator>

#include <system_error>

#include <thread>

#include <tuple>

#include <typeindex>

#include <type_traits>

#include <unordered_map>

#include <unordered_set>

#endif

using namespace std;

using i64 = long long;

using i128 = __int128;

#define MAXN 2005

#define K 38

#define MAXP 25

#define MAXK 55

#define MAXC 255

#define MAXERR 105

#define MAXLEN 105

#define MDIR 10

#define MAXR 705

#define BASE 102240

#define MAXA 28

#define MAXT 100005

#define LIMIT 86400

#define MAXV 305

#define OP 0

#define CLO 1

#define DIG 1

#define C1 0

#define C2 1

#define PLUS 0

#define MINUS 1

#define MUL 2

#define CLO 1

#define VERT 1

#define B 31

#define B2 1007

#define W 1

#define H 18

#define SPEC 1

#define MUL 2

#define CNT 3

#define ITER 1000

#define INF 1e18

#define EPS 1e-9

#define MOD 998244353

#define CONST 998244353

#define FACT 100000000000000

#define PI 3.14159265358979

#define SRC 0

#define pb push_back

#define eb emplace_back

typedef long double ll;

typedef long double ld;

typedef pair<int,int> ii;

typedef pair<int,ll> li;

typedef tuple<ll,ll,ll> iii;

typedef vector<vector<int>> vv;

typedef vector<int> vi;

typedef pair<ld,int> iv;

typedef vector<ii> vii;

typedef complex<double> cd;

typedef int hash_t;

#define sc second

#define fr first

#define rep(i,x,y) for (int i = (x); i < (y); ++i)

#define rev(i,x,y) for (int i = (x); i >= (y); --i)

#define LSOne(S) (S & (-S))

#define trav(i,v) for (auto &i : v)

#define foreach(it,v) for (auto it = begin(v); it != end(v); ++it)

#define sortarr(v) sort(begin(v), end(v))



struct Line {

    mutable ll k, m, p;

    bool operator<(const Line& o) const { return k < o.k; }

    bool operator<(ll x) const { return p < x; }

};

 

//this Line Container finds the maximum cost, negate lines to find the minimum cost!

struct LineContainer : multiset<Line, less<>> {

    // (for doubles, use inf = 1/.0, div(a,b) = a/b)

    const ll inf = 1/.0;

    ll div(ll a, ll b) { // floored division

        return a / b; }

    bool isect(iterator x, iterator y) {

        if (y == end()) { x->p = inf; return 0; }

        if (x->k == y->k) x->p = x->m > y->m ? inf : -inf;

        else x->p = div(y->m - x->m, x->k - y->k);

        return x->p >= y->p;

    }

    void add(ll k, ll m) {

        auto z = insert({k, m, 0}), y = z++, x = y;

        while (isect(y, z)) z = erase(z);

        if (x != begin() && isect(--x, y)) isect(x, y = erase(y));

        while ((y = x) != begin() && (--x)->p >= y->p)

            isect(x, erase(y));

    }

    ll query(ll x) {

        assert(!empty());

        auto l = *lower_bound(x);

        return l.k * x + l.m;

    }

};



int xl[MAXN], xr[MAXN], y[MAXN];

int main() {

    int n; cin >> n;

    rep(i,0,n) cin >> xl[i] >> xr[i] >> y[i];

    LineContainer lmin, lmax;

    //invert (x,y) into (y,x) take all angles w.r.t to the y-axis

    //normalise to change all x coordinates as the intercept when y = 0

    //symmetric projection when scaled by a particular vector

    rep(i,0,n) {

        lmin.add(-y[i], -xl[i]);

        lmin.add(-y[i], -xr[i]);

        

        lmax.add(y[i], xl[i]);

        lmax.add(y[i], xr[i]);

    }

    vector<pair<double,double>> segs;

    //calculate all possible ranges on angles here

    rep(i,0,n) {

        rep(j,i+1,n) {

            if (y[i] == y[j]) continue;

            double a = xr[j] - xl[i];

            double b = xl[j] - xr[i];

            a /= (double)(y[i] - y[j]);

            b /= (double)(y[i] - y[j]);

            segs.eb(min(a,b), max(a,b));

        }

    }

    //sort by start points

    sort(segs.begin(), segs.end());

    vector<pair<double,double>> mr;

    for (auto &[st,ed] : segs) {

        if (mr.empty() || mr.back().sc <= st) {

            mr.eb(st,ed);

        } else mr.back().sc = max(mr.back().sc, ed);

    }

    if (mr.empty()) mr.eb(0,0);

    ld ans = INF;

    for (auto &[st,ed] : mr) {

        ans = min(ans, lmin.query(st) + lmax.query(st));

        ans = min(ans, lmin.query(ed) + lmax.query(ed));

    }

    cout << fixed << setprecision(15) << ans;

}
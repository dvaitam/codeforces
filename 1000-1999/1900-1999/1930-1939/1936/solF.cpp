#pragma gcc optimize("O3,Ofast,no-stack-protector,fast-math,section-anchors,inline,unroll-loops")
#pragma gcc target("abm,mmx,sse,sse2,sse3,sse4,ssse3,tune=native,lzcnt,popcnt,avx,avx2")
//#include<immintrin.h>
#include<iostream>
#include<iomanip>
#include<algorithm>
#include<array>
#include<vector>
#include<set>
#include<string>
#include<map>
#include<tuple>
#include<cmath>
#include<deque>
#include<ctime>
#include<random>
#include<bitset>
#include<unordered_map>
#include<unordered_set>
#include<chrono>
#include<queue>
#include<fstream>
#include<cassert>
#include<thread>
using namespace std;

/*
//#include<bits/stdc++.h>
#include<ext/pb_ds/assoc_container.hpp>
#include<ext/pb_ds/tree_policy.hpp>
using namespace __gnu_pbds;
using ordered_set = tree<int,null_type,std::less<int>,rb_tree_tag,tree_order_statistics_node_update>;
*/
using ll = long long;
using ld = long double;
//using li = __int128;
using ui = unsigned int;
using ul = unsigned long long;
#define ft first
#define sc second
#define all(x) x.begin(), x.end()
template <typename T1, typename T2> ostream& operator << (ostream &c, const pair<T1, T2> &v);
template <typename T> ostream& operator << (ostream &c, const vector<T> &v) {
    c << '{';
    for (int i = 0; i < v.size(); i++) {
        c << v[i];
        if (i+1 != v.size()) c << ", ";
    }
    c << '}';
    return c;
}
template <typename T1, typename T2> ostream& operator << (ostream &c, const pair<T1, T2> &v) {
    c << '<' << v.ft << ", " << v.sc << '>';
    return c;
}
#define debug(x) cout << (#x) << " = " << x << endl;
const ll inf = 1e18, mod = 1e9+7, N = 1e6;
const long double pi = acos(-1);
struct vect {
    ld x, y;
    vect(ld x_ = 0, ld y_ = 0) {
        x = x_;
        y = y_;
    }
    ld len() {
        return sqrt(x*x+y*y);
    }
    ld lensq() {
        return x*x+y*y;
    }
};
istream& operator>>(istream& c, vect& v) {
    c >> v.x >> v.y;
    return c;
}
ostream& operator<<(ostream& c, vect v) {
    c << v.x << " " << v.y;
    return c;
}
vect operator+(const vect& a, const vect& b) {
    return vect(a.x + b.x, a.y + b.y);
}
vect operator-(const vect& a, const vect& b) {
    return vect(a.x - b.x, a.y - b.y);
}
vect operator* (const vect& a, ld k) {
    return vect(a.x * k, a.y * k);
}
void operator*= (vect& a, ld k) {
    a = a*k;
}
ld operator*(vect a, vect b) {
    return a.x * b.x + a.y * b.y;
}
ld operator^(vect a, vect b) {
    return a.x * b.y - a.y * b.x;
}
void operator+=(vect& a, vect b) {
    a = a + b;
}
void operator-=(vect& a, vect b) {
    a = a - b;
}
bool operator == (vect a, vect b) {
    return (a.x == b.x && a.y == b.y);
}
vect operator-(const vect& a) {
    return vect(-a.x, -a.y);
}
struct line {
    ld a, b, c;
    line(ld a_ = 0, ld b_ = 0, ld c_ = 0) {
        a = a_;
        b = b_;
        c = c_;
    }
};
istream& operator>>(istream& c, line& l) {
    c >> l.a >> l.b >> l.c;
    return c;
}
ostream& operator<<(ostream& c, line l) {
    c << l.a << " " << l.b << " " << l.c;
    return c;
}
line line_by_vects(vect a, vect b) {
    return line(b.y-a.y, a.x-b.x, b.x*a.y-b.y*a.x);
}
pair<int, vect> ill(line x, line y) {
    if (x.a*y.b == y.a*x.b && x.b*y.c == x.c*y.b) return {2, vect(0, 0)};
    if (x.a*y.b == y.a*x.b) return {0, vect(0, 0)};
    return {1, vect((y.c*x.b-x.c*y.b)/(x.a*y.b-y.a*x.b), (y.c*x.a-x.c*y.a)/(y.a*x.b-x.a*y.b))};
}
ld eps = 1e-10; // height, isin, f3 - triangle
pair<ld, vect> height_gr(line l, vect v) {
    ld d = abs((l.a*v.x+l.b*v.y+l.c)/sqrt(l.a*l.a+l.b*l.b));
    vect n(l.a, l.b), n1 = n*(d/n.len()), n2 = n*(-d/n.len());
    if (abs((v.x+n2.x)*l.a+(v.y+n2.y)*l.b+l.c) <= eps) n1 = n2;
    return {d, v+n1};
}
line mid_per(vect u, vect v) {
    line l = line_by_vects(u, v);
    vect m = (u+v)*0.5;
    return line(-l.b, l.a, l.b*m.x-l.a*m.y);
}
struct circle {
    vect v;
    ld r;
    circle(vect v_ = {0, 0}, ld _r = 0) {
        v = v_;
        r = _r;
    }
};
istream& operator>>(istream& c, circle& r) {
    c >> r.v.x >> r.v.y >> r.r;
    return c;
}
ostream& operator<<(ostream& c, circle r) {
    c << r.v.x << " " << r.v.y << " " << r.r;
    return c;
}
vector<vect> ilc(line l, circle c) {
    vect col{-l.b, l.a};
    auto [d, g] = height_gr(l, c.v);
    ld r = sqrt(c.r*c.r-d*d);
    col *= (r/col.len());
    if ((g-c.v).len() > c.r) return {};
    if ((g-c.v).len() == c.r) return {g};
    return {g-col, g+col};
}
vector<vect> icc(circle c1, circle c2) {
    if (c1.v == c2.v) {
        if (c1.r != c2.r) {
            return {};
        } else {
            return {{0, 0}, {0, 0}, {0, 0}};
        }
    }
    ld cc1 = c1.r*c1.r-c1.v.lensq(), cc2 = c2.r*c2.r-c2.v.lensq();
    vect r = c2.v-c1.v;
    return ilc(line(2*r.x, 2*r.y, cc2-cc1), c1);
}
pair<bool, circle> circle_by_vects(vect u, vect v, vect w) {
    line m1 = mid_per(u, v), m2 = mid_per(v, w);
    auto [ch, O] = ill(m1, m2);
    if (ch == 1) {
        return {1, circle(O, vect(O - v).len())};
    } else {
        return {0, circle({0, 0}, 0)};
    }
}
bool isin(circle a, circle b) {
    return (a.v-b.v).len() <= b.r-a.r+eps;
}
circle f2(circle a, circle b) {
    vect v = b.v - a.v, va = v*(a.r/v.len()), vb = v*(1-b.r/v.len());
    return circle(a.v+((va+vb)*0.5), (va-vb).len()/2);
}
circle inv(circle a, circle b) {
    vect v1 = a.v + vect(0, a.r), v2 = a.v + vect(0, -a.r), v3 = a.v + vect(a.r, 0);
    v1 -= b.v; v2 -= b.v; v3 -= b.v;
    v1 = v1*(b.r*b.r/v1.lensq()); v2 = v2*(b.r*b.r/v2.lensq()); v3 = v3*(b.r*b.r/v3.lensq());
    return circle_by_vects(v1+b.v, v2+b.v, v3+b.v).sc;
}
line inv2(circle a, circle b) {
    vect v = (a.v-b.v)*2;
    v = v*(b.r*b.r/v.lensq());
    return line_by_vects(b.v+v, b.v+v+vect(-v.y, v.x));
}
vect incenter(vect a, vect b, vect c) {
    ld lc = (a-b).len(), la = (b-c).len(), lb = (a-c).len();
    return ill(line_by_vects(a, ((b*lb)+(c*lc))*(1/(lb+lc))),
               line_by_vects(b, ((a*la)+(c*lc))*(1/(la+lc)))).sc;
}
circle f3(circle a, circle b, circle c) {
    ld x1 = a.v.x, y1 = a.v.y, r1 = a.r;
    ld x2 = b.v.x, y2 = b.v.y, r2 = b.r;
    ld x3 = c.v.x, y3 = c.v.y, r3 = c.r;
    ld a2 = x1 - x2, a3 = x1 - x3, b2 = y1 - y2, b3 = y1 - y3, c2 = r2 - r1, c3 = r3 - r1;
    ld d1 = x1 * x1 + y1 * y1 - r1 * r1, d2 = d1 - x2 * x2 - y2 * y2 + r2 * r2,
            d3 = d1 - x3 * x3 - y3 * y3 + r3 * r3;
    ld ab = a3 * b2 - a2 * b3;
    ld xa = (b2 * d3 - b3 * d2) / (ab * 2) - x1;
    ld xb = (b3 * c2 - b2 * c3) / ab;
    ld ya = (a3 * d2 - a2 * d3) / (ab * 2) - y1;
    ld yb = (a2 * c3 - a3 * c2) / ab;
    ld A = xb * xb + yb * yb - 1;
    ld B = 2 * (r1 + xa * xb + ya * yb);
    ld C = xa * xa + ya * ya - r1 * r1;
    ld r = -(A ? (B - sqrtl(B * B - 4 * A * C)) / (2 * A) : C / B);
    return circle(vect(x1 + xa + xb * r, y1 + ya + yb * r), r);
}
signed main() {
    ios_base::sync_with_stdio(0);
    cin.tie(0);
    //freopen("1B/tests/152", "r", stdin);
    //freopen("macloren.txt", "w", stdout)
    mt19937 rnd(std::chrono::steady_clock::now().time_since_epoch().count());
    int n;
    cin >> n;
    vector<circle> v(n);
    for (int i = 0; i < n; i++) cin >> v[i];
    shuffle(all(v), default_random_engine(rnd()));
    circle ans = {{0, 0}, 1e18};
    for (int i = 0; i < n; i++) {
        if (!isin(ans, v[i])) {
            ans = v[i];
            for (int j = 0; j < i; j++) {
                if (!isin(ans, v[j])) {
                    ans = f2(v[i], v[j]);
                    for (int y = 0; y < j; y++) {
                        if (!isin(ans, v[y])) {
                            ans = f3(v[i], v[j], v[y]);
                        }
                    }
                }
            }
        }
    }
    cout.precision(21);
    cout << ans;
}
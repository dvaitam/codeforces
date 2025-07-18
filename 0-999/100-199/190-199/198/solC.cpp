#include <fstream>
#include <stdio.h>
#include <iostream>
#include <string>
#include <complex>
#include <math.h>
#include <cmath>
#include <cstring>
#include <cstdio>
#include <cassert>
#include <set>
#include <vector>
#include <map>
#include <queue>
#include <deque>
#include <stack>
#include <algorithm>
#include <bitset>
#include <list>

using namespace std;

#define x first
#define y second
#define pb push_back
#define mp make_pair
#define all(a) (a).begin(), (a).end()
#define rall(a) (a).rbegin(), (a).rend()
#define sz(a) int((a).size())
#define sqr(x) ((x)*(x))
#define forn(i, n) for (int i = 0; i < int(n); i++)
#define NAME "numpalindrome"
#define FREOPEN freopen(NAME".in", "r", stdin); freopen(NAME".out", "w", stdout);
#define Freopen freopen("input.txt", "r", stdin); freopen("output.txt", "w", stdout);

#define y0 sdkfaslhagaklsldk
#define y1 aasdfasdfasdf
#define yn askfhwqriuperikldjk
#define j1 assdgsdgasghsf
#define tm sdfjahlfasfh
#define lr asgasgashqt
#define free afdshjioeyqtw
#define next qpowityqwopqw

typedef unsigned int unt;
typedef long long ll;
typedef unsigned long long ull;
typedef double ld;
typedef pair < int, int > pii;
typedef pair < ll, ll > pll;
typedef double geom;

const ll MOD = 1e9 + 7;
const int INF = 2e9 + 1;
const ll INF64 = 1e18;
const ld EPS = 1e-9;
const ld PI = acos(-1);
const ll MD = 1551513443;
const ll T = 25923;
const int N = 250100;
const int M = 510;

ld xx, yy, vv, tx, ty;
ld x, y, v, rr;

ld get(ld m) {
    ld qx = x + (tx - x) * m;
    ld qy = y + (ty - y) * m;
    return sqrt(sqr(qx) + sqr(qy));
}

pair < ld, ld > get(ld x, ld y) {
    ld c = sqrt(sqr(x) + sqr(y)), a = rr;
    ld b = sqrt(sqr(c) - sqr(a));
    return mp(b, acos(a / c));
};

ld check(ld t) {
    ld l = 0, r = 1, m1, m2;
    forn(u, 400) {
        ld m = (r - l) / 3;
        m1 = l + m; m2 = r - m;
        if (get(m1) <= get(m2))
            r = m2;
        else
            l = m1;
    }
    if (get(l) >= rr)
        return sqrt(sqr(x - tx) + sqr(y - ty));
    ld spang = abs(atan2(y, x) - atan2(ty, tx));
    if (2 * PI - spang < spang)
        spang = 2 * PI - spang;
    pair < ld, ld > p1 = get(x, y), p2 = get(tx, ty);
    spang -= p1.y + p2.y;
    return p1.x + p2.x + spang * rr;
}

int main() {
    //Freopen;
    //ios::sync_with_stdio(false);
    cin >> xx >> yy >> vv;
    cin >> x >> y >> v >> rr;
    ld stang = atan2(yy, xx);
    ld rr = sqrt(sqr(xx) + sqr(yy));
    ld l = 0, r = 1e7;
    forn(u, 200) {
        ld m = (l + r) / 2;
        ld nang = stang + (vv * m / rr);
        tx = cos(nang) * rr;
        ty = sin(nang) * rr;
        if (check(m) <= v * m)
            r = m;
        else
            l = m;
    }
    cout.precision(16);
    cout << fixed << r;
    return 0;
}


/*




*/
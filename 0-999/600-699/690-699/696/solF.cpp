#include<bits/stdc++.h>
#define rep(i, j, k) for(int i = int(j); i < int(k); i++)
#define all(v) (v).begin(), (v).end()
#define pb push_back
#define ft first
#define sd second
typedef long long ll;
typedef unsigned long long ull;
typedef long double ld;
typedef unsigned int uint;
const long long INF = 4e18L + 1;
const int IINF = 2e9 + 1;
const int limit = 1048576;
using namespace std;

typedef double R;

typedef complex<R> com;

const R eps = 1e-9;
const R fail = 4.4242315e-12;
const com fail2(fail, fail);

bool eq(R r1, R r2) { return fabs(r1 - r2) < eps; }

bool eq(com c1, com c2) { return eq(c1.real(), c2.real()) and eq(c1.imag(), c2.imag()); }

// checklist
// long longi
// treść źle przeczytana
// przypadki brzegowe(min max wejście)

R dot(com c1, com c2) { return c1.real() * c2.real() + c1.imag() * c2.imag();}
R det(com c1, com c2) { return c1.real() * c2.imag() - c1.imag() * c2.real();}

struct line{
    com n;
    R c;
    line(com n1, R c1)
        :n{n1 / abs(n1)}, c{c1} {}

    line(com p1, com p2)
        :n{((p2 - p1) / abs(p2 - p1)) * com(0, 1)}, c{dot(p1, n)} {}

    com dir() const { return n * com(0, 1); }

    com val(R t) const { return c * n + t * dir(); }

    R intersect(const line& other) const {
        if(fabs(det(n, other.n)) < eps) {
            return fail;
        } else {
            return (other.c - c * dot(n, other.n)) / dot(dir(), other.n);
        }
    }
};

struct seg{
    line l;
    R start, end;
    
    seg(com p1, com p2)
        :l{p1, p2}, start{dot(l.dir(), p1)}, end{dot(l.dir(), p2)} {}

    com p1() const { return l.val(start); }
    com p2() const { return l.val(end); }

    R len() const { return fabs(end - start); }
};

R dist(const line& l, const com& p) { return fabs(dot(p, l.n) - l.c); }

com intersect(const line& a, const line& b) {
    R t = a.intersect(b);
    if(eq(t, fail)) {
        return fail2;
    } else {
        return a.val(t);
    }
}

com intersect(const seg& s, const line& l) {
    R t = s.l.intersect(l);
    if(eq(t, fail) or !(eq(fabs(t - s.start) + fabs(t - s.end), s.len()))) {
        return fail2;
    } else {
        return s.l.val(t);
    }
}

int n;
vector<com> polygon;

bool check_r( R radius, com& candB, com& candL) {
    int i1 = 2, m = 0;
    vector<int> best(n);
    vector<com> cand(n);
    rep(i, 0, n) {
        cand[i] = polygon[i + 1];
    }
    rep(i, 0, n) {
        i1 = max(i1, i + 2);
        m = max(m, i);
        best[i] = i1 - i;
        line l1(polygon[i], polygon[i + 1]);
        while(i1 != i + n) {
            line l2(polygon[i1], polygon[i1 + 1]);
            if(det(polygon[i + 1] - polygon[i], polygon[i1 + 1] - polygon[i1]) < eps) {
                break;
            }
            while(dist(l1, polygon[m + 1]) < dist(l2, polygon[m + 1])) {
                m++;
            }
            com low = polygon[m], high = polygon[m + 1];
            rep(_, 0, 42) {
                com mid = (low + high) / 2.0;
                if(dist(l1, mid) < dist(l2, mid)) {
                    low = mid;
                } else {
                    high = mid;
                }
            }
            if(dist(l1, low) < radius) {
                i1++;
                best[i] = i1 - i;
                cand[i] = low;
            } else {
                break;
            }
        }
    }
    rep(i, 0, n) {
        if(best[i] + best[(i + best[i]) % n] >= n) {
            candB = cand[i]; candL = cand[(i + best[i]) % n];
            return true;
        }
    }
    return false;
}

int main()
{
    ios_base::sync_with_stdio(0);
//     cin.tie(0);
    
    cin >> n;
    rep(_, 0, n) {
        R t1, t2; cin >> t1 >> t2;
        polygon.pb(com(t1, t2));
    }
    rep(i, n, 2 * n) {
        polygon.pb(polygon[i - n]);
    }
    R start = 0.0, fin = 100000;
    com Bar, Lya;
    rep(_, 0, 84) {
        R m = (start + fin) / 2.0;
        com candB, candL;
        if(check_r(m, candB, candL)) {
            fin = m;
            Bar = candB; Lya = candL;
        } else {
            start = m;
        }
        /* cout << endl << endl; */
    }
    cout << setprecision(10) << fixed << fin << "\n";
    cout << Bar.real() << " " << Bar.imag() << "\n";
    cout << Lya.real() << " " << Lya.imag() << "\n";
    return 0;
}
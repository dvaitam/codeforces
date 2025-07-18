#include <bits/stdc++.h>
 
#define fi first
#define se second
#define mp make_pair
#define mt make_tuple
#define pb push_back
#define INF  (1<<30)
#define INFL (1LL<<60)
#define EPS ((ld)(1e-9))
 
#define sz(x) ((int)(x).size())
#define setz(x) memset(x, 0, sizeof(x))
#define all(x) (x).begin(), (x).end()
#define rep(i, e) for (int i = 0, _##i = (e); i < _##i; i++)
#define repp(i, s, e) for (int i = (s), _##i = (e); i < _##i; i++)
#define repr(i, s, e) for (int i = (s)-1, _##i = (e); i >= _##i; i--)
#define repi(i, x) for (auto &i : (x))
#define ARI(...) vector<int>({__VA_ARGS__})
#define ARS(...) vector<string>({__VA_ARGS__})
 
 
using namespace std;
 
typedef long long ll;
typedef unsigned long long ull;
typedef long double ld;
typedef complex<double> base;
typedef pair<int, int> pii;
typedef pair<double, double> pdd;
typedef pair<ll, ll> pll;
 
template<typename T, typename V>
ostream &operator<<(ostream &os, const pair<T, V> pai) { 
    return os << '(' << pai.first << ' ' << pai.second << ')';
}
 
template<typename T>
ostream &operator<<(ostream &os, const vector<T> v) {
    cout << '[';
    for (auto p : v) cout << p << ",";
    cout << "]";
    return os;
}
 
template<typename T, typename V>
ostream &operator<<(ostream &os, const set<T, V> v) {
    cout << "{";
    for (auto p : v) cout << p << ",";
    cout << "}";
    return os;
}
 
template<typename T, typename V>
ostream &operator<<(ostream &os, const map<T, V> v) {
    cout << "{";
    for (auto p : v) cout << p << ",";
    cout << "}";
    return os;
}
 
#ifndef __SOULTCH
#define debug(...) 0
#define endl '\n'
#else
#define debug(...) cout << " [-] ", _dbg(#__VA_ARGS__, __VA_ARGS__)
template<class TH> void _dbg(const char *sdbg, TH h){ cout << sdbg << '=' << h << endl; }
template<class TH, class... TA> void _dbg(const char *sdbg, TH h, TA... a) {
    while(*sdbg != ',') cout << *sdbg++;
    cout << '=' << (h) << ','; 
    _dbg(sdbg+1, a...);
}
#endif
 
template<typename T> void get_max(T &a, T b) {a = max(a, b);}
template<typename T> void get_min(T &a, T b) {a = min(a, b);}

namespace SEG {
    const int N = 1<<19;
    
    struct node {
        ll pt, x, maxx, tot;
    } tree[N*2];

    ostream& operator<< (ostream& os, node a) {
        os << a.pt << ' ' << a.x << ' ' << a.maxx << ' ' << a.tot;
        return os;
    }

    node operator+(node a, node b) {
        if (a.tot < b.tot+a.pt) return {a.pt+b.pt, b.x, b.maxx, max(a.tot, b.tot+a.pt)};
        else                    return {a.pt+b.pt, b.x, a.maxx, max(a.tot, b.tot+a.pt)};
    }
    
    void set_v(int p, int v, int x) {
        if (tree[p+N].x != x) {
            tree[p+N].x = x;
            tree[p+N].maxx = x;
            tree[p+N].tot -= x;
        }
        tree[p+N].pt += v;
        tree[p+N].tot += v;
    }

    void build() {
        repr(i, N, 0) tree[i] = tree[i<<1]+tree[i<<1|1];
        rep(i, 6) debug(tree[N+i]);
    }

    void update(int p, int v) {
        p += N;
        tree[p].pt -= v;
        tree[p].tot -= v;

        for (p >>= 1; p; p >>= 1) tree[p] = tree[p<<1]+tree[p<<1|1];
    }

    node query(int l, int r) {
        l += N, r += N;
        node lsum = {0, 0, 0, -INFL};
        node rsum = {0, INFL, 0, -INFL};
        while(l <= r) {
            if (l&1) lsum = lsum+tree[l++];
            if (~r&1) rsum = tree[r--]+rsum;
            l >>= 1;
            r >>= 1;
        }
        return lsum+rsum;
    }
}

int N;
vector<int> X;
struct point {
    int a, b, c;
} A[500000];

int main(void) {
    iostream::sync_with_stdio(false);
    cin.tie(nullptr);
    
    cin >> N;
    rep(i, N) {
        int a, b, c; cin >> a >> b >> c;
        if (a > b) swap(a, b);
        A[i] = {a, b, c};
        debug(a, b, c);
        X.push_back(b);
    }

    sort(all(X));
    X.erase(unique(all(X)), X.end());

    rep(i, N) A[i].b = lower_bound(all(X), A[i].b)-X.begin();
    rep(i, N) SEG::set_v(A[i].b, A[i].c, X[A[i].b]);
    SEG::build();

    sort(A, A+N, [](auto a, auto b) {
        return a.a < b.a;
    });

    pair<ll, pii> res = {0, {1e9+1, 1e9+1}};
    rep(i, N) {
        while(i > 0 and i < N and A[i-1].a == A[i].a) {
            SEG::update(A[i].b, A[i].c);
            i++;
        }
        if (i == N) break;
        
        int p = lower_bound(all(X), A[i].a)-X.begin();
        auto r = SEG::query(p, sz(X)-1);
        debug(r);
        r.tot += A[i].a;
        if (res.fi < r.tot) res = {r.tot, {A[i].a, r.maxx}};
        SEG::update(A[i].b, A[i].c);
    }

    cout << res.fi << endl;
    cout << res.se.fi << ' ' << res.se.fi << ' ' << res.se.se << ' ' << res.se.se << endl;
}
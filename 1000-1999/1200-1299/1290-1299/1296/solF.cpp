#include <iostream>

#include <vector>

#include <bitset>

#include <map>

#include <queue>

#include <stack>

#include <algorithm>

#include <cmath>

#include <set>

#include <cstring>

#include <array>

#pragma GCC optimize(2)

#include <bits/stdc++.h>

#define rep(i,from,to) for(int i=from;i<to;i++)

#define ite2(x,y,arr) for(auto [x,y]:arr)

#define pdd pair<double, double>

#define ite(i,arr) for(auto &i:arr)

#define MID(l,r) int mid=(l+r)>>1

#define ALL(arr) arr.begin(),arr.end()

#define AXY(a,x,y) int x=a.first,y=a.second

#define vc vector

#define vi vector<int>

#define vll vector<long long>

#define pii pair<int,int>



typedef long long ll;

typedef unsigned long long ull;

using namespace std;

namespace dbg {

#ifdef FTY

    template <typename T>

    void __print_var(string_view name, const T& x) { std::cerr << name << " = " << x; }

    template <typename T>

    void __print_var(string_view name, const vector<T>& x) {

        std::cerr << name << " = ";

        bool is_first = true;

        for (auto& ele : x) std::cerr << (is_first ? (is_first = false, "[") : ", ") << ele;

        std::cerr << "]";

    }

    template <typename T>

    void __print_var(string_view name, const set<T>& x) {

        std::cerr << name << " = ";

        bool is_first = true;

        for (auto& ele : x) std::cerr << (is_first ? (is_first = false, "{") : ", ") << ele;

        std::cerr << "}";

    }

    template <typename K, typename V>

    void __print_var(string_view name, const map<K, V>& x) {

        std::cerr << name << " = ";

        bool is_first = true;

        for (auto& [k, v] : x) std::cerr << (is_first ? (is_first = false, "{") : ", ") << "(" << k << ": " << v << ")";

        std::cerr << "}";

    }

    template <typename T>

    void __log(string_view name, const T& x) {

        __print_var(name, x); std::cerr << endl;

    }

#define LOG(args)\

    { std::cerr << "line " << __LINE__ << ": " << __func__ << "(): ";\

    __log(#args, ##args); }

#else

#define LOG(...)

#endif

}

//#define int ll

const double eps = 1e-8;

const double PI = acos(-1);

struct Point {

    double x, y;

    Point() {

        x = 0, y = 0;

    }

    Point(double x, double y) {

        this->x = x;

        this->y = y;

    }

};

Point operator +(Point a, Point b) {

    return Point(a.x + b.x, a.y + b.y);

}

Point operator -(Point a, Point b) {

    return Point(a.x - b.x, a.y - b.y);

}

Point operator *(double a, Point b) {

    return Point(a * b.x, a * b.y);

}

Point operator *(Point b, double a) {

    return Point(a * b.x, a * b.y);

}

Point operator /(Point b, double a) {

    return Point(b.x / a, b.y / a);

}

double len(Point a) {

    return sqrt(a.x * a.x + a.y * a.y);

}

double dis(Point a, Point b) {

    return len(a - b);

}

bool operator ==(Point a, Point b) {

    return dis(a, b) <= eps;

}

bool operator !=(Point a, Point b) {

    return !(a == b);

}

double operator *(Point a, Point b) {

    return a.x * b.x + a.y * b.y;

}



double operator ^(Point a, Point b) {

    return a.x * b.y - a.y * b.x;

}



double getAngel(double b, double a, double c) {

    return acos((a * a + c * c - b * b) / (2 * a * c));

}

double getAngel(Point a, Point b) {

    return acos(a * b / len(a) / len(b));

}







int mod = 1e9 + 7;

using namespace dbg;

template <class T>

T __gcd(T a, T b) {

    if (a < b) swap(a, b);

    return b ? __gcd(b, a % b) : a;

}

template <class T>

T __lcm(T a, T b) {

    T num = __gcd(a, b);

    return a / num * b;

}

inline int lowbit(int num) { return num & (-num); }

inline int qmi(int a, int b) {

    a %= mod;

    ll res = 1;

    while (b) {

        if (b & 1) res = (ll)res * a % mod;

        a = (ll)a * a % mod;

        b >>= 1;

    }

    return res;

}

int inv(int num) {

    return qmi(num, mod - 2);

}

const int N = 5e5 + 10;

ll fact[(int)N + 5];

ll inv_fact[(int)N + 5];

int prime[(int)N + 5];

int valid[(int)N + 5];

int pn = 0;

inline void getPrime() {

    rep(i, 2, N + 1) {

        if (!valid[i]) {

            valid[i] = i;

            prime[pn++] = i;

        }

        for (int j = 0; j < pn && i * prime[j] <= N; j++) {

            valid[i * prime[j]] = prime[j];

            if (i % prime[j] == 0) break;

        }

    }

}

inline void getFact() {

    fact[0] = fact[1] = 1;

    for (int i = 2; i <= N; i++) {

        fact[i] = i * fact[i - 1] % mod;

    }

}

inline void getInv() {

    inv_fact[N] = inv(fact[N]);

    for (int i = N - 1; i >= 0; i--) {

        inv_fact[i] = inv_fact[i + 1] * (i + 1) % mod;

    }

}

inline ll CC(int n, int m) {

    if (m > n) {

        return 0;

    }//不输出0是为了好debug  此处会有RE ！！！

    if (m < 0) {

        return 1;

    }

    ll res = fact[n] * inv_fact[m] % mod * inv_fact[n - m] % mod;

    return res;

}



int n;

vc<vi> edge;

vi fa;

vi h;

vi value;

map<pii,pii> res;

void dfs(int p, int father,int height) {

    fa[p] = father;

    h[p] = height;

    ite(u, edge[p]) {

        if (u != father) {

            dfs(u, p,height+1);

        }

    }

}







bool cherk(int u, int v, int w) {

    if (h[u] < h[v]) {

        swap(u, v);

    }

    

    while (h[u] > h[v]) {

        if (value[u] == w) return true;

        u = fa[u];

    }

    while (u != v) {

        if (value[u] == w) return true;

        u = fa[u];

        if (value[v] == w) return true;

        v = fa[v];

    }

    return false;

}

void func(int u, int v, int w) {

    if (h[u] < h[v]) {

        swap(u, v);

    }

    while (h[u] > h[v]) {

        value[u] = max(value[u], w);

        u = fa[u];

    }

    while (u != v) {

        value[u] = max(value[u], w);

        u = fa[u];

        value[v] = max(value[v], w);

        v = fa[v];

    }



}





signed main() {

    ios::sync_with_stdio(false);

    cin.tie(0);

    cout.tie(0);



    int q = 1;

    //cin >> q;

    while (q--) {

        int n;

        cin >> n;

        edge.resize(n);

        fa.resize(n,-1);

        h.resize(n);

        value.resize(n, 1);

        rep(i, 0, n - 1) {

            int u, v;

            cin >> u >> v;

            u--,v--;

            edge[u].push_back(v);

            edge[v].push_back(u);

            res[{min(u, v), max(u, v) }] = { 1,i };

        }

        dfs(0, -1,1);



        int m;

        cin >> m;

        queue<array<int, 3>> que;

        rep(i, 0, m) {

            int u, v;

            cin >> u >> v;

            u--, v--;

            int w;

            cin >> w;

            func(u, v, w);

            array<int, 3> tmp = { u,v,w };

            que.push(tmp);

        }

        int flag = 1;

        while (que.size()) {

            auto cur = que.front();

            que.pop();

            int u = cur[0], v = cur[1], w = cur[2];

            if (!cherk(u, v, w)) {

                flag = 0;

                break;

            }

        }

        rep(i, 1, n) {

            res[{min(i, fa[i]), max(i, fa[i])}].first = value[i];

        }

        if (flag) {

            vi ans(n - 1);

            ite2(__, v, res) {

                ans[v.second] = v.first;

            }

            ite(i, ans) cout << i << ' ';

            cout << endl;

        }

        else {

            cout << -1 << endl;

        }

        

    }

    return 0;

}
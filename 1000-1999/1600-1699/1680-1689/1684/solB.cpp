#pragma GCC optimize("Ofast")
#pragma GCC optimization ("unroll-loops")
#pragma GCC optimization ("-fsanitize=address")
#include <bits/stdc++.h>
using namespace std;
/*_____________________Containers_____________________*/
typedef long long int ll;
typedef unsigned long long int ull;
typedef long double ld;
typedef pair<int, int> PII;
typedef array<int, 3> ARR;
typedef vector<int> VI;
typedef vector<double> VDB;
typedef vector<string> VS;
typedef vector<PII> VPII;
typedef vector<ARR> VARR;
typedef vector<VI> VVI;
typedef vector<VDB> VVDB;
typedef vector<VPII> VVPII;
typedef unordered_map<int,int> UMPII;
typedef map<int,int> MPII;
typedef set<int> SETI;
typedef stack<int> STI;
typedef queue<int> QI;
typedef multiset<int> MSETI;
typedef pair<long, long> PLL;
typedef vector<long> VL;
typedef vector<PLL> VPLL;
typedef vector<VL> VVL;
typedef vector<VARR> VVARR;
typedef map<long, long> MPLL;
typedef set<long> SETL;
typedef queue<long> QL;
typedef multiset<long> MSETL;
typedef pair<ll, ll> PLLLL;
typedef vector<ll> VLL;
typedef vector<PLLLL> VPLLLL;
typedef vector<VLL> VVLL;
typedef map<ll, ll> MPLLLL;
typedef set<ll> SETLL;
typedef queue<ll> QLL;
typedef multiset<ll> MSETLL;
/*__________________________________________________________*/
/*____________________Limits & Constants____________________*/
constexpr int inf = 0x3f3f3f3f;
constexpr ll INF = 0x3F3F3F3F3F3F3F3F;
/*_______________________________________________________________*/
/*_____________________________Macros____________________________*/
#define MOD int(998244353)
#define PI (2 * acos(0.0))
#define ff(i, j, k) for(ll i = j; i < k; i++)
#define FF(i, j, k) for(ll i = j; i >= k; i--)
#define mem(a,n) memset(a, n, sizeof(a) )
#define gcd(a,b) __gcd(a,b)
#define lcm(a,b) (a*(b/gcd(a,b)))
#define spc(n) setprecision(n)
#define sz(x) int(x.size())
#define pb push_back
#define eb emplace_back
#define mk(x,y) make_pair(x,y)
#define fi first
#define se second
#define bit_popcount __builtin_popcount
#define bit_ctz __builtin_ctz
#define bit_clz __builtin_clz
#define bit_popcountll __builtin_popcountll
#define bit_ctzll __builtin_ctzll
#define bit_clzll __builtin_clzll
#define csor(a,compare) sort(a.begin(),a.end(),compare)
#define asor(a) sort(a.begin(),a.end())
#define dsor(a) sort(a.rbegin(),a.rend())
#define mxe(v) *max_element(v.begin(), v.end())
int ask(const string &s) { int r; return cout << "? " << s <<endl, cin >> r, r; }
/*________________________________________________________________*/
/*___________________________printf out___________________________*/
#ifndef ONLINE_JUDGE
template<class T> void _R(T &x) { cin >> x; }
void _R(int &x) { (void)scanf("%d", &x); }
void _R(ll &x) { (void)scanf("%lld", &x); }
void _R(double &x) { (void)scanf("%lf", &x); }
void _R(char &x) { (void)scanf(" %c", &x); }
void _R(char *x) { (void)scanf("%s", x); }
void R() {}
template<class T, class... U> void R(T &head, U &... tail) { _R(head); R(tail...); }
template<class T> void _W(const T &x) { cout << x; }
void _W(const int &x) { printf("%d", x); }
void _W(const ll &x) { printf("%lld", x); }
void _W(const double &x) { printf("%.16f", x); }
void _W(const char &x) { putchar(x); }
void _W(const char *x) { printf("%s", x); }
template<class T,class U> void _W(const pair<T,U> &x) {_W(x.fi); putchar(' '); _W(x.se); putchar('\n'); }
template<class T> void _W(const vector<T> &x) { for (auto i = x.begin(); i != x.end(); _W(*i++)) putchar(" \n"[i == x.cbegin()]); }
void W() {}
template<class T, class... U> void W(const T &head, const U &... tail) { _W(head); putchar(sizeof...(tail) ? ' ' : '\n'); W(tail...); }
#endif
/*_______________________________________________________________*/
/*_________________Combinations number___________________________*/
const int Combinations_N = 400005;
ll fact[Combinations_N];
ll mult(ll a,ll b) {
    return ((a % MOD)*(b % MOD)) % MOD;
}
ll add(ll a,ll b) {
    return ((a % MOD)+(b % MOD)) % MOD;
}
ll subt(ll a,ll b) {
    ll ans = ((a % MOD) - (b % MOD)) % MOD;
    ans %= MOD;
    ans = (ans + MOD) % MOD;
    return ans;
}

ll exp(ll a, ll b) {
    ll ans = 1;
    ll po = a;
    while (b != 0){
        if(b % 2){
            ans = ((ans % MOD) * (po % MOD)) % MOD;
        }
        po = ((po % MOD) * (po % MOD)) % MOD;
        b /= 2;
    }
    return ans;
}        
ll inv(ll x) {
	return exp(x, MOD - 2); // 公式：1 / B = B ^ (MOD - 2)
}
ll divide(ll x, ll y) {
	return mult(x, inv(y));
}
void factorial(void) {
	fact[0] = 1;
	for(int i = 1; i < Combinations_N; i++)
		fact[i] = mult(i, fact[i - 1]);
}
// Cm n = C(n-m) n
// Cm (n+1) = Cm n + C(m-1) n
// Σk i=0 Cn (n+i) = C(n+1) (n+k+1)
ll C(ll n, ll k) {
	if(k > n) return 0;
	return divide(fact[n], mult(fact[n - k], fact[k]));
}

/*_______________________________________________________________*/
/*______________________Carmichael Numbers_______________________*/
ll Pow(ll a, ll n) {
    if (n == 0) return 0;
    if (n == 1) return a;
    ll c = Pow(a, n / 2) % MOD;
    if (n % 2 == 0) return (c % MOD) * (c % MOD) % MOD;
    return (c % MOD) * (c % MOD) * (a % MOD)% MOD;
}

/*_______________________________________________________________*/
/*___________________________Modular_____________________________*/
template <typename T>
T inverse(T a, T m) {
    T u = 0, v = 1;
    while (a != 0) {
        T t = m / a;
        m -= t * a; swap(a, m);
        u -= t * v; swap(u, v);
    }
    assert(m == 1);
    return u;
}
 
template <typename T>
class Modular {
    public:
    using Type = typename decay<decltype(T::value)>::type;
    
    constexpr Modular() : value() {}
    template <typename U>
    Modular(const U& x) {
        value = normalize(x);
    }
    
    template <typename U>
    static Type normalize(const U& x) {
        Type v;
        if (-mod() <= x && x < mod()) v = static_cast<Type>(x);
        else v = static_cast<Type>(x % mod());
        if (v < 0) v += mod();
        return v;
    }
    
    const Type& operator()() const { return value; }
    template <typename U>
    explicit operator U() const { return static_cast<U>(value); }
    constexpr static Type mod() { return T::value; }
    
    Modular& operator+=(const Modular& other) { if ((value += other.value) >= mod()) value -= mod(); return *this; }
    Modular& operator-=(const Modular& other) { if ((value -= other.value) < 0) value += mod(); return *this; }
    template <typename U> Modular& operator+=(const U& other) { return *this += Modular(other); }
    template <typename U> Modular& operator-=(const U& other) { return *this -= Modular(other); }
    Modular& operator++() { return *this += 1; }
    Modular& operator--() { return *this -= 1; }
    Modular operator++(int) { Modular result(*this); *this += 1; return result; }
    Modular operator--(int) { Modular result(*this); *this -= 1; return result; }
    Modular operator-() const { return Modular(-value); }
    
    template <typename U = T>
    typename enable_if<is_same<typename Modular<U>::Type, int>::value, Modular>::type& operator*=(const Modular& rhs) {
#ifdef _WIN32
    uint64_t x = static_cast<int64_t>(value) * static_cast<int64_t>(rhs.value);
    uint32_t xh = static_cast<uint32_t>(x >> 32), xl = static_cast<uint32_t>(x), d, m;
    asm(
      "divl %4; \n\t"
      : "=a" (d), "=d" (m)
      : "d" (xh), "a" (xl), "r" (mod())
    );
    value = m;
#else
    value = normalize(static_cast<int64_t>(value) * static_cast<int64_t>(rhs.value));
#endif
    return *this;
    }
    template <typename U = T>
    typename enable_if<is_same<typename Modular<U>::Type, long long>::value, Modular>::type& operator*=(const Modular& rhs) {
        long long q = static_cast<long long>(static_cast<long double>(value) * rhs.value / mod());
        value = normalize(value * rhs.value - q * mod());
        return *this;
    }
    template <typename U = T>
    typename enable_if<!is_integral<typename Modular<U>::Type>::value, Modular>::type& operator*=(const Modular& rhs) {
        value = normalize(value * rhs.value);
        return *this;
    }
    
    Modular& operator/=(const Modular& other) { return *this *= Modular(inverse(other.value, mod())); }
    
    friend const Type& abs(const Modular& x) { return x.value; }
    
    template <typename U>
    friend bool operator==(const Modular<U>& lhs, const Modular<U>& rhs);
    
    template <typename U>
    friend bool operator<(const Modular<U>& lhs, const Modular<U>& rhs);
    
    template <typename V, typename U>
    friend V& operator>>(V& stream, Modular<U>& number);
    
private:
    Type value;
};
 
template <typename T> bool operator==(const Modular<T>& lhs, const Modular<T>& rhs) { return lhs.value == rhs.value; }
template <typename T, typename U> bool operator==(const Modular<T>& lhs, U rhs) { return lhs == Modular<T>(rhs); }
template <typename T, typename U> bool operator==(U lhs, const Modular<T>& rhs) { return Modular<T>(lhs) == rhs; }
 
template <typename T> bool operator!=(const Modular<T>& lhs, const Modular<T>& rhs) { return !(lhs == rhs); }
template <typename T, typename U> bool operator!=(const Modular<T>& lhs, U rhs) { return !(lhs == rhs); }
template <typename T, typename U> bool operator!=(U lhs, const Modular<T>& rhs) { return !(lhs == rhs); }
 
template <typename T> bool operator<(const Modular<T>& lhs, const Modular<T>& rhs) { return lhs.value < rhs.value; }
 
template <typename T> Modular<T> operator+(const Modular<T>& lhs, const Modular<T>& rhs) { return Modular<T>(lhs) += rhs; }
template <typename T, typename U> Modular<T> operator+(const Modular<T>& lhs, U rhs) { return Modular<T>(lhs) += rhs; }
template <typename T, typename U> Modular<T> operator+(U lhs, const Modular<T>& rhs) { return Modular<T>(lhs) += rhs; }
 
template <typename T> Modular<T> operator-(const Modular<T>& lhs, const Modular<T>& rhs) { return Modular<T>(lhs) -= rhs; }
template <typename T, typename U> Modular<T> operator-(const Modular<T>& lhs, U rhs) { return Modular<T>(lhs) -= rhs; }
template <typename T, typename U> Modular<T> operator-(U lhs, const Modular<T>& rhs) { return Modular<T>(lhs) -= rhs; }
 
template <typename T> Modular<T> operator*(const Modular<T>& lhs, const Modular<T>& rhs) { return Modular<T>(lhs) *= rhs; }
template <typename T, typename U> Modular<T> operator*(const Modular<T>& lhs, U rhs) { return Modular<T>(lhs) *= rhs; }
template <typename T, typename U> Modular<T> operator*(U lhs, const Modular<T>& rhs) { return Modular<T>(lhs) *= rhs; }
 
template <typename T> Modular<T> operator/(const Modular<T>& lhs, const Modular<T>& rhs) { return Modular<T>(lhs) /= rhs; }
template <typename T, typename U> Modular<T> operator/(const Modular<T>& lhs, U rhs) { return Modular<T>(lhs) /= rhs; }
template <typename T, typename U> Modular<T> operator/(U lhs, const Modular<T>& rhs) { return Modular<T>(lhs) /= rhs; }
 
template<typename T, typename U>
Modular<T> power(const Modular<T>& a, const U& b) {
    assert(b >= 0);
    Modular<T> x = a, res = 1;
    U p = b;
    while (p > 0) {
        if (p & 1) res *= x;
        x *= x;
        p >>= 1;
    }
    return res;
}
 
template <typename T>
bool IsZero(const Modular<T>& number) {
    return number() == 0;
}
 
template <typename T>
string to_string(const Modular<T>& number) {
    return to_string(number());
}
 
// U == std::ostream? but done this way because of fastoutput
template <typename U, typename T>
U& operator<<(U& stream, const Modular<T>& number) {
    return stream << number();
}
 
// U == std::istream? but done this way because of fastinput
template <typename U, typename T>
U& operator>>(U& stream, Modular<T>& number) {
    typename common_type<typename Modular<T>::Type, long long>::type x;
    stream >> x;
    number.value = Modular<T>::normalize(x);
    return stream;
}
 
/*
using ModType = int;
 
struct VarMod { static ModType value; };
ModType VarMod::value;
ModType& MOD = VarMod::value;
using Mint = Modular<VarMod>;
*/

using Mint = Modular<std::integral_constant<decay<decltype(MOD)>::type, MOD>>;
typedef vector<Mint> VMint;
typedef vector<VMint> VVMint;

/*vector<Mint> fact(1, 1);
vector<Mint> inv_fact(1, 1);

Mint C(int n, int k) {
    if (k < 0 || k > n) {
        return 0;
    }
    while ((int) fact.size() < n + 1) {
        fact.push_back(fact.back() * (int) fact.size());
        inv_fact.push_back(1 / fact.back());
    }
    return fact[n] * inv_fact[k] * inv_fact[n - k];
}*/

/*_______________________________________________________________*/
/*_________________________segment tree__________________________*/
/*______________________权值线段树可参考1705E_____________________*/
class segtree {
    public:
    struct node {
        ll mn = -1;
        ll mx = -1;
        // ll val = 0;
        // ll sum_v = 0;
        ll add = 0;
        ll sum_a = 0;
        ll idx;
 
        void apply(int l, int r, ll v) {
            mn = mx = v;
            // val = v;
            // sum_v = (ll)v * (r - l + 1);
            add += v;
            sum_a += (ll)v * (r - l + 1);
            idx = l;
        }
    };

    node unite(const node &a, const node &b) const {
        node res;
        res.mn = min(a.mn, b.mn);
        res.mx = max(a.mx, b.mx);
        // res.sum_v = a.sum_v + b.sum_v;
        res.sum_a = a.sum_a + b.sum_a;
        if (a.mx > b.mx) res.idx = a.idx;
        else res.idx = b.idx;
        return res;
    }
 
    inline void push(int x, int l, int r) {
        int y = (l + r) >> 1;
        int z = x + ((y - l + 1) << 1);
        // if (tree[x].val != 0) {
        //     tree[x + 1].apply(l, y, tree[x].val);
        //     tree[z].apply(y + 1, r, tree[x].val);
        //     tree[x].val = 0;
        // }
        if (tree[x].add != 0) {
            tree[x + 1].apply(l, y, tree[x].add);
            tree[z].apply(y + 1, r, tree[x].add);
            tree[x].add = 0;
        }
    }
    
    inline void pull(int x, int z) {
        tree[x] = unite(tree[x + 1], tree[z]);
    }
    
    int n;
    vector<node> tree;
    
    void build(int x, int l, int r) {
        if (l == r) {
            return;
        }
        int y = (l + r) >> 1;
        int z = x + ((y - l + 1) << 1);
        build(x + 1, l, y);
        build(z, y + 1, r);
        pull(x, z);
    }
    
    template <typename M>
    void build(int x, int l, int r, const vector<M> &v) {
        if (l == r) {
            tree[x].apply(l, r, v[l]);
            return;
        }
        int y = (l + r) >> 1;
        int z = x + ((y - l + 1) << 1);
        build(x + 1, l, y, v);
        build(z, y + 1, r, v);
        pull(x, z);
    }
    
    node get(int x, int l, int r, int ll, int rr) {
        if (ll <= l && r <= rr) {
            return tree[x];
        }
        int y = (l + r) >> 1;
        int z = x + ((y - l + 1) << 1);
        push(x, l, r);
        node res{};
        if (rr <= y) {
            res = get(x + 1, l, y, ll, rr);
        } else {
            if (ll > y) {
                res = get(z, y + 1, r, ll, rr);
            } else {
                res = unite(get(x + 1, l, y, ll, rr), get(z, y + 1, r, ll, rr));
            }
        }
        pull(x, z);
        return res;
    }
    
    template <typename... M>
    void modify(int x, int l, int r, int ll, int rr, const M&... v) {
        if (ll <= l && r <= rr) {
            tree[x].apply(l, r, v...);
            return;
        }
        int y = (l + r) >> 1;
        int z = x + ((y - l + 1) << 1);
        push(x, l, r);
        if (ll <= y) {
            modify(x + 1, l, y, ll, rr, v...);
        }
        if (rr > y) {
            modify(z, y + 1, r, ll, rr, v...);
        }
        pull(x, z);
    }
    
    int find_first_knowingly(int x, int l, int r, const function<bool(const node&)> &f) {
        if (l == r) {
            return l;
        }
        push(x, l, r);
        int y = (l + r) >> 1;
        int z = x + ((y - l + 1) << 1);
        int res;
        if (f(tree[x + 1])) {
            res = find_first_knowingly(x + 1, l, y, f);
        } else {
            res = find_first_knowingly(z, y + 1, r, f);
        }
        pull(x, z);
        return res;
    }
    
    int find_first(int x, int l, int r, int ll, int rr, const function<bool(const node&)> &f) {
        if (ll <= l && r <= rr) {
            if (!f(tree[x])) {
                return -1;
            }
            return find_first_knowingly(x, l, r, f);
        }
        push(x, l, r);
        int y = (l + r) >> 1;
        int z = x + ((y - l + 1) << 1);
        int res = -1;
        if (ll <= y) {
            res = find_first(x + 1, l, y, ll, rr, f);
        }
        if (rr > y && res == -1) {
            res = find_first(z, y + 1, r, ll, rr, f);
        }
        pull(x, z);
        return res;
    }
    
    int find_last_knowingly(int x, int l, int r, const function<bool(const node&)> &f) {
        if (l == r) {
            return l;
        }
        push(x, l, r);
        int y = (l + r) >> 1;
        int z = x + ((y - l + 1) << 1);
        int res;
        if (f(tree[z])) {
            res = find_last_knowingly(z, y + 1, r, f);
        } else {
            res = find_last_knowingly(x + 1, l, y, f);
        }
        pull(x, z);
        return res;
    }
    
    int find_last(int x, int l, int r, int ll, int rr, const function<bool(const node&)> &f) {
        if (ll <= l && r <= rr) {
            if (!f(tree[x])) {
                return -1;
            }
            return find_last_knowingly(x, l, r, f);
        }
        push(x, l, r);
        int y = (l + r) >> 1;
        int z = x + ((y - l + 1) << 1);
        int res = -1;
        if (rr > y) {
            res = find_last(z, y + 1, r, ll, rr, f);
        }
        if (ll <= y && res == -1) {
            res = find_last(x + 1, l, y, ll, rr, f);
        }
        pull(x, z);
        return res;
    }
    
    segtree(int _n) : n(_n) {
        assert(n > 0);
        tree.resize(2 * n - 1);
        build(0, 0, n - 1);
    }
    
    template <typename M>
    segtree(const vector<M> &v) {
        n = v.size();
        assert(n > 0);
        tree.resize(2 * n - 1);
        build(0, 0, n - 1, v);
    }
    
    node get(int ll, int rr) {
        assert(0 <= ll);
        assert(ll <= rr);
        assert(rr <= n - 1);
        return get(0, 0, n - 1, ll, rr);
    }
    
    node get(int p) {
        assert(0 <= p && p <= n - 1);
        return get(0, 0, n - 1, p, p);
    }
    
    template <typename... M>
    void modify(int ll, int rr, const M&... v) {
        assert(0 <= ll);
        assert(ll <= rr);
        assert(rr <= n - 1);
        modify(0, 0, n - 1, ll, rr, v...);
    }
    
    // find_first and find_last call all FALSE elements
    // to the left (right) of the sought position exactly once
    
    int find_first(int ll, int rr, const function<bool(const node&)> &f) {
        assert(0 <= ll);
        assert(ll <= rr);
        assert(rr <= n - 1);
        return find_first(0, 0, n - 1, ll, rr, f);
    }
    
    int find_last(int ll, int rr, const function<bool(const node&)> &f) {
        assert(0 <= ll);
        assert(ll <= rr);
        assert(rr <= n - 1);
        return find_last(0, 0, n - 1, ll, rr, f);
    }
};

/*_______________________________________________________________*/
/*______________________Disjoint Set Union_______________________*/
class dsu {
    public:
    VI p;
    VLL ml;
    VLL mr;
    VLL mx;
    VLL s;
    int n;

    dsu (int _n) : n(_n) {
        assert(n > 0);
        p.resize(n + 1);
        iota(p.begin(), p.end(), 0);
        ml.resize(n + 1);
        iota(ml.begin(), ml.end(), 0);
        mr.resize(n + 1);
        iota(mr.begin(), mr.end(), 0);
        mx.resize(n + 1);
        iota(mx.begin(), mx.end(), 0);
        s.resize(n + 1);
        iota(s.begin(), s.end(), 0);
    }

    template <typename M>
    dsu (const vector<M> &v) {
        n = v.size();
        assert(n > 0);
        p.resize(n + 1);
        iota(p.begin(), p.end(), 0);
        ml.resize(n + 1);
        mr.resize(n + 1);
        mx.resize(n + 1);
        s.resize(n + 1);
        ff (i, 1, n + 1) {
            ml[i] = mr[i] = mx[i] = s[i] = v[i];
        }
    }
    
    inline int find (int x) {
        return (x == p[x] ? x : (p[x] = find(p[x])));
    }
    
    inline bool unite (int x, int y) {
        x = find(x);
        y = find(y);
        if (x == y) return false;
        p[x] = y;
        mx[y] = max({mx[x], mx[y], max(0LL, mr[x]) + ml[y], mr[x] + max(0LL, ml[y])});
        mr[y] = max(mr[y], max(0LL, mr[x]) + s[y]);
        ml[y] = max(ml[x], max(0LL, ml[y]) + s[x]);
        s[y] += s[x];
        return true;
    }
};

/*_______________________________________________________________*/
/*_________________________Sparse Table__________________________*/
/*__________________________可参考1691D___________________________*/
template <typename T, class F = function<T(const T&, const T&)>>
class SparseTable {
    public:
    int n;
    vector<vector<T>> mat;
    F func;
    
    SparseTable(const vector<T>& a, const F& f) : func(f) {
        n = static_cast<int>(a.size());
        int max_log = 32 - __builtin_clz(n);
        mat.resize(max_log);
        mat[0] = a;
        for (int j = 1; j < max_log; j++) {
            mat[j].resize(n - (1 << j) + 1);
            for (int i = 0; i <= n - (1 << j); i++) {
                mat[j][i] = func(mat[j - 1][i], mat[j - 1][i + (1 << (j - 1))]);
            }
        }
    }
    
    T get(int from, int to) const {
        assert(0 <= from && from <= to && to <= n - 1);
        int lg = 32 - __builtin_clz(to - from + 1) - 1;
        return func(mat[lg][from], mat[lg][to - (1 << lg) + 1]);
    }
};

/*_______________________________________________________________*/
/*____________________________Graph______________________________*/
template <typename T>
class graph {
    public:
    struct edge {
        int from;
        int to;
        T cost;
    };
    
    vector<edge> edges;
    vector<vector<int>> g;
    int n;
    
    graph(int _n) : n(_n) {
        g.resize(n);
    }
    
    virtual int add(int from, int to, T cost) = 0;
};

template <typename T>
class digraph : public graph<T> {
    public:
    using graph<T>::edges;
    using graph<T>::g;
    using graph<T>::n;

    digraph(int _n) : graph<T>(_n) {
    }

    int add(int from, int to, T cost = 1) {
        assert(0 <= from && from < n && 0 <= to && to < n);
        int id = (int) edges.size();
        g[from].push_back(id);
        edges.push_back({from, to, cost});
        return id;
    }

    digraph<T> reverse() const {
        digraph<T> rev(n);
        for (auto &e : edges) {
        rev.add(e.to, e.from, e.cost);
        }
        return rev;
    }
};
    
template <typename T>
class forest : public graph<T> {
    public:
    using graph<T>::edges;
    using graph<T>::g;
    using graph<T>::n;
    
    forest(int _n) : graph<T>(_n) {
    }
    
    int add(int from, int to, T cost = 1) {
        assert(0 <= from && from < n && 0 <= to && to < n);
        int id = (int)edges.size();
        assert(id < n - 1);
        g[from].push_back(id);
        g[to].push_back(id);
        edges.push_back({from, to, cost});
        return id;
    }
};
    
template <typename T>
class dfs_forest : public forest<T> {
    public:
    using forest<T>::edges;
    using forest<T>::g;
    using forest<T>::n;
    
    vector<int> pv;
    vector<int> pe;
    vector<int> order;
    vector<int> pos;
    vector<int> end;
    vector<int> sz;
    vector<int> root;
    vector<int> depth;
    vector<T> dist;
    
    dfs_forest(int _n) : forest<T>(_n) {
    }
    
    void init() {
        pv = vector<int>(n, -1);
        pe = vector<int>(n, -1);
        order.clear();
        pos = vector<int>(n, -1);
        end = vector<int>(n, -1);
        sz = vector<int>(n, 0);
        root = vector<int>(n, -1);
        depth = vector<int>(n, -1);
        dist = vector<T>(n);
    }
    
    void clear() {
        pv.clear();
        pe.clear();
        order.clear();
        pos.clear();
        end.clear();
        sz.clear();
        root.clear();
        depth.clear();
        dist.clear();
    }

    private:
    void do_dfs(int v) {
        pos[v] = (int) order.size();
        order.push_back(v);
        sz[v] = 1;
        for (int id : g[v]) {
            if (id == pe[v]) {
                continue;
            }
            auto &e = edges[id];
            int to = e.from ^ e.to ^ v;
            depth[to] = depth[v] + 1;
            dist[to] = dist[v] + e.cost;
            pv[to] = v;
            pe[to] = id;
            root[to] = (root[v] != -1 ? root[v] : to);
            do_dfs(to);
            sz[v] += sz[to];
        }
        end[v] = (int) order.size() - 1;
    }
    
    void do_dfs_from(int v) {
        depth[v] = 0;
        dist[v] = T{};
        root[v] = v;
        pv[v] = pe[v] = -1;
        do_dfs(v);
    }
    
    public:
    void dfs(int v, bool clear_order = true) {
        if (pv.empty()) {
            init();
        } else {
            if (clear_order) {
                order.clear();
            }
        }
        do_dfs_from(v);
    }
    
    void dfs_all() {
        init();
        for (int v = 0; v < n; v++) {
            if (depth[v] == -1) {
                do_dfs_from(v);
            }
        }
        assert((int) order.size() == n);
    }
};
    
template <typename T>
class lca_forest : public dfs_forest<T> {
    public:
    using dfs_forest<T>::edges;
    using dfs_forest<T>::g;
    using dfs_forest<T>::n;
    using dfs_forest<T>::pv;
    using dfs_forest<T>::pos;
    using dfs_forest<T>::end;
    using dfs_forest<T>::depth;
    
    int h;
    vector<vector<int>> pr;
    
    lca_forest(int _n) : dfs_forest<T>(_n) {
    }
    
    inline void build_lca() {
        assert(!pv.empty());
        int max_depth = 0;
        for (int i = 0; i < n; i++) {
            max_depth = max(max_depth, depth[i]);
        }
        h = 1;
        while ((1 << h) <= max_depth) {
            h++;
        }
        pr.resize(n);
        for (int i = 0; i < n; i++) {
            pr[i].resize(h);
            pr[i][0] = pv[i];
        }
        for (int j = 1; j < h; j++) {
            for (int i = 0; i < n; i++) {
                pr[i][j] = (pr[i][j - 1] == -1 ? -1 : pr[pr[i][j - 1]][j - 1]);
            }
        }
    }
    
    inline bool anc(int x, int y) {
        return (pos[x] <= pos[y] && end[y] <= end[x]);
    }
    
    inline int go_up(int x, int up) {
        assert(!pr.empty());
        up = min(up, (1 << h) - 1);
        for (int j = h - 1; j >= 0; j--) {
            if (up & (1 << j)) {
                x = pr[x][j];
                if (x == -1) {
                break;
                }
            }
        }
        return x;
    }
    
    inline int lca(int x, int y) {
        assert(!pr.empty());
        if (anc(x, y)) {
            return x;
        }
        if (anc(y, x)) {
            return y;
        }
        for (int j = h - 1; j >= 0; j--) {
            if (pr[x][j] != -1 && !anc(pr[x][j], y)) {
                x = pr[x][j];
            }
        }
        return pr[x][0];
    }
};

/*_______________________________________________________________*/
/*_______________________Topological Sort________________________*/
template <typename T>
vector<int> find_topsort(const digraph<T> &g) {
    vector<int> deg(g.n, 0);
    for (int id = 0; id < (int) g.edges.size(); id++) {
        deg[g.edges[id].to]++;
    }
    vector<int> x;
    for (int i = 0; i < g.n; i++) {
        if (deg[i] == 0) {
            x.push_back(i);
        }
    }
    for (int ptr = 0; ptr < (int) x.size(); ptr++) {
        int i = x[ptr];
        for (int id : g.g[i]) {
            auto &e = g.edges[id];
            int to = e.to;
            if (--deg[to] == 0) {
                x.push_back(to);
            }
        }
    }
    if ((int) x.size() != g.n) {
        return vector<int>();
    }
    return x;
}

/*_______________________________________________________________*/
/*______________________________SCC______________________________*/



void init_code() {
    ios_base::sync_with_stdio(false); cin.tie(0); cout.tie(0);
    // 本地调试交互式题目可参考1705F
#ifndef ONLINE_JUDGE

#else

#endif
}
/*
 * 求最长严格递增子序列长度
 */
int LIS(VL v) {
    int n = v.size();
    VI dp(n);
    dp[0] = 1;
    int ans = 0;
    ff (i, 1, n) {
        int cnt = 1;
        FF (j, i - 1, 0) {
            if (v[i] > v[j]) {
                cnt = max(cnt, dp[j] + 1);
            }
        }
        dp[i] = cnt;
    }
    ff (i, 0, n) {
        ans = max(ans, dp[i]);
    }
    return ans;
}

/*
 *
 */
void solve() {
    ll a, b, c;
    cin >> a >> b >> c;
    cout << (a + b * c) << " " << b << " " << c << endl;
}

/*_______________________________________________________________*/
signed main()
{
    init_code();
	ll t = 1;
    cin >> t;
    cout << spc(20) << endl;
	while(t--){
        solve();
	}
}
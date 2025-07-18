/*

  Freedom of action

  If there is no soul, there is still a body

*/

#pragma GCC optimize(3)

#include <bits/stdc++.h>

#include <ext/pb_ds/assoc_container.hpp>

#include <ext/pb_ds/tree_policy.hpp>

#include <ext/pb_ds/trie_policy.hpp>

#include <ext/pb_ds/priority_queue.hpp>

using namespace __gnu_pbds;

using namespace std;

namespace Template {

    #define IOS ios::sync_with_stdio(false);cin.tie(nullptr);cout.tie(nullptr)

    #define FREOPEN freopen("in.in", "r", stdin);freopen("out.out", "w", stdout)

    #ifdef A_king

    #include "E:\VS\PROJECT\Another\debug.h"

    #else

    #define out(x...) void(0);

    #endif

    #define endl '\n'

    #define pb push_back

    #define pf push_front

    #define len(container) (int)(container).size()

    #define all(container) (container).begin(), (container).end()

    #define rall(container) (container).rbegin(), (container).rend()

    #define YES cout << "YES" << endl

    #define Yes cout << "Yes" << endl

    #define yes cout << "yes" << endl

    #define NO cout << "NO" << endl

    #define No cout << "No" << endl

    #define no cout << "no" << endl

    const int mod = 998244353;const double PI = acos(-1.0);const double eps = 1e-8;const int inf = 0x3f3f3f3f;const long long INF = 0x3f3f3f3f3f3f3f3f;

    template<typename T> T ksm(T a, long long b) {T res = 1;while (b) {if (b & 1) res *= a;a *= a;b >>= 1;}return res;}

    template<const int MOD = mod> struct Modular {

        int x;

        int Mod(int x) {if (x < 0) x += MOD;if (x >= MOD) x -= MOD;return x;}

        Modular(int x = 0) : x(Mod(x)) {}

        Modular(long long x) : x(Mod(x % MOD)) {}

        int val() const {return x;}

        Modular operator -() const {Modular T(MOD - x);return T;}

        Modular inv() const {assert(x != 0);return ksm(*this, MOD - 2);}

        Modular inverse() const {long long a = x, b = MOD, u = 1, v = 0;while (b) {long long t = a / b;a -= t * b;a ^= b ^= a ^= b;u -= t * v;u ^= v ^= u ^= v;}if (u < 0) u += MOD;return u;}

        Modular &operator *= (const Modular &T) {x = (long long)(x) * T.x % MOD;return *this;}

        Modular &operator += (const Modular &T) {x = Mod(x + T.x);return *this;}

        Modular &operator -= (const Modular &T) {x = Mod(x - T.x);return *this;}

        Modular &operator /= (const Modular &T) {return *this *= T.inverse();}

        Modular &operator ++ (int) {return *this = *this + 1;}

        Modular &operator -- (int) {return *this = *this - 1;}

        bool operator == (const Modular &T) const {return x == T.x;}

        bool operator != (const Modular &T) const {return x != T.x;}

        bool operator <= (const Modular &T) const {return x <= T.x;}

        bool operator >= (const Modular &T) const {return x >= T.x;}

        bool operator < (const Modular &T) const {return x < T.x;}

        bool operator > (const Modular &T) const {return x > T.x;}

        friend Modular operator * (const Modular &T, const Modular &Y) {Modular res = T;res *= Y;return res;}

        friend Modular operator + (const Modular &T, const Modular &Y) {Modular res = T;res += Y;return res;}

        friend Modular operator - (const Modular &T, const Modular &Y) {Modular res = T;res -= Y;return res;}

        friend Modular operator / (const Modular &T, const Modular &Y) {Modular res = T;res /= Y;return res;}

        friend istream &operator >> (istream &in, Modular& T) {long long val;in >> val;T = val;return in;}

        friend ostream &operator << (ostream &os, const Modular& T) {return os << T.x;}

    };

    typedef trie<string, null_type, trie_string_access_traits<>, pat_trie_tag, trie_prefix_search_node_update> TRIE;

    template<typename T> using RB = tree<T, null_type, less<T>, rb_tree_tag, tree_order_statistics_node_update>;

    template<typename T> using Q = __gnu_pbds::priority_queue<T, greater<T>, pairing_heap_tag>;// point_iterator

    typedef long long ll;typedef unsigned long long ull;typedef pair<double, double> PDD;typedef pair<ll, ll> PLL;typedef pair<int, int> PII;typedef pair<int, PII> PIII;

    template<typename T> using vc = vector<T>;template<typename T> using vvc = vc<vc<T>>;

    template<typename T> inline T Abs(T a) {if (a < 0) a = -1 * a;return a;}

    template<typename T> T Gcd(const T &a, const T &b) {return b ? Gcd(b, a % b) : a;}

    template<typename T> T Lcm(const T &a, const T &b) {return a / Gcd(a, b) * b;}

    template<typename T> inline void Swap(T &a, T &b) {a ^= b ^= a ^= b;}

    template<typename T> inline bool Max(T &a, const T &b) {return a < b ? a = b, 1 : 0;}template<typename T> inline void Max(T &a, const T &b, const T &c){Max(a, b);Max(a, c);}

    template<typename T> inline bool Min(T &a, const T &b) {return a > b ? a = b, 1 : 0;}template<typename T> inline void Min(T &a, const T &b, const T &c){Min(a, b);Min(a, c);}

    template<typename T> inline T Sum(const vector<T> &x, int pos = 0) {return accumulate(x.begin() + pos, x.end(), 0ll);}

    template<typename T> inline T Maxe(const vector<T> &x, int pos = 0) {return *max_element(x.begin() + pos, x.end());}template<typename T> inline T Maxi(vector<T> x, int pos = 0) {return max_element(x.begin() + pos, x.end()) - x.begin();}

    template<typename T> inline T Mine(const vector<T> &x, int pos = 0) {return *min_element(x.begin() + pos, x.end());}template<typename T> inline T Mini(vector<T> x, int pos = 0) {return min_element(x.begin() + pos, x.end()) - x.begin();}

    template<typename T> inline void Disperse(vc<T> &v) {sort(all(v));v.erase(unique(all(v)), v.end());}

    template<typename T, typename U> istream &operator >> (istream &in, pair<T, U> &Arg) {return in >> Arg.first >> Arg.second;}

    template<typename T> istream &operator >> (istream &in, vc<T> &v) {int n = len(v) - 1;for (int i = 1;i <= n;++ i) in >> v[i];return in;}

    template<typename T> istream &operator >> (istream &in, vvc<T> &v) {int n = len(v) - 1, m = len(v[0]) - 1;for (int i = 1;i <= n;++ i) for (int j = 1;j <= m;++ j) in >> v[i][j];return in;}

    inline int popcount(const long long &x) {return __builtin_popcountll(x);}inline int clz(const long long &x) {return __builtin_clzll(x);}inline int parity(const long long &x) {return __builtin_parityll(x);}

}using namespace Template;

using Z = Modular<>;// 998244353 1000000007

const int N = 1e6 + 10, M = 2e6 + 10;







// #define int long long

bool multicase = true;

void solve(int group_Id) {

    int n, m, k1;

    cin >> n >> m >> k1;

    int k2 = n * m / 2 - k1, flag = 0;

    if (n & 1) {

        flag = 1;

        int now = m / 2;

        if (k1 < now) {NO;return ;}

        k1 -= now;

        n --;

    }else if (m & 1) {

        flag = 2;

        int now = n / 2;

        if (k2 < now) {NO;return ;}

        k2 -= now;

        m --;

    }

    if ((k1 & 1) or (k2 & 1)) {NO;return ;}

    YES;

    vvc<char> g(n + 5, vc<char> (m + 5));

    if (flag == 1) {

        for (int i = 1;i <= m;i += 4) {

            g[n + 1][i] = 'a';

            g[n + 1][i + 1] = 'a';

            g[n + 1][i + 2] = 'b';

            g[n + 1][i + 3] = 'b';

        }

    }else if (flag == 2) {

        for (int i = 1;i <= n;i += 4) {

            g[i][m + 1] = 'a';

            g[i + 1][m + 1] = 'a';

            g[i + 2][m + 1] = 'b';

            g[i + 3][m + 1] = 'b';

        }

    }

    char c = 'c';

    for (int i = 1;i <= n;i += 2) {

        for (int j = 1;j <= m;j += 2) {

            if (k1) {

                k1 -= 2;

                g[i][j] = c;

                g[i][j + 1] = c;

                g[i + 1][j] = c + 1;

                g[i + 1][j + 1] = c + 1;

                c += 2;

                if (c - 'a' == 26) c = 'c';

                while (g[i][j] == g[i - 1][j] or g[i][j + 1] == g[i - 1][j + 1]

                or g[i][j] == g[i][j - 1] or g[i + 1][j] == g[i + 1][j - 1]) {

                    g[i][j] = c;

                    g[i][j + 1] = c;

                    g[i + 1][j] = c + 1;

                    g[i + 1][j + 1] = c + 1;

                    c += 2;

                    if (c - 'a' == 26) c = 'c';

                }

            }else {

                g[i][j] = c;

                g[i + 1][j] = c;

                g[i][j + 1] = c + 1;

                g[i + 1][j + 1] = c + 1;

                c += 2;

                if (c - 'a' == 26) c = 'c';

                while (g[i][j] == g[i - 1][j] or g[i][j + 1] == g[i - 1][j + 1]

                or g[i][j] == g[i][j - 1] or g[i + 1][j] == g[i + 1][j - 1]) {

                    g[i][j] = c;

                    g[i + 1][j] = c;

                    g[i][j + 1] = c + 1;

                    g[i + 1][j + 1] = c + 1;

                    c += 2;

                    if (c - 'a' == 26) c = 'c';

                }

            }

        }

    }

    if (flag == 1) n ++;

    else if (flag == 2) m ++;

    for (int i = 1;i <= n;i ++) {

        for (int j = 1;j <= m;j ++) {

            cout << g[i][j];

        }

        cout << endl;

    }

}



signed main(signed argc, char const *argv[])

{

#ifdef A_king

    FREOPEN;

    auto clock_start = chrono::high_resolution_clock::now();

    cerr << __FILE__ << ", " << __DATE__ << ", " << __TIME__ << endl;

#endif

    IOS;

    cout << fixed << setprecision(16);

    int T_case = 1;

    if (multicase) cin >> T_case;

    for (int i = 1;i <= T_case;i ++) solve(i);

#ifdef A_king

    auto clock_end = chrono::high_resolution_clock::now();

    cerr << "Run Time: " << chrono::duration_cast<chrono::milliseconds>(clock_end - clock_start).count() << " ms" << endl;

#endif

    return 0;

}
#include <bits/stdc++.h>

using namespace std;

#define int         long long int

#define ll          long long

#define tc          int test;cin>>test;while(test--)

#define rep1(n)     for(int i=1;i<=n;++i)

#define rep0(n)     for(int i=0;i<n;++i)

#define rrep(i,z,a) for(int i=z;i>=a;--i)

#define rep(i,a,b)  for(int i=a;i<b;++i)

#define kill(x)     return cout<<x<<endl,0

#define all(v)      v.begin(),v.end()

#define pb          push_back

#define ss          second

#define ff          first

#define endl        "\n"

#define spc         " "

#define mod         1000000007

#define maxn        998244353

#define inf         5e18

template<typename T1, typename T2>istream& operator>>(istream& in, pair<T1, T2> &a) {in >> a.ff >> a.ss; return in;}

template<typename T1, typename T2>ostream& operator<<(ostream& out, pair<T1, T2> a) {out << a.ff << " " << a.ss; return out;}

template<class T>ostream &operator <<(ostream& out, vector<T>&v) {for (auto &x : v)cout << x << " "; return out;}

template<class T>istream &operator >>(istream& in, vector<T>&v) {for (auto &x : v)cin >> x; return in;}

template<typename T, typename T1>T amax(T &a, T1 b) {if (b > a)a = b; return a;}

template<typename T, typename T1>T amin(T &a, T1 b) {if (b < a)a = b; return a;}

typedef pair<int, int> pii;

typedef vector<int> vi;

typedef vector<pii> vii;

//debugger

void _print(int t) {cerr << t;}

void _print(string t) {cerr << t;}

void _print(char t) {cerr << t;}

void _print(long double t) {cerr << t;}

void _print(double t) {cerr << t;}

void _print(unsigned long long t) {cerr << t;}

template <class T, class V> void _print(pair <T, V> p);

template <class T> void _print(vector <T> v);

template <class T> void _print(set <T> v);

template <class T, class V> void _print(map <T, V> v);

template <class T> void _print(multiset <T> v);

template <class T, class V> void _print(pair <T, V> p) {cerr << "{"; _print(p.ff); cerr << ","; _print(p.ss); cerr << "}";}

template <class T> void _print(vector <T> v) {cerr << "[ "; for (T i : v) {_print(i); cerr << " ";} cerr << "]";}

template <class T> void _print(set <T> v) {cerr << "[ "; for (T i : v) {_print(i); cerr << " ";} cerr << "]";}

template <class T> void _print(multiset <T> v) {cerr << "[ "; for (T i : v) {_print(i); cerr << " ";} cerr << "]";}

template <class T, class V> void _print(map <T, V> v) {cerr << "[ "; for (auto i : v) {_print(i); cerr << " ";} cerr << "]";}

#ifndef ONLINE_JUDGE

#define debug(...) __f(#__VA_ARGS__, __VA_ARGS__)

template <typename Arg1>

void __f(const char* name, Arg1&& arg1) {cerr << name << " : "; _print(arg1); cerr << endl;}

template <typename Arg1, typename... Args>

void __f(const char* names, Arg1&& arg1, Args&&... args) {

    const char* comma = strchr(names + 1, ','); cerr.write(names, comma - names) << " : ";

    _print(arg1); cerr << endl; __f(comma + 2, args...);

}

#else

#define debug(...)

#endif

//debugger



const int N = 5000005;

const double eps = 1e-9;

double v;

double f(double c, double m, double p)

{

    if (abs(p - 1) < eps)return 1;

    double e1(0), e2(0);

    debug(c, m, p);

    if (abs(c) < eps) {

        e2 += f(c, m - min(m, v), p + min(m, v));

        return 1 + m * e2;

    }

    if (abs(m) < eps) {

        e1 += f(c - min(c, v), m, p + min(c, v));

        return 1 + c * e1;

    }



    e1 += f(c - min(c, v), m + min(c, v) / 2, p + min(c, v) / 2);

    e2 += f(c + min(m, v) / 2, m - min(m, v), p + min(m, v) / 2);

    return 1 + c * e1 + m * e2;

}

int solve()

{

    double c, p, m;

    cin >> c >> m >> p >> v;

    cout << fixed << setprecision(12) << f(c, m, p);

    return 0;

}

signed main()

{

    ios::sync_with_stdio(0); cin.tie(0); cout.tie(0);

#ifndef ONLINE_JUDGE

    freopen("input.txt", "r", stdin);

    freopen("output.txt", "w", stdout);

    freopen("error.txt", "w", stderr);

#endif

    tc

    solve(), cout << endl;

    cerr << 1000.0 * (float)clock() / CLOCKS_PER_SEC << " ms\n";

    return 0;

}
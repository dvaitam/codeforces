#include <bits/stdc++.h>

using namespace std;



using i64 = long long;



#define ll long long

#define endl "\n"

#define FOR(i,n) for(long long i = 0; i < n; i++)

#define FOR1(i,s,e,d) for(long long i = s; i < e; i += d)

#define FOR2(i,s,e,d) for(long long i = s; i > e; i -= d)

#define FOR3(x,v) for (auto x : v)

#define pll pair<long long, long long>

#define vi vector<int>

#define vl vector<long long>

#define vll vector<pair<long long, long long>>

#define vvl vector<vector<long long>>

#define vs vector<string>

#define vc vector<char>

#define mll map<long long, long long>

#define mlc map<long long,char>

#define mls map<long long,string>

#define mlvl map<long long, vector<long long>>

#define pql priority_queue<long long>

#define over(ans) {cout<<ans<<endl;return;}

#define all(v) (v).begin(), (v).end()





const ll INF = 1e9;

const ll mod = 998244353;

const ll mod1 = 1e9 + 7;



template<class T, class T1> bool ckmin(T& a, const T1& b) { return b < a ? a = b, 1 : 0; }

template<class T, class T1> bool ckmax(T& a, const T1& b) { return a < b ? a = b, 1 : 0; }



template<class c>

void display(c vect) {

    FOR3(x, vect) cout << x << ' ';

    cout << "\n";

}



void display(vll v) {

    FOR3(x, v) cout << "(" << x.first << ", " << x.second << ") ";

    cout << endl;

}



void display(vvl vect2d) {

    FOR3(x, vect2d) display(x);

}



void display(mll m) {

    FOR3(x, m) cout << x.first << ' ' << x.second << endl;

}



void display(mlvl  m) {

    FOR3(x, m) {

        cout << x.first << " : ";

        display(x.second);

    }

}



vl vin(ll n) {

    vl v(n);

    FOR(i,n) cin >> v[i];

    return v;

}



vi vin(int n) {

    vi v(n);

    FOR(i, n) cin >> v[i];

    return v;

}



int to_int(string s) {

    return stoi(s);

}



int to_int(char c) {

    string s = {c};

    return stoi(s);

}



char to_char(int i) {

    return char('0' + i);

}



bool is_prime(ll n) {

    if (n == 1) return false;

    if (n == 2 || n == 3) return true;

    if ( n % 6 == 1 || n % 6 == 5) {

        for (ll k = 5; k*k <= n; k += 2) if (n % k == 0) return false;

        return true;

    }

    return false;

}



ll gcd(ll a, ll b) {

    return b == 0 ? a : gcd(b, a % b);   

}



ll inv(ll a, ll b){

    return 1<a ? b - inv(b%a,a)*b/a : 1;

}



template<class T> 

T power(T a, ll b) {

    T res = 1;

    for(; b; b /= 2, a *= a) if (b%2) {res *= a;}

    return res;

}



template<class T>

T powermod(T a, ll b, ll P) {

    a %= P; ll res = 1;

    while (b > 0) {

        if (b & 1) res = res * a % P;

        a = a * a % P;

        b >>= 1;

    }

    return res;

}



vl factors(ll n) {

    vl factors;

    ll i = 1;

    while(i*i <= n) {

        if (n%i == 0) {

            factors.push_back(i);

            if (i*i != n) factors.push_back(n/i);

        }

        i++;

    }

    sort(factors.begin(), factors.end());

    return factors;

}



ll lsb(ll a) {return (a & -a);}



template < typename T>

struct Fenwick {

    const ll n;

    vector<T> a;

    Fenwick(ll n) : n(n), a(n) {}

    void add(ll x, T v) {

        for (ll i = x+1; i <= n; i += i & -i) {

            a[i-1] += v;

        }

    }

    T sum(ll x) {

        T ans = 0;

        for (ll i = x; i > 0; i -= i & -i) {

            ans += a[i-1];

        }

        return ans;

    }

    T rangeSum(ll l, ll r) {

        return sum(r) - sum(l);

    }

};



struct DSU {

    vl f, siz;

    DSU(ll n) : f(n), siz(n-1) { iota(f.begin(), f.end(), 0); }

    ll leader(ll x) {

        while (x != f[x]) x = f[x] = f[f[x]];

        return x;

    }

    bool same(ll x, ll y) { return leader(x) == leader(y); }

    bool merge(ll x, ll y) {

        x = leader(x);

        y = leader(y);

        if (x == y) return false;

        siz[x] += siz[y];

        f[y] = x;

        return true;

    }

    ll size(ll x) { return siz[leader(x)]; }

};



ll ncrmodP(ll n, ll r, ll P) {

    ll fac[n+1];

    fac[0] = 1;

    FOR1(i, 1, n+1, 1) {

        fac[i] = (fac[i-1] * i) % P;

    }



    return (fac[n] * inv(fac[r], P) % P * inv(fac[n-r], P) % P) % P;

}



vi primes_n(int n) {

    vi primes;

    FOR1(i, 1, n+1, 1) {

        if (is_prime(i)) primes.push_back(i);

    }

    return primes;

}



// (condition ? value if yes : value if no)

// sort(v.begin(), v.end(), [](data_type x, data_type y) -> bool {

//     return condition;

// });

// use 1.0 for decimal calculations

// cout << setprecision(12) << fixed << ans << endl;





void init() {



}



void solve() {

    int n; cin >> n;

    double ans = 0;

    FOR1(i, 1, n+1, 1) ans += 1.0 / i;

    cout << setprecision(12) << fixed << ans << endl;

}   



int main() {

    ios::sync_with_stdio(false), cin.tie(nullptr);

    init();

    int t = 1;

    // cin >> t;

    FOR(i,t) {

        // cout << "Case #" << i+1 << ": ";

        solve();

    }

}
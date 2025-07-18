#include "bits/stdc++.h"



using namespace std;

using ll = long long;

using pll = pair<ll, ll>;

using vll = vector<ll>;

using vpl = vector<pll>;

using ld = long double;

using str = string;

using big_int = __int128_t;



const ld eps = 1e-7;

const ld PI = acos(-1);



#define all(c) (c).begin(), (c).end()

#define rall(c) ((c).rbegin()), ((c).rend())

#define ff first

#define ss second

#define pb push_back

#define pf push_front

#define fast ios_base::sync_with_stdio(0); cin.tie(0)

#define forn(i, n) for (ll i = 0; i < n; ++i)

#define sz(a) (ll)a.size()

#define endl '\n'

#define u_map unordered_map

#define mset multiset

//#define x first

//#define y second



#ifdef ONLINE_JUDGE

#define debug(x);

#else

#define debug(x) cerr << #x << ": " << x << endl;

#endif



str IO[2] = {"NO\n", "YES\n"};

str io[2] = {"no\n", "yes\n"};

str Io[2] = {"No\n", "Yes\n"};



mt19937 rnd(chrono::steady_clock::now().time_since_epoch().count());

//mt19937 rnd(1232423);



template<class T> bool inmin(T& x_, T y_) {return y_ < x_ ? (x_ = y_, true) : false;}

template<class T> bool inmax(T& x_, T y_) {return y_ > x_ ? (x_ = y_, true) : false;}



//const ll Mod = 1e9 + 7;

const ll Mod = 998244353;

//const ll Mod = 1e9 + 9;

//const ll Mod = 1234567891;

const ll INF = 1e16;



template<class T> void add(T& x_, T y_) { x_ = (x_ + y_) % Mod; };

template<class T> void sub(T& x_, T y_) { x_ = (x_ + Mod - y_) % Mod; };

template<class T> void mul(T& x_, T y_) { x_ = (x_ * y_) % Mod; };



inline void solve() {

    ll n;

    cin >> n;

    vll a(n);

    forn (i, n) cin >> a[i];

    reverse(all(a));

    forn (i, n) {

        if (2 * i < n) cout << -a[i] << ' ';

        else cout << a[i] << ' ';

    }

    cout << endl;

}



int main() {

    fast;

    //cout << fixed << setprecision(10);

    ll test;

    test = 1;

    cin >> test;

    for (ll id = 0; id < test; ++id) solve();

    return 0;

}
#include <bits/stdc++.h>

using namespace std;



#ifdef DEBUG

#include "debug.h"

#else

#define debug(...) 1

#endif



using ll = int;

using db = long double;

using VS = vector<string>;

using VLL = vector<ll>;

using VVLL = vector<VLL>;

using VVVLL = vector<VVLL>;

using PLL = pair<ll, ll>;

using MLL = map<ll, ll>;

using SLL = set<ll>;

using QLL = queue<ll>;

using SS = stringstream;



#define rep(x, l, u) for (ll x = l; x < u; x++)

#define rrep(x, l, u) for (ll x = l; x >= u; x--)

#define fe(x, a) for (auto x : a)

#define all(x) x.begin(), x.end()

#define rall(x) x.rbegin(), x.rend()

#define mst(x, v) memset(x, v, sizeof(x))

#define sz(x) (ll) x.size()

#define ios ios_base::sync_with_stdio(0), cin.tie(0)



#define umap unordered_map

#define uset unordered_set

#define mset multiset



// clang-format off



ll ob(ll i, ll n) { return i < 0 || i >= n; }

ll tp(ll x) { return ( 1LL << x ); }

ll rup(ll a, ll b) { return a % b ? a/b + 1 : a/b; }

ll sign(ll x) {	return x == 0 ? 0 : x / abs(x); }

void makemod(ll& x, ll m) { x %= m; if (x < 0) { x += m; } }

ll getmod(ll x, ll m) { makemod(x, m); return x; }

ll powmod(ll a, ll b, ll m) { if (b == 0) return 1; ll h = powmod(a, b/2, m); ll ans = h*h%m; return b%2 ? ans*a%m : ans; }

ll invmod(ll a, ll m) { return powmod(a, m - 2, m); }



template <typename A, typename B> bool upmin(A& x, B v) { if (v >= x) return false; return x = v, true; }

template <typename A, typename B> bool upmax(A& x, B v) { if (v <= x) return false; return x = v, true; }



// clang-format on



const VLL di = {0, 0, 1, -1, 1, -1, 1, -1}, dj = {1, -1, 0, 0, -1, -1, 1, 1};

mt19937 rng(chrono::steady_clock::now().time_since_epoch().count());



ll mod = 3;



ll invThree(ll x) {

    return x;  // lol

}



// system of equations

namespace SOE {

vector<vector<short>> A;  // matrix

vector<short> values;     // rhs

int n, m;



void clear() {

    A.clear();

    values.clear();

    n = 0, m = 0;

}



void add(vector<short>& lhs, short rhs) {

    A.push_back(lhs);

    values.push_back(rhs);

    n++;

    m = lhs.size();

}



void multK(ll ind, ll k) {

    // ind = ind * k

    for (auto& a : A[ind]) {

        a *= k;

        a %= 3;

    }

    values[ind] *= k;

    values[ind] %= 3;

}



void addMultK(ll from, ll to, ll k) {

    // to = to + from * k

    rep(j, 0, m) {

        A[to][j] += A[from][j] * k;

        A[to][j] %= 3;

    }

    values[to] += values[from] * k;

    values[to] %= 3;

}



void swapTwo(ll from, ll to) {

    swap(A[from], A[to]);

    swap(values[from], values[to]);

}



void solve() {

    ll ind = 0;  // row

    rep(j, 0, m) {

        // look for the first row with j'th col non-zero

        ll found = 0;

        rep(ii, ind, n) {

            if (A[ii][j]) {

                swapTwo(ii, ind);

                found = 1;

                break;

            }

        }



        if (!found) continue;



        // use ind

        // divide to get to 1

        multK(ind, invThree(A[ind][j]));



        // for every other row, get it to 0

        rep(otherRow, 0, n) {

            if (otherRow == ind) continue;

            if (A[otherRow][j] == 0) continue;



            addMultK(ind, otherRow, 3 - A[otherRow][j]);

        }



        ind++;

    }

}



}  // namespace SOE



void solve() {

    SOE::clear();



    ll n, m;

    cin >> n >> m;



    VVLL edges(m);

    VVLL adj(n, VLL(n, -1));

    rep(i, 0, m) {

        ll a, b, c;

        cin >> a >> b >> c;

        a--, b--;



        edges[i] = {a, b, c};

        adj[a][b] = i;

        adj[b][a] = i;



        if (c != -1) {

            vector<short> row(m);

            row[i] = 1;

            SOE::add(row, c - 1);

        }

    }



    // find triangles

    rep(i, 0, n) {

        rep(j, 0, i) {

            rep(k, 0, j) {

                ll a = adj[i][j];

                ll b = adj[j][k];

                ll c = adj[k][i];



                if (a != -1 && b != -1 && c != -1) {

                    // triangle! a + b + c = 0

                    vector<short> row(m);

                    row[a] = 1;

                    row[b] = 1;

                    row[c] = 1;



                    SOE::add(row, 0);

                }

            }

        }

    }



    SOE::solve();



    VLL ans(m);

    rep(j, 0, m) {

        rep(i, 0, SOE::n) {

            if (SOE::A[i][j]) {

                ans[j] = SOE::values[i];

                SOE::A[i][j] = 0;

                SOE::values[i] = 0;



                break;

            }

        }

    }

    rep(i, 0, SOE::n) {

        if (SOE::values[i]) {

            cout << -1 << endl;

            return;

        }

    }



    fe(a, ans) cout << a + 1 << ' ';

    cout << endl;

}



int main() {

    ios;

    ll t = 1;

    cin >> t;

    rep(i, 0, t) solve();

    return 0;

}
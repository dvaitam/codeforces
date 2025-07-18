//#pragma GCC optimize("O3")
//#pragma GCC target("sse4.1")

#ifdef LOCAL
#define _GLIBCXX_DEBUG
#endif // LOCAL
#include <bits/stdc++.h>
using namespace std;
#define pb push_back
#define mp make_pair
#define fst first
#define snd second
#define forn(i, n) for (int i = 0; i < (int)(n); ++i)
#define ford(i, n) for (int i = ((int)(n) - 1); i >= 0; --i)
typedef long long ll;
typedef vector<int> vi;
typedef vector<vi> vvi;
typedef vector<ll> vll;
typedef vector<char> vc;
typedef vector<vc> vvc;
typedef pair<int, int> pii;
typedef vector<pii> vii;
#define sz(c) (int)(c).size()
#define all(c) (c).begin(), (c).end()

#ifdef LOCAL
#define eprintf(args...) fprintf(stderr, args), fflush(stderr);
#else
#define eprintf(args...) ;
#endif

ll mult (ll a, ll b, ll mod)
{
    using T = long double;
    ll quot = (ll)((T)a * (T)b / (T)mod);
    ll rem = a * b - quot * mod;
    return (rem + 5 * mod) % mod;
}

ll powmod (ll a, ll pw, ll mod)
{
    ll res = 1;
    while (pw)
    {
        if (pw & 1)
            res = mult(res, a, mod);
        a = mult(a, a, mod);
        pw >>= 1;
    }

    return res;
}

ll m, x;

bool read ()
{
    if (!(cin >> m >> x))
        return false;
    return true;
}

vector<pair<ll, int>> get (ll n)
{
    vector<pair<ll, int>> ans;

    for (int i = 2; (ll)i * i <= n; i++) if (n % i == 0)
    {
        int cnt = 0;
        while (n % i == 0)
            cnt++, n /= i;

        ans.pb(mp(i, cnt));
    }

    if (n > 1)
        ans.pb(mp(n, 1));

    return ans;
}

ll lcm (ll a, ll b)
{
    return a / __gcd(a, b) * b;
}

ll rec (int pos, ll cur_lcm, ll cur_phi, const vector<vll> &orders, const vi &degs, const vll &ps)
{
    if (pos == sz(degs))
    {
        assert(cur_phi % cur_lcm == 0);
        return cur_phi / cur_lcm;
    }

    ll ans = 0;

    forn (i, degs[pos] + 1)
    {
        ans += rec(pos + 1, lcm(cur_lcm, orders[pos][i]), cur_phi, orders, degs, ps);
        cur_phi *= (i == 0 ? ps[pos] - 1 : ps[pos]);
    }

    return ans;
}

ll find_order (ll x, ll q, ll p, const vector<pair<ll, int>> &fact_prs)
{
    vll ps(sz(fact_prs));
    forn (i, sz(fact_prs))
        ps[i] = fact_prs[i].fst;

    ll phi = 1;
    forn (i, sz(fact_prs)) forn (j, fact_prs[i].snd)
        phi *= fact_prs[i].fst;

    x %= q;
    assert(powmod(x, phi, q) == 1);

    while (true)
    {
        bool ok = false;
        for (ll pp : ps) if (phi % pp == 0 && powmod(x, phi / pp, q) == 1)
            phi /= pp, ok = true;

        if (!ok)
            break;
    }

    return phi;
}

ll solve ()
{
    vector<pair<ll, int>> fact = get(m);

    int divs = 1;
    for (const auto &pr : fact)
        divs *= pr.snd + 1;

    vector<vector<pair<ll, int>>> fact_prs (sz(fact));
    forn (i, sz(fact))
        fact_prs[i] = get(fact[i].fst - 1);

    vector<vll> orders(sz(fact));
    forn (i, sz(fact))
    {
        orders[i].resize(fact[i].snd + 1);
        orders[i][0] = 1;

        ll q = fact[i].fst;

        forn (j, fact[i].snd + 1) if (j)
        {
            orders[i][j] = find_order(x, q, fact[i].fst, fact_prs[i]);

            q *= fact[i].fst;
            if (j == 1)
                fact_prs[i].pb(mp(fact[i].fst, 0));
            fact_prs[i].back().snd++;
        }
    }

    vi degs(sz(fact));
    vll ps(sz(fact));
    forn (i, sz(fact))
        tie(ps[i], degs[i]) = fact[i];

    ll ans = rec(0, 1LL, 1LL, orders, degs, ps);
    return ans;
}

int main ()
{
#ifdef LOCAL
    freopen("input.txt", "rt", stdin);
    freopen("output.txt", "w", stdout);
#endif // LOCAL

    ios_base::sync_with_stdio(false);
    cin.tie(nullptr);

    while (read())
    {
        auto ans = solve();
        cout << ans << endl;
    }
}
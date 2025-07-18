#include <bits/stdc++.h>
#include <numeric>
#include <algorithm>
#include <functional>
#define nl '\n'
#define ll long long
#define vi vector<int>
#define vl vector<ll>
#define pi pair<int, int>
#define pl pair<ll, ll>
#define setl set<ll>
#define mstl multiset<ll>
#define mpvl map<int, vector<char>>
#define map_i map<int, int>
#define map_l map<ll, ll>
#define map_c map<char, int>
#define vc vector<char>
#define vs vector<string>
#define vcprll vector<pair<ll, ll>>
#define mpvcll_mp map<ll, vector<ll>> mp;
#define ump unordered_map<ll, ll>
#define pb push_back
#define pbb pop_back
#define m_p make_pair
#define mnv(v) *min_element(v.begin(), v.end())
#define mxv(v) *max_element(v.begin(), v.end())
#define countv(v, a) count(v.begin(), v.end(), a)
#define acumlt(x) accumulate(x.begin(), x.end(), 0LL);
#define srt(v) sort(v.begin(), v.end())
#define srt_r(a) sort(a.rbegin(), a.rend());
#define bct(x) __builtin_popcountll(x)
#define full(x) x.begin(), x.end()
#define grtall(v) sort(v.begin(), v.end(), greater<int>())
#define rev(v) reverse(v.begin(), v.end())
#define lower(v, x) lower_bound(all(v), x) - v.bg + 1
#define upper(v, x) upper_bound(all(v), x) - v.bg + 1
#define lwrbnd(x) lower_bound(x)
#define upbnd(x) upper_bound(x)
#define multiset_pr_i multiset<pair<int, int>>
#define sum(v, x) accumulate(all(v), x)
#define flI(i, a, n) for (ll i = a; i < n; i++)
#define fld(i, a, n) for (ll i = a; i <= n; i++)
#define flrd(i, a, n) for (ll i = a; i >= n; i--)
#define flrc(i, a, x) for (char c = a; c < x; c++)
#define fl2(i, a, n) for (ll i = a; i <= n; i += 2)
#define s second
#define f first
#define noo cout << "NO" << nl;
#define yss cout << "YES" << nl;
#define em empty()
#define in insert
#define ers erase
// cout << fixed << setprecision(9);
using namespace std;
#define endcde return 0
const int mx = 1e5 + 1;
const int mod = 1e9 + 7;
ll kll()
{
    ll n;
    cin >> n;
    return n;
}
float fll()
{
    float n;
    cin >> n;
    return n;
}
string sll()
{
    string n;
    cin >> n;
    return n;
}
ll fact(ll n)
{
    if (n == 0)
        return 1;
    return n * fact(n - 1);
}
void murad()
{
    ll n=kll();
    cout<<0<<nl;
}
int main()
{
    int tstcse = 1;
    // cin >> tstcse;
    while (tstcse-- > 0)
    {
        murad();
    }
    endcde;
}
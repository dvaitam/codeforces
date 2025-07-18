//#include <bits/stdc++.h>

#include <iostream>

#include <vector>

#include <algorithm>

#include <deque>

#include <set>

#include <map>

#include <unordered_set>

#include <unordered_map>

#include <ext/pb_ds/assoc_container.hpp>

#include <ext/pb_ds/tree_policy.hpp>

using namespace std;

using namespace __gnu_pbds;





//KDK

//#define int long long

#define ff first

#define ss second

#define pb push_back

#define ppb pop_back

#define pf push_front

#define ppf pop_front

#define all(a) (a).begin(), (a).end()

#define rall(a) (a).rbegin(), (a).rend()

#define uset unordered_set

#define umap unordered_map

#define rs resize

#define sz(x) long(x.size())

#define mp make_pair

#define qmn(a, b) a = min(a, b)

#define qmx(a, b) a = max(a, b)

#define in freopen("filename.in", "r", stdin)

#define out freopen("filename.out", "w", stdout)

#define fastio cin.tie(nullptr)->sync_with_stdio(false)

#define YES "YES\n"

#define NO "NO\n"

typedef long long ll;

typedef unsigned long long ull;

typedef long double ld;

typedef pair<int, int> pii;

typedef tree<int, null_type, less<int>, rb_tree_tag, tree_order_statistics_node_update> ordered_set;

typedef tree<int, null_type, less_equal<int>, rb_tree_tag, tree_order_statistics_node_update> ordered_multiset;

typedef set<int> :: iterator iter;

typedef multiset<int> :: iterator miter;

typedef set<pii> :: iterator piter;

//mt19937 rnd(chrono::steady_clock::now().time_since_epoch().count());

//tree.insert(x);

//tree.find_by_order(k);

//tree.order_of_key(el);

const int INF = 2e9 + 1;

const int MOD = 998244353;

const int MXN = 3e5 + 2022;

const int MX = 2e5;

const int P = 179;

const int Q = 1791791791;





void solve()

{

    //solve here

    int a, b;

    cin >> a >> b;

    string s = "";

    while (a && b)

    {

        if (a > b)

        {

            s += '0';

            s += '1';

        }

        else

        {

            s += '1';

            s += '0';

        }

        a--;

        b--;

    }

    while (a--) s += '0';

    while (b--) s += '1';

    cout << s << "\n";

}





signed main()

{

    //write your code here

    fastio;

    int tests = 1;

    cin >> tests;

    while (tests--) solve();

    return 0;

}
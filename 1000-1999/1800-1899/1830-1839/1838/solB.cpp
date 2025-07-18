#include <bits/stdc++.h>
#include <ext/pb_ds/tree_policy.hpp>
#include <ext/pb_ds/assoc_container.hpp>

#define el '\n'
#define FIO ios_base::sync_with_stdio(0), cin.tie(0), cout.tie(0);

using namespace std;
using namespace __gnu_pbds;

typedef long long ll;
typedef unsigned uint;
typedef __int128 bint;
typedef long double ld;
typedef complex<ld> pt;
typedef unsigned long long ull;

template<typename T, typename X>
using hashTable = gp_hash_table<T, X>;
template<typename T>
using ordered_set = tree<T, null_type, less<T>, rb_tree_tag, tree_order_statistics_node_update>;
template<typename T>
using ordered_multiset = tree<T, null_type, less_equal<T>, rb_tree_tag, tree_order_statistics_node_update>;

// mt19937_64 for unsigned long long
mt19937 rng(std::chrono::system_clock::now().time_since_epoch().count());

void doWork()
{
    int n;
    cin >> n;
    vector<int> v(n);
    int pos1 = -1, pos2 = -1, posn = -1;
    for(int i = 0; i < n; i++)
    {
        cin >> v[i];
        pos1 = (v[i] == 1 ? i + 1 : pos1);
        pos2 = (v[i] == 2 ? i + 1 : pos2);
        posn = (v[i] == n ? i + 1 : posn);
    }

    if(min(pos1, pos2) < posn && posn < max(pos1, pos2))
        cout << pos1 << ' ' << pos2 << '\n';
    else if(posn < min(pos1, pos2))
        cout << posn << ' ' << min(pos1, pos2) << '\n';
    else
        cout << max(pos1, pos2) << ' ' << posn << '\n';
}

signed main()
{
    FIO
    int T = 1;
    cin >> T;
    for(int i = 1; i <= T; i++)
        doWork();
}
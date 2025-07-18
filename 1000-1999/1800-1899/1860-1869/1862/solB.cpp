#include <bits/stdc++.h>
#include <ext/pb_ds/assoc_container.hpp> // Common file
#include <ext/pb_ds/tree_policy.hpp>

using namespace std;
using namespace __gnu_pbds;
typedef long long int ll;
typedef vector<int> vi;
typedef pair<int, int> pi;

typedef tree<pair<int, int>, null_type,
             less<pair<int, int>>, rb_tree_tag,
             tree_order_statistics_node_update>
    ordered_set_pair;
#define PRIME 1000000007
#define REP(i, a, b) for (int i = a; i < b; ++i)
#define REPR(i, a, b) for (int i = a; i >= b; i--)
const int SIZE = 1e5 + 5;
int t = 1;
int t_orig;

void test_case() {
    int n;
    cin >> n;
    vi a(n);
    REP(i, 0, n) cin >> a[i];
    vi res;
    res.push_back(a[0]);
    REP(i, 1, n) {
        if (a[i] < a[i - 1]) {
            res.push_back(a[i]);
        }
        res.push_back(a[i]);
    }
    cout << res.size() << "\n";
    REP(i, 0, res.size()) cout << res[i] << " ";
}

int main()
{
    ios::sync_with_stdio(0);
    cin.tie(0);
    cin >> t;
    t_orig = t;
    while (t--)
    {
        test_case();
        cout << "\n";
    }
    return 0;
}
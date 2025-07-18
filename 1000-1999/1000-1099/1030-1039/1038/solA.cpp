#include<bits/stdc++.h>
#define maxn 1000001
#define oo 1048577
#define fi first
#define se second
#define vti vector <int>
#define pb(x) push_back(x)
#define pf(x) push_front(x)
#define reset(a) memset(a, 0, sizeof(a))
#define ll long long
#define pii pair<int, int>
#define MOD 1000000007
#define mp make_pair
#define PI acos(-1)
#define name ""

using namespace std;

int dem[27];
int n, k;
string s;
int minx = 1e9;
int main()
{
    ios_base::sync_with_stdio(0);
    cin.tie(0);
    cout.tie(0);

    //freopen(name".inp", "r", stdin);
    //freopen(name".out", "w", stdout);
    cin >> n >> k>> s;
    for(int i = 0; i < s.size(); i++)
        dem[ s[i] - 65 ]++;
    for(int i = 0; i < k; i++)
        minx = min(minx, dem[i]);
    cout << minx * k;
}
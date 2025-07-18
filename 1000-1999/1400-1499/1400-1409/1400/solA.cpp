#include <bits/stdc++.h>

using namespace std; 

#define endl "\n"

#define inf 0x3f3f3f3f

#define mod7 1000000007

#define mod9 998244353

#define rep(i,n,m) for(int i=n;i<m;i++)

#define pb push_back

#define debug(a) cout << "Debuging...|" << #a << ": " << a << "\n";

#define f first

#define s second

#define int long long

#define ld long double

#define pi pair<int,int>

#define pld pair<ld,ld>

typedef long long ll;



void solve()

{

    int n;

    string s;

    cin >> n >> s;

    string res;

    int t = 0;

    rep(i, 0, n){

        res += s[i + t];

        t ++;

    }

    cout << res << endl;

}

signed main()

{

    cin.tie(0)->sync_with_stdio(0);



    int _; cin >> _;

    while(_ --){

        solve();

    }

    return 0;

}
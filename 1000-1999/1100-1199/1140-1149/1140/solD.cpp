// time is tough.....so don't stop ...keep learning ...keep growing :)

#include <bits/stdc++.h>
#include<ext/pb_ds/assoc_container.hpp>
#include<ext/pb_ds/tree_policy.hpp>

using namespace std;
using namespace __gnu_pbds;

typedef pair<int,int> ii;
typedef long long ll;
typedef long double ld;
typedef unsigned long long ull;
typedef tree<pair<int,int>,null_type,less< pair<int,int> >,rb_tree_tag,tree_order_statistics_node_update > pbds;

#define debug(x) cout<<#x<<" :: "<<x<<endl;
#define debug2(x,y) cout<<#x<<" :: "<<x<<"\t"<<#y<<" :: "<<y<<endl;
#define debug3(x,y,z) cout<<#x<<" :: "<<x<<"\t"<<#y<<" :: "<<y<<"\t"<<#z<<" :: "<<z<<endl;
#define pb push_back
#define ft first
#define sd second
#define all(c)    (c).begin(),(c).end()
#define tr(c,i)   for(typeof((c).begin() i = (c).begin(); i != (c).end(); i++)
#define present(c,x)   ((c).find(x) != (c).end())
#define IOS ios_base::sync_with_stdio(false); cin.tie(0); cout.tie(0);

ll poww(ll a,ll b){ll res=1;while(b){if(b&1){res*=a;}a=a*a;b>>=1;}return res;}
ll poww(ll a,ll b,ll mod){ll res=1;while(b){if(b&1){res*=a;res%=mod;}a=a*a;a%=mod;b>>=1;}return res;}

const int INF = 2e9 ;
const int dx[] = {-1,0,0,1};
const int dy[] = {0,-1,1,0};
const int lg = 21;
const ll mod = 1e9 + 7;
const int N = 2e5 + 6 ;

signed main() {
    // IOS;
    ll n;
    cin >> n;
    ll ans =0 ;
    for(int i=2;i<n;i++) {
        ans += i*1ll*(i+1);
    }
    cout << ans << '\n';
    return 0;
}
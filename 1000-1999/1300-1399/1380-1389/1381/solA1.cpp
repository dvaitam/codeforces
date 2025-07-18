#include<bits/stdc++.h>
using namespace std;
//DataTypes
using ll  = long  long int;
using ld  = long double;

#define all(x)    x.begin(), x.end()
#define allr(x)   x.rbegin(),x.rend()
#define fl(i,n)   for(ll i = 0; i < n; ++i)

typedef vector<ll> vl;

//algo
#define Vmax(x)   *max_element(all(x))
#define Vmin(x)   *min_element(all(x))
#define Vsum(x)   accumulate(all(x),0ll)

//IO
#define nl      cout<< "\n";
#define ya      cout<<"YES\n";
#define na      cout << "NO\n";
#define inpt(v) fl(i,v.size()) cin >> v[i];
#define prt(v)  for(auto i:v) cout << i << " "; cout << "\n";
#define pr(x)   cout<<x;nl;
#define yn(ok)  cout << (ok?"Yes\n" :"No\n");
//Constants
const int M = 1e9+7; 
const int N = 2e5+10;
const ll INF = 9e18;
const ld eps = 1e-20;

ll n,m,k,a,b,c,q,x,y,l,r,ans,ans1,ans2,mx,mn,sum;

void solve()
{
    //START
    cin>>n;
    string s,t;
    cin>>s>>t;
    vl ans;
    if(s[0]!=t[0])
    {
        ans.push_back(1);
        s[0]=t[0];
    }
    for(int i=1;i<n;i++)
    {
        if(s[i]==t[i]) continue;
        ans.push_back(i+1);
        ans.push_back(1);
        ans.push_back(i+1);
        // s[i]=t[i];
        // if(s[i]==s[0])
        // {
        //     ans.push_back(i+1);
        //     ans.push_back(1);
        //     ans.push_back(i+1);
        //     s[i]=t[i];
        // }
        // else
        // {
        //     ans.push_back(1);
        //     ans.push_back(i+1);
        //     s[i]=t[i];
        // }
    }
    k=ans.size();
    if(k==0){pr(0);return;}
    cout<<k<<" ";
    prt(ans);
    //END
}
int main(){
    ios_base::sync_with_stdio(false); cin.tie(NULL);
    int t=1;
    cin >> t;
    for(int i=1;i<=t;i++){
        //cout<<"Case #"<<i<<": ";
        solve();
    }
    return 0;
}
/*__builtin_popcountll(x) , */
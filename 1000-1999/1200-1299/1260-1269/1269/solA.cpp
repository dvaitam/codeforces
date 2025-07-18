#include <iostream>
#include <climits>
#include <cstdio>
#include <cmath>
#include <vector>
#include <set>
#include <map>
#include <algorithm>
#include <string>
#include <queue>
#include <iomanip>
#include <cassert>
#define ins insert
#define er erase
#define pb push_back
#define epb emplace_back
#define mp make_pair
#define pr pair<int,int>
#define all(x) x.begin(),x.end()
#define rt return
#define f(x,y,z) for(x=y;x<z;++x)
#define rev(x,y,z) for(x=z-1;x>=y;--x)
#define pf(x) printf("%d\n",x)
#define pfs(s) printf("%s\n",s)
#define out(s) cout<<s<<"\n"
#define in(s) cin>>s
#define st string
#define X first
#define Y second
#define si set<int>
using namespace std;
typedef long long ll;typedef vector<int> vi;typedef vector<ll> vl;typedef vector<bool> vb;typedef vector<vi> vvi;
ll mod=1e9+7;
ll modpow(ll x, ll y, ll m){ll z=1;while(y){if(y&1)z=(z*x)%m;x=(x*x)%m;y>>=1;}return z;}
int egcd(int a, int b, int& x, int& y){if(a == 0){x= 0;y= 1;return b;}int x1, y1;int gcd= egcd(b%a, a, x1, y1);x= y1 - (b/a) * x1;y= x1;return gcd;}
string z[2]={"NO","YES"};
ll modinv(ll x,ll mod){rt modpow(x,mod-2,mod);}
int main()
{
    ios_base::sync_with_stdio(false);cin.tie(NULL);cout.tie(NULL);
    //cout.precision(8);cout<<fixed;
    int t=1;//cin>>t;
    int i,j,k,n,m;
    st s;
    while(t--)
    {
        //vi a(n);
        in(n);
        if(n==1)
            out("9 8");
            else
            out(3*n<<" "<<2*n);   
    }
    rt 0;
}
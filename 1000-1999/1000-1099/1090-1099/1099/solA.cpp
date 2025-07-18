#include <bits/stdc++.h>

using namespace std;
 
#define ull unsigned long long
#define ll long long
#define ld long double

#define pb push_back
#define mp make_pair
#define fi first
#define se second

#define MX LLONG_MAX
#define MN LLONG_MIN

const int N=2e5 + 5;
const ll mod=1e9+7;

ll powermodm(ll x,ll n,ll M){ll result=1;while(n>0){if(n % 2 ==1)result=(result * x)%M;x=(x*x)%M;n=n/2;}return result;}
ll power(ll _a,ll _b){ll _r=1;while(_b){if(_b%2==1)_r=(_r*_a);_b/=2;_a=(_a*_a);}return _r;}
ll gcd(ll a,ll b){if(a==0)return b;return gcd(b%a,a);}          
ll lcm(ll a,ll b){return (max(a,b)/gcd(a,b))*min(a,b);}

int main()                  
{
    ios_base::sync_with_stdio(false);
    cin.tie(NULL);
    ll n,h;
    cin>>n>>h;
    ll w1,h1;
    cin>>w1>>h1;
    ll w2,h2;
    cin>>w2>>h2;
    ll curr=h,ans=n;
    while(1)
    {
        if(curr==0)
            break;
        ans += curr;
        if(curr==h1)
        {
            ans -= w1;
            if(ans<0)
                ans = 0;
        }
        if(curr==h2)
        {
            ans -= w2;
            if(ans<0)
                ans = 0;
        }
        curr--;
    }
    ll co=0;
    cout<<max(ans,co);
    return 0;
}
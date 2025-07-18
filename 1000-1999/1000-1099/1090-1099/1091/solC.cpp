#include<bits/stdc++.h>
using namespace std;
#define ll long long int
#define pb(x) push_back(x)
#define f(i,st,n) for(i=st;i<n;i++)
#define qu queue
#define ps(x) push(x)
#define vec vector
#define ft() front()
#define itr vector<ll>::iterator
#define mst(a,b) memset(a,b,sizeof(a))
#define mod 1000000007
#define sfl(x) scanf("%lld",&x)
#define pfl(x) printf("%lld\n",x)
#define pf(x) printf("%lld",x)
ll gcd(ll a,ll b){
    if(a==0){
        return b;
    }
    return gcd(b%a,a);
}
ll power(ll x,ll y,ll m){
    ll res=1;
    while(y>0){
        if(y%2==1){
            res=(res%m*x%m)%m;
        }
        x=(x%m*x%m)%m;
        y=y/2;
    }
    return res;
}


int main()
{
ll n,i,x,y,j,l;
sfl(n);
vec<ll>v;
for(i=1;i*i<=n;i++){
if(n%i==0){
    if(i*i!=n){v.pb(i);v.pb(n/i);}
    else v.pb(i);
}
}
vec<ll>ans;
f(i,0,v.size()){
l=n+1-v[i];
x=(l-1);
x/=v[i];
x+=1;
y=(1+l);
y*=x;
y/=2;
ans.pb(y);
}
sort(ans.begin(),ans.end());
f(i,0,ans.size()){
cout<<ans[i]<<" ";
}

    return 0;
}
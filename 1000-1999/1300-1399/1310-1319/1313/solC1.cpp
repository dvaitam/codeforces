#include<iostream>
#include<math.h>
#include<algorithm>
#include<string>
#include<cstring>
#include<string.h>
#include<vector>
#include<map>
#include<cmath>
#include<cstdio>
#include<set>
#include<deque>
#include<queue>
#include<stack>
using namespace std;
const int inf=0x3f3f3f3f,hamod=1e9+7,HAmod=1e9+9,mod=200907,N=1e6+10,maxn=1e5+10;
const double pi=acos(-1.0);
typedef long long ll;
#define mp make_pair
#define eps 1e-8
#define lighting ios::sync_with_stdio(false);cin.tie(0),cout.tie(0);
//#define lighting ios::sync_with_stdio(false);cin.tie(nullptr),cout.tie(nullptr);
stack<ll>l,r;
ll n,maxx=-1,pos,m[N],suml[N],sumr[N];
signed main()
{
    lighting;
    cin>>n;
    for(ll i=1;i<=n;i++)
    {
        cin>>m[i];
        while(!l.empty()&&m[l.top()]>m[i]) l.pop();
        if(l.empty()) suml[i]=suml[0]+i*m[i];
        else suml[i]=suml[l.top()]+(i-l.top())*m[i];
        l.push(i);
    }
    for(ll i=n;i>=1;i--)
    {
        while(!r.empty()&&m[r.top()]>m[i]) r.pop();
        if(r.empty()) sumr[i]=sumr[n+1]+(n-i+1)*m[i];
        else sumr[i]=sumr[r.top()]+(r.top()-i)*m[i];
        r.push(i);
        if(suml[i]+sumr[i]-m[i]>maxx)
        {
            maxx=suml[i]+sumr[i]-m[i];
            pos=i;
        }
    }
    //cout<<maxx<<endl;
    for(ll i=pos-1;i>=1;i--)
    {
        m[i]=min(m[i],m[i+1]);
    }
    for(ll i=pos+1;i<=n;i++)
    {
        m[i]=min(m[i],m[i-1]);
    }
    for(ll i=1;i<=n;i++) cout<<m[i]<<" ";
}
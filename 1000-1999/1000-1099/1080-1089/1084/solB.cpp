#include <iostream>
#include <vector>
#include <stdio.h>
#include <bits/stdc++.h>
#include <algorithm>
#include <functional>
#include <queue>
using namespace std;

 
#define ull unsigned long long
#define ll long long
#define mod 1000000007
#define pb push_back
#define fi first
#define se second
#define PI acos(-1.0)
#define ld long double
#define sl(x) scanf("%lld",&x)
#define sl2(x,y) scanf("%lld %lld",&x,&y)
#define sl3(x,y,z) scanf("%lld %lld %lld",&x,&y,&z)
#define sld(x) scanf("%lf",&x)
#define sld2(x,y) scanf("%lf %lf",&x,&y)
#define sld3(x,y,z) scanf("%lf %lf %lf",&x,&y,&z)
#define FOR(i,a,b) for(ll i=a;i<=b;i++)
#define parr(a,n) FOR(i,0,n-1) printf("%lld ",a[i])
const int N=1e5 + 5;
const int M=1e5 + 5;
#define MX LLONG_MAX
#define MN LLONG_MIN
bool compare(const pair<int ,int > &a,const pair<int ,int > &b)
{
    if(a.second == b.second)
    {
        return a.first < b.first;
    }
    return a.second > b.second;
}
ll gcd(ll a, ll b){return __gcd(a,b);}
ll lcm(ll a, ll b){return (a*b)/gcd(a, b);}
ll power(ll a,ll b){ll p=1;for(ll i=0;i<b;i++){p = p * a;}return p;}  
 
using namespace std;
 
int main()                  
{
    ios_base::sync_with_stdio(false);
    cin.tie(NULL);
    ll n,s,ans=0;
    cin>>n>>s;
    ll a[n],sum=0;
    for(ll i=0;i<n;i++)
    {
        cin>>a[i];
        sum += a[i];
    }
    if(sum<s)
    {
        cout<<-1;
    }
    else
    {
        ll l=0,r=1e9;
        while(l<=r)
        {
            
            ll t = 0,flag=1,mn=MX;
            ll mid = (l+r)>>1;
            for(ll i=0;i<n;i++)
            {
                if(a[i]-mid>=0)
                {
                    t += (a[i] - mid);
                    mn = min(mn,mid);
                }
                else
                {
                    mn = min(mn,a[i]);
                }
            }
            //cout<<l<<" "<<r<<" "<<mid<<" "<<t<<" "<<mn<<" "<<ans<<endl;
            
                if(t>=s)
                {
                    l = mid+1;
                    ans = max(ans,mn);
                }
                else
                {
                    r = mid-1;
                }
                        
        }
        cout<<ans;
    }
    
    return 0;
}
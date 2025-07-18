#include <bits/stdc++.h>
using namespace std;
 #define ll long long
#define rep(i,a,n) for(int i=a;i<n;i++)
#define sf(n) scanf("%d",&n)
#define pii pair<int,int>
#define p2ii pair<pii,int>
#define ff first 
#define ss second 
int main()
{
    int n,k;sf(n);sf(k);int a[n];
    rep(i,0,n)sf(a[i]);int sum[n];
    sum[0]=a[0];double avg=0;
    rep(i,1,n)sum[i]=sum[i-1]+a[i];
    rep(range_element,k,n+1){
        int mx=sum[range_element-1];
        rep(j,range_element,n)
        {
            mx=max(mx,sum[j]-sum[j-range_element]);
        }
        avg=max(avg,(mx*1.0)/range_element);
    }
    cout<<fixed<<setprecision(9)<<avg<<"\n";
       
    
}
// LUOGU_RID: 168636739
#include<bits/stdc++.h>
#define int long long
using namespace std;int n,c[5001][5001],t;signed main(){cin>>t;while(t--){cin>>n;for(int i=0;i<=n;i++)for(int j=0;j<=i;j++)if(!j)c[i][j]=1;else c[i][j]=(c[i-1][j-1]+c[i-1][j])%1000000007;int ans=1;for(int i=1;i<=n;i++)for(int j=i+1;j<=2*i+1;j++)if(min(n,j-1)>=j-i-1&&max(0ll,n-j)>=2*i-j+1)ans=(ans+c[min(n,j-1)][j-i-1]*c[max(0ll,n-j)][2*i-j+1]%1000000007*j)%1000000007;cout<<ans<<endl;}return 0;}
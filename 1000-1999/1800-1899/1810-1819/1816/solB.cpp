#include<bits/stdc++.h>
#define int long long
#define endl '\n'
#define debug(x) cout<<#x<<"="<<x<<endl;
using namespace std;
void solve()
{
    int n;cin>>n;
    int arr[2][n];
    memset(arr,-1,sizeof arr);
    int mx=2*n,t=0;
    arr[0][0]=mx;
    mx--;
    arr[1][n-1]=mx;
    mx--;
    for(int i=n-2;i>0;i--)
    {
        arr[t][i]=mx--;
        t^=1;
    }
    mx=1,t=1;
    for(int i=0;i<n;i++)
    {
        arr[t][i]=mx++;
        t^=1;
    }
    for(int i=0;i<2;i++)
    {
        for(int j=0;j<n;j++)
            cout<<arr[i][j]<<" ";
        cout<<endl;
    }
    
    
    
}
signed main()
{
    ios::sync_with_stdio(false);
    cin.tie(0);
    int t;cin>>t;
    while(t--)
        solve();
}
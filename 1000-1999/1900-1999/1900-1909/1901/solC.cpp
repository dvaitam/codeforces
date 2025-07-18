#include <bits/stdc++.h>
#define int long long
using namespace std;
const int maxn=1e6+5;
int a[maxn];
vector<int>ans;
void solve()
{
    ans.clear();
    int n;
    cin>>n;
    for(int i=1;i<=n;i++)
        cin>>a[i];
    sort(a+1,a+n+1);
    int val=a[n]-a[1];
    while(a[n]!=a[1])
    {
        if(a[n]/2-a[1]/2<(a[n]+1)/2-(a[1]+1)/2)
        {
            ans.push_back(0);
            a[n]/=2;
            a[1]/=2;
        }
        else
        {
            ans.push_back(1);
            a[n]=(a[n]+1)/2;
            a[1]=(a[1]+1)/2;
        }
    }
    cout<<ans.size()<<"\n";
    if(ans.size()<=n)
    {
        for(auto e:ans)
            cout<<e<<" ";
        if(ans.size()!=0)
        cout<<"\n";
    }
}
signed main()
{
    ios::sync_with_stdio(false);
    cin.tie(0);
    cout.tie(0);
    int t=1;
    cin>>t;
    while(t--)
    {
        solve();
    }
    return 0;
}
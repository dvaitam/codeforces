#include<bits/stdc++.h>
using namespace std;
struct Node
{
    int id,val;
}a[1007];
bool cmp(Node a,Node b)
{
    return a.val<b.val;
}
int main()
{
    int n;cin>>n;
    for(int i=1;i<=n;i++)cin>>a[i].val,a[i].id=i;
    sort(a+1,a+1+n,cmp);
    int ans[1007];
    long long res=0;
    int t=1;
    for(int i=n;i>=1;i--)
    ans[t]=a[i].id,res+=(t-1)*a[i].val+1,t++;
    cout<<res<<endl;
    for(int i=1;i<=n;i++)
    cout<<ans[i]<<' ';
}
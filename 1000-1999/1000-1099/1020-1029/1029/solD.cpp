#include <bits/stdc++.h>
#define ft first
#define sd second
#define maxn 200005
#define mod
#define PI 3.14159265
#define MP make_pair
#define PB push_back
#define heap priority_queue
#define NguyenThiTham ""
using namespace std;
vector < long > v[12];
long n,k,a[maxn];
long dem(long u)
{
    long res=0;
    while(u)
    {
        res++;
        u/=10;
    }
    return res;
}
int main()
{
    ios_base::sync_with_stdio(0);cin.tie(0);cout.tie(0);
    cin >> n >> k;
    for(long i=1;i<=n;++i)
    {
        cin >> a[i];
        v[dem(a[i])].PB(a[i]%k);
    }
    for(long i=1;i<=10;++i) if(v[i].size()) sort(v[i].begin(),v[i].end());
    long long res=0;
    for(long i=1;i<=n;++i)
    {
        long long x=a[i]%k;
        for(long j=1;j<=10;++j)
        {
            x=(x*10)%k;
            if(v[j].size()==0) continue;
            //x=(x*10)%k;
            long y=(k-x)%k;
            long r=upper_bound(v[j].begin(),v[j].end(),y)-v[j].begin();
            long l=lower_bound(v[j].begin(),v[j].end(),y)-v[j].begin();
            res+=r-l;
            if(dem(a[i])==j && (a[i]%k)==y) res--;
        }
    }
    cout << res;
    return 0;
}
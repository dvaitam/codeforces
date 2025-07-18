#include<bitset>
#include<map>
#include<vector>
#include<cstdio>
#include<iostream>
#include<cstring>
#include<string>
#include<algorithm>
#include<cmath>
#include<stack>
#include<queue>
#include<set>
#define inf 0x3f3f3f3f
#define mem(a,x) memset(a,x,sizeof(a))

using namespace std;

typedef long long ll;
typedef unsigned long long ull;
typedef pair<int,int> pii;

inline int in()
{
    int res=0;char c;int f=1;
    while((c=getchar())<'0' || c>'9')if(c=='-')f=-1;
    while(c>='0' && c<='9')res=res*10+c-'0',c=getchar();
    return res*f;
}
const int N=200010;
int a[N];
int dp[N];

int main()
{
    int n=in(),m=in();
    for(int i=1;i<=n;i++)
    {
        a[i]=in();

    }
    for(int i=n;i>=1;i--)
    {
        if(a[i]==a[i+1]) dp[i]=dp[i+1]+1;

    }
    while(m--)
    {
        int l=in(),r=in(),x=in(),ans=-1;
        for(int i=l;i<=r;i+=dp[i]+1)
        {
            if(a[i]!=x)
            {
                ans=i;
                break;
            }
        }

        printf("%d\n",ans);
    }
    return 0;
}
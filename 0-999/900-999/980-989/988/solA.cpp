#include<cstdio>
#include<cstring>
#include<algorithm>
using namespace std;
const int maxn = 1e2 + 5;
int a[maxn];
int p[maxn];
int ans[maxn];
bool cmp(int x,int y)
{
    return a[x] < a[y];
}
int main()
{
    int n,k;
    scanf("%d %d",&n,&k);
    for(int i = 0;i < n; ++i)
    {
        scanf("%d",&a[i]);
        p[i] = i;
    }
    a[n] = -1;
    p[n] = n;
    sort(p,p + n,cmp);
    int sum = 0;
    for(int i = 1;i <= n; ++i)
    {
        if(a[p[i]] != a[p[i - 1]]) ans[sum++] = p[i - 1] + 1;
    }
    if(sum >= k)
    {
        printf("YES\n");
        for(int i = 0;i < k; ++i)
        {
            printf("%d",ans[i]);
            if(i != k) printf(" ");
        }
        printf("\n");
    }
    else printf("NO\n");
    return 0;
}
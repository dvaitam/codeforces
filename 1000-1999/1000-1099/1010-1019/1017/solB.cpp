#include<cstdio>
#include<cstring>
#include<cctype>
#include<algorithm>
#define maxn 100005
#define LL long long
using namespace std;

int n,cnt[2][2];
char a[maxn],b[maxn];
LL ans=0;

int main()
{
	scanf("%d",&n);
    scanf("%s",a+1);
    scanf("%s",b+1);
    for(int i=1;i<=n;i++)
		cnt[a[i]-'0'][b[i]-'0']++;
	for(int i=1;i<=n;i++)
		if(b[i]=='0')
			ans += cnt[(int(a[i]-'0'))^1][1] + cnt[(int(a[i]-'0'))^1][0];
		else
			ans += cnt[(int(a[i]-'0'))^1][0];
	printf("%lld\n",ans/2);
}
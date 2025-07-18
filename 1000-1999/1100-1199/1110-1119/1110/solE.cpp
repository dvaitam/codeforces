#include<cstdio>
#include<algorithm>
#define N 100010
using namespace std;
int a[N],n,c[N];

int read()
{
    int ans=0,fu=1;
    char j=getchar();
    for (;j<'0' || j>'9';j=getchar()) if (j=='-') fu=-1;
    for (;j>='0' && j<='9';j=getchar()) ans*=10,ans+=j-'0';
    return ans*fu;
}

int main()
{
    n=read();
    for (int i=1;i<=n;i++)
	a[i]=read();
    for (int i=1;i<=n;i++)
	c[i]=read();
    if (a[1]!=c[1] || a[n]!=c[n])
    {
	printf("No");
	return 0;
    }
    for (int i=n;i>1;i--)
	a[i]-=a[i-1],c[i]-=c[i-1];
    sort(a+1,a+n+1);
    sort(c+1,c+n+1);
    for (int i=1;i<=n;i++)
	if (a[i]!=c[i])
	{
	    printf("No");
	    return 0;
	}
    /*for (int i=1;i<=n;i++) printf("%d ",a[i]);
    putchar('\n');
    for (int i=1;i<=n;i++) printf("%d ",c[i]);*/
    printf("Yes");
    return 0;
}
#include<cstdio>

#include<cstring>

#include<cstdlib>

#include<cmath>

#include<iostream>

#include<algorithm>  

#include<queue>

#include<map>

#include<set>

#include<utility>

#include<vector>

#define N 100005

#define inf 0x7fffffff

#define C 1000000000LL

#define LL long long

using namespace std;

inline int ra()

{

    int x=0,f=1; char ch=getchar();

    while (ch<'0' || ch>'9') {if (ch=='-') f=-1; ch=getchar();}

    while (ch>='0' && ch<='9') {x=x*10+ch-'0'; ch=getchar();}

    return x*f;

}

int a[1005],b[1005],c[1005],cnt,tot;

void pre()

{

	for (int i=2; i<=1000; i++)

	{

		bool flag=0;

		for (int j=2; j<=sqrt(i); j++)

			if (i%j==0)

			{

				flag=1;

				break;

			}

		if (!flag) c[++tot]=i;

	}

}

int main()

{

	int n=ra(),ans=0;

	pre();

	for (int i=1; i<=n; i++)

		a[i]=ra();

	for (int i=1; i<n; i++)

		if (__gcd(a[i],a[i+1])!=1)

		{

			for (int j=1; j<=tot; j++)

				if (__gcd(a[i],c[j])==1 && __gcd(c[j],a[i+1])==1)

				{

					b[i]=c[j];

					break;

				}

			++ans;

		}

	cout<<ans<<endl;

	for (int i=1; i<=n; i++)

	{

		printf("%d ",a[i]);

		if (b[i]) printf("%d ",b[i]);

	}

}
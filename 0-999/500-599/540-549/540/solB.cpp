#include<iostream>

#include<cstdio>

#include<algorithm>

#include<cmath>

#define maxn 1005

#define inf 0x7FFFFFFE

using namespace std;

typedef long long ll;

template<typename T>

inline void _r(T &d){char t=getchar();bool f=false;while(t<48||t>57){if(t==45)f=true;t=getchar();}for(d=0;t>=48&&t<=57;t=getchar())d=d*10+t-48;if(f)d=-d;}

int n,k,p,x,y,a[maxn],b[maxn],tot,cur,siz;

int main()

{

	int i;

	_r(n);_r(k);_r(p);_r(x);_r(y);

	for(i=1;i<=k;i++)

	{

		_r(a[i]);tot+=a[i];

	}

	if(tot+(n-k)>x)puts("-1\n");

	else

	{

		cur=x-tot-(n-k);

		for(i=k+1;i<=n;i++)

		{

			if(cur>=y-1)b[++siz]=y,cur-=y-1;

			else b[++siz]=1;

			a[i]=b[siz];

		}

		sort(a+1,a+1+n);

		if(a[n/2+1]<y||a[n]>p)puts("-1\n");

		else

		{

			for(i=1;i<=siz;i++)printf("%d ",b[i]);

			puts("\n");

		}

	}

	return 0;

}
#include<cstdio>

using namespace std;

typedef long long ll;

const int N=1001000;

int i,j,k,n,ch,tt,l;

struct cc {int x,y;} A[N];

ll t,a[N],b[N];

void R(int &x)

{

	x=0;ch=getchar();

	while (ch<'0' || '9'<ch) ch=getchar();

	while ('0'<=ch && ch<='9') x=x*10+ch-'0',ch=getchar();

}

void W(int x)

{

	if (x>=10) W(x/10);

	putchar(x%10+'0');

}

int main()

{

	R(n);

	for (i=1;i<=n;i++)

	{

		R(k);A[i].x=-1;

		a[i]=a[i-1]+k;

	}

	for (i=1;i<=n;i++)

	{

		R(k);

		b[i]=b[i-1]+k;

	}

	l=0;

	if (a[n]<=b[n])

	{

		for (i=1;i<=n;i++)

		{

			while (l<n && b[l+1]<=a[i]) l++;

			tt=a[i]-b[l];

			if (A[tt].x>-1)

			{

				W(i-A[tt].x);puts("");

				for (j=A[tt].x+1;j<=i;j++) W(j),putchar(' ');

				puts("");

				W(l-A[tt].y);puts("");

				for (j=A[tt].y+1;j<=l;j++) W(j),putchar(' ');

				puts("");

				return 0;

			}

			A[tt].x=i;

			A[tt].y=l;

		} 

	}

	else

	{

		for (i=1;i<=n;i++)

		{

			while (l<n && a[l+1]<=b[i]) l++;

			tt=b[i]-a[l];

			if (A[tt].x>-1)

			{

				W(l-A[tt].y);puts("");

				for (j=A[tt].y+1;j<=l;j++) W(j),putchar(' ');

				puts("");

				W(i-A[tt].x);puts("");

				for (j=A[tt].x+1;j<=i;j++) W(j),putchar(' ');

				puts("");

				return 0;

			}

			A[tt].x=i;

			A[tt].y=l;

		} 

	}

}
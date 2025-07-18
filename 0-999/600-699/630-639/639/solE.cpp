#include<cstdio>

#include<algorithm>

using namespace std;

typedef long long ll;

const int N=150015;

int i,j,k,n,ch;

ll T,suml,sumr;

double l,r,mid;

struct cc {int p,t;ll l,r;} A[N];

bool cmp(const cc &a,const cc &b)

{

	return 1ll*a.t*b.p<1ll*b.t*a.p;

}

bool cmp_p(const cc &a,const cc &b)

{

	return a.p<b.p;

}

void R(int &x)

{

	x=0;ch=getchar();

	while (ch<'0' || '9'<ch) ch=getchar();

	while ('0'<=ch && ch<='9') x=x*10+ch-'0',ch=getchar();

}

bool check(double x)

{

	double max1=0.0,max2=0.0,t;

	int i,j,k;

	for (i=1;i<=n;i++)

	{

		j=i;t=A[i].p*(1.0-x*A[i].r/(1.0*T));

		if (max1>t+1e-12) return 0;

		t=A[i].p*(1.0-x*A[i].l/(1.0*T));

		if (t>max2) max2=t;

		while (A[j].p==A[j+1].p)

		{

			j++;

			t=A[j].p*(1.0-x*A[j].r/(1.0*T));

			if (max1>t+1e-12) return 0;

			t=A[j].p*(1.0-x*A[j].l/(1.0*T));

			if (t>max2) max2=t;

		}

		max1=max2;i=j;

	}

	return 1;

}

int main()

{

	R(n);

	for (i=1;i<=n;i++) R(A[i].p);

	for (i=1;i<=n;i++) R(A[i].t),T+=A[i].t;

	sort(A+1,A+n+1,cmp);

	suml=sumr=0;

	for (i=1;i<=n;i++)

	{

		j=i;sumr+=A[i].t;

		while (j<n && 1ll*A[j].t*A[j+1].p==1ll*A[j+1].t*A[j].p) sumr+=A[++j].t;

		for (k=i;k<=j;k++) A[k].l=suml+A[k].t,A[k].r=sumr;

		suml=sumr;i=j;

	}

	sort(A+1,A+n+1,cmp_p);

	l=0.0;r=1.0;

	for (i=1;i<=50;i++)

	{

		mid=0.5*(l+r);

		if (check(mid)) l=mid;

		else r=mid;

	}

	printf("%.12lf\n",l);

}
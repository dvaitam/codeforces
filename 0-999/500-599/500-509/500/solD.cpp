/*

._._._._._._._._._._._._._._._._._._._._._._._._._._._._._._._._._._._._._._._._

	ABHINANDAN AGARWAL

	MNNIT ALLAHABAD

	CSE

._._._._._._._._._._._._._._._._._._._._._._._._._._._._._._._._._._._._._._._._

*/

#include<bits/stdc++.h>

using namespace std;

#define pc putchar

#define gc getchar

typedef long long ll;

typedef unsigned long long llu;

#define mp make_pair

#define pb push_back

#define sc(x) scanf("%d",&x);

#define fl(i,n) for (i=0;i<n;i++)

#define fle(i,n) for (i=1;i<=n;i++)

#define fla(i,a,n) for (i=a;i<n;i++)

#define mem(a,i) memset(a,i,sizeof(a))

#define fi first

#define se second

#define pii pair<int,int> 

#define piii pair<int,pair<int,int> >

#define wl(n) while (n--)



int fs()

{

	int x=0;

	char c;

	c=getchar();

	while (c<'0'||c>'9')

		c=getchar();

	while ('0'<=c&&c<='9')

	{

		x=(x<<3)+(x<<1)+c-'0';

		c=getchar();

	}

	return x; 

}

void pps(const char *s)

{

	int i;

	for (i=0;s[i]!='\0';i++)

	{

		putchar(s[i]);

	}

}

void _ppd(int a)

{

	if (a==0)

		return ;

	_ppd(a/10);

	putchar(a%10+'0');

}

void ppd(int a) // Printing integer using utchar unlocked!

{

	if (a==0)

	{

		pc('0');

		return;

	}

	if (a<0)

	{

		pc('-');

		a=-a;

	}

	_ppd(a);

}

void _pplld(ll a)

{

	if (a==0)

		return ;

	_pplld(a/10);

	pc(a%10+'0');

}

void pplld(ll a)

{

	if (a==0)

	{

		pc('0');

		return ;

	}

	if (a<0)

	{

		pc('-');

		a=-a;

	}

	_pplld(a);

}

int ggs(char *s) // String inputting

{

	int x=0;

	char c=gc();

	while (!('a'<=c&&c<='z')&&c!=-1) // change conditions to whatever strings can input !

		c=gc();

	while ('a'<=c&&c<='z')

	{

		s[x++]=c;

		c=gc();

	}

	s[x]='\0';

	return x;

}



int fsn()

{

	int x=0;

	char c;

	c=getchar();

	while (!('0'<=c&&c<='9'||c=='-'))

		c=getchar();

	int neg=0;

	if (c=='-'){

		neg=1;

		c=getchar();

	}

	while ('0'<=c&&c<='9')

	{

		x=(x<<3)+(x<<1)+c-'0';

		c=getchar();

	}

	return (neg==0?x:-x); 

}

//------------------------------------------------------------------------------------

int n;



typedef long double LF;

vector<pii>gr[100010];

double hag[100010]; // permutations

ll we[100010];

int sz[100010];ll poop;

int dfs(int c,int p)

{

	int i,z=gr[c].size();

	ll ret=1;

	for (i=0;i<z;i++)

	{

		if (gr[c][i].fi!=p)

		{	

			ll opar=n,neeche=0;

			neeche=dfs(gr[c][i].fi,c);

			opar=n-neeche;

			hag[gr[c][i].se]=(opar*(opar-1)/2)*(neeche)/(LF)poop+(neeche*(neeche-1)/2)*(LF)opar/(LF)poop;

			hag[gr[c][i].se]*=2.0;

			ret+=neeche;

		}

	}

	

	return (int)ret;

}

int main()

{

	int i,j;

	n=fs();

	poop=(ll)n*((ll)n-1)*((ll)n-2)/(ll)6;

	for (i=0;i<n-1;i++)

	{

		int x,y,w;

		x=fs();y=fs();w=fs();

		we[i]=w;

		gr[x].pb(mp(y,i));

		gr[y].pb(mp(x,i));

	}

	dfs(1,-1);

	

	LF ans=0;

	for (i=0;i<n-1;i++)

	{

		ans+=(LF)hag[i]*we[i];

	}

	/*for (i=0;i<n-1;i++)

	{

		printf("%d: %lld %lld\n",i,hag[i],we[i]);

	}*/

	int q;

	q=fs();

	wl(q)

	{

		int x,nw;

		x=fs();nw=fs();x--;

		ans-=we[x]*hag[x];

		we[x]=nw;

		ans+=we[x]*hag[x];

		printf("%.15lf\n",(double)ans);

	}



	

	return 0;

}
#include<bits/stdc++.h>
#define tr(i,l,r) for(i=(l);i<=(r);i++)
#define rtr(i,r,l) for(i=(r);i>=(l);i--)
#define ll long long
#define oo 0x7F7F7F7F
#define maxn 200050
using namespace std;

int tt,n,i,x,xx;
bool p[maxn];
char s[maxn];

bool check(int x)
{
	int l,x2=0;
	while (x) {x2+=x%10; x/=10;}
	x=x2; x2=0;
	while (x) {x2+=x%10; x/=10;}
	return(x2<10);
}
int main()
{
	scanf("%d%s",&n,s+1);
	//freopen("t.txt","w",stdout);
	//srand(2); n=800; tr(i,1,n) s[i]=rand()%2+'8';
	x=0; tr(i,1,n) x+=s[i]-'0';
	xx=x;
	if (!check(x))
	{
		memset(p,0,sizeof(p));
		x=xx;
		tr(i,1,n) if ((i&1)&&(i+1<=n))
		{
			p[i]=1;
			x+=(s[i]-'0')*9;
			if (check(x)) break;
		}
	}
	if (!check(x))
	{
		memset(p,0,sizeof(p));
		x=xx;
		tr(i,1,n) if ((!(i&1))&&(i+1<=n))
		{
			p[i]=1;
			x+=(s[i]-'0')*9;
			if (check(x)) break;
		}
	}
	if (!check(x))
	{
		memset(p,0,sizeof(p));
		x=xx;
		rtr(i,n,1) if ((i&1)&&(i+1<=n))
		{
			p[i]=1;
			x+=(s[i]-'0')*9;
			if (check(x)) break;
		}
	}
	if (!check(x))
	{
		memset(p,0,sizeof(p));
		x=xx;
		rtr(i,n,1) if ((!(i&1))&&(i+1<=n))
		{
			p[i]=1;
			x+=(s[i]-'0')*9;
			if (check(x)) break;
		}
	}
	//
	if (!check(x))
	{
		memset(p,0,sizeof(p));
		x=xx;
		tr(i,1,n) if ((i%3==1)&&(i+2<=n))
		{
			p[i]=1; p[i+1]=1;
			x+=(s[i]-'0')*99+(s[i+1]-'0')*9;
			if (check(x)) break;
		}
	}
	if (!check(x))
	{
		memset(p,0,sizeof(p));
		x=xx;
		tr(i,1,n) if ((i%3==2)&&(i+2<=n))
		{
			p[i]=1; p[i+1]=1;
			x+=(s[i]-'0')*99+(s[i+1]-'0')*9;
			if (check(x)) break;
		}
	}
	if (!check(x))
	{
		memset(p,0,sizeof(p));
		x=xx;
		rtr(i,n,1) if ((i%3==1)&&(i+2<=n))
		{
			p[i]=1; p[i+1]=1;
			x+=(s[i]-'0')*99+(s[i+1]-'0')*9;
			if (check(x)) break;
		}
	}
	if (!check(x))
	{
		memset(p,0,sizeof(p));
		x=xx;
		rtr(i,n,1) if ((i%3==2)&&(i+2<=n))
		{
			p[i]=1; p[i+1]=1;
			x+=(s[i]-'0')*99+(s[i+1]-'0')*9;
			if (check(x)) break;
		}
	}
	//
	assert(check(x));
	tr(i,1,n) 
	{
		putchar(s[i]);
		if (!p[i]&&i<n) putchar('+');
	}
	puts("");
	
	sprintf(s+1,"%d",x); n=strlen(s+1);
	x=0;
	tr(i,1,n) 
	{
		putchar(s[i]);
		x+=s[i]-'0';
		if (i<n) putchar('+'); else puts("");
	}
	
	sprintf(s+1,"%d",x); n=strlen(s+1);
	x=0;
	tr(i,1,n) 
	{
		putchar(s[i]);
		x+=s[i]-'0';
		if (i<n) putchar('+'); else puts("");
	}
	
	return 0;
}
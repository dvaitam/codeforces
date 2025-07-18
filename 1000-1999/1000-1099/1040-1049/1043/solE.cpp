#include <bits/stdc++.h>
using namespace std;

//#define getc getc_unlocked
//#define putc putc_unlocked
#define eb emplace_back
#define lb lower_bound
#define ub upper_bound
#define mp make_pair
#define ff first
#define ss second

typedef long long int ll;
typedef long double ld;
typedef short int sh;

inline void readI(int *i)
{
	int t=0;
	register char z='a';
	int znak=1;
	z=getc(stdin);
	if (z=='-')
	{
		znak=-1;
	}
	while ((z<'0') || ('9'<z))
	{
		z=getc(stdin);
		if (z=='-')
		{
			znak=-1;
		}
	}
	while (('0'<=z) && (z<='9'))
	{
		t=(t<<3)+(t<<1)+z-'0';
		z=getc(stdin);
	}
	*i=(t*znak);
}
inline void readUI(int *i)
{
	int t=0;
	register char z='a';
	z=getc(stdin);
	while ((z<'0') || ('9'<z))
	{
		z=getc(stdin);
	}
	while (('0'<=z) && (z<='9'))
	{
		t=(t<<3)+(t<<1)+z-'0';
		z=getc(stdin);
	}
	*i=t;
}
inline void readL(ll *l)
{
	ll t=0;
	register char z='a';
	int znak=1;
	z=getc(stdin);
	if (z=='-')
	{
		znak*=-1;
	}
	while ((z<'0') || ('9'<z))
	{
		z=getc(stdin);
		if (z=='-')
		{
			znak=-1;
		}
	}
	while (('0'<=z) && (z<='9'))
	{
		t=(t<<3)+(t<<1)+z-'0';
		z=getc(stdin);
	}
	*l=(t*znak);
}
inline void readUL(ll *l)
{
	ll t=0;
	register char z='a';
	z=getc(stdin);
	while ((z<'0') || ('9'<z))
	{
		z=getc(stdin);
	}
	while (('0'<=z) && (z<='9'))
	{
		t=(t<<3)+(t<<1)+z-'0';
		z=getc(stdin);
	}
	*l=t;
}
inline void writeI(int i)
{
	if (i==0)
	{
		putc(48, stdout);
	}
	else
	{
		if (i<0)
		{
		   i*=-1;
		   putc(45, stdout);
		}
		int _tab[12];
		int wsk=0;
		while (i>0)
		{
			++wsk;
			_tab[wsk]=(i%10)+48;
			i/=10;
		}
		for (int j=wsk; j>=1; --j)
		{
			putc(_tab[j], stdout);
		}
	}
}
inline void writeL(ll l)
{
	if (l==0)
	{
		putc(48, stdout);
	}
	else
	{
		if (l<0)
		{
		   l*=-1;
		   putc(45, stdout);
		}
		int _tab[21];
		int wsk=0;
		while (l>0)
		{
			++wsk;
			_tab[wsk]=(l%10)+48;
			l/=10;
		}
		for (int j=wsk; j>=1; --j)
		{
			putc(_tab[j], stdout);
		}
	}
}
inline void writeS(string s)
{
	int l=s.length();
	for (int i=0; i<l; ++i)
	{
		putc(s[i], stdout);
	}
}
inline void writeC(char c)
{
	putc(c, stdout);
}
inline void space()
{
	putc(32, stdout);
}
inline void endl()
{
	putc(10, stdout);
}

#define debug if(0)
#define debug2 if(1)
#define debug3 if(1)
#define debug4 if(1)
#define MAXN 300005
#define ZAK 1048580

int n, m, r, c, d;
ll a[MAXN], b[MAXN];
ll wynik[MAXN];
pair <ll, int> tab[MAXN];
ll pref[MAXN], suf[MAXN];

int main()
{
	readI(&n);
	r=1;
	while (r<n)
	{
		r<<=1;
	}
	readI(&m);
	for (int i=1; i<=n; ++i)
	{
		readL(&a[i]);
		readL(&b[i]);
		tab[i]=mp(b[i]-a[i], i);
	}
	sort(tab+1, tab+1+n);
	for (int i=1; i<=n; ++i)
	{
		pref[i]=pref[i-1]+b[tab[i].ss];
	}
	for (int i=n; i>=1; --i)
	{
		suf[i]=suf[i+1]+a[tab[i].ss];
	}
	for (int i=1; i<=n; ++i)
	{
		int nr=tab[i].ss;
		wynik[nr]=b[nr]*(n-i);
		wynik[nr]+=a[nr]*(i-1);
		wynik[nr]+=pref[i-1];
		wynik[nr]+=suf[i+1];
	}
	for (int i=1; i<=m; ++i)
	{
		readI(&c);
		readI(&d);
		ll wyn=min(a[c]+b[d], a[d]+b[c]);
		wynik[c]-=wyn;
		wynik[d]-=wyn;
	}
	for (int i=1; i<=n; ++i)
	{
		writeL(wynik[i]);
		space();
	}
	return 0;
}
/*
*/
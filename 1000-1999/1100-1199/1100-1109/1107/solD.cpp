#include <bits/stdc++.h>
#define FOR(i, a, b) for (int i=a; i<(b); i++)
#define FORB(y, m) for (int y,_=m; y=-_&_; _^=y)
#define PPC(x) __builtin_popcount(x)
#define ALL(x) (x).begin(), (x).end()
#define INT(x)	int x; scanf ("%d", &x)
#define LL(x)	long long x; scanf ("%lld", &x)
#define pb push_back
using namespace std;

const int maxN = 5242;

bool tofail[16];
int ps[16], pln[16], ss[16];
char in[maxN][maxN / 4];

int todig(char a)
{
	if (a >= '0' and a <= '9')
		return a-'0';
	return (a - 'A') + 10;
}

void init()
{
	tofail[2] = true;
	tofail[4] = true;
	tofail[5] = true;
	tofail[6] = true;
	tofail[9] = true;
	tofail[10] = true;
	tofail[11] = true;
	tofail[13] = true;
	
	ps[0] = ps[1] = ps[3] = ps[7] = 0;
	ps[8] = ps[12] = ps[14] = ps[15] = 1;
	pln[0] = pln[15] = 0;
	pln[1] = pln[14] = 3;
	pln[3] = pln[12] = 2;
	pln[7] = pln[8] = 1;
	
	ss[0] = ss[8] = ss[12] = ss[14] = 0;
	ss[1] = ss[3] = ss[7] = ss[15] = 1;
}

void fail()
{
	printf("1\n");
	exit(0);
}

int main()
{
	init();
	INT(n);	
	int res = 0, rows = 0;
	in[0][1] = '#';
	
	FOR(i, 1, n+1)
		scanf ("%s", &in[i][1]);
	FOR(i, 1, n+1)
	{
		rows++;
		if (i == n or strcmp(&in[i][1], &in[i+1][1]) != 0)
			res = __gcd(res, rows), rows = 0;
		
		int len = 0, sn = ps[todig(in[i][1])];
		
		FOR(j, 1, n/4+1)
		{
			int cs = todig(in[i][j]);
				
			if (tofail[cs])
				fail();
			
			if (cs == 0 and sn == 0)
			{
				len += 4;
				continue;
			}
			
			if (cs == 15 and sn == 1)
			{
				len += 4;
				continue;
			}
			
			if (sn == ps[cs])
				len += pln[cs];
			res = __gcd(res, len);
			len = 4 - pln[cs];
			sn = ss[cs];
		}
		
		res = __gcd(res, len);
	}
	
	printf("%d\n", res);			
	return 0;
}
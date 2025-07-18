#include<cstdio>

#include<cstring>

#include<algorithm>

#include<cstdlib>

#include<map>

#include<queue>

#include<iterator>

using namespace std;

#define FOR(i,s,e) for(int i = (s); i < (e); i++)

#define FOE(i,s,e) for(int i = (s); i <= (e); i++)

#define FOD(i,s,e) for(int i = (s); i >= (e); i--)

#define CLR(a) memset(a,0,sizeof(a))

#define ll long long

#include<ctime>

#include<cmath>

#include<vector>

#include<iostream>



int n, m, x, y, z, k, w, a, b;

int l[105], LINK[10005], ed[10005], v[105];

ll L[10005], G[10005], A[105], tmp, tmp2;

int B[10005], vis[105];

bool ok;



ll gcd(ll a, ll b) { if (a > b) return gcd(b, a); if (a == 0ll) return b; return gcd(b % a, a); }

ll lcm(ll a, ll b) { ll c = a * b; c /= gcd(a, b); return c; }



void bd(int x, int y, int a, int b) { 

	LINK[w] = l[x]; l[x] = w; ed[w] = y; G[w] = (ll)(a); L[w] = (ll)(b); w++;

}



void DFS(int x)

{

	v[x] = 1; vis[k++] = x;

	for (int i = l[x]; i != -1; i = LINK[i])

	{

		if (v[ed[i]] == 1)

		{

			if (lcm(A[x], A[ed[i]]) != L[i] || gcd(A[x], A[ed[i]]) != G[i])

				ok = 0;

			continue;

		}

		if (L[i] % A[x] != 0ll) { ok = 0; continue; }

		A[ed[i]] = (L[i] / A[x]) * G[i];

		DFS(ed[i]);

	}

}	



int main ()

{

	memset(l, -1, sizeof(l));

	scanf("%d %d", &n, &m);

	

	FOR(i, 0, m) 

	{

		scanf("%d %d %d %d", &x, &y, &a, &b); 

		bd(x, y, a, b);

		bd(y, x, a, b);

	}

	

	memset(A, 0, sizeof(A));

	

	FOE(i, 1, n)

	{

		if (v[i]) continue;

		if (l[i] == -1) continue;

		

		ok = 1;

		

		for (int j = l[i]; j != -1; j = LINK[j])

		{

			if (j == l[i]) tmp = L[j];

			else tmp = gcd(tmp, L[j]);

		}

		

		for (int j = l[i]; j != -1; j = LINK[j])

		{

			if (j == l[i]) tmp2 = G[j];

			else tmp2 = lcm(tmp2, G[j]);

			if (tmp2 > tmp) { ok = 0; break; }

		}

		

		if (ok == 0) { printf("NO\n"); return 0; }

		

		w = 0;

		for (int j = 1; j * j <= (int)(tmp); j++)

		{

			if ((int)(tmp) % j == 0) 

			{

				if (j % (int)(tmp2) == 0) B[w++] = j;

				if (j != (int)(tmp) / j && ((int)(tmp) / j) % tmp2 == 0ll) 

					B[w++] = (int)(tmp) / j;

			}

		}

		

		ok = 0;

		

		FOR(j, 0, w)

		{

			A[i] = (ll)(B[j]);

			ok = 1; k = 0;

			DFS(i);

			if (ok == 1) break;

			FOR(ii, 0, k) v[vis[ii]] = 0;

		}

		

		if (ok == 0) { printf("NO\n"); return 0; }

	}

	

	FOE(i, 1, n) if (A[i] == 0ll) A[i] = 1ll;

	

	printf("YES\n");

	FOE(i, 1, n) printf("%I64d ", A[i]);

	printf("\n");

	

	return 0;

}
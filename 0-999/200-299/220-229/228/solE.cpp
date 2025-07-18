#include <bits/stdc++.h>



#define FI(i,a,b) for(int i=(a);i<=(b);i++)

#define FD(i,a,b) for(int i=(a);i>=(b);i--)



#define LL long long

#define Ldouble long double

#define PI 3.1415926535897932384626



#define PII pair<int,int>

#define PLL pair<LL,LL>

#define mp make_pair

#define fi first

#define se second



using namespace std;



int n, m, par[105], dist[105];

bool use[105];



int fp(int x){

	if(par[x] == x) return x;

	else{

		int pa = par[x];

		par[x] = fp(pa);

		dist[x] ^= dist[pa];

		return par[x];

	}

}

int main(){

	scanf("%d %d",&n,&m);

	FI(i,1,n) par[i] = i;

	FI(i,1,m){

		int a, b, c;

		scanf("%d %d %d",&a,&b,&c);

		c ^= 1;

		int pa = fp(a), pb = fp(b);

		if(pa == pb){

			if(dist[a] ^ dist[b] ^ c){

				printf("Impossible\n");

				return 0;

			}

			continue;

		}

		par[pa] = pb;

		dist[pa] = c ^ dist[a] ^ dist[b];

	}

	

	int cnt = 0;

	

	FI(i,1,n){

		fp(i);

		if(dist[i] & 1) use[i] = 1, cnt++;

	}

	

	printf("%d\n",cnt);

	FI(i,1,n) if(use[i]){

		cnt--;

		printf("%d%c",i,cnt?' ':'\n');

	}

	return 0;

}
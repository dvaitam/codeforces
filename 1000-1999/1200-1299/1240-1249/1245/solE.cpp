#include "bits/stdc++.h"
using namespace std;
typedef long long ll;
const int INF = (1<<30);
const ll INFLL = (1ll<<60);
const ll MOD = (ll)(1e9+7);

#define l_ength size

void mul_mod(ll& a, ll b){
	a *= b;
	a %= MOD;
}

void add_mod(ll& a, ll b){
	a = (a<MOD)?a:(a-MOD);
	b = (b<MOD)?b:(b-MOD);
	a += b;
	a = (a<MOD)?a:(a-MOD);
}

int h[10][10],b[10][10],x[100],y[100];
bool done[100][2];
long double memo[100][2];

long double solve(int p, int k){
	int s;
	long double c=6.0;
	if(done[p][k]){
		return memo[p][k];
	}
	done[p][k] = true;
	if(!(p||k)){
		return memo[p][k] = 0.0;
	}
	if(!k){
		memo[p][k] = 6.0;
		for(s=1; s<7; ++s){
			if(p-s<0){
				c -= 1.0;
			}else{
				memo[p][k] += solve(p-s,1);
			}
		}
		memo[p][k] /= c;
		return memo[p][k];
	}
	memo[p][k] = solve(p,0);
	if(h[x[p]][y[p]]){
		memo[p][k] = min(memo[p][k],solve(b[x[p]-h[x[p]][y[p]]][y[p]],0));
	}
	return memo[p][k];
}

int main(void){
	int n=10,i,j,m=0;
	for(i=0; i<n; ++i){
		for(j=0; j<n; ++j){
			scanf("%d",&h[i][j]);
		}
		if(i%2){
			for(j=n-1; j>=0; --j){
				b[i][j] = m;
				x[m] = i; y[m] = j;
				++m;
			}
		}else{
			for(j=0; j<n; ++j){
				b[i][j] = m;
				x[m] = i; y[m] = j;
				++m;
			}
		}
	}
	cout << fixed << setprecision(50) << solve(99,1) << endl;
	return 0;
}
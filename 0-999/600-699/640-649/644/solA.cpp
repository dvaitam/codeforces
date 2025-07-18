#include <cstdio>
#include <iostream>
#include <vector>
#include <queue>
#include <map>
#include <cmath>
#include <cstring>
#include <algorithm>
using namespace std;

#define pb push_back
#define ri(x) scanf("%d",&x)
#define rii(x,y) ri(x),ri(y)
#define ms(obj,val) memset(obj,val,sizeof(obj))
#define ms2(obj,val,sz) memset(obj,val,sizeof(obj[0])*sz)
#define FOR(i,f,t) for(int i=f; i<(int)t; i++)
#define FORR(i,f,t) for(int i=f; i>(int)t; i--)

typedef long long ll;
typedef vector<int> vi;

int N,A,B,M[100][100];

int main() {
	ri(N),rii(A,B);
	if(N>A*B) printf("-1\n");
	else {
		int cnt=1;
		FOR(i,0,A) {
			FOR(j,0,B) {
				M[i][j]=cnt;
				if(B%2==0 && i%2!=0 && j%2==0) M[i][j]++;
				if(B%2==0 && i%2!=0 && j%2!=0) M[i][j]--;
				if(M[i][j]>N) M[i][j]=0;
				cnt++;
				if(cnt>N+1) break;
			}
			if(cnt>N+1) break;
		}
		FOR(i,0,A) FOR(j,0,B) printf("%d%c",M[i][j]," \n"[j==B-1]);
	}
}
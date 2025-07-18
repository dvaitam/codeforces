#include <stdio.h>
#include <string>
#include <sstream>
#include <vector>
#include <algorithm>
#include <utility>
#include <math.h>

using namespace std;

#define lint long long

#define ss stringstream
#define sz size()
#define pb push_back
#define mp make_pair

#define FOR(i,n) SFOR(i,0,n)
#define SFOR(i,m,n) for(i=m;i<n;i++)
#define FORD(i,n) for(i=n-1;i>=0;i--)

int A[200000];
int n;

bool Try(int q) {
	if (n < 3*q) return false;
	int i,j,k;
	k = 0;
	if (A[q-1] == A[n-q]) return false;
	if (A[q-1] == A[q]) {
		FOR(j,q) if (A[j] == A[q]) break;
		SFOR(i,q,n) if (A[i] != A[q]) break;
		i -= q;
		if (i > j) k += i-j;
	}
	if (n - k < 3*q) return false;
	if (A[n-q-1] == A[n-q]) {
		SFOR(j,q,n-q) if (A[j] == A[n-q]) break;
		j -= k+q;
		SFOR(i,n-q,n) if (A[i] != A[n-q]) break;
		i -= (n-q);
		if (i > j) return false;
	}
	return true;
}

void Out(int q) {
	int i,j,k;
	i = 0;
	j = q;
	k = n-q;
	printf("%d\n",q);
	while (q != 0) {
		while (A[i] == A[j]) j++;
		printf("%d %d %d\n",A[k],A[j],A[i]);
		i++;j++;k++;
		q--;
	}
}

int main() {
	int i;
	scanf("%d",&n);
	FOR(i,n) scanf("%d",&A[i]);
	sort(A,A+n);
	int l,r,m;
	l = 0;
	r = n;
	while (r > l+1) {
		m = (r+l)/2;
		if (Try(m)) l = m; else r = m;
	}
	Out(l);
}
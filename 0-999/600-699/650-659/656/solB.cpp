#include <bits/stdc++.h>
using namespace std;
int m[20], r[20];
int gcd(int a, int b){
	int e=1;
	while(b){
		e = a%b;
		a=b, b=e;
	}
	return a;
}
int main(){
	int n;
	scanf("%d", &n);
	for(int i=0; i<n; i++) scanf("%d", &m[i]);
	for(int i=0; i<n; i++) scanf("%d", &r[i]);
	int ans=1;
	for(int i=0; i<n; i++){
		int p = gcd(ans, m[i]);
		ans = ans/p*m[i];
	}
	int num=0;
	for(int i=0; i<ans; i++){
		for(int j=0; j<n; j++){
			if(i%m[j]==r[j]) {
				num++;break;
			}
		}
	}
	double fff = num*1.0/ans;
	printf("%.10lf\n", fff);
	return 0;
}
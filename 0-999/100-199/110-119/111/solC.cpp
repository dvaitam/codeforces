#include <stdio.h>
#include <iostream>
using namespace std;
int n,m;
int main( ){
	cin>>n>>m;
	if (n>m) swap(n,m);
	if (n==1) cout<<m*n-(m+2)/(4-n)<<endl;
	if (n==2) cout<<m*n-(m+2)/(4-n)<<endl;
	if (n==3) cout<<m*n-m/4*3-(m%=4)-(m==0)<<endl;
	if (n==4) cout<<m*n-m-(m==5 || m==6 || m==9)<<endl;
	if (n==5) cout<<m*n-m/5*6+(m==7)-(m%=5)-1-(m>1)<<endl;
	if (n==6) cout<<m*n-10<<endl;
	return 0;
}
#include<bits/stdc++.h>

using namespace std;

long long a;

int main(){

	scanf("%lld",&a);

	long long now = 100000000000000000ll%a,P=0;

	for (int i=1;i<=9;i++) P = (P + (now * i %a)) % a;

	now=0;

	for (int i=1;i<=18;i++) now = (now + P)%a; 

	now = (now + 1)  % a;

	long long l=1,r=1000000000000000000ll;

	l+=a-now;

	r+=a-now;

	printf("%lld %lld",l,r);

	return 0;

}
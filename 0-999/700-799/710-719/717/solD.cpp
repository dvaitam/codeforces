#include<stdio.h>

#include<bits/stdc++.h>

using namespace std;

#define ll long long

#define f(i,x,n) for (int i = x;i < n;++i)



typedef vector<double> mt;



mt operator *(mt &a, mt &b) {

	mt t(128, 0.0);

	f(i, 0, 128)f(j, 0, 128)t[i ^ j] += a[i] * b[j];

	return t;

}



mt pw(mt &a, int p) {

	if (p == 1)return a;

	mt t = pw(a, p >> 1);

	t = t * t;

	if (p & 1)t = t * a;

	return t;

}



int main() {

	mt z(128, 0.0);

	int n, m;

	scanf("%d%d", &n, &m);

	f(i, 0, m + 1) {

		double t;

		scanf("%lf", &t);

		z[i] = t;

	}

	printf("%.7lf", 1.0 - pw(z, n)[0]);

}
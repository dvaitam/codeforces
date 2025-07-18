#include<iostream>

#include<stdio.h>

#include<stdlib.h>

#include<string.h>

#include<string>

#include<vector>

#include<set>

#include<algorithm>

#include<stack>

#include<queue>

#include<map>

#include<limits>

#include<list>

#include<math.h>

using namespace std;

#define ll long long

#define f(i,x,n) for (int i = x;i < n;++i)



int n, k, x[10000];



bool cn(double a) {

	double b = 0.0;

	f(i, 0, n) {

		if (x[i] > a)b += (x[i] - a) / 100 * (100 - k);

		else b -= a - x[i];

	}

	if (b > 0)return true;

	else return false;

}



int main() {

	//freopen("in.txt", "r", stdin);

	scanf("%d%d", &n, &k);

	f(i, 0, n)scanf("%d", x + i);

	double l = 0, r = 1000, m;

	while (r > l) {

		m = (l + r) / 2.0;

		if (cn(m))l = m + 1e-8;

		else r = m - 1e-8;

	}

	printf("%.8lf", (r + l) / 2);

}
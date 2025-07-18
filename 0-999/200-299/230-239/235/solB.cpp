#include <cstdio>
#include <cmath>
#include <cstring>
#include <algorithm>
#include <set>
#include <vector>
#include <map>
#include <string>
#include <functional>
#include <numeric>
#include <iostream>
#include <sstream>
#include <ctime>
#include <cassert>

using namespace std;

typedef long long i64;
typedef unsigned char u8;

int main() {
#ifdef pperm
	freopen("input.txt", "r", stdin);
	//freopen("output.txt", "w", stdout);
#endif
	int n;
	scanf("%d", &n);
	double m = 0;
	double r = 0;
	for (int i = 0; i < n; i++) {
		double p;
		scanf("%lf", &p);
		r += p * (m * 2 + 1);
		m = (m + 1) * p;
	}
	printf("%.10lf\n", r);
#ifdef pperm
	fprintf(stderr, "%.3lf\n", clock() / double (CLOCKS_PER_SEC));
#endif
	return 0;	
}
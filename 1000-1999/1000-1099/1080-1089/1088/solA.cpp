#include <iostream>
#include <stdio.h>
#include <cmath>
#include <vector>
#include <string>
#include <functional>
#include <set>
#include <cstdlib>
#include <map>
#include <cctype>
#include <algorithm>
#define pii pair<int, int>
#define ll long long
#define pis pair<int, string>
#define mp make_pair

const ll MOD = 998244353;
const int INF = 1 << 30;
using namespace std;

int main() {
	//freopen("input.txt", "r", stdin);
	int n;
	scanf("%d", &n);
	if (n == 1) {
		printf("-1");
	}
	else {
		printf("%d %d", n, n);
	}
	return 0;
}
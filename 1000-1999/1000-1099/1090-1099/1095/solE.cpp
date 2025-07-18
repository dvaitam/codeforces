#include <iostream>
#include <set>
#include <map>
#include <list>
#include <queue>
#include <stack>
#include <cmath>
#include <string>
#include <vector>
#include <cstring>
#include <cstdio>
#include <utility>
#include <algorithm>
#include <functional>
using namespace std;
typedef long long ll;
typedef pair<int, int> P;
typedef vector<int> vi;
char s[1000100];
int a[1000100];
int main() {
	int l;
	scanf("%d%s", &l, s);
	if (l & 1) {
		return !printf("0\n");
	}
	int num = 0;
	for (int i = 0; i < l; i++) {
		if (s[i] == '(') {
			num++;
			a[i] = num;
		}
		else {
			num--;
			a[i] = num;
		}
	}
	if (a[l - 1] == 0)printf("0\n");
	else if (a[l - 1] == -2) {
		int right = l - 1;
		int wa = 0;
		bool ok = 1;
		for (int i = 0; i < l; i++) {
			if (a[i] < -2)ok = 0;
		}
		if (!ok)return !printf("0\n");
		for (int i = 0; i < l; i++) {
			if (a[i] < 0) {
				right = i;
				break;
			}
		}
		for (int i = right; i >= 0; i--) {
			if (s[i] == ')')wa++;
		}
		printf("%d\n", wa);
	}
	else if (a[l - 1] == 2) {
		int left = 0;
		int wa = 0;
		bool ok = 1;
		for (int i = 0; i < l; i++) {
			if (a[i] < 0)ok = 0;
		}
		if (!ok)return !printf("0\n");
		for (int i = l - 1; i >= 0; i--) {
			if (a[i] < 2) {
				left = i + 1;
				break;
			}
		}
		for (int i = left; i < l; i++) {
			if (s[i] == '(')wa++;
		}
		printf("%d\n", wa);
	}
	else printf("0\n");
}
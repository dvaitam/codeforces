#include <bits/stdc++.h>

#include <cmath>

#include <cstdio>

#include <vector>

#include <queue>

#include <iostream>

#include <algorithm>

#include <cstdlib>

#include <stack>

#define REP(i,x,y) for(int i = x; i<=y; i++)

using namespace std;



int main() {

	int n, k;

	cin >> n >> k;

	int a[n+1];

	int long long sum = 0;

	REP(i,1,n) {

		cin >> a[i];

		sum += a[i];

	}

	int minDiff = sum%n? 1 : 0;

	int currDiff;

	int minH = 1, maxH = 1;

	REP(i,1,n) {

		if(a[i]<a[minH]) minH = i;

		if(a[maxH]<a[i]) maxH = i;

	}

	currDiff = a[maxH]-a[minH];

	if(currDiff==minDiff) {

		cout << minDiff << " 0\n";

		return 0;

	}

	int t;

	vector<int> minP, maxP;

	maxP.push_back(0);

	minP.push_back(0);

	for(t = 1; t<=k; t++) {

		a[maxH]--; maxP.push_back(maxH);

		a[minH]++; minP.push_back(minH);

		REP(i,1,n) {

			if(a[i]<a[minH]) minH = i;

			if(a[maxH]<a[i]) maxH = i;

		}

		currDiff = a[maxH]-a[minH];

		if(currDiff==minDiff) break;

	}

	if(t>k) {

		cout << currDiff << " " << k << endl;

		REP(i,1,k) {

			cout << maxP[i] << " " << minP[i] << endl;

		}

		return 0;

	}

	cout << currDiff << " " << t << endl;

	REP(i,1,t) {

		cout << maxP[i] << " " << minP[i] << endl;

	}

	return 0;

}
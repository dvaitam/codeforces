#include <iostream>

#include <cstdio>

#include <algorithm>

#include <cmath>

#include <string>

#include <vector>

#include <stack>

#include <queue>

#include <set>

#include <cstring>

#include <map>

#include <cstdlib>

#include <ctime>

#include <cassert>

#include <bitset>

#define f first

#define s second

#define ll long long

#define ull unsigned long long

#define mp make_pair

#define pb push_back

#define vi vector <int>

#define ld long double

#define pii pair<int, int>

#define y1 sdasd

using namespace std;    

const int N = int(3e5), mod = int(1e9)  + 7; 



int x1,y1,x2,y2,v,t,vx,vy,wx,wy;



double sqr(double x) {

	return x * x;

}



bool check(double m){

	if(m <= t){

		return (sqrt(sqr(x1 + vx * m - x2) + sqr(y1 + vy * m - y2)) <= v * m);	

	}

	else{

		return (sqrt(sqr(x1 + vx * t + (m - t) * wx - x2) + sqr(y1 + vy * t + (m - t) * wy - y2)) <= v * m);

	}

}



int main () {

	cin >> x1 >> y1 >> x2 >> y2;

	cin >> v >> t;

	cin >> vx >> vy >> wx >> wy;



	double l =0, r = 1e8,mid;

	for(int it = 0; it < 200; it++){

		mid = (r + l) * 0.5;

		if(check(mid)) r = mid;

		else l = mid;

	}

	printf("%.12lf",r);



return 0;

}
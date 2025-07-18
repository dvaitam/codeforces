// #include <bits/stdc++.h>
#include <iostream>
#include <limits>
#include <iomanip>
#include <cstdio>
#include <algorithm>
#include <cstring>
#include <cassert>
#include <vector>
#include <list>
#include <stack>
#include <queue>
#include <set>
#include <map>
#include <cmath>
#include <sstream>

using namespace std;

#define MAXINT numeric_limits<int>::max()
#define MININT numeric_limits<int>::min()
#define PI 3.141592653589793
#define rep(i, a, n) for(int i = (a); i < (n); i++)
#define per(i, a, n) for(int i = (n)-1; i >= (a); i--)
#define pb push_back
#define mp make_pair
#define fi first
#define sn second
#define vi vector<int>
#define si set<int>
#define li list<int>

	

typedef long long ll;
typedef pair<int, int> pii;
typedef pair<long long, long long> pll;
typedef pair<double, double> pdd;

template <typename T>
void print_array(T arr[], int size_arr){
	rep(i, 0, size_arr){
		cout << arr[i];
		if(i == size_arr - 1) cout << endl;
		else cout << " ";
	}
}

template<typename T>
void print_vector(vector<T> v){
	rep(i, 0, v.size()){
		cout << v[i];
		if(i == v.size() - 1) cout << endl;
		else cout << " ";
	}
}

template <typename T>
void memset_array(T arr[], T value, int size_arr){
	rep(i, 0, size_arr){
		arr[i] = value;
	}
}

bool compare_lexical_string(string a, string b){
	rep(i, 0, min(a.length(), b.length())){
		if(a[i] != b[i]) return a[i] < b[i];
	}
	return a.length() < b.length();
}
struct lex_string
{
	bool operator()(string a, string b){
		return compare_lexical_string(a, b);
	}
};



/////////////////////////////////////////////////////////////////////////////
// NEVER USE NAME : DISTANCE FOR FUNCTION
// NEVER USE NAME : FIND FOR FUNCTION
// AVOID DOUBLE IF POSSIBLE
////////////////////////////////////////////////////////////////////////////
// compare_lexical_string("huy", "huy") will return false
long double compute(long double x1, long double y1, long double r1, long double x2, long double y2, long double r2){
	if(r1 < r2) return compute(x2, y2, r2, x1, y1, r1);
	//
	long double dis_square = (x1 - x2)*(x1-x2) + (y1-y2)*(y1-y2);
	if(dis_square <= (r1 - r2)*(r1-r2)){
		return r2*r2*PI;
	}
	else if(dis_square >= (r1 + r2)*(r1+r2)){
		return 0;
	}
	else{
		long double cos_2 = (r2*r2 + dis_square - r1*r1)/(2*r2*sqrt(dis_square));
		// cout << cos_2 << endl;
		long double cos_1 = (r1*r1 + dis_square - r2*r2)/(2*r1*sqrt(dis_square));
		// cout << cos_1 << endl;
		long double p1 = r1*r1*(acos(cos_1) - cos_1*sin(acos(cos_1)));
		long double p2 = r2*r2*(acos(cos_2)- cos_2*sin(acos(cos_2)));
		return p1 + p2;
	}
}

int main(){
	ios_base::sync_with_stdio(0); cin.tie(0);
	long double x1, y1, r1, x2, y2, r2;
	cin >> x1 >> y1 >> r1 >> x2 >> y2 >> r2;
	// cout << acos(-0.5) << endl;
	cout << fixed << setprecision(20) << compute(x1, y1, r1, x2, y2, r2) << endl;
	
	return 0;
}
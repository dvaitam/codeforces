#include <iostream>
#include <algorithm>
#include <vector>
#include <string>
#include <map>
#include <set>
#include <queue>
#include <stack>
#include <cmath>
#include <cassert>
#include <cstring>
#include <ext/numeric>
using namespace std ;
using namespace __gnu_cxx ;
typedef long long LL ;
typedef pair<int,int> PII ;
typedef vector<int> VI ;
const int INF = 1000*1000*1000+1 ;
const LL INFLL = (LL)INF * (LL)INF ;
#define REP(i,n) for(i=0;i<(n);++i)
#define ALL(c) c.begin(),c.end()
#define VAR(v,i) __typeof(i) v=(i)
#define FOREACH(i,c) for(VAR(i,(c).begin());i!=(c).end();++i)
#define CLEAR(t) memset(t,0,sizeof(t))
#define PB push_back
#define MP make_pair
#define FI first
#define SE second

const int MAXN = 310 ;

struct line {
	int a, b ;
	double x0, x1 ;
	line(int aa, int bb, double xx0, double xx1) : a(aa), b(bb), x0(xx0), x1(xx1) {}
} ;
bool operator<(const line &X, const line &Y) {
	return X.a < Y.a ;
}

set<line> S[MAXN] ;

double eval(int a, int b, double x) {
	return a*x+b ;
}

double P(int a, int b, double x0, double x1) {
	return (eval(a,b,x0)+eval(a,b,x1))*(x1-x0)/2.0 ;
}

int main()
{
	ios_base::sync_with_stdio(0) ;
	int n, k, i, j, y0, y1 ;
	cin >> n >> k ;
	REP(i,n) {
		cin >> y0 ;
		double area=0.0 ;
		REP(j,k) {
			cin >> y1 ;
			int a = y1-y0 ;
			int b = y0 ;
			if(S[j].empty()) {
				S[j].insert(line(a, b, 0,1)) ;
				area += P(a,b, 0,1) ;
			}
			else {
				double left = INF, right=-INF ;
				
				set<line>::iterator q = S[j].lower_bound(line(a,b,0,1)) ;
				while(q != S[j].end()) {
					if(eval(a,b, q->x0) >= eval(q->a, q->b, q->x0) &&
						eval(a,b, q->x1) >= eval(q->a, q->b, q->x1)) {
							area += P(a,b, q->x0, q->x1) - P(q->a, q->b, q->x0, q->x1) ;
							left = min(left, q->x0) ;
							right = max(right, q->x1) ;
							S[j].erase(q++) ;
						}
					else {
						if(eval(a,b, q->x0) > eval(q->a, q->b, q->x0)) {
							double x_sr = (double)(b - q->b)/(q->a - a) ;
							area += P(a,b, q->x0, x_sr) - P(q->a, q->b, q->x0, x_sr) ;
							left = min(left, q->x0) ;
							right = max(right, x_sr) ;
							
							int aa=q->a, bb=q->b ;
							double xx1 = q->x1 ;
							S[j].erase(q) ;
							S[j].insert(line(aa,bb, x_sr, xx1)) ;
						}
						break ;
					}
				}
				
				q = S[j].lower_bound(line(a,b,0,1)) ;
				if(q != S[j].begin()) {
					q-- ;
					while(1) {
						if(eval(a,b, q->x0) >= eval(q->a, q->b, q->x0) &&
							eval(a,b, q->x1) >= eval(q->a, q->b, q->x1)) {
								area += P(a,b, q->x0, q->x1) - P(q->a, q->b, q->x0, q->x1) ;
								left = min(left, q->x0) ;
								right = max(right, q->x1) ;
								if(q == S[j].begin()) {
									S[j].erase(q) ;
									break ;
								}
								else S[j].erase(q--) ;
						}
						else {
							if(eval(a,b, q->x1) > eval(q->a, q->b, q->x1)) {
								double x_sr = (double)(b - q->b)/(q->a - a) ;
								area += P(a,b, x_sr, q->x1) - P(q->a, q->b, x_sr, q->x1) ;
								left = min(left, x_sr) ;
								right = max(right, q->x1) ;
								
								int aa=q->a, bb=q->b ;
								double xx0 = q->x0 ;
								S[j].erase(q) ;
								S[j].insert(line(aa,bb, xx0, x_sr)) ;
							}
							break ;
						}
					}
				}
				
				if(left<=right) S[j].insert(line(a,b, left, right)) ;
			}
			y0 = y1 ;
		}
		cout.setf(ios::fixed) ;
		cout.precision(10) ;
		cout << area << endl ;
	}
}
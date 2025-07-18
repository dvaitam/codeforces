#include<iostream>
#include<vector>
#include<iomanip>
#include<cmath>
using namespace std;

double PI = atan2(0, -1);
double a = 72.0 / 180.0 * PI;

struct Vec {
	double y, x;
	Vec(double ny = 0, double nx = 0) {
		y = ny;
		x = nx;
	}
	
	Vec& operator+=(Vec r) {
		y += r.y, x += r.x;
		return *this;
	}
	Vec& operator-=(Vec r) {
		y -= r.y, x -= r.x;
		return *this;
	}
	Vec operator+(Vec r) { return Vec(*this)+=r;}
	Vec operator-(Vec r) { return Vec(*this)-=r;}
};

Vec rotate(Vec l, double r) {
	return Vec(l.y*cos(r) + l.x*sin(r), l.x*cos(r) - l.y*sin(r) );
}


//-----------------------------------------------------
//Main function


int main() {
	ios_base::sync_with_stdio(false);
	cin.tie(0);
	
	Vec base = Vec(0, 10);
	vector<Vec> pts(1, Vec(0, 0) );
	
	int n;
	cin >> n;
	
	for(int i=0;i<n;i++) {
		int c = pts.size()-1;
		pts.push_back(pts[c] - rotate(base, -2*a) );
		pts.push_back(pts[c] + rotate(base, -a) );
		pts.push_back(pts.back() + base);
		pts.push_back(pts.back() + rotate(base, a) );
	}
	
	cout << pts.size() << '\n';
	for(auto c : pts)
		cout << fixed << setprecision(9) << c.x << ' ' << c.y << '\n';
	
	for(int i=0;i<n;i++) {
		cout << i*4+1 << ' ' << i*4+3 << ' ' << i*4+4 << ' ' << i*4+5 << ' ' << i*4+2 << '\n';
	}
	
	cout << "1";
	for(int i=0;i<n;i++) {
		cout << ' ' << i*4+4 << ' ' << i*4+2 << ' ' << i*4+3 << ' ' << i*4+5;
	}
	for(int i=n-1;i>=0;i--)
		cout << ' ' << i*4+1;
	cout << '\n';
	
	return 0;
}
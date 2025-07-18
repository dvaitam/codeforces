#include <iostream>
#include <memory>
#include <cstdio>
#include <iomanip>
#include <cstring>
#include <algorithm>
#define _USE_MATH_DEFINES 
#include <cmath>
#include <string>
using namespace std;

int main(){
	ios_base::sync_with_stdio(0);
	int n,x,y;
	double d;
	cin >> n;
	double *a=new double[n+1];
	for (int i=0;i<n;++i) {
		cin >> x >> y;
		if (x==0&&y>=0) a[i]=90;
		else if (x==0&&y<=0) a[i]=270; 
		else {
			d=atan(abs(((double)y)/x))*180/M_PI;
			if (y>=0&&x>=0) a[i]=d;
			else if (y>=0&&x<=0) a[i]=180-d;
			else if (y<=0&&x<=0) a[i]=180+d;
			else a[i]=360-d;
		}
	}
	sort(a,a+n);
	d=-1;
	for (int i=1;i<n;++i) if (d==-1||a[i]-a[i-1]>d) d=a[i]-a[i-1];
	if (d==-1||360-a[n-1]+a[0]>d) d=360-a[n-1]+a[0];
	cout << setprecision(8) << fixed << 360-d << endl;
    return 0;
}
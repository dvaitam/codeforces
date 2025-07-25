#include <cstdio>
#include <cstdlib>
#include <cstring>
#include <algorithm>
#include <vector>
#include <cmath>
using namespace std;


#define fr(i,a,b) for(int i=a;i<b;++i)



struct P{
	double x,y,z;
	P(){}
	P(double X, double Y, double Z){x = X, y = Y, z = Z; }
	P operator +(P a){return P(x+a.x, y+a.y, z+a.z); }
	P operator -(P a){return P(x-a.x, y-a.y, z-a.z); }
	P operator *(double a){return P(x*a, y*a, z*a); }
	double operator !(){
		return sqrt(x*x + y*y + z*z);
	}
}v[10100];

P onde = P(-1e10,-1e10,-1e10);
P h;
double vs, vh;
int n;
bool can(P a, double x){
	return (!(h-a)/vh) < x+1e-11;
}





int main(){
	scanf("%d",&n);
	fr(i,0,n+1) scanf("%lf %lf %lf",&v[i].x,&v[i].y,&v[i].z);
	double tini = 0.0;
	scanf("%lf %lf",&vh,&vs);
	scanf("%lf %lf %lf",&h.x,&h.y,&h.z);
	int end;
	fr(i,0,n){
		if(can(v[i+1], tini + !(v[i]-v[i+1])/vs)){
			double ini = 0.0, fim = 1.0;
			fr(j,0,100){
				double meio = (ini+fim)/2.0;
				P m = v[i] + (v[i+1]-v[i])*meio;
				if(can(m, tini + !(v[i]-m)/vs)) fim = meio;
				else ini = meio;
			}
			onde = v[i] + (v[i+1]-v[i])*ini;
			end = i;
			break;
		}
		tini += !(v[i]-v[i+1])/vs;
	}
	if(onde.x > -1e8){
		printf("YES\n");
		printf("%.10lf\n",max(!(onde-h)/vh, tini + !(v[end]-onde)/vs));
		printf("%.10lf %.10lf %.10lf\n",onde.x,onde.y,onde.z);
	}
	else printf("NO\n");
	return 0;
}
#include <algorithm>

#include <cmath>

#include <cstring>

#include <cassert>

#include <cstdio>

#include <iostream>

#define fst first

#define snd second

#define pb push_back

#define mp make_pair

#define fore(i,a,n) for(int i = a, to = n; i < to;i++)



using namespace std;



typedef long long ll;

typedef pair<int,int> ii;



//hola

#define eps 1e-7



struct point {

	double x,y;

	point(double a, double b) : x(a), y(b) {}

	point() {}

	point operator*(double t){return point(t*x,t*y);}

	point operator+(point p){return point(x+p.x,y+p.y);}

	point operator-(point p){return point(x-p.x,y-p.y);}

	double operator%(point p){return x*p.y-y*p.x;}

	double operator*(point p) {return x*p.x+y*p.y;}

	bool operator==(point p){return abs(x- p.x)<eps && abs(y-p.y)<eps;}

	bool operator<(point p)const {

		if(x == p.x)

			return y < p.y-eps;

		return x < p.x-eps;

	}

};



point rota90h(point p) {

	point q = point(-p.y,p.x);

	return q;

}



point rota90ah(point p) {

	point q = point(p.y,-p.x);

	return q;

}



bool esrec(point a, point b, point c, point d) {

	bool B = ((b-a)*(b-c) == 0) && ((b-a)*(a-d)==0) && (d-c)*(c-b)==0;

	return B;

}



bool escua(point a, point b, point c, point d) {

	bool B = ((b-a)*(b-c) == 0) && ((b-a)*(a-d)==0) && (c-a)*(d-b)==0;

	return B;

}



bool escuadrado(point a,point b, point c, point d) {

	bool B = escua(a,b,c,d) || escua(a,b,d,c) || escua(a,c,b,d) || escua(a,c,d,b) || escua(a,d,b,c) || escua(a,d,c,b);

	return B;

}



bool esrectangulo(point a,point b, point c, point d) {

	bool B = esrec(a,b,c,d) || esrec(a,b,d,c) || esrec(a,c,b,d) || esrec(a,c,d,b) || esrec(a,d,b,c) || esrec(a,d,c,b);

	return B;

}





point P[10];

bool Q[10];

point R[10];

int C[10];



int main(){

	double a,b;

	for (int i = 0; i < 8; i++) {

		scanf("%lf %lf",&a,&b);

		P[i] = point(a,b);

	}

	memset(Q,true,sizeof(Q));

	

	for(int i = 0; i<5; i++) {

		for(int j = i+1; j<6; j++) {

			for (int k = j+1; k<7; k++) {

				for (int l = k+1; l<8; l++) {

					if (escuadrado(P[i],P[j],P[k],P[l])) {

						Q[i]=false;

						Q[j]=false;

						Q[k]=false;

						Q[l]=false;

						int n = 0;

						for (int m=0; m<8; m++) {

							if (Q[m]) {R[n]=P[m]; C[n]=m; n++;}

						}

						if (esrectangulo(R[0],R[1],R[2],R[3])) {

							printf("YES\n");

							printf("%d %d %d %d\n",i+1,j+1,k+1,l+1);

							printf("%d %d %d %d\n",C[0]+1,C[1]+1,C[2]+1,C[3]+1);

							return 0;

						}

						memset(Q,true,sizeof(Q));

					}

				}

			}

		}

	}

	printf("NO\n");

	return 0;

}
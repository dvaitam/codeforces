#include <stdio.h>
#include <math.h>
#include <algorithm>

using namespace std;

#define EPS 1E-8

inline int SG(double x){
    if(x>EPS) return 1;
    if(x<-EPS) return -1;
    return 0;
}

class PT{ public:
    double x,y;
    PT(){}
    PT(double _x,double _y){
	x=_x; y=_y;
    }
    PT operator+(const PT& p) const{ return PT(x+p.x,y+p.y); }
    PT operator-(const PT& p) const{ return PT(x-p.x,y-p.y); }
    double operator^(const PT& p) const{ return x*p.y-y*p.x; }
    double operator*(const PT& p) const{ return x*p.x+y*p.y; }
    double disTo(PT& p) const{ return sqrt((p.x-x)*(p.x-x)+(p.y-y)*(p.y-y)); }
};

inline double tri(PT& p1,PT& p2,PT& p3){
    return (p2-p1)^(p3-p1);
}

class PY{ public:
    int n;
    PT pt[5];
    PT& operator[](const int x){
	return pt[x];
    }
    void input(){
	int i;
	n=4;
	for(i=0;i<n;i++){
	    scanf("%lf %lf",&pt[i].x,&pt[i].y);
	}
    }
    double getArea(){
	int i;
	double s=pt[n-1]^pt[0];
	for(i=0;i<n-1;i++){
	    s+=pt[i]^pt[i+1];
	}
	return s/2;
    }
};

PY py[500];
pair<double,int> c[5000];

inline double segP(PT &p,PT &p1,PT &p2){
    if(SG(p1.x-p2.x)==0) return (p.y-p1.y)/(p2.y-p1.y);
    return (p.x-p1.x)/(p2.x-p1.x);
}

double polyUnion(int n){
    int i,j,ii,jj,ta,tb,r,d;
    double z,w,s,sum,tc,td;
    for(i=0;i<n;i++){
	py[i][py[i].n]=py[i][0];
    }
    sum=0;
    for(i=0;i<n;i++){
	for(ii=0;ii<py[i].n;ii++){
	    r=0;
	    c[r++]=make_pair(0.0,0);
	    c[r++]=make_pair(1.0,0);
	    for(j=0;j<n;j++){
		if(i==j) continue;
		for(jj=0;jj<py[j].n;jj++){
		    ta=SG(tri(py[i][ii],py[i][ii+1],py[j][jj]));
		    tb=SG(tri(py[i][ii],py[i][ii+1],py[j][jj+1]));
		    if(ta==0 && tb==0){
			if((py[j][jj+1]-py[j][jj])*(py[i][ii+1]-py[i][ii])>0 && j<i){
			    c[r++]=make_pair(segP(py[j][jj],py[i][ii],py[i][ii+1]),1);
			    c[r++]=make_pair(segP(py[j][jj+1],py[i][ii],py[i][ii+1]),-1);
			}
		    }else if(ta>=0 && tb<0){
			tc=tri(py[j][jj],py[j][jj+1],py[i][ii]);
			td=tri(py[j][jj],py[j][jj+1],py[i][ii+1]);
			c[r++]=make_pair(tc/(tc-td),1);
		    }else if(ta<0 && tb>=0){
			tc=tri(py[j][jj],py[j][jj+1],py[i][ii]);
			td=tri(py[j][jj],py[j][jj+1],py[i][ii+1]);
			c[r++]=make_pair(tc/(tc-td),-1);
		    }
		}
	    }
	    sort(c,c+r);
	    z=min(max(c[0].first,0.0),1.0);
	    d=c[0].second;
	    s=0;
	    for(j=1;j<r;j++){
		w=min(max(c[j].first,0.0),1.0);
		if(!d) s+=w-z;
		d+=c[j].second;
		z=w;
	    }
	    //printf("%.5f\n",s);
	    sum+=(py[i][ii]^py[i][ii+1])*s;
	}
    }
    return sum/2;
}

int main(){
    int n,i,j,k;
    double sum,ds;
    scanf("%d",&n);
    sum=0;
    for(i=0;i<n;i++){
	py[i].input();
	ds=py[i].getArea();
	if(ds<0){
	    for(j=0,k=py[i].n-1;j<k;j++,k--){
		swap(py[i][j],py[i][k]);
	    }
	    ds=-ds;
	}
	sum+=ds;
    }
    printf("%.9f\n",sum/polyUnion(n));
    return 0;
}
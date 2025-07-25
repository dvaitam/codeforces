#include <cstdio>
#include <algorithm>
#include <cmath>
#include <cstring>
#define sz(a) int((a).size())
#define pb push_back
#define rep(i,j,k) for (int i=j; i<=k; i++)
#define online1
#define sqr(x) (x)*(x)
using namespace std;
typedef long long LL;

const int maxn=1100;

const double eps=1e-8;

int n,m;
bool fuck;

int sign(double x){
    if (x<-eps) return -1; return x>eps;
}

struct point_3d{
    double x,y,z;
    point_3d(double _x,double _y,double _z):x(_x),y(_y),z(_z){
    }
    point_3d(){
    }
    point_3d operator+(point_3d &a){
        return point_3d(a.x+x, a.y+y, a.z+z);
    }
    point_3d operator-(point_3d &a){
        return point_3d(x-a.x, y-a.y, z-a.z);
    }
    point_3d operator/(double f){
        return point_3d(x/f, y/f, z/f);
    }
    double len(){
        return x*x+y*y+z*z;        
    }
};

double dis(point_3d &a, point_3d &b){
    return sqr(a.x-b.x)+sqr(a.y-b.y)+sqr(a.z-b.z);   
}

point_3d project(double A,double B,double C,double x1,double y1,double z1){
    double lambda=(A*x1+B*y1+C*z1)/(A*A+B*B+C*C);
    return point_3d(x1-lambda*A,y1-lambda*B,z1-lambda*C);
}

void calc2(double a,double b,double c,double d,double e,double f,double &x,double &y){
    if (sign(a)==0) swap(a,d),swap(b,e),swap(c,f);
    e+=(-b/a*d); f+=(-c/a*d);
    y=-f/e;
    x=(-b*y-c)/a;   
}

point_3d cir_cen(point_3d &a, point_3d& b, point_3d& c, double A3, double B3, double C3){
    double x,y,z,A1,B1,C1,D1,A2,B2,C2,D2,D3=0;
    
    A1=-2*a.x+2*b.x;
    B1=-2*a.y+2*b.y;
    C1=-2*a.z+2*b.z;
    D1=-(b.x*b.x+b.y*b.y+b.z*b.z-a.x*a.x-a.y*a.y-a.z*a.z);
    
    A2=-2*c.x+2*b.x;
    B2=-2*c.y+2*b.y;
    C2=-2*c.z+2*b.z;
    D2=-(b.x*b.x+b.y*b.y+b.z*b.z-c.x*c.x-c.y*c.y-c.z*c.z);
    
    if (sign(A2)!=0) swap(A2,A1),swap(B2,B1),swap(C2,C1),swap(D2,D1);
    if (sign(A3)!=0) swap(A3,A1),swap(B3,B1),swap(C3,C1),swap(D3,D1);
    
    if (sign(A1)==0) {fuck=true; return a;}
    
    B2+=(-B1/A1*A2); C2+=(-C1/A1*A2); D2+=(-D1/A1*A2);
    B3+=(-B1/A1*A3); C3+=(-C1/A1*A3); D3+=(-D1/A1*A3);
    
    calc2(B2,C2,D2,B3,C3,D3,y,z);
    
    return point_3d((-B1*y-C1*z-D1)/A1,y,z);   
}

point_3d p[maxn],a[maxn],o;
double r;

int main(){
    scanf("%d%d",&n,&m);
    rep(i,1,n) scanf("%lf%lf%lf",&p[i].x,&p[i].y,&p[i].z);
    
    rep(mm,1,m){
        double A,B,C;
        scanf("%lf%lf%lf",&A,&B,&C);    
        
        rep(i,1,n) a[i]=project(A,B,C,p[i].x,p[i].y,p[i].z);
            
        random_shuffle(a+1,a+n+1);
        o=point_3d(0.00,0.00,0.00); r=0;
        rep(i,1,n)
            if ( sign( dis(o,a[i])-r )==1 ){
                o=a[i]; r=0;
                rep(j,1,i-1)
                    if ( sign( dis(o,a[j])-r )==1 ){
                        o=(a[i]+a[j])/2; r=(o-a[i]).len(); 
                        rep(k,1,j-1)
                            if ( sign( dis(o,a[k])-r )==1 ){
                                o=cir_cen(a[i],a[j],a[k],A,B,C);  r=(o-a[k]).len();
                            }
                    }
            }  
        printf("%.9lf\n",sqrt(r));
    }   

    return 0;
}
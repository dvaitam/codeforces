#include <bits/stdc++.h>
using namespace std;
#define fr(x, y, z) for(int x=y;x<z;x++)
typedef long long ll;
typedef long double ld;
typedef pair<ll,ll> pll;
typedef pair<pll,ll> ppl;
const int ms=1e4+10;

const double eps = 1e-8;
const double pi = acos(-1.0);
 
int cmp(double a, double b = 0) {
    if (fabs(a-b) < eps) return 0;
    else if (a < b) return -1;
    return 1;
}
 
struct PT {
    double x, y;
 
    PT () {}
    PT(double x, double y) : x(x), y(y) {}
     
    //Chamemos (*this) de p
    PT operator+(const PT &q) const { return PT(x+q.x, y+q.y); }
    PT operator-(const PT &q) const { return PT(x-q.x, y-q.y); }
    PT operator*(double t) const { return PT(x*t, y*t); }
    PT operator/(double t) const { return PT(x/t, y/t); }
    PT operator-() const { return PT(-x, -y); }
    PT operator[](double t) const { return PT(x*cos(t) - y*sin(t), x*sin(t) + y*cos(t)); } //rotaciona p em t radianos anti-horario
    double operator*(const PT &q) const { return x*q.x + y*q.y; } //produto escalar entre p e q
    double operator%(const PT &q) const { return x*q.y - y*q.x; } //produto cruzado entre p e q
    double operator!() const { return sqrt(x*x + y*y); } //norma de p
    double operator^(const PT &q) const { return atan2(*this%q,*this*q); } //pega o angulo entre p e q
    double operator>(const PT &q) const { return ((*this*q)/(q*q)); } //pega o k da projeção de p em q
};

PT cen, pts[ms];

void print(PT p){
    return;
    cout<<"("<<p.x<<", "<<p.y<<")"<<endl;
}
int main(){
    ll n,a=0,b=1,q;
    ll ans=0;
    scanf("%lld%lld",&n,&q);
    fr(i,0,n){
        scanf("%lf%lf",&pts[i].x,&pts[i].y);
        //cout<<pts[i].x<<" "<<pts[i].y<<endl;
    }
    double mass=0;
    fr(i,2,n){
        PT temp = (pts[0]+pts[i-1]+pts[i])/3;
        double tmass=fabs((pts[i-1]-pts[0])%(pts[i]-pts[0]));
        
    //cout<<cen.x<<" "<<cen.y<<endl;
    //cout<<temp.x<<" "<<temp.y<<endl;
        cen = cen * mass + temp * tmass;
        mass += tmass;
        cen = cen / mass;
    }
    fr(i,0,n){
        pts[i]=pts[i]-cen;
    }
    ld ang=0;
    fr(i,0,q){
        int c;
        scanf("%d",&c);
        if(c==1){
            scanf("%d",&c);
            c--;
            if(b==c)swap(a,b);
            cen=cen+pts[b][ang];
            print(cen);
            ld tang=(pts[b][ang])^PT(0,1);
            
            ang+=tang;
            if(ang>2*pi)ang-=2*pi;
            if(ang<2*pi)ang+=2*pi;
            cen=cen-pts[b][ang];
            print(cen);
            
            scanf("%d",&c);
            c--;
            a=c;
        }else{
            
            scanf("%d",&c);
            c--;
            PT temp=pts[c][ang]+cen;
            printf("%.08lf %.08lf\n",temp.x,temp.y);
        }
    }
    
   // cout<<cen.x<<" "<<cen.y<<endl;
}
#include<bits/stdc++.h>
using namespace std;
const double eps = 1e-15;
double d(double x1, double y1, double x2, double y2){
    return sqrt((x2-x1)*(x2-x1) + (y2-y1)*(y2-y1));
}
int main(){
    double r, x1, x2, y1, y2;
    scanf("%lf%lf%lf%lf%lf", &r, &x1, &y1, &x2, &y2);
    //cout << r << " " << x1 << " " << y1 << " " << x2 << " " << y2 << endl;
    if( ((x1-x2==0) || (fabs(x1-x2)<eps)) && ((y1-y2==0) || (fabs(y1-y2)<eps))){
        //cout << 1 << endl;
        double R = r / 2;
        double x0, y0;
        x0 = (2*x1-r)/2;
        y0 = y1;
        return 0 * printf("%.16f %.16f %.16f\n", x0, y0, R);
    }
    double dd = d(x1, y1, x2, y2);
    if(dd >= r || fabs(dd - r) < eps){
        printf("%.16f %.16f %.16f\n", x1, y1, r);
    }else{
        double R = (r + dd) / 2;
        double x0, y0;
        x0 = (x1-x2)*(r-R)/dd + x1;
        y0 = (y1-y2)*(r-R)/dd + y1;
        printf("%.16f %.16f %.16f\n", x0, y0, R);
    }
    return 0;
}
#include<cstdio>
#include<cmath>
#include<cstring>
using namespace std;
const int size = 110;
const double eps = 1e-12;
struct point3D{double x, y, z;} p[size], outer[4], ret;
int nouter; double radius;
double dissqr3D (point3D p1, point3D p2){
double dx=p1.x - p2.x, dy=p1.y - p2.y, dz=p1.z - p2.z; return dx * dx + dy * dy + dz * dz;}
double dot (point3D p1, point3D p2){return p1.x*p2.x + p1.y*p2.y + p1.z*p2.z;}
void ball (){
point3D q[3]; double m[3][3], sol[3], L[3], det; int i, j;
ret.x = ret.y = ret.z = radius = 0;
switch (nouter){
case 1: ret = outer[0]; break;
case 2:
ret.x = (outer[0].x + outer[1].x) / 2;
ret.y = (outer[0].y + outer[1].y) / 2;
ret.z = (outer[0].z + outer[1].z) / 2;
radius = dissqr3D(ret, outer[0]); break;
case 3:
for (i = 0; i < 2; i++){
q[i].x = outer[i + 1].x - outer[0].x;
q[i].y = outer[i + 1].y - outer[0].y;
q[i].z = outer[i + 1].z - outer[0].z; }
for (i = 0; i < 2; i++) for (j = 0; j < 2; j++) m[i][j] = dot(q[i], q[j]) * 2;
for (i = 0; i < 2; i++) sol[i] = dot(q[i], q[i]);
if (fabs(det = m[0][0]*m[1][1] -m[0][1]*m[1][0]) < eps) return;
L[0] = (sol[0]*m[1][1] - sol[1]*m[0][1]) / det;
L[1] = (sol[1]*m[0][0] - sol[0]*m[1][0]) / det;
ret.x = outer[0].x + q[0].x*L[0] + q[1].x*L[1];
ret.y = outer[0].y + q[0].y*L[0] + q[1].y*L[1];
ret.z = outer[0].z + q[0].z*L[0] + q[1].z*L[1];
radius = dissqr3D(ret, outer[0]); break;
case 4:
for (i = 0; i < 3; i++){
q[i].x = outer[i + 1].x - outer[0].x;
q[i].y = outer[i + 1].y - outer[0].y;
q[i].z = outer[i + 1].z - outer[0].z; sol[i] = dot(q[i], q[i]); }
for (i = 0; i < 3; i++) for (j = 0; j < 3; j++) m[i][j] = dot(q[i], q[j]) * 2;
det = m[0][0]*m[1][1]*m[2][2] + m[0][1]*m[1][2]*m[2][0]
+ m[0][2]*m[2][1]*m[1][0] - m[0][2]*m[1][1]*m[2][0]
- m[0][1]*m[1][0]*m[2][2] - m[0][0]*m[1][2]*m[2][1];
if (fabs(det) < eps) return;
for (j = 0; j < 3; j++){ for (i = 0; i < 3; i++) m[i][j] = sol[i];
L[j] = ( m[0][0]*m[1][1]*m[2][2] + m[0][1]*m[1][2]*m[2][0]
+ m[0][2]*m[2][1]*m[1][0] - m[0][2]*m[1][1]*m[2][0]
- m[0][1]*m[1][0]*m[2][2] - m[0][0]*m[1][2]*m[2][1]) / det;
for (i = 0; i < 3; i++) m[i][j] = dot(q[i], q[j]) * 2; }
ret = outer[0];
for (i = 0; i < 3; i++){
ret.x += q[i].x * L[i];
ret.y += q[i].y * L[i];
ret.z += q[i].z * L[i]; }
radius = dissqr3D(ret, outer[0]); } }

void minball (int n){
    ball ();
    if (nouter < 4)
    for (int i = 0; i < n; i++)
        if (dissqr3D(ret, p[i]) - radius > eps){
            outer[nouter] = p[i];
            ++ nouter;
            minball (i);
            -- nouter;
            if (i > 0){
                point3D tmp = p[i];
                memmove(&p[1], &p[0], sizeof(point3D) * i);p[0] = tmp;
            }
        }
}

int main (){
int n; while (scanf ("%d", &n)!=EOF){
for (int i = 0; i < n; i++) scanf ("%lf%lf%lf", &p[i].x, &p[i].y, &p[i].z);
minball (n); printf ("%.8lf %.8lf %.8lf\n", ret.x+eps,ret.y+eps,ret.z+eps); } }
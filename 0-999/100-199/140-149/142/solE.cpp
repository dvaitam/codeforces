#include <iostream>
#include <cstdio>
#include <cstdlib>
#include <cmath>
#include <climits>
#include <cstring>
#include <string>
#include <vector>
#include <stack>
#include <queue>
#include <set>
#include <map>
#include <algorithm>

using namespace std;

double h, r;
double pi = 3.141592653589793;

bool isGreaterThan(double a, double b)
{
    return a - b > 0.000001;
}

bool isZero(double d)
{
    return d < 0.000001 && d > -0.000001;
}

pair<double, double>  findCoordinatesOnPlane(double x, double y, double z)
{
    double angleOnCircle, theta, sinTheta, cosTheta;
    double d, tmp;

    angleOnCircle = atan2(y, x);
    angleOnCircle += angleOnCircle < 0.0 ? pi + pi : 0.0;

    theta = angleOnCircle * r / sqrt(r * r + h * h);
    sinTheta = sin(theta);
    cosTheta = cos(theta);
    d = (h-z);
    d = (z?sqrt(x * x + y * y + d * d):sqrt(r*r + h*h));

    return make_pair(theta, d);
}

double euclideanDistance(double x1, double y1, double x2, double y2)
{
    double tmp1 = x1 - x2, tmp2 = y1 - y2;
    return sqrt(tmp1 * tmp1 + tmp2 * tmp2);
}

double findConeSurfaceDistance(double x1, double y1, double z1,
                               double x2, double y2, double z2)
{
    pair<double, double> a, b, tmp;
    double tmp1, tmp2, aTheta, bTheta, alpha;
    tmp = findCoordinatesOnPlane(x1, y1, z1);
    b = findCoordinatesOnPlane(x2, y2, z2);
    a = max(tmp, b);
    b = min(tmp, b);
    alpha = 2 * pi * r/ sqrt(r * r + h * h);
    aTheta = a.first;
    bTheta = b.first;
/*    printf("aTheta=%.9lf\n", aTheta * 180.0 / pi);
    printf("bTheta=%.9lf\n", bTheta * 180.0 / pi);
    printf("alpha=%.9lf\n", alpha * 180.0 / pi);
*/
    if(isGreaterThan((aTheta - bTheta), alpha / 2.0))
    {
        // find mirror and do shieat
        aTheta = bTheta + (alpha - aTheta + bTheta);
    }

    return euclideanDistance(a.second * cos(aTheta), a.second * sin(aTheta),
                             b.second * cos(bTheta), b.second * sin(bTheta));
}

double sortaBinarySearch(double left, double right, double xHi, double yHi, double zHi,
                                                    double xLo, double yLo, double zLo)
{
    double x, y, leftVal, rightVal, middle, midVal;
    x = r * cos(left);
    y = r * sin(left);
    leftVal = euclideanDistance(x, y, xLo, yLo) +
             findConeSurfaceDistance(xHi, yHi, zHi, x, y, 0.0);
    x = r * cos(right);
    y = r * sin(right);
    rightVal = euclideanDistance(x, y, xLo, yLo) +
             findConeSurfaceDistance(xHi, yHi, zHi, x, y, 0.0);

    while(!isZero(right - left))
    {
        middle = (left + right) / 2.0;
        x = r * cos(middle);
        y = r * sin(middle);
        midVal = euclideanDistance(x, y, xLo, yLo) +
                 findConeSurfaceDistance(xHi, yHi, zHi, x, y, 0.0);
        if(leftVal > rightVal)
        {
            leftVal = midVal;
            left = middle;
        }
        else
        {
            rightVal = midVal;
            right = middle;
        }
    }
    return min(leftVal, rightVal);
}

double findOptimalSingleBrinkPointDistance(double xHi, double yHi, double zHi,
                                           double xLo, double yLo, double zLo)
{
    pair<double, double> tmpCoord;
    double minDist;
    double left = 0, right = 2 * pi, middle;
    double leftVal, rightVal, midVal, prevMinDegree;
    double curr;
    double x, y;
    double aDegree = pi / 180.0, twoDegrees = pi / 90.0;
    double minDegree, minVal;

    x = r * cos(left);
    y = r * sin(left);
    minVal = euclideanDistance(x, y, xLo, yLo) +
             findConeSurfaceDistance(xHi, yHi, zHi, x, y, 0.0);
    minDegree = left;

    for(left = aDegree; left < right; left += aDegree)
    {
        x = r * cos(left);
        y = r * sin(left);
        curr = euclideanDistance(x, y, xLo, yLo) + 
               findConeSurfaceDistance(x, y, 0, xHi, yHi, zHi);
        if(curr < minVal)
        {
            prevMinDegree = minDegree;
            minVal = curr;
            minDegree = left;
        }
    }

    left = minDegree - aDegree * 2;
    right = minDegree + aDegree * 2;

//    minDegree = atan2(yHi, xHi);
    minVal = sortaBinarySearch(minDegree - aDegree, minDegree + aDegree, xHi, yHi, zHi, xLo, yLo, zLo);
    minVal = min(minVal, sortaBinarySearch(minDegree - 2 * aDegree, minDegree + 2 * aDegree, xHi, yHi, zHi, xLo, yLo, zLo));

    return minVal;
}

double sortaBinarySearch2(double left, double right, double x1, double y1, double z1,
                                                     double x2, double y2, double z2)
{
    double leftVal, rightVal, midVal, middle, x, y;

    x = r * cos(left);
    y = r * sin(left);
    leftVal = findConeSurfaceDistance(x1, y1, z1, x, y, 0.0) +
              findOptimalSingleBrinkPointDistance(x2, y2, z2, x, y, 0.0);
    x = r * cos(right);
    y = r * sin(right);
    rightVal = findConeSurfaceDistance(x1, y1, z1, x, y, 0.0) +
               findOptimalSingleBrinkPointDistance(x2, y2, z2, x, y, 0.0);

    while(!isZero(right - left))
    {
        middle = (left + right) / 2.0;
        x = r * cos(middle);
        y = r * sin(middle);
        midVal = findConeSurfaceDistance(x1, y1, z1, x, y, 0.0);
        midVal += findOptimalSingleBrinkPointDistance(x2, y2, z2, x, y, 0.0);
//        printf("middle=%.9lf - left=%.9lf - right=%.9lf\n", middle * 180.0 / pi, left * 180.0/pi, right * 180.0/ pi);
//        printf("x=%.9lf - y=%.9lf\n", x, y);
//        cout << "leftVal=" << leftVal << " - rightVal=" << rightVal << " - midVal=" << midVal << endl;
        if(leftVal > rightVal)
        {
            leftVal = midVal;
            left = middle;
        }
        else
        {
            rightVal = midVal;
            right = middle;
        }
    }
    return min(leftVal, rightVal);
}

double findOptimalDoubleBrinkPointDistance(double x1, double y1, double z1,
                                           double x2, double y2, double z2)
{
    double minDegree, minVal;
    double left = 0, right = 2 * pi, middle;
    double leftVal, rightVal, midVal;
    double curr;
    double x, y;
    double aDegree = pi / 180.0, twoDegrees = pi / 90.0;

    minDegree = atan2(y1, x1);
    minVal = sortaBinarySearch2(minDegree - aDegree, minDegree + aDegree, x1, y1, z1, x2, y2, z2);
    minVal = min(minVal, sortaBinarySearch2(minDegree - 2 * aDegree, minDegree + 2 * aDegree, x1, y1, z1, x2, y2, z2));
    minVal = min(minVal, sortaBinarySearch2(minDegree - 3 * aDegree, minDegree + 3 * aDegree, x1, y1, z1, x2, y2, z2));
    return minVal;
}

int main()
{
    pair<double, double> a, b;
    double x1, y1, z1, x2, y2, z2;
    double tmp1, tmp2, res;
    bool pt1OnBase, pt2OnBase;

    cin >> r >> h;

    cin >> x1 >> y1 >> z1;
    cin >> x2 >> y2 >> z2;
    pt1OnBase = isZero(z1);
    pt2OnBase = isZero(z2);

    if(pt1OnBase)
    {
        if(pt2OnBase)
        {
            tmp1 = x1 - x2;
            tmp2 = y1 - y2;
            res = sqrt(tmp1 * tmp1 + tmp2 * tmp2);
        }
        else
        {
            res = findOptimalSingleBrinkPointDistance(x2, y2, z2, x1, y1, z1);
            if(isZero(x1 * x1 + y1 * y1 - r * r))
                res = min(res, findConeSurfaceDistance(x1, y1, z1, x2, y2, z2));
        }
    }
    else if(pt2OnBase)
    {
        res = findOptimalSingleBrinkPointDistance(x1, y1, z1, x2, y2, z2);
        if(isZero(x2 * x2 + y2 * y2 - r * r))
            res = min(res, findConeSurfaceDistance(x1, y1, z1, x2, y2, z2));
    }
    else
    {
        res = findConeSurfaceDistance(x1, y1, z1, x2, y2, z2);
        res = min(res, findOptimalDoubleBrinkPointDistance(x1, y1, z1, x2, y2, z2));
    }

    printf("%.9lf\n", res);

    return 0;
}
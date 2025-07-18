#include <iostream>

#include <cstring>

#include <string>

#include <stdio.h>

#include <algorithm>

#include <vector>

#include <set>

#include <map>

#include <queue>

#include <fstream>

#include <math.h>

using namespace std;	

	

		

#define ull unsigned long long



const int bval = 200000 + 10;



const int UPLIMIT = 1e9 + 3;







int main(){



	int a, b;

	cin >> a >> b;

	

	int n;

	cin >> n;

	int x, y, v;

	double mintime = 1e9;

	for (int i = 0; i < n; i++){

		cin >> x >> y >> v;

		double sqr = sqrt((x - a)  * (x - a) +  (y - b) * (y - b));

		if (sqr / v < mintime)

			mintime = sqr / v;

		

	}

	

	printf("%f",mintime);

	

	return 0;

	

}
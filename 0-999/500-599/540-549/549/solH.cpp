#include <stdio.h>

#include <algorithm>

#include <iostream>

#include <vector>

#include <map>

#include <set>

#include <stdlib.h>

#include <math.h>

#include <time.h>

#include <string.h>

using namespace std;

typedef long long int ll;

ll a, b, c, d;

int main(void)

{

	cin>>a>>b>>c>>d;

	ll t;

	t=0;

	t=max(t,abs(a+b+c+d));

	t=max(t,abs(a-b+c-d));

	t=max(t,abs(a+b-c-d));

	t=max(t,abs(a-b-c+d));

	if(t==0)

	{

		cout<<0;

	}

	else

	{

		printf("%.10lf",(double)(abs(a*d-b*c))/t);

	}

}
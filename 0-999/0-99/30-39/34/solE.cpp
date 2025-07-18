#include <cstdio>

#include <algorithm>

#include <cmath>



using namespace std;



const int N = 15;

const double ERROR = 1e-9;

const double INFINITE = 1e9;



double x[N];

double v[N];

int m[N];

int n;



double calTime(int i, int j)

{

	if (fabs(v[i] - v[j]) > ERROR && fabs(x[j] - x[i]) > ERROR && (v[i] - v[j]) * (x[j] - x[i]) > 0)

	{

		return (x[j] - x[i]) / (v[i] - v[j]);

	}

	

	return INFINITE;

}



void collid(int i, int j)

{

	double v1 = v[i];

	double v2 = v[j];

	int m1 = m[i];

	int m2 = m[j];

	v[i] = ((m1 - m2) * v1 + 2*m2*v2) / (m1 + m2);

	v[j] = ((m2 - m1) * v2 + 2*m1*v1) / (m1 + m2);

}



void checkCollision()

{

	int i, j;

	for (i = 0; i < n; i++)

	{

		for (j = i + 1; j < n; j++)

		{

			if (fabs(x[i] - x[j]) < ERROR)

			{

				collid(i, j);

			}

		}

	}

}



void passTime(double t)

{

	for (int i = 0; i < n; i++)

	{

		x[i] += v[i] * t;

	}

}



int main()

{

	int i, j;

	double T;

	scanf("%d%lf", &n, &T);

	for (i = 0; i < n; i++)

	{

		scanf("%lf%lf%d", &x[i], &v[i], &m[i]);

	}

	

	double colTime = INFINITE;

	while (T > 0)

	{

		double t;

		colTime = INFINITE;

		for (i = 0; i < n; i++)

		{

			for (j = i + 1; j < n; j++)

			{

				t = calTime(i, j);

				if (t < colTime && t < T)

				{

					colTime = t;

				}

			}

		}



		if (colTime == INFINITE)

		{

			passTime(T);

			break;

		}	

		

		passTime(colTime);

		checkCollision();

		T -= colTime;

	}

	

	for (i = 0; i < n; i++)

	{

		printf("%.6lf\n", x[i]);

	}

	printf("\n");

}
# Pseudorandom number generators

## Description
This is simple application to test implementations of random number generators.

For now it supports LGC, and two ways of generating random numbers in the Cauchy distribution.

By default the app generates 10 000 000 numbers for each generator, prints relevant statistics for the distribution tested and creates a histogram of these numbers.

## LGC
The implementation uses `glibc` parameters and of course tries to provide numbers from the uniform distribution. Generated numbers are positive integers up to $2^{31}$. Statistics are calculated for normalized values the are in range of [0,1].

Theoretical value of the mean should be 0.5, and standard deviation should be $\frac{1}{2\sqrt{3}}$ which is around 0,288675134594813.

Below are statistics calculated by the program.
```
Generating numbers using LGC took 40.918709ms
LGC stats (normalized):
	Mean: 0.499957
	Standard deviation: 0.288711
```

Statistics are quite good and close to theoretical values.

Below is the histogram of generated numbers.

![](./plots/Linear%20Congruential%20Generator.png)

Theoretically the histogram should be flatter but of course the generator is not perfect and some numbers repeat more often than others.

## Cauchy

To generate numbers in the Cauchy distribution two algorithms have been used: simple approach using tangent and an alternative that leverages the probability density function - tangent is not needed here.

`X0` and `Gamma` parameters of the distribution are set to 0 and 1 respectively - they can be adjusted in the code.

Since the Cauchy distribution has no mean and standard deviation the 1st quartile, median, 3rd quartile and interquartile range are calculated. With above mentioned parameters the theoretical values for these statistics should be -1, 0, 1 and 2 respectively.

Below are calculated statistics.
```
Generating Cauchy numbers took 343.272875ms
Cauchy stats:
	1st quartile: -0.999538
	Median: -0.000865
	3rd quartile: 0.998840
	Interquartile range: 1.998378

Generating Cauchy numbers without using tangent took 280.397542ms
Cauchy no tangent stats:
	1st quartile: -0.998990
	Median: 0.001031
	3rd quartile: 1.001328
	Interquartile range: 2.000317
```

Calculated statistics for the method using tangent are closer to theoretical values but noticeably generating numbers using the alternative method was ~60ms faster.

It's important to mention that the Cauchy distribution has extreme values and the "heavy tails" that this distribution has would cause the histogram to consist of a single high bar. This is because of the limited amount of bins that can be used. Due to this issue the histogram is created only for values in the [-4, 4] range. 

Below are histograms for the tangent method and alternative methods.

![](plots/Simple%20Cauchy%20Generator.png)

![](plots/Fast%20Cauchy%20Generator.png)

Histograms are similar to each other and close to the actual distribution but of course they are not perfect.
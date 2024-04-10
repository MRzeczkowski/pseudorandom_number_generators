# Pseudorandom number generators

## Description
This is a simple application to test implementations of random number generators.

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

To generate numbers in the Cauchy distribution two algorithms have been used: simple approach using tangent and an alternative that leverages the probability density function and some clever maths - tangent is not needed here so it's a bit faster.

`X0` and `Gamma` parameters of the distribution are by default set to 0 and 1 respectively - they can be adjusted in the code.

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

Another important note is that the fast method uses an algorithm from this book "Wieczorkowski R. - Komputerowe generatory liczb losowych", specifically "Algorytm 3.37" on pages 75-76. This algorithm is used to generate numbers in range [-1, 1] for default values of `X0` and `Gamma`. Because it does not generate tails I've added some code inspired by [this blog post](https://devzine.pl/2011/02/21/generator-liczb-pseudolosowych-cz-3-rozklad-cauchyego/). Maybe there is a more performant way of doing this but I didn't find one online. I modify the output for the algorithm and `X0` and `Gamma` work with it correctly.

Below are histograms for the tangent method and the alternative, faster method.

![](plots/Simple%20Cauchy%20Generator.png)

![](plots/Fast%20Cauchy%20Generator.png)

Histograms are similar to each other and close to the actual distribution but of course they are not perfect.
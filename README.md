# GOvPython

This is a PoC (Proof of Concept) that explores the possibilities of implementing Machine Learning in Go and Benchmarking that against a more straight forward implementation in Python.

This PoC explores the use case of a Random Forest algorithm to build a model that can predict whether a flight will be delayed based on historical flight data from 2015.

Data is sourced from Kaggle -> you can download data here https://www.kaggle.com/datasets/usdot/flight-delays

NOTE! Data contains millions of records for the purpose of this mini project I have limited the data to 100000 records for ease of processing (I recommend you do the same unless you have a custom Nvidia Server :D)

Documentation for packages/libraries use:
Python:
Pandas -> https://pandas.pydata.org/pandas-docs/stable/index.html
Sklearn -> https://scikit-learn.org/stable/modules/classes.html

Go (These are just the tools i found, I'm sure there might be even better tools out there):
Gota (Go's alternative to Pandas) -> https://pkg.go.dev/github.com/go-gota/gota
Gonum (Go's alternative to Numpy/Scipy) -> https://github.com/gonum/matrix -> This requires the CGO GCC compiler to be available on your local machine, GCC allows you to run native C++ code in your Go code/compiler
Golearn (Go's alternative to SciKit-Learn) -> https://pkg.go.dev/github.com/sjwhitworth/golearn

Benchmarking (I tested this on my Windows Intel(x64) laptop so you might need to investigate tools for other computer hardware architectures and Operating Systems):
Intel Power Gadget -> https://www.intel.com/content/www/us/en/developer/articles/tool/power-gadget.html
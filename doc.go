/*
Package retry manages the execution of certain pieces of code that must run a
specific number of times before being considered as failed.

The execution of a function will be considered failed only if all the tries
returned an error, when it happens the 'With*' function will return all the
errors in a ErrorGroup.

If the TryFunc return a nil error at any moment it will be considered a
successful execution and nil will be returned by the 'With*' function.

To know more about the ErrorGroup go to https://github.com/txgruppi/errorgroup-go

How many times TryFunc will execute?

The general rule is 'len(BackoffArray) + 1'

This happens because the BackoffArray is an array of intervals and makes no
sense to have an interval after the last execution. Because of this every
time the last interval is extracted from the BackoffArray the TryFunc is
executed one more time. If you want to run the function only one time you
should give a BackoffArray with zero items.
*/
package retry

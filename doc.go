/*
Package registry defines primitives for keeping centralized repositories
of qri types (peers, datasets, etc). It uses classical client/server patterns,
arranging types into cannonical stores.

At first glance, this seems to run against the grain of "decentralize or die"
principles espoused by those of us interested in reducing points of failure in
a network. Consider this package testiment that nothing is absolute.

It is a long term goal at qri that it be *possible* to fully decentralize all
aspects, of qri this isn't practical short-term, and isn't always a desired
property.

As an example, associating human-readable usernames with crypto keypairs is an
order of magnitude easier if you just put the damn thing in a list. So that's
what this registry does.

Long term, we intended to implement a distributed hash table (DHT) to make it
possible to operate fully-decentralized, and provide registry support as a
configurable detail.

This base package provides common primitives that other packages can import to
work with a registry, and subpackages for turning these primitives into usable
tools like servers & (eventually) command-line clients
*/
package registry

# Set practice: coding & using sets in Go

## Elevator Pitch

Boolean logic and sets are closely related. That's why set operations like intersection, difference and sub-set testing can replace a lot of complex logic with nested loops and ifs. We'll discuss why and how to implement sets in Go, from the explicit use of 'map' to generic code with the 'gen' tool.


## Description

This talk has two main objectives: to show the practical use of set operations to simplify Go programming and to discuss alternatives for implementing a set collection type in Go.


### Set Advocacy

Consider this story in a project: “Mark all favored products, except those already in the shopping cart.” Or this one: “List only products whose description contains all search words”. These functionalities can be thought algorithmically through loops and conditionals, but they can also be understood declaratively, briefly and precisely, as relations between sets: the first is the difference `F – C` (favored minus cart), and the second is a sub-set test, `B ⊆ D` (search words are contained in description).

Sets support operations like union, intersection, difference, and containment tests. Databases, query languages, and the standard library of many languages provide those operations because they are very useful for selecting and filtering data.

The lack of sets in the Go standard library means that many teams are coding the logic of sets in various parts of their projects, sometimes without realizing they are reinventing the wheel. Understanding such functions as operations between sets simplifies communication within the team, and creates the conditions to leverage correct, well-tested, efficient implementations.


### Discuss ways of implementing sets

The only generic collection types in Go are the built-ins: array, slice, map, and channel. Without user-defined generics, the first question when implementing a set is: a set of what? What is the type of the set elements? As long as Go does not add a generic set built-in we need to decide how to implement set operations or choose one of the many existing set implementations.

There are several ways to get sets in Go, ranging from explicit use of a `map` all the way to using code-generation to produce type-safe custom implementations. Another consideration: do you need a `Set` that that is safe for concurrent use?

There is no single ideal solution for all use cases, so we`ll discuss the pros and cons of these approaches to make good choices when we want to leverage the clarity and expressiveness of sets to simplify logic in our projects.


### Talk outline

* Sets and boolean logic
* Sets in other languages
* Practical use cases for sets in Go code
* Implementation alternatives
   * Basic set operations on a `map`
   * Encapsulating `map` in a custom `Set`
   * `Set` type for `interface {}`
   * Survey of existing implementations
   * Using the `gen` tool to automatically generate custom implementations.
   * Pros and cons matrix

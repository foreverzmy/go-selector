go-selector
===========

[![Build Status](https://travis-ci.org/blendlabs/go-selector.svg?branch=master)](https://travis-ci.org/blendlabs/go-selector)

Selector is a library that matches as closely as possible the intent and semantics of kubernetes selectors.

It supports unicode in names (such as `함=수`), but does not support escaped symbols.

## BNF
```
  <selector-syntax>         ::= <requirement> | <requirement> "," <selector-syntax>
  <requirement>             ::= [!] KEY [ <set-based-restriction> | <exact-match-restriction> ]
  <set-based-restriction>   ::= "" | <inclusion-exclusion> <value-set>
  <inclusion-exclusion>     ::= <inclusion> | <exclusion>
  <exclusion>               ::= "notin"
  <inclusion>               ::= "in"
  <value-set>               ::= "(" <values> ")"
  <values>                  ::= VALUE | VALUE "," <values>
  <exact-match-restriction> ::= ["="|"=="|"!="] VALUE
```

## Usage

Fetch the package as normal:
```bash
> go get -u github.com/blendlabs/go-selector
```

Include in your project:
```golang
import selector "github.com/blendlabs/go-selector"
```

## Example

Given a label collection:
```golang
valid := selector.Labels{
  "zoo":   "mar",
  "moo":   "lar",
  "thing": "map",
}
```

We can then compile a selector:

```golang
selector, _ := selector.Parse("zoo in (mar,lar,dar),moo,thing == map,!thingy")
fmt.Println(selector.Matches(valid)) //prints `true`
```

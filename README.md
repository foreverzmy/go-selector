go-selector
===========

[![Build Status](https://travis-ci.org/blendlabs/go-selector.svg?branch=master)](https://travis-ci.org/blendlabs/go-selector)

Selector is a library that matches as closely as possible the intent and semantics of kubernetes selectors.

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

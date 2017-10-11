go-selector
===========

Selector is a library that matches as closely as possible the intent and symantics of kubernetes selectors.

## BNF
  <selector-syntax>         ::= <requirement> | <requirement> "," <selector-syntax>
  <requirement>             ::= [!] KEY [ <set-based-restriction> | <exact-match-restriction> ]
  <set-based-restriction>   ::= "" | <inclusion-exclusion> <value-set>
  <inclusion-exclusion>     ::= <inclusion> | <exclusion>
  <exclusion>               ::= "notin"
  <inclusion>               ::= "in"
  <value-set>               ::= "(" <values> ")"
  <values>                  ::= VALUE | VALUE "," <values>
  <exact-match-restriction> ::= ["="|"=="|"!="] VALUE
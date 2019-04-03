
# Babelfish UAST

_[Babelfish documentation](https://docs.sourced.tech/babelfish/): specifications, usage, examples._

## Definition

One of the most important components of source{d} Engine is the UAST, which stands for:
[Universal Abstract Syntax Tree](https://docs.sourced.tech/babelfish/uast/uast-specification).

UASTs are a normalized form of a programming language's AST, annotated with language-agnostic roles and transformed with language-agnostic concepts (e.g. Functions, Imports etc.).

These enable advanced static analysis of code and easy feature extraction for statistics or [Machine Learning on Code](https://github.com/src-d/awesome-machine-learning-on-source-code).

## UAST Usage

To parse a file for a UAST using source{d} Engine, head to the [Parsing Code section](#parsing-code) of this document.

## Supported Languages

To see which languages are available, check the table of [Babelfish supported languages](https://docs.sourced.tech/babelfish/languages).

## Clients and Connectors

For connecting to the language parsing server (Babelfish) and analyzing the UAST, there are several language clients currently supported and maintained:

- [Babelfish Go Client](https://github.com/bblfsh/client-go)
- [Babelfish Python Client](https://github.com/bblfsh/client-python)
- [Babelfish Scala Client](https://github.com/bblfsh/client-scala)

The Gitbase Spark connector is under development, which aims to allow for an easy integration with Spark & PySpark:

- [Gitbase Spark Connector](https://github.com/src-d/gitbase-spark-connector)- more coming soon!

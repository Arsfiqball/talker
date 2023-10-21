# Power Context (Poco)

Poco is a package that focuses on powering the propagated context of a program.
It enables the program to simply specify what it does without needing to worry
about how to report it to the instrumenting system.

> Poco is inspired by Powercode Talker, a cyberse monster from Yu-Gi-Oh! VRAINS.

## Features

- [x] Embeddable observer abstraction (to tell span, event, error)
- [x] Error typing with stack trace, error code, and error message
- [x] Utility to recover from panic and report it as error

## Installation

This package is part of [Talker](https://github.com/Arsfiqball/talker)
project. To install it, simply run:

```bash
go get -u github.com/Arsfiqball/talker
```

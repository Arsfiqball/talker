# Execution Controller (Exco)

Exco is a package that focuses on controlling the execution of a program,
specifically managing the execution order of sequential processes, parallel
operations, and application runtime. These types of tasks typically occur
during initialization, health checks, and graceful shutdown. Therefore, this
package helps to structure them in a more organized manner.

> Exco is inspired by Excode Talker, a cyberse monster from Yu-Gi-Oh! VRAINS.

## Features

- [x] Sequential execution
- [x] Parallel execution
- [x] Process management (init, health check, graceful shutdown)

## Installation

This package is part of [Talker](https://github.com/Arsfiqball/talker)
project. To install it, simply run:

```bash
go get -u github.com/Arsfiqball/talker
```

## Usage

### Sequential Execution

Sequential execution is a process that runs a series of tasks in a sequential
order. The next task will only be executed if the previous task is successful.

```go
cb := exco.Sequential(
    func(ctx context.Context) error {
        // do something
        return nil
    },
    exco.Timeout(
        func(ctx context.Context) error {
            // do something
            return nil
        },
        30 * time.Second, // timeout after 30 seconds
    ),
    exco.Retry(
        func(ctx context.Context) error {
            // do something
            return nil
        },
        5, // retry 5 times
        5 * time.Second, // wait 5 seconds between each retry
    ),
)

err := cb(context.Background())
```

### Parallel Execution

Parallel execution is a process that runs a series of tasks in parallel. The
next task will only be executed if all previous tasks are successful.

```go
cb := exco.Parallel(
    func(ctx context.Context) error {
        // do something
        return nil
    },
    exco.Timeout(
        func(ctx context.Context) error {
            // do something
            return nil
        },
        30 * time.Second, // timeout after 30 seconds
    ),
    exco.IgnoreError(func(ctx context.Context) error {
        // do something
        return errors.New("this error will be ignored")
    }),
)

err := cb(context.Background())
```

### Process Management

```go
proc := exco.Process{
    MonitorAddr: ":8086",
    Start: func(ctx context.Context) error {
        // do something
        return nil
    },
    Live: func(ctx context.Context) error {
        // do something
        return nil
    },
    Ready: exco.Parallel(
        func(ctx context.Context) error {
            // do something
            return nil
        },
        // readiness check of http server at port 8080
        exco.HttpGetCheck("http://localhost:8080/readiness"),
    ),
    Stop: func(ctx context.Context) error {
        // do something
        return nil
    },
}

sig := make(chan os.Signal, 1)
signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

exco.Serve(proc, sig)
```

## Maintainer

- Iqbal Mohammad Abdul Ghoni - [Arsfiqball](https://github.com/Arsfiqball)

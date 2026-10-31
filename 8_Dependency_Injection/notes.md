# Go: Buffers, io.Writer, and Testing Output

## 1. The problem: How do we test printed output?

Consider a simple greeting function:

```go
func Greet(name string) {
    fmt.Printf("Hello, %s", name)
}
```

`fmt.Printf` prints directly to **standard output (stdout)**, which usually means the terminal.

The problem is that when we test this function, we want to capture its output and check whether it matches what we expect.

For example:

- Expected: `"Hello, Chris"`
- Actual: `"Hello, Chris"`

But how do we capture the output instead of sending it directly to the terminal?

The solution is to **inject the destination where the output should be written**.

Instead of making `Greet` decide where to print, we let the caller provide a destination.

This makes the function more flexible and easier to test.

---

## 2. Understanding bytes.Buffer

A `bytes.Buffer` is an in-memory storage area that holds a sequence of bytes.

Think of it as a temporary container into which you can write output and retrieve it later.

```go
var buffer bytes.Buffer

buffer.WriteString("Hello")
buffer.WriteString(", World!")

fmt.Println(buffer.String())
```

Output:

```text
Hello, World!
```

What happened?

1. We created an empty buffer in memory.
2. `WriteString()` added data to the buffer.
3. The buffer retained the written data.
4. `String()` allowed us to retrieve the accumulated contents.

### Why is this useful?

Compare two destinations:

| stdout (terminal) | bytes.Buffer |
|---|---|
| Sends output to the terminal | Stores output in memory |
| Primarily for displaying output | Useful for capturing output |
| Not convenient to retrieve programmatically | Can retrieve contents later |
| Output goes to an external destination | Output stays in memory |

The important intuition:

**A buffer gives us a place to accumulate output before deciding what to do with it.**

We can inspect it, read it, append more data, or pass its contents somewhere else.

### Useful bytes.Buffer methods

```go
var buffer bytes.Buffer

buffer.WriteString("Hello")

// Retrieve contents as a string
fmt.Println(buffer.String())

// Retrieve contents as a byte slice
fmt.Println(buffer.Bytes())

// Remove existing contents
buffer.Reset()
```

`String()` returns the buffer's contents as a string, whereas `Bytes()` returns the contents as a byte slice (`[]byte`).

---

## 3. Understanding io.Writer

Before changing our `Greet` function, we need to understand `io.Writer`.

Go defines an interface called `io.Writer`:

```go
type Writer interface {
    Write(p []byte) (n int, err error)
}
```

An `io.Writer` represents anything that can accept bytes.

In other words, it provides a standard way of writing data somewhere.

Different destinations can implement this same interface:

| Destination | Where the data goes |
|---|---|
| `os.Stdout` | Terminal |
| `bytes.Buffer` | In-memory storage |
| `os.File` | File on disk |
| `http.ResponseWriter` | HTTP response |

All these types implement the `Write([]byte)` method, so they satisfy the `io.Writer` interface.

### Why does this matter?

Go allows us to write functions that accept an `io.Writer` instead of requiring a specific destination type.

This means our function does not need to know where the output is going.

It only needs to know that the destination can accept bytes.

**Remember:**

- `io.Writer` describes the capability of writing data.
- `bytes.Buffer` is one implementation that stores written data in memory.
- `os.Stdout` is another implementation that writes data to the terminal.

---

## 4. How fmt.Printf and fmt.Fprintf work

This is an important connection.

Consider:

```go
fmt.Printf("Hello, World!")
```

Under the hood, `fmt.Printf` uses `fmt.Fprintf`:

```go
func Printf(format string, a ...interface{}) (n int, err error) {
    return Fprintf(os.Stdout, format, a...)
}
```

Notice that `Printf` passes `os.Stdout` as the destination.

Now look at `Fprintf`:

```go
func Fprintf(w io.Writer, format string, a ...interface{}) (n int, err error) {
    // Formats the data and writes it to w
}
```

The first argument is an `io.Writer`.

This explains the difference:

```go
fmt.Printf("Hello")
```

Writes to standard output by default.

```go
fmt.Fprintf(writer, "Hello")
```

Writes to the destination supplied through `writer`.

For example:

```go
var buffer bytes.Buffer

fmt.Fprintf(&buffer, "Hello")

fmt.Println(buffer.String())
```

Output:

```text
Hello
```

Instead of printing directly to the terminal, we directed the output into a buffer.

---

## 5. Making Greet testable using dependency injection

### Step 1: Write the test first

```go
func TestGreet(t *testing.T) {
    buffer := bytes.Buffer{}

    Greet(&buffer, "Chris")

    got := buffer.String()
    want := "Hello, Chris"

    if got != want {
        t.Errorf("got %q want %q", got, want)
    }
}
```

Let's understand what happens:

1. We create an empty `bytes.Buffer`.
2. We pass its pointer to `Greet`.
3. `Greet` writes the greeting into the buffer.
4. We retrieve the contents using `buffer.String()`.
5. We compare the actual output against the expected output.

Because the output is stored in memory, we can inspect it after calling the function.

### Step 2: Implement Greet

Initially, we might write:

```go
func Greet(writer *bytes.Buffer, name string) {
    fmt.Fprintf(writer, "Hello, %s", name)
}
```

This works with our test because `bytes.Buffer` implements `io.Writer`.

However, there is a problem.

What if we want to use the same function in our actual application?

```go
func main() {
    Greet(os.Stdout, "Elodie")
}
```

The compiler complains because `os.Stdout` is a `*os.File`, not a `*bytes.Buffer`.

Our function is unnecessarily restrictive. It only accepts one specific type.

### Step 3: Use the io.Writer interface

Instead of accepting `*bytes.Buffer`, accept the more general `io.Writer` interface.

```go
func Greet(writer io.Writer, name string) {
    fmt.Fprintf(writer, "Hello, %s", name)
}
```

Now both of these work:

**In our test:**

```go
func TestGreet(t *testing.T) {
    buffer := bytes.Buffer{}

    Greet(&buffer, "Chris")

    got := buffer.String()
    want := "Hello, Chris"

    if got != want {
        t.Errorf("got %q want %q", got, want)
    }
}
```

**In our actual application:**

```go
func main() {
    Greet(os.Stdout, "Elodie")
}
```

The same function works with both destinations.

This is because both `*bytes.Buffer` and `*os.File` implement `io.Writer`.

### Complete implementation

```go
package main

import (
    "fmt"
    "io"
    "os"
)

func Greet(writer io.Writer, name string) {
    fmt.Fprintf(writer, "Hello, %s", name)
}

func main() {
    Greet(os.Stdout, "Elodie")
}
```

---

## 6. What is dependency injection?

Dependency injection simply means **providing a function with the dependencies it needs instead of making it create or decide those dependencies itself.**

In our example, `Greet` needs somewhere to write its output.

Instead of deciding that the output must go to the terminal, we inject the destination:

```go
func Greet(writer io.Writer, name string) {
    fmt.Fprintf(writer, "Hello, %s", name)
}
```

The caller decides where the output goes.

```go
// Write to memory
Greet(&buffer, "Chris")

// Write to terminal
Greet(os.Stdout, "Elodie")
```

This gives us two important benefits:

1. **Testability:** We can inject a buffer and inspect the output.
2. **Reusability:** We can use the same function with different destinations without changing its implementation.

The key idea is to depend on an interface describing what something can do, rather than a concrete type describing exactly what it is.

---

## 7. Using io.Writer with the Internet (HTTP)

One interesting application of `io.Writer` is HTTP servers.

Consider this code:

```go
func Greet(writer io.Writer, name string) {
    fmt.Fprintf(writer, "Hello, %s", name)
}

func MyGreeterHandler(w http.ResponseWriter, r *http.Request) {
    Greet(w, "world")
}

func main() {
    log.Fatal(http.ListenAndServe(":5001",
        http.HandlerFunc(MyGreeterHandler)))
}
```

Let's understand what happens.

### What is an HTTP handler?

When someone visits:

```text
http://localhost:5001
```

Their browser sends an HTTP request to our Go server.

Go calls our handler:

```go
func MyGreeterHandler(w http.ResponseWriter, r *http.Request)
```

The handler receives two important things:

- `r`: The incoming HTTP request.
- `w`: The destination where we write the HTTP response.

### What is http.ResponseWriter?

`http.ResponseWriter` is an interface provided by Go's `net/http` package.

It allows an HTTP handler to construct and send a response back to the client.

Importantly, `http.ResponseWriter` also implements `io.Writer`.

This means we can pass it directly to our `Greet` function:

```go
Greet(w, "world")
```

Remember, `Greet` only requires an `io.Writer`.

It doesn't care whether the destination is:
- A terminal.
- An in-memory buffer.
- An HTTP response.

As long as the destination implements `io.Writer`, it works.

### What happens internally?

The flow looks like this:

```text
Browser sends HTTP request
          |
          v
     Go HTTP Server
          |
          v
   MyGreeterHandler
          |
          v
     Greet(w, "world")
          |
          v
    fmt.Fprintf(writer, ...)
          |
          v
   http.ResponseWriter
          |
          v
   HTTP response sent
          |
          v
Browser displays "Hello, world"
```

The important connection:

We originally created `Greet` to make printing testable.

But because we used the general-purpose `io.Writer` interface, we can also reuse the exact same function to write HTTP responses.

No changes to `Greet` are necessary.

---

## 8. The bigger picture: One function, multiple destinations

This is the central lesson of the entire concept.

```go
func Greet(writer io.Writer, name string) {
    fmt.Fprintf(writer, "Hello, %s", name)
}
```

Our function doesn't care where the greeting goes.

The caller decides:

```go
// 1. Store output in memory
var buffer bytes.Buffer
Greet(&buffer, "Chris")

// 2. Print output to terminal
Greet(os.Stdout, "Elodie")

// 3. Write output to an HTTP response
Greet(httpResponseWriter, "world")
```

All three destinations satisfy the same interface: `io.Writer`.

### Final mental model

| Concept | Intuition |
|---|---|
| `bytes.Buffer` | A temporary in-memory container for accumulating and retrieving bytes. |
| `io.Writer` | A common interface for anything that accepts bytes. |
| `os.Stdout` | A writer that sends bytes to the terminal. |
| `fmt.Printf` | Prints to standard output by default. |
| `fmt.Fprintf` | Writes formatted output to a supplied writer. |
| Dependency injection | Passing a dependency into a function rather than deciding it inside the function. |
| `http.ResponseWriter` | A writer that sends data as an HTTP response. |

**The most important takeaway:**

Don't make a function responsible for deciding where its output goes if that decision can be made by the caller.

Accept an `io.Writer`, and your function becomes easier to test, more flexible, and reusable in different contexts.
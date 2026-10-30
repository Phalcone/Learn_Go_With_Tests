
# Error Handling in Go

## 1. What is an error in Go?

In Go, an error is simply a value that represents something going wrong during the execution of a function or method.

Unlike some languages that use exceptions, Go commonly handles errors by returning them from functions.

For example, a wallet might not have enough Bitcoin for a withdrawal.

```go
func (w *Wallet) Withdraw(amount Bitcoin) error {
    if amount > w.balance {
        return errors.New("Oh, no!")
    }

    w.balance -= amount
    return nil
}
```

The method has two possible outcomes:

- If the withdrawal fails, it returns an error.
- If the withdrawal succeeds, it returns `nil`.

**Remember:**
- `nil` means no error occurred.
- A non-nil error means something went wrong.

## 2. Understanding the error interface

In Go, `error` is a built-in interface.

Its definition is:

```go
type error interface {
    Error() string
}
```

This means any type that implements the `Error() string` method satisfies the `error` interface.

For example:

```go
err := errors.New("insufficient funds")
```

`errors.New()` creates an error value containing the supplied message.

You can retrieve its message using:

```go
fmt.Println(err.Error())
```

Output:

```text
insufficient funds
```

The important idea is that `error` is an interface, not a string.

It represents a value that can provide an error message through its `Error()` method.

## 3. Creating and returning errors

The `errors` package provides the `New()` function to create simple errors.

```go
func (w *Wallet) Withdraw(amount Bitcoin) error {
    if amount > w.balance {
        return errors.New("insufficient funds")
    }

    w.balance -= amount
    return nil
}
```

Here:

```go
return errors.New("insufficient funds")
```

Creates an error and returns it to whoever called `Withdraw()`.

And:

```go
return nil
```

Indicates that the withdrawal was successful.

## 4. Receiving and checking errors

When calling a function that returns an error, we can store that error in a variable.

```go
err := wallet.Withdraw(Bitcoin(100))
```

If the withdrawal fails, `err` contains the returned error.

If it succeeds, `err` is `nil`.

We can check it:

```go
if err != nil {
    fmt.Println("Withdrawal failed:", err)
}
```

This is the standard Go error-handling pattern.

**Important:** An error is just a value. We can return it, store it, pass it to another function, and inspect it.

## 5. Testing errors

Suppose we want to test what happens when someone withdraws more Bitcoin than their wallet contains.

```go
t.Run("withdraw insufficient funds", func(t *testing.T) {
    startingBalance := Bitcoin(20)
    wallet := Wallet{"", startingBalance}

    err := wallet.Withdraw(Bitcoin(100))

    assertError(t, err)
    assertBalance(t, wallet, startingBalance)
})
```

Our helper:

```go
assertError := func(t testing.TB, err error) {
    t.Helper()

    if err == nil {
        t.Error("wanted an error but didn't get one")
    }
}
```

### Execution flow

1. The wallet starts with 20 Bitcoin.
2. We attempt to withdraw 100 Bitcoin.
3. `Withdraw()` checks whether the requested amount exceeds the balance.
4. It creates and returns an error.
5. The error gets stored in `err`.
6. `assertError()` checks that `err` is not nil.
7. `assertBalance()` verifies that the balance remains unchanged.

Notice that the withdrawal never happens because the method returns immediately when it encounters insufficient funds.

**The test checks two things:**
- An error was returned.
- The wallet's state was not modified.

## 6. Sentinel errors: Reusing a specific error

Imagine we have multiple places in our application that need to report insufficient funds.

Instead of creating a new error every time, we can define a reusable error.

This is commonly called a **sentinel error**.

```go
var ErrInsufficientFunds = errors.New("insufficient funds")
```

Now our Withdraw method becomes:

```go
func (w *Wallet) Withdraw(amount Bitcoin) error {
    if amount > w.balance {
        return ErrInsufficientFunds
    }

    w.balance -= amount
    return nil
}
```

Notice the difference:

Previously:

```go
return errors.New("insufficient funds")
```

Now:

```go
return ErrInsufficientFunds
```

We are returning the same predefined error value rather than creating a new one.

### Why is this useful?

Because we can identify the exact kind of error.

```go
err := wallet.Withdraw(Bitcoin(100))

if err == ErrInsufficientFunds {
    fmt.Println("The wallet doesn't have enough funds.")
}
```

Here, `==` works because `Withdraw()` returns the exact sentinel error value.

But this approach has a limitation when errors are wrapped.

## 7. Wrapping errors: Adding context without losing the original error

Imagine our application has multiple layers:

- Wallet handles the actual withdrawal.
- Bank processes the withdrawal for a particular account.
- An HTTP handler might call the Bank.

Suppose the wallet returns:

```text
insufficient funds
```

The Bank receives this error.

If the Bank simply returns it unchanged, the caller only knows that funds were insufficient.

It doesn't know which account or operation caused the problem.

We want to add useful context without losing the original error.

### Example: ProcessWithdrawal

```go
func ProcessWithdrawal(
    wallet *Wallet,
    accountID string,
    amount Bitcoin,
) error {

    if err := wallet.Withdraw(amount); err != nil {
        return fmt.Errorf(
            "processing withdrawal for account %s: %w",
            accountID,
            err,
        )
    }

    return nil
}
```

Notice this part:

```go
fmt.Errorf("processing withdrawal for account %s: %w", accountID, err)
```

We are creating a new error that contains:

1. Additional information about the operation.
2. The original error wrapped inside it.

The `%w` verb is what makes this wrapping possible.

### What does the output look like?

```go
wallet := Wallet{balance: Bitcoin(10)}

err := ProcessWithdrawal(
    &wallet,
    "acc-123",
    Bitcoin(100),
)

fmt.Println(err)
```

Output:

```text
processing withdrawal for account acc-123: insufficient funds
```

We now know:
- Which operation failed: withdrawal.
- Which account was involved: acc-123.
- What actually went wrong: insufficient funds.

**The important distinction:**

`%v` formats an error's message into a new error, but does not preserve the original error relationship.

`%w` wraps the original error, preserving its identity and allowing Go to inspect it later.

## 8. Why errors.Is() is needed

Now consider this situation.

Our original error is:

```go
var ErrInsufficientFunds = errors.New("insufficient funds")
```

The wallet returns it:

```go
return ErrInsufficientFunds
```

But `ProcessWithdrawal()` wraps it:

```go
return fmt.Errorf(
    "processing withdrawal for account %s: %w",
    accountID,
    err,
)
```

The returned error is now a different outer error value.

So this comparison:

```go
err == ErrInsufficientFunds
```

will be false.

Why?

Because `err` is the newly created wrapped error, not the original sentinel error itself.

However, the original error still exists inside the wrapped error.

This is exactly what `errors.Is()` solves.

```go
if errors.Is(err, ErrInsufficientFunds) {
    fmt.Println("Insufficient funds!")
}
```

Even though the outer error is different, `errors.Is()` can find the original error inside it.

### How does errors.Is() work?

Conceptually, it follows the chain of wrapped errors.

Imagine this structure:

```text
Outer error:
"processing withdrawal for account acc-123"
        |
        v
Wrapped original error:
"insufficient funds"
```

`errors.Is()` checks whether the target error exists anywhere in this chain.

It repeatedly unwraps errors until:
- It finds a matching error.
- Or there are no more wrapped errors to inspect.

The standard `fmt.Errorf()` function supports this because `%w` creates an error that implements:

```go
Unwrap() error
```

That method returns the original wrapped error.

**Important:** `errors.Is()` is generally safer than comparing errors with `==` when errors might be wrapped. It also works with unwrapped sentinel errors.

## 9. Understanding errors.Unwrap()

Go provides:

```go
errors.Unwrap(err)
```

It retrieves the next error underneath a wrapped error.

For example:

```go
originalErr := ErrInsufficientFunds

wrappedErr := fmt.Errorf(
    "processing withdrawal: %w",
    originalErr,
)

fmt.Println(wrappedErr)
fmt.Println(errors.Unwrap(wrappedErr))
```

Output:

```text
processing withdrawal: insufficient funds
insufficient funds
```

`errors.Unwrap()` removes one wrapping layer.

You generally don't need to manually unwrap errors when checking them. Use `errors.Is()` instead.

## 10. Understanding errors.As()

Sometimes we don't just want to check whether a particular error exists.

We want to extract an error of a particular type from an error chain.

This is what `errors.As()` does.

For example, imagine we have a custom error type:

```go
type WithdrawalError struct {
    AccountID string
    Reason    string
}

func (e *WithdrawalError) Error() string {
    return e.Reason + " for account " + e.AccountID
}
```

We can create and return it:

```go
return &WithdrawalError{
    AccountID: accountID,
    Reason:    "insufficient funds",
}
```

Now suppose this error gets wrapped by another function.

We can still retrieve the original typed error:

```go
var withdrawalErr *WithdrawalError

if errors.As(err, &withdrawalErr) {
    fmt.Println(withdrawalErr.AccountID)
    fmt.Println(withdrawalErr.Reason)
}
```

`errors.As()` searches the error chain and extracts the first matching error into the target variable.

### Difference between Is and As

| Function | Purpose |
|---|---|
| `errors.Is()` | Checks whether a particular error exists in the error chain. |
| `errors.As()` | Finds an error of a particular type and extracts it. |
| `errors.Unwrap()` | Retrieves the next underlying wrapped error. |

Think of it this way:

- **Is:** Is this the error I'm looking for?
- **As:** Can I extract an error of this type?
- **Unwrap:** What error is underneath this one?

---

## 11. Complete mental model

Consider this entire flow:

```text
Test / HTTP Handler
        |
        v
ProcessWithdrawal()
        |
        v
Wallet.Withdraw()
        |
        v
Insufficient balance detected
        |
        v
Returns ErrInsufficientFunds
        |
        v
ProcessWithdrawal wraps the error using %w
        |
        v
Caller receives the wrapped error
        |
        v
errors.Is(err, ErrInsufficientFunds)
        |
        v
Finds the original error inside the chain
```

## Key takeaways

1. **error interface:** Requires an `Error() string` method.
2. **errors.New():** Creates a simple error value.
3. **Returning errors:** Functions return errors to communicate failure.
4. **nil:** Means no error occurred.
5. **Sentinel errors:** Predefined reusable errors that can be identified.
6. **Wrapping with %w:** Adds context while preserving the original error.
7. **errors.Is():** Checks whether a particular error exists in the error chain.
8. **errors.Unwrap():** Retrieves the next underlying error.
9. **errors.As():** Extracts an error of a particular type from the error chain.

**The central idea:** Errors in Go are values that travel back through your program. Each layer can add useful context using `%w`, while the original cause remains available through `errors.Is()` and `errors.As()`.

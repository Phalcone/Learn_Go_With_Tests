
# Go Methods, and Pointer Receivers

## 1. Structs group related data

A struct combines fields that belong to one concept.

```go
type BankAccount struct {
    name    string
    balance float64
}
```

A value of `BankAccount` contains both `name` and `balance`.

```go
account := BankAccount{
    name:    "John",
    balance: 1000,
}
```

## 2. What is state?

**State** means the current values stored in a struct.

For this account, its state includes the current name and balance. If the balance changes, the account's state changes.

## 3. Why use methods?

If fields are exported, code in another package can access and change them directly:

```go
type BankAccount struct {
    Balance float64
}

account.Balance = -500 // Direct change; no validation
```

That can be a problem. A caller could put the struct into a state that your program should not allow.

Instead, make fields unexported (start their names with a lowercase letter) and provide methods for the operations callers are allowed to perform:

```go
type BankAccount struct {
    balance float64
}

func (a *BankAccount) Deposit(amount float64) {
    if amount > 0 {
        a.balance += amount
    }
}

func (a *BankAccount) Withdraw(amount float64) {
    if amount > 0 && amount <= a.balance {
        a.balance -= amount
    }
}

func (a *BankAccount) Balance() float64 {
    return a.balance
}
```

Callers can use the account like this:

```go
account := BankAccount{balance: 1000}

account.Deposit(500)
account.Withdraw(200)

fmt.Println(account.Balance()) // 1300
```

They can request a deposit or withdrawal and read the balance, but code in another package cannot directly access `account.balance`. They can only use the method call 
not access the variables of the struct. And we can code the methods in such a way that no wrongdoing can occur.

### How this protects the data

- `balance` is unexported, so it cannot be accessed directly from another package.
- `Deposit` and `Withdraw` are exported methods, so callers can use them.
- The methods check the request before changing the balance.
- `Balance` returns the value without giving callers a way to modify the field.

This is **encapsulation**: keeping data and the operations that manage it together, while controlling access to the data.

**Important Go detail:** Unexported fields are still accessible to other code in the same package. Go's visibility rules work at the package level, not the individual-file or individual-struct level.

## 4. Pointer receivers: why `*BankAccount`?

Look at this method:

```go
func (a *BankAccount) Deposit(amount float64) {
    a.balance += amount
}
```

The part `(a *BankAccount)` is the **receiver**. It says that `a` is a pointer to a `BankAccount`.

A pointer holds the location of a value in memory. It lets the method reach and modify the original account rather than working on a separate copy.

### Without a pointer receiver

```go
func (a BankAccount) Deposit(amount float64) {
    a.balance += amount
}
```

Here, `a BankAccount` receives a **copy** of the account. The method changes that copy. When the method finishes, the original account's balance remains unchanged.

For example, if the original balance is `1000`, calling this version of `Deposit(500)` changes the copy to `1500`, but the original stays at `1000`.

### With a pointer receiver

```go
func (a *BankAccount) Deposit(amount float64) {
    a.balance += amount
}
```

Here, `a` points to the original account. The method changes the original balance from `1000` to `1500`.

In short:

| Receiver | What the method works with | Does a field change affect the original? |
|---|---|---|
| `a BankAccount` | A copy of the struct | No |
| `a *BankAccount` | A pointer to the original struct | Yes |

Use a pointer receiver when a method needs to modify the struct's state. It can also avoid copying a large struct.

Go often lets you call a pointer-receiver method on an addressable variable without writing `&` yourself:

```go
account := BankAccount{balance: 1000}
account.Deposit(500) // Go takes the address of account for the call
```

The equivalent explicit call is:

```go
(&account).Deposit(500)
```

## 5. Reading versus changing state

The receiver also helps communicate what a method does:

```go
func (a *BankAccount) Deposit(amount float64) {
    // Changes the account
}

func (a *BankAccount) Balance() float64 {
    // Reads and returns the balance
}
```

A method that only needs to read the struct can usually use a value receiver or a pointer receiver. Go does not have a special `const` receiver.

The `Balance` method above uses a pointer receiver for consistency with the other account methods; it does not modify the balance.

## Key points

1. A **struct** groups related fields.
2. The struct's **state** is the current value of those fields.
3. **Methods** define the operations callers can perform on a struct.
4. **Unexported fields** (lowercase names) cannot be accessed directly from another package.
5. **Pointer receivers** (`*Type`) let methods modify the original struct instead of a copy.
6. Methods can validate requests before changing state, helping keep data valid.

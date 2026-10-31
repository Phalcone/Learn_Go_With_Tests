# Go Dictionary --- Notes

These notes follow the same order as the pasted lesson. The focus is on
what each piece of code does, why it is introduced, how the tests guide
the implementation, and how the four dictionary operations work.

## 1. Start with a map and test Search

A Go map stores key--value pairs. Here, each word is a key and its
definition is the value:

``` go
dictionary := map[string]string{
    "test": "this is just a test",
}
```

The type `map[string]string` means: - `map` declares a map. - The type
inside `[]` (`string`) is the key type. - The type after `[]` (`string`)
is the value type.

Map keys must be **comparable**. Go must be able to compare two keys to
determine whether they are equal and therefore whether a requested key
matches an entry. Strings, numbers and booleans are comparable; slices,
maps and functions are not. The value type can be any type, including
another map.

The first test checks that searching for a known word returns its
definition:

``` go
func TestSearch(t *testing.T) {
    dictionary := map[string]string{"test": "this is just a test"}

    got := Search(dictionary, "test")
    want := "this is just a test"

    if got != want {
        t.Errorf("got %q want %q given, %q", got, want, "test")
    }
}
```

`got` is the actual result and `want` is the expected result. `t.Errorf`
reports a test failure. `%q` prints a quoted string, making the actual
and expected strings easier to distinguish.

### Follow the test-first process

The lesson uses a small test-driven development loop:

1.  Write a test for the behaviour you want.
2.  Run `go test` and see what fails.
3.  Write the minimum code needed to make the test compile.
4.  Run the test again to see the behaviour failure.
5.  Implement enough code to make it pass.
6.  Refactor if the code or test can be made clearer.

At first, `go test` fails because `Search` has not been declared. Add
the smallest implementation:

``` go
func Search(dictionary map[string]string, word string) string {
    return ""
}
```

Now the code compiles, but the test fails because the function returns
an empty string. That failure tells us the missing behaviour: the
function must retrieve the definition.

``` go
func Search(dictionary map[string]string, word string) string {
    return dictionary[word]
}
```

A map lookup uses `dictionary[key]`, just as an array or slice is
indexed, except that a map uses a key rather than a numeric index.

## 2. Refactor the test with an assertion helper

The same string comparison will be needed in later tests, so it is moved
into a helper:

``` go
func assertStrings(t testing.TB, got, want string) {
    t.Helper()

    if got != want {
        t.Errorf("got %q want %q", got, want)
    }
}
```

-   `testing.TB` lets the helper work with both tests and benchmarks.
-   `t.Helper()` tells Go that this function is a test helper. If the
    assertion fails, the reported location points to the test that
    called it.
-   `t.Errorf` reports the mismatch without immediately stopping the
    test.

The test becomes:

``` go
func TestSearch(t *testing.T) {
    dictionary := map[string]string{"test": "this is just a test"}

    got := Search(dictionary, "test")
    want := "this is just a test"

    assertStrings(t, got, want)
}
```

## 3. Create a Dictionary type and make Search a method

Passing a map into every function works, but the code is specifically
modelling a dictionary. Give the map a meaningful type:

``` go
type Dictionary map[string]string
```

This is a named type whose underlying data structure is still a map from
strings to strings. It lets us attach methods to the dictionary.

The test now uses the new type and calls `Search` as a method:

``` go
func TestSearch(t *testing.T) {
    dictionary := Dictionary{"test": "this is just a test"}

    got := dictionary.Search("test")
    want := "this is just a test"

    assertStrings(t, got, want)
}
```

The method is:

``` go
func (d Dictionary) Search(word string) string {
    return d[word]
}
```

`(d Dictionary)` is the **receiver**. It makes `Search` a method
belonging to `Dictionary`, so it is called with
`dictionary.Search("test")`. The receiver is named `d` inside the
method.

## 4. Make Search report when a word is missing

The current `Search` returns an empty string when a word is absent. But
an empty string could also be a real definition. The caller cannot tell
whether the word is missing or its definition is empty.

Go map lookup can return two values:

``` go
definition, ok := d[word]
```

-   `definition` is the value associated with the key.
-   `ok` is `true` if the key exists and `false` if it does not.

Use this to make `Search` return both a definition and an error:

``` go
func (d Dictionary) Search(word string) (string, error) {
    definition, ok := d[word]
    if !ok {
        return "", errors.New("could not find the word you were looking for")
    }

    return definition, nil
}
```

Go functions commonly return an `error` as an additional result: - On
success, return the useful result and `nil` for the error. - On failure,
return an error describing what went wrong. The other result is usually
its zero value.

The caller should check the error before using the result:

``` go
definition, err := dictionary.Search("test")
if err != nil {
    // Handle the error.
    return
}

// It is now safe to use definition.
```

### Test known and unknown words separately

``` go
func TestSearch(t *testing.T) {
    dictionary := Dictionary{"test": "this is just a test"}

    t.Run("known word", func(t *testing.T) {
        got, _ := dictionary.Search("test")
        want := "this is just a test"

        assertStrings(t, got, want)
    })

    t.Run("unknown word", func(t *testing.T) {
        _, err := dictionary.Search("unknown")
        want := "could not find the word you were looking for"

        if err == nil {
            t.Fatal("expected to get an error.")
        }

        assertStrings(t, err.Error(), want)
    })
}
```

`t.Run` creates named subtests. This makes the known-word and
unknown-word behaviours independently visible in test output.

In the known-word test, `_` ignores the error because the test is
focused on the returned definition. In the unknown-word test, `_`
ignores the definition because the test is focused on the error.

The test first checks `err == nil`. Only if an error exists does it call
`err.Error()` to get its message. `t.Fatal` stops that subtest
immediately if the expected error was not returned.

## 5. Give the missing-word error a reusable name

Instead of creating an error inside `Search` each time, define a named
error:

``` go
var ErrNotFound = errors.New(
    "could not find the word you were looking for",
)
```

Then return it:

``` go
func (d Dictionary) Search(word string) (string, error) {
    definition, ok := d[word]
    if !ok {
        return "", ErrNotFound
    }

    return definition, nil
}
```

The test can now compare against `ErrNotFound`, rather than repeating
the error message:

``` go
t.Run("unknown word", func(t *testing.T) {
    _, got := dictionary.Search("unknown")

    if got == nil {
        t.Fatal("expected to get an error.")
    }

    assertError(t, got, ErrNotFound)
})
```

The helper is:

``` go
func assertError(t testing.TB, got, want error) {
    t.Helper()

    if !errors.Is(got, want) {
        t.Errorf("got error %q want %q", got, want)
    }
}
```

`errors.Is` checks whether the error matches the expected error. It is
preferable to comparing error-message strings: the test depends on the
named error, not on the exact wording. It also continues to work if an
error is later wrapped.

## 6. Add a word to the dictionary

The next requirement is to add a word. Write the test first:

``` go
func TestAdd(t *testing.T) {
    dictionary := Dictionary{}
    dictionary.Add("test", "this is just a test")

    want := "this is just a test"
    got, err := dictionary.Search("test")
    if err != nil {
        t.Fatal("should find added word:", err)
    }

    assertStrings(t, got, want)
}
```

The test uses `Search` to check that the added word can be found. This
reuses behaviour that has already been tested.

First add a method that compiles but does nothing:

``` go
func (d Dictionary) Add(word, definition string) {
}
```

The test now runs but fails because the word was not added. Implement
the map assignment:

``` go
func (d Dictionary) Add(word, definition string) {
    d[word] = definition
}
```

`d[word] = definition` inserts the key and value into the map.

### Important: map assignment overwrites existing values

If the key already exists, map assignment replaces its old value. A map
does not automatically return an error for duplicate keys.

That is not the behaviour wanted from this `Add` method: `Add` should
add a new word, not silently replace an existing definition. So the
tests are expanded to cover both cases:

``` go
func TestAdd(t *testing.T) {
    t.Run("new word", func(t *testing.T) {
        dictionary := Dictionary{}
        word := "test"
        definition := "this is just a test"

        err := dictionary.Add(word, definition)

        assertError(t, err, nil)
        assertDefinition(t, dictionary, word, definition)
    })

    t.Run("existing word", func(t *testing.T) {
        word := "test"
        definition := "this is just a test"
        dictionary := Dictionary{word: definition}

        err := dictionary.Add(word, "new test")

        assertError(t, err, ErrWordExists)
        assertDefinition(t, dictionary, word, definition)
    })
}
```

The existing-word test checks both that an error is returned and that
the original definition remains unchanged.

Because the tests now use the return value from `Add`, its signature
must return an error. Define the new error and update the method
signature:

``` go
var (
    ErrNotFound   = errors.New("could not find the word you were looking for")
    ErrWordExists = errors.New("cannot add word because it already exists")
)

func (d Dictionary) Add(word, definition string) error {
    d[word] = definition
    return nil
}
```

This is only the minimum code needed to compile. The existing-word test
still fails because the implementation overwrites the old value and
returns `nil`.

### Implement the duplicate check using Search

``` go
func (d Dictionary) Add(word, definition string) error {
    _, err := d.Search(word)

    switch err {
    case ErrNotFound:
        d[word] = definition
    case nil:
        return ErrWordExists
    default:
        return err
    }

    return nil
}
```

The switch handles the possible outcomes of `Search`:

-   `ErrNotFound`: the word is not in the dictionary, so add it.
-   `nil`: the word already exists, so return `ErrWordExists` and leave
    the original definition unchanged.
-   `default`: an unexpected error occurred, so return it instead of
    treating it as a normal missing-word case.

After a successful insertion, return `nil` to show that `Add` succeeded.

### Simplify the test with assertDefinition

The repeated lookup and comparison are moved into a helper:

``` go
func assertDefinition(
    t testing.TB,
    dictionary Dictionary,
    word, definition string,
) {
    t.Helper()

    got, err := dictionary.Search(word)
    if err != nil {
        t.Fatal("should find added word:", err)
    }

    assertStrings(t, got, definition)
}
```

This helper searches for the word, stops the test if it cannot be found,
and then checks that its definition is correct.

## 6.1 The central idea: errors become part of the normal design

One of the strongest ideas in this lesson is how error handling is
introduced gradually and then reused throughout the dictionary. Errors
are not added as an afterthought. The tests first reveal where the
program needs to communicate failure, and the code then returns a
specific error that the tests and other methods can understand.

Think of an error as an extra piece of information returned by an
operation:

``` go
definition, err := dictionary.Search("test")
```

The two return values answer two different questions:

-   `definition`: what did the search find?
-   `err`: did the operation succeed, and if not, what went wrong?

This gives the caller an explicit way to distinguish success from
failure. A successful search returns the definition and `nil`. A failed
search returns an empty definition and `ErrNotFound`.

### Step 1: Let the test reveal that a result alone is not enough

The original `Search` returned only a string:

``` go
func (d Dictionary) Search(word string) string {
    return d[word]
}
```

That is enough to find a known word, but it cannot distinguish a missing
word from a word whose definition is an empty string. The test for an
unknown word exposes this gap.

The method is therefore changed to return `(string, error)`. The test
checks that an error is present and then checks its message. This is the
first important connection: **the test defines what failure should look
like, and the function's return signature is changed to make that
failure observable.**

### Step 2: Turn the error into a named condition

Initially, `Search` creates an error directly:

``` go
return "", errors.New("could not find the word you were looking for")
```

The lesson then extracts that error into `ErrNotFound`:

``` go
var ErrNotFound = errors.New(
    "could not find the word you were looking for",
)
```

Now the function returns a known error value:

``` go
if !ok {
    return "", ErrNotFound
}
```

And the test checks for that same named condition:

``` go
assertError(t, got, ErrNotFound)
```

This is cleaner than making the test compare message text. The function
and test share the meaning `ErrNotFound`; the message is simply how that
condition is displayed to a person.

### Step 3: Use the same error to make other methods intelligent

The particularly useful part is that `Add`, `Update`, and `Delete` call
`Search` rather than each implementing their own map-existence check.

For example, `Add` does this:

``` go
_, err := d.Search(word)

switch err {
case ErrNotFound:
    d[word] = definition
case nil:
    return ErrWordExists
default:
    return err
}
```

The result from `Search` tells `Add` what action is appropriate:

-   `ErrNotFound` means the word is absent, so adding it is allowed.
-   `nil` means the word exists, so adding it would overwrite something;
    return `ErrWordExists`.
-   Any other error is passed back to the caller.

The same lookup result drives the other methods, but each method
interprets it according to its own purpose:

``` go
// Update
switch err {
case ErrNotFound:
    return ErrWordDoesNotExist
case nil:
    d[word] = definition
default:
    return err
}
```

``` go
// Delete
switch err {
case ErrNotFound:
    return ErrWordDoesNotExist
case nil:
    delete(d, word)
default:
    return err
}
```

This is a neat separation of responsibilities: - `Search` knows how to
determine whether a word exists. - `Add`, `Update`, and `Delete` decide
what that result means for their own operation.

The existence check is written once, and the methods build their own
rules on top of it.

### Step 4: Give each operation the error that describes its failure

Notice that a missing word does not always mean the same thing to the
caller:

  Operation   What went wrong?                        Error returned
  ----------- --------------------------------------- -----------------------
  Search      The requested word could not be found   `ErrNotFound`
  Add         The word is already present             `ErrWordExists`
  Update      There is no existing word to update     `ErrWordDoesNotExist`
  Delete      There is no existing word to delete     `ErrWordDoesNotExist`

This is why `Update` and `Delete` return `ErrWordDoesNotExist` rather
than simply returning `ErrNotFound`. The error describes the failed
operation, not just the lower-level lookup result.

That distinction is useful to callers. A caller can decide what to do
based on the operation-specific failure, instead of having to infer what
happened from a generic message.

### Step 5: The tests verify both the error and the state

A good test does not stop at checking that an error was returned. It
also checks that the dictionary was left in the correct state.

For example, the existing-word `Add` test checks:

``` go
err := dictionary.Add(word, "new test")

assertError(t, err, ErrWordExists)
assertDefinition(t, dictionary, word, definition)
```

These two assertions prove two separate parts of the behaviour:

1.  The method correctly reports that the word already exists.
2.  The original definition has not been overwritten.

Likewise, the `Update` test checks that a missing word produces
`ErrWordDoesNotExist`, while the `Delete` tests check both that a
successful deletion makes the word unsearchable and that deleting a
missing word returns the expected error.

The tests and implementation therefore reinforce each other. The tests
state the contract; the methods implement it; and the returned errors
make that contract observable.

### Step 6: Why the `default` branch matters

The methods use a switch with a `default` branch:

``` go
default:
    return err
```

This is a small but important safety measure. The method expects either
`ErrNotFound` or `nil` from `Search`. If some other error is returned,
it should not silently treat that error as one of the expected cases.

Returning it preserves the original failure for the caller. It also
means the method does not accidentally modify the dictionary when the
lookup did not complete in an expected way.

### Step 7: Make the shared errors stable with `DictionaryErr` constants

The final refactor turns the named errors into constants:

``` go
type DictionaryErr string

func (e DictionaryErr) Error() string {
    return string(e)
}

const (
    ErrNotFound         = DictionaryErr("could not find the word you were looking for")
    ErrWordExists       = DictionaryErr("cannot add word because it already exists")
    ErrWordDoesNotExist = DictionaryErr("cannot perform operation on word because it does not exist")
)
```

This keeps the same error-handling design, but makes the error values
immutable. The custom type implements the `error` interface through its
`Error() string` method, so the constants can still be returned as
errors and passed to `errors.Is`.

The key point is that the refactor does not change what the methods or
tests mean. `Search` still returns `ErrNotFound`; `Add` still returns
`ErrWordExists` for a duplicate; and `Update` and `Delete` still return
`ErrWordDoesNotExist` when their target is absent. The representation
becomes more stable, while the behaviour and tests remain the same.

**The overall pattern:** a method returns an error to make failure
explicit; the caller or another method interprets that named error;
tests assert the exact failure condition; and shared lookup logic lets
all operations use the same source of truth. That is how error handling
becomes integrated into the design rather than scattered across the
code.

## 8. Why maps do not need pointer receivers

The methods use value receivers, for example:

``` go
func (d Dictionary) Add(word, definition string) error
```

Even though `Add` changes the dictionary, it does not need a pointer
receiver. The reason is how maps are represented and passed in Go.

A map value is a small value that refers to the underlying map data
structure. When a map is passed to a method, Go copies that small map
value, but the copy still refers to the **same underlying map data**.

Conceptually:

``` text
dictionary ──┐
             ├──> same underlying map
method's d ──┘
```

Therefore, this changes the shared map:

``` go
d[word] = definition
```

The caller sees the change. A pointer receiver is not required just to
insert, update, or delete map entries.

The important distinction is: Go copies the map value, not all the
key--value entries into a separate map. Both the caller and the method
still refer to the same underlying map.

### Nil maps versus empty initialized maps

A map can be `nil`:

``` go
var dictionary map[string]string
```

A nil map is safe to read. A lookup behaves like a lookup in an empty
map and returns the value type's zero value:

``` go
value := dictionary["missing"] // Safe; value is ""
```

But writing to a nil map causes a runtime panic:

``` go
dictionary["test"] = "definition" // Panic
```

So, before writing to a map, initialize it. Either of these creates an
empty, usable map:

``` go
dictionary := map[string]string{}

// Or:
dictionary := make(map[string]string)
```

For the custom type, this is also initialized:

``` go
dictionary := Dictionary{}
```

This is why the tests use `Dictionary{}` when they need an empty
dictionary. It is not nil; it is an empty map that can accept new
entries.

## 9. Update an existing definition

`Update` should change the definition of a word that is already present.

Start with a test:

``` go
func TestUpdate(t *testing.T) {
    word := "test"
    definition := "this is just a test"
    dictionary := Dictionary{word: definition}
    newDefinition := "new definition"

    dictionary.Update(word, newDefinition)

    assertDefinition(t, dictionary, word, newDefinition)
}
```

Add the minimum method needed to compile:

``` go
func (d Dictionary) Update(word, definition string) {}
```

The test fails because the old definition is still there. As with `Add`,
map assignment performs the change:

``` go
func (d Dictionary) Update(word, definition string) {
    d[word] = definition
}
```

But this also inserts a word if it does not exist. That is not what
`Update` should do. Add a test for a missing word and make `Update`
return an error:

``` go
t.Run("existing word", func(t *testing.T) {
    word := "test"
    definition := "this is just a test"
    dictionary := Dictionary{word: definition}
    newDefinition := "new definition"

    err := dictionary.Update(word, newDefinition)

    assertError(t, err, nil)
    assertDefinition(t, dictionary, word, newDefinition)
})

t.Run("new word", func(t *testing.T) {
    word := "test"
    definition := "this is just a test"
    dictionary := Dictionary{}

    err := dictionary.Update(word, definition)

    assertError(t, err, ErrWordDoesNotExist)
})
```

Define the more specific error and change the signature:

``` go
var (
    ErrNotFound         = errors.New("could not find the word you were looking for")
    ErrWordExists       = errors.New("cannot add word because it already exists")
    ErrWordDoesNotExist = errors.New("cannot perform operation on word because it does not exist")
)

func (d Dictionary) Update(word, definition string) error {
    d[word] = definition
    return nil
}
```

This compiles, but the missing-word test fails because the method still
inserts the word. The correct implementation checks with `Search` first:

``` go
func (d Dictionary) Update(word, definition string) error {
    _, err := d.Search(word)

    switch err {
    case ErrNotFound:
        return ErrWordDoesNotExist
    case nil:
        d[word] = definition
    default:
        return err
    }

    return nil
}
```

-   If `Search` returns `ErrNotFound`, there is no existing entry to
    update, so return `ErrWordDoesNotExist`.
-   If `Search` returns `nil`, the word exists, so replace its
    definition.
-   If another error occurs, pass it back to the caller.

This is similar to `Add`, but the actions for the two normal cases are
reversed.

### Why use ErrWordDoesNotExist instead of ErrNotFound?

The errors describe different situations:

-   `ErrNotFound` is used when a search cannot find a word.
-   `ErrWordDoesNotExist` is used when an operation such as update or
    delete cannot be performed because its target does not exist.

A specific error gives the caller more information about what failed.
The caller can respond differently to a failed search and a failed
update.

## 10. Delete a word

`Delete` should remove an existing word. Start with a test:

``` go
func TestDelete(t *testing.T) {
    word := "test"
    dictionary := Dictionary{word: "test definition"}

    dictionary.Delete(word)

    _, err := dictionary.Search(word)
    assertError(t, err, ErrNotFound)
}
```

The test deletes the word and then uses `Search` to verify that it is
gone.

First add a method stub:

``` go
func (d Dictionary) Delete(word string) {
}
```

The test fails because the word is still present. Go provides a built-in
`delete` function for maps:

``` go
func (d Dictionary) Delete(word string) {
    delete(d, word)
}
```

`delete` takes the map and the key to remove. It does not return a
value. If the key is absent, it does nothing.

As with `Update`, we also want to report when the word does not exist.
The tests now cover both cases:

``` go
func TestDelete(t *testing.T) {
    t.Run("existing word", func(t *testing.T) {
        word := "test"
        dictionary := Dictionary{word: "test definition"}

        err := dictionary.Delete(word)
        assertError(t, err, nil)

        _, err = dictionary.Search(word)
        assertError(t, err, ErrNotFound)
    })

    t.Run("non-existing word", func(t *testing.T) {
        word := "test"
        dictionary := Dictionary{}

        err := dictionary.Delete(word)
        assertError(t, err, ErrWordDoesNotExist)
    })
}
```

The method must return an error, so its final implementation is:

``` go
func (d Dictionary) Delete(word string) error {
    _, err := d.Search(word)

    switch err {
    case ErrNotFound:
        return ErrWordDoesNotExist
    case nil:
        delete(d, word)
    default:
        return err
    }

    return nil
}
```

-   If the word is missing, return `ErrWordDoesNotExist`.
-   If it exists, delete it.
-   If another error occurs, return that error.
-   If deletion succeeds, return `nil`.

The test checks both the returned error and the state of the dictionary
after deletion.

## 11. Refactor the errors into constants

The lesson then changes the error declarations from `var` values created
by `errors.New` to constants:

``` go
const (
    ErrNotFound         = DictionaryErr("could not find the word you were looking for")
    ErrWordExists       = DictionaryErr("cannot add word because it already exists")
    ErrWordDoesNotExist = DictionaryErr("cannot perform operation on word because it does not exist")
)

type DictionaryErr string

func (e DictionaryErr) Error() string {
    return string(e)
}
```

### How can a constant be an error?

Go's `error` is an interface. A type implements that interface if it has
an `Error` method that returns a string:

``` go
Error() string
```

`DictionaryErr` is a custom type whose underlying type is `string`. Its
method:

``` go
func (e DictionaryErr) Error() string {
    return string(e)
}
```

converts the `DictionaryErr` value into a regular string and returns it
as the error message. Because it has this method, a `DictionaryErr`
value can be used wherever an `error` is expected.

### Why can these errors be constants?

A constant must be created from a constant expression. `errors.New(...)`
is a function call, so its result cannot be declared with `const`.

A string conversion to the custom type is allowed in a constant
declaration:

``` go
ErrNotFound = DictionaryErr("could not find the word you were looking for")
```

The custom type also supplies the `Error()` method, so the constant
still behaves as an error.

Constants are immutable: they cannot be reassigned. This is useful for
shared error values because code cannot accidentally replace a standard
error with a different value. The named errors are stable conditions
that functions can return and tests can check.

The tests continue to use `errors.Is`:

``` go
assertError(t, got, ErrNotFound)
```

This keeps the test focused on which error occurred, rather than on the
exact text of its message.

## 12. How Search is used by all the operations

`Search` is used both by callers who want a definition and internally by
`Add`, `Update`, and `Delete` to check whether a word exists.

  -----------------------------------------------------------------------
  Method                  If Search finds the     If Search returns
                          word (`err == nil`)     `ErrNotFound`
  ----------------------- ----------------------- -----------------------
  `Search`                Return the definition   Return `ErrNotFound`

  `Add`                   Return `ErrWordExists`  Add the new word

  `Update`                Replace the definition  Return
                                                  `ErrWordDoesNotExist`

  `Delete`                Delete the word         Return
                                                  `ErrWordDoesNotExist`
  -----------------------------------------------------------------------

This reuse keeps the existence check in one place. The three modifying
methods all use a `switch` so they can handle the expected outcomes
explicitly. Their `default` case returns any unexpected error rather
than silently treating it as a missing word or a successful lookup.

## 13. Test helpers and test behaviour

The lesson uses three helpers to keep tests focused:

``` go
func assertStrings(t testing.TB, got, want string) {
    t.Helper()

    if got != want {
        t.Errorf("got %q want %q", got, want)
    }
}
```

Checks that two strings match.

``` go
func assertError(t testing.TB, got, want error) {
    t.Helper()

    if !errors.Is(got, want) {
        t.Errorf("got error %q want %q", got, want)
    }
}
```

Checks that the returned error matches the expected error. Passing `nil`
as `want` checks that no error occurred.

``` go
func assertDefinition(
    t testing.TB,
    dictionary Dictionary,
    word, definition string,
) {
    t.Helper()

    got, err := dictionary.Search(word)
    if err != nil {
        t.Fatal("should find added word:", err)
    }

    assertStrings(t, got, definition)
}
```

Checks that a word can be found and that its definition is correct.

The helpers use `t.Helper()` so a failure points to the test that called
the helper. They use `t.Errorf` for a comparison failure so the test can
continue, and `t.Fatal` when a failed lookup means the remaining
definition check cannot be performed meaningfully.

The tests use `t.Run` to separate behaviours such as `"new word"` and
`"existing word"`. For each operation, they check not just the returned
error but also the resulting dictionary contents where relevant.

## 14. Run the tests

Run all tests in the package directory:

``` bash
go test
```

For named subtests and more detailed output:

``` bash
go test -v
```

The tests cover these behaviours:

-   **Search:** finds a known word and reports `ErrNotFound` for an
    unknown word.
-   **Add:** inserts a new word, but refuses to overwrite an existing
    word.
-   **Update:** changes an existing definition, but does not create a
    missing word.
-   **Delete:** removes an existing word and reports when the word was
    not present.

The overall pattern is: describe the behaviour in a test, use the
failure to guide the smallest implementation, then refactor without
changing the tested behaviour.

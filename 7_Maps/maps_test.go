package main

import (
	"errors"
	"testing"
)

func TestSearch(t *testing.T) {
	dictionary := Dictionary{"test": "this is just a test"}
	t.Run("Known word", func(t *testing.T) {
		got, err := dictionary.Search("test")
		want := "this is just a test"

		if err != nil {
			t.Fatal("Didn't expect an error.")
		}

		assertStrings(t, got, want)
	})
	t.Run("Unknown word", func(t *testing.T) {
		_, err := dictionary.Search("Unknown")

		if err == nil {
			t.Fatal("expected to get an error.")
		}

		assertError(t, err, ErrWordNotFound)
	})
}

func assertStrings(t testing.TB, got, want string) {
	t.Helper()

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func assertError(t testing.TB, got, want error) {
	t.Helper()

	if !errors.Is(got, want) {
		t.Errorf("Got error: %q; Want: %q", got, want)
	}
}

func assertDefinition(t testing.TB, dictionary Dictionary, word, definition string) {
	t.Helper()

	got, err := dictionary.Search(word)
	if err != nil {
		t.Fatal("should find added word:", err)
	}
	assertStrings(t, got, definition)
}

func TestAdd(t *testing.T) {
	t.Run("New word", func(t *testing.T) {
		dictionary := Dictionary{}
		word := "test"
		definition := "this is just a test."

		dictionary.Add(word, definition)

		assertDefinition(t, dictionary, word, definition)
	})
	t.Run("Existing word", func(t *testing.T) {
		word := "test"
		definition := "This is just a test."

		dictionary := Dictionary{word: definition}
		err := dictionary.Add(word, "New test")

		assertError(t, err, ErrWordExists)
		assertDefinition(t, dictionary, word, definition)
	})
}

func TestUpdate(t *testing.T) {
	t.Run("Updating present word", func(t *testing.T) {
		word := "test"
		definition := "Initial test."
		dictionary := Dictionary{word: definition}

		newDefinition := "Update test"
		dictionary.Update(word, newDefinition)

		assertDefinition(t, dictionary, word, newDefinition)
	})
	t.Run("Word is not present and is being added due to update", func(t *testing.T) {
		word := "test"
		definition := "Initial test."
		dictionary := Dictionary{}

		err := dictionary.Update(word, definition)

		assertError(t, err, ErrUpdateMissingWord)
	})
}

func TestDelete(t *testing.T) {
	t.Run("Word you're trying to delete exists.", func(t *testing.T) {
		word := "test"
		definition := "This is just a test."
		dictionary := Dictionary{word: definition}
		dictionary.Delete(word)

		_, err := dictionary.Search(word)

		assertError(t, err, ErrWordNotFound)
	})
	t.Run("Word you're trying to delete exists.", func(t *testing.T) {
		word := "test"
		definition := "This is just a test."
		dictionary := Dictionary{word: definition}
		dictionary.Delete(word)

		_, err := dictionary.Search(word)

		assertError(t, err, ErrWordNotFound)
	})
}

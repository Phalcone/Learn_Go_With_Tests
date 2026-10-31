package main

const (
	ErrWordNotFound       = DictionaryErr("could not find the word you were looking for.")
	ErrWordExists         = DictionaryErr("cannot add word because it already exists.")
	ErrUpdateMissingWord  = DictionaryErr("The word you're trying to update doesn't exist in the dictionary.")
	ErrDeleteNonExistWord = DictionaryErr("The word your'e trying to delete is not present in the dictionary.")
)

type DictionaryErr string

func (e DictionaryErr) Error() string {
	return string(e)
}

type Dictionary map[string]string

func (d Dictionary) Search(word string) (string, error) {
	definition, existence := d[word]

	if !existence {
		return "", ErrWordNotFound
	}
	return definition, nil
}

func (d Dictionary) Add(word, definition string) error {
	_, existence := d[word]
	if !existence {
		d[word] = definition
		return nil
	}
	return ErrWordExists
}

func (d Dictionary) Update(word, definition string) error {
	_, err := d.Search(word)

	switch err {
	case ErrWordNotFound:
		return ErrUpdateMissingWord
	case nil:
		d[word] = definition
		return nil
	default:
		return err
	}
}

func (d Dictionary) Delete(word string) error {
	_, err := d.Search(word)

	switch err {
	case nil:
		delete(d, word)
		return nil
	case ErrWordNotFound:
		return ErrDeleteNonExistWord
	default:
		return err
	}
}

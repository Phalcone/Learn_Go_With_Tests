package main

import "fmt"

const (
	spanish = "Spanish"
	french  = "French"

	englishHelloPrefix  = "Hello, "
	spanishhHelloPrefix = "Hola, "
	frenchHelloPrefix   = "Bonjour, "
)

func Hello(name, language string) string {
	switch name {
	case "":
		name = "World"
		return greetingPrefix(language) + name + "!"
	default:
		return greetingPrefix(language) + name + "!"
	}
}

func greetingPrefix(language string) (prefix string) {
	switch language {
	case spanish:
		prefix = spanishhHelloPrefix
	case french:
		prefix = frenchHelloPrefix
	default:
		prefix = englishHelloPrefix
	}
	return
}

func main() {
	fmt.Println(Hello("Akshaj", ""))
}

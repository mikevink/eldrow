package main

import (
	"bufio"
	"flag"
	"log"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func main() {
	var dict string
	var length int
	var skip string
	var include string
	var match string
	var all bool
	flag.StringVar(&dict, "d", "/usr/share/myspell/en_US.dic", "dictionary file to use")
	flag.IntVar(&length, "l", 5, "word length to use")
	flag.StringVar(&skip, "s", "-", "which letters to skip")
	flag.StringVar(&include, "i", "-", "which letters to include")
	flag.StringVar(&match, "m", "-", "what regular expression to match")
	flag.BoolVar(&all, "all", false, "print all the words")

	flag.Parse()

	file, err := os.Open(dict)

	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}

	matcher, _ := regexp.Compile(`(?i)^` + match + `$`)
	alpha, _ := regexp.Compile(`(?i)^[a-z]{` + strconv.Itoa(length) + `}$`)

	var includes []rune
	if "-" != include {
		includes = []rune(include)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var words []string

	for scanner.Scan() {
		word := scanner.Text()
		if strings.Contains(word, "/") {
			word = strings.Split(word, "/")[0]
		}
		ok := true
		if "-" != match {
			ok = matcher.MatchString(word)
		} else {
			ok = alpha.MatchString(word)
		}

		if ok && "-" != include {
			for _, r := range includes {
				ok = ok && strings.ContainsRune(word, r)
			}
		}

		if ok && "-" != skip {
			ok = !strings.ContainsAny(word, skip)
		}

		if ok {
			words = append(words, word)
		}
	}

	_ = file.Close()

	if "-" == skip && "-" == include && "-" == match {
		if all {
			log.Printf("All words of length %d:", length)
			for _, word := range words {
				log.Printf("\t%s", word)
			}
		} else {
			rng := rand.New(rand.NewSource(time.Now().UnixNano()))
			inx := rng.Intn(len(words))
			log.Printf("Returning a random %d letter word: %s", length, words[inx])
		}
	} else {
		log.Printf("Words that match the restrictions [length: %d, include: '%s', skip: '%s', match: '%s']:", length, include, skip, match)
		for _, word := range words {
			log.Printf("\t%s", word)
		}
	}
}
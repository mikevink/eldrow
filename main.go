package main

import (
	"bufio"
	"flag"
	"log"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"
)

func configDir() string {
	dir, err := os.UserConfigDir()
	if nil != err {
		log.Fatal("Could not get user config dir", err)
	}
	return dir + "/eldrow.yaml"
}

var DefaultDictionary = "/usr/share/myspell/en_US.dic"

type Args struct {
	Dict         string
	Length       string
	RegexpLength string
	Skip         string
	Include      string
	Match        string
	All          bool
}

func parseArgs() Args {
	args := Args{
		Dict:         DefaultDictionary,
		Length:       "*",
		RegexpLength: "*",
		Skip:         "",
		Include:      "",
		Match:        "",
		All:          false,
	}

	flag.StringVar(&args.Dict, "d", args.Dict, "dictionary file to use")
	flag.StringVar(&args.Length, "l", args.Length, "word length to use")
	flag.StringVar(&args.Skip, "s", args.Skip, "which letters to skip")
	flag.StringVar(&args.Include, "i", args.Include, "which letters to include")
	flag.StringVar(&args.Match, "m", args.Match, "what regular expression to match")
	flag.BoolVar(&args.All, "all", args.All, "print all the words")

	flag.Parse()

	if "*" != args.Length {
		args.RegexpLength = "{" + args.Length + "}"
	}

	return args
}

func main() {
	args := parseArgs()

	file, err := os.Open(args.Dict)

	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}

	matcher, _ := regexp.Compile(`(?i)^` + args.Match + `$`)
	alpha, _ := regexp.Compile(`(?i)^[a-z]` + args.RegexpLength + `$`)

	var includes []rune
	if "" != args.Include {
		includes = []rune(args.Include)
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
		if "" != args.Match {
			ok = matcher.MatchString(word)
		} else {
			ok = alpha.MatchString(word)
		}

		if ok && "" != args.Include {
			for _, r := range includes {
				ok = ok && strings.ContainsRune(word, r)
			}
		}

		if ok && "" != args.Skip {
			ok = !strings.ContainsAny(word, args.Skip)
		}

		if ok {
			words = append(words, word)
		}
	}

	_ = file.Close()

	if "" == args.Skip && "" == args.Include && "" == args.Match {
		if args.All {
			log.Printf("All words of length %s:", args.Length)
			for _, word := range words {
				log.Printf("\t%s", word)
			}
		} else {
			rng := rand.New(rand.NewSource(time.Now().UnixNano()))
			inx := rng.Intn(len(words))
			log.Printf("Returning a random %s letter word: %s", args.Length, words[inx])
		}
	} else {
		log.Printf("Words that match the restrictions [length: %s, include: '%s', skip: '%s', match: '%s']:", args.Length, args.Include, args.Skip, args.Match)
		for _, word := range words {
			log.Printf("\t%s", word)
		}
	}
}

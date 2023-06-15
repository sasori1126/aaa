package trackerGenerator

import (
	"errors"
	"math"
	"math/rand"
	"strings"
	"time"
)

var charset = map[string]string{
	"numeric":               "0123456789",
	"alphabetic":            "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ",
	"alphabeticUpperCase":   "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
	"alphabeticLoweCase":    "abcdefghijklmnopqrstuvwxyz",
	"alphanumeric":          "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ",
	"alphanumericUpperCase": "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ",
	"alphanumericLowerCase": "0123456789abcdefghijklmnopqrstuvwxyz",
}

type Generator struct {
	config *Config
}

type Config struct {
	Length  int
	Count   int
	Charset string
	Pattern string
	Prefix  string
	Postfix string
}

func randomInt(min, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	index := rand.Intn(max-min) + min
	return index
}

func (g *Generator) randomElement() byte {
	cs := []byte(g.config.Charset)
	index := randomInt(0, len(cs)-1)
	return cs[index]
}

func (g *Generator) GenerateOne() (string, error) {
	var code []byte
	pattern := g.config.Pattern
	for _, i := range strings.Split(pattern, "") {
		if i == "#" {
			c := g.randomElement()
			code = append(code, c)
		} else {
			changeByte := []byte(i)
			code = append(code, changeByte[0])
		}
	}

	return g.config.Prefix + string(code) + g.config.Postfix, nil
}

func (g *Generator) Generate(count int) ([]string, error) {
	g.config.Count = count
	var codes []string
	for i := 0; i < g.config.Count; i++ {
		code, err := g.GenerateOne()
		if err != nil {
			return nil, err
		}

		codes = append(codes, code)
	}

	return codes, nil
}

func uniqueCharset(charset string) string {
	inStack := make(map[byte]bool)
	var chars []byte
	for _, b := range []byte(charset) {
		if !inStack[b] {
			chars = append(chars, b)
			inStack[b] = true
		}
	}

	return string(chars)
}

func isFeasible(config Config) bool {
	charsetLength := len(config.Charset)
	exp := math.Pow(float64(charsetLength), float64(config.Length))
	return exp >= float64(config.Count)
}

func createConfig(config Config) (*Config, error) {
	con := &Config{
		Length:  config.Length,
		Count:   config.Count,
		Charset: config.Charset,
		Pattern: config.Pattern,
		Prefix:  config.Prefix,
		Postfix: config.Postfix,
	}
	if config.Count == 0 {
		con.Count = 1
	}

	if config.Length == 0 {
		con.Length = 8
	}

	charsetMap := map[string]bool{
		"numeric":               true,
		"alphanumeric":          true,
		"alphabetic":            true,
		"alphabeticUpperCase":   true,
		"alphabeticLoweCase":    true,
		"alphanumericUpperCase": true,
		"alphanumericLowerCase": true,
	}

	if !charsetMap[config.Charset] {
		con.Charset = charset["alphanumeric"]
	} else {
		con.Charset = charset[config.Charset]
	}
	con.Charset = uniqueCharset(con.Charset)

	if config.Pattern == "" {
		pattern := []byte("#")
		var p []byte

		for i := 0; i < config.Length; i++ {
			p = append(p, pattern[0])
		}

		con.Pattern = string(p)
	} else {
		if !patternIsValid(config.Pattern) {
			return nil, errors.New("pattern is invalid")
		}

		con.Length = len(config.Pattern)
	}

	if !isFeasible(*con) {
		return nil, errors.New("not possible to generate requested codes")
	}
	return con, nil
}

func patternIsValid(pattern string) bool {
	for _, s := range strings.Split(pattern, "") {
		if s == "#" {
			return true
		}
	}

	return false
}

func New(config Config) (*Generator, error) {
	con, err := createConfig(config)
	if err != nil {
		return nil, err
	}

	return &Generator{config: con}, nil
}

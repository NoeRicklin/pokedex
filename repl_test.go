package main

import(
	"testing"
	"time"
	"github.com/NoeRicklin/pokedex/internal/pokecache"
)

type Case struct{
	input		string
	expected	[]string
}

func TestCleanInput(t *testing.T) {
	cases := []Case {
		{
			input:		"   hello world   ",
			expected:	[]string{"hello", "world"},
		},
		{
			input:		"word",
			expected:	[]string{"word"},
		},
		{
			input:		"",
			expected:	[]string{""},
		},
		{
			input:		"This, is a text. ",
			expected:	[]string{"this,", "is", "a", "text."},
		},
		{
			input:		"NO UPPER",
			expected:	[]string{"no", "upper"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		
		if len(actual) != len(c.expected) {
			t.Errorf("Slice lengths not equal. Expected %d words, got %d.",
			len(c.expected), len(actual))
		}

		for i := range actual {
			if actual[i] != c.expected[i] {
				t.Errorf("Incorrect word. Expected %s, got %s.",
				c.expected[i], actual[i])
			}
		}
	}
}

func TestCache(t *testing.T) {
	c := pokecache.NewCache(time.Millisecond * 100)

	key := "1"
	val := []byte("val")

	c.Add(key, val)
	_, cacheHit := c.Get(key)
	if !cacheHit { t.Errorf("Missed cache") }

	tm := time.NewTimer(time.Millisecond * 110)
	<-tm.C

	_, newCacheHit := c.Get(key)
	if newCacheHit { t.Errorf("Cache was not reaped") }
}


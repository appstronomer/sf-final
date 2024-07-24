// Пакет для работы с RSS-потоками.
package rss

import (
	"os"
	"reflect"
	"sf-news/pkg/storage"
	"testing"
)

func TestParse(t *testing.T) {
	postsWant := []storage.PostFull{
		{
			ID:      0,
			PubTime: 1717068660,
			Link:    "https://fake-site.com/news/1?utm_medium=rss",
			Title:   "Article 1",
			Content: "Article 1 first paragraph\r\n\r\nArticle 1 second paragraph",
		},
		{
			ID:      0,
			PubTime: 1714435200,
			Link:    "https://fake-site.com/news/2?utm_medium=rss",
			Title:   "Article 2",
			Content: "Article 2 first paragraph\r\n\r\nArticle 2 second paragraph",
		},
	}

	rssReader, err := os.Open("testdata/rss.xml")
	if err != nil {
		t.Errorf("read testdata/rss.xml error = %v", err)
	}
	defer rssReader.Close()

	postsGot, err := Parse(rssReader)
	if err != nil {
		t.Errorf("Parse() error = %v", err)
	}
	if len(postsGot) != len(postsWant) {
		t.Errorf("posts count got %v; want %v", len(postsGot), len(postsWant))
	}
	for i := 0; i < len(postsWant); i++ {
		if !reflect.DeepEqual(postsGot[i], postsWant[i]) {
			t.Errorf("posts[%v] got %#v; want %#v", i, postsGot[i], postsWant[i])
		}
	}
}

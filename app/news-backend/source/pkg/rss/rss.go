// Пакет для работы с RSS-потоками.
package rss

import (
	"encoding/xml"
	"io"
	"strings"
	"time"

	"sf-news/pkg/storage"

	"github.com/k3a/html2text"
)

type RssData struct {
	Chanel Channel `xml:"channel"`
}

type Channel struct {
	Items []Item `xml:"item"`
}

type Item struct {
	PubDate     string `xml:"pubDate"`
	Link        string `xml:"link"`
	Title       string `xml:"title"`
	Description string `xml:"description"`
}

func Parse(dataReader io.Reader) ([]storage.PostFull, error) {
	bytesData, err := io.ReadAll(dataReader)
	if err != nil {
		return nil, err
	}
	var rssData RssData
	err = xml.Unmarshal(bytesData, &rssData)
	if err != nil {
		return nil, err
	}
	var dataPosts []storage.PostFull
	for _, item := range rssData.Chanel.Items {
		var post storage.PostFull
		post.Title = item.Title
		post.Content = html2text.HTML2Text(item.Description)
		post.Link = item.Link
		item.PubDate = strings.ReplaceAll(item.PubDate, ",", "")
		pubTime, err := time.Parse("Mon 2 Jan 2006 15:04:05 -0700", item.PubDate)
		if err != nil {
			pubTime, err = time.Parse("Mon 2 Jan 2006 15:04:05 GMT", item.PubDate)
		}
		if err == nil {
			post.PubTime = pubTime.Unix()
		}
		dataPosts = append(dataPosts, post)
	}
	return dataPosts, nil
}

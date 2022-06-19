package main

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"os"
	"sort"
)

const RssFile = "rss.xml"

type Rss struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Chan    Channel  `xml:"channel"`
}

type Channel struct {
	Title string `xml:"title"`
	Link  string `xml:"link"`
	Desc  string `xml:"description"`
	Items []Item `xml:"item"`
}

type Item struct {
	Title   string `xml:"title"`
	Link    string `xml:"link"`
	PubDate string `xml:"pubDate"`
	Desc    string `xml:"description"`
	Guid    string `xml:"guid"`
}

func writeRssFile(rss Rss) {
	data, marshalErr := xml.MarshalIndent(rss, "", "  ")

	if marshalErr != nil {
		fmt.Println("failed to marshal XML")
		os.Exit(0)
	}

	os.WriteFile(RssFile, data, 0644)
}

func loadRss() Rss {
	data, readErr := os.ReadFile(RssFile)

	if readErr != nil {
		fmt.Println("Failed to open RSS file. Make sure to run 'rss new' to generate a file.")
		os.Exit(0)
	}

	var rss Rss
	xml.Unmarshal(data, &rss)

	return rss
}

func createRssFile() {
	in := bufio.NewReader(os.Stdin)

	fmt.Print("Channel title: ")
	channelTitle, _ := in.ReadString('\n')
	fmt.Print("Channel link: ")
	channelLink, _ := in.ReadString('\n')
	fmt.Print("Channel description: ")
	channelDesc, _ := in.ReadString('\n')

	channelTitle = channelTitle[:len(channelTitle)-1]
	channelLink = channelLink[:len(channelLink)-1]
	channelDesc = channelDesc[:len(channelDesc)-1]

	rss := Rss{
		Version: "2.0",
		Chan: Channel{
			Title: channelTitle,
			Link:  channelLink,
			Desc:  channelDesc,
			Items: []Item{},
		},
	}

	writeRssFile(rss)
}

func addItem() {
	rss := loadRss()

	in := bufio.NewReader(os.Stdin)

	fmt.Print("Item title: ")
	itemTitle, _ := in.ReadString('\n')
	fmt.Print("Item link: ")
	itemLink, _ := in.ReadString('\n')
	fmt.Print("Item publication date: ")
	itemPubDate, _ := in.ReadString('\n')
	fmt.Print("Item description: ")
	itemDesc, _ := in.ReadString('\n')

	// remove newline characters
	itemTitle = itemTitle[:len(itemTitle)-1]
	itemLink = itemLink[:len(itemLink)-1]
	itemPubDate = itemPubDate[:len(itemPubDate)-1]
	itemDesc = itemDesc[:len(itemDesc)-1]

	newItem := Item{
		Title:   itemTitle,
		Link:    itemLink,
		PubDate: itemPubDate,
		Desc:    itemDesc,
		Guid:    itemLink,
	}

	rss.Chan.Items = append(rss.Chan.Items, newItem)

	// sort items by date in descending order
	sort.Slice(rss.Chan.Items, func(i, j int) bool {
		return rss.Chan.Items[i].PubDate > rss.Chan.Items[j].PubDate
	})

	writeRssFile(rss)
}

func removeItem() {
	rss := loadRss()

	for i := range rss.Chan.Items {
		fmt.Printf("%2d: %s\n", i, rss.Chan.Items[i].Title)
	}

	var indxToRemove int
	fmt.Print("Item to remove: ")
	fmt.Scanln(&indxToRemove)

	rss.Chan.Items = append(rss.Chan.Items[:indxToRemove], rss.Chan.Items[indxToRemove+1:]...)

	writeRssFile(rss)
}

func displayHelp() {
	fmt.Printf("%6s: Creates a new rss.xml file. You will be prompted for the information needed.\n", "new")
	fmt.Printf("%6s: Adds and item to the rss file. You will be prompted for information. "+
		"The new items gets inserted in order by date with newest items first.\n", "add")
	fmt.Printf("%6s: Displays all items in the file. You can enter the number of the item to remove.\n", "remove")
	fmt.Printf("%6s: Displays information about available commands.\n", "help")
}

func main() {
	args := os.Args

	if len(args) < 2 {
		fmt.Println("provide at least one argument")
		os.Exit(0)
	}

	switch args[1] {
	case "new":
		createRssFile()
	case "add":
		addItem()
	case "remove":
		removeItem()
	case "help":
		displayHelp()
	default:
		fmt.Printf("argument %s not recognized\n\n", args[1])
		displayHelp()
	}
}

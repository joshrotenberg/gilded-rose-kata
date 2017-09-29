package main

import "fmt"

// Item represents an item in the inventory
type Item struct {
	name            string
	sellIn, quality int
}

var items = []Item{
	Item{"+5 Dexterity Vest", 10, 20},
	Item{"Aged Brie", 2, 0},
	Item{"Elixir of the Mongoose", 5, 7},
	Item{"Sulfuras, Hand of Ragnaros", 0, 80},
	Item{"Backstage passes to a TAFKAL80ETC concert", 15, 20},
	Item{"Conjured Mana Cake", 3, 6},
}

func incrementQuality(item *Item, amount int, max int) {
	item.quality += amount
	if item.quality > max {
		item.quality = max
	}
}
func decrementQuality(item *Item, amount int) {
	item.quality -= amount
	if item.quality < 0 {
		item.quality = 0
	}
}

// UpdateItem updates the sell-in and quality for the given item
func UpdateItem(item *Item) {
	switch item.name {
	case "Aged Brie":
		item.sellIn--
		incrementQuality(item, 1, 50)

	case "Sulfuras, Hand of Ragnaros":
		// force 80 quality if it isn't already. does "never has to be sold" mean sellIn is always ... 1? 0? does it matter?
		item.quality = 80
	case "Backstage passes to a TAFKAL80ETC concert":
		if item.sellIn <= 5 && item.sellIn > 0 {
			incrementQuality(item, 3, 50)
		} else if item.sellIn <= 10 && item.sellIn > 0 {
			incrementQuality(item, 2, 50)
		} else if item.sellIn <= 0 {
			item.quality = 0
		} else {
			incrementQuality(item, 1, 50)
		}
		item.sellIn--

	default:
		if item.sellIn < 1 {
			decrementQuality(item, 2)
		} else {
			decrementQuality(item, 1)
		}
		item.sellIn--
	}
}

func main() {
	for i := 0; i < len(items); i++ {
		UpdateItem(&items[i])
	}
	fmt.Println(items)
}

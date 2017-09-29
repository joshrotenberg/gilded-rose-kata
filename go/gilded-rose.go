package main

import (
	"fmt"
	"regexp"
)

// itemNameRegex matches the item's name and conjured status
var itemNameRegex = regexp.MustCompile(`(?i)^(conjured\s+)?(.*)`)

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

// isConjured returns true if the item is conjured.
func isConjured(item Item) bool {
	match := itemNameRegex.FindStringSubmatch(item.name)
	return match[1] != ""
}

// incrementQuality increments the item's quality if the new quality doesn't exceed the max.
func incrementQuality(item *Item, amount int, max int) {
	item.quality += amount
	if item.quality > max {
		item.quality = max
	}
}

// decrementQuality decrements the item's quality, testing to see if the item is conjured, and respects
// the constraint that quality is never negative.
func decrementQuality(item *Item, amount int) {
	item.quality -= amount
	if isConjured(*item) {
		item.quality--
	}
	if item.quality < 0 {
		item.quality = 0
	}
}

// UpdateItem updates the sell-in and quality for the given item. Adding an item with special rules requires
// adding a new case matching the item's name. Items without special rules are handled by the default case.
func UpdateItem(item *Item) {
	switch item.name {
	// "Aged Brie" actually increases in quality the older it gets
	case "Aged Brie":
		item.sellIn--
		incrementQuality(item, 1, 50)

	// "Sulfuras", being a legendary item, never has to be sold or decreases in quality
	case "Sulfuras, Hand of Ragnaros":
		// force 80 quality if it isn't already. does "never has to be sold" mean sellIn is always ... 1? 0? does it matter?
		item.quality = 80

	// "Backstage passes", like aged brie, increases in quality as it's sell-in value approaches; quality increases by 2 when
	// there are 10 days or less and by 3 when there are 5 days or less but quality drops to 0 after the concert
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

	// All other items will be handled here.
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

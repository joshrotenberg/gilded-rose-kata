package main

import (
	"fmt"
	"testing"
)

const (
	vest     string = "+5 Dexterity Vest"
	brie     string = "Aged Brie"
	elixir   string = "Elixir of the Mongoose"
	sulfuras string = "Sulfuras, Hand of Ragnaros"
	passes   string = "Backstage passes to a TAFKAL80ETC concert"
	cake     string = "Conjured Mana Cake"
)

type TestItem struct {
	item            Item
	expcetedSellIn  int
	expectedQuality int
}

func runTestCases(t *testing.T, name string, testItems ...TestItem) {

	for _, testItem := range testItems {
		t.Run(fmt.Sprintf("%s: failure for item '%s", name, testItem.item.name), func(t *testing.T) {
			UpdateItem(&testItem.item)
			if testItem.item.quality != testItem.expectedQuality {
				t.Fatalf("expceted quality update for %s: expceted %d, got %d", testItem.item.name, testItem.expectedQuality, testItem.item.quality)
			}
			if testItem.item.sellIn != testItem.expcetedSellIn {
				t.Fatalf("unexpceted sellIn update for %s: expceted %d, got %d", testItem.item.name, testItem.expcetedSellIn, testItem.item.sellIn)
			}
		})
	}
}
func Test_CurrentInventory(t *testing.T) {

	runTestCases(
		t,
		"Current Inventory",
		TestItem{Item{vest, 10, 20}, 9, 19},
		TestItem{Item{brie, 2, 0}, 1, 1},
		TestItem{Item{elixir, 5, 7}, 4, 6},
		TestItem{Item{sulfuras, 0, 80}, 0, 80},
		TestItem{Item{passes, 15, 20}, 14, 21},
		// TestItem{Item{cake, 3, 6}, 3, 4}, // currently not implemented
	)
}

func Test_QualityDegradesTwiceAsFast(t *testing.T) {

	runTestCases(
		t,
		"Once the sell by date has passed, quality degrades twice as fast",
		TestItem{Item{vest, -1, 6}, -2, 4},
		TestItem{Item{brie, -1, 20}, -2, 21}, // XXX: docs say quality degrades twice as fast, but not increases twice as fast, so 21?
		TestItem{Item{elixir, -1, 7}, -2, 5},
		TestItem{Item{sulfuras, -2, 8}, -2, 80},
		TestItem{Item{passes, -3, 0}, -4, 0},
	)
}

func Test_AgedBrieIncreases(t *testing.T) {

	runTestCases(
		t,
		"'Aged Brie' actually increases in quality the older it gets",
		TestItem{Item{brie, 3, 20}, 2, 21},
		TestItem{Item{brie, 2, 5}, 1, 6},
	)
}

func Test_Sulfuras(t *testing.T) {

	runTestCases(
		t,
		"'Sulfuras', being a legendary item, never has to be sold or decreases in quality",
		TestItem{Item{sulfuras, 5, 20}, 5, 80},
		TestItem{Item{sulfuras, 50, 200}, 50, 80},
	)
}

func Test_BackstatePasses(t *testing.T) {

	runTestCases(
		t,
		"'Backstage passes', like aged brie, increases in quality as it's sell-in value approaches ...",
		TestItem{Item{passes, 11, 7}, 10, 8},
		TestItem{Item{passes, 10, 7}, 9, 9},
		TestItem{Item{passes, 6, 7}, 5, 9},
		TestItem{Item{passes, 5, 7}, 4, 10},
		TestItem{Item{passes, 0, 7}, -1, 0},
		TestItem{Item{passes, -2, 7}, -3, 0},
	)
}

func Test_ConjuredItems(t *testing.T) {

	runTestCases(
		t,
		"'Conjured' items degrade in quality twice as fast as normal items",
		TestItem{Item{cake, 5, 6}, 4, 4},
	)
}

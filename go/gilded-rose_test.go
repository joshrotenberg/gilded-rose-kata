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
type TestCase struct {
	name  string
	items []TestItem
}

func generateTestCase(name string, items ...TestItem) TestCase {

	return TestCase{name, items}
}

func runTestCases(t *testing.T, testCase TestCase) {

	for _, testItem := range testCase.items {
		t.Run(fmt.Sprintf("%s: failure for item '%s", testCase.name, testItem.item.name), func(t *testing.T) {
			err := UpdateItem(&testItem.item)
			if err != nil {
				t.Fatalf("failed to update %s: %+v", testItem.item.name, err)
			}
			if testItem.item.quality != testItem.expectedQuality {
				t.Fatalf("unexpceted quality update for %s: expceted %d, got %d", testItem.item.name, testItem.expectedQuality, testItem.item.quality)
			}
			if testItem.item.sellIn != testItem.expcetedSellIn {
				t.Fatalf("unexpceted sellIn update for %s: expceted %d, got %d", testItem.item.name, testItem.expcetedSellIn, testItem.item.sellIn)
			}
		})
	}
}
func Test_CurrentInventory(t *testing.T) {

	testCase := generateTestCase(
		"Current Inventory",
		TestItem{Item{vest, 10, 20}, 9, 19},
		TestItem{Item{brie, 2, 0}, 1, 1},
		TestItem{Item{elixir, 5, 7}, 4, 6},
		TestItem{Item{sulfuras, 0, 80}, 0, 80},
		TestItem{Item{passes, 15, 20}, 14, 21},
		// TestItem{Item{cake, 3, 6}, 3, 4}, // currently not implemented
	)
	runTestCases(t, testCase)
}

(ns gilded-rose.core-spec
  (:require [speclj.core :refer :all]
            [gilded-rose.core :refer [item
                                      update-quality
                                      update-current-inventory]]))
;; reduce typo likelihood
(def pantry {:vest "+5 Dexterity Vest"
             :brie "Aged Brie"
             :elixir "Elixir of the Mongoose"
             :hand "Sulfuras, Hand of Ragnaros"
             :passes "Backstage passes to a TAFKAL80ETC concert"})

;; wrap update unwrap for a single item convenience fn
(def single-update (comp first update-quality vector))
(println (single-update (item (:brie pantry) 10 59)))
(defn accurate-update?
  "Given an item name and current/expected sell-in and quality parameters,
  returns true if the update matches expectations."
  [name & {:keys [sell-in quality expected-sell-in expected-quality]}]
  (let [in (item name sell-in quality)
        out (item name expected-sell-in expected-quality)]
    (= (single-update in) out)))

(describe "gilded rose"
          (it "returns items in the expected format"
              (should (= (item (:vest pantry) 1 2)
                         {:name (:vest pantry) :sell-in 1 :quality 2})))

          (it "updates the item correctly based on current inventory"

              (should (accurate-update? (:vest pantry)
                                        :sell-in 10 :quality 20
                                        :expected-sell-in 9 :expected-quality 19))

              (should (accurate-update? (:brie pantry)
                                        :sell-in 2 :quality 0
                                        :expected-sell-in 1 :expected-quality 1))

              (should (accurate-update? (:elixir pantry)
                                        :sell-in 5 :quality 7
                                        :expected-sell-in 4 :expected-quality 6))

              (should (accurate-update? (:hand pantry)
                                        :sell-in 0 :quality 80
                                        :expected-sell-in 0 :expected-quality 80))

              (should (accurate-update? (:passes pantry)
                                        :sell-in 15 :quality 20
                                        :expected-sell-in 14 :expected-quality 21)))

          (it "quality degrades twice as fast once sell by has passed"

              (should (accurate-update? (:vest pantry)
                                        :sell-in -1 :quality 6
                                        :expected-sell-in -2 :expected-quality 4))

              ;; increases in quality no matter what
              (should (accurate-update? (:brie pantry)
                                        :sell-in -1 :quality 20
                                        :expected-sell-in -2 :expected-quality 21))

              (should (accurate-update? (:elixir pantry)
                                        :sell-in -1 :quality 7
                                        :expected-sell-in -2 :expected-quality 5))

              ;; never decreases in quality
              (should (accurate-update? (:hand pantry)
                                        :sell-in -2 :quality 8
                                        :expected-sell-in -2 :expected-quality 8))

              ;; passes are useless (quality is 0) after sell-by
              (should (accurate-update? (:passes pantry)
                                        :sell-in -3 :quality 0
                                        :expected-sell-in -4 :expected-quality 0)))

          (it "Aged Brie increases in quality the older it gets"
              (should (accurate-update? (:brie pantry)
                                        :sell-in 3 :quality 20
                                        :expected-sell-in 2 :expected-quality 21))
              (should (accurate-update? (:brie pantry)
                                        :sell-in 2 :quality 5
                                        :expected-sell-in 1 :expected-quality 6)))


          (it "sulfuras never has to be sold nor does it decrease in quality"
              (should (accurate-update? (:hand pantry)
                                        :sell-in 5 :quality 20
                                        :expected-sell-in 5 :expected-quality 20))
              (should (accurate-update? (:hand pantry)
                                        :sell-in 50 :quality 200
                                        :expected-sell-in 50 :expected-quality 200)))

          (it "backstage passes have all kinds of special rules"
              ;; more than 10 days
              (should (accurate-update? (:passes pantry)
                                       :sell-in 11 :quality 7
                                       :expected-sell-in 10 :expected-quality 8))
              ;; between 6 and 10 days
              (should (accurate-update? (:passes pantry)
                                        :sell-in 10 :quality 7
                                        :expected-sell-in 9 :expected-quality 9))
              (should (accurate-update? (:passes pantry)
                                        :sell-in 6 :quality 7
                                        :expected-sell-in 5 :expected-quality 9))
              ;; 5 or less
              (should (accurate-update? (:passes pantry)
                                        :sell-in 5 :quality 7
                                        :expected-sell-in 4 :expected-quality 10))
              (should (accurate-update? (:passes pantry)
                                        :sell-in 1 :quality 7
                                        :expected-sell-in 0 :expected-quality 10))
              ;; missed the concert, they are worthless now
              (should (accurate-update? (:passes pantry)
                                        :sell-in 0 :quality 7
                                        :expected-sell-in -1 :expected-quality 0))
              (should (accurate-update? (:passes pantry)
                                        :sell-in -2 :quality 7
                                        :expected-sell-in -3 :expected-quality 0)))

          ;; these next two may be implicit (i.e. the inventory will never start with their quality
          ;; being outside the paramters) but they don't actually appear to be enforced explicitly,
          ;; or at least not correctly before the updates are made.
          (it "the quality of an item is never more than 50"
              (should (accurate-update? (:vest pantry)
                                        :sell-in 5 :quality 50
                                        :expected-sell-in 4 :expected-quality 51))
              (should (accurate-update? (:brie pantry)
                                        :sell-in 5 :quality 50
                                        :expected-sell-in 4 :expected-quality 50)))

          (it "quality is never negative"
              (should (accurate-update? (:elixir pantry)
                                        :sell-in 1 :quality 0
                                        :expected-sell-in -4 :expected-quality 0)))
          (it "conjured items should degrade twice as fast as normal items")

          )

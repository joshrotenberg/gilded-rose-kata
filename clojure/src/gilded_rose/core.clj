(ns gilded-rose.core)

(def entry-pattern #"(?i)^(conjured\s+)?(.*)")

(defn parse-item
  "Takes an inventory item and returns a vector of match information on the name based around
  whether the item is conjured or not."
  [item]
  (re-matches entry-pattern (:name item)))

(defn item-name
  "Convenience function to pull out only the standard name of an item."
  [item]
  (last (parse-item (:name item))))

(defn conjured?
  "Returns true if the item matches our understanding of a conjured item."
  [item]
  (if (second (parse-item item))
    true
    false))

(defn item-details
  "Returns the details of an item, namely its sell-in and quality features."
  [item]
  [(:sell-in item) (:quality item)])

(defn degrade-quality
  "Given the expiration and conjured status and the current quality of an item, returns an
  updated quality."
  [expired? conjured? quality]
  (let [updated-quality (if (or conjured? expired?)
                          (- quality 2)
                          (dec quality))]
    (if (neg? updated-quality)
      0
      updated-quality)))

(defmulti update-item
  "Handle item updates, dispatched on their name (or default to the standard update rules)."
  (fn [item] (last (parse-item item))))

;; "Aged Brie" actually increases in quality the older it gets
(defmethod update-item "Aged Brie" [item]
  (let [[sell-in quality] (item-details item)]
    (merge item
           {:sell-in (dec sell-in)
            :quality (if (< quality 50)
                       (inc quality)
                       quality)})))

;;"Backstage passes", like aged brie, increases in quality as it's sell-in value approaches;
;; quality increases by 2 when there are 10 days or less and by 3 when there are 5 days or
;; less but quality drops to 0 after the concert
(defmethod update-item "Backstage passes to a TAFKAL80ETC concert" [item]
  (let [[sell-in quality] (item-details item)
        expired? (zero? sell-in)
        conjured? (conjured? item)]
    (merge item
           {:sell-in (dec sell-in)
            :quality (cond
                       (>= 0 sell-in) 0
                       (<= sell-in 5) (+ quality 3)
                       (<= sell-in 10) (+ quality 2)
                       :else (inc quality))})))

;; "Sulfuras", being a legendary item, never has to be sold or decreases in quality
(defmethod update-item "Sulfuras, Hand of Ragnaros" [item] item)

;; covers all other products that don't have special rules
(defmethod update-item :default [item]
  (let [[sell-in quality] (item-details item)
        expired? (> 1 sell-in)
        conjured? (conjured? item)]
    (merge item
           {:sell-in (dec sell-in)
            :quality (degrade-quality expired? conjured? quality)})))

(defn update-quality
  [items]
  (map update-item items))

(defn item [item-name, sell-in, quality]
  {:name item-name, :sell-in sell-in, :quality quality})

(defn update-current-inventory[]
  (let [inventory
    [
     (item "+5 Dexterity Vest" 10 20)
     (item "Aged Brie" 2 0)
     (item "Elixir of the Mongoose" 5 7)
     (item "Sulfuras, Hand Of Ragnaros" 0 80)
     (item "Backstage passes to a TAFKAL80ETC concert" 15 20)
     (item "Conjured Mana Cake" 3 6)
     ]]
    (update-quality inventory)
    ))

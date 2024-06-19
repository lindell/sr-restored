package scraper

import "github.com/PuerkitoBio/goquery"

func docMap[E any](sel *goquery.Selection, fn func(val *goquery.Selection) (E, error)) ([]E, error) {
	var err error
	var val E
	var ret []E
	sel.EachWithBreak(func(i int, s *goquery.Selection) bool {
		val, err = fn(s)
		if err != nil {
			return false
		}
		ret = append(ret, val)
		return true
	})
	return ret, err
}

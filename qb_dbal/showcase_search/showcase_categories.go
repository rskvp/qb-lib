package showcase_search

import (
	qbc "github.com/rskvp/qb-core"
	"github.com/rskvp/qb-lib/qb_dbal/commons"
)

type ShowcaseCategoryWeight struct {
	WeightInDate  int `json:"weight_in_date"`
	WeightOutDate int `json:"weight_out_date"`
}

type ShowcaseCategories struct {
	root string
	data map[string]ShowcaseCategoryWeight

	filename string
}

func NewShowcaseCategories(root string) *ShowcaseCategories {
	instance := new(ShowcaseCategories)
	instance.root = root
	instance.data = make(map[string]ShowcaseCategoryWeight)

	instance.init()

	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *ShowcaseCategories) Clear() {
	filename := instance.filename
	if len(filename) > 0 {
		if b, _ := qbc.Paths.Exists(filename); b {
			_ = qbc.IO.Remove(filename)
		}
	}
}

func (instance *ShowcaseCategories) Get(name string) ShowcaseCategoryWeight {
	return GetCategoryWeight(instance.data, name)
}

func (instance *ShowcaseCategories) SetWeight(category string, inTime bool, weight int) ShowcaseCategoryWeight {
	if nil != instance {
		if v, b := instance.data[category]; b {
			if inTime {
				v.WeightInDate = weight
			} else {
				v.WeightOutDate = weight
			}
			instance.data[category] = v

			instance.save()

			return v
		}
	}
	return ShowcaseCategoryWeight{
		WeightInDate:  1,
		WeightOutDate: 1,
	}
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *ShowcaseCategories) init() {
	instance.filename = qbc.Paths.Concat(instance.root, "categories.json")
	if b, _ := qbc.Paths.Exists(instance.filename); b {
		// load
		_ = qbc.JSON.ReadFromFile(instance.filename, &instance.data)
	}
	if len(instance.data) == 0 {
		// add defaults
		instance.data[commons.CAT_PERSON] = ShowcaseCategoryWeight{WeightInDate: 2, WeightOutDate: 2}
		instance.data[commons.CAT_EVENT] = ShowcaseCategoryWeight{WeightInDate: 4, WeightOutDate: 1}
		instance.data[commons.CAT_ADV] = ShowcaseCategoryWeight{WeightInDate: 1, WeightOutDate: 1}
		instance.data[commons.CAT_DOCUMENT] = ShowcaseCategoryWeight{WeightInDate: 2, WeightOutDate: 1}
		instance.data[commons.CAT_POST] = ShowcaseCategoryWeight{WeightInDate: 4, WeightOutDate: 1}
	}

	instance.save()
}

func (instance *ShowcaseCategories) save() {
	// save to file
	if len(instance.filename) > 0 {
		_, _ = qbc.IO.WriteTextToFile(qbc.JSON.Stringify(instance.data), instance.filename)
	}
}

//----------------------------------------------------------------------------------------------------------------------
//	S T A T I C
//----------------------------------------------------------------------------------------------------------------------

func GetCategoryWeight(data map[string]ShowcaseCategoryWeight, name string) ShowcaseCategoryWeight {
	if v, b := data[name]; b {
		return v
	}
	return ShowcaseCategoryWeight{
		WeightInDate:  1,
		WeightOutDate: 1,
	}
}

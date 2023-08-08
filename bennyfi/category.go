// The MIT License (MIT)

// Copyright (c) 2020, Digital Scarcity

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.
package bennyfi

import (
	"encoding/json"
	"fmt"

	eos "github.com/sebastianmontero/eos-go"
)

type Category struct {
	Category            eos.Name `json:"category"`
	CategoryName        string   `json:"category_name"`
	CategoryDescription string   `json:"category_description"`
	CategoryImage       string   `json:"category_image"`
}

func (m *Category) String() string {
	result, err := json.Marshal(m)
	if err != nil {
		panic(fmt.Sprintf("Failed marshalling round: %v", err))
	}
	return string(result)
}

type SetCategoryArgs struct {
	*Category
	Authorizer eos.AccountName `json:"authorizer"`
}

func (m *SetCategoryArgs) String() string {
	result, err := json.Marshal(m)
	if err != nil {
		panic(fmt.Sprintf("Failed marshalling round: %v", err))
	}
	return string(result)
}

type EraseCategoryArgs struct {
	Category   eos.Name        `json:"category"`
	Authorizer eos.AccountName `json:"authorizer"`
}

func (m *EraseCategoryArgs) String() string {
	result, err := json.Marshal(m)
	if err != nil {
		panic(fmt.Sprintf("Failed marshalling round: %v", err))
	}
	return string(result)
}
func (m *BennyfiContract) SetCategory(categoryArgs *SetCategoryArgs) (string, error) {
	// actionData := make(map[string]interface{})
	// actionData["authorizer"] = categoryArgs.Authorizer
	// actionData["beneficiary"] = categoryArgs.Beneficiary
	// actionData["attributes"] = categoryArgs.Attributes
	return m.ExecAction(categoryArgs.Authorizer, "setcategory", categoryArgs)
}

func (m *BennyfiContract) EraseCategory(category eos.Name, authorizer eos.AccountName) (string, error) {
	actionData := &EraseCategoryArgs{
		Category:   category,
		Authorizer: authorizer,
	}
	return m.ExecAction(authorizer, "erasectgry", actionData)
}

func (m *BennyfiContract) GetCategoriesReq(req *eos.GetTableRowsRequest) ([]*Category, error) {

	var categories []*Category
	if req == nil {
		req = &eos.GetTableRowsRequest{}
	}
	req.Table = "categories"
	err := m.GetTableRows(*req, &categories)
	if err != nil {
		return nil, fmt.Errorf("get table rows %v", err)
	}
	return categories, nil
}

func (m *BennyfiContract) GetAllCategoriesAsMap() ([]map[string]interface{}, error) {
	req := eos.GetTableRowsRequest{
		Table: "categories",
	}
	return m.GetAllTableRowsAsMap(req, "category")
}

func (m *BennyfiContract) GetCategoryById(category eos.Name) (*Category, error) {
	request := &eos.GetTableRowsRequest{}
	m.FilterCategoriesById(request, category)
	categories, err := m.GetCategoriesReq(request)
	if err != nil {
		return nil, err
	}
	if len(categories) > 0 {
		return categories[0], nil
	}
	return nil, nil
}

func (m *BennyfiContract) FilterCategoriesById(req *eos.GetTableRowsRequest, category eos.Name) {
	req.LowerBound = string(category)
	req.UpperBound = string(category)
}

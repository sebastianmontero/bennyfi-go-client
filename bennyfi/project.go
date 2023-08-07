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
	"fmt"

	eos "github.com/sebastianmontero/eos-go"
	"github.com/sebastianmontero/eos-go-toolbox/dto"
)

var (
	ProjectAttrName        = "project_name"
	ProjectAttrArtifactCID = "artifact_cid"
)

type Attribute struct {
	Key   string         `json:"first"`
	Value *dto.FlexValue `json:"second"`
}

type Attributes []*Attribute

func (m Attributes) FindPos(key string) int {
	for i, attr := range m {
		if attr.Key == key {
			return i
		}
	}
	return -1
}

func (m Attributes) Find(key string) *Attribute {
	pos := m.FindPos(key)
	if pos >= 0 {
		return m[pos]
	}
	return nil
}

type Project struct {
	ProjectID   uint64          `json:"project_id"`
	Authorizer  eos.AccountName `json:"authorizer"`
	Beneficiary eos.AccountName `json:"beneficiary"`
	Attributes  Attributes      `json:"attributes"`
	CreatedDate eos.TimePoint   `json:"created_date"`
	UpdatedDate eos.TimePoint   `json:"updated_date"`
}

type SetProjectArgs struct {
	ProjectID   uint64          `json:"project_id"`
	Authorizer  eos.AccountName `json:"authorizer"`
	Beneficiary eos.AccountName `json:"beneficiary"`
	Attributes  Attributes      `json:"attributes"`
}

func (m *Project) ToSetProjectArgs() *SetProjectArgs {
	return &SetProjectArgs{
		ProjectID:   m.ProjectID,
		Authorizer:  m.Authorizer,
		Beneficiary: m.Beneficiary,
		Attributes:  m.Attributes,
	}
}

func (m *BennyfiContract) SetProject(projectArgs *SetProjectArgs) (string, error) {
	return m.ExecAction(projectArgs.Authorizer, "setproject", projectArgs)
}

func (m *BennyfiContract) SetProjectFromProject(project *Project) (string, error) {
	return m.SetProject(project.ToSetProjectArgs())
}

func (m *BennyfiContract) GetProjectsReq(req *eos.GetTableRowsRequest) ([]Project, error) {

	var projects []Project
	if req == nil {
		req = &eos.GetTableRowsRequest{}
	}
	req.Table = "projects"
	err := m.GetTableRows(*req, &projects)
	if err != nil {
		return nil, fmt.Errorf("get table rows %v", err)
	}
	return projects, nil
}

func (m *BennyfiContract) GetLastProject() (*Project, error) {
	projects, err := m.GetProjectsReq(&eos.GetTableRowsRequest{
		Reverse: true,
		Limit:   1,
	})
	if err != nil {
		return nil, err
	}
	if len(projects) > 0 {
		return &projects[0], nil
	}
	return nil, nil
}

func (m *BennyfiContract) GetProjectsByAuthorizerAndId(authorizer interface{}, projectIDUpperBound uint64) ([]Project, error) {
	request := &eos.GetTableRowsRequest{}
	err := m.FilterProjectsByAuthorizerAndId(request, authorizer, projectIDUpperBound)
	if err != nil {
		return nil, err
	}
	return m.GetProjectsReq(request)
}

func (m *BennyfiContract) FilterProjectsByAuthorizerAndId(req *eos.GetTableRowsRequest, authorizer interface{}, projectIDUpperBound uint64) error {

	req.Index = "2"
	req.KeyType = "i128"
	req.Reverse = true
	authAndRndLB, err := m.EOS.GetComposedIndexValue(authorizer, 0)
	if err != nil {
		return fmt.Errorf("failed to generate lower bound composed index, err: %v", err)
	}
	if projectIDUpperBound == 0 {
		projectIDUpperBound = 18446744073709551615
	}
	authAndRndUB, err := m.EOS.GetComposedIndexValue(authorizer, projectIDUpperBound)
	if err != nil {
		return fmt.Errorf("failed to generate upper bound composed index, err: %v", err)
	}
	fmt.Println("LB: ", authAndRndLB, "UB: ", authAndRndUB)
	req.LowerBound = authAndRndLB
	req.UpperBound = authAndRndUB
	return err
}

func (m *BennyfiContract) GetProjectsByBeneficiaryAndId(beneficiary interface{}, projectIDUpperBound uint64) ([]Project, error) {
	request := &eos.GetTableRowsRequest{}
	err := m.FilterProjectsByBeneficiaryAndId(request, beneficiary, projectIDUpperBound)
	if err != nil {
		return nil, err
	}
	return m.GetProjectsReq(request)
}

func (m *BennyfiContract) FilterProjectsByBeneficiaryAndId(req *eos.GetTableRowsRequest, beneficiary interface{}, projectIDUpperBound uint64) error {

	req.Index = "3"
	req.KeyType = "i128"
	req.Reverse = true
	beneAndRndLB, err := m.EOS.GetComposedIndexValue(beneficiary, 0)
	if err != nil {
		return fmt.Errorf("failed to generate lower bound composed index, err: %v", err)
	}
	if projectIDUpperBound == 0 {
		projectIDUpperBound = 18446744073709551615
	}
	beneAndRndUB, err := m.EOS.GetComposedIndexValue(beneficiary, projectIDUpperBound)
	if err != nil {
		return fmt.Errorf("failed to generate upper bound composed index, err: %v", err)
	}
	fmt.Println("LB: ", beneAndRndLB, "UB: ", beneAndRndUB)
	req.LowerBound = beneAndRndLB
	req.UpperBound = beneAndRndUB
	return err
}

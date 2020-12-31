package models

import (
	"strconv"
	"time"
)

type CardOption struct {
	Name     string
	Value    string
	Selected bool
}

type ViewModel struct {
	PathPrefix   string
	Birthday     string
	Expiration   string
	Manipulation bool
	CardOptions  [3]CardOption
	Parts        []string
}

const (
	NewID    = "newid"
	OldID    = "oldid"
	Passport = "passport"
)

func NewViewModel() ViewModel {
	model := ViewModel{}

	// prepare routing prefix
	model.PathPrefix = "/"

	// prepare birthday and expiration
	model.Birthday = "1980-01-01"
	model.Expiration = "2030-12-31"
	model.Manipulation = false
	model.Parts = nil

	// prepare default card selection
	model.CardOptions[0] = CardOption{Name: "Neuer Personalausweis", Value: NewID, Selected: false}
	model.CardOptions[1] = CardOption{Name: "Alter Personalausweis", Value: OldID, Selected: false}
	model.CardOptions[2] = CardOption{Name: "EU-Reisepass", Value: Passport, Selected: false}

	return model
}

func (vm *ViewModel) CalculateNewID() {
	cardNumberBlock := "T220001293"

	bdBlock := transformDate(vm.Birthday)
	bdBlock += calculateChecksumOfBlock(bdBlock, false)

	expBlock := transformDate(vm.Expiration)
	expBlock += calculateChecksumOfBlock(expBlock, false)

	// overall checksum
	checksum := calculateChecksumOfBlock(cardNumberBlock+bdBlock+expBlock, vm.Manipulation)

	result := make([]string, 4)
	result[0] = cardNumberBlock
	result[1] = bdBlock
	result[2] = expBlock + "D"
	result[3] = checksum
	vm.Parts = result
}

func (vm *ViewModel) CalculateOldID() {
	cardNumberBlock := "1220001297"

	bdBlock := transformDate(vm.Birthday)
	bdBlock += calculateChecksumOfBlock(bdBlock, false)

	expBlock := transformDate(vm.Expiration)
	expBlock += calculateChecksumOfBlock(expBlock, false)

	// overall checksum
	checksum := calculateChecksumOfBlock(cardNumberBlock+bdBlock+expBlock, vm.Manipulation)

	result := make([]string, 4)
	result[0] = cardNumberBlock + "D"
	result[1] = bdBlock
	result[2] = expBlock
	result[3] = checksum
	vm.Parts = result
}

func (vm *ViewModel) CalculatePassport() {
	cardNumberBlock := "C01X00T478"

	bdBlock := transformDate(vm.Birthday)
	bdBlock += calculateChecksumOfBlock(bdBlock, false)

	expBlock := transformDate(vm.Expiration)
	expBlock += calculateChecksumOfBlock(expBlock, false)

	// overall checksum
	checksum := calculateChecksumOfBlock(cardNumberBlock+bdBlock+expBlock, vm.Manipulation)

	// build result in single line
	result := make([]string, 1)
	result[0] = cardNumberBlock + "D<<" + bdBlock + "F" + expBlock + "<<<<<<<<<<<<<<<" + checksum
	vm.Parts = result
}

func transformDate(input string) string {
	inputLayout := "2006-01-02"
	outputLayout := "060102"

	t, _ := time.Parse(inputLayout, input)

	return t.Format(outputLayout)
}

func calculateChecksumOfBlock(block string, manipulate bool) string {
	weights := NewWeightsGenerator()
	sum := 0

	for _, char := range block {
		no := transformToNumber(char)
		w := weights.next()
		sum += (no * w) % 10
	}

	if manipulate {
		sum++
	}

	return strconv.Itoa(sum % 10)
}

func transformToNumber(character rune) int {
	if '0' <= character && character <= '9' {
		return int(character - '0')
	}

	if 'A' <= character && character <= 'Z' {
		return int(character-'A') + 10
	}

	return -999
}

type WeightsGenerator struct {
	state int
}

func NewWeightsGenerator() WeightsGenerator {
	return WeightsGenerator{state: 1}
}

func (g *WeightsGenerator) next() int {
	switch g.state {
	case 7:
		g.state = 3
	case 3:
		g.state = 1
	case 1:
		g.state = 7
	}

	return g.state
}

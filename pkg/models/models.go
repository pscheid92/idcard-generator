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
	Birthday     string
	Expiration   string
	Manipulation bool
	CardOptions  [3]CardOption
	Parts        []string
}

func NewViewModel() ViewModel {
	model := ViewModel{}

	// prepare birthday and expiration
	model.Birthday = "1980-01-01"
	model.Expiration = "2030-12-31"
	model.Manipulation = false
	model.Parts = nil

	// prepare default card selection
	model.CardOptions[0] = CardOption{"Neuer Personalausweis", "newid", false}
	model.CardOptions[1] = CardOption{"Alter Personalausweis", "oldid", false}
	model.CardOptions[2] = CardOption{"EU-Reisepass", "passport", false}

	return model
}

func (vm *ViewModel) CalculateNewId() {
	cardNumberBlock := "T220001293"

	bdBlock := transformDate(vm.Birthday)
	bdBlock = bdBlock + calculateChecksumOfBlock(bdBlock, false)

	expBlock := transformDate(vm.Expiration)
	expBlock = expBlock + calculateChecksumOfBlock(expBlock, false)

	// overall checksum
	checksum := calculateChecksumOfBlock(cardNumberBlock+bdBlock+expBlock, vm.Manipulation)

	result := make([]string, 4)
	result[0] = cardNumberBlock
	result[1] = bdBlock
	result[2] = expBlock + "D"
	result[3] = checksum
	vm.Parts = result
}

func (vm *ViewModel) CalculateOldId() {
	cardNumberBlock := "1220001297"

	bdBlock := transformDate(vm.Birthday)
	bdBlock = bdBlock + calculateChecksumOfBlock(bdBlock, false)

	expBlock := transformDate(vm.Expiration)
	expBlock = expBlock + calculateChecksumOfBlock(expBlock, false)

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
	bdBlock = bdBlock + calculateChecksumOfBlock(bdBlock, false)

	expBlock := transformDate(vm.Expiration)
	expBlock = expBlock + calculateChecksumOfBlock(expBlock, false)

	// overall checksum
	checksum := calculateChecksumOfBlock(cardNumberBlock+bdBlock+expBlock, vm.Manipulation)

	result := make([]string, 3)
	result[0] = cardNumberBlock + "D"
	result[1] = bdBlock + "F" + expBlock
	result[2] = checksum
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

	var sum = 0
	for _, char := range block {
		no := transformToNumber(char)
		w := weights.next()
		sum += (no * w) % 10
	}

	if manipulate {
		sum += 1
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

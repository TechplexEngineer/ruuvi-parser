package parser

import "fmt"

type RuuviMessageTooShortError struct {
	ExpectedLength int
	ActualLength   int
}

func (e *RuuviMessageTooShortError) Error() string {
	return fmt.Sprintf("data is too short expected: %d got %d", e.ExpectedLength, e.ActualLength)
}

type RuuviMessageNotManufacturerSpecificDataError struct {
	Expected byte
	Actual   byte
}

func (e *RuuviMessageNotManufacturerSpecificDataError) Error() string {
	return fmt.Sprintf("data is not manufacturer specific data. expected: %x got %x", e.Expected, e.Actual)
}

type RuuviWrongCompanyIdentError struct {
	Expected []byte
	Actual   []byte
}

func (e *RuuviWrongCompanyIdentError) Error() string {
	return fmt.Sprintf("data has wrong company identifier. expected: %x got %x", e.Expected, e.Actual)
}

type RuuviUnsupportedFormatError struct {
	Expected byte
	Actual   byte
}

func (e *RuuviUnsupportedFormatError) Error() string {
	return fmt.Sprintf("data format is %d expected data format %d", e.Actual, e.Expected)
}

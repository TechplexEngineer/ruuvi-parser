package parser

import (
	"encoding/binary"
	"encoding/hex"
)

const RuuviV3MinMessageLength = 21
const DataFormat3 byte = 3

func ParseFormat3(inputHex string) (Measurement, error) {
	var m Measurement
	data, err := hex.DecodeString(inputHex)
	if err != nil {
		return m, err
	}
	if len(data) < RuuviV3MinMessageLength {
		return m, &RuuviMessageTooShortError{ExpectedLength: RuuviV3MinMessageLength, ActualLength: len(data)}
	}

	if data[4] != BluetoothMfgSpecificBoundary {
		return m, &RuuviMessageNotManufacturerSpecificDataError{Expected: BluetoothMfgSpecificBoundary, Actual: data[4]}
	}

	if data[5] != RuuviCompanyIdentifier()[0] || data[6] != RuuviCompanyIdentifier()[1] {
		return m, &RuuviWrongCompanyIdentError{Expected: RuuviCompanyIdentifier(), Actual: data[5:6]}
	}

	// strip off prefix
	data = data[7:]

	if data[0] != DataFormat3 { // data format
		return m, &RuuviUnsupportedFormatError{Expected: DataFormat3, Actual: data[0]}
	}

	m.DataFormat = int64(data[0])
	m.Humidity = f64(float64(data[1]) / 2)
	temperatureSign := (data[2] >> 7) & 1
	temperatureBase := data[2] & 0x7F
	temperatureFraction := float64(int8(data[3])) / 100
	temperature := float64(temperatureBase) + temperatureFraction
	if temperatureSign == 1 {
		temperature *= -1
	}
	m.Temperature = f64(temperature)
	m.Pressure = f64(float64(binary.BigEndian.Uint16(data[4:])) + 50_000)
	m.AccelerationX = f64(float64(int16(binary.BigEndian.Uint16(data[6:]))) / 1000)
	m.AccelerationY = f64(float64(int16(binary.BigEndian.Uint16(data[8:]))) / 1000)
	m.AccelerationZ = f64(float64(int16(binary.BigEndian.Uint16(data[10:]))) / 1000)
	m.BatteryVoltage = f64(float64(binary.BigEndian.Uint16(data[12:])) / 1000)

	return m, nil
}

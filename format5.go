package parser

import (
	"encoding/binary"
	"encoding/hex"
)

const RuuviV5MinMessageLength = 31
const DataFormat5 byte = 5

func ParseFormat5(inputHex string) (Measurement, error) {
	var m Measurement
	data, err := hex.DecodeString(inputHex)
	if err != nil {
		return m, err
	}
	if len(data) < RuuviV5MinMessageLength {
		return m, &RuuviMessageTooShortError{ExpectedLength: RuuviV5MinMessageLength, ActualLength: len(data)}
	}

	if data[4] != BluetoothMfgSpecificBoundary { // manufacturer specific data
		return m, &RuuviMessageNotManufacturerSpecificDataError{Expected: BluetoothMfgSpecificBoundary, Actual: data[4]}
	}

	if data[5] != RuuviCompanyIdentifier()[0] || data[6] != RuuviCompanyIdentifier()[1] {
		return m, &RuuviWrongCompanyIdentError{Expected: RuuviCompanyIdentifier(), Actual: data[5:6]}
	}

	// strip off prefix
	data = data[7:]

	if data[0] != DataFormat5 {
		return m, &RuuviUnsupportedFormatError{Expected: DataFormat5, Actual: data[0]}
	}

	m.DataFormat = int64(data[0])
	m.Temperature = f64(float64(int16(binary.BigEndian.Uint16(data[1:]))) / 200)
	m.Humidity = f64(float64(binary.BigEndian.Uint16(data[3:])) / 400)
	m.Pressure = f64(float64(binary.BigEndian.Uint16(data[5:])) + 50_000)
	m.AccelerationX = f64(float64(int16(binary.BigEndian.Uint16(data[7:]))) / 1000)
	m.AccelerationY = f64(float64(int16(binary.BigEndian.Uint16(data[9:]))) / 1000)
	m.AccelerationZ = f64(float64(int16(binary.BigEndian.Uint16(data[11:]))) / 1000)

	powerInfo := binary.BigEndian.Uint16(data[13:])

	m.BatteryVoltage = f64(float64(powerInfo>>5)/1000.0 + 1.6)
	m.TxPower = i64(int64(powerInfo&0b11111)*2 - 40)

	m.MovementCounter = i64(int64(data[15]))
	m.MeasurementSequenceNumber = i64(int64(binary.BigEndian.Uint16(data[16:])))

	return m, nil
}

package parser

import (
	"encoding/hex"
	"encoding/json"
	"log"
	"reflect"
	"testing"
)

const Preamble = "02010011FF9904"

func TestParseFormat3(t *testing.T) {
	type args struct {
		inputHex string
	}
	tests := []struct {
		name    string
		args    args
		want    Measurement
		wantErr error
	}{
		{
			name: "valid data",
			args: args{inputHex: Preamble + "03291A1ECE1EFC18F94202CA0B53"},
			want: Measurement{
				DataFormat:     int64(DataFormat3),
				Temperature:    f64(26.3),
				Pressure:       f64(102766),
				Humidity:       f64(20.5),
				AccelerationX:  f64(-1),
				AccelerationY:  f64(-1.726),
				AccelerationZ:  f64(0.714),
				BatteryVoltage: f64(2.899),
			},
			wantErr: nil,
		},
		{
			name: "maximum values",
			args: args{inputHex: Preamble + "03FF7F63FFFF7FFF7FFF7FFFFFFF"},
			want: Measurement{
				DataFormat:     int64(DataFormat3),
				Temperature:    f64(127.99),
				Pressure:       f64(115535),
				Humidity:       f64(127.5),
				AccelerationX:  f64(32.767),
				AccelerationY:  f64(32.767),
				AccelerationZ:  f64(32.767),
				BatteryVoltage: f64(65.535),
			},
			wantErr: nil,
		},
		{
			name: "minimum values",
			args: args{inputHex: Preamble + "0300FF6300008001800180010000"},
			want: Measurement{
				DataFormat:     int64(DataFormat3),
				Temperature:    f64(-127.99),
				Pressure:       f64(50000),
				Humidity:       f64(0),
				AccelerationX:  f64(-32.767),
				AccelerationY:  f64(-32.767),
				AccelerationZ:  f64(-32.767),
				BatteryVoltage: f64(0),
			},
			wantErr: nil,
		},
		{
			name:    "invalid format",
			args:    args{inputHex: Preamble + "0900FF6300008001800180010000"},
			want:    Measurement{},
			wantErr: &RuuviUnsupportedFormatError{},
		},
		{
			name:    "wrong company",
			args:    args{inputHex: "02010011FF9F04" + "0300FF6300008001800180010000"},
			want:    Measurement{},
			wantErr: &RuuviWrongCompanyIdentError{},
		},
		{
			name:    "Missing MFG Specific Data",
			args:    args{inputHex: "02010011F09004" + "0300FF6300008001800180010000"},
			want:    Measurement{},
			wantErr: &RuuviMessageNotManufacturerSpecificDataError{},
		},
		{
			name:    "message too short",
			args:    args{inputHex: Preamble + "0900FF"},
			want:    Measurement{},
			wantErr: &RuuviMessageTooShortError{},
		},
		{
			name:    "invalid hex data",
			args:    args{inputHex: Preamble + "("},
			want:    Measurement{},
			wantErr: hex.InvalidByteError(0x00),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseFormat3(tt.args.inputHex)

			if tt.wantErr == nil && err == nil {
				// expected no error, and got no error
			} else if (tt.wantErr == nil && err != nil) || (tt.wantErr != nil && err == nil) {
				// expected does not match got
				t.Errorf("ParseFormat3() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				// expected error
				if reflect.TypeOf(err) != reflect.TypeOf(tt.wantErr) {
					t.Errorf("ParseFormat3() error = '%v', wantErr '%v'", err, tt.wantErr)
					return
				}
			}
			if !reflect.DeepEqual(got, tt.want) {
				dat, _ := json.MarshalIndent(got, "", "    ")
				log.Printf("Got: %s", dat)

				dat2, _ := json.MarshalIndent(tt.want, "", "    ")
				log.Printf("Want:%s", dat2)
				t.Errorf("ParseFormat3() got = %v, want %v", got, tt.want)
			}
		})
	}
}

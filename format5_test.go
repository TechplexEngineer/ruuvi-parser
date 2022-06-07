package parser

import (
	"encoding/hex"
	"encoding/json"
	"log"
	"reflect"
	"testing"
)

func TestParseFormat5(t *testing.T) {
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
			args: args{inputHex: Preamble + "0512FC5394C37C0004FFFC040CAC364200CDCBB8334C884F"},
			want: Measurement{
				DataFormat:                int64(DataFormat5),
				Temperature:               f64(24.3),
				Pressure:                  f64(100044),
				Humidity:                  f64(53.49),
				AccelerationX:             f64(0.004),
				AccelerationY:             f64(-0.004),
				AccelerationZ:             f64(1.036),
				TxPower:                   i64(4),
				BatteryVoltage:            f64(2.977),
				MovementCounter:           i64(66),
				MeasurementSequenceNumber: i64(205),
				//Mac:                       "CB:B8:33:4C:88:4F",
			},
			wantErr: nil,
		},
		{
			name: "maximum values",
			args: args{inputHex: Preamble + "057FFFFFFEFFFE7FFF7FFF7FFFFFDEFEFFFECBB8334C884F"},
			want: Measurement{
				DataFormat:                int64(DataFormat5),
				Temperature:               f64(163.835),
				Pressure:                  f64(115534),
				Humidity:                  f64(163.8350),
				AccelerationX:             f64(32.767),
				AccelerationY:             f64(32.767),
				AccelerationZ:             f64(32.767),
				TxPower:                   i64(20),
				BatteryVoltage:            f64(3.646),
				MovementCounter:           i64(254),
				MeasurementSequenceNumber: i64(65534),
				//Mac:                       "CB:B8:33:4C:88:4F",
			},
			wantErr: nil,
		},
		{
			name: "minimum values",
			args: args{inputHex: Preamble + "058001000000008001800180010000000000CBB8334C884F"},
			want: Measurement{
				DataFormat:                int64(DataFormat5),
				Temperature:               f64(-163.835),
				Pressure:                  f64(50000),
				Humidity:                  f64(0.000),
				AccelerationX:             f64(-32.767),
				AccelerationY:             f64(-32.767),
				AccelerationZ:             f64(-32.767),
				TxPower:                   i64(-40),
				BatteryVoltage:            f64(1.6),
				MovementCounter:           i64(0),
				MeasurementSequenceNumber: i64(0),
				//Mac:                       "CB:B8:33:4C:88:4F",
			},
			wantErr: nil,
		},
		{
			name:    "invalid format",
			args:    args{inputHex: Preamble + "098001000000008001800180010000000000CBB8334C884F"},
			want:    Measurement{},
			wantErr: &RuuviUnsupportedFormatError{},
		},
		{
			name:    "wrong company",
			args:    args{inputHex: "02010011FF9F04" + "058001000000008001800180010000000000CBB8334C884F"},
			want:    Measurement{},
			wantErr: &RuuviWrongCompanyIdentError{},
		},
		{
			name:    "Missing MFG Specific Data",
			args:    args{inputHex: "02010011F09004" + "058001000000008001800180010000000000CBB8334C884F"},
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
			got, err := ParseFormat5(tt.args.inputHex)

			if tt.wantErr == nil && err == nil {
				// expected no error, and got no error
			} else if (tt.wantErr == nil && err != nil) || (tt.wantErr != nil && err == nil) {
				// expected does not match got
				t.Errorf("ParseFormat5() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				// expected error
				if reflect.TypeOf(err) != reflect.TypeOf(tt.wantErr) {
					t.Errorf("ParseFormat5() error = '%v', wantErr '%v'", err, tt.wantErr)
					return
				}
			}
			if !reflect.DeepEqual(got, tt.want) {
				dat, _ := json.MarshalIndent(got, "", "    ")
				log.Printf("Got: %s", dat)

				dat2, _ := json.MarshalIndent(tt.want, "", "    ")
				log.Printf("Want:%s", dat2)
				t.Errorf("ParseFormat5() got = %#v, want %#v", got, tt.want)
			}
		})
	}
}

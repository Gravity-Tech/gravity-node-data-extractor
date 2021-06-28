package bridge

import (
	"context"
	"testing"
	"time"

	"github.com/Gravity-Tech/gateway/abi/ethereum/ibport"
	"github.com/Gravity-Tech/gravity-node-data-extractor/v2/extractors"
	"github.com/Gravity-Tech/gravity-node-data-extractor/v2/helpers"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/wavesplatform/gowaves/pkg/client"
)

//var bridge := &WavesToEthereumExtractionBridge{}

func TestWavesToEthereumExtractionBridge_Configure(t *testing.T) {
	type fields struct {
		config         ConfigureCommand
		configured     bool
		ethClient      *ethclient.Client
		wavesClient    *client.Client
		wavesHelper    helpers.ClientHelper
		ibPortContract *ibport.IBPort
	}
	type args struct {
		config ConfigureCommand
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{name: "empty args", fields: fields{}, args: args{}, wantErr: true},
		{name: "empty args with destination node", fields: fields{}, args: args{config: ConfigureCommand{DestinationNodeUrl: "https://api.avax-test.network/ext/bc/C/rpc"}}, wantErr: true},
		{name: "already configured", fields: fields{configured: true}, args: args{}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := &WavesToEthereumExtractionBridge{
				config:     tt.fields.config,
				configured: tt.fields.configured,
				//cache:          tt.fields.cache,
				ethClient:      tt.fields.ethClient,
				wavesClient:    tt.fields.wavesClient,
				wavesHelper:    tt.fields.wavesHelper,
				ibPortContract: tt.fields.ibPortContract,
			}
			if err := provider.Configure(tt.args.config); (err != nil) != tt.wantErr {
				t.Errorf("WavesToEthereumExtractionBridge.Configure() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWavesToEthereumExtractionBridge_ExtractReverseTransferRequest(t *testing.T) {
	type fields struct {
		config         ConfigureCommand
		configured     bool
		cache          map[RequestId]time.Time
		ethClient      *ethclient.Client
		wavesClient    *client.Client
		wavesHelper    helpers.ClientHelper
		ibPortContract *ibport.IBPort
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *extractors.Data
		wantErr bool
	}{
		{name: "simple", fields: fields{}, args: args{ctx: context.Background()}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := &WavesToEthereumExtractionBridge{}
			cmd := ConfigureCommand{
				SourceDecimals:      6,
				DestinationDecimals: 18,
				SourceNodeUrl:       "https://nodes.wavesexplorer.com",
				DestinationNodeUrl:  "https://bsc-dataseed1.binance.org",
				IBPortAddress:       "0x59622815BADB181a2c37052136a9480C6A4a4eA6",
				LUPortAddress:       "3PPUsj1yjMMAAg2hihdebK7n8zkAagHqdNT",
				//SourceChainAssetID: "6nSpVyNH7yM69eg446wrQR94ipbbcmZMU1ENPwanC97g",
				//DestinationChainAssetID: "0x496d451dDAB0F79346f773CbC2eb7Aee58446019"
			}

			provider.Configure(cmd)
			_, err := provider.ExtractReverseTransferRequest(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("WavesToEthereumExtractionBridge.ExtractReverseTransferRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("WavesToEthereumExtractionBridge.ExtractReverseTransferRequest() = %v, want %v", got, tt.want)
			// }
		})
	}
}

func TestWavesToEthereumExtractionBridge_ExtractDirectTransferRequest(t *testing.T) {
	type fields struct {
		config         ConfigureCommand
		configured     bool
		cache          map[RequestId]time.Time
		ethClient      *ethclient.Client
		wavesClient    *client.Client
		wavesHelper    helpers.ClientHelper
		ibPortContract *ibport.IBPort
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *extractors.Data
		wantErr bool
	}{
		{name: "simple", fields: fields{}, args: args{ctx: context.Background()}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := &WavesToEthereumExtractionBridge{}
			cmd := ConfigureCommand{
				SourceDecimals:      6,
				DestinationDecimals: 18,
				SourceNodeUrl:       "https://nodes.wavesexplorer.com",
				DestinationNodeUrl:  "https://bsc-dataseed1.binance.org",
				IBPortAddress:       "0x59622815BADB181a2c37052136a9480C6A4a4eA6",
				LUPortAddress:       "3PPUsj1yjMMAAg2hihdebK7n8zkAagHqdNT",
				//SourceChainAssetID: "6nSpVyNH7yM69eg446wrQR94ipbbcmZMU1ENPwanC97g",
				//DestinationChainAssetID: "0x496d451dDAB0F79346f773CbC2eb7Aee58446019"
			}

			provider.Configure(cmd)
			_, err := provider.ExtractDirectTransferRequest(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("WavesToEthereumExtractionBridge.ExtractDirectTransferRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("WavesToEthereumExtractionBridge.ExtractDirectTransferRequest() = %v, want %v", got, tt.want)
			// }
		})
	}
}

package settlement

import "errors"

// Config for the DymensionLayerClient
type Config struct {
	KeyringBackend string `mapstructure:"keyring_backend"`
	NodeAddress    string `mapstructure:"node_address"`
	KeyringHomeDir string `mapstructure:"keyring_home_dir"`
	DymAccountName string `mapstructure:"dym_account_name"`
	RollappID      string `mapstructure:"rollapp_id"`
	GasLimit       uint64 `mapstructure:"gas_limit"`
	GasPrices      string `mapstructure:"gas_prices"`
	GasFees        string `mapstructure:"gas_fees"`
	AddressPrefix  string `mapstructure:"address_prefix"`
	Contract       string `mapstructure:"contract"`

	//For testing only. probably should be refactored
	ProposerPubKey string `json:"proposer_pub_key"`
}

func (c Config) Validate() error {
	if c.GasPrices != "" && c.GasFees != "" {
		return errors.New("cannot provide both fees and gas prices")
	}

	if c.GasPrices == "" && c.GasFees == "" {
		return errors.New("must provide either fees or gas prices")
	}

	if c.RollappID == "" {
		return errors.New("must provide rollapp id")
	}

	return nil
}

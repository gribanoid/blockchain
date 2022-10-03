package config

type Config struct {
	BitcoinNode
}

type BitcoinNode struct {
	NodeAddress string `envconfig:"NODE_ADDRESS" default:""`
	RPCUser     string `envconfig:"RPC_USER" default:""`
	RPCPassword string `envconfig:"RPC_PASSWORD" default:""`
}

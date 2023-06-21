package cmd

import "strings"

type bcosSolidityCmdV060 struct{}

func (c *bcosSolidityCmdV060) Cmd(args Args) (string, []string, error) {
	if err := args.Valid(); err != nil {
		return "", nil, err
	}
	cmd := []string{"solc", "--bin", "--abi", "--overwrite", args.TargetDir + "/Setter.sol", "-o", args.TargetDir}
	return strings.Join(cmd, " "), cmd, nil
}

func (c *bcosSolidityCmdV060) Organization() string {
	return "fisco-bcos"
}

func (c *bcosSolidityCmdV060) Images() string {
	return "ethereum/solc:0.6.0"
}

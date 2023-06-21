package cmd

import "strings"

type chainmakerSolidityCmdV060 struct{}

func (c *chainmakerSolidityCmdV060) Cmd(args Args) (string, []string, error) {
	if err := args.Valid(); err != nil {
		return "", nil, err
	}
	cmd := []string{"solc", "--bin", "--abi", "--overwrite", args.TargetDir + "/Setter.sol", "-o", args.TargetDir}
	return strings.Join(cmd, " "), cmd, nil
}

func (c *chainmakerSolidityCmdV060) Organization() string {
	return "chainmaker"
}

func (c *chainmakerSolidityCmdV060) Images() string {
	return "ethereum/solc:0.6.0"
}

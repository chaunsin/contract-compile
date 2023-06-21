package cmd

import "strings"

type xuperSolidityCmdV060 struct{}

func (c *xuperSolidityCmdV060) Cmd(args Args) (string, []string, error) {
	if err := args.Valid(); err != nil {
		return "", nil, err
	}
	cmd := []string{"solc", "--bin", "--abi", "--overwrite", args.TargetDir + "/Setter.sol", "-o", args.TargetDir}
	return strings.Join(cmd, " "), cmd, nil
}

func (c *xuperSolidityCmdV060) Organization() string {
	return "xuperchain"
}

func (c *xuperSolidityCmdV060) Images() string {
	return "ethereum/solc:0.6.0"
}

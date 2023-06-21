package cmd

import "strings"

type ethSolidityCmdV060 struct{}

func (c *ethSolidityCmdV060) Cmd(args Args) (string, []string, error) {
	if err := args.Valid(); err != nil {
		return "", nil, err
	}
	// todo: 考虑编译多个合约文件
	//cmd := []string{"solc", "--bin", "--abi", "--overwrite", "/data/Setter.sol", "-o", "/data"}
	cmd := []string{"solc", "--bin", "--abi", "--overwrite", args.TargetDir + "/Setter.sol", "-o", args.TargetDir}
	return strings.Join(cmd, " "), cmd, nil
}

func (c *ethSolidityCmdV060) Organization() string {
	return "eth"
}

func (c *ethSolidityCmdV060) Images() string {
	return "ethereum/solc:0.6.0"
}

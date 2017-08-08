package main

import "fmt"

var upgradeTemplate = `%s
Total Price: %d
Gold: %d
Wood: %d
Stone: %d
Price per Warrior: %d
`

func renderUpgradeName(name string, upgrade Upgrade) string {
	str := ""
	if upgrade.UpgradeStorage {
		str = str + "Storage + "
	}
	if upgrade.UpgradeHouses {
		str = str + "Houses + "
	}
	return str + name
}

func renderUpgrade(name string, upgrade Upgrade) string {
	str := renderUpgradeName(name, upgrade)
	str = fmt.Sprintf(upgradeTemplate, str, upgrade.Price, upgrade.Gold, upgrade.Wood, upgrade.Stone, upgrade.PricePerWarrior)
	return str
}

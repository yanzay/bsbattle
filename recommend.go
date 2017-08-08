package main

import (
	"fmt"
	"log"
)

type Upgrade struct {
	PricePerWarrior int
	UpgradeStorage  bool
	UpgradeHouses   bool
	Price           int
	Gold            int
	Wood            int
	Stone           int
}

const (
	Storage int = iota
	Barracks
	Wall
	Trebuchet
	Houses
)

type Coef struct {
	Gold  int
	Wood  int
	Stone int
}

var coefs = map[int]Coef{
	Barracks:  Coef{200, 100, 100},
	Storage:   Coef{200, 100, 100},
	Wall:      Coef{5000, 500, 1500},
	Trebuchet: Coef{8000, 1000, 300},
	Houses:    Coef{200, 100, 100},
}

func recommend(builds *Buildings) string {
	if builds.Houses < builds.Barracks {
		return fmt.Sprintf("Upgrade Houses up to %d level, then come back for better advice.", builds.Barracks)
	}
	var upgradeBarracks Upgrade
	if builds.Houses == builds.Barracks {
		upgradeBarracks = addUpgrades(calcHouses(builds), calcBarracks(builds))
		upgradeBarracks.UpgradeHouses = true
		upgradeBarracks.PricePerWarrior = upgradeBarracks.Price / 40
	} else {
		upgradeBarracks = calcBarracks(builds)
	}
	upgradeWall := calcWall(builds)
	upgradeTreb := calcTreb(builds)
	upgrades := fmt.Sprintf("%s\n%s\n%s",
		renderUpgrade("Barracks", upgradeBarracks),
		renderUpgrade("Wall", upgradeWall),
		renderUpgrade("Trebuchet", upgradeTreb))
	var recommendation string
	if upgradeBarracks.PricePerWarrior < upgradeTreb.PricePerWarrior && upgradeBarracks.PricePerWarrior < upgradeWall.PricePerWarrior {
		recommendation = renderUpgradeName("Barracks", upgradeBarracks)
	}
	if upgradeWall.PricePerWarrior < upgradeBarracks.PricePerWarrior && upgradeWall.PricePerWarrior < upgradeTreb.PricePerWarrior {
		recommendation = renderUpgradeName("Wall", upgradeWall)
	}
	if upgradeTreb.PricePerWarrior < upgradeBarracks.PricePerWarrior && upgradeTreb.PricePerWarrior < upgradeWall.PricePerWarrior {
		recommendation = renderUpgradeName("Trebuchet", upgradeTreb)
	}
	return upgrades + "\nRecommended Upgrade: " + recommendation
}

func calcBarracks(builds *Buildings) Upgrade {
	update := calcUpdatePrice(builds.Barracks, coefs[Barracks])
	if StorageCap(builds) < update.Wood {
		update.Price += calcStorage(builds)
		update.UpgradeStorage = true
	}
	update.PricePerWarrior = update.Price / 40
	return update
}

func calcHouses(builds *Buildings) Upgrade {
	return calcUpdatePrice(builds.Houses, coefs[Houses])
}

func calcWall(builds *Buildings) Upgrade {
	update1 := calcUpdatePrice(builds.Wall, coefs[Wall])
	update2 := calcUpdatePrice(builds.Wall+1, coefs[Wall])
	if StorageCap(builds) < update2.Stone {
		log.Println("StorageCap", builds, StorageCap(builds), update2.Stone)
		update2.Price += calcStorage(builds)
		update2.UpgradeStorage = true
	}
	update := addUpgrades(update1, update2)
	update.PricePerWarrior = update.Price / 108
	return update
}

func addUpgrades(update1, update2 Upgrade) Upgrade {
	return Upgrade{
		Price:          update1.Price + update2.Price,
		UpgradeStorage: update2.UpgradeStorage,
		Gold:           update1.Gold + update2.Gold,
		Wood:           update1.Wood + update2.Wood,
		Stone:          update1.Stone + update2.Stone,
	}
}

func calcTreb(builds *Buildings) Upgrade {
	update1 := calcUpdatePrice(builds.Trebuchet, coefs[Trebuchet])
	update2 := calcUpdatePrice(builds.Trebuchet+1, coefs[Trebuchet])
	if StorageCap(builds) < update1.Wood {
		update2.Price += calcStorage(builds)
		update2.UpgradeStorage = true
	}
	update := addUpgrades(update1, update2)
	update.PricePerWarrior = update.Price / 100
	return update
}

func StorageCap(builds *Buildings) int {
	return builds.Storage * (builds.Storage*50 + 1000)
}

func calcStorage(builds *Buildings) int {
	return calcUpdatePrice(builds.Storage, coefs[Storage]).Price
}

func calcUpdatePrice(level int, coef Coef) Upgrade {
	rise := (level + 1) * (level + 2) / 2
	updateGold := coef.Gold * rise
	updateWood := coef.Wood * rise
	updateStone := coef.Stone * rise
	return Upgrade{
		Price: updateGold + updateWood*2 + updateStone*2,
		Gold:  updateGold,
		Wood:  updateWood,
		Stone: updateStone,
	}
}

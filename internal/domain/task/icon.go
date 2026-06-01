package task

import "errors"

type Icon int32

const (
	IconUnknown Icon = iota
	IconMark
	IconHome
	IconJob
	IconSupermarket
	IconCafe
	IconActivity
	IconDrive
	IconFlight
	IconStar
	IconFlag
	IconHospital
	IconOutdoor
)

var iconStrMapper = map[Icon]string{
	IconMark:        "mark",
	IconHome:        "home",
	IconJob:         "job",
	IconSupermarket: "supermarket",
	IconCafe:        "cafe",
	IconActivity:    "activity",
	IconDrive:       "drive",
	IconFlight:      "flight",
	IconStar:        "star",
	IconFlag:        "flag",
	IconHospital:    "hospital",
	IconOutdoor:     "outdoor",
}

var stringIconMapper = map[string]Icon{
	"mark":        IconMark,
	"home":        IconHome,
	"job":         IconJob,
	"supermarket": IconSupermarket,
	"cafe":        IconCafe,
	"activity":    IconActivity,
	"drive":       IconDrive,
	"flight":      IconFlight,
	"star":        IconStar,
	"flag":        IconFlag,
	"hospital":    IconHospital,
	"outdoor":     IconOutdoor,
}

func (i Icon) String() string {
	if value, ok := iconStrMapper[i]; ok {
		return value
	}
	return "unknown"
}

func (i Icon) Valid() bool {
	return IconUnknown < i && i <= IconOutdoor
}

func ParseIcon(s string) (Icon, error) {
	if value, ok := stringIconMapper[s]; ok {
		return value, nil
	}
	return IconUnknown, errors.New("unknown icon")
}

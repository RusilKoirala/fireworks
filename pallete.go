package main

import "math/rand"

type palette []string

var (
	red    = palette{"#fff5f5", "#ff8787", "#f03e3e", "#a61e1e", "#4a0e0e"}
	gold   = palette{"#fffbe6", "#ffe066", "#fcc419", "#e8590c", "#5c2604"}
	green  = palette{"#f4fff0", "#8ce99a", "#40c057", "#2b8a3e", "#0f3318"}
	blue   = palette{"#f0f8ff", "#74c0fc", "#339af0", "#1864ab", "#0a2340"}
	purple = palette{"#fbf0ff", "#e599f7", "#be4bdb", "#862e9c", "#33103d"}
	cyan   = palette{"#effffe", "#66d9e8", "#22b8cf", "#0b7285", "#052d33"}
	pink   = palette{"#fff0f6", "#faa2c1", "#f06595", "#c2255c", "#4d0d23"}
	silver = palette{"#ffffff", "#e9ecef", "#adb5bd", "#6c757d", "#2b3035"}

	allPalletes  = []palette{red, gold, green, blue, purple, cyan, pink, silver}
	emberPalette = palette{"#ffe8a3", "#f59f00", "#d9480f", "#7a2e05", "#3a1503"}
	heartPals    = []palette{red, pink}

	tailGray = []string{"#ffffff", "#dee2e6", "#adb5bd", "#868e96", "#5c636a", "#3d434a", "#2b3035"}
)

func randPal() palette {
	return allPalletes[rand.Intn(len((allPalletes)))]
}

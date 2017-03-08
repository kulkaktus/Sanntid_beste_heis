package channels

func Sensor(floor int) int {
	return [...]int{0x204, 0x205, 0x206, 0x207}[floor-1]
}

func Up_button(floor int) int {
	return [...]int{0x311, 0x310, 0x201, 0}[floor-1]
}

func Down_button(floor int) int {
	return [...]int{0, 0x200, 0x202, 0x203}[floor-1]
}

func Internal_button(floor int) int {
	return [...]int{0x315, 0x314, 0x313, 0x312}[floor-1]
}

func Internal_light(floor int) int {
	return [...]int{0x30D, 0x30C, 0x30B, 0x30A}[floor-1]
}

func Up_light(floor int) int {
	return [...]int{0x309, 0x308, 0x306, 0}[floor-1]
}

func Down_light(floor int) int {
	return [...]int{0, 0x307, 0x305, 0x304}[floor-1]
}

const (
	Floor_light_0 = 0x300
	Floor_light_1 = 0x301
	Door          = 0x303
	Stop_light    = 0x30E
	Motor_value   = 0x100
	Motor_dir     = 0x30F
)

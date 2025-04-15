package go_telemetry_course

type GoatStudent struct {
	Name    string
	Age     int
	Classes []string
}

type GoatClass struct {
	Name      string
	StartTime int // 1 to 8, representing school hours
}

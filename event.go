package disperse

type Event struct {
	ConsumerId string
	Kind       string
	Data       []byte
}

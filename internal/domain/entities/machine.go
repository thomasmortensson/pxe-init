package entities

type Machine struct {
	Name  string
	MAC   string
	UUID  string
	Image *Image
}

// NewMachine generates a new entity Machine structure
// NOTE, this function does not set an initial Image pointer. Please call SetImage
func NewMachine(name, mac, uuid string) *Machine {
	return &Machine{
		Name: name,
		MAC:  mac,
		UUID: uuid,
	}
}

// SetImage sets the specified image as the active image for machine
func (m *Machine) SetImage(image *Image) {
	m.Image = image
}

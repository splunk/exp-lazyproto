package lazyproto

const FlagsMessageModified = 1

type ProtoMessage struct {
	Bytes  []byte
	Flags  uint64
	Parent *ProtoMessage
}

func (m *ProtoMessage) MarkModified() {
	if m.Flags&FlagsMessageModified != 0 {
		return
	}

	m.Flags |= FlagsMessageModified
	parent := m.Parent
	for parent != nil {
		if parent.Flags&FlagsMessageModified != 0 {
			break
		}
		parent.Flags |= FlagsMessageModified
		parent = parent.Parent
	}
}

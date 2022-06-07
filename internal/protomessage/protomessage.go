package protomessage

type ProtoMessage struct {
	// Bytes are set to nil when the message is modified (i.e. marshaling will do
	// full field-by-field encoding).
	Bytes  BytesView
	Parent *ProtoMessage
}

func (m *ProtoMessage) IsModified() bool {
	return m.Bytes.Data == nil
}

func (m *ProtoMessage) MarkModified() {
	if m.Bytes.Data != nil {
		m.markModified()
	}
}

func (m *ProtoMessage) markModified() {
	m.Bytes.Data = nil
	m.Bytes.Len = 0
	parent := m.Parent
	for parent != nil {
		if parent.IsModified() {
			break
		}
		parent.Bytes.Data = nil
		parent.Bytes.Len = 0
		parent = parent.Parent
	}
}

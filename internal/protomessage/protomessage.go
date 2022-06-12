package protomessage

type ProtoMessage struct {
	// Bytes are a view into the original bytes from which this message was
	// unmarshalled.
	Bytes BytesView

	// Parent is the parent message of this message (or nill if no parent).
	Parent *ProtoMessage
}

func (m *ProtoMessage) IsModified() bool {
	// Bytes are set to nil when the message is modified (i.e. marshaling will do
	// full field-by-field encoding) or if the message was created a new (was
	// not unmarshalled).
	return m.Bytes.Data == nil
}

// MarkModified marks this message as modified, so that future marshalling operations
// will no longer try to re-use the original bytes.
func (m *ProtoMessage) MarkModified() {
	if m.Bytes.Data != nil {
		m.markModified()
	}
}

func (m *ProtoMessage) markModified() {
	// Reset the pointer to original bytes.
	m.Bytes.Data = nil
	m.Bytes.Len = 0

	// Traverse up the parents chain and mark all parents as modified too.
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

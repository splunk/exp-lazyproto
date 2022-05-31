package simple

type ScopeLogsPool struct {
	freedScopeLogs []ScopeLogs
}

var scopeLogsPool = ScopeLogsPool{}

func (p *ScopeLogsPool) GetScopeLogss(count int) []ScopeLogs {
	if len(p.freedScopeLogs) >= count {
		r := p.freedScopeLogs[len(p.freedScopeLogs)-count:]
		p.freedScopeLogs = p.freedScopeLogs[:len(p.freedScopeLogs)-count]
		return r
	}

	r := make([]ScopeLogs, count)
	i := 0
	for ; i < len(p.freedScopeLogs); i++ {
		r[i] = p.freedScopeLogs[i]
	}
	p.freedScopeLogs = nil

	return r
}

func (p *ScopeLogsPool) Release(l *LogsData) {
	for _, rl := range l.resourceLogs {
		p.freedScopeLogs = append(p.freedScopeLogs, rl.scopeLogs...)
	}
}

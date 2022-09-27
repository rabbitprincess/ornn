package bp

type T_BP__config__base_info struct {
	TD_n4_version TD_N4_bp_version
	// ???
}

func (t *T_BP__config__base_info) Set_version(_td_n4_bp_version TD_N4_bp_version) {
	t.TD_n4_version = _td_n4_bp_version
}

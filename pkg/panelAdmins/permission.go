package panelAdmins

type ReadWrite struct {
	Write bool
	Read  bool
}

type Permission struct {
	Onboarding ReadWrite
	Role       ReadWrite
	Team       ReadWrite
	Tenant     ReadWrite
	Billing    ReadWrite
}

func (p *Permission) UpdatePermission(pm Permission) error  {
	*p = pm
	return nil
}

func (p *Permission) ConvertToJson() {}

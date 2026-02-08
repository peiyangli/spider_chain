package types

type OperatorI interface {
	GetRole() uint64
	GetPermissions() uint64
}

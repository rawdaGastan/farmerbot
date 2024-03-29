// Package models for farmerbot models.
package models

// ConsumableResources for node resources
type ConsumableResources struct {
	OverProvisionCPU float64  `json:"overProvisionCPU,omitempty"` // how much we allow over provisioning the CPU range: [1;4]
	Total            Capacity `json:"total"`
	Used             Capacity `json:"used,omitempty"`
}

// Capacity is node resource capacity
type Capacity struct {
	HRU  uint64 `json:"HRU"`
	SRU  uint64 `json:"SRU"`
	CRU  uint64 `json:"CRU"`
	MRU  uint64 `json:"MRU"`
	Ipv4 uint64 `json:"ipv4"`
}

// IsEmpty checks empty capacity
func (cap *Capacity) isEmpty() bool {
	return cap.CRU == 0 && cap.MRU == 0 && cap.SRU == 0 && cap.HRU == 0
}

// Add adds a new for capacity
func (cap *Capacity) Add(add Capacity) {
	cap.CRU += add.CRU
	cap.MRU += add.MRU
	cap.SRU += add.SRU
	cap.HRU += add.HRU
}

// Subtract subtracts a new capacity
func (cap *Capacity) subtract(sub Capacity) (result Capacity) {
	result.CRU = cap.CRU - sub.CRU
	result.MRU = cap.MRU - sub.MRU
	result.SRU = cap.SRU - sub.SRU
	result.HRU = cap.HRU - sub.HRU
	return result
}

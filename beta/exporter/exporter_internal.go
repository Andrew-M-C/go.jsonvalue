package exporter

import (
	"reflect"
	"sync"
)

// g stands for Global
type g struct {
	lock sync.RWMutex

	exportersByType map[reflect.Type]omnipotentExporter

	debugf func(string, ...any)
}

var internal = &g{}

func init() {
	internal.exportersByType = make(map[reflect.Type]omnipotentExporter)
	internal.debugf = func(string, ...any) {}
}

func (i *g) loadExportersByType(typ reflect.Type) (omnipotentExporter, bool) {
	i.lock.RLock()
	defer i.lock.RUnlock()

	e, exist := i.exportersByType[typ]
	return e, exist
}

func (i *g) storeExporterByType(typ reflect.Type, e omnipotentExporter) {
	internal.lock.Lock()
	defer internal.lock.Unlock()

	internal.exportersByType[typ] = e
}
